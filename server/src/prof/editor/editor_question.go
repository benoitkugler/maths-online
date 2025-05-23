// Package editor provides functionnality for a frontend
// to edit and preview math questions.
package editor

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"

	"github.com/benoitkugler/maths-online/server/src/maths/expression"
	"github.com/benoitkugler/maths-online/server/src/maths/questions"
	"github.com/benoitkugler/maths-online/server/src/prof/preview"
	tcAPI "github.com/benoitkugler/maths-online/server/src/prof/teacher"
	ed "github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/sql/homework"
	"github.com/benoitkugler/maths-online/server/src/sql/reviews"
	"github.com/benoitkugler/maths-online/server/src/sql/tasks"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	taAPI "github.com/benoitkugler/maths-online/server/src/tasks"
	"github.com/benoitkugler/maths-online/server/src/utils"
	"github.com/labstack/echo/v4"
)

var errAccessForbidden = errors.New("access fordidden")

// Controller is the global object responsible to
// handle incoming requests regarding the editor for questions and exercices
type Controller struct {
	db *sql.DB

	admin teacher.Teacher
}

func NewController(db *sql.DB, admin teacher.Teacher) *Controller {
	return &Controller{
		db:    db,
		admin: admin,
	}
}

type OriginKind uint8

const (
	All OriginKind = iota
	OnlyPersonnal
	OnlyAdmin
)

