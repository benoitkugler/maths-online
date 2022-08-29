package trivial

import (
	"context"
	"sort"
	"testing"
	"time"

	ed "github.com/benoitkugler/maths-online/sql/editor"
	tv "github.com/benoitkugler/maths-online/trivial"
)

var dummyQuestions = tv.QuestionPool{
	tv.WeigthedQuestions{Questions: []ed.Question{qu(1)}, Weights: []float64{1}},
	tv.WeigthedQuestions{Questions: []ed.Question{qu(1)}, Weights: []float64{1}},
	tv.WeigthedQuestions{Questions: []ed.Question{qu(1)}, Weights: []float64{1}},
	tv.WeigthedQuestions{Questions: []ed.Question{qu(1)}, Weights: []float64{1}},
	tv.WeigthedQuestions{Questions: []ed.Question{qu(1)}, Weights: []float64{1}},
}

func TestGameID(t *testing.T) {
	s := make([]string, 20)
	for i := range s {
		s[i] = string(gameIDFromSerial("test", i))
	}

	if !sort.StringsAreSorted(s) {
		t.Fatal("game ids are not sorted")
	}
}

func TestSession(t *testing.T) {
	session := newGameSession("test", -1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go session.mainLoop(ctx)

	session.createGameEvents <- createGame{ID: "g1", Options: tv.Options{PlayersNumber: 2, Questions: dummyQuestions}}
	session.createGameEvents <- createGame{ID: "g2", Options: tv.Options{PlayersNumber: 2, Questions: dummyQuestions}}

	time.Sleep(time.Millisecond * 10)

	session.lock.Lock()
	if L := len(session.games); L != 2 {
		t.Fatal(L)
	}
	session.lock.Unlock()

	// try to remove an inexisting game
	session.stopGameEvents <- stopGame{ID: "xxx"}
	time.Sleep(time.Millisecond * 10)

	session.lock.Lock()
	if L := len(session.games); L != 2 {
		t.Fatal(L)
	}
	session.lock.Unlock()

	session.stopGameEvents <- stopGame{ID: "g1"}
	time.Sleep(time.Millisecond * 10)

	session.lock.Lock()
	if L := len(session.games); L != 1 {
		t.Fatal(L)
	}
	session.lock.Unlock()

	// test restart
	session.stopGameEvents <- stopGame{ID: "g2", Restart: true}
	time.Sleep(time.Millisecond * 10)

	session.lock.Lock()
	if L := len(session.games); L != 1 {
		t.Fatal(L)
	}
	session.lock.Unlock()
}
