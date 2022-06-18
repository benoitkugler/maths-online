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
	"github.com/benoitkugler/maths-online/maths/questions/client"
	"github.com/benoitkugler/maths-online/prof/teacher"
	"github.com/benoitkugler/maths-online/utils"
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

type ListQuestionsOut struct {
	Questions []QuestionGroup // limited by `pagination`
	Size      int             // total number of groups
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
	Title        string
	Tags         []string
	Id           int64
	Difficulty   DifficultyTag // deduced from the tags
	IsInGroup    bool          // true if the question is in an implicit group, ignoring the current filter
	Origin       teacher.Origin
	NeedExercice bool
}

func normalizeTitle(title string) string {
	return removeAccents(strings.TrimSpace(strings.ToLower(title)))
}

func (ct *Controller) searchQuestions(query ListQuestionsIn, userID int64) (out ListQuestionsOut, err error) {
	const pagination = 100

	// to find implicit groups, we need all the questions
	questions, err := SelectAllQuestions(ct.db)
	if err != nil {
		return out, utils.SQLError(err)
	}
	questions.RestrictVisible(userID)

	// the group are not modified by the title query though

	queryTitle := normalizeTitle(query.TitleQuery)
	var (
		ids      IDs
		ownerIDs IDs
		groups   = make(map[string][]int64)
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
			vis, _ := teacher.NewVisibility(qu.IdTeacher, userID, ct.admin.Id, qu.Public)

			question := QuestionHeader{
				Id:         id,
				Title:      title,
				Difficulty: crible.Difficulty(),
				IsInGroup:  len(ids) > 1,
				Tags:       tagsMap[id].List(),
				Origin: teacher.Origin{
					AllowPublish: userID == ct.admin.Id,
					IsPublic:     qu.Public,
					Visibility:   vis,
				},
				NeedExercice: qu.NeedExercice,
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
	}

	// sort before pagination
	sort.Slice(out.Questions, func(i, j int) bool { return out.Questions[i].Title < out.Questions[j].Title })

	out.Size = len(out.Questions)
	if len(out.Questions) > pagination {
		out.Questions = out.Questions[:pagination]
	}

	return out, nil
}

// duplicateQuestion duplicate the given question, returning
// the newly created one
func (ct *Controller) duplicateQuestion(idQuestion, userID int64) (Question, error) {
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
func (ct *Controller) duplicateQuestionWithDifficulty(idQuestion, userID int64) error {
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
func updateTags(tx *sql.Tx, tags QuestionTags, idQuestion int64) error {
	var nbDiff int
	for _, tag := range tags {
		switch DifficultyTag(tag.Tag) {
		case Diff1, Diff2, Diff3:
			nbDiff++
		}
	}
	if nbDiff > 1 {
		return errors.New("Un seul niveau de difficulté est autorisé par question.")
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

func (ct *Controller) updateTags(params UpdateTagsIn, userID int64) error {
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
	Tags map[int64][]string
}

func (ct *Controller) updateGroupTags(params UpdateGroupTagsIn, userID int64) (UpdateGroupTagsOut, error) {
	questions, err := SelectAllQuestions(ct.db)
	if err != nil {
		return UpdateGroupTagsOut{}, utils.SQLError(err)
	}

	var groupIDs IDs
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
	out := UpdateGroupTagsOut{Tags: make(map[int64][]string)}
	for idQuestion, tags := range tagsByQuestion {
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
		err := updateTags(tx, newTags, idQuestion)
		if err != nil {
			_ = tx.Rollback()
			return out, err
		}

		out.Tags[idQuestion] = newTags.List()
	}

	err = tx.Commit()
	return out, err
}

func (ct *Controller) checkParameters(params CheckParametersIn) CheckParametersOut {
	err := params.Parameters.Validate()
	if err != nil {
		return CheckParametersOut{ErrDefinition: err.(questions.ErrParameters)}
	}

	var out CheckParametersOut
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

	loopback.unsetQuestion()
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

func (ct *Controller) saveAndPreview(params SaveAndPreviewIn, userID int64) (SaveAndPreviewOut, error) {
	qu, err := SelectQuestion(ct.db, params.Question.Id)
	if err != nil {
		return SaveAndPreviewOut{}, err
	}

	if !qu.IsVisibleBy(userID) {
		return SaveAndPreviewOut{}, accessForbidden
	}

	if err := params.Question.Page.Validate(); err != nil {
		return SaveAndPreviewOut{Error: err.(questions.ErrQuestionInvalid)}, nil
	}

	// if the question is owned : save it, else only preview
	if qu.IdTeacher == userID {
		_, err := params.Question.Update(ct.db)
		if err != nil {
			return SaveAndPreviewOut{}, utils.SQLError(err)
		}
	}

	question := params.Question.Page.Instantiate()

	ct.lock.Lock()
	defer ct.lock.Unlock()

	loopback, ok := ct.sessions[params.SessionID]
	if !ok {
		return SaveAndPreviewOut{}, fmt.Errorf("invalid session ID %s", params.SessionID)
	}

	loopback.setQuestion(question)
	return SaveAndPreviewOut{IsValid: true}, nil
}

// ------------------------------------------------------------------------------

type VarEntry struct {
	Variable expression.Variable
	Resolved string
}

type InstantiatedQuestion struct {
	Id       int64
	Question client.Question `dart-extern:"exercices/types.gen.dart"`
	Params   []VarEntry
}

type InstantiateQuestionsOut []InstantiatedQuestion

// InstantiateQuestions loads and instantiates the given questions,
// also returning the paramerters used to do so.
func (ct *Controller) InstantiateQuestions(ids []int64) (InstantiateQuestionsOut, error) {
	questions, err := SelectQuestions(ct.db, ids...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	out := make(InstantiateQuestionsOut, len(ids))
	for index, id := range ids {
		qu := questions[id]
		vars, err := qu.Page.Parameters.ToMap().Instantiate()
		if err != nil {
			return nil, err
		}
		varList := make([]VarEntry, 0, len(vars))
		for k, v := range vars {
			varList = append(varList, VarEntry{Variable: k, Resolved: v.Serialize()})
		}
		instance, err := qu.Page.InstantiateWith(vars)
		if err != nil {
			return nil, err
		}
		out[index] = InstantiatedQuestion{
			Id:       id,
			Question: instance.ToClient(),
			Params:   varList,
		}
	}

	return out, nil
}

// EvaluateQuestion instantiate the given question with the given parameters,
// and evaluate the given answer.
func (ct *Controller) EvaluateQuestion(id int64, params []VarEntry, answer client.QuestionAnswersIn) (client.QuestionAnswersOut, error) {
	qu, err := SelectQuestion(ct.db, id)
	if err != nil {
		return client.QuestionAnswersOut{}, utils.SQLError(err)
	}

	paramsDict := make(expression.Vars)
	for _, entry := range params {
		paramsDict[entry.Variable], err = expression.Parse(entry.Resolved)
		if err != nil {
			return client.QuestionAnswersOut{}, err
		}
	}

	instance, err := qu.Page.InstantiateWith(paramsDict)
	if err != nil {
		return client.QuestionAnswersOut{}, err
	}

	return instance.EvaluateAnswer(answer), nil
}
