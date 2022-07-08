// Package trivialpoursuit implements a backend for
// a multi player trivial poursuit game, where questions
// are (short) maths questions.
package game

import (
	"fmt"
	"sort"
	"time"

	"github.com/benoitkugler/maths-online/maths/questions"
	"github.com/benoitkugler/maths-online/prof/editor"
	"github.com/benoitkugler/maths-online/utils"
)

const defautQuestionTimeout = time.Minute / 10

// NbCategories is the number of categories of question
const NbCategories = int(nbCategories)

// PlayerSerial identifies a player in the game
type PlayerSerial = int

type WeigthedQuestions struct {
	Questions []editor.Question
	Weights   []float64 // same length as `Questions`
}

func (wq WeigthedQuestions) sample() editor.Question {
	index := utils.SampleIndex(wq.Weights)
	return wq.Questions[index]
}

type QuestionPool [NbCategories]WeigthedQuestions

type currentQuestion struct {
	Question  questions.QuestionInstance // the instantiated version
	categorie categorie                  // the origin
	ID        int64                      // the origin
}

// MaybeUpdate is an optional update
type MaybeUpdate = *StateUpdate

type Game struct {
	// GameState is the exposed game state, shared by clients
	GameState

	// QuestionPool is the list of the question
	// being asked, for each category
	QuestionPool QuestionPool

	// QuestionTimeout is started when emitting a new question
	// and cleared early if all players have answered
	// If reached, it should trigger QuestionTimeoutAction.
	QuestionTimeout *time.Timer // may be nil

	// refreshed for each question
	currentAnswers map[PlayerSerial]bool
	question       currentQuestion // the question to answer, or empty

	// refreshed for each new turn
	currentWantNextTurn map[PlayerSerial]bool

	dice DiceThrow // last dice thrown

	QuestionDurationLimit time.Duration
	ShowDecrassage        bool

	hasStarted bool
}

// NewGame returns an empty game, using the given `questions`, waiting for players to be
// added.
// `questionTimeout` is an optionnal parameter which default to one minute
func NewGame(questionTimeout time.Duration, showDecrassage bool, questions QuestionPool) Game {
	if questionTimeout == 0 {
		questionTimeout = defautQuestionTimeout
	}

	timer := time.NewTimer(time.Second)
	timer.Stop()
	return Game{
		GameState: GameState{
			Players: make(map[int]*PlayerStatus),
			Player:  -1,
		},
		currentAnswers:        make(map[int]bool),
		currentWantNextTurn:   make(map[int]bool),
		QuestionTimeout:       timer,
		QuestionDurationLimit: questionTimeout,
		QuestionPool:          questions,
		ShowDecrassage:        showDecrassage,
	}
}

// NumberPlayers return the number of players actually in the game.
// If `onlyActive` is true, inactive players are not considered.
func (gs GameState) NumberPlayers(onlyActive bool) int {
	if onlyActive {
		var nb int
		for _, pl := range gs.Players {
			if pl.IsInactive {
				continue
			}
			nb++
		}
		return nb
	}
	return len(gs.Players)
}

func (g *Game) HasStarted() bool { return g.hasStarted }

// IsPlaying returns true if the game has started and is not finished yet.
func (g *Game) IsPlaying() bool {
	return g.HasStarted() && g.Player != -1 &&
		(len(g.winners()) == 0 ||
			!g.arePlayersReadyForNextTurn())
}

// return true when the game is at a question or question result panel
func (g *Game) inQuestion() bool { return g.question.ID != 0 }

// AddPlayer add a player to the game and returns
// its id.
// If name is empty, a name is generated.
func (g *Game) AddPlayer(name string) (PlayerSerial, LobbyUpdate) {
	max := -1
	for id := range g.Players {
		if id > max {
			max = id
		}
	}
	playerSerial := max + 1

	if name == "" {
		name = generatePlayerName(playerSerial)
	}

	g.Players[playerSerial] = &PlayerStatus{
		Name: name,
	}

	return playerSerial, LobbyUpdate{
		Player:     playerSerial,
		Names:      g.playerNames(),
		PlayerName: name,
		IsJoining:  true,
	}
}

