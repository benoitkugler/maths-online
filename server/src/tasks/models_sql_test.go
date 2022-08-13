package tasks

import (
	"database/sql"
	"testing"

	"github.com/benoitkugler/maths-online/prof/editor"
	"github.com/benoitkugler/maths-online/prof/teacher"
	tu "github.com/benoitkugler/maths-online/utils/testutils"
)

func createEx(t *testing.T, db *sql.DB, idTeacher teacher.IdTeacher) (editor.Exercice, editor.ExerciceQuestions) {
	qu, err := editor.Question{IdTeacher: 1}.Insert(db)
	tu.Assert(t, err == nil)

	ex, err := editor.Exercice{IdTeacher: 1}.Insert(db)
	tu.Assert(t, err == nil)

	tx, err := db.Begin()
	tu.Assert(t, err == nil)

	questions := editor.ExerciceQuestions{
		{IdExercice: ex.Id, IdQuestion: qu.Id, Index: 0, Bareme: 2},
		{IdExercice: ex.Id, IdQuestion: qu.Id, Index: 1, Bareme: 1},
		{IdExercice: ex.Id, IdQuestion: qu.Id, Index: 2, Bareme: 3},
	}

	err = editor.InsertManyExerciceQuestions(tx, questions...)
	tu.Assert(t, err == nil)

	err = tx.Commit()
	tu.Assert(t, err == nil)

	return ex, questions
}

func TestProgression(t *testing.T) {
	db := tu.NewTestDB(t, "../prof/teacher/gen_create.sql", "../prof/editor/gen_create.sql", "gen_create.sql")
	defer db.Remove()

	_, err := teacher.Teacher{IsAdmin: true}.Insert(db)
	tu.Assert(t, err == nil)

	cl, err := teacher.Classroom{IdTeacher: 1}.Insert(db)
	tu.Assert(t, err == nil)

	student, err := teacher.Student{IdClassroom: cl.Id}.Insert(db)
	tu.Assert(t, err == nil)

	ex, _ := createEx(t, db.DB, 1)

	task, err := Task{IdExercice: ex.Id}.Insert(db)
	tu.Assert(t, err == nil)

	tx, err := db.Begin()
	tu.Assert(t, err == nil)

	prog, err := Progression{IdTask: task.Id, IdStudent: student.Id, IdExercice: ex.Id}.Insert(tx)
	tu.Assert(t, err == nil)

	err = InsertManyProgressionQuestions(tx,
		ProgressionQuestion{IdProgression: prog.Id, IdExercice: prog.IdExercice, Index: 0, History: QuestionHistory{false, true}},
		ProgressionQuestion{IdProgression: prog.Id, IdExercice: prog.IdExercice, Index: 2, History: QuestionHistory{}},
	)
	tu.Assert(t, err == nil)
	err = tx.Commit()
	tu.Assert(t, err == nil)

	out, err := loadProgressions(db, Progressions{prog.Id: prog})
	tu.Assert(t, err == nil)
	if out[prog.Id].NextQuestion != 1 {
		t.Fatal(out)
	}
}
