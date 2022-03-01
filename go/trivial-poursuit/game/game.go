// Package trivialpoursuit implements a backend for
// a multi player trivial poursuit game, where questions
// are (short) maths questions.
package game

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"
)

// PlayerID identifies a player in the game
type PlayerID = int

var DebugLogger = log.New(os.Stdout, "game-debug:", log.LstdFlags)

type Game struct {
	// GameState is the exposed game state, shared by clients
	GameState

	// QuestionTimeout is started when emitting a new question
	// and cleared early if all players have answered
	// If reached, it should trigger QuestionTimeoutAction.
	QuestionTimeout *time.Timer // may be nil

	// refreshed for each question
	currentAnswers map[PlayerID]playerAnswerResult
	question       showQuestion // the question to answer, or empty

	dice diceThrow // last dice thrown

	questionDurationLimit time.Duration
}

// NewGame returns an empty game, waiting for players to be
// added.
// `questionTimeout` is an optionnal parameter which default to one minute
func NewGame(questionTimeout time.Duration) *Game {
	if questionTimeout == 0 {
		questionTimeout = time.Minute
	}

	timer := time.NewTimer(time.Second)
	timer.Stop()
	return &Game{
		GameState: GameState{
			Successes: make(map[int]*success),
			Player:    -1,
		},
		currentAnswers:        make(map[int]playerAnswerResult),
		QuestionTimeout:       timer,
		questionDurationLimit: questionTimeout,
	}
}

// NumberPlayers return the number of players actually in the game.
func (gs GameState) NumberPlayers() int { return len(gs.Successes) }

// AddPlayer add a player to the game and returns
// its id.
func (g *Game) AddPlayer() PlayerID {
	max := -1
	for id := range g.Successes {
		if id > max {
			max = id
		}
	}
	playerID := max + 1
	g.Successes[playerID] = &success{}
	return playerID
}

// RemovePlayer remove `player` from the game.
func (g *Game) RemovePlayer(player PlayerID) {
	delete(g.Successes, player)
}

// panic if no players are present
func (g *Game) nextPlayer() PlayerID {
	var sortedIds []int
	for player := range g.Successes {
		sortedIds = append(sortedIds, player)
	}
	sort.Ints(sortedIds)

	for _, player := range sortedIds {
		if player > g.Player {
			return player
		}
	}
	return sortedIds[0]
}

// convenient method to start a new turn,
// launch the dice, and compute the possible moves
func (g *Game) startTurn() GameEvents {
	g.Player = g.nextPlayer()
	g.dice = newDiceThrow()
	choices := Board.choices(g.PawnTile, int(g.dice.Face)).list()
	return GameEvents{
		Events: Events{
			playerTurn{g.Player, g.playerName(g.Player)},
			g.dice,
			possibleMoves{CurrentPlayer: g.Player, Tiles: choices},
		},
		State: g.GameState,
	}
}

// StartGame actually launch the game with the players
// registred so far, which must not be empty.
func (g *Game) StartGame() GameEvents {
	evs := g.startTurn()
	evs.Events = append(Events{gameStart{}}, evs.Events...)
	return evs
}

// HandleClientEvent handles the given `event`, or returns
// an error if the `event` is not valid with respect to the current
// state (enforcing rules).
// Caller should check and ignore empty return values.
func (g *Game) HandleClientEvent(event ClientEvent) (EventList, error) {
	switch eventData := event.Event.(type) {
	case move:
		evs, err := g.handleMove(eventData, event.Player)
		if err != nil {
			return nil, err
		}
		return EventList{{Events: evs, State: g.GameState}}, nil
	case answer:
		evs := g.handleAnswer(eventData, event.Player)
		return evs, nil
	case Ping:
		// safely ignore the event
		DebugLogger.Printf("PING event (from player %d): %s", event.Player, eventData.Info)
		return nil, nil
	}
	return nil, fmt.Errorf("invalid client event %T", event.Event)
}

