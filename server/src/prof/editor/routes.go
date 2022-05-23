package editor

import (
	"fmt"

	"github.com/benoitkugler/maths-online/maths/exercice"
	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/prof/teacher"
	"github.com/benoitkugler/maths-online/utils"
	"github.com/labstack/echo/v4"
)

const LoopbackEndpoint = "/prof-loopback/:session_id"

// EditorStartSession setup a new editing session.
// In particular, it launches in the background a
// `loopbackController` instance to handle preview requests.
func (ct *Controller) EditorStartSession(c echo.Context) error {
	out := ct.startSession()

	return c.JSON(200, out)
}

// EditorGetTags return all tags currently used by questions.
// It also add the special difficulty tags.
func (ct *Controller) EditorGetTags(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	tags, err := SelectAllQuestionTags(ct.db)
	if err != nil {
		return err
	}

	// only return tags used by visible questions
	questions, err := SelectAllQuestions(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	// add the special difficulty tags among the proposition,
	// in first choices
	var filtred []string
	for _, tag := range tags {
		if !questions[tag.IdQuestion].IsVisibleBy(user.Id) {
			continue
		}

		switch DifficultyTag(tag.Tag) {
		case Diff1, Diff2, Diff3:
		default:
			filtred = append(filtred, tag.Tag)
		}
	}

	filtred = append([]string{
		string(Diff1), string(Diff2), string(Diff3),
	}, filtred...)

	return c.JSON(200, filtred)
}

type ListQuestionsIn struct {
	TitleQuery string // empty means all
	Tags       []string
}

func (ct *Controller) EditorSearchQuestions(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	var args ListQuestionsIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	out, err := ct.searchQuestions(args, user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) EditorCreateQuestion(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	question := Question{IdTeacher: user.Id, Public: false}
	question, err := question.Insert(ct.db)
	if err != nil {
		return err
	}

	return c.JSON(200, question)
}

func (ct *Controller) EditorDuplicateQuestion(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	out, err := ct.duplicateQuestion(id, user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) EditorDuplicateQuestionWithDifficulty(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	err = ct.duplicateQuestionWithDifficulty(id, user.Id)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) EditorDeleteQuestion(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	qu, err := SelectQuestion(ct.db, id)
	if err != nil {
		return utils.SQLError(err)
	}
	if qu.IdTeacher != user.Id {
		return accessForbidden
	}

	_, err = DeleteQuestionById(ct.db, id)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) EditorGetQuestion(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	question, err := SelectQuestion(ct.db, id)
	if err != nil {
		return err
	}

	if !question.IsVisibleBy(user.Id) {
		return accessForbidden
	}

	return c.JSON(200, question)
}

type CheckParametersIn struct {
	SessionID  string
	Parameters exercice.Parameters
}

type CheckParametersOut struct {
	ErrDefinition exercice.ErrParameters
	// Variables is the list of the variables defined
	// in the parameteres (intrinsics included)
	Variables []expression.Variable
}

func (ct *Controller) EditorCheckParameters(c echo.Context) error {
	var args CheckParametersIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	out := ct.checkParameters(args)

	return c.JSON(200, out)
}

type SaveAndPreviewIn struct {
	SessionID string
	Question  Question
}

type SaveAndPreviewOut struct {
	Error   exercice.ErrQuestionInvalid
	IsValid bool
}

// For non personnal questions, only preview.
func (ct *Controller) EditorSaveAndPreview(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	var args SaveAndPreviewIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	out, err := ct.saveAndPreview(args, user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) EditorPausePreview(c echo.Context) error {
	sessionID := c.QueryParam("sessionID")

	err := ct.pausePreview(sessionID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

// EditorEndPreview cleanly remove the loopback controller instead
// of waiting for it to timeout.
func (ct *Controller) EditorEndPreview(c echo.Context) error {
	sessionID := c.Param("sessionID")

	err := ct.endPreview(sessionID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

type UpdateTagsIn struct {
	Tags       []string
	IdQuestion int64
}

func (ct *Controller) EditorUpdateTags(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	var args UpdateTagsIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	err := ct.updateTags(args, user.Id)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

// AccessLoopback establish a connection with the embedded preview app
func (ct *Controller) AccessLoopback(c echo.Context) error {
	sessionID := c.Param("session_id")

	loopback, ok := ct.sessions[sessionID]
	if !ok {
		return fmt.Errorf("invalid session ID %s", sessionID)
	}

	// connect to the websocket handler, which handle errors
	loopback.setupWebSocket(c.Response().Writer, c.Request())

	return nil
}
