package trivial

import (
	"context"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/server/src/sql/events"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

type noOpSuccesHandler struct{}

func (noOpSuccesHandler) OnQuestion(player PlayerID, correct, hasStreak3 bool) events.EventNotification {
	return events.EventNotification{}
}

func (noOpSuccesHandler) OnWin(player PlayerID) events.EventNotification {
	return events.EventNotification{}
}

func TestPanics(t *testing.T) {
	tu.ShouldPanic(t, func() { _ = (pGameOver + 1).String() })

	tu.ShouldPanic(t, func() {
		r := Room{game: game{phase: pGameOver + 1}, players: map[PlayerID]*playerConn{"": {}}}
		r.removePlayer(Player{})
	})
}

func TestSummary(t *testing.T) {
	r := NewRoom("", Options{Launch: LaunchStrategy{Max: 3}, Questions: exPool, QuestionTimeout: time.Minute}, noOpSuccesHandler{})
	go r.Listen(context.Background())

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
	tu.Assert(t, len(r.Summary().InQuestionStudents) == 0)

	r.throwAndMove("p1")

	time.Sleep(10 * time.Millisecond)

	tu.Assert(t, r.Summary().LatestQuestion.ID != 0)
	tu.Assert(t, len(r.Summary().InQuestionStudents) == 3)

	r.Event <- ClientEvent{Event: Answer{}, Player: "p1"}

	time.Sleep(10 * time.Millisecond)

	tu.Assert(t, len(r.Summary().InQuestionStudents) == 2)

	r.Event <- ClientEvent{Event: Answer{}, Player: "p2"}
	r.Event <- ClientEvent{Event: Answer{}, Player: "p3"}

	time.Sleep(10 * time.Millisecond)

	tu.Assert(t, len(r.Summary().InQuestionStudents) == 3)
}

func TestReview(t *testing.T) {
	r := NewRoom("", Options{Launch: LaunchStrategy{Max: 2}, Questions: exPool, QuestionTimeout: time.Minute}, noOpSuccesHandler{})

	naturalEnding := make(chan bool)
	go func() {
		_, isNaturalEnd := r.Listen(context.Background())
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
