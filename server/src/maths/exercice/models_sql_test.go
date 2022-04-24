package exercice

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/benoitkugler/maths-online/utils/testutils"
)

func TestRoot(t *testing.T) {
	// create a DB shared by all tests
	db := testutils.CreateDBDev(t, "create_gen.sql")
	defer testutils.RemoveDBDev()
	defer db.Close()

	// t.Run("CRUD for Exercice", func(t *testing.T) { testExercice(t, db) })
	t.Run("CRUD for Question", func(t *testing.T) { testQuestion(t, db) })
}

func testQuestion(t *testing.T, db *sql.DB) {
	questions, err := SelectAllQuestions(db)
	if err != nil {
		t.Fatal(err)
	}
	L := len(questions)

	qu := randQuestion()
	qu, err = qu.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	questions, err = SelectAllQuestions(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(questions) != L+1 {
		t.Fatal()
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	err = InsertManyQuestionTags(tx, QuestionTag{IdQuestion: qu.Id, Tag: "seconde"}, QuestionTag{IdQuestion: qu.Id, Tag: "calcul"})
	if err != nil {
		t.Fatal(err)
	}
	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}

	tags, err := SelectAllTags(db)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(tags, []string{"calcul", "seconde"}) {
		t.Fatal()
	}

	_, err = DeleteQuestionById(db, qu.Id)
	if err != nil {
		t.Fatal(err)
	}
}