// ReconnectPlayer marks `player` has active again, updating its name.
func (g *Game) ReconnectPlayer(player PlayerSerial, pseudo string) PlayerReconnected {
	g.Players[player].Name = pseudo
	g.Players[player].IsInactive = false
	return PlayerReconnected{
		PlayerID:   player,
		PlayerName: g.playerName(player),
	}
}

// RemovePlayer remove `player` from the game.
// The exact behavior depends on whether or not the game has started.
// If so, the player is only put in inactive mode.
// It not, the player is entirely removed.
// If the player was currently throwing or choosing the tile,
// the turn is reset to the next player.
// If a question is being answered and the `player` was the
// last answering, the question is concluded.
func (g *Game) RemovePlayer(player PlayerSerial) StateUpdate {
	playerName := g.playerName(player)

	if g.HasStarted() {
		g.GameState.Players[player].IsInactive = true
	} else {
		delete(g.GameState.Players, player)
	}

	out := Events{LobbyUpdate{
		Player:     player,
		Names:      g.playerNames(),
		PlayerName: playerName,
		IsJoining:  false,
	}}

	isInQuestion := g.inQuestion()

	if !isInQuestion && g.hasStarted && g.GameState.Player == player && g.NumberPlayers(true) > 0 {
		g.GameState.Player = g.nextPlayer()
		resetTurn := PlayerTurn{
			Player:     g.GameState.Player,
			PlayerName: g.playerName(g.GameState.Player),
		}
		out = append(out, resetTurn)
	}

	if isInQuestion {
		if endQuestion := g.concludeQuestion(false); endQuestion != nil {
			return StateUpdate{
				State:  g.GameState,
				Events: append(out, endQuestion.Events...),
			}
		}

		if endTurn := g.tryEndTurn(); endTurn != nil {
			return StateUpdate{
				State:  g.GameState,
				Events: append(out, endTurn.Events...),
			}
		}
	}

	return StateUpdate{State: g.GameState, Events: out}
}

