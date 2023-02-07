package trivial

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	"github.com/benoitkugler/maths-online/server/src/sql/editor"
)

// interaction with the client

// NbCategories is the number of categories of question
const NbCategories = int(nbCategories)

// serial identifies a player in the game
type serial = PlayerID

// GameState represents an on-going game.
type GameState struct {
	Players    map[serial]PlayerStatus // per-player advance
	PawnTile   int                     // position of the pawn
	PlayerTurn serial                  // the player currently playing (choosing where to move)
}

type QR struct {
	IdQuestion editor.IdQuestion
	Success    bool
}

// QuestionReview stores the results of one player
// against the questions asked during the game
type QuestionReview struct {
	QuestionHistory []QR
	// Ids of question the player wants to mark for further work
	MarkedQuestions []editor.IdQuestion
}

// Success are the categories completed by a player
type Success [NbCategories]bool

func (sc Success) isDone() bool {
	for _, b := range sc {
		if !b {
			return false
		}
	}
	return true
}

// PlayerStatus exposes the information about one player
type PlayerStatus struct {
	Name    string
	Review  QuestionReview
	Success Success
	// Has the player disconnect ?
	IsInactive bool
}

// StateUpdate describes a list of events yielding
// a new game state.
// Clients should animate the events and update the state.
type StateUpdate struct {
	Events Events
	State  GameState
}

func (su StateUpdate) String() string {
	var events []string
	for _, ev := range su.Events {
		events = append(events, fmt.Sprintf("\t%T: %v", ev, ev))
	}
	return fmt.Sprintf("Events:\n %s\n--> New state: %v\n", strings.Join(events, "\n"), su.State)
}

type Events []ServerEvent

// ServerEvent is an action (created by the server) advancing the game
// or requiring to update the UI
type ServerEvent interface {
	isServerEvent()
}

func (PlayerJoin) isServerEvent()                   {}
func (PlayerReconnected) isServerEvent()            {}
func (LobbyUpdate) isServerEvent()                  {}
func (GameStart) isServerEvent()                    {}
func (PlayerLeft) isServerEvent()                   {}
func (PlayerTurn) isServerEvent()                   {}
func (DiceThrow) isServerEvent()                    {}
func (Move) isServerEvent()                         {}
func (PossibleMoves) isServerEvent()                {}
func (ShowQuestion) isServerEvent()                 {}
func (PlayerAnswerResults) isServerEvent()          {}
func (PlayersStillInQuestionResult) isServerEvent() {}
func (GameEnd) isServerEvent()                      {}
func (GameTerminated) isServerEvent()               {}

// PlayerJoin is only emitted to the actual player
// who join the game
type PlayerJoin struct {
	Player serial
}

type PlayerReconnected struct {
	ID     serial
	Pseudo string
}

type LobbyUpdate struct {
	PlayerPseudos map[serial]string // the new players in the lobby
	Pseudo        string
	ID            serial // the player who joined or left
	IsJoining     bool   // false for leaving
}

type GameStart struct{}

type PlayerLeft struct {
	Player serial
}

// PlayerTurn is emitted at the start of
// a player
type PlayerTurn struct {
	PlayerName string
	Player     serial
}

// DiceThrow represents the result obtained
// when throwing a dice
type DiceThrow struct {
	Face uint8
}

func newDiceThrow() DiceThrow {
	const maxFaceNumber = 3
	return DiceThrow{uint8(rand.Int31n(maxFaceNumber) + 1)}
}

// Move is emitted when a player choose to Move the
// pawn
type Move struct {
	// the tiles to go through to animate the move
	// (only valid when send by the server)
	Path []int
	Tile int
}

// PossibleMoves is emitted after a diceThrow
type PossibleMoves struct {
	PlayerName string
	Tiles      []int  // the tile indices where the current player may move
	Player     serial // the player allowed to play
}

// ShowQuestion is emitted when a player
// should answer a question
type ShowQuestion struct {
	TimeoutSeconds int
	Categorie      Categorie
	ID             editor.IdQuestion // to facilitate the tracking of the question results
	Question       client.Question   // the actual question
}

// PlayerAnswerResults indicates
// if the players have answered correctly to the
// current question
type PlayerAnswerResults struct {
	Categorie Categorie
	Results   map[serial]playerAnswerResult
}

type playerAnswerResult struct {
	Success    bool
	AskForMask bool // true if Success is false and if the player has not already marked 3
}

// PlayersStillInQuestionResult is emitted when some players are
// still in the [pQuestionResult] step
type PlayersStillInQuestionResult struct {
	Players     []serial
	PlayerNames []string
}

// GameEnd is emitted when at least one player has won
type GameEnd struct {
	QuestionDecrassageIds map[serial][]editor.IdQuestion // player->questions
	Winners               []serial
	WinnerNames           []string
}

// GameTerminated is emitted when the game
// is manually terminated by the teacher
type GameTerminated struct{}

// ClientEventITF is the common interface for
// events send by a student client to the game server
type ClientEventITF interface {
	isClientEvent()
}

func (ClientMove) isClientEvent()   {}
func (Answer) isClientEvent()       {}
func (DiceClicked) isClientEvent()  {}
func (WantNextTurn) isClientEvent() {}
func (Ping) isClientEvent()         {}

type ClientMove Move

// the proposition of a client to a question
type Answer struct {
	Answer client.QuestionAnswersIn
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
