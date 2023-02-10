package trivial

import (
	"errors"
	"log"
	"os"
	"time"
)

var (
	WarningLogger  = log.New(os.Stdout, "tv-game:ERROR: ", 0)
	ProgressLogger = log.New(os.Stdout, "tv-game:INFO : ", 0)
)

var gameStartDelay = time.Second // reduced in tests to avoid latency

func (pc playerConn) send(events StateUpdate) {
	err := pc.conn.WriteJSON(events)
	if err != nil {
		WarningLogger.Printf("Sending to client %s failed: %s", pc.pl.ID, err)
	}
}

func (r *Room) broadcastEvents(events Events) {
	ProgressLogger.Printf("Game %s : broadcasting...", r.ID)

	state := r.state()
	for _, pc := range r.players {
		if pc.conn == nil { // ignore disconnected players
			continue
		}
		pc.send(StateUpdate{Events: events, State: state})
	}
}

// Listen starts the main game loop, listening
// on the game channels.
// Note that it does start the game itself : the game state is
// initially in the lobby.
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
		case <-r.game.questionTimer.C:
			r.onQuestionTimeout()
		case message := <-r.Event:
			isGameOver := r.onEvent(message)
			if isGameOver {
				ProgressLogger.Printf("Game %s is over: exiting game loop.", r.ID)

				return r.replay(), true
			}
		}
	}
}

func (r *Room) onTerminate() {
	r.lock.Lock()
	defer r.lock.Unlock()

	ProgressLogger.Printf("Game %s : terminating game...", r.ID)

	r.broadcastEvents(Events{
		GameTerminated{},
	},
	)
}

// ErrGameStarted is returned from `Join` when the game
// has already started
var ErrGameStarted = errors.New("game already started")

// Join should be used on a new connection, on one the of the following cases:
//   - totally fresh connection
//   - reconnection in an already started game
//
// It is safe for concurrent uses.
// It returns an error for instance if the game has already started.
func (r *Room) Join(player Player, connection Connection) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	// check if it is a reconnection
	_, isKnownPlayer := r.players[player.ID]
	if !isKnownPlayer {
		// fresh connection : only allow it on waiting games
		if r.game.hasStarted() {
			return ErrGameStarted
		}

		ProgressLogger.Printf("Game %s : adding new player %s...", r.ID, player.ID)

		// register the player
		pc := playerConn{pl: player, conn: connection, advance: playerAdvance{} /* zero value is enough */}
		r.players[player.ID] = &pc

		// notify the player who joined to show the lobby ...
		pc.send(StateUpdate{Events: Events{PlayerJoin{Player: player.ID}}, State: r.state()})

		// ... also notify the other players ...
		r.broadcastEvents(Events{LobbyUpdate{
			ID:            player.ID,
			Pseudo:        player.Pseudo,
			IsJoining:     true,
			PlayerPseudos: r.playerPseudos(),
		}})

		// ... and check if the new player triggers a game start, after a brief pause
		if events := r.tryStartGame(); len(events) != 0 {
			time.Sleep(gameStartDelay)
			r.broadcastEvents(events)
		}
	} else { // reconnection
		ProgressLogger.Printf("Game %s : reconnecting player %s...", r.ID, player.ID)

		events := r.reconnectPlayer(player, connection)

		r.broadcastEvents(events)
	}

	return nil
}

// StartGame launches the game, returning an error if the game
// is in auto mode, or if no players are in the game yet.
//
// It is safe for concurrent use.
func (r *Room) StartGame() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if !r.game.options.Launch.Manual {
		return errors.New("Une partie en lancement automatique ne peut pas être démarrée manuellement.")
	}

	if len(r.players) == 0 {
		return errors.New("Aucun joueur n'est présent dans la partie.")
	}

	events := r.startGame()
	r.broadcastEvents(events)

	return nil
}

func (r *Room) onLeave(playerID PlayerID) {
	r.lock.Lock()
	defer r.lock.Unlock()

	pc, in := r.players[playerID]
	if !in { // defensive check
		return
	}

	ProgressLogger.Printf("Game %s : removing player %s...", r.ID, pc.pl.ID)

	update := r.removePlayer(pc.pl)
	r.broadcastEvents(update)
}

// ClientEvent is a game event send by a
// known player
type ClientEvent struct {
	Event  ClientEventITF // as send by the client
	Player PlayerID       // filled by the server
}

func (r *Room) onEvent(event ClientEvent) (isGameOver bool) {
	r.lock.Lock()
	defer r.lock.Unlock()

	ProgressLogger.Printf("Game %s : handling client event (%T)...", r.ID, event.Event)

	player, ok := r.players[event.Player]
	if !ok { // defensive check
		return
	}

	events, isGameOver, err := r.handleClientEvent(event.Event, player.pl)
	if err != nil { // malicious client: ignore the query
		WarningLogger.Println(err)
		return false
	}

	if len(events) != 0 { // do not update the state for nothing
		r.broadcastEvents(events)
	}

	return isGameOver
}

func (r *Room) onQuestionTimeout() {
	r.lock.Lock()
	defer r.lock.Unlock()

	ProgressLogger.Printf("Game %s : questionTimeoutAction...", r.ID)

	events := r.tryEndQuestion(true)
	r.broadcastEvents(events)
}
