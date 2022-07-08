package trivial

import (
	"errors"
	"log"
	"os"

	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
)

var (
	WarningLogger  = log.New(os.Stdout, "tv-game:ERROR: ", 0)
	ProgressLogger = log.New(os.Stdout, "tv-game:INFO : ", 0)
)

func (pc playerConn) send(events game.StateUpdate) {
	err := pc.conn.WriteJSON(events)
	if err != nil {
		WarningLogger.Printf("Sending to client %s failed: %s", pc.pl.ID, err)
	}
}

// TODO: probably use r.game.State
func (r *Room) broadcastEvents(events game.StateUpdate) {
	ProgressLogger.Printf("Game %s : broadcasting...", r.ID)

	for _, pc := range r.currentPlayers {
		if pc.conn == nil { // ignore disconnected players
			continue
		}
		pc.send(events)
	}
}

// Listen starts the main game loop, listening
// on the game channels.
// It returns `true` when the game is over, or `false` when it
// is manually terminated.
// Care should be taken to make sure no more events are send on the
// channels when this method has returned.
func (r *Room) Listen() (replay Replay, naturalEnding bool) {
	for {
		select {
		case <-r.Terminate:
			r.onTerminate()
			return Replay{}, false
		case client := <-r.Leave:
			r.onLeave(client)
		case <-r.game.QuestionTimeout.C:
			r.onQuestionTimeout()
		case message := <-r.Event:
			isGameOver := r.onEvent(message)
			if isGameOver {
				ProgressLogger.Printf("Game %s is over: exiting game loop.", r.ID)

				return r.review(), true
			}
		}
	}
}

func (r *Room) onTerminate() {
	r.lock.Lock()
	defer r.lock.Unlock()

	ProgressLogger.Printf("Game %s : terminating game...", r.ID)

	r.broadcastEvents(game.StateUpdate{
		Events: game.Events{
			game.GameTerminated{},
		},
		State: r.game.GameState,
	})
}

// ErrGameStarted is returned from `Join` when the game
// has already started
var ErrGameStarted = errors.New("game started")

// Join should be used on a new connection, on one the of the following cases:
//	- totally fresh connection
//	- reconnection in an already started game
// It is safe for concurrent uses.
// It returns an error for instance if the game has already started.
func (r *Room) Join(player Player, connection Connection) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	// check if it is a reconnection
	_, isKnownPlayer := r.currentPlayers[player.ID]
	if !isKnownPlayer {
		// fresh connection : only allow it on waiting games
		if r.game.HasStarted() {
			return ErrGameStarted
		}

		ProgressLogger.Printf("Game %s : adding new player %s...", r.ID, player.ID)

		// TODO: merge AddPlayer and StartGame
		serial, event := r.game.AddPlayer(player.Pseudo)
		// register the serial for futur uses ...
		player.serial = serial
		// ... and register the player
		pc := playerConn{pl: player, conn: connection}
		r.currentPlayers[player.ID] = pc

		// only notifie the player who joined ...
		pc.send(game.StateUpdate{
			Events: game.Events{game.PlayerJoin{Player: serial}},
			State:  r.game.GameState,
		})

		// ... and check if the new player triggers a game start
		var events game.StateUpdate
		if r.game.NumberPlayers(true) >= r.Options().PlayersNumber {
			ProgressLogger.Printf("Game %s : starting...", r.ID)
			// TODO: refactor StartGame
			events = r.game.StartGame()
		} else { // update the lobby
			events = game.StateUpdate{
				Events: game.Events{event},
				State:  r.game.GameState,
			}
		}
		r.broadcastEvents(events)
	} else { // reconnection
		ProgressLogger.Printf("Game %s : reconnecting player %s...", r.ID, player.ID)

		pc := r.currentPlayers[player.ID]
		pc.conn = connection
		pc.pl.Pseudo = player.Pseudo
		r.currentPlayers[player.ID] = pc // register the new client connection

		// TODO: refactor ReconnectPlayer
		event := r.game.ReconnectPlayer(pc.pl.serial, pc.pl.Pseudo)

		r.broadcastEvents(game.StateUpdate{
			Events: game.Events{event},
			State:  r.game.GameState,
		})
	}

	return nil
}

func (r *Room) onLeave(playerID PlayerID) {
	r.lock.Lock()
	defer r.lock.Unlock()

	pc, in := r.currentPlayers[playerID]
	// defensive check
	if !in {
		return
	}

	ProgressLogger.Printf("Game %s : removing player %d (ID: %s)...", r.ID, pc.pl.serial, pc.pl.ID)

	// TODO: refactor RemovePlayer
	update := r.game.RemovePlayer(pc.pl.serial)
	delete(r.currentPlayers, playerID)
	r.broadcastEvents(update)
}

func (r *Room) onEvent(event game.ClientEvent) (isGameOver bool) {
	r.lock.Lock()
	defer r.lock.Unlock()

	ProgressLogger.Printf("Game %s : handleClientEvent...", r.ID)

	events, isGameOver, err := r.game.HandleClientEvent(event)
	if err != nil { // malicious client: ignore the query
		WarningLogger.Println(err)
		return false
	}

	if events != nil {
		r.broadcastEvents(*events)
	}

	return isGameOver
}

func (r *Room) onQuestionTimeout() {
	r.lock.Lock()
	defer r.lock.Unlock()

	ProgressLogger.Printf("Game %s : questionTimeoutAction...", r.ID)

	events := r.game.QuestionTimeoutAction()
	if events != nil {
		r.broadcastEvents(*events)
	}
}
