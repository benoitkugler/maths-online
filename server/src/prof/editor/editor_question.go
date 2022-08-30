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
	tcAPI "github.com/benoitkugler/maths-online/prof/teacher"
	ed "github.com/benoitkugler/maths-online/sql/editor"
	"github.com/benoitkugler/maths-online/sql/teacher"
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
	user := tcAPI.JWTTeacher(c)

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

// EditorCreateQuestion creates a group with one question.
func (ct *Controller) EditorCreateQuestiongroup(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	out, err := ct.createQuestion(user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) createQuestion(userID uID) (ed.Questiongroup, error) {
	tx, err := ct.db.Begin()
	if err != nil {
		return ed.Questiongroup{}, utils.SQLError(err)
	}

	group, err := ed.Questiongroup{IdTeacher: userID, Public: false}.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return ed.Questiongroup{}, utils.SQLError(err)
	}

	_, err = ed.Question{IdGroup: group.Id.AsOptional()}.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return ed.Questiongroup{}, utils.SQLError(err)
	}

	err = tx.Commit()
	if err != nil {
		return ed.Questiongroup{}, utils.SQLError(err)
	}

	return group, nil
}

func (ct *Controller) getGroup(qu ed.Question) (ed.Questiongroup, error) {
	if !qu.IdGroup.Valid {
		return ed.Questiongroup{}, errors.New("internal error: question without group")
	}

	group, err := ed.SelectQuestiongroup(ct.db, qu.IdGroup.ID)
	if err != nil {
		return ed.Questiongroup{}, utils.SQLError(err)
	}
	return group, nil
}

