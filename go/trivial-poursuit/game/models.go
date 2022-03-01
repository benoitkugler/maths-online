package game

import (
	"encoding/json"
	"math/rand"
)

// interaction with the client

//go:generate ../../../../structgen/structgen -source=models.go -mode=dart:../../../eleve/lib/trivialpoursuit/events.gen.dart -mode=itfs-json:gen.go

// GameState represents an on-going game.
type GameState struct {
	Successes map[PlayerID]*success // per-player advance
	PawnTile  int                   // position of the pawn
	Player    int                   // the player currently playing (choosing where to move)
}

// GameEvents describes a list of events yielding
// a game state.
// Clients should animate the returned events and
// update the state.
type GameEvents struct {
	Events events
	State  GameState
}

// IsEmpty returns `true` if no events are actually emitted.
func (ge GameEvents) IsEmpty() bool { return len(ge.Events) == 0 }

type events []gameEvent

func (evs events) MarshalJSON() ([]byte, error) {
	tmp := make([]gameEventWrapper, len(evs))
	for i, v := range evs {
		tmp[i] = gameEventWrapper{Data: v}
	}
	return json.Marshal(tmp)
}

func (evs *events) UnmarshalJSON(data []byte) error {
	var tmp []gameEventWrapper
	err := json.Unmarshal(data, &tmp)
	*evs = make(events, len(tmp))
	for i, v := range tmp {
		(*evs)[i] = v.Data
	}
	return err
}

// gameEvent is an action (created by the server) advancing the game
// or requiring to update the UI
type gameEvent interface {
	isGameEvent()
}

func (gameStart) isGameEvent()          {}
func (playerLeft) isGameEvent()         {}
func (playerTurn) isGameEvent()         {}
func (diceThrow) isGameEvent()          {}
func (move) isGameEvent()               {}
func (possibleMoves) isGameEvent()      {}
func (showQuestion) isGameEvent()       {}
func (playerAnswerResult) isGameEvent() {}
func (gameEnd) isGameEvent()            {}

type gameStart struct{}

type playerLeft struct {
	Player PlayerID
}

// playerTurn is emitted at the start of
// a player
type playerTurn struct {
	Player PlayerID
}

// diceThrow represents the result obtained
// when throwing a dice
type diceThrow struct {
	Face uint8
}

func newDiceThrow() diceThrow {
	const maxFaceNumber = 3
	return diceThrow{uint8(rand.Int31n(maxFaceNumber) + 1)}
}

// move is emitted when a player choose to move the
// pawn
type move struct {
	Tile int
}

// possibleMoves is emitted after a diceThrow
type possibleMoves struct {
	Tiles []int // the tile indices where the clurrent player may move
}

// showQuestion is emitted when a player
// should answer a question
type showQuestion struct {
	Question  string
	Categorie categorie
}

// playerAnswerResult indicates
// if the player has answered correctly to the
// current question
type playerAnswerResult struct {
	Player  int
	Success bool
}

// gameEnd is emitted when at least one player has won
type gameEnd struct {
	Winners []string
}

// clientEvent is send by a client to the game server
type clientEvent interface {
	isClientEvent()
}

func (move) isClientEvent()   {}
func (answer) isClientEvent() {}
func (Ping) isClientEvent()   {}

// the proposition of a client to a question
type answer struct {
	Content string
}

type Ping string // DEBUG

type ClientEvent struct {
	Event  clientEvent
	Player PlayerID
}

func (ev ClientEvent) MarshalJSON() ([]byte, error) {
	tmp := struct {
		Event clientEventWrapper
	}{
		Event: clientEventWrapper{ev.Event},
	}
	return json.Marshal(tmp)
}

func (ev *ClientEvent) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Event clientEventWrapper
	}
	err := json.Unmarshal(data, &tmp)
	ev.Event = tmp.Event.Data
	return err
}
