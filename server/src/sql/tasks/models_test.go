package tasks

import (
	"database/sql"
	"testing"

	ed "github.com/benoitkugler/maths-online/sql/editor"
	"github.com/benoitkugler/maths-online/sql/teacher"
	tu "github.com/benoitkugler/maths-online/utils/testutils"
)

func createEx(t *testing.T, db *sql.DB, idTeacher teacher.IdTeacher) (ed.Exercice, ed.ExerciceQuestions) {
	group, err := ed.Exercicegroup{IdTeacher: idTeacher}.Insert(db)
	tu.Assert(t, err == nil)

	ex, err := ed.Exercice{IdGroup: group.Id}.Insert(db)
	tu.Assert(t, err == nil)

	qu1, err := ed.Question{NeedExercice: ex.Id.AsOptional()}.Insert(db)
	tu.Assert(t, err == nil)
	qu2, err := ed.Question{NeedExercice: ex.Id.AsOptional()}.Insert(db)
	tu.Assert(t, err == nil)
	qu3, err := ed.Question{NeedExercice: ex.Id.AsOptional()}.Insert(db)
	tu.Assert(t, err == nil)

	tx, err := db.Begin()
	tu.Assert(t, err == nil)

	questions := ed.ExerciceQuestions{
		{IdExercice: ex.Id, IdQuestion: qu1.Id, Index: 0, Bareme: 2},
		{IdExercice: ex.Id, IdQuestion: qu2.Id, Index: 1, Bareme: 1},
		{IdExercice: ex.Id, IdQuestion: qu3.Id, Index: 2, Bareme: 3},
	}

	err = ed.InsertManyExerciceQuestions(tx, questions...)
	tu.Assert(t, err == nil)

	err = tx.Commit()
	tu.Assert(t, err == nil)

	return ex, questions
}

func TestTaskConstraint(t *testing.T) {
	db := tu.NewTestDB(t, "../teacher/gen_create.sql", "../editor/gen_create.sql", "gen_create.sql")
	defer db.Remove()

	tc, err := teacher.Teacher{IsAdmin: true}.Insert(db)
	tu.Assert(t, err == nil)

	ex, questions := createEx(t, db.DB, tc.Id)

	mono, err := Monoquestion{IdQuestion: questions[0].IdQuestion}.Insert(db)
	tu.Assert(t, err == nil)

	// exactly one target must be given
	_, err = Task{}.Insert(db)
	tu.Assert(t, err != nil)
	_, err = Task{IdExercice: ex.Id.AsOptional(), IdMonoquestion: mono.Id.AsOptional()}.Insert(db)
	tu.Assert(t, err != nil)
}
