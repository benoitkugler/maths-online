// Package trivial implements a game controller
// for one group of players, handling concurrency.
// It is meant to be used by a server with websocket connections
// for the players, but is actually network agnostic.
package trivial

import (
	"sync"
	"time"
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

// Options is the configuration of one game.
// All fields are required.
type Options struct {
	// QuestionPool is the list of the question
	// being asked, for each category
	Questions QuestionPool
	// PlayersNumber iss the required number for the game, used to trigger a game start
	PlayersNumber int
	// QuestionTimeout is the time limit for one question
	QuestionTimeout time.Duration

	ShowDecrassage bool
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
}

// playerConn stores a player profile and the underlying connection,
// which is nil for inactive (disconnected) ones
type playerConn struct {
	pl      Player
	conn    Connection
	advance playerAdvance
}

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
	Event chan ClientEvent

	// protect external access to the game state and players
	lock sync.Mutex

	// game state, whose accesses are protected
	// by the channels and `lock`
	game game

	// currentPlayers stores the actual players in the game,
	// including the inactive (disconnected) ones, for which
	// `Connection` is nil.
	// we always have len(currentPlayers) <= expectedPlayers
	currentPlayers map[string]*playerConn
}

func NewRoom(ID RoomID, options Options) *Room {
	return &Room{
		ID:             ID,
		Terminate:      make(chan bool),
		Leave:          make(chan PlayerID),
		Event:          make(chan ClientEvent),
		game:           newGame(options),
		currentPlayers: make(map[PlayerID]*playerConn),
	}
}

// Options returns the (readonly) configuration used by
// the game.
func (r *Room) Options() Options { return r.game.options }

// Replay exposes some information to be persisted
// after the game end, such as the successes of the players
type Replay struct {
	QuestionHistory map[Player]QuestionReview
	ID              RoomID
}

// return the current game review, without locking
func (r *Room) review() Replay {
	out := Replay{
		ID:              r.ID,
		QuestionHistory: make(map[Player]QuestionReview),
	}

	for _, pl := range r.currentPlayers {
		if pl.conn == nil { // player not connected anymore
			continue
		}
		out.QuestionHistory[pl.pl] = pl.advance.review
	}
	return out
}

// Summary provides an high level overview of the game,
// and may be emitted back to the teacher monitor.
type Summary struct {
	PlayerTurn *Player // nil before game start
	// Successes does not contains disconnected players
	Successes map[Player]Success
	ID        RoomID
	RoomSize  int // Number of player expected
}

// Summary locks and returns the current game summary.
func (r *Room) Summary() Summary {
	r.lock.Lock()
	defer r.lock.Unlock()

	successes := make(map[Player]Success)
	for _, v := range r.currentPlayers {
		successes[v.pl] = v.advance.success
	}
	out := Summary{
		ID:        r.ID,
		Successes: successes,
		RoomSize:  r.game.options.PlayersNumber,
	}

	if se := r.game.playerTurn; se != "" {
		if pl, has := r.currentPlayers[se]; has {
			out.PlayerTurn = &pl.pl
		}
	}

	return out
}
