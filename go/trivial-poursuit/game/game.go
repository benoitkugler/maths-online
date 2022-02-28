// Package trivialpoursuit implements a backend for
// a multi player trivial poursuit game, where questions
// are (short) maths questions.
package game

import (
	"fmt"
	"log"
	"sort"
	"time"
)

// PlayerID identifies a player in the game
type PlayerID = int

var QuestionDurationLimit = time.Minute

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
}

// NewGame returns an empty game, waiting for players to be
// added.
func NewGame() *Game {
	return &Game{
		GameState: GameState{
			Successes: make(map[int]*success),
			Player:    -1,
		},
		currentAnswers: make(map[int]playerAnswerResult),
	}
}

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
func (g *Game) startTurn() events {
	g.Player = g.nextPlayer()
	g.dice = newDiceThrow()
	choices := Board.choices(g.PawnTile, int(g.dice.Face)).list()
	return events{
		playerTurn{g.Player},
		g.dice,
		possibleMoves{choices},
	}
}

// StartGame actually launch the game with the players
// registred so far, which must not be empty.
func (g *Game) StartGame() GameEvents {
	evs := g.startTurn()
	return GameEvents{
		Events: append(events{gameStart{}}, evs...),
		State:  g.GameState,
	}
}

// HandleClientEvent handles the given `event`, or returns
// an error if the `event` is not valid with respect to the current
// state (enforcing rules).
func (g *Game) HandleClientEvent(event ClientEvent) (GameEvents, error) {
	switch eventData := event.Event.(type) {
	case move:
		evs, err := g.handleMove(eventData, event.Player)
		if err != nil {
			return GameEvents{}, err
		}
		return GameEvents{Events: evs, State: g.GameState}, nil
	case answer:
		evs := g.handleAnswer(eventData, event.Player)
		return GameEvents{Events: evs, State: g.GameState}, nil
	case ping:
		// safely ignore the event
		log.Printf("Client event; ping from player %d: %s", event.Player, eventData)
		return GameEvents{}, nil
	}
	return GameEvents{}, fmt.Errorf("invalid client event %T", event.Event)
}

func (g *Game) handleMove(m move, player PlayerID) (events, error) {
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
	question := g.emitQuestion()
	return events{
		m, // now valid
		question,
	}, nil
}

func (g *Game) handleAnswer(a answer, player PlayerID) events {
	isValid := g.isAnswerValid(a)
	g.Successes[player][g.question.Categorie] = isValid
	g.currentAnswers[player] = playerAnswerResult{
		Player:  player,
		Success: isValid,
	}

	return g.concludeTurn(false) // wait for other players if needed
}

// emitQuestion generate a question with the right categorie,
// and reset the current answers
func (gs *Game) emitQuestion() showQuestion {
	cat := categories[gs.PawnTile]
	question := showQuestion{
		Question:  fmt.Sprintf("Quelle est la catégorie %d", cat),
		Categorie: cat,
	}
	gs.question = question

	gs.QuestionTimeout = time.NewTimer(QuestionDurationLimit)

	return question
}

func (gs *Game) QuestionTimeoutAction() GameEvents {
	evs := gs.concludeTurn(true)
	return GameEvents{
		Events: evs,
		State:  gs.GameState,
	}
}

// isAnswerValid validdate `a` against the current question
func (gs *Game) isAnswerValid(a answer) bool {
	// TODO: à implémenter
	return a.Content == fmt.Sprintf("%d", gs.question.Categorie)
}

func (gs *Game) concludeTurn(force bool) events {
	evs := gs.endQuestion(force)
	if len(evs) == 0 { // nothing has changed
		return nil
	}

	// check for winners
	winners := gs.winners()
	if len(winners) != 0 { // end the game
		evs = append(evs, gameEnd{gs.playerNames(winners)})
	} else { // start a new turn
		evs = append(evs, gs.startTurn()...)
	}

	return evs
}

// endQuestion close the current question
// if `force` is false, it only does so if every player have answered
func (gs *Game) endQuestion(force bool) events {
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
	var out events
	for _, re := range gs.currentAnswers {
		out = append(out, re)
	}

	// cleanup
	if gs.QuestionTimeout != nil {
		gs.QuestionTimeout.Stop()
		gs.QuestionTimeout = nil
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

func (gs *Game) playerNames(players []int) []string {
	out := make([]string, len(players))
	for i, id := range players {
		out[i] = fmt.Sprintf("Joueur %d", id+1)
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
