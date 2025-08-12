package trivial

import (
	"encoding/json"
	"io"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/server/src/pass"
	ed "github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/sql/events"
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
	var gs gameStore
	s := make([]string, 20)
	for i := range s {
		s[i] = string(gs.newTeacherGameID("test").String())
	}

	if !sort.StringsAreSorted(s) {
		t.Fatal("game ids are not sorted")
	}
}

func TestSession(t *testing.T) {
	gs := newGameStore(nil, pass.Encrypter{}, "test")

	gs.createGame(createGame{ID: teacherCode{"teacherXXX", "g1"}, Options: tv.Options{Launch: tv.LaunchStrategy{Max: 2}, Questions: dummyQuestions}})
	gs.createGame(createGame{ID: teacherCode{"teacherXXX", "g2"}, Options: tv.Options{Launch: tv.LaunchStrategy{Max: 2}, Questions: dummyQuestions}})

	time.Sleep(time.Millisecond * 10)

	gs.lock.Lock()
	if L := len(gs.games); L != 2 {
		t.Fatal(L)
	}
	gs.lock.Unlock()

	for _, s := range gs.collectSummaries("teacherXXX") {
		testutils.Assert(t, s.LatestQuestion.ID == 0)
	}

	// try to remove an inexisting game
	gs.stopGame(teacherCode{"teacherXXX", "xxx"}, false)

	time.Sleep(time.Millisecond * 10)

	gs.lock.Lock()
	if L := len(gs.games); L != 2 {
		t.Fatal(L)
	}
	gs.lock.Unlock()

	gs.stopGame(teacherCode{"teacherXXX", "g1"}, false)
	time.Sleep(time.Millisecond * 10)

	gs.lock.Lock()
	if L := len(gs.games); L != 1 {
		t.Fatal(L)
	}
	gs.lock.Unlock()

	// test restart
	gs.stopGame(teacherCode{"teacherXXX", "g2"}, true)
	time.Sleep(time.Millisecond * 10)

	gs.lock.Lock()
	if L := len(gs.games); L != 1 {
		t.Fatal(L)
	}
	gs.lock.Unlock()
}

func TestSessionManual(t *testing.T) {
	db, err := tu.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", tu.DB, err)
		return
	}

	ct := NewController(db, pass.Encrypter{}, "", teacher.Teacher{Id: 1})

	l, err := ct.getTrivialPoursuits(1, teacher.Mathematiques)
	tu.AssertNoErr(t, err)

	config := l[0]

	out, err := ct.launchConfig(LaunchSessionIn{IdConfig: config.Config.Id, Groups: GroupsStrategyManual{4}}, 1)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(out.GameIDs) == 4)

	gID, err := ct.store.parseCode(string(out.GameIDs[0]))
	tu.AssertNoErr(t, err)

	err = ct.store.startGame(gID)
	tu.Assert(t, err != nil) // no players

	out, err = ct.launchConfig(LaunchSessionIn{IdConfig: config.Config.Id, Groups: GroupsStrategyAuto{[]int{1, 2, 3}}}, 1)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(out.GameIDs) == 3)
}

func TestMonitor(t *testing.T) {
	db, err := tu.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", tu.DB, err)
		return
	}

	ct := NewController(db, pass.Encrypter{}, "", teacher.Teacher{Id: 1})

	l, err := ct.getTrivialPoursuits(1, teacher.Mathematiques)
	tu.AssertNoErr(t, err)

	config := l[0]

	_, err = ct.launchConfig(LaunchSessionIn{IdConfig: config.Config.Id, Groups: GroupsStrategyManual{4}}, 1)
	tu.AssertNoErr(t, err)

	session := ct.store.getSessionID(1)
	sums := ct.store.collectSummaries(session)
	tu.Assert(t, len(sums) == 4)
}

type clientOut struct {
	updates []tv.StateUpdate
}

func (c *clientOut) WriteJSON(v interface{}) error {
	err := json.NewEncoder(io.Discard).Encode(v)
	c.updates = append(c.updates, v.(tv.StateUpdate))
	return err
}

func (c *clientOut) lastU(lock *sync.Mutex) tv.StateUpdate {
	lock.Lock()
	defer lock.Unlock()
	return c.updates[len(c.updates)-1]
}

func TestAdvance(t *testing.T) {
	tv.GameStartDelay = time.Millisecond

	db := tu.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../sql/events/gen_create.sql")
	defer db.Remove()

	key := pass.Encrypter{}

	gs := newGameStore(db.DB, key, "")

	cl, err := teacher.Classroom{}.Insert(db)
	tu.AssertNoErr(t, err)
	student, err := teacher.Student{IdClassroom: cl.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	encID := key.EncryptID(int64(student.Id))

	gameID := selfaccessCode("g1")
	gs.createGame(createGame{ID: gameID, Options: tv.Options{
		Questions:       dummyQuestions,
		Launch:          tv.LaunchStrategy{Max: 1},
		QuestionTimeout: time.Minute,
	}})
	game := gs.games[gameID]

	meta, err := gs.setupStudent(encID, gameID, key)
	tu.AssertNoErr(t, err)

	var out clientOut
	err = game.Join(tv.Player{ID: meta.PlayerID}, &out)
	tu.AssertNoErr(t, err)

	playTurn := func() events.Events {
		game.Event <- tv.ClientEvent{Player: meta.PlayerID, Event: tv.DiceClicked{}}
		time.Sleep(time.Millisecond)

		tile := out.lastU(&gs.lock).Events[1].(tv.PossibleMoves).Tiles[0]
		game.Event <- tv.ClientEvent{Player: meta.PlayerID, Event: tv.ClientMove{Tile: tile}}
		time.Sleep(time.Millisecond)

		game.Event <- tv.ClientEvent{Player: meta.PlayerID, Event: tv.Answer{}}
		time.Sleep(time.Millisecond)

		game.Event <- tv.ClientEvent{Player: meta.PlayerID, Event: tv.WantNextTurn{}}

		evList, err := events.SelectEventsByIdStudents(db, student.Id)
		tu.AssertNoErr(t, err)

		return evList
	}

	evList := playTurn()

	// the questions are empty so the answer if always true
	tu.Assert(t, len(evList) == 1)
	tu.Assert(t, evList[0].Event == events.E_All_QuestionRight)

	evList = playTurn()
	tu.Assert(t, len(evList) == 2)

	evList = playTurn()
	tu.Assert(t, len(evList) == 4) // streak

	evList = playTurn()
	tu.Assert(t, len(evList) == 5)
}
