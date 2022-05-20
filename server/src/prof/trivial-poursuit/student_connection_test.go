package trivialpoursuit

import (
	"testing"

	"github.com/benoitkugler/maths-online/pass"
)

func TestController_setupStudentClientDemo(t *testing.T) {
	creds := pass.DB{
		Host:     "localhost",
		User:     "benoit",
		Password: "dummy",
		Name:     "isyro_prod",
	}
	db, err := creds.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", creds, err)
		return
	}

	ct := NewController(db, pass.Encrypter{}, "1234")
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
	if session.playerIDs[out.PlayerID] != -1 {
		t.Fatal()
	}
}

func TestController_setupStudentClient(t *testing.T) {
	creds := pass.DB{
		Host:     "localhost",
		User:     "benoit",
		Password: "dummy",
		Name:     "isyro_prod",
	}
	db, err := creds.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", creds, err)
		return
	}

	ct := NewController(db, pass.Encrypter{}, "1234")

	gs, err := ct.createGameSession("7894", TrivialConfig{Id: -1, Questions: demoQuestions, QuestionTimeout: 120}, RandomGroupStrategy{2, 2})
	if err != nil {
		t.Fatal(err)
	}

	out, err := ct.setupStudentClient(gs.id, "", "")
	if err != nil {
		t.Fatal(err)
	}

	if len(ct.sessions) != 1 {
		t.Fatal()
	}
	session := ct.sessions[gs.id]
	if len(session.games) != 1 {
		t.Fatal()
	}
	if session.playerIDs[out.PlayerID] != -1 {
		t.Fatal()
	}
}
