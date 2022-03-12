package game

import (
	"encoding/json"
	"math/rand"
)

// interaction with the client

//go:generate ../../../../../structgen/structgen -source=models.go -mode=dart:../../../../eleve/lib/trivialpoursuit/events.gen.dart -mode=itfs-json:gen.go

// GameState represents an on-going game.
type GameState struct {
	Players  map[PlayerID]*PlayerStatus // per-player advance
	PawnTile int                        // position of the pawn
	Player   int                        // the player currently playing (choosing where to move)
}

// PlayerStatus exposes the information about one player
type PlayerStatus struct {
	Name    string
	Success success
}

type StateUpdates = []StateUpdate

// StateUpdate describes a list of events yielding
// a new game state.
// Clients should animate the events and
// update the state.
type StateUpdate struct {
	Events Events
	State  GameState
}

type Events []GameEvent

func (evs Events) MarshalJSON() ([]byte, error) {
	tmp := make([]GameEventWrapper, len(evs))
	for i, v := range evs {
		tmp[i] = GameEventWrapper{Data: v}
	}
	return json.Marshal(tmp)
}

func (evs *Events) UnmarshalJSON(data []byte) error {
	var tmp []GameEventWrapper
	err := json.Unmarshal(data, &tmp)
	*evs = make(Events, len(tmp))
	for i, v := range tmp {
		(*evs)[i] = v.Data
	}
	return err
}

// GameEvent is an action (created by the server) advancing the game
// or requiring to update the UI
type GameEvent interface {
	isGameEvent()
}

func (PlayerJoin) isGameEvent()         {}
func (LobbyUpdate) isGameEvent()        {}
func (gameStart) isGameEvent()          {}
func (playerLeft) isGameEvent()         {}
func (playerTurn) isGameEvent()         {}
func (diceThrow) isGameEvent()          {}
func (move) isGameEvent()               {}
func (possibleMoves) isGameEvent()      {}
func (showQuestion) isGameEvent()       {}
func (playerAnswerResult) isGameEvent() {}
func (gameEnd) isGameEvent()            {}

// PlayerJoin is only emitted to the actual player
// who join the game
type PlayerJoin struct {
	Player PlayerID
}

type LobbyUpdate struct {
	Names      map[PlayerID]string // the new players in the lobby
	PlayerName string
	Player     PlayerID // the player who joined or left
	IsJoining  bool     // false for leaving
}

type gameStart struct{}

type playerLeft struct {
	Player PlayerID
}

// playerTurn is emitted at the start of
// a player
type playerTurn struct {
	PlayerName string
	Player     PlayerID
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
	// the tiles to go through to animate the move
	// (only valid when send by the server)
	Path []int
	Tile int
}

// possibleMoves is emitted after a diceThrow
type possibleMoves struct {
	PlayerName string
	Tiles      []int    // the tile indices where the current player may move
	Player     PlayerID // the player allowed to play
}

// showQuestion is emitted when a player
// should answer a question
type showQuestion struct {
	Question       string
	TimeoutSeconds int
	Categorie      categorie
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
	Winners     []int
	WinnerNames []string
}

// clientEventData is send by a client to the game server
type clientEventData interface {
	isClientEvent()
}

func (move) isClientEvent()        {}
func (answer) isClientEvent()      {}
func (diceClicked) isClientEvent() {}
func (Ping) isClientEvent()        {}

// the proposition of a client to a question
type answer struct {
	Content string
}

// diceClicked is emitted when the current player
// throws the dice
type diceClicked struct{}

// Ping is used to maintain the client connection
// openned
type Ping struct {
	Info string
}

type ClientEvent struct {
	Event  clientEventData
	Player PlayerID
}

func (ev ClientEvent) MarshalJSON() ([]byte, error) {
	tmp := struct {
		Event clientEventDataWrapper
	}{
		Event: clientEventDataWrapper{ev.Event},
	}
	return json.Marshal(tmp)
}

func (ev *ClientEvent) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Event clientEventDataWrapper
	}
	err := json.Unmarshal(data, &tmp)
	ev.Event = tmp.Event.Data
	return err
}