// panic if no active players are present
func (g *Game) nextPlayer() PlayerSerial {
	var sortedIds []int
	for player := range g.Players {
		if g.Players[player].IsInactive { // ignore inactive players
			continue
		}
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

// startTurn starts a new turn, updating the state
func (g *Game) startTurn() StateUpdate {
	for k := range g.currentWantNextTurn { // reset the "ready for next turn" map
		delete(g.currentWantNextTurn, k)
	}

	g.Player = g.nextPlayer()
	return StateUpdate{
		Events: Events{PlayerTurn{g.playerName(g.Player), g.Player}},
		State:  g.GameState,
	}
}

// handleDiceClicked launches the dice, and compute the possible moves
// returns an error if the player is not allowed to click
func (g *Game) handleDiceClicked(player PlayerSerial) (StateUpdate, error) {
	// check if the player is allowed to throw the dice
	if g.Player != player {
		return StateUpdate{}, fmt.Errorf("player %d is not allowed to throw the dice during turn of player %d", player, g.Player)
	}

	g.dice = newDiceThrow()
	choices := Board.choices(g.PawnTile, int(g.dice.Face)).list()
	return StateUpdate{
		Events: Events{
			g.dice,
			PossibleMoves{PlayerName: g.playerName(g.Player), Player: g.Player, Tiles: choices},
		},
		State: g.GameState,
	}, nil
}

// StartGame actually launch the game with the players
// registred so far, which must not be empty.
// It also starts the first turn.
func (g *Game) StartGame() StateUpdate {
	g.hasStarted = true
	evs := g.startTurn()
	evs.Events = append(Events{GameStart{}}, evs.Events...)
	return evs
}

// HandleClientEvent handles the given `event`, or returns
// an error if the `event` is not valid with respect to the current
// state (enforcing rules).
// Caller should check and ignore empty return values, which mean
// nothing should happen.
func (g *Game) HandleClientEvent(event ClientEvent) (updates MaybeUpdate, isGameOver bool, err error) {
	switch eventData := event.Event.(type) {
	case DiceClicked:
		update, err := g.handleDiceClicked(event.Player)
		return &update, false, err
	case ClientMove:
		evs, err := g.handleMove(eventData, event.Player)
		if err != nil {
			return nil, false, err
		}
		return &StateUpdate{Events: evs, State: g.GameState}, false, nil
	case Answer:
		updates = g.handleAnswer(eventData, event.Player)
		return updates, false, nil
	case WantNextTurn:
		updates = g.handleWantNextTurn(eventData, event.Player)
		isGameOver = !g.IsPlaying()
		return updates, isGameOver, nil
	case Ping:
		// safely ignore the event
		return nil, false, nil
	}
	return nil, false, fmt.Errorf("invalid client event %T", event.Event)
}

func (g *Game) handleMove(m ClientMove, player PlayerSerial) (Events, error) {
	// check if the player is allowed to move
	if g.Player != player {
		return nil, fmt.Errorf("player %d is not allowed to move during turn of player %d", player, g.Player)
	}
	// check if the tile is actually reachable
	choices := Board.choices(g.PawnTile, int(g.dice.Face))
	if _, has := choices[m.Tile]; !has {
		return nil, fmt.Errorf("pawn is not allowed to move to %d", m.Tile)
	}

	g.PawnTile = m.Tile
	g.dice = DiceThrow{}
	question := g.EmitQuestion()
	return Events{
		Move{
			Tile: m.Tile, // now valid
			Path: choices[m.Tile],
		},
		question,
	}, nil
}

func (g *Game) handleAnswer(a Answer, player PlayerSerial) MaybeUpdate {
	isValid := g.isAnswerValid(a)

	// we defer the state update to the end of the question
	g.currentAnswers[player] = isValid

	return g.concludeQuestion(false) // wait for other players if needed
}

// EmitQuestion generate a question with the right categorie
func (gs *Game) EmitQuestion() ShowQuestion {
	// select the category
	cat := categories[gs.PawnTile]
	// select the question among the pool...
	question := gs.QuestionPool[cat].sample()
	// ...and instantiate it
	instance := question.Page.Instantiate()

	gs.question = currentQuestion{
		categorie: cat,
		ID:        question.Id,
		Question:  instance,
	}

	out := ShowQuestion{
		TimeoutSeconds: int(gs.QuestionDurationLimit.Seconds()),
		ID:             question.Id,
		Categorie:      cat,
		Question:       instance.ToClient(),
	}

	gs.QuestionTimeout.Reset(gs.QuestionDurationLimit)

	return out
}

// QuestionTimeoutAction closes the current question session,
// but doest not start a new turn.
func (gs *Game) QuestionTimeoutAction() MaybeUpdate {
	return gs.concludeQuestion(true)
}

func (gs *Game) concludeQuestion(force bool) MaybeUpdate {
	evs := gs.endQuestion(force)
	if len(evs) == 0 { // nothing has changed
		return nil
	}
	return &StateUpdate{Events: evs, State: gs.GameState}
}

// isAnswerValid validate `a` against the current question
func (gs *Game) isAnswerValid(a Answer) bool {
	result := gs.question.Question.EvaluateAnswer(a.Answer).IsCorrect()
	return result
}

func (gs *Game) arePlayersReadyForNextTurn() bool {
	for player, pl := range gs.Players {
		if pl.IsInactive {
			continue
		}
		if ok := gs.currentWantNextTurn[player]; !ok {
			return false
		}
	}
	return true
}

// if all the players are ready, go to the next turn (or end the game if needed)
// otherwise, it is a no-op
func (gs *Game) tryEndTurn() (updates MaybeUpdate) {
	if !gs.arePlayersReadyForNextTurn() { // do nothing
		return nil
	}

	// reset the question
	gs.question = currentQuestion{}

	// check for winners
	winners := gs.winners()
	isGameOver := len(winners) != 0
	if isGameOver { // end the game
		updates = &StateUpdate{
			Events: Events{GameEnd{
				Winners:               winners,
				WinnerNames:           gs.idToNames(winners),
				QuestionDecrassageIds: gs.decrassage(),
			}},
			State: gs.GameState,
		}
	} else { // start a new turn
		v := gs.startTurn()
		updates = &v
	}

	return updates
}

func (gs *Game) handleWantNextTurn(event WantNextTurn, player PlayerSerial) (updates MaybeUpdate) {
	gs.currentWantNextTurn[player] = true

	pReview := &gs.Players[player].Review
	if event.MarkQuestion {
		pReview.MarkedQuestions = append(pReview.MarkedQuestions, gs.question.ID)
	}

	return gs.tryEndTurn()
}

// endQuestion close the current question
// if `force` is false, it only does so if every player have answered
func (gs *Game) endQuestion(force bool) Events {
	hasAllAnswered := true
	for player, pl := range gs.Players {
		if _, has := gs.currentAnswers[player]; !has && !pl.IsInactive {
			hasAllAnswered = false
			break
		}
	}
	if !hasAllAnswered && !force { // abort closing
		return nil
	}

	out := PlayerAnswerResults{
		Categorie: gs.question.categorie,
		Results:   make(map[int]playerAnswerResult),
	}

	// return the answers event, defaulting to
	// false for no answer
	for player, state := range gs.Players {
		// we still mark invalid answsers for inactive player,
		// to avoid cheating by leaving before right before the question

		isValid, _ := gs.currentAnswers[player]
		// update the success
		state.Success[gs.question.categorie] = isValid // false if not answered
		state.Review.QuestionHistory = append(state.Review.QuestionHistory, QR{
			IdQuestion: gs.question.ID,
			Success:    isValid,
		})
		askForMark := !isValid && len(state.Review.MarkedQuestions) < 3

		out.Results[player] = playerAnswerResult{Success: isValid, AskForMask: askForMark}
	}

	// cleanup
	stopped := gs.QuestionTimeout.Stop()
	if !stopped && !force {
		<-gs.QuestionTimeout.C // drain the channel
	}

	// question is used in wantNextTurn
	for k := range gs.currentAnswers {
		delete(gs.currentAnswers, k)
	}

	return Events{out}
}

// winners returns the players who win, or an empty slice
// use it to check if the game is over
func (gs *Game) winners() (out []int) {
	for player, state := range gs.Players {
		if state.Success.isDone() {
			out = append(out, player)
		}
	}
	sort.Ints(out)
	return out
}

// returns nil if !showDecrassage
func (gs *Game) decrassage() (ids map[int][]int64) {
	if !gs.ShowDecrassage {
		return nil
	}

	const nbMax = 3
	ids = make(map[int][]int64)
	for player, state := range gs.Players {
		questions := state.Review.MarkedQuestions
		// add from the failed questions
		for _, question := range state.Review.QuestionHistory {
			if len(questions) >= nbMax {
				break
			}
			if !question.Success {
				questions = append(questions, question.IdQuestion)
			}
		}
		ids[player] = questions
	}
	return ids
}

func generatePlayerName(player PlayerSerial) string {
	return fmt.Sprintf("Joueur %d", player+1)
}

func (g *Game) playerName(player PlayerSerial) string {
	return g.Players[player].Name
}

func (g *Game) playerNames() map[PlayerSerial]string {
	out := make(map[int]string, len(g.Players))
	for player := range g.Players {
		out[player] = g.playerName(player)
	}
	return out
}

func (gs *Game) idToNames(players []PlayerSerial) []string {
	out := make([]string, len(players))
	for i, id := range players {
		out[i] = gs.playerName(id)
	}
	return out
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
