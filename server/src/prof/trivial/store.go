package trivial

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	tv "github.com/benoitkugler/maths-online/server/src/trivial"
	"github.com/benoitkugler/maths-online/server/src/utils"
)

// gameStore is the global storage for game rooms,
// safe for concurrent use
type gameStore struct {
	// demoPin is used to create testing games on the fly
	demoPin string

	lock sync.Mutex

	games map[gameID]*tv.Room

	// additional map use to link teacherCode with teacher DB id,
	// to keep sync with games
	teacherSessions map[sessionID]teacher.IdTeacher

	// map registred players to their game room
	playerIDs map[tv.PlayerID]gameID
}

// initialize the maps
func newGameStore(demoPin string) gameStore {
	return gameStore{
		games:           make(map[gameID]*tv.Room),
		teacherSessions: make(map[sessionID]teacher.IdTeacher),
		playerIDs:       make(map[tv.PlayerID]gameID),
		demoPin:         demoPin,
	}
}

// getSession returns the games associated to the user, or nil
// DO NOT LOCK
func (gs *gameStore) getSession(sessionID sessionID) (out []*tv.Room) {
	for id, room := range gs.games {
		id, ok := id.(teacherCode)
		if ok && id.sessionID == sessionID {
			out = append(out, room)
		}
	}

	return out
}

// return empty if not found
func (ct *gameStore) getSessionID(userID uID) sessionID {
	ct.lock.Lock()
	defer ct.lock.Unlock()

	for k, v := range ct.teacherSessions {
		if v == userID {
			return k
		}
	}
	return ""
}

// getOrCreateSession returns the session for userID or creates a new one if needed
func (ct *gameStore) getOrCreateSession(userID uID) sessionID {
	if sID := ct.getSessionID(userID); sID != "" {
		return sID
	}

	ct.lock.Lock()
	defer ct.lock.Unlock()

	// create a new session ID...
	sID := utils.RandomID(true, 4, func(s string) bool {
		_, taken := ct.teacherSessions[s]
		return taken
	})
	// ... and register it
	ct.teacherSessions[sID] = userID

	return sID
}

// locks, and make sure game id are properly sorted when
// created from groups
func (gs *gameStore) newTeacherGameID(session sessionID) teacherCode {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	serial := 0
	gameID := func() teacherCode { return teacherCode{session, fmt.Sprintf("%02d", serial+1)} }
	newID := gameID()
	for gs.games[newID] != nil {
		serial++
		newID = gameID()
	}
	return newID
}

// locks and generate a new id
func (gs *gameStore) newSelfaccessGameID() selfaccessCode {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	id := utils.RandomID(true, 5, func(s string) bool {
		_, taken := gs.games[selfaccessCode(s)]
		return taken
	})
	return selfaccessCode(id)
}

type createGame struct {
	ID      gameID
	Options tv.Options
}

// createGame locks, creates, registers and starts the eveng loop of new game
func (gs *gameStore) createGame(params createGame) {
	game := tv.NewRoom(tv.RoomID(params.ID.String()), params.Options)

	// register the controller...
	gs.lock.Lock()
	gs.games[params.ID] = game
	gs.lock.Unlock()

	// ...and starts it
	ctx, cancelFunc := context.WithTimeout(context.Background(), gameTimeout)
	go func() {
		replay, naturalEnding := game.Listen(ctx)
		cancelFunc()
		if naturalEnding { // exploit the review
			gs.exploitReplay(replay)
		}
		ProgressLogger.Printf("Game %s is done, cleaning up...", params.ID)

		// if the game was terminated explicitely (by stopGame),
		// do not perform the cleanup to avoid interfering with restart
		if naturalEnding {
			gs.afterGameEnd(params.ID)
		}
	}()

	ProgressLogger.Printf("Creating game %s (%T, launch: %s)", params.ID, params.ID, params.Options.Launch)
}

// TODO:
func (gs *gameStore) exploitReplay(review tv.Replay) {
	ProgressLogger.Printf("GAME REPLAY: %v", review)
}

// cleanup the ressource associated with the game
func (gs *gameStore) afterGameEnd(gameID gameID) {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	delete(gs.games, gameID)

	// cleanup session map if needed
	if tc, ok := gameID.(teacherCode); ok {
		if len(gs.getSession(tc.sessionID)) == 0 { // no more session
			delete(gs.teacherSessions, tc.sessionID)
			ProgressLogger.Printf("Removing session %s", tc.sessionID)
		}
	}
}

// locks and start the given game
func (gs *gameStore) startGame(gameID gameID) error {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	game, ok := gs.games[gameID]
	if !ok {
		return fmt.Errorf("internal error: no game with ID %s", gameID)
	}

	return game.StartGame()
}

func (gs *gameStore) stopGame(id gameID, restart bool) {
	gs.lock.Lock()
	game := gs.games[id]
	gs.lock.Unlock()
	if game == nil {
		return
	}

	// copy the current configuration
	create := createGame{
		ID:      id,
		Options: game.Options(),
	}

	game.Terminate <- true
	// restart if needed
	if restart {
		time.Sleep(time.Millisecond)
		gs.createGame(create)
	} else { // cleanup
		gs.afterGameEnd(id)
	}
}

// locks and add a new player in the map player -> games
func (gs *gameStore) registerPlayer(gameID gameID) tv.PlayerID {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	playerID := tv.PlayerID(utils.RandomID(false, 20, func(s string) bool {
		_, has := gs.playerIDs[tv.PlayerID(s)]
		return has
	}))
	gs.playerIDs[playerID] = gameID

	return playerID
}

func (gs *gameStore) generateName() string {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	allPlayers := make(map[tv.Player]bool)
	for _, game := range gs.games {
		for p := range game.Summary().Successes {
			allPlayers[p] = true
		}
	}

	nameFromID := func(s string) string { return fmt.Sprintf("Joueur %s", s) }

	id := utils.RandomID(true, 6, func(s string) bool {
		return allPlayers[tv.Player{Pseudo: nameFromID(s)}]
	})

	return nameFromID(id)
}
