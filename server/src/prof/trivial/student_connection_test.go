package trivial

import (
	"context"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/server/src/pass"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	tv "github.com/benoitkugler/maths-online/server/src/trivial"
	"github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func TestController_setupStudentClientDemo(t *testing.T) {
	db, err := testutils.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", testutils.DB, err)
		return
	}

	ct := NewController(db, pass.Encrypter{}, "1234", teacher.Teacher{})
	out, err := ct.setupStudentClient("1234.12.2", "", "")
	if err != nil {
		t.Fatal(err)
	}

	if len(ct.sessions) != 1 {
		t.Fatal()
	}
	session := ct.sessions["1234.12.2"]
	if len(session.games) != 1 {
		t.Fatal()
	}
	if session.playerIDs[out.PlayerID] != "1234.12.2" {
		t.Fatal()
	}
}

func TestController_setupStudentClient(t *testing.T) {
	db, err := testutils.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", testutils.DB, err)
		return
	}

	ct := NewController(db, pass.Encrypter{}, "1234", teacher.Teacher{})

	gs := ct.createSession("7894", -1)
	if len(ct.sessions) != 1 {
		t.Fatal()
	}

	go gs.mainLoop(context.Background())

	questionPool, err := selectQuestions(ct.db, demoQuestions, ct.admin.Id)
	if err != nil {
		t.Fatal(err)
	}

	options := tv.Options{
		PlayersNumber:   2,
		QuestionTimeout: time.Second * 120,
		ShowDecrassage:  true,
		Questions:       questionPool,
	}

	gameID := gs.newGameID()
	gs.createGameEvents <- createGame{
		ID:      gameID,
		Options: options,
	}

	time.Sleep(time.Millisecond)

	out, err := ct.setupStudentClient(string(gameID), "", "")
	if err != nil {
		t.Fatal(err)
	}

	session := ct.sessions[gs.id]
	if len(session.games) != 1 {
		t.Fatal()
	}
	if session.playerIDs[out.PlayerID] != gameID {
		t.Fatal()
	}
}
