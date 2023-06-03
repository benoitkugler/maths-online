package trivial

import (
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/server/src/pass"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	tv "github.com/benoitkugler/maths-online/server/src/trivial"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func TestController_parseCode(t *testing.T) {
	const demoPin = "1234"
	gs := newGameStore(demoPin)

	for _, test := range []struct {
		code     string
		expected gameID
		wantErr  bool
	}{
		{"1234.12.2", demoCode{demoPin, "12", 2}, false},
		{"1234.12.ax", nil, true},
		{"1238.12.4", nil, true},
		{"1235.12", teacherCode{"1235", "12"}, false},
		{"7896.127", teacherCode{"7896", "127"}, false},
		{"7896.1", nil, true},
		{"12312", selfaccessCode("12312"), false},
		{"1238.12.4.8", nil, true},
		{"1234.abc.4", demoCode{demoPin, "abc", 4}, false},
		{"1234.12.1", demoCode{demoPin, "12", 1}, false},
		{"1234.1.1", nil, true},
		{"", nil, true},
		{"789456qsd", nil, true},
		{"1234.a", nil, true},
	} {
		got, err := gs.parseCode(test.code)
		tu.Assert(t, (err != nil) == test.wantErr)
		tu.Assert(t, got == test.expected)
	}
}

func TestController_setupStudentClientDemo(t *testing.T) {
	db, err := tu.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", tu.DB, err)
		return
	}

	ct := NewController(db, pass.Encrypter{}, "1234", teacher.Teacher{})
	out, err := ct.setupStudentClient("1234.12.2", "", "")
	tu.AssertNoErr(t, err)

	tu.Assert(t, len(ct.store.games) == 1)
	_, ok := ct.store.games[demoCode{"1234", "12", 2}]
	tu.Assert(t, ok)
	tu.Assert(t, ct.store.playerIDs[out.PlayerID] == demoCode{"1234", "12", 2})
}

func TestController_setupStudentClient(t *testing.T) {
	db, err := tu.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", tu.DB, err)
		return
	}

	ct := NewController(db, pass.Encrypter{}, "1234", teacher.Teacher{})

	sessionID := ct.store.getOrCreateSession(10)

	questionPool, err := selectQuestions(ct.db, demoQuestions, ct.admin.Id)
	tu.AssertNoErr(t, err)

	options := tv.Options{
		Launch:          tv.LaunchStrategy{Max: 2},
		QuestionTimeout: time.Second * 120,
		ShowDecrassage:  true,
		Questions:       questionPool,
	}

	gameID := ct.store.newTeacherGameID(sessionID)
	ct.store.createGame(createGame{
		ID:      gameID,
		Options: options,
	})

	time.Sleep(time.Millisecond)

	out, err := ct.setupStudentClient(string(gameID.roomID()), "", "")
	tu.AssertNoErr(t, err)

	tu.Assert(t, len(ct.store.games) == 1)
	tu.Assert(t, ct.store.playerIDs[out.PlayerID] == gameID)
}