func (ct *Controller) EditorGetQuestionsIndex(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	out, err := ct.loadQuestionsIndex(userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) loadQuestionsIndex(userID uID) (Index, error) {
	user, err := teacher.SelectTeacher(ct.db, userID)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	groups, err := ed.SelectAllQuestiongroups(ct.db)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	groups.RestrictVisible(userID)

	// load the tags ...
	tags, err := ed.SelectQuestiongroupTagsByIdQuestiongroups(ct.db, groups.IDs()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	return buildIndexFor(questionsToIndex(groups, tags), user.FavoriteMatiere), nil
}

type SearchQuestionsIn = Query

func (ct *Controller) EditorSearchQuestions(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args SearchQuestionsIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	out, err := ct.searchQuestions(args, userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

// EditorCreateQuestion creates a group with one question.
func (ct *Controller) EditorCreateQuestiongroup(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	out, err := ct.createQuestion(userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) createQuestion(userID uID) (QuestiongroupExt, error) {
	user, err := teacher.SelectTeacher(ct.db, userID)
	if err != nil {
		return QuestiongroupExt{}, utils.SQLError(err)
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return QuestiongroupExt{}, utils.SQLError(err)
	}

	group, err := ed.Questiongroup{IdTeacher: userID, Public: false}.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return QuestiongroupExt{}, utils.SQLError(err)
	}

	qu, err := ed.Question{
		IdGroup: group.Id.AsOptional(),
		Enonce:  questions.Enonce{questions.TextBlock{}}, // add a text block, very common in practice
	}.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return QuestiongroupExt{}, utils.SQLError(err)
	}

	// add the favorite matiere as tag
	ts := ed.TagSection{Section: ed.Matiere, Tag: string(user.FavoriteMatiere)}
	err = ed.InsertQuestiongroupTag(tx, ed.QuestiongroupTag{
		Tag:             ts.Tag,
		Section:         ts.Section,
		IdQuestiongroup: group.Id,
	})
	if err != nil {
		_ = tx.Rollback()
		return QuestiongroupExt{}, utils.SQLError(err)
	}

	err = tx.Commit()
	if err != nil {
		return QuestiongroupExt{}, utils.SQLError(err)
	}

	origin := questionOrigin(group, tcAPI.OptionalIdReview{}, userID, ct.admin.Id)
	return QuestiongroupExt{
		Group:    group,
		Tags:     ed.Tags{ts},
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
	userID := tcAPI.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	out, err := ct.duplicateQuestion(ed.IdQuestion(id), userID)
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
		return ed.Question{}, errAccessForbidden
	}

	// shallow copy is enough
	newQuestion := qu
	newQuestion.Subtitle += " (2)" // to distinguish from origin
	newQuestion, err = newQuestion.Insert(ct.db)
	if err != nil {
		return ed.Question{}, utils.SQLError(err)
	}

	return newQuestion, nil
}

// EditorDuplicateQuestiongroup duplicates the whole group, deep copying
// every question, and assigns it to the current user.
func (ct *Controller) EditorDuplicateQuestiongroup(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	err = ct.duplicateQuestiongroup(ed.IdQuestiongroup(id), userID)
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
		return errAccessForbidden
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
	err = ed.UpdateQuestiongroupTags(tx, tags.Tags(), newGroup.Id)
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
	userID := tcAPI.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	out, err := ct.deleteQuestion(ed.IdQuestion(id), userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

// EditorDeleteQuestiongroup remove the whole group.
// An information is returned if the question is used in monoquestions (tasks)
func (ct *Controller) EditorDeleteQuestiongroup(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	out, err := ct.deleteQuestiongroup(ed.IdQuestiongroup(id), userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

type TaskDetails struct {
	Id    tasks.IdTask
	Sheet homework.Sheet
}

type TaskUses []TaskDetails // the task containing the [Monoquestions]

func loadTaskDetails(db ed.DB, idTasks []tasks.IdTask) ([]TaskDetails, error) {
	links, err := homework.SelectSheetTasksByIdTasks(db, idTasks...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	sheets, err := homework.SelectSheets(db, links.IdSheets()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	taskToSheet := links.ByIdTask()

	out := make([]TaskDetails, len(idTasks))
	for i, id := range idTasks {
		sheet := sheets[taskToSheet[id].IdSheet]
		out[i] = TaskDetails{
			Id:    id,
			Sheet: sheet,
		}
	}

	return out, nil
}

// getQuestionUses returns the item using the given question (Monoquestion and RandomMonoquestion)
// exercices are not considered since questions in exercices can't be accessed directly
func getQuestionUses(db ed.DB, id ed.IdQuestion, idGroup ed.IdQuestiongroup) (out TaskUses, err error) {
	monoquestions, err := tasks.SelectMonoquestionsByIdQuestions(db, id)
	if err != nil {
		return out, utils.SQLError(err)
	}

	ts1, err := tasks.SelectTasksByIdMonoquestions(db, monoquestions.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	randMono, err := tasks.SelectRandomMonoquestionsByIdQuestiongroups(db, idGroup)
	if err != nil {
		return out, utils.SQLError(err)
	}

	ts2, err := tasks.SelectTasksByIdRandomMonoquestions(db, randMono.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	return loadTaskDetails(db, append(ts1.IDs(), ts2.IDs()...))
}

// getQuestiongroupUses returns the item using the given questiongroup
// or its variants
func getQuestiongroupUses(db ed.DB, id ed.IdQuestiongroup) (out TaskUses, err error) {
	// load the variants
	variants, err := ed.SelectQuestionsByIdGroups(db, id)
	if err != nil {
		return out, utils.SQLError(err)
	}

	monoquestions, err := tasks.SelectMonoquestionsByIdQuestions(db, variants.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	randomMonoquestions, err := tasks.SelectRandomMonoquestionsByIdQuestiongroups(db, id)
	if err != nil {
		return out, utils.SQLError(err)
	}

	tasks1, err := tasks.SelectTasksByIdMonoquestions(db, monoquestions.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	tasks2, err := tasks.SelectTasksByIdRandomMonoquestions(db, randomMonoquestions.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	// tasks1 and tasks2 are disjoint, by design
	return loadTaskDetails(db, append(tasks1.IDs(), tasks2.IDs()...))
}

type DeleteQuestionOut struct {
	Deleted   bool
	BlockedBy TaskUses // non empty iff Deleted == false
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
		return DeleteQuestionOut{}, errAccessForbidden
	}

	uses, err := getQuestionUses(ct.db, id, group.Id)
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

func (ct *Controller) deleteQuestiongroup(id ed.IdQuestiongroup, userID uID) (DeleteQuestionOut, error) {
	group, err := ed.SelectQuestiongroup(ct.db, id)
	if err != nil {
		return DeleteQuestionOut{}, utils.SQLError(err)
	}
	if group.IdTeacher != userID {
		return DeleteQuestionOut{}, errAccessForbidden
	}
	uses, err := getQuestiongroupUses(ct.db, id)
	if err != nil {
		return DeleteQuestionOut{}, err
	}
	if len(uses) != 0 {
		return DeleteQuestionOut{
			Deleted:   false,
			BlockedBy: uses,
		}, nil
	}

	// delete the variants then the group
	tx, err := ct.db.Begin()
	if err != nil {
		return DeleteQuestionOut{}, err
	}
	_, err = ed.DeleteQuestionsByIdGroups(tx, id)
	if err != nil {
		_ = tx.Rollback()
		return DeleteQuestionOut{}, err
	}
	_, err = ed.DeleteQuestiongroupById(tx, id)
	if err != nil {
		_ = tx.Rollback()
		return DeleteQuestionOut{}, err
	}
	err = tx.Commit()
	if err != nil {
		return DeleteQuestionOut{}, err
	}

	return DeleteQuestionOut{Deleted: true}, nil
}

// LoadQuestionVariants returns the question of the group [id],
// sorted by Id
func LoadQuestionVariants(db ed.DB, id ed.IdQuestiongroup) ([]ed.Question, error) {
	dict, err := ed.SelectQuestionsByIdGroups(db, id)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	var out []ed.Question
	for _, qu := range dict {
		out = append(out, qu)
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Id < out[j].Id })
	sort.SliceStable(out, func(i, j int) bool { return out[i].Difficulty < out[j].Difficulty })

	return out, nil
}

// EditorGetQuestions returns the questions for the given group
func (ct *Controller) EditorGetQuestions(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	idGroup, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	// check the access
	group, err := ed.SelectQuestiongroup(ct.db, ed.IdQuestiongroup(idGroup))
	if err != nil {
		return utils.SQLError(err)
	}

	if !group.IsVisibleBy(userID) {
		return errAccessForbidden
	}

	out, err := LoadQuestionVariants(ct.db, ed.IdQuestiongroup(idGroup))
	if err != nil {
		return utils.SQLError(err)
	}

	return c.JSON(200, out)
}

type CheckQuestionParametersIn struct {
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
	userID := tcAPI.JWTTeacher(c)

	var args ed.Questiongroup
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	group, err := ed.SelectQuestiongroup(ct.db, args.Id)
	if err != nil {
		return utils.SQLError(err)
	}

	if group.IdTeacher != userID {
		return errAccessForbidden
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
	userID := tcAPI.JWTTeacher(c)

	// we only accept public question from admin account
	if userID != ct.admin.Id {
		return errAccessForbidden
	}

	var args QuestionUpdateVisiblityIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	group, err := ed.SelectQuestiongroup(ct.db, args.ID)
	if err != nil {
		return utils.SQLError(err)
	}
	if group.IdTeacher != userID {
		return errAccessForbidden
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
	userID := tcAPI.JWTTeacher(c)

	var args SaveQuestionMetaIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	err := ct.saveQuestionMeta(args, userID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

type SaveQuestionAndPreviewIn struct {
	Id   ed.IdQuestion
	Page questions.QuestionPage

	// Set the initial view to display the correction,
	// instead of the enonce.
	ShowCorrection bool
}

type ListQuestionsOut struct {
	Groups      []QuestiongroupExt
	NbQuestions int // Number of variants in [Groups]
}

// QuestiongroupExt adds the question and tags to a QuestionGroup
// Standalone question are represented by a group of length one.
type QuestiongroupExt struct {
	Group    ed.Questiongroup
	Origin   tcAPI.Origin
	Tags     ed.Tags
	Variants []QuestionHeader
}

func NewQuestiongroupExt(group ed.Questiongroup, variants []ed.Question, tags ed.Tags,
	inReview tcAPI.OptionalIdReview, userID, adminID uID,
) QuestiongroupExt {
	origin := questionOrigin(group, inReview, userID, adminID)
	groupExt := QuestiongroupExt{
		Group:  group,
		Origin: origin,
		Tags:   tags,
	}

	for _, question := range variants {
		groupExt.Variants = append(groupExt.Variants, newQuestionHeader(question))
	}

	// sort to make sure the display is consistent between two queries
	sort.Slice(groupExt.Variants, func(i, j int) bool { return groupExt.Variants[i].Id < groupExt.Variants[j].Id })
	sort.SliceStable(groupExt.Variants, func(i, j int) bool { return groupExt.Variants[i].Difficulty < groupExt.Variants[j].Difficulty })

	return groupExt
}

// QuestionHeader is a summary of the meta data of a question
type QuestionHeader struct {
	Id            ed.IdQuestion
	Subtitle      string
	Difficulty    ed.DifficultyTag
	HasCorrection bool // has the question a non empty [Correction] content
}

func newQuestionHeader(question ed.Question) QuestionHeader {
	return QuestionHeader{
		Id:            question.Id,
		Subtitle:      question.Subtitle,
		Difficulty:    question.Difficulty,
		HasCorrection: len(question.Correction) != 0,
	}
}

func questionOrigin(qu ed.Questiongroup, inReview tcAPI.OptionalIdReview, userID, adminID uID) tcAPI.Origin {
	return tcAPI.Origin{
		Visibility:   tcAPI.NewVisibility(qu.IdTeacher, userID, adminID, qu.Public),
		IsInReview:   inReview,
		PublicStatus: tcAPI.NewPublicStatus(qu.IdTeacher, userID, adminID, qu.Public),
	}
}

func (ct *Controller) searchQuestions(query Query, userID uID) (out ListQuestionsOut, err error) {
	query.normalize()

	var groups ed.Questiongroups
	if isQueryTODO(query.TitleQuery) {
		questions, err := ed.SelectAllQuestions(ct.db)
		if err != nil {
			return out, utils.SQLError(err)
		}
		ids := ed.IdQuestiongroupSet{}
		for _, question := range questions {
			if question.NeedExercice.Valid {
				continue // ignore exercices
			}
			if question.Parameters.HasTODO() {
				ids.Add(question.IdGroup.ID)
			}
		}
		groups, err = ed.SelectQuestiongroups(ct.db, ids.Keys()...)
		if err != nil {
			return out, utils.SQLError(err)
		}
	} else if idVariant, isVariante := isQueryVariant(query.TitleQuery); isVariante {
		qu, err := ed.SelectQuestion(ct.db, ed.IdQuestion(idVariant))
		if err == sql.ErrNoRows {
			return out, fmt.Errorf("La question %d n'existe pas.", idVariant)
		} else if err != nil {
			return out, utils.SQLError(err)
		}

		if !qu.IdGroup.Valid {
			return out, errAccessForbidden
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

	revs, err := reviews.SelectAllReviewQuestions(ct.db)
	if err != nil {
		return out, utils.SQLError(err)
	}
	revsMap := revs.ByIdQuestion()

	// .. and build the groups, restricting to the ones matching the given tags
	for _, group := range groups {
		tags := tagsMap[group.Id].Tags()
		tagIndex := tags.BySection()
		if !query.matchTags(tagIndex) {
			continue
		}

		questions := questionsByGroup[group.Id]
		if len(questions) == 0 { // ignore empty groupExts
			continue
		}

		var inReview tcAPI.OptionalIdReview
		link, isInReview := revsMap[group.Id]
		if isInReview {
			inReview = tcAPI.OptionalIdReview{InReview: true, Id: link.IdReview}
		}

		groupExt := NewQuestiongroupExt(group, questions, tags, inReview, userID, ct.admin.Id)

		out.NbQuestions += len(groupExt.Variants)

		out.Groups = append(out.Groups, groupExt)
	}

	sort.Slice(out.Groups, func(i, j int) bool { return out.Groups[i].Group.Title < out.Groups[j].Group.Title })

	return out, nil
}

type UpdateQuestiongroupTagsIn struct {
	Id   ed.IdQuestiongroup
	Tags ed.Tags
}

func (ct *Controller) EditorUpdateQuestionTags(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args UpdateQuestiongroupTagsIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	err := ct.updateQuestionTags(args, userID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) updateQuestionTags(params UpdateQuestiongroupTagsIn, userID uID) error {
	group, err := ed.SelectQuestiongroup(ct.db, params.Id)
	if err != nil {
		return utils.SQLError(err)
	}
	if group.IdTeacher != userID {
		return errAccessForbidden
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return err
	}

	err = ed.UpdateQuestiongroupTags(tx, params.Tags, params.Id)
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
	out.Variables = params.Parameters.ToMap().DefinedVariables()
	sort.Slice(out.Variables, func(i, j int) bool {
		return out.Variables[i].String() < out.Variables[j].String()
	})

	return out
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
		return errAccessForbidden
	}

	qu.Subtitle = params.Question.Subtitle
	qu.Difficulty = params.Question.Difficulty

	_, err = params.Question.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}

type SaveQuestionAndPreviewOut struct {
	Error   questions.ErrQuestionInvalid
	IsValid bool
	Preview preview.LoopbackShowQuestion
}

// For non personnal questions, only preview.
func (ct *Controller) EditorSaveQuestionAndPreview(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args SaveQuestionAndPreviewIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	out, err := ct.saveQuestionAndPreview(args, userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) saveQuestionAndPreview(params SaveQuestionAndPreviewIn, userID uID) (SaveQuestionAndPreviewOut, error) {
	qu, err := ed.SelectQuestion(ct.db, params.Id)
	if err != nil {
		return SaveQuestionAndPreviewOut{}, utils.SQLError(err)
	}

	group, err := ct.getGroup(qu)
	if err != nil {
		return SaveQuestionAndPreviewOut{}, err
	}

	// if the question is in review, allow external user to preview it
	_, inReview, err := reviews.SelectReviewQuestionByIdQuestion(ct.db, group.Id)
	if err != nil {
		return SaveQuestionAndPreviewOut{}, utils.SQLError(err)
	}

	if !inReview && !group.IsVisibleBy(userID) {
		return SaveQuestionAndPreviewOut{}, errAccessForbidden
	}

	if err := params.Page.Validate(); err != nil {
		return SaveQuestionAndPreviewOut{Error: err.(questions.ErrQuestionInvalid)}, nil
	}

	// if the question is owned : save it, else only preview
	if group.IdTeacher == userID {
		qu.Enonce = params.Page.Enonce
		qu.Correction = params.Page.Correction
		qu.Parameters = params.Page.Parameters
		_, err := qu.Update(ct.db)
		if err != nil {
			return SaveQuestionAndPreviewOut{}, utils.SQLError(err)
		}
	}

	question, instanceParams, err := params.Page.InstantiateErr()
	if err != nil {
		return SaveQuestionAndPreviewOut{}, err
	}
	questionOut := preview.LoopbackShowQuestion{
		Question:       question.ToClient(),
		Params:         taAPI.NewParams(instanceParams),
		Origin:         params.Page,
		ShowCorrection: params.ShowCorrection,
	}

	return SaveQuestionAndPreviewOut{IsValid: true, Preview: questionOut}, nil
}

// EditorQuestionExportLateX instantiate the given question and generates a LaTeX version,
// returning the code as a string
func (ct *Controller) EditorQuestionExportLateX(c echo.Context) error {
	var args questions.QuestionPage
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	out, err := exportQuestionLatex(args)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

type ExportQuestionLatexOut struct {
	Error   questions.ErrQuestionInvalid
	IsValid bool
	Latex   string
}

func exportQuestionLatex(question questions.QuestionPage) (ExportQuestionLatexOut, error) {
	if err := question.Validate(); err != nil {
		return ExportQuestionLatexOut{Error: err.(questions.ErrQuestionInvalid)}, nil
	}

	instance, _, err := question.InstantiateErr()
	if err != nil {
		return ExportQuestionLatexOut{}, err
	}

	return ExportQuestionLatexOut{
		IsValid: true,
		Latex:   instance.Enonce.ToLatex(),
	}, nil
}