// EditorDuplicateQuestion duplicate one variant inside a group.
func (ct *Controller) EditorDuplicateQuestion(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	out, err := ct.duplicateQuestion(ed.IdQuestion(id), user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

// duplicateQuestion duplicate the given question, returning
// the newly created one
func (ct *Controller) duplicateQuestion(idQuestion ed.IdQuestion, userID uID) (ed.Question, error) {
	qu, err := ed.SelectQuestion(ct.db, idQuestion)
	if err != nil {
		return ed.Question{}, utils.SQLError(err)
	}

	group, err := ct.getGroup(qu)
	if err != nil {
		return ed.Question{}, err
	}

	if !group.IsVisibleBy(userID) {
		return ed.Question{}, accessForbidden
	}

	// shallow copy is enough
	newQuestion := qu
	newQuestion, err = newQuestion.Insert(ct.db)
	if err != nil {
		return ed.Question{}, utils.SQLError(err)
	}

	return newQuestion, nil
}

// EditorDuplicateQuestiongroup duplicates the whole group, deep copying
// every question.
func (ct *Controller) EditorDuplicateQuestiongroup(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	err = ct.duplicateQuestiongroup(ed.IdQuestiongroup(id), user.Id)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

// EditorDeleteQuestion remove the given question,
// also removing the group if needed.
// TODO: check usage in exercices and tasks
func (ct *Controller) EditorDeleteQuestion(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	err = ct.deleteQuestion(ed.IdQuestion(id), user.Id)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) deleteQuestion(id ed.IdQuestion, userID uID) error {
	qu, err := ed.SelectQuestion(ct.db, id)
	if err != nil {
		return utils.SQLError(err)
	}

	group, err := ct.getGroup(qu)
	if err != nil {
		return err
	}

	if group.IdTeacher != userID {
		return accessForbidden
	}

	links, err := ed.SelectExerciceQuestionsByIdQuestions(ct.db, id)
	if err != nil {
		return utils.SQLError(err)
	}

	if len(links) != 0 {
		ex, err := ed.SelectExercice(ct.db, links[0].IdExercice)
		if err != nil {
			return utils.SQLError(err)
		}

		return fmt.Errorf("La question est utilisée dans l'exercice %s : %d.", ex.Subtitle, ex.Id)
	}

	_, err = ed.DeleteQuestionById(ct.db, id)
	if err != nil {
		return utils.SQLError(err)
	}

	// check if group is empty
	ques, err := ed.SelectQuestionsByIdGroups(ct.db, group.Id)
	if err != nil {
		return utils.SQLError(err)
	}
	if len(ques) == 0 {
		_, err = ed.DeleteQuestiongroupById(ct.db, group.Id)
		if err != nil {
			return utils.SQLError(err)
		}
	}

	return nil
}

func (ct *Controller) EditorGetQuestion(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	question, err := ed.SelectQuestion(ct.db, ed.IdQuestion(id))
	if err != nil {
		return err
	}

	group, err := ct.getGroup(question)
	if err != nil {
		return err
	}
	if !group.IsVisibleBy(user.Id) {
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
	ID     ed.IdQuestiongroup
	Public bool
}

func (ct *Controller) QuestionUpdateVisiblity(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	// we only accept public question from admin account
	if user.Id != ct.admin.Id {
		return accessForbidden
	}

	var args QuestionUpdateVisiblityIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	group, err := ed.SelectQuestiongroup(ct.db, args.ID)
	if err != nil {
		return utils.SQLError(err)
	}
	if group.IdTeacher != user.Id {
		return accessForbidden
	}

	group.Public = args.Public
	group, err = group.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	return c.NoContent(200)
}

type SaveQuestionAndPreviewIn struct {
	SessionID string
	Question  ed.Question
}

type SaveQuestionAndPreviewOut struct {
	Error   questions.ErrQuestionInvalid
	IsValid bool
}

// For non personnal questions, only preview.
func (ct *Controller) EditorSaveQuestionAndPreview(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

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
	Groups      []QuestiongroupExt // limited by `pagination`
	NbGroups    int                // total number of groups (passing the given filter)
	NbQuestions int                // total number of questions
}

// QuestiongroupExt adds the question and tags to a QuestionGroup
// Standalone question are represented by a group of length one.
type QuestiongroupExt struct {
	Group     ed.Questiongroup
	Origin    tcAPI.Origin
	Tags      []string
	Questions []QuestionHeader
}

// QuestionHeader is a sumary of the meta data of a question
type QuestionHeader struct {
	Id         ed.IdQuestion
	Subtitle   string
	Difficulty ed.DifficultyTag // deduced from the tags
}

func normalizeTitle(title string) string {
	return utils.RemoveAccents(strings.TrimSpace(strings.ToLower(title)))
}

func questionOrigin(qu ed.Questiongroup, userID, adminID uID) (tcAPI.Origin, bool) {
	vis := tcAPI.NewVisibility(qu.IdTeacher, userID, adminID, qu.Public)
	if vis.Restricted() {
		return tcAPI.Origin{}, false
	}
	return tcAPI.Origin{
		AllowPublish: userID == adminID,
		IsPublic:     qu.Public,
		Visibility:   vis,
	}, true
}

func (ct *Controller) searchQuestions(query ListQuestionsIn, userID uID) (out ListQuestionsOut, err error) {
	const pagination = 30 // number of groups

	groups, err := ed.SelectAllQuestiongroups(ct.db)
	if err != nil {
		return out, utils.SQLError(err)
	}
	groups.RestrictVisible(userID)

	// restrict the groups to matching title
	queryTitle := normalizeTitle(query.TitleQuery)
	for _, group := range groups {
		thisTitle := normalizeTitle(group.Title)
		if queryTitle != "" && !strings.Contains(thisTitle, queryTitle) {
			delete(groups, group.Id)
		}
	}

	// load the tags ...
	tags, err := ed.SelectQuestiongroupTagsByIdQuestiongroups(ct.db, groups.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}
	tagsMap := tags.ByIdQuestiongroup()

	// ... and the tmp
	tmp, err := ed.SelectQuestionsByIdGroups(ct.db, groups.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}
	questionsByGroup := tmp.ByGroup()

	// normalize query
	for i, t := range query.Tags {
		query.Tags[i] = ed.NormalizeTag(t)
	}

	// .. and build the groups, restricting to the ones matching the given tags
	for _, group := range groups {
		crible := tagsMap[group.Id].Crible()
		if !crible.HasAll(query.Tags) {
			continue
		}
		questions := questionsByGroup[group.Id]
		if len(questions) == 0 { // ignore empty groupExts
			continue
		}

		origin, _ := questionOrigin(group, userID, ct.admin.Id)
		groupExt := QuestiongroupExt{
			Group:  group,
			Origin: origin,
			Tags:   tagsMap[group.Id].List(),
		}

		for _, question := range questions {
			question := QuestionHeader{
				Id:         question.Id,
				Subtitle:   question.Subtitle,
				Difficulty: question.Difficulty,
			}
			groupExt.Questions = append(groupExt.Questions, question)
		}

		// sort to make sure the display is consistent between two queries
		sort.Slice(groupExt.Questions, func(i, j int) bool { return groupExt.Questions[i].Id < groupExt.Questions[j].Id })
		sort.SliceStable(groupExt.Questions, func(i, j int) bool { return groupExt.Questions[i].Difficulty < groupExt.Questions[j].Difficulty })

		out.NbQuestions += len(groupExt.Questions)
	}

	// sort before pagination
	sort.Slice(out.Groups, func(i, j int) bool { return out.Groups[i].Group.Title < out.Groups[j].Group.Title })

	out.NbGroups = len(out.Groups)
	if len(out.Groups) > pagination {
		out.Groups = out.Groups[:pagination]
	}

	return out, nil
}

// duplicateQuestiongroup creates a new group with the same title, questions and tags
// only personnal questions are allowed
func (ct *Controller) duplicateQuestiongroup(idGroup ed.IdQuestiongroup, userID uID) error {
	group, err := ed.SelectQuestiongroup(ct.db, idGroup)
	if err != nil {
		return utils.SQLError(err)
	}

	if group.IdTeacher != userID {
		return accessForbidden
	}

	tags, err := ed.SelectQuestiongroupTagsByIdQuestiongroups(ct.db, group.Id)
	if err != nil {
		return utils.SQLError(err)
	}

	questions, err := ed.SelectQuestionsByIdGroups(ct.db, group.Id)
	if err != nil {
		return utils.SQLError(err)
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}

	// start by inserting a new group ...
	newGroup, err := group.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}
	// .. then add its tags ..
	err = ed.UpdateTags(tx, tags, newGroup.Id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	// finaly, copy the questions...
	for _, question := range questions {
		question.IdGroup = newGroup.Id.AsOptional() // re-direct the group
		_, err = question.Insert(tx)
		if err != nil {
			_ = tx.Rollback()
			return utils.SQLError(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}

type UpdateQuestiongroupTagsIn struct {
	Id   ed.IdQuestiongroup
	Tags []string
}

func (ct *Controller) EditorUpdateTags(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	var args UpdateQuestiongroupTagsIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	err := ct.updateTags(args, user.Id)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) updateTags(params UpdateQuestiongroupTagsIn, userID uID) error {
	group, err := ed.SelectQuestiongroup(ct.db, params.Id)
	if err != nil {
		return utils.SQLError(err)
	}
	if group.IdTeacher != userID {
		return accessForbidden
	}

	var tags ed.QuestiongroupTags
	for _, tag := range params.Tags {
		// enforce proper tags
		tag = ed.NormalizeTag(tag)
		if tag == "" {
			continue
		}

		tags = append(tags, ed.QuestiongroupTag{IdQuestiongroup: params.Id, Tag: tag})
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return err
	}

	err = ed.UpdateTags(tx, tags, params.Id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
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
	qu, err := ed.SelectQuestion(ct.db, params.Question.Id)
	if err != nil {
		return SaveQuestionAndPreviewOut{}, err
	}

	group, err := ct.getGroup(qu)
	if err != nil {
		return SaveQuestionAndPreviewOut{}, err
	}

	if !group.IsVisibleBy(userID) {
		return SaveQuestionAndPreviewOut{}, accessForbidden
	}

	if err := params.Question.Page.Validate(); err != nil {
		return SaveQuestionAndPreviewOut{Error: err.(questions.ErrQuestionInvalid)}, nil
	}

	// if the question is owned : save it, else only preview
	if group.IdTeacher == userID {
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
