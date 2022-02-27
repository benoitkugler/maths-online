package trivialpoursuit

import (
	"encoding/json"
	"math/rand"
)

// interaction with the client

//go:generate ../../../structgen/structgen -source=models.go -mode=dart:../../eleve/lib/trivialpoursuit/events.gen.dart -mode=itfs-json:gen.go

// GameState represents the state of the game at one point in time
type GameState struct {
	Successes []success // per-player advance
	PawnTile  int       // position of the pawn
	Player    int       // the player currently playing (choosing where to move)

	// not exported to the client

	question showQuestion // the question to answer, or empty
}

// EventRange describes a list of events yielding
// a game state.
// Clients should animate the returned events, then
// update the state, then process the next range, etc..
type EventRange struct {
	Events events
	State  GameState
	Start  int
}

type events []event

func (evs events) MarshalJSON() ([]byte, error) {
	tmp := make([]eventWrapper, len(evs))
	for i, v := range evs {
		tmp[i] = eventWrapper{data: v}
	}
	return json.Marshal(tmp)
}

func (evs *events) UnmarshalJSON(data []byte) error {
	var tmp []eventWrapper
	err := json.Unmarshal(data, &tmp)
	*evs = make(events, len(tmp))
	for i, v := range tmp {
		(*evs)[i] = v.data
	}
	return err
}

// event is an action advancing the game
// or requiring to update the UI
type event interface {
	apply(*GameState)
}

// playerTurn is emitted at the start of
// a player
type playerTurn struct {
	Player int
}

// also remove the current question
func (p playerTurn) apply(state *GameState) {
	state.Player = p.Player
	state.question = showQuestion{}
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

// apply overwrite the current dice result
func (dice diceThrow) apply(state *GameState) {}

// possibleMoves is emitted after a diceThrow
type possibleMoves struct {
	Tiles []int // the tile indices where the clurrent player may move
}

func (possibleMoves) apply(state *GameState) {}

// move is emitted when a player choose to move the
// pawn
type move struct {
	Tile int
}

// also clear the current dice
func (m move) apply(state *GameState) {
	state.PawnTile = int(m.Tile)
}

// showQuestion is emitted when a player
// should answer a question
type showQuestion struct {
	Question  string
	Categorie uint8
}

func (qu showQuestion) apply(state *GameState) {
	state.question = qu
}

// playerAnswerResult indicates
// if the player has answered correctly to the
// current question
type playerAnswerResult struct {
	Player  int
	Success bool
}

func (pa playerAnswerResult) apply(state *GameState) {
	// wrong answers remove the potential previous success
	state.Successes[pa.Player][state.question.Categorie] = pa.Success
}
