// Package editor provides functionnality for a frontend
// to edit and preview math questions.
package editor

import (
	"context"
	"fmt"
	"math/rand"
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

	sessions map[string]*loopbackController
}

func NewController() *Controller {
	return &Controller{
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

// TODO actually validate and save
func (ct *Controller) saveAndPreview(sessionID string) error {
	index := rand.Intn(len(exercice.PredefinedQuestions))
	question := exercice.PredefinedQuestions[index].ToClient()

	ct.lock.Lock()
	defer ct.lock.Unlock()

	loopback, ok := ct.sessions[sessionID]
	if !ok {
		return fmt.Errorf("invalid session ID %s", sessionID)
	}

	loopback.broadcast <- question
	return nil
}
