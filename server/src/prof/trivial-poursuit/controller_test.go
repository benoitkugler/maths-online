package trivialpoursuit

import (
	"os"
	"sort"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/pass"
	tv "github.com/benoitkugler/maths-online/trivial-poursuit"
	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
)

func TestGameTermination(t *testing.T) {
	tv.ProgressLogger.SetOutput(os.Stdout)

	ct := newGameSession("test", nil, TrivialConfig{}, RandomGroupStrategy{2, 2}, game.QuestionPool{})

	id := ct.createGame(2)

	if len(ct.games) != 1 {
		t.Fatal("expected one game")
	}
	ct.games[id].Terminate <- true

	time.Sleep(20 * time.Millisecond)

	ct.lock.Lock()
	if len(ct.games) != 0 {
		t.Fatal("game should have been removed")
	}
	ct.lock.Unlock()
}

func TestGameID(t *testing.T) {
	s := make([]string, 20)
	for i := range s {
		s[i] = gameIDFromSerial(i)
	}

	if !sort.StringsAreSorted(s) {
		t.Fatal("game ids are not sorted")
	}
}

func TestMissingQuestions(t *testing.T) {
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

	ct := NewController(db, pass.Encrypter{}, "")

	criteria := CategoriesQuestions{
		{
			{
				"POURCENTAGES",
				"VIOLET",
			},
			{
				"POURCENTAGES",
			},
		},
		{
			{
				"POURCENTAGES",
				"VERT",
			},
		},
		{
			{
				"POURCENTAGES",
				"ORANGE",
			},
		},
		{
			{
				"POURCENTAGES",
				"JAUNE",
			},
		},
		{
			{
				"POURCENTAGES",
				"BLEU",
			},
		},
	}
	out, err := ct.checkMissingQuestions(criteria)
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Missing) != 0 {
		t.Fatal()
	}

	criteria = CategoriesQuestions{
		{
			{
				"Pourcentages",
				"Valeur finale",
			},
		},
		{
			{
				"Pourcentages",
				"Taux réciproque",
			},
		},
		{
			{
				"Pourcentages",
				"Proportion",
			},
			{
				"Pourcentages",
				"Proportion de proportion",
			},
		},
		{
			{
				"Pourcentages",
				"Evolutions identiques",
			},
			{
				"Pourcentages",
				"Evolutions successives",
			},
		},
		{
			{
				"Pourcentages",
				"Coefficient multiplicateur",
			},
			{
				"Pourcentages",
				"Taux d'évolution",
			},
		},
	}
	out, err = ct.checkMissingQuestions(criteria)
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Missing) == 0 {
		t.Fatal("categories should be missing")
	}
}

func TestGetTrivials(t *testing.T) {
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

	ct := NewController(db, pass.Encrypter{}, "")

	for range [10]int{} {
		t.Run("", func(t *testing.T) {
			_, err := ct.getTrivialPoursuits()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestController_isDemoSessionID(t *testing.T) {
	const demoPin = "1234"
	tests := []struct {
		args          string
		wantRoom      string
		wantNbPlayers int
	}{
		{"1234.abc.4", "abc", 4},
		{"1234.12.1", "12", 1},
		{"1234.1.1", "", 0},
		{"", "", 0},
		{"789456qsd", "", 0},
		{"1234.a", "", 0},
	}
	for _, tt := range tests {
		ct := &Controller{
			demoPin: demoPin,
		}
		if gotRoom, gotNbPlayers := ct.isDemoSessionID(tt.args); gotRoom != tt.wantRoom || gotNbPlayers != tt.wantNbPlayers {
			t.Errorf("Controller.isDemoSessionID() = %v, want %v", gotNbPlayers, tt.wantNbPlayers)
		}
	}
}
