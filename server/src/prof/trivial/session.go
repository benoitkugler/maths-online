package trivial

import (
	"context"
	"fmt"
	"sync"
	"time"

	tv "github.com/benoitkugler/maths-online/server/src/trivial"
	"github.com/benoitkugler/maths-online/server/src/utils"
)

type stopGame struct {
	ID      tv.RoomID
	Restart bool // if false, definitively close the game
}

type createGame struct {
	ID      tv.RoomID
	Options tv.Options
}

// gameSession is the list of the current active games
// for one teacher
// it is created at the first game launch and either explicitly shut down
// or timed out
type gameSession struct {
	id        SessionID
	idTeacher uID // optional for demonstration games

	lock sync.Mutex

	// active games
	games map[tv.RoomID]*tv.Room

	// map registred players to their game room
	playerIDs map[tv.PlayerID]tv.RoomID

	// stopGameEvents and createGameEvents are used by the teacher to control
	// the current session
	createGameEvents   chan createGame // calls createGame()
	stopGameEvents     chan stopGame   // calls stopGame()
	afterGameEndEvents chan tv.RoomID  // calls afterGameEnd()
}

func newGameSession(id SessionID, idTeacher uID) *gameSession {
	return &gameSession{
		id:                 id,
		idTeacher:          idTeacher,
		games:              make(map[tv.RoomID]*tv.Room),
		playerIDs:          make(map[tv.PlayerID]tv.RoomID),
		createGameEvents:   make(chan createGame),
		stopGameEvents:     make(chan stopGame),
		afterGameEndEvents: make(chan tv.RoomID),
	}
}

// make sure game id are properly sorted when
// created from groups
func (gs *gameSession) newGameID() tv.RoomID {
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
func gameIDFromSerial(sessionID string, serial int) tv.RoomID {
	return tv.RoomID(sessionID + fmt.Sprintf("%02d", serial+1))
}

// createGame registers and start a new game
func (gs *gameSession) createGame(params createGame) {
	game := tv.NewRoom(params.ID, params.Options)

	// register the controller...
	gs.lock.Lock()
	gs.games[params.ID] = game
	gs.lock.Unlock()

	// ...and starts it
	go func() {
		replay, naturalEnding := game.Listen()
		if naturalEnding { // exploit the review
			gs.exploitReplay(replay)
		}
		ProgressLogger.Printf("Game %s is done, cleaning up...", params.ID)

		// if the game was terminated explicitely (by stopGame),
		// do not perform the cleanup to avoid interfering with restart
		if naturalEnding {
			gs.afterGameEndEvents <- params.ID
		}
	}()

	ProgressLogger.Printf("Creating game %s (launch: %s)", params.ID, params.Options.Launch)
}

// TODO:
func (gs *gameSession) exploitReplay(review tv.Replay) {
	ProgressLogger.Printf("GAME REPLAY: %v", review)
}

func (gs *gameSession) afterGameEnd(gameID tv.RoomID) {
	gs.lock.Lock()
	defer gs.lock.Unlock()
	delete(gs.games, gameID)
}

func (gs *gameSession) startGame(gameID tv.RoomID) error {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	game, ok := gs.games[gameID]
	if !ok {
		return fmt.Errorf("internal error: no game with ID %s", gameID)
	}

	return game.StartGame()
}

func (gs *gameSession) stopGame(params stopGame) {
	game := gs.games[params.ID]
	if game == nil {
		return
	}

	// copy the current configuration
	create := createGame{
		ID:      game.ID,
		Options: game.Options(),
	}

	game.Terminate <- true

	// restart if needed
	if params.Restart {
		time.Sleep(time.Millisecond)
		gs.createGame(create)
	} else { // cleanup
		gs.afterGameEnd(params.ID)
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
			// terminate the session if there is no more games and no restart it asked for
			if !sg.Restart && len(gs.games) == 0 {
				return
			}
		case gameID := <-gs.afterGameEndEvents:
			gs.afterGameEnd(gameID)
		}
	}
}

// pass -1 as userID for anonymous sessions
func (ct *Controller) createSession(id SessionID, userID uID) *gameSession {
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
func (ct *Controller) getOrCreateSession(userID uID) *gameSession {
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

// locks and add a new player in the map player -> games
func (gs *gameSession) registerPlayer(gameID tv.RoomID) tv.PlayerID {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	playerID := tv.PlayerID(utils.RandomID(false, 20, func(s string) bool {
		_, has := gs.playerIDs[tv.PlayerID(s)]
		return has
	}))
	gs.playerIDs[playerID] = gameID

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
