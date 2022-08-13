// Package editor provides functionnality for a frontend
// to edit and preview math questions.
package editor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/questions"
	"github.com/benoitkugler/maths-online/prof/teacher"
	"github.com/benoitkugler/maths-online/utils"
	"github.com/labstack/echo/v4"
)

const sessionTimeout = 6 * time.Hour

var accessForbidden = errors.New("access fordidden")

// Controller is the global object responsible to
// handle incoming requests regarding the editor.
type Controller struct {
	lock sync.Mutex

	db *sql.DB

	sessions map[string]*loopbackController

	admin teacher.Teacher
}

func NewController(db *sql.DB, admin teacher.Teacher) *Controller {
	return &Controller{
		db:       db,
		sessions: make(map[string]*loopbackController),
		admin:    admin,
	}
}

type StartSessionOut struct {
	ID string
}

// startSession setup a new editing session.
// In particular, it launches in the background a
// `loopbackController` instance to handle preview requests.
func (ct *Controller) startSession() StartSessionOut {
	ct.lock.Lock()
	defer ct.lock.Unlock()

	// generate a new session ID
	newID := utils.RandomID(false, 40, func(s string) bool {
		_, has := ct.sessions[s]
		return has
	})

	// create and register the loopback controller
	loopback := newLoopbackController(newID)
	ct.sessions[newID] = loopback

	// start the websocket for the loopback
	go func() {
		ctx, cancelFunc := context.WithTimeout(context.Background(), sessionTimeout)
		loopback.startLoop(ctx) // block

		cancelFunc() // cancel the timer if needed

		// remove the loopback controller when the session is over
		ct.lock.Lock()
		defer ct.lock.Unlock()
		delete(ct.sessions, newID)
	}()

	return StartSessionOut{ID: newID}
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

	out, err := ct.duplicateQuestion(IdQuestion(id), user.Id)
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

	err = ct.duplicateQuestionWithDifficulty(IdQuestion(id), user.Id)
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

	err = ct.deleteQuestion(IdQuestion(id), user.Id)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) deleteQuestion(id IdQuestion, userID uID) error {
	qu, err := SelectQuestion(ct.db, id)
	if err != nil {
		return utils.SQLError(err)
	}
	if qu.IdTeacher != userID {
		return accessForbidden
	}

	links, err := SelectExerciceQuestionsByIdQuestions(ct.db, id)
	if err != nil {
		return utils.SQLError(err)
	}

	if len(links) != 0 {
		ex, err := SelectExercice(ct.db, links[0].IdExercice)
		if err != nil {
			return utils.SQLError(err)
		}

		return fmt.Errorf("La question est utilisée dans l'exercice %s : %d.", ex.Title, ex.Id)
	}

	_, err = DeleteQuestionById(ct.db, id)
	if err != nil {
		return err
	}

	return nil
}

func (ct *Controller) EditorGetQuestion(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	question, err := SelectQuestion(ct.db, IdQuestion(id))
	if err != nil {
		return err
	}

	if !question.IsVisibleBy(user.Id) {
		return accessForbidden
	}

	return c.JSON(200, question)
}

type CheckQuestionParametersIn struct {
	SessionID  string
	Parameters questions.Parameters
}

type CheckQuestionParametersOut struct {
	ErrDefinition questions.ErrParameters
	// Variables is the list of the variables defined
	// in the parameteres (intrinsics included)
	Variables []expression.Variable
}

func (ct *Controller) EditorCheckQuestionParameters(c echo.Context) error {
	var args CheckQuestionParametersIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	out := ct.checkQuestionParameters(args)

	return c.JSON(200, out)
}

type QuestionUpdateVisiblityIn struct {
	QuestionID IdQuestion
	Public     bool
}

func (ct *Controller) QuestionUpdateVisiblity(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	// we only accept public question from admin account
	if user.Id != ct.admin.Id {
		return accessForbidden
	}

	var args QuestionUpdateVisiblityIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	qu, err := SelectQuestion(ct.db, args.QuestionID)
	if err != nil {
		return utils.SQLError(err)
	}
	if qu.IdTeacher != user.Id {
		return accessForbidden
	}

	if !args.Public {
		// check that it is not harmful to hide the question again,
		// that is exercices using this question are owned by the admin
		links, err := SelectExerciceQuestionsByIdQuestions(ct.db, qu.Id)
		if err != nil {
			return utils.SQLError(err)
		}
		exercices, err := SelectExercices(ct.db, links.IdExercices()...)
		if err != nil {
			return utils.SQLError(err)
		}
		usedExercices := make(Exercices)
		for _, link := range links {
			ex := exercices[link.IdExercice]
			if ex.IdTeacher != user.Id {
				usedExercices[ex.Id] = ex
			}
		}
		if L := len(usedExercices); L != 0 {
			return fmt.Errorf("La question est utilisée dans %d exercice(s).", L)
		}
	}
	qu.Public = args.Public
	qu, err = qu.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	return c.NoContent(200)
}

