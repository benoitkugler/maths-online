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

	ex "github.com/benoitkugler/maths-online/maths/exercice"
	"github.com/benoitkugler/maths-online/maths/exercice/client"
	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/utils"
)

const sessionTimeout = 6 * time.Hour

// Controller is the global object responsible to
// handle incoming requests regarding the editor.
type Controller struct {
	lock sync.Mutex

	db *sql.DB

	sessions map[string]*loopbackController
}

func NewController(db *sql.DB) *Controller {
	return &Controller{
		db:       db,
		sessions: make(map[string]*loopbackController),
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
	Title      string
	Tags       []string
	Id         int64
	Difficulty ex.DifficultyTag // deduced from the tags
	IsInGroup  bool             // true if the question is in an implicit group, ignoring the current filter
}

func (ct *Controller) searchQuestions(query ListQuestionsIn) (out ListQuestionsOut, err error) {
	const pagination = 100

	// to find implicit groups, we need all the questions
	questions, err := ex.SelectAllQuestions(ct.db)
	if err != nil {
		return out, utils.SQLError(err)
	}

	// the group are not modified by the title query though

	queryTitle := strings.TrimSpace(strings.ToLower(query.TitleQuery))
	var (
		ids    ex.IDs
		groups = make(map[string][]int64)
	)
	for _, question := range questions {
		thisTitle := strings.TrimSpace(strings.ToLower(question.Title))
		if strings.Contains(thisTitle, queryTitle) {
			groups[question.Title] = append(groups[question.Title], question.Id)
			ids = append(ids, question.Id)
		}
	}

	// load the tags ...
	tags, err := ex.SelectQuestionTagsByIdQuestions(ct.db, ids...)
	if err != nil {
		return out, utils.SQLError(err)
	}
	tagsMap := tags.ByIdQuestion()

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

			question := QuestionHeader{
				Id:         id,
				Title:      title,
				Difficulty: crible.Difficulty(),
				IsInGroup:  len(ids) > 1,
				Tags:       tagsMap[id].List(),
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
func (ct *Controller) duplicateQuestion(idQuestion int64) (ex.Question, error) {
	qu, err := ex.SelectQuestion(ct.db, idQuestion)
	if err != nil {
		return ex.Question{}, utils.SQLError(err)
	}
	tags, err := ex.SelectQuestionTagsByIdQuestions(ct.db, qu.Id)
	if err != nil {
		return ex.Question{}, utils.SQLError(err)
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return ex.Question{}, utils.SQLError(err)
	}

	newQuestion := qu // shallow copy is enough
	newQuestion, err = newQuestion.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return ex.Question{}, utils.SQLError(err)
	}

	for i := range tags {
		tags[i].IdQuestion = newQuestion.Id
	}
	err = updateTags(tx, tags, newQuestion.Id)
	if err != nil {
		_ = tx.Rollback()
		return ex.Question{}, utils.SQLError(err)
	}

	err = tx.Commit()
	if err != nil {
		return ex.Question{}, utils.SQLError(err)
	}

	return newQuestion, nil
}

// duplicateQuestionWithDifficulty creates new questions with the same title
// and content as the given question, but with difficulty levels
func (ct *Controller) duplicateQuestionWithDifficulty(idQuestion int64) error {
	qu, err := ex.SelectQuestion(ct.db, idQuestion)
	if err != nil {
		return utils.SQLError(err)
	}
	tags, err := ex.SelectQuestionTagsByIdQuestions(ct.db, qu.Id)
	if err != nil {
		return utils.SQLError(err)
	}

	// if the question already has a difficulty, respect it
	// otherwise, attribute the difficulty one
	currentDifficulty := tags.Crible().Difficulty()
	var newDifficulties [2]ex.DifficultyTag
	switch currentDifficulty {
	case ex.Diff1:
		newDifficulties = [2]ex.DifficultyTag{ex.Diff2, ex.Diff3}
	case ex.Diff2:
		newDifficulties = [2]ex.DifficultyTag{ex.Diff1, ex.Diff3}
	case ex.Diff3:
		newDifficulties = [2]ex.DifficultyTag{ex.Diff1, ex.Diff2}
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}

	if currentDifficulty == "" {
		// update the current question
		newTags := append(tags, ex.QuestionTag{IdQuestion: idQuestion, Tag: string(ex.Diff1)})
		err = updateTags(tx, newTags, idQuestion)
		if err != nil {
			_ = tx.Rollback()
			return utils.SQLError(err)
		}
		newDifficulties = [2]ex.DifficultyTag{ex.Diff2, ex.Diff3}
	}

	for _, diff := range newDifficulties {
		newQuestion := qu // shallow copy is enough
		newQuestion, err = newQuestion.Insert(tx)
		if err != nil {
			_ = tx.Rollback()
			return utils.SQLError(err)
		}
		var newTags ex.QuestionTags
		for _, t := range tags {
			// do not add existing difficulties
			switch ex.DifficultyTag(t.Tag) {
			case ex.Diff1, ex.Diff2, ex.Diff3:
				continue
			}

			t.IdQuestion = newQuestion.Id
			newTags = append(newTags, t)
		}
		newTags = append(newTags, ex.QuestionTag{IdQuestion: newQuestion.Id, Tag: string(diff)})
		err = updateTags(tx, newTags, newQuestion.Id)
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

// do NOT commit or rollback
func updateTags(tx *sql.Tx, tags ex.QuestionTags, idQuestion int64) error {
	var nbDiff int
	for _, tag := range tags {
		switch ex.DifficultyTag(tag.Tag) {
		case ex.Diff1, ex.Diff2, ex.Diff3:
			nbDiff++
		}
	}
	if nbDiff > 1 {
		return errors.New("Un seul niveau de difficulté est autorisé par question.")
	}

	_, err := ex.DeleteQuestionTagsByIdQuestions(tx, idQuestion)
	if err != nil {
		return err
	}
	err = ex.InsertManyQuestionTags(tx, tags...)
	if err != nil {
		return err
	}
	return nil
}

func (ct *Controller) updateTags(params UpdateTagsIn) error {
	var tags ex.QuestionTags
	for _, tag := range params.Tags {
		// enforce proper tags
		tag = strings.ToUpper(strings.TrimSpace(tag))
		if tag == "" {
			continue
		}

		tags = append(tags, ex.QuestionTag{IdQuestion: params.IdQuestion, Tag: tag})
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

func (ct *Controller) checkParameters(params CheckParametersIn) CheckParametersOut {
	err := params.Parameters.Validate()
	if err != nil {
		return CheckParametersOut{ErrDefinition: err.(ex.ErrParameters)}
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

func (ct *Controller) saveAndPreview(params SaveAndPreviewIn) (SaveAndPreviewOut, error) {
	if err := params.Question.Validate(); err != nil {
		return SaveAndPreviewOut{Error: err.(ex.ErrQuestionInvalid)}, nil
	}

	_, err := params.Question.Update(ct.db)
	if err != nil {
		return SaveAndPreviewOut{}, err
	}

	question := params.Question.Instantiate()

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

type InstantiatedQuestion struct {
	Id       int64
	Question client.Question `dart-extern:"exercices/types.gen.dart"`
	Params   expression.Variables
}

type InstantiateQuestionsOut []InstantiatedQuestion

// InstantiateQuestions loads and instantiates the given questions,
// also returning the paramerters used to do so.
func (ct *Controller) InstantiateQuestions(ids []int64) (InstantiateQuestionsOut, error) {
	questions, err := ex.SelectQuestions(ct.db, ids...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	out := make(InstantiateQuestionsOut, len(ids))
	for index, id := range ids {
		qu := questions[id]
		vars, err := qu.Parameters.ToMap().Instantiate()
		if err != nil {
			return nil, err
		}
		instance, err := qu.InstantiateWith(vars)
		if err != nil {
			return nil, err
		}
		out[index] = InstantiatedQuestion{
			Id:       id,
			Question: instance.ToClient(),
			Params:   vars,
		}
	}

	return out, nil
}

// EvaluateQuestion instantiate the given question with the given parameters,
// and evaluate the given answer.
func (ct *Controller) EvaluateQuestion(id int64, params expression.Variables, answer client.QuestionAnswersIn) (client.QuestionAnswersOut, error) {
	qu, err := ex.SelectQuestion(ct.db, id)
	if err != nil {
		return client.QuestionAnswersOut{}, utils.SQLError(err)
	}

	instance, err := qu.InstantiateWith(params)
	if err != nil {
		return client.QuestionAnswersOut{}, err
	}

	return instance.EvaluateAnswer(answer), nil
}