func (g *Game) handleMove(m move, player PlayerID) (Events, error) {
	// check if the player is allowed to move
	if g.Player != player {
		return nil, fmt.Errorf("player %d is not allowed to move during turn of player %d", player, g.Player)
	}
	// check if the tile is actually reachable
	choices := Board.choices(g.PawnTile, int(g.dice.Face))
	if !choices[m.Tile] {
		return nil, fmt.Errorf("pawn is not allowed to move to %d", m.Tile)
	}

	g.PawnTile = m.Tile
	g.dice = diceThrow{}
	question := g.EmitQuestion()
	return Events{
		m, // now valid
		question,
	}, nil
}

func (g *Game) handleAnswer(a answer, player PlayerID) EventList {
	isValid := g.isAnswerValid(a)
	g.Successes[player][g.question.Categorie] = isValid
	g.currentAnswers[player] = playerAnswerResult{
		Player:  player,
		Success: isValid,
	}

	return g.concludeTurn(false) // wait for other players if needed
}

// EmitQuestion generate a question with the right categorie,
// and reset the current answers
func (gs *Game) EmitQuestion() showQuestion {
	cat := categories[gs.PawnTile]
	question := showQuestion{
		Question:  fmt.Sprintf("Quelle est la catégorie %d", cat),
		Categorie: cat,
	}
	gs.question = question

	gs.QuestionTimeout.Reset(gs.questionDurationLimit)

	return question
}

func (gs *Game) QuestionTimeoutAction() EventList {
	return gs.concludeTurn(true)
}

// isAnswerValid validdate `a` against the current question
func (gs *Game) isAnswerValid(a answer) bool {
	// TODO: à implémenter
	return a.Content == fmt.Sprintf("%d", gs.question.Categorie)
}

func (gs *Game) concludeTurn(force bool) EventList {
	evs := gs.endQuestion(force)
	if len(evs) == 0 { // nothing has changed
		return nil
	}

	out := EventList{{Events: evs, State: gs.GameState}}

	// check for winners
	winners := gs.winners()
	if len(winners) != 0 { // end the game
		out = append(out, GameEvents{
			Events: Events{gameEnd{gs.playerNames(winners)}},
			State:  gs.GameState,
		})
	} else { // start a new turn
		out = append(out, gs.startTurn())
	}

	return out
}

// endQuestion close the current question
// if `force` is false, it only does so if every player have answered
func (gs *Game) endQuestion(force bool) Events {
	hasAllAnswered := true
	for player := range gs.Successes {
		if _, has := gs.currentAnswers[player]; !has {
			hasAllAnswered = false
			break
		}
	}
	if !hasAllAnswered && !force { // abort closing
		return nil
	}

	// return the answers event
	var out Events
	for _, re := range gs.currentAnswers {
		out = append(out, re)
	}

	// cleanup
	stopped := gs.QuestionTimeout.Stop()
	if !stopped && !force {
		<-gs.QuestionTimeout.C // drain the channel
	}

	gs.question = showQuestion{}
	for k := range gs.currentAnswers {
		delete(gs.currentAnswers, k)
	}

	return out
}

// winners returns the players who win, or an empty slice
// use it to check if the game is over
func (gs *Game) winners() (out []int) {
	for player, success := range gs.Successes {
		if success.isDone() {
			out = append(out, player)
		}
	}
	sort.Ints(out)
	return out
}

func (gs *Game) playerName(player PlayerID) string {
	return fmt.Sprintf("Joueur %d", player+1)
}

func (gs *Game) playerNames(players []PlayerID) []string {
	out := make([]string, len(players))
	for i, id := range players {
		out[i] = gs.playerName(id)
	}
	return out
}

// the categories completed by a player
type success [nbCategories]bool

func (sc success) isDone() bool {
	for _, b := range sc {
		if !b {
			return false
		}
	}
	return true
}