type SaveQuestionAndPreviewIn struct {
	SessionID string
	Question  Question
}

type SaveQuestionAndPreviewOut struct {
	Error   questions.ErrQuestionInvalid
	IsValid bool
}

// For non personnal questions, only preview.
func (ct *Controller) EditorSaveQuestionAndPreview(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	var args SaveQuestionAndPreviewIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	out, err := ct.saveQuestionAndPreview(args, user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

type ListQuestionsOut struct {
	Questions   []QuestionGroup // limited by `pagination`
	NbGroups    int             // total number of groups (passing the given filter)
	NbQuestions int             // total number of questions (passing the given filter)
}

// QuestionGroup groups the question forming an implicit
// group, defined by a shared title
// Standalone question are represented by a group of length one.
type QuestionGroup struct {
	Title     string
	Questions []QuestionHeader
	Size      int // the total size of the group, regardless of the current filter
}

// QuestionHeader is a sumary of the meta data of a question
type QuestionHeader struct {
	Title      string
	Tags       []string
	Id         IdQuestion
	Difficulty DifficultyTag // deduced from the tags
	IsInGroup  bool          // true if the question is in an implicit group, ignoring the current filter
	Origin     teacher.Origin
}

func normalizeTitle(title string) string {
	return removeAccents(strings.TrimSpace(strings.ToLower(title)))
}

func (qu Question) origin(userID, adminID uID) (teacher.Origin, bool) {
	vis, ok := teacher.NewVisibility(qu.IdTeacher, userID, adminID, qu.Public)
	if !ok {
		return teacher.Origin{}, false
	}
	return teacher.Origin{
		AllowPublish: userID == adminID,
		IsPublic:     qu.Public,
		Visibility:   vis,
	}, true
}

func (ct *Controller) searchQuestions(query ListQuestionsIn, userID uID) (out ListQuestionsOut, err error) {
	const pagination = 30 // number of groups

	// to find implicit groups, we need all the questions
	questions, err := SelectAllQuestions(ct.db)
	if err != nil {
		return out, utils.SQLError(err)
	}
	questions.RestrictVisible(userID)
	questions.RestrictNeedExercice()

	// the group are not modified by the title query though

	queryTitle := normalizeTitle(query.TitleQuery)
	var (
		ids      []IdQuestion
		ownerIDs []uID
		groups   = make(map[string][]IdQuestion)
	)
	for _, question := range questions {
		thisTitle := normalizeTitle(question.Page.Title)
		if strings.Contains(thisTitle, queryTitle) {
			groups[question.Page.Title] = append(groups[question.Page.Title], question.Id)
			ids = append(ids, question.Id)
			ownerIDs = append(ownerIDs, question.IdTeacher)
		}
	}

	// load the tags ...
	tags, err := SelectQuestionTagsByIdQuestions(ct.db, ids...)
	if err != nil {
		return out, utils.SQLError(err)
	}
	tagsMap := tags.ByIdQuestion()

	// normalize query
	for i, t := range query.Tags {
		query.Tags[i] = NormalizeTag(t)
	}

	// .. and build the group, restricting the questions matching the given tags
	out.Questions = make([]QuestionGroup, 0, len(groups))
	for title, ids := range groups {
		group := QuestionGroup{
			Title: title,
			Size:  len(ids),
		}

		// select the questions
		for _, id := range ids {
			crible := tagsMap[id].Crible()

			if !crible.HasAll(query.Tags) {
				continue
			}

			qu := questions[id]
			origin, _ := qu.origin(userID, ct.admin.Id)
			question := QuestionHeader{
				Id:         id,
				Title:      title,
				Difficulty: crible.Difficulty(),
				IsInGroup:  len(ids) > 1,
				Tags:       tagsMap[id].List(),
				Origin:     origin,
			}
			group.Questions = append(group.Questions, question)
		}

		// sort to make sure the display is consistent between two queries
		sort.Slice(group.Questions, func(i, j int) bool { return group.Questions[i].Id < group.Questions[j].Id })
		sort.SliceStable(group.Questions, func(i, j int) bool { return group.Questions[i].Difficulty < group.Questions[j].Difficulty })

		// ignore empty groups
		if len(group.Questions) != 0 {
			out.Questions = append(out.Questions, group)
		}

		out.NbQuestions += len(group.Questions)
	}

	// sort before pagination
	sort.Slice(out.Questions, func(i, j int) bool { return out.Questions[i].Title < out.Questions[j].Title })

	out.NbGroups = len(out.Questions)
	if len(out.Questions) > pagination {
		out.Questions = out.Questions[:pagination]
	}

	return out, nil
}

// duplicateQuestion duplicate the given question, returning
// the newly created one
func (ct *Controller) duplicateQuestion(idQuestion IdQuestion, userID uID) (Question, error) {
	qu, err := SelectQuestion(ct.db, idQuestion)
	if err != nil {
		return Question{}, utils.SQLError(err)
	}

	if !qu.IsVisibleBy(userID) {
		return Question{}, accessForbidden
	}

	tags, err := SelectQuestionTagsByIdQuestions(ct.db, qu.Id)
	if err != nil {
		return Question{}, utils.SQLError(err)
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return Question{}, utils.SQLError(err)
	}

	// shallow copy is enough; make the question private
	newQuestion := qu
	newQuestion.IdTeacher = userID
	newQuestion.Public = false
	newQuestion, err = newQuestion.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return Question{}, utils.SQLError(err)
	}

	for i := range tags {
		tags[i].IdQuestion = newQuestion.Id
	}
	err = updateTags(tx, tags, newQuestion.Id)
	if err != nil {
		_ = tx.Rollback()
		return Question{}, err
	}

	err = tx.Commit()
	if err != nil {
		return Question{}, utils.SQLError(err)
	}

	return newQuestion, nil
}

// duplicateQuestionWithDifficulty creates new questions with the same title
// and content as the given question, but with difficulty levels
// only personnal questions are allowed
func (ct *Controller) duplicateQuestionWithDifficulty(idQuestion IdQuestion, userID uID) error {
	qu, err := SelectQuestion(ct.db, idQuestion)
	if err != nil {
		return utils.SQLError(err)
	}

	if qu.IdTeacher != userID {
		return accessForbidden
	}

	tags, err := SelectQuestionTagsByIdQuestions(ct.db, qu.Id)
	if err != nil {
		return utils.SQLError(err)
	}

	// if the question already has a difficulty, respect it
	// otherwise, attribute the difficulty one
	currentDifficulty := tags.Crible().Difficulty()
	var newDifficulties [2]DifficultyTag
	switch currentDifficulty {
	case Diff1:
		newDifficulties = [2]DifficultyTag{Diff2, Diff3}
	case Diff2:
		newDifficulties = [2]DifficultyTag{Diff1, Diff3}
	case Diff3:
		newDifficulties = [2]DifficultyTag{Diff1, Diff2}
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}

	if currentDifficulty == "" {
		// update the current question
		newTags := append(tags, QuestionTag{IdQuestion: idQuestion, Tag: string(Diff1)})
		err = updateTags(tx, newTags, idQuestion)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		newDifficulties = [2]DifficultyTag{Diff2, Diff3}
	}

	for _, diff := range newDifficulties {
		newQuestion := qu // shallow copy is enough
		newQuestion, err = newQuestion.Insert(tx)
		if err != nil {
			_ = tx.Rollback()
			return utils.SQLError(err)
		}
		var newTags QuestionTags
		for _, t := range tags {
			// do not add existing difficulties
			switch DifficultyTag(t.Tag) {
			case Diff1, Diff2, Diff3:
				continue
			}

			t.IdQuestion = newQuestion.Id
			newTags = append(newTags, t)
		}
		newTags = append(newTags, QuestionTag{IdQuestion: newQuestion.Id, Tag: string(diff)})
		err = updateTags(tx, newTags, newQuestion.Id)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}

// do NOT commit or rollback
func updateTags(tx *sql.Tx, tags QuestionTags, idQuestion IdQuestion) error {
	var nbDiff, nbLevel int
	for _, tag := range tags {
		switch tag.Tag {
		case string(Diff1), string(Diff2), string(Diff3):
			nbDiff++
		case string(Seconde), string(Premiere), string(Terminale):
			nbLevel++
		}
	}
	if nbDiff > 1 {
		return errors.New("Un seul niveau de difficulté est autorisé par question.")
	}

	if nbLevel > 1 {
		return errors.New("Une seule classe est autorisée par question.")
	}

	_, err := DeleteQuestionTagsByIdQuestions(tx, idQuestion)
	if err != nil {
		return utils.SQLError(err)
	}
	err = InsertManyQuestionTags(tx, tags...)
	if err != nil {
		return utils.SQLError(err)
	}
	return nil
}

func (ct *Controller) updateTags(params UpdateTagsIn, userID uID) error {
	question, err := SelectQuestion(ct.db, params.IdQuestion)
	if err != nil {
		return utils.SQLError(err)
	}
	if question.IdTeacher != userID {
		return accessForbidden
	}

	var tags QuestionTags
	for _, tag := range params.Tags {
		// enforce proper tags
		tag = NormalizeTag(tag)
		if tag == "" {
			continue
		}

		tags = append(tags, QuestionTag{IdQuestion: params.IdQuestion, Tag: tag})
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return err
	}

	err = updateTags(tx, tags, params.IdQuestion)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

type UpdateGroupTagsOut struct {
	Tags map[IdQuestion][]string
}

func (ct *Controller) updateGroupTags(params UpdateGroupTagsIn, userID uID) (UpdateGroupTagsOut, error) {
	questions, err := SelectAllQuestions(ct.db)
	if err != nil {
		return UpdateGroupTagsOut{}, utils.SQLError(err)
	}

	var groupIDs []IdQuestion
	for _, question := range questions {
		if question.Page.Title == params.GroupTitle && question.IdTeacher == userID {
			groupIDs = append(groupIDs, question.Id)
		}
	}

	// compute the current common tags
	tags, err := SelectQuestionTagsByIdQuestions(ct.db, groupIDs...)
	if err != nil {
		return UpdateGroupTagsOut{}, utils.SQLError(err)
	}
	tagsByQuestion := tags.ByIdQuestion()
	var allTags [][]string
	for _, qus := range tagsByQuestion {
		allTags = append(allTags, qus.List())
	}
	commonTags := CommonTags(allTags)

	NormalizeTags(params.CommonTags)

	// replace commonTags by the input query
	crible := NewCrible(commonTags)
	tx, err := ct.db.Begin()
	if err != nil {
		return UpdateGroupTagsOut{}, utils.SQLError(err)
	}
	out := UpdateGroupTagsOut{Tags: make(map[IdQuestion][]string)}
	for _, idQuestion := range groupIDs {
		tags := tagsByQuestion[idQuestion]

		var newTags QuestionTags
		// start with the "exclusive" tags
		for _, tag := range tags {
			if !crible[tag.Tag] {
				newTags = append(newTags, tag)
			}
		}

		exclusive := newTags.Crible()
		// then add the new common tags, making sure
		// no duplicate is added
		for _, tag := range params.CommonTags {
			if exclusive[tag] {
				continue
			}
			newTags = append(newTags, QuestionTag{IdQuestion: idQuestion, Tag: tag})
		}

		// finally udpate the tags on DB
		err = updateTags(tx, newTags, idQuestion)
		if err != nil {
			_ = tx.Rollback()
			return out, err
		}

		out.Tags[idQuestion] = newTags.List()
	}

	err = tx.Commit()
	return out, err
}

func (ct *Controller) checkQuestionParameters(params CheckQuestionParametersIn) CheckQuestionParametersOut {
	err := params.Parameters.Validate()
	if err != nil {
		return CheckQuestionParametersOut{ErrDefinition: err.(questions.ErrParameters)}
	}

	var out CheckQuestionParametersOut
	for vr := range params.Parameters.ToMap() {
		out.Variables = append(out.Variables, vr)
	}
	sort.Slice(out.Variables, func(i, j int) bool {
		return out.Variables[i].String() < out.Variables[j].String()
	})

	return out
}

func (ct *Controller) pausePreview(sessionID string) error {
	ct.lock.Lock()
	defer ct.lock.Unlock()

	loopback, ok := ct.sessions[sessionID]
	if !ok {
		return fmt.Errorf("invalid session ID %s", sessionID)
	}

	loopback.pause()
	return nil
}

// endPreview terminates the current session
func (ct *Controller) endPreview(sessionID string) error {
	ct.lock.Lock()
	defer ct.lock.Unlock()

	loopback, ok := ct.sessions[sessionID]
	if !ok {
		return fmt.Errorf("invalid session ID %s", sessionID)
	}

	loopback.clientLeft <- true
	return nil
}

func (ct *Controller) saveQuestionAndPreview(params SaveQuestionAndPreviewIn, userID uID) (SaveQuestionAndPreviewOut, error) {
	qu, err := SelectQuestion(ct.db, params.Question.Id)
	if err != nil {
		return SaveQuestionAndPreviewOut{}, err
	}

	if !qu.IsVisibleBy(userID) {
		return SaveQuestionAndPreviewOut{}, accessForbidden
	}

	if err := params.Question.Page.Validate(); err != nil {
		return SaveQuestionAndPreviewOut{Error: err.(questions.ErrQuestionInvalid)}, nil
	}

	// if the question is owned : save it, else only preview
	if qu.IdTeacher == userID {
		_, err := params.Question.Update(ct.db)
		if err != nil {
			return SaveQuestionAndPreviewOut{}, utils.SQLError(err)
		}
	}

	question := params.Question.Page.Instantiate()

	ct.lock.Lock()
	defer ct.lock.Unlock()

	loopback, ok := ct.sessions[params.SessionID]
	if !ok {
		return SaveQuestionAndPreviewOut{}, fmt.Errorf("invalid session ID %s", params.SessionID)
	}

	loopback.setQuestion(question)
	return SaveQuestionAndPreviewOut{IsValid: true}, nil
}
