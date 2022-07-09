package trivial

import (
	"testing"
	"time"
)

func TestSummary(t *testing.T) {
	r := NewRoom("", Options{PlayersNumber: 3, Questions: exPool, QuestionTimeout: time.Minute})
	go r.Listen()

	if sum := r.Summary(); len(sum.Successes) != 0 || sum.PlayerTurn != nil {
		t.Fatal(sum)
	}

	r.mustJoin(t, "p1")
	r.mustJoin(t, "p2")

	if sum := r.Summary(); len(sum.Successes) != 2 || sum.PlayerTurn != nil {
		t.Fatal(sum)
	}

	r.mustJoin(t, "p3") // trigger a game start

	if sum := r.Summary(); len(sum.Successes) != 3 || sum.PlayerTurn.ID != "p1" {
		t.Fatal(sum)
	}
}

func TestReview(t *testing.T) {
	r := NewRoom("", Options{PlayersNumber: 2, Questions: exPool, QuestionTimeout: time.Minute})
	r.game.dice = DiceThrow{1}

	naturalEnding := make(chan bool)
	go func() {
		_, isNaturalEnd := r.Listen()
		naturalEnding <- isNaturalEnd
	}()

	r.mustJoin(t, "p1")
	r.mustJoin(t, "p2")

	r.Event <- ClientEvent{Event: ClientMove{Tile: 1}, Player: "p1"}

	r.Event <- ClientEvent{Event: Answer{}, Player: "p1"}
	r.Event <- ClientEvent{Event: Answer{}, Player: "p2"}

	time.Sleep(time.Millisecond)

	r.Event <- ClientEvent{Event: WantNextTurn{MarkQuestion: true}, Player: "p1"}
	r.Event <- ClientEvent{Event: WantNextTurn{}, Player: "p3"}

	r.Terminate <- true

	if isNat := <-naturalEnding; isNat {
		t.Fatal()
	}

	history := r.replay().QuestionHistory[Player{ID: "p1", Pseudo: ""}]
	if len(history.MarkedQuestions) != 1 || len(history.QuestionHistory) != 1 {
		t.Fatal(history)
	}
}
