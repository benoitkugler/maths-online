package exercice

import (
	"testing"

	"github.com/benoitkugler/maths-online/prof/editor"
	"github.com/benoitkugler/maths-online/utils/testutils"
)

func TestCRUD(t *testing.T) {
	db := testutils.CreateDBDev(t, "../editor/gen_create.sql", "gen_create.sql")
	defer testutils.RemoveDBDev()
	defer db.Close()

	ex, err := randExercice().Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	qu1, err := editor.Question{}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}
	qu2, err := editor.Question{}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	err = InsertManyExerciceQuestions(tx,
		ExerciceQuestion{IdExercice: ex.Id, IdQuestion: qu1.Id, Bareme: 4},
		ExerciceQuestion{IdExercice: ex.Id, IdQuestion: qu2.Id, Bareme: 5},
	)
	if err != nil {
		t.Fatal(err)
	}

	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}
}
