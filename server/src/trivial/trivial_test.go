package trivial

import (
	"testing"
	"time"

	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func TestPanics(t *testing.T) {
	tu.ShouldPanic(t, func() { _ = (pGameOver + 1).String() })

	tu.ShouldPanic(t, func() {
		r := Room{game: game{phase: pGameOver + 1}, players: map[PlayerID]*playerConn{"": {}}}
		r.removePlayer(Player{})
	})
}

func TestSummary(t *testing.T) {
	r := NewRoom("", Options{Launch: LaunchStrategy{Max: 3}, Questions: exPool, QuestionTimeout: time.Minute})
	go r.Listen()

	_ = r.Options()

	if sum := r.Summary(); len(sum.Successes) != 0 || sum.PlayerTurn != nil {
		t.Fatal(sum)
	}

	r.mustJoin(t, "p1")
	r.mustJoin(t, "p2")

	if sum := r.Summary(); len(sum.Successes) != 2 || sum.PlayerTurn != nil {
		t.Fatal(sum)
	}

	r.mustJoin(t, "p3") // trigger a game start

	sum := r.Summary()
	tu.Assert(t, len(sum.Successes) == 3 && sum.PlayerTurn.ID == "p1")
	tu.Assert(t, sum.LatestQuestion.ID == 0)

	r.throwAndMove("p1")

	time.Sleep(10 * time.Millisecond)

	tu.Assert(t, r.Summary().LatestQuestion.ID != 0)
}

func TestReview(t *testing.T) {
	r := NewRoom("", Options{Launch: LaunchStrategy{Max: 2}, Questions: exPool, QuestionTimeout: time.Minute})

	naturalEnding := make(chan bool)
	go func() {
		_, isNaturalEnd := r.Listen()
		naturalEnding <- isNaturalEnd
	}()

	r.mustJoin(t, "p1")
	r.mustJoin(t, "p2")

	r.throwAndMove("p1")

	r.Event <- ClientEvent{Event: Answer{}, Player: "p1"}
	r.Event <- ClientEvent{Event: Answer{}, Player: "p2"}

	time.Sleep(time.Millisecond)

	r.Event <- ClientEvent{Event: WantNextTurn{MarkQuestion: true}, Player: "p1"}
	r.Event <- ClientEvent{Event: WantNextTurn{}, Player: "p3"}

	r.Terminate <- true

	isNat := <-naturalEnding
	tu.Assert(t, !isNat)

	history := r.replay().QuestionHistory[Player{ID: "p1", Pseudo: ""}]
	tu.Assert(t, len(history.MarkedQuestions) == 1 && len(history.QuestionHistory) == 1)
}
