// Package trivial implements a game controller
// for one group of players, handling concurrency.
// It is meant to be used by a server with websocket connections
// for the players, but is actually network agnostic.
package trivial

import (
	"sync"
	"time"

	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
)

// Connection abstracts away the network technology used to
// communicate back to the client
// `WriteJSON` will be called from the event goroutine started
// by `Room.Listen`, so that implementations should support one concurrent write,
// as it is the case for gorilla/websocket.Connection
type Connection interface {
	WriteJSON(v interface{}) error
}

// RoomID is the public, full identifier of a game room,
// usually of the form <sessionID><gameID> (excepted for demonstration games).
// It is not used internally by this package but associated
// to every room object, and will be used by consuming packages.
type RoomID string

// Options is the configuration of one game
type Options struct {
	Questions       game.QuestionPool
	PlayersNumber   int
	QuestionTimeout time.Duration
	ShowDecrassage  bool
}

// PlayerID is a unique identifier of each player,
// usually generated at the first connection.
// It is used to handle reconnection and external monitoring
// of the players
type PlayerID = string

// Player represents one player.
type Player struct {
	ID PlayerID

	// Pseudo is the display name of each player,
	// which may change during a game (upon deconnection/reconnection)
	Pseudo string

	// serial is the number of the player during the game,
	// used by the client and to handle rotation on new turns
	serial game.PlayerSerial
}

// playerConn stores a player meta and the underlying connection,
// which is nil for inactive (disconnected) ones
type playerConn struct {
	pl   Player
	conn Connection
}

// Phase identifies the current phase of the game
type Phase uint8

const (
	PWaiting  Phase = iota // not started yet
	PThrowing              // start of turn, waiting for dice throw
	PMoving                // dice was thrown, waiting for player move
	PQuestion              // question is being answered
	PResult                // players are consulting answer results
)

// Room is the game host, and the main entry point
// of a game.
// All exported methods are safe for concurrent use; events are
// send on the exposed channel fields.
type Room struct {
	// ID is the readonly ID for this game.
	ID RoomID

	// Terminate is used to cleanly exit the game,
	// noticing clients and exiting the main goroutine.
	// It is however not considered as a normal exit,
	// so that the `Replay` is not emitted.
	Terminate chan bool

	// Leave is used when the player leave the game
	// (either on purpose or when its connection breaks)
	// For started games, the player is only set inactive,
	// whereas for waiting (in lobby) games, the player is totaly removed.
	Leave chan PlayerID

	// Event is used when a client send an event
	Event chan game.ClientEvent

	// protect external access to the game state
	lock sync.Mutex

	// actual game logic, whose accesses are protected
	// by the channels and `lock`
	game game.Game

	// used for instance to trigger the correct event
	// when a player disconnect
	phase Phase

	// required number for the game, used to trigger a game start
	expectedPlayers int

	// currentPlayers stores the actual players in the game,
	// including the inactive (disconnected) ones, for which
	// `Connection` is nil.
	// we always have len(currentPlayers) <= expectedPlayers
	currentPlayers map[string]playerConn
}

func NewRoom(ID RoomID, options Options) *Room {
	return &Room{
		ID:              ID,
		Terminate:       make(chan bool),
		Leave:           make(chan PlayerID),
		Event:           make(chan game.ClientEvent),
		game:            game.NewGame(options.QuestionTimeout, options.ShowDecrassage, options.Questions),
		expectedPlayers: options.PlayersNumber,
		currentPlayers:  make(map[PlayerID]playerConn),
	}
}

// Options returns the (readonly) configuration used by
// the game.
func (r *Room) Options() Options {
	return Options{
		Questions:       r.game.QuestionPool,
		PlayersNumber:   r.expectedPlayers,
		QuestionTimeout: r.game.QuestionDurationLimit,
		ShowDecrassage:  r.game.ShowDecrassage,
	}
}

// Replay exposes some information to be persisted
// after the game end, such as the successes of the players
type Replay struct {
	QuestionHistory map[Player]game.QuestionReview
	ID              RoomID
}

// return the current game review, without locking
func (r *Room) review() Replay {
	out := Replay{
		ID:              r.ID,
		QuestionHistory: make(map[Player]game.QuestionReview),
	}

	players := r.reversePlayers()
	for k, v := range r.game.Players {
		p, has := players[k]
		if !has { // player not connected anymore
			continue
		}
		out.QuestionHistory[p] = v.Review
	}
	return out
}

// Summary provides an high level overview of the game,
// and may be emitted back to the teacher monitor.
type Summary struct {
	PlayerTurn *Player // nil before game start
	// Successes does not contains disconnected players
	Successes map[Player]game.Success
	ID        RoomID
	RoomSize  int // Number of player expected
}

// does not include inactive players
func (r *Room) reversePlayers() map[game.PlayerSerial]Player {
	players := make(map[game.PlayerSerial]Player)
	for _, pc := range r.currentPlayers {
		if pc.conn == nil {
			continue
		}
		players[pc.pl.serial] = pc.pl
	}
	return players
}

// Summary locks and returns the current game summary.
func (r *Room) Summary() Summary {
	r.lock.Lock()
	defer r.lock.Unlock()

	state := r.game.GameState
	players := r.reversePlayers()

	successes := make(map[Player]game.Success)
	for k, v := range state.Players {
		client := players[k]
		successes[client] = v.Success
	}
	out := Summary{
		ID:        r.ID,
		Successes: successes,
		RoomSize:  r.expectedPlayers,
	}
	if id := state.Player; id != -1 {
		if pl, has := players[id]; has {
			out.PlayerTurn = &pl
		}
	}

	return out
}
