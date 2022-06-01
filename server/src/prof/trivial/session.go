package trivial

import (
	"context"
	"fmt"
	"sync"

	tv "github.com/benoitkugler/maths-online/trivial-poursuit"
	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
	"github.com/benoitkugler/maths-online/utils"
)

type stopGame struct {
	ID      tv.GameID
	Restart bool // if false, definitively close the game
}

type createGame struct {
	ID        tv.GameID
	Questions game.QuestionPool
	Options   tv.GameOptions
}

// gameSession is the list of the current active games
// for one teacher
// it is created at the first game launch and either explicitly shut down
// or timed out
type gameSession struct {
	id        SessionID
	idTeacher int64 // optional for demonstration games

	lock sync.Mutex

	// active games
	games map[tv.GameID]*tv.GameController

	// registred players
	playerIDs map[PlayerID]gamePosition

	// stopGameEvents and createGameEvents are used by the teacher to control
	// the current session
	createGameEvents chan createGame // calls createGame()
	stopGameEvents   chan stopGame   // calls stopGame()

	// channel receiving game progress
	monitorSummary chan tv.GameSummary

	// clients to which send the content of `monitor`
	teacherClients map[*teacherClient]bool
}

func newGameSession(id SessionID, idTeacher int64) *gameSession {
	return &gameSession{
		id:               id,
		idTeacher:        idTeacher,
		games:            make(map[tv.GameID]*tv.GameController),
		playerIDs:        make(map[PlayerID]gamePosition),
		createGameEvents: make(chan createGame),
		stopGameEvents:   make(chan stopGame),
		monitorSummary:   make(chan tv.GameSummary),
		teacherClients:   make(map[*teacherClient]bool),
	}
}

// make sure game id are properly sorted when
// created from groups
func (gs *gameSession) newGameID() tv.GameID {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	serial := len(gs.games)
	newID := gameIDFromSerial(gs.id, serial)
	for gs.games[newID] != nil {
		serial++
		newID = gameIDFromSerial(gs.id, serial)
	}
	return newID
}

// make sure game id are properly sorted when
// created from groups
func gameIDFromSerial(sessionID string, serial int) tv.GameID {
	return tv.GameID(sessionID + fmt.Sprintf("%02d", serial+1))
}

// createGame registers and start a new game
func (gs *gameSession) createGame(params createGame) {
	game := tv.NewGameController(params.ID, params.Questions, params.Options, gs.monitorSummary)

	// register the controller...
	gs.lock.Lock()
	gs.games[params.ID] = game
	gs.lock.Unlock()

	// ...and starts it
	go func() {
		review, ok := game.StartLoop()
		if ok { // exploit the review
			gs.exploitReview(review)
		}
		ProgressLogger.Printf("Game %s is done", params.ID)
	}()

	ProgressLogger.Printf("Creating game %s for %d players", params.ID, params.Options.PlayersNumber)
}

// TODO:
func (gs *gameSession) exploitReview(review tv.Review) {
	ProgressLogger.Printf("GAME REVIEW: %v", review)
}

func (gs *gameSession) afterGameEnd(gameID tv.GameID) {
	gs.lock.Lock()
	delete(gs.games, gameID)
	gs.lock.Unlock()

	gs.monitorRemoveGame(gameID)
}

func (gs *gameSession) stopGame(params stopGame) {
	game := gs.games[params.ID]
	if game == nil {
		return
	}

	// copy the current configuration
	create := createGame{
		ID:        game.ID,
		Questions: game.Game.QuestionPool,
		Options:   game.Options,
	}

	game.Terminate <- true

	// restart if needed
	if params.Restart {
		gs.createGame(create)
	} else {
		gs.afterGameEnd(game.ID)
	}
}

// mainLoop blocks until all games are terminated or `ctx` is Done
func (gs *gameSession) mainLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			for _, game := range gs.games {
				game.Terminate <- true
			}
			return
		case cg := <-gs.createGameEvents:
			gs.createGame(cg)
		case sg := <-gs.stopGameEvents:
			gs.stopGame(sg)
			// terminate the session if there is no more games
			if len(gs.games) == 0 {
				return
			}
		case summary := <-gs.monitorSummary:
			for client := range gs.teacherClients {
				client.sendSummary(summary)
			}
		}
	}
}

func (gs *gameSession) monitorRemoveGame(gameID tv.GameID) {
	// notify the monitors
	for client := range gs.teacherClients {
		client.removeGame(gameID)
	}
}

// pass -1 for anonymous sessions
func (ct *Controller) createSession(id SessionID, userID int64) *gameSession {
	// init the session for this teacher
	session := newGameSession(id, userID)

	ct.lock.Lock()
	// register the controller...
	ct.sessions[id] = session
	ct.lock.Unlock()

	ProgressLogger.Printf("Launching session %s", id)
	// ...and start it
	go func() {
		ctx, cancelFunc := context.WithTimeout(context.Background(), sessionTimeout)
		session.mainLoop(ctx)

		cancelFunc()

		// remove the game controller when the game is over
		ct.lock.Lock()
		defer ct.lock.Unlock()
		delete(ct.sessions, id)

		ProgressLogger.Printf("Removing session %s", id)
	}()

	return session
}

// getOrCreateSession returns the session for userID or creates a new one if needed
func (ct *Controller) getOrCreateSession(userID int64) *gameSession {
	if session := ct.getSession(userID); session != nil {
		return session
	}

	ct.lock.Lock()
	sessionID := utils.RandomID(true, 4, func(s string) bool {
		_, taken := ct.sessions[s]
		return taken
	})
	ct.lock.Unlock()

	return ct.createSession(sessionID, userID)
}

// locks and add a new player entry with a sentinel PlayerID
func (gs *gameSession) registerPlayer(gameID tv.GameID) PlayerID {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	playerID := PlayerID(utils.RandomID(false, 20, func(s string) bool {
		_, has := gs.playerIDs[PlayerID(s)]
		return has
	}))
	gs.playerIDs[playerID] = gamePosition{Player: -1, Game: gameID} // register with an initial sentinel value
	return playerID
}

func (gs *gameSession) generateName() string {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	allPlayers := make(map[tv.Player]bool)
	for _, game := range gs.games {
		for p := range game.Summary().Successes {
			allPlayers[p] = true
		}
	}

	nameFromID := func(s string) string {
		return fmt.Sprintf("Joueur %s", s)
	}

	id := utils.RandomID(true, 5, func(s string) bool {
		return allPlayers[tv.Player{Pseudo: nameFromID(s)}]
	})

	return nameFromID(id)
}
