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

func (pc playerConn) send(events StateUpdate) {
	err := pc.conn.WriteJSON(events)
	if err != nil {
		WarningLogger.Printf("Sending to client %s failed: %s", pc.pl.ID, err)
	}
}

func (r *Room) broadcastEvents(events Events) {
	ProgressLogger.Printf("Game %s : broadcasting...", r.ID)

	state := r.state()
	for _, pc := range r.currentPlayers {
		if pc.conn == nil { // ignore disconnected players
			continue
		}
		pc.send(StateUpdate{Events: events, State: state})
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
		case <-r.game.questionTimer.C:
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

	r.broadcastEvents(Events{
		GameTerminated{},
	},
	)
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
		if r.game.hasStarted() {
			return ErrGameStarted
		}

		ProgressLogger.Printf("Game %s : adding new player %s...", r.ID, player.ID)

		serial := r.newSerial()
		// register the serial for futur uses ...
		player.serial = serial
		// ... and register the player
		pc := playerConn{pl: player, conn: connection, advance: playerAdvance{} /* zero value is enough */}
		r.currentPlayers[player.ID] = &pc

		// notify the player who joined to show the lobby ...
		pc.send(StateUpdate{Events: Events{PlayerJoin{Player: serial}}, State: r.state()})

		// ... also notify the other players ...
		r.broadcastEvents(Events{LobbyUpdate{
			Player:     serial,
			PlayerName: player.Pseudo,
			IsJoining:  true,
			Names:      r.playerPseudos(),
		}})

		// ... and check if the new player triggers a game start, after a brief pause
		time.Sleep(time.Second)
		if events := r.tryStartGame(); len(events) != 0 {
			r.broadcastEvents(events)
		}
	} else { // reconnection
		ProgressLogger.Printf("Game %s : reconnecting player %s...", r.ID, player.ID)

		pc := r.currentPlayers[player.ID]
		pc.conn = connection // use the new client connection
		pc.pl.Pseudo = player.Pseudo

		event := PlayerReconnected{
			PlayerID:   pc.pl.serial,
			PlayerName: player.Pseudo,
		}
		r.broadcastEvents(Events{event})
	}

	return nil
}

func (r *Room) onLeave(playerID PlayerID) {
	r.lock.Lock()
	defer r.lock.Unlock()

	pc, in := r.currentPlayers[playerID]
	if !in { // defensive check
		return
	}

	ProgressLogger.Printf("Game %s : removing player %d (ID: %s)...", r.ID, pc.pl.serial, pc.pl.ID)

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

	player, ok := r.currentPlayers[event.Player]
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
