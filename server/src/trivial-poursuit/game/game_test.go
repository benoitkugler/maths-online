package game

import (
	"reflect"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/maths/exercice/client"
	"github.com/benoitkugler/maths-online/prof/editor"
)

var exQu = WeigthedQuestions{
	Questions: []editor.Question{{Id: 1}, {Id: 2}},
	Weights:   []float64{1. / 2, 1. / 2},
}

func playersFromSuccess(scs ...Success) map[int]*PlayerStatus {
	out := make(map[int]*PlayerStatus)
	for i, s := range scs {
		out[i] = &PlayerStatus{Success: s}
	}
	return out
}

func Test_gameState_winners(t *testing.T) {
	tests := []struct {
		sc      map[int]*PlayerStatus
		wantOut []int
	}{
		{playersFromSuccess(Success{true}, Success{true}, Success{true, true, true, true, true}), []int{2}},
		{playersFromSuccess(Success{true}, Success{true}, Success{true, true, true, true}), nil},
		{playersFromSuccess(Success{true, true, true, true, true}, Success{true}, Success{true, true, true, true, true}), []int{0, 2}},
	}
	for _, tt := range tests {
		gs := &Game{
			GameState: GameState{
				Players: tt.sc,
			},
		}
		if gotOut := gs.winners(); !reflect.DeepEqual(gotOut, tt.wantOut) {
			t.Errorf("gameState.winners() = %v, want %v", gotOut, tt.wantOut)
		}
	}
}

func TestGameState_nextPlayer(t *testing.T) {
	type fields struct {
		Successes map[PlayerID]*PlayerStatus
		Player    int
	}
	tests := []struct {
		fields fields
		want   PlayerID
	}{
		{
			fields{
				Successes: map[int]*PlayerStatus{0: {}, 1: {}, 4: {}},
				Player:    0,
			},
			1,
		},
		{
			fields{
				Successes: map[int]*PlayerStatus{0: {}, 1: {}, 4: {}},
				Player:    1,
			},
			4,
		},
		{
			fields{
				Successes: map[int]*PlayerStatus{0: {}, 1: {}, 4: {}},
				Player:    2,
			},
			4,
		},
		{
			fields{
				Successes: map[int]*PlayerStatus{0: {}, 1: {}, 4: {}},
				Player:    4,
			},
			0,
		},
	}
	for _, tt := range tests {
		g := &Game{
			GameState: GameState{
				Players: tt.fields.Successes,
				Player:  tt.fields.Player,
			},
		}
		if got := g.nextPlayer(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("GameState.nextPlayer() = %v, want %v", got, tt.want)
		}
	}
}

func TestStart(t *testing.T) {
	g := NewGame(0, true, QuestionPool{})
	p1, _ := g.AddPlayer("")
	if p1 != 0 {
		t.Fatalf("unexpected player id %d", p1)
	}

	p2, _ := g.AddPlayer("")
	p3, _ := g.AddPlayer("")
	if p1 == p2 || p2 == p3 {
		t.Fatal()
	}

	if g.NumberPlayers(true) != 3 {
		t.Fatal()
	}

	g.RemovePlayer(p2)
	g.StartGame()
	if g.Player != p1 {
		t.Fatalf("unexpected first player %d", g.Player)
	}

	g.RemovePlayer(p3)
	if g.NumberPlayers(false) != 2 {
		t.Fatal()
	}
}

func TestDisconnect(t *testing.T) {
	g := NewGame(0, true, QuestionPool{})
	p1, _ := g.AddPlayer("")
	p2, _ := g.AddPlayer("")

	g.StartGame()
	events := g.RemovePlayer(p1).Events
	if len(events) != 2 {
		t.Fatal()
	}
	resetTurn := events[1].(PlayerTurn)
	if resetTurn.Player != p2 {
		t.Fatal(resetTurn)
	}
}

func TestDisconnectInQuestion(t *testing.T) {
	g := NewGame(time.Second/2, true, QuestionPool{exQu, exQu, exQu, exQu, exQu})

	p1, _ := g.AddPlayer("")
	p2, _ := g.AddPlayer("")

	g.EmitQuestion()

	g.handleAnswer(Answer{}, p1)

	// player1 has answered, waiting for player2

	g.RemovePlayer(p2)

	if g.QuestionTimeout.Stop() {
		t.Fatal("question should have been closed")
	}
}

