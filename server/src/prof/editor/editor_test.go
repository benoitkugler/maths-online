package editor

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/benoitkugler/maths-online/prof/teacher"
	"github.com/benoitkugler/maths-online/utils/testutils"
)

func TestInstantiateQuestions(t *testing.T) {
	db, err := testutils.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", testutils.DB, err)
		return
	}

	ct := NewController(db, teacher.Teacher{})
	out, err := ct.InstantiateQuestions([]int64{24, 29, 37})
	if err != nil {
		t.Fatal(err)
	}
	s, _ := json.MarshalIndent(out, " ", " ")
	fmt.Println(string(s))
}

func TestExerciceCRUD(t *testing.T) {
	db := testutils.CreateDBDev(t, "../teacher/gen_create.sql", "gen_create.sql")
	defer testutils.RemoveDBDev()
	defer db.Close()

	_, err := teacher.Teacher{IsAdmin: true}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	ct := NewController(db, teacher.Teacher{Id: 1})

	ex, err := ct.createExercice(1)
	if err != nil {
		t.Fatal(err)
	}

	l, err := ct.addQuestionToExercice(ExerciceAddQuestionIn{IdExercice: ex.Exercice.Id, IdQuestion: -1}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 1 {
		t.Fatal(err)
	}

	qu, err := Question{IdTeacher: 1}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	l, err = ct.addQuestionToExercice(ExerciceAddQuestionIn{IdExercice: ex.Exercice.Id, IdQuestion: qu.Id}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 2 {
		t.Fatal(err)
	}

	l, err = ct.removeQuestionFromExercice(ExerciceRemoveQuestionIn{IdExercice: ex.Exercice.Id, IdQuestion: qu.Id}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 1 {
		t.Fatal(err)
	}

	exe, err := ct.updateExercice(Exercice{Id: ex.Exercice.Id, Description: "test", Title: "test2", Flow: Sequencial}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if exe.Flow != Sequencial {
		t.Fatal(exe)
	}

	err = ct.deleteExercice(ex.Exercice.Id, true, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDB(t *testing.T) {
	db, err := testutils.DB.ConnectPostgres()
	if err != nil {
		t.Skip("DB not available")
	}

	ct := NewController(db, teacher.Teacher{Id: 1})
	_, err = ct.getExercices(1)
	if err != nil {
		t.Fatal(err)
	}
}
