package tasks

import (
	"database/sql"
	"testing"

	ed "github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func createEx(t *testing.T, db *sql.DB, idTeacher teacher.IdTeacher) (ed.Exercice, ed.ExerciceQuestions, ed.Questiongroup) {
	group, err := ed.Exercicegroup{IdTeacher: idTeacher}.Insert(db)
	tu.AssertNoErr(t, err)

	ex, err := ed.Exercice{IdGroup: group.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	qu1, err := ed.Question{NeedExercice: ex.Id.AsOptional()}.Insert(db)
	tu.AssertNoErr(t, err)
	qu2, err := ed.Question{NeedExercice: ex.Id.AsOptional()}.Insert(db)
	tu.AssertNoErr(t, err)
	qu3, err := ed.Question{NeedExercice: ex.Id.AsOptional()}.Insert(db)
	tu.AssertNoErr(t, err)

	tx, err := db.Begin()
	tu.AssertNoErr(t, err)

	questions := ed.ExerciceQuestions{
		{IdExercice: ex.Id, IdQuestion: qu1.Id, Index: 0, Bareme: 2},
		{IdExercice: ex.Id, IdQuestion: qu2.Id, Index: 1, Bareme: 1},
		{IdExercice: ex.Id, IdQuestion: qu3.Id, Index: 2, Bareme: 3},
	}

	err = ed.InsertManyExerciceQuestions(tx, questions...)
	tu.AssertNoErr(t, err)

	err = tx.Commit()
	tu.AssertNoErr(t, err)

	qg, err := ed.Questiongroup{IdTeacher: idTeacher}.Insert(db)
	tu.AssertNoErr(t, err)

	_, err = ed.Question{IdGroup: qg.Id.AsOptional()}.Insert(db)
	tu.AssertNoErr(t, err)
	_, err = ed.Question{IdGroup: qg.Id.AsOptional()}.Insert(db)
	tu.AssertNoErr(t, err)
	_, err = ed.Question{IdGroup: qg.Id.AsOptional()}.Insert(db)
	tu.AssertNoErr(t, err)

	return ex, questions, qg
}

func TestTaskConstraint(t *testing.T) {
	db := tu.NewTestDB(t, "../teacher/gen_create.sql", "../editor/gen_create.sql", "gen_create.sql")
	defer db.Remove()

	tc, err := teacher.Teacher{IsAdmin: true, FavoriteMatiere: teacher.Mathematiques}.Insert(db)
	tu.AssertNoErr(t, err)

	ex, questions, _ := createEx(t, db.DB, tc.Id)

	mono, err := Monoquestion{IdQuestion: questions[0].IdQuestion, NbRepeat: 1}.Insert(db)
	tu.AssertNoErr(t, err)

	// exactly one target must be given
	_, err = Task{}.Insert(db)
	tu.Assert(t, err != nil)
	_, err = Task{IdExercice: ex.Id.AsOptional(), IdMonoquestion: mono.Id.AsOptional()}.Insert(db)
	tu.Assert(t, err != nil)
}

func TestResizeProgression(t *testing.T) {
	db := tu.NewTestDB(t, "../teacher/gen_create.sql", "../editor/gen_create.sql", "gen_create.sql")
	defer db.Remove()

	tc, err := teacher.Teacher{IsAdmin: true, FavoriteMatiere: teacher.Mathematiques}.Insert(db)
	tu.AssertNoErr(t, err)
	cl, err := teacher.Classroom{}.Insert(db)
	tu.AssertNoErr(t, err)
	st, err := teacher.Student{IdClassroom: cl.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	_, _, group := createEx(t, db.DB, tc.Id)

	mono, err := RandomMonoquestion{IdQuestiongroup: group.Id, NbRepeat: 3}.Insert(db)
	tu.AssertNoErr(t, err)

	task, err := Task{IdRandomMonoquestion: mono.Id.AsOptional()}.Insert(db)
	tu.AssertNoErr(t, err)

	tx, err := db.Begin()
	tu.AssertNoErr(t, err)
	err = InsertManyProgressions(tx,
		Progression{IdStudent: st.Id, IdTask: task.Id, Index: 0},
		Progression{IdStudent: st.Id, IdTask: task.Id, Index: 1},
		Progression{IdStudent: st.Id, IdTask: task.Id, Index: 2},
	)
	tu.AssertNoErr(t, err)
	err = tx.Commit()
	tu.AssertNoErr(t, err)

	err = ResizeProgressions(db, task.Id, 1)
	tu.AssertNoErr(t, err)

	pgs, err := SelectAllProgressions(db)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(pgs) == 1)
}