func TestDisconnectInBeforeNextTurn(t *testing.T) {
	g := NewGame(time.Second/2, true, QuestionPool{exQu, exQu, exQu, exQu, exQu})

	p1, _ := g.AddPlayer("")
	p2, _ := g.AddPlayer("")
	p3, _ := g.AddPlayer("")
	g.StartGame()

	g.EmitQuestion()

	g.handleAnswer(Answer{}, p1)
	g.handleAnswer(Answer{}, p2)
	g.handleAnswer(Answer{}, p3)

	g.handleWantNextTurn(WantNextTurn{}, p1)
	g.handleWantNextTurn(WantNextTurn{}, p3)

	// player1 wants next turn, waiting for player2
	if g.GameState.Player != 0 {
		t.Fatal(g.GameState.Player)
	}

	g.RemovePlayer(p2)

	if g.GameState.Player != 2 {
		t.Fatal(g.GameState.Player)
	}
}

func TestEmitQuestion(t *testing.T) {
	g := NewGame(time.Second/2, true, QuestionPool{exQu, exQu, exQu, exQu, exQu})

	g.AddPlayer("")

	if g.QuestionTimeout.Stop() {
		t.Fatal("timer should not being running")
	}

	g.EmitQuestion()

outer:
	for {
		select {
		case <-g.QuestionTimeout.C:
			g.QuestionTimeoutAction()
			break outer
		}
	}

	g.EmitQuestion()
	g.handleAnswer(Answer{client.QuestionAnswersIn{}}, 0)
	g.endQuestion(false)
	if g.QuestionTimeout.Stop() {
		t.Fatal("timer should have been stopped")
	}
}

func TestHandleClientEvent(t *testing.T) {
	g := NewGame(0, true, QuestionPool{exQu, exQu, exQu, exQu, exQu})
	g.AddPlayer("")
	g.StartGame()

	// check nextTurn is properly reset
	g.currentWantNextTurn = map[int]bool{2: true}
	g.startTurn()
	if len(g.currentWantNextTurn) != 0 {
		t.Fatal("currentWantNextTurn should be reset")
	}

	up, isOver, err := g.HandleClientEvent(ClientEvent{Event: Ping{}})
	if err != nil {
		t.Fatal(err)
	}
	if up != nil || isOver {
		t.Fatal("Ping should be ignored")
	}

	_, _, err = g.HandleClientEvent(ClientEvent{})
	if err == nil {
		t.Fatal("expected error for invalid client event type")
	}

	_, _, err = g.HandleClientEvent(ClientEvent{Event: DiceClicked{}, Player: 2})
	if err == nil {
		t.Fatal("expected error for invalid click")
	}

	g.HandleClientEvent(ClientEvent{Event: DiceClicked{}})

	_, _, err = g.HandleClientEvent(ClientEvent{Event: ClientMove{}, Player: 2})
	if err == nil {
		t.Fatal("expected error for invalid move")
	}

	_, _, err = g.HandleClientEvent(ClientEvent{Event: ClientMove{Tile: 89}})
	if err == nil {
		t.Fatal("expected error for invalid tile")
	}

	expected := g.dice.Face
	_, _, err = g.HandleClientEvent(ClientEvent{Event: ClientMove{Tile: int(g.dice.Face)}})
	if err != nil {
		t.Fatal(err)
	}
	if g.PawnTile != int(expected) {
		t.Fatalf("expected %d, got %d", expected, g.PawnTile)
	}
	if g.dice != (diceThrow{}) { // dice is reset on move
		t.Fatal()
	}

	_, _, err = g.HandleClientEvent(ClientEvent{Event: Answer{client.QuestionAnswersIn{}}})
	if err != nil {
		t.Fatal(err)
	}

	for cat := categorie(0); cat < nbCategories; cat++ {
		g.EmitQuestion()
		g.question.categorie = cat
		_, _, err = g.HandleClientEvent(ClientEvent{Event: Answer{client.QuestionAnswersIn{}}})
		if err != nil {
			t.Fatal(err)
		}
	}

	if !reflect.DeepEqual(g.winners(), []int{0}) {
		t.Fatal()
	}

	if !g.IsPlaying() {
		t.Fatal("game should still be active")
	}

	g.HandleClientEvent(ClientEvent{Player: 0, Event: WantNextTurn{MarkQuestion: true}})

	if g.IsPlaying() {
		t.Fatal("game should be over")
	}
}

func TestRemovePieOnTimeout(t *testing.T) {
	g := NewGame(0, true, QuestionPool{exQu, exQu, exQu, exQu, exQu})
	g.AddPlayer("")
	g.StartGame()
	g.Players[0].Success = Success{true, true, true, true, true}

	qu := g.EmitQuestion()
	time.Sleep(10 * time.Millisecond)
	g.QuestionTimeoutAction()

	if g.Players[0].Success[qu.Categorie] {
		t.Fatal("success must have been lost")
	}
	if len(g.Players[0].Review.QuestionHistory) != 1 {
		t.Fatal("missing question in history")
	}
}
