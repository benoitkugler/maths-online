package trivial

import (
	"os"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/pass"
	"github.com/benoitkugler/maths-online/prof/teacher"
	tv "github.com/benoitkugler/maths-online/trivial-poursuit"
	"github.com/benoitkugler/maths-online/utils/testutils"
)

func TestCreateConfig(t *testing.T) {
	db, err := testutils.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", testutils.DB, err)
		return
	}

	out, err := TrivialConfig{
		QuestionTimeout: 120,
		ShowDecrassage:  true,
		IdTeacher:       1,
	}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := DeleteTrivialConfigById(db, out.Id); err != nil {
		t.Fatal(err)
	}
}

func TestGetConfig(t *testing.T) {
	db := testutils.CreateDBDev(t, "../teacher/gen_create.sql", "../editor/gen_create.sql", "gen_create.sql")
	defer testutils.RemoveDBDev()
	defer db.Close()

	user1, err := teacher.Teacher{Mail: "1"}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}
	user2, err := teacher.Teacher{Mail: "2"}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	_, err = TrivialConfig{IdTeacher: user1.Id}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	c2, err := TrivialConfig{IdTeacher: user2.Id}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	ct := NewController(db, pass.Encrypter{}, "", teacher.Teacher{})
	l, err := ct.getTrivialPoursuits(user1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 1 {
		t.Fatal()
	}

	c2.Public = true
	if _, err := c2.Update(db); err != nil {
		t.Fatal(err)
	}

	l, err = ct.getTrivialPoursuits(user1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 2 {
		t.Fatal()
	}
}

func TestGameTermination(t *testing.T) {
	tv.ProgressLogger.SetOutput(os.Stdout)

	ct := newGameSession("test", -1)

	ct.createGame(createGame{ID: "Game1"})

	if len(ct.games) != 1 {
		t.Fatal("expected one game")
	}
	ct.stopGame(stopGame{ID: "Game1"})

	time.Sleep(20 * time.Millisecond)

	ct.lock.Lock()
	if len(ct.games) != 0 {
		t.Fatal("game should have been removed")
	}
	ct.lock.Unlock()
}

func TestMissingQuestions(t *testing.T) {
	db, err := testutils.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", testutils.DB, err)
		return
	}

	ct := NewController(db, pass.Encrypter{}, "", teacher.Teacher{})

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
	out, err := ct.checkMissingQuestions(criteria, 1)
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
	out, err = ct.checkMissingQuestions(criteria, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Missing) == 0 {
		t.Fatal("categories should be missing")
	}
}

func TestGetTrivials(t *testing.T) {
	db, err := testutils.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", testutils.DB, err)
		return
	}

	ct := NewController(db, pass.Encrypter{}, "", teacher.Teacher{})

	for range [10]int{} {
		t.Run("", func(t *testing.T) {
			_, err := ct.getTrivialPoursuits(0)
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
