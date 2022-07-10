package trivial

import (
	"fmt"
	"sort"
	"time"

	"github.com/benoitkugler/maths-online/maths/questions"
	"github.com/benoitkugler/maths-online/prof/editor"
	"github.com/benoitkugler/maths-online/utils"
)

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

// Phase identifies the current phase of the game
type Phase uint8

const (
	PWaiting  Phase = iota // not started yet
	PThrowing              // start of turn, waiting for dice throw
	PMoving                // dice was thrown, waiting for player move
	PQuestion              // question is being answered
	PResult                // players are consulting answer results
	POver                  // game has finished
)

// game is the internal state of the game
type game struct {
	phase      Phase  // phase of the game
	pawnTile   int    // position of the pawn on the board
	playerTurn serial // the player currently playing (for instance, choosing where to move)

	options Options

	// questionTimer is started when emitting a new question
	// and cleared early if all players have answered
	// If reached, it should trigger QuestionTimeoutAction.
	questionTimer *time.Timer

	// refreshed for each question
	currentAnswers map[serial]bool
	question       currentQuestion // the question to answer, or empty

	// refreshed for each new turn
	currentWantNextTurn map[serial]bool

	dice DiceThrow // last dice thrown
}

// newGame returns an empty game, using the given `options`
func newGame(options Options) game {
	timer := time.NewTimer(time.Second /* ignored */)
	timer.Stop()
	return game{
		options:             options,
		playerTurn:          "",
		currentAnswers:      make(map[serial]bool),
		currentWantNextTurn: make(map[serial]bool),
		questionTimer:       timer,
	}
}

func (g *game) hasStarted() bool { return g.phase != PWaiting }

func (r *Room) nbActivePlayers() int {
	var out int
	for _, pl := range r.players {
		if pl.conn != nil {
			out++
		}
	}
	return out
}

func (r *Room) playerPseudos() map[serial]string {
	out := make(map[serial]string, len(r.players))
	for _, player := range r.players {
		out[player.pl.ID] = player.pl.Pseudo
	}
	return out
}

func (r *Room) serialsToPseudos(players []serial) []string {
	dict := r.playerPseudos()
	out := make([]string, len(players))
	for i, se := range players {
		out[i] = dict[se]
	}
	return out
}

func (r *Room) serialToPseudo(se serial) string { return r.players[se].pl.Pseudo }

func (r *Room) tryStartGame() Events {
	// before starting, all players are active since deconnecting
	// exlude them from the lobby (see `removePlayer`)
	if len(r.players) >= r.game.options.PlayersNumber {
		ProgressLogger.Printf("Game %s : starting...", r.ID)

		// also starts the first turn.
		eventNewTurn := r.startTurn()
		return Events{GameStart{}, eventNewTurn}
	}

	return nil
}

// removePlayer remove `player` from the game.
// The exact behavior depends on whether or not the game has started.
// If so, the player is only put in inactive mode.
// It not, the player is entirely removed.
// If the player was currently throwing or choosing the tile,
// the turn is reset to the next player.
// If a question is being answered and the `player` was the
// last answering, the question is concluded.
func (r *Room) removePlayer(player Player) Events {
	playerName := player.Pseudo

	if r.game.hasStarted() {
		r.players[player.ID].conn = nil
	} else {
		delete(r.players, player.ID)
	}

	out := Events{LobbyUpdate{
		ID:            player.ID,
		Pseudo:        playerName,
		IsJoining:     false,
		PlayerPseudos: r.playerPseudos(),
	}}

	switch r.game.phase {
	case PWaiting, POver:
		// nothing more to be done
	case PThrowing, PMoving: // if it is the current player, reset the turn
		if r.game.playerTurn == player.ID && r.nbActivePlayers() > 0 {
			resetTurn := r.startTurn()
			out = append(out, resetTurn)
		}
	case PQuestion:
		endQuestion := r.tryEndQuestion(false)
		out = append(out, endQuestion...)
	case PResult:
		endTurn := r.tryEndTurn()
		out = append(out, endTurn...)
	default:
		panic("exhaustive switch")
	}

	return out
}

// isAnswerValid validate `a` against the current question
func (g *game) isAnswerValid(a Answer) bool {
	result := g.question.Question.EvaluateAnswer(a.Answer).IsCorrect()
	return result
}

