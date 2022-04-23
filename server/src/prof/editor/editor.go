// Package editor provides functionnality for a frontend
// to edit and preview math questions.
package editor

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	ex "github.com/benoitkugler/maths-online/maths/exercice"
	"github.com/benoitkugler/maths-online/utils"
)

const sessionTimeout = 12 * time.Hour

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

type QuestionHeader struct {
	Id    int64
	Title string
	Tags  []string
}

func (ct *Controller) searchQuestions(query ListQuestionsIn) (out []QuestionHeader, err error) {
	const pagination = 30

	var questions ex.Questions
	if len(query.Tags) != 0 {
		questions, err = ex.SelectQuestionByTags(ct.db, query.Tags...)
	} else {
		questions, err = ex.SelectAllQuestions(ct.db)
	}
	if err != nil {
		return nil, err
	}

	queryTitle := strings.TrimSpace(strings.ToLower(query.TitleQuery))
	var ids ex.IDs
	for _, question := range questions {
		thisTitle := strings.TrimSpace(strings.ToLower(question.Title))
		if strings.Contains(thisTitle, queryTitle) {
			out = append(out, QuestionHeader{Id: question.Id, Title: question.Title})
			ids = append(ids, question.Id)
		}
	}

	tags, err := ex.SelectQuestionTagsByIdQuestions(ct.db, ids...)
	if err != nil {
		return nil, err
	}
	tagsMap := tags.ByIdQuestion()
	for i, question := range out {
		for _, tag := range tagsMap[question.Id] {
			out[i].Tags = append(out[i].Tags, tag.Tag)
		}
	}

	if len(out) > pagination {
		out = out[:pagination]
	}

	return out, nil
}

// duplicateWithDifficulty creates new questions with the same title
// and content as the given question, but with difficulty levels
func (ct *Controller) duplicateWithDifficulty(idQuestion int64) error {
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
	crible := tags.Crible()
	var currentDifficulty string
	var newDifficulties [2]ex.DifficultyTag
	if d := string(ex.Diff1); crible[d] {
		currentDifficulty = d
		newDifficulties = [2]ex.DifficultyTag{ex.Diff2, ex.Diff3}
	} else if d := string(ex.Diff2); crible[d] {
		currentDifficulty = d
		newDifficulties = [2]ex.DifficultyTag{ex.Diff1, ex.Diff3}
	} else if d := string(ex.Diff3); crible[d] {
		currentDifficulty = d
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
