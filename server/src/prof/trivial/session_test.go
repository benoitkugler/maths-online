package trivial

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/server/src/pass"
	ed "github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	tv "github.com/benoitkugler/maths-online/server/src/trivial"
	"github.com/benoitkugler/maths-online/server/src/utils/testutils"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
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

	session.createGameEvents <- createGame{ID: "g1", Options: tv.Options{Launch: tv.LaunchStrategy{Max: 2}, Questions: dummyQuestions}}
	session.createGameEvents <- createGame{ID: "g2", Options: tv.Options{Launch: tv.LaunchStrategy{Max: 2}, Questions: dummyQuestions}}

	time.Sleep(time.Millisecond * 10)

	session.lock.Lock()
	if L := len(session.games); L != 2 {
		t.Fatal(L)
	}
	session.lock.Unlock()

	for _, s := range session.collectSummaries() {
		testutils.Assert(t, s.LatestQuestion.ID == 0)
	}

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

func TestSessionManual(t *testing.T) {
	db, err := tu.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", tu.DB, err)
		return
	}

	ct := NewController(db, pass.Encrypter{}, "", teacher.Teacher{Id: 1})

	l, err := ct.getTrivialPoursuits(1)
	tu.AssertNoErr(t, err)

	config := l[0]

	out, err := ct.launchConfig(LaunchSessionIn{IdConfig: config.Config.Id, Groups: GroupsStrategyManual{4}}, 1)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(out.GameIDs) == 4)

	session := ct.getSession(1)
	err = session.startGame(out.GameIDs[0])
	tu.Assert(t, err != nil) // no players

	out, err = ct.launchConfig(LaunchSessionIn{IdConfig: config.Config.Id, Groups: GroupsStrategyAuto{[]int{1, 2, 3}}}, 1)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(out.GameIDs) == 3)
}