// tryEndQuestion close the current question
// if `force` is false, it only does so if every player have answered
func (r *Room) tryEndQuestion(force bool) Events {
	hasAllAnswered := true
	for _, pl := range r.players {
		if _, has := r.game.currentAnswers[pl.pl.ID]; !has && pl.conn != nil {
			hasAllAnswered = false
			break
		}
	}
	if !hasAllAnswered && !force { // abort closing
		return nil
	}

	out := PlayerAnswerResults{
		Categorie: r.game.question.categorie,
		Results:   make(map[serial]playerAnswerResult),
	}

	// return the answers event, defaulting to
	// false for no answer
	for _, player := range r.players {
		// we still mark invalid answsers for inactive player,
		// to avoid cheating by leaving before right before the question

		isValid, _ := r.game.currentAnswers[player.pl.ID]
		// update the success
		player.advance.success[r.game.question.categorie] = isValid // false if not answered
		player.advance.review.QuestionHistory = append(player.advance.review.QuestionHistory, QR{
			IdQuestion: r.game.question.ID,
			Success:    isValid,
		})
		askForMark := !isValid && len(player.advance.review.MarkedQuestions) < 3

		out.Results[player.pl.ID] = playerAnswerResult{Success: isValid, AskForMask: askForMark}
	}

	// cleanup
	stopped := r.game.questionTimer.Stop()
	if !stopped && !force {
		<-r.game.questionTimer.C // drain the channel
	}

	// question is used in wantNextTurn
	for k := range r.game.currentAnswers {
		delete(r.game.currentAnswers, k)
	}

	r.game.phase = PResult

	return Events{out}
}

func (r *Room) arePlayersReadyForNextTurn() bool {
	for _, pl := range r.players {
		if pl.conn == nil {
			continue
		}
		if ok := r.game.currentWantNextTurn[pl.pl.ID]; !ok {
			return false
		}
	}
	return true
}

// if all the players are ready, go to the next turn (or end the game if needed)
// otherwise, it is a no-op
func (r *Room) tryEndTurn() Events {
	if !r.arePlayersReadyForNextTurn() { // do nothing
		return nil
	}

	// reset the question
	r.game.question = currentQuestion{}

	// check for winners
	winners := r.winners()
	isGameOver := len(winners) != 0
	if isGameOver { // end the game
		r.game.phase = POver
		return Events{
			GameEnd{
				Winners:               winners,
				WinnerNames:           r.serialsToPseudos(winners),
				QuestionDecrassageIds: r.decrassage(),
			},
		}
	}

	// else, start a new turn
	v := r.startTurn()
	return Events{v}
}

