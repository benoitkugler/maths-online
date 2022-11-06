// Package editor provides functionnality for a frontend
// to edit and preview math questions.
package editor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/questions"
	tcAPI "github.com/benoitkugler/maths-online/prof/teacher"
	ed "github.com/benoitkugler/maths-online/sql/editor"
	"github.com/benoitkugler/maths-online/sql/homework"
	"github.com/benoitkugler/maths-online/sql/tasks"
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

type OriginKind uint8

const (
	All OriginKind = iota
	OnlyPersonnal
	OnlyAdmin
)

type SearchQuestionsIn = Query

type Query struct {
	TitleQuery string // empty means all
	Tags       []string
	Origin     OriginKind
}

func (query Query) normalize() {
	// normalize query
	for i, t := range query.Tags {
		query.Tags[i] = ed.NormalizeTag(t)
	}
}

func (query Query) matchOrigin(vis tcAPI.Visibility) bool {
	switch query.Origin {
	case OnlyPersonnal:
		return vis == tcAPI.Personnal
	case OnlyAdmin:
		return vis == tcAPI.Admin
	default:
		return true
	}
}

func (ct *Controller) EditorSearchQuestions(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	var args SearchQuestionsIn
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

func (ct *Controller) createQuestion(userID uID) (QuestiongroupExt, error) {
	tx, err := ct.db.Begin()
	if err != nil {
		return QuestiongroupExt{}, utils.SQLError(err)
	}

	group, err := ed.Questiongroup{IdTeacher: userID, Public: false}.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return QuestiongroupExt{}, utils.SQLError(err)
	}

	qu, err := ed.Question{IdGroup: group.Id.AsOptional()}.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return QuestiongroupExt{}, utils.SQLError(err)
	}

	err = tx.Commit()
	if err != nil {
		return QuestiongroupExt{}, utils.SQLError(err)
	}

	origin, _ := questionOrigin(group, userID, ct.admin.Id)
	return QuestiongroupExt{
		Group:    group,
		Tags:     nil,
		Origin:   origin,
		Variants: []QuestionHeader{newQuestionHeader(qu)},
	}, nil
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
// every question, and assigns it to the current user.
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

