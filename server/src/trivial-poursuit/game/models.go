package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"github.com/benoitkugler/maths-online/maths/exercice/client"
)

// interaction with the client

//go:generate ../../../../../structgen/structgen -source=models.go -mode=dart:../../../../eleve/lib/trivialpoursuit/events.gen.dart -mode=itfs-json:gen_itfs.go

// GameState represents an on-going game.
type GameState struct {
	Players  map[PlayerID]*PlayerStatus // per-player advance
	PawnTile int                        // position of the pawn
	Player   PlayerID                   // the player currently playing (choosing where to move)
}

type QR struct {
	IdQuestion int64
	Success    bool
}

// QuestionReview stores the results of one player
// against the questions asked during the game
type QuestionReview struct {
	QuestionHistory []QR
	// Ids of question the player wants to mark for further work
	MarkedQuestions []int64
}

// PlayerStatus exposes the information about one player
type PlayerStatus struct {
	Name    string
	Review  QuestionReview
	Success Success
}

// StateUpdate describes a list of events yielding
// a new game state.
// Clients should animate the events and
// update the state.
type StateUpdate struct {
	Events Events
	State  GameState
}

func (su StateUpdate) String() string {
	var events []string
	for _, ev := range su.Events {
		events = append(events, fmt.Sprintf("\t%T: %v", ev, ev))
	}
	return fmt.Sprintf("Events:\n %s\n--> New state: %v", strings.Join(events, "\n"), su.State)
}

type Events []GameEvent

// GameEvent is an action (created by the server) advancing the game
// or requiring to update the UI
type GameEvent interface {
	isGameEvent()
}

func (PlayerJoin) isGameEvent()          {}
func (LobbyUpdate) isGameEvent()         {}
func (gameStart) isGameEvent()           {}
func (playerLeft) isGameEvent()          {}
func (playerTurn) isGameEvent()          {}
func (diceThrow) isGameEvent()           {}
func (move) isGameEvent()                {}
func (PossibleMoves) isGameEvent()       {}
func (showQuestion) isGameEvent()        {}
func (playerAnswerResults) isGameEvent() {}
func (gameEnd) isGameEvent()             {}
func (GameTerminated) isGameEvent()      {}

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

// PossibleMoves is emitted after a diceThrow
type PossibleMoves struct {
	PlayerName string
	Tiles      []int    // the tile indices where the current player may move
	Player     PlayerID // the player allowed to play
}

// showQuestion is emitted when a player
// should answer a question
type showQuestion struct {
	TimeoutSeconds int
	Categorie      categorie
	ID             int64           // to facilitate the tracking of the question results
	Question       client.Question `dart-extern:"../exercices/types.gen.dart"` // the actual question
}

// playerAnswerResults indicates
// if the players have answered correctly to the
// current question
type playerAnswerResults struct {
	Categorie categorie
	Results   map[PlayerID]playerAnswerResult
}

type playerAnswerResult struct {
	Success    bool
	AskForMask bool // true if Success is false and if the player has not already marked 3
}

// gameEnd is emitted when at least one player has won
type gameEnd struct {
	Winners     []int
	WinnerNames []string
}

// GameTerminated is emitted when the game
// is manually terminated by the teacher
type GameTerminated struct{}

// clientEventData is send by a client to the game server
type clientEventData interface {
	isClientEvent()
}

func (ClientMove) isClientEvent()   {}
func (Answer) isClientEvent()       {}
func (DiceClicked) isClientEvent()  {}
func (WantNextTurn) isClientEvent() {}
func (Ping) isClientEvent()         {}

type ClientMove move

// the proposition of a client to a question
type Answer struct {
	Answer client.QuestionAnswersIn `dart-extern:"../exercices/types.gen.dart"`
}

// DiceClicked is emitted when the current player
// throws the dice
type DiceClicked struct{}

// WantNextTurn is emitted when a player is done
// looking at the question answer panel
type WantNextTurn struct {
	// MarkQuestion is true if the player wants to
	// keep the question for following trainings
	MarkQuestion bool
}

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