// winners returns the players who win, or an empty slice
// use it to check if the game is over
func (r *Room) winners() (out []serial) {
	for _, player := range r.players {
		if player.advance.success.isDone() {
			out = append(out, player.pl.ID)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

// returns nil if `ShowDecrassage` is false
func (r *Room) decrassage() (ids map[serial][]int64) {
	if !r.game.options.ShowDecrassage {
		return nil
	}

	const nbMax = 3
	ids = make(map[serial][]int64)
	for _, player := range r.players {
		quIds := editor.NewSetFromSlice(player.advance.review.MarkedQuestions)

		// add from the failed questions
		for _, question := range player.advance.review.QuestionHistory {
			if !question.Success {
				quIds.Add(question.IdQuestion)
			}
		}

		questions := quIds.Keys()
		if len(questions) >= nbMax {
			questions = questions[0:nbMax]
		}
		ids[player.pl.ID] = questions
	}
	return ids
}

// panic if no active players are present
func (r *Room) nextPlayer() serial {
	var sortedIds []string
	for _, player := range r.players {
		if player.conn == nil { // ignore inactive players
			continue
		}
		sortedIds = append(sortedIds, string(player.pl.ID))
	}
	sort.Strings(sortedIds)

	for _, player := range sortedIds {
		if serial(player) > r.game.playerTurn {
			return serial(player)
		}
	}
	return serial(sortedIds[0])
}

// startTurn starts a new turn, updating the state
// if will panic if no active players are present
func (r *Room) startTurn() PlayerTurn {
	for k := range r.game.currentWantNextTurn { // reset the "ready for next turn" map
		delete(r.game.currentWantNextTurn, k)
	}

	r.game.phase = PThrowing
	r.game.playerTurn = r.nextPlayer()
	return PlayerTurn{
		Player:     r.game.playerTurn,
		PlayerName: r.serialToPseudo(r.game.playerTurn),
	}
}

// handleClientEvent handles the given `event`, or returns
// an error if the `event` is not valid with respect to the current
// state (enforcing rules).
// Caller should check and ignore empty return values, which mean
// nothing should happen.
func (r *Room) handleClientEvent(event ClientEventITF, player Player) (events Events, isGameOver bool, err error) {
	switch eventData := event.(type) {
	case DiceClicked:
		events, err := r.handleDiceClicked(player.ID)
		return events, false, err
	case ClientMove:
		events, err := r.handleMove(eventData, player.ID)
		return events, false, err
	case Answer:
		events := r.handleAnswer(eventData, player.ID)
		return events, false, nil
	case WantNextTurn:
		events := r.handleWantNextTurn(eventData, player)
		return events, r.game.phase == POver, nil
	case Ping:
		// safely ignore the event
		return nil, false, nil
	default:
		return nil, false, fmt.Errorf("invalid client event %T", event)
	}
}

// handleDiceClicked launches the dice, and compute the possible moves
// returns an error if the player is not allowed to click
func (r *Room) handleDiceClicked(player serial) (Events, error) {
	g := &r.game
	// check if the player is allowed to throw the dice
	if g.playerTurn != player {
		return nil, fmt.Errorf("player %s is not allowed to throw the dice during turn of player %s", player, g.playerTurn)
	}

	g.dice = newDiceThrow()
	choices := Board.choices(g.pawnTile, int(g.dice.Face)).list()
	g.phase = PMoving
	return Events{
		g.dice,
		PossibleMoves{PlayerName: r.serialToPseudo(player), Player: player, Tiles: choices},
	}, nil
}

func (r *Room) handleMove(m ClientMove, player serial) (Events, error) {
	g := &r.game
	// check if the player is allowed to move
	if g.playerTurn != player {
		return nil, fmt.Errorf("player %s is not allowed to move during turn of player %s", player, g.playerTurn)
	}
	// check if the tile is actually reachable
	choices := Board.choices(g.pawnTile, int(g.dice.Face))
	if _, has := choices[m.Tile]; !has {
		return nil, fmt.Errorf("pawn is not allowed to move to %d", m.Tile)
	}

	g.pawnTile = m.Tile
	g.dice = DiceThrow{}
	question := g.emitQuestion()
	return Events{
		Move{
			Tile: m.Tile, // now valid
			Path: choices[m.Tile],
		},
		question,
	}, nil
}

// emitQuestion generate a question with the right categorie
func (gs *game) emitQuestion() ShowQuestion {
	gs.phase = PQuestion

	// select the category
	cat := categories[gs.pawnTile]
	// select the question among the pool...
	question := gs.options.Questions[cat].sample()
	// ...and instantiate it
	instance := question.Page.Instantiate()

	gs.question = currentQuestion{
		categorie: cat,
		ID:        question.Id,
		Question:  instance,
	}

	out := ShowQuestion{
		TimeoutSeconds: int(gs.options.QuestionTimeout.Seconds()),
		ID:             question.Id,
		Categorie:      cat,
		Question:       instance.ToClient(),
	}

	gs.questionTimer.Reset(gs.options.QuestionTimeout)

	return out
}

func (r *Room) handleAnswer(a Answer, player serial) Events {
	isValid := r.game.isAnswerValid(a)

	// we defer the state update to the end of the question
	r.game.currentAnswers[player] = isValid

	return r.tryEndQuestion(false) // wait for other players if needed
}

func (r *Room) handleWantNextTurn(event WantNextTurn, player Player) Events {
	g := &r.game
	g.currentWantNextTurn[player.ID] = true

	pReview := &r.players[player.ID].advance.review
	if event.MarkQuestion {
		pReview.MarkedQuestions = append(pReview.MarkedQuestions, g.question.ID)
	}

	return r.tryEndTurn()
}

type playerAdvance struct {
	review  QuestionReview
	success Success
}

func (r *Room) state() GameState {
	out := GameState{
		Players:    make(map[serial]PlayerStatus),
		PawnTile:   r.game.pawnTile,
		PlayerTurn: r.game.playerTurn,
	}
	for _, pl := range r.players {
		out.Players[pl.pl.ID] = PlayerStatus{
			Name:       pl.pl.Pseudo,
			Review:     pl.advance.review,
			Success:    pl.advance.success,
			IsInactive: pl.conn == nil,
		}
	}
	return out
}
