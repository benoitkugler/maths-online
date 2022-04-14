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

	"github.com/benoitkugler/maths-online/maths/exercice"
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

	var questions exercice.Questions
	if len(query.Tags) != 0 {
		questions, err = exercice.SelectQuestionByTags(ct.db, query.Tags...)
	} else {
		questions, err = exercice.SelectAllQuestions(ct.db)
	}
	if err != nil {
		return nil, err
	}

	queryTitle := strings.TrimSpace(strings.ToLower(query.TitleQuery))
	var ids exercice.IDs
	for _, question := range questions {
		thisTitle := strings.TrimSpace(strings.ToLower(question.Title))
		if strings.Contains(thisTitle, queryTitle) {
			out = append(out, QuestionHeader{Id: question.Id, Title: question.Title})
			ids = append(ids, question.Id)
		}
	}

	tags, err := exercice.SelectQuestionTagsByIdQuestions(ct.db, ids...)
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

func (ct *Controller) updateTags(params UpdateTagsIn) error {
	tx, err := ct.db.Begin()
	if err != nil {
		return err
	}
	_, err = exercice.DeleteQuestionTagsByIdQuestions(tx, params.IdQuestion)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	var tags exercice.QuestionTags
	for _, tag := range params.Tags {
		tags = append(tags, exercice.QuestionTag{IdQuestion: params.IdQuestion, Tag: tag})
	}
	err = exercice.InsertManyQuestionTags(tx, tags...)
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
		return CheckParametersOut{ErrDefinition: err.(exercice.ErrParameters)}
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

	loopback.broadcast <- LoopbackState{IsPaused: true}
	return nil
}

func (ct *Controller) saveAndPreview(params SaveAndPreviewIn) error {
	// TODO validation step before saving
	_, err := params.Question.Update(ct.db)
	if err != nil {
		return err
	}

	question := params.Question.Instantiate().ToClient()

	ct.lock.Lock()
	defer ct.lock.Unlock()

	loopback, ok := ct.sessions[params.SessionID]
	if !ok {
		return fmt.Errorf("invalid session ID %s", params.SessionID)
	}

	loopback.broadcast <- LoopbackState{Question: question}
	return nil
}