// duplicateQuestiongroup creates a new group with the same title, (copied) questions and tags
func (ct *Controller) duplicateQuestiongroup(idGroup ed.IdQuestiongroup, userID uID) error {
	group, err := ed.SelectQuestiongroup(ct.db, idGroup)
	if err != nil {
		return utils.SQLError(err)
	}

	if !group.IsVisibleBy(userID) {
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
	newGroup := group
	newGroup.IdTeacher = userID
	newGroup.Public = false
	newGroup, err = newGroup.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}
	// .. then add its tags ..
	err = ed.UpdateQuestiongroupTags(tx, tags, newGroup.Id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	// finaly, copy the questions...
	for _, question := range questions {
		question.IdGroup = newGroup.Id.AsOptional() // re-direct to the new group
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

// EditorDeleteQuestion remove the given question,
// also removing the group if needed.
// An information is returned if the question is used in monoquestions (tasks)
func (ct *Controller) EditorDeleteQuestion(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	out, err := ct.deleteQuestion(ed.IdQuestion(id), user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

type TaskDetails struct {
	Id        tasks.IdTask
	Sheet     homework.Sheet
	Classroom teacher.Classroom
}

type QuestionExerciceUses []TaskDetails // the task containing the [Monoquestions]

func newQuestionExericeUses(db ed.DB, idTasks []tasks.IdTask) ([]TaskDetails, error) {
	links, err := homework.SelectSheetTasksByIdTasks(db, idTasks...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	sheets, err := homework.SelectSheets(db, links.IdSheets()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	taskToSheet := links.ByIdTask()

	classrooms, err := teacher.SelectClassrooms(db, sheets.IdClassrooms()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	out := make([]TaskDetails, len(idTasks))
	for i, id := range idTasks {
		sheet := sheets[taskToSheet[id].IdSheet]
		out[i] = TaskDetails{
			Id:        id,
			Sheet:     sheet,
			Classroom: classrooms[sheet.IdClassroom],
		}
	}

	return out, nil
}

// getQuestionUses returns the item using the given question
// exercices are not considered since questions in exercices can't be accessed directly
func getQuestionUses(db ed.DB, id ed.IdQuestion) (out QuestionExerciceUses, err error) {
	monoquestions, err := tasks.SelectMonoquestionsByIdQuestions(db, id)
	if err != nil {
		return out, utils.SQLError(err)
	}

	tasks, err := tasks.SelectTasksByIdMonoquestions(db, monoquestions.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	return newQuestionExericeUses(db, tasks.IDs())
}

type DeleteQuestionOut struct {
	Deleted   bool
	BlockedBy QuestionExerciceUses // non empty iff Deleted == false
}

func (ct *Controller) deleteQuestion(id ed.IdQuestion, userID uID) (DeleteQuestionOut, error) {
	qu, err := ed.SelectQuestion(ct.db, id)
	if err != nil {
		return DeleteQuestionOut{}, utils.SQLError(err)
	}

	group, err := ct.getGroup(qu)
	if err != nil {
		return DeleteQuestionOut{}, err
	}

	if group.IdTeacher != userID {
		return DeleteQuestionOut{}, accessForbidden
	}

	uses, err := getQuestionUses(ct.db, id)
	if err != nil {
		return DeleteQuestionOut{}, err
	}
	if len(uses) != 0 {
		return DeleteQuestionOut{
			Deleted:   false,
			BlockedBy: uses,
		}, nil
	}

	_, err = ed.DeleteQuestionById(ct.db, id)
	if err != nil {
		return DeleteQuestionOut{}, utils.SQLError(err)
	}

	// check if group is empty
	ques, err := ed.SelectQuestionsByIdGroups(ct.db, group.Id)
	if err != nil {
		return DeleteQuestionOut{}, utils.SQLError(err)
	}
	if len(ques) == 0 {
		_, err = ed.DeleteQuestiongroupById(ct.db, group.Id)
		if err != nil {
			return DeleteQuestionOut{}, utils.SQLError(err)
		}
	}

	return DeleteQuestionOut{Deleted: true}, nil
}

// EditorGetQuestions returns the questions for the given group
func (ct *Controller) EditorGetQuestions(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	idGroup, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	// check the access
	group, err := ed.SelectQuestiongroup(ct.db, ed.IdQuestiongroup(idGroup))
	if err != nil {
		return utils.SQLError(err)
	}

	if !group.IsVisibleBy(user.Id) {
		return accessForbidden
	}

	dict, err := ed.SelectQuestionsByIdGroups(ct.db, ed.IdQuestiongroup(idGroup))
	if err != nil {
		return utils.SQLError(err)
	}

	var out []ed.Question
	for _, qu := range dict {
		out = append(out, qu)
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Id < out[j].Id })
	sort.SliceStable(out, func(i, j int) bool { return out[i].Difficulty < out[j].Difficulty })

	return c.JSON(200, out)
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

func (ct *Controller) EditorUpdateQuestiongroup(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	var args ed.Questiongroup
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	group, err := ed.SelectQuestiongroup(ct.db, args.Id)
	if err != nil {
		return utils.SQLError(err)
	}

	if group.IdTeacher != user.Id {
		return accessForbidden
	}

	group.Title = args.Title
	_, err = group.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	return c.NoContent(200)
}

type QuestionUpdateVisiblityIn struct {
	ID     ed.IdQuestiongroup
	Public bool
}

func (ct *Controller) EditorUpdateQuestiongroupVis(c echo.Context) error {
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

type SaveQuestionMetaIn struct {
	Question ed.Question
}

// EditorSaveQuestionMeta only save the meta data of the question,
// not its content.
func (ct *Controller) EditorSaveQuestionMeta(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	var args SaveQuestionMetaIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	err := ct.saveQuestionMeta(args, user.Id)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

type SaveQuestionAndPreviewIn struct {
	SessionID string
	Id        ed.IdQuestion
	Page      questions.QuestionPage
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
	Group    ed.Questiongroup
	Origin   tcAPI.Origin
	Tags     []string
	Variants []QuestionHeader
}

// QuestionHeader is a summary of the meta data of a question
type QuestionHeader struct {
	Id         ed.IdQuestion
	Subtitle   string
	Difficulty ed.DifficultyTag
}

func newQuestionHeader(question ed.Question) QuestionHeader {
	return QuestionHeader{
		Id:         question.Id,
		Subtitle:   question.Subtitle,
		Difficulty: question.Difficulty,
	}
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

func isQueryVariant(query string) (int64, bool) {
	if strings.HasPrefix(query, "variante:") {
		idS := strings.TrimPrefix(query, "variante:")
		id, err := strconv.ParseInt(idS, 10, 64)
		return id, err == nil
	}
	return 0, false
}

type itemGroupQuery interface {
	match(id int64, title string) bool
}

func newQuery(title string) (itemGroupQuery, error) {
	// special case for pattern id:id1, id2, ...
	if strings.HasPrefix(title, "id:") {
		idsString := strings.TrimSpace(title[len("id:"):])
		if len(idsString) != 0 {
			out := make(queryByIds)
			ids := strings.Split(idsString, ",")

			for _, id := range ids {
				idV, err := strconv.Atoi(id)
				if err != nil {
					return nil, fmt.Errorf("Requête invalide: entier attendu (%s)", err)
				}
				out[int64(idV)] = true
			}
			return out, nil
		}
	}
	return queryByTitle(normalizeTitle(title)), nil
}

type queryByTitle string

func (qt queryByTitle) match(_ int64, title string) bool {
	itemTitle := normalizeTitle(title)
	return qt == "" || strings.Contains(itemTitle, string(qt))
}

type queryByIds map[int64]bool

func (qt queryByIds) match(id int64, _ string) bool { return qt[id] }

func normalizeTitle(title string) string {
	return utils.RemoveAccents(strings.TrimSpace(strings.ToLower(title)))
}

func (ct *Controller) searchQuestions(query Query, userID uID) (out ListQuestionsOut, err error) {
	const pagination = 30 // number of groups

	query.normalize()

	var groups ed.Questiongroups
	if idVariant, isVariante := isQueryVariant(query.TitleQuery); isVariante {
		qu, err := ed.SelectQuestion(ct.db, ed.IdQuestion(idVariant))
		if err != nil {
			return out, utils.SQLError(err)
		}
		if !qu.IdGroup.Valid {
			return out, accessForbidden
		}
		groups, err = ed.SelectQuestiongroups(ct.db, qu.IdGroup.ID)
		if err != nil {
			return out, utils.SQLError(err)
		}
	} else {
		groups, err = ed.SelectAllQuestiongroups(ct.db)
		if err != nil {
			return out, utils.SQLError(err)
		}
		matcher, err := newQuery(query.TitleQuery)
		if err != nil {
			return out, err
		}

		// restrict the groups to matching title and origin
		for _, group := range groups {
			vis := tcAPI.NewVisibility(group.IdTeacher, userID, ct.admin.Id, group.Public)

			keep := query.matchOrigin(vis) && matcher.match(int64(group.Id), group.Title)
			if !keep {
				delete(groups, group.Id)
			}
		}
	}

	groups.RestrictVisible(userID)

	// load the tags ...
	tags, err := ed.SelectQuestiongroupTagsByIdQuestiongroups(ct.db, groups.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}
	tagsMap := tags.ByIdQuestiongroup()

	// ... and the questions
	tmp, err := ed.SelectQuestionsByIdGroups(ct.db, groups.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}
	questionsByGroup := tmp.ByGroup()

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
			groupExt.Variants = append(groupExt.Variants, newQuestionHeader(question))
		}

		// sort to make sure the display is consistent between two queries
		sort.Slice(groupExt.Variants, func(i, j int) bool { return groupExt.Variants[i].Id < groupExt.Variants[j].Id })
		sort.SliceStable(groupExt.Variants, func(i, j int) bool { return groupExt.Variants[i].Difficulty < groupExt.Variants[j].Difficulty })

		out.NbQuestions += len(groupExt.Variants)

		out.Groups = append(out.Groups, groupExt)
	}

	// sort before pagination
	sort.Slice(out.Groups, func(i, j int) bool { return out.Groups[i].Group.Title < out.Groups[j].Group.Title })

	out.NbGroups = len(out.Groups)
	if len(out.Groups) > pagination {
		out.Groups = out.Groups[:pagination]
	}

	return out, nil
}

type UpdateQuestiongroupTagsIn struct {
	Id   ed.IdQuestiongroup
	Tags []string
}

func (ct *Controller) EditorUpdateQuestionTags(c echo.Context) error {
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

	err = ed.UpdateQuestiongroupTags(tx, tags, params.Id)
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

// return the owner of the group of of the exercice
func (ct *Controller) getQuestionOwner(question ed.Question) (teacher.IdTeacher, error) {
	if question.IdGroup.Valid {
		group, err := ct.getGroup(question)
		if err != nil {
			return 0, err
		}
		return group.IdTeacher, nil
	}
	ex, err := ed.SelectExercice(ct.db, question.NeedExercice.ID)
	if err != nil {
		return 0, utils.SQLError(err)
	}
	group, err := ed.SelectExercicegroup(ct.db, ex.IdGroup)
	if err != nil {
		return 0, utils.SQLError(err)
	}
	return group.IdTeacher, nil
}

func (ct *Controller) saveQuestionMeta(params SaveQuestionMetaIn, userID uID) error {
	qu, err := ed.SelectQuestion(ct.db, params.Question.Id)
	if err != nil {
		return err
	}

	owner, err := ct.getQuestionOwner(qu)
	if err != nil {
		return err
	}

	if owner != userID {
		return accessForbidden
	}

	qu.Description = params.Question.Description
	qu.Subtitle = params.Question.Subtitle
	qu.Difficulty = params.Question.Difficulty

	_, err = params.Question.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}

func (ct *Controller) saveQuestionAndPreview(params SaveQuestionAndPreviewIn, userID uID) (SaveQuestionAndPreviewOut, error) {
	qu, err := ed.SelectQuestion(ct.db, params.Id)
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

	if err := params.Page.Validate(); err != nil {
		return SaveQuestionAndPreviewOut{Error: err.(questions.ErrQuestionInvalid)}, nil
	}

	// if the question is owned : save it, else only preview
	if group.IdTeacher == userID {
		qu.Page = params.Page
		_, err := qu.Update(ct.db)
		if err != nil {
			return SaveQuestionAndPreviewOut{}, utils.SQLError(err)
		}
	}

	question := params.Page.Instantiate()

	ct.lock.Lock()
	defer ct.lock.Unlock()

	loopback, ok := ct.sessions[params.SessionID]
	if !ok {
		return SaveQuestionAndPreviewOut{}, fmt.Errorf("invalid session ID %s", params.SessionID)
	}

	loopback.setQuestion(question)
	return SaveQuestionAndPreviewOut{IsValid: true}, nil
}
