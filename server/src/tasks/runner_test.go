package tasks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/benoitkugler/maths-online/maths/questions"
	"github.com/benoitkugler/maths-online/maths/questions/client"
	ed "github.com/benoitkugler/maths-online/sql/editor"
	ta "github.com/benoitkugler/maths-online/sql/tasks"
	"github.com/benoitkugler/maths-online/sql/teacher"
	"github.com/benoitkugler/maths-online/utils/testutils"
	tu "github.com/benoitkugler/maths-online/utils/testutils"
)

func TestInstantiateQuestions(t *testing.T) {
	db, err := tu.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", tu.DB, err)
		return
	}

	out, err := InstantiateQuestions(db, []ed.IdQuestion{24, 29, 37})
	tu.Assert(t, err == nil)
	s, _ := json.MarshalIndent(out, " ", " ")
	fmt.Println(string(s)) // may be used as reference for client tests
}

func createEx(t *testing.T, db *sql.DB, idTeacher teacher.IdTeacher) (ed.Exercice, ed.ExerciceQuestions) {
	group, err := ed.Exercicegroup{IdTeacher: idTeacher}.Insert(db)
	tu.Assert(t, err == nil)

	ex, err := ed.Exercice{IdGroup: group.Id, Flow: ed.Parallel}.Insert(db)
	tu.Assert(t, err == nil)

	qu, err := ed.Question{
		NeedExercice: ex.Id.AsOptional(),
		Page: questions.QuestionPage{
			Enonce: questions.Enonce{
				questions.NumberFieldBlock{Expression: "1"},
			},
		},
	}.Insert(db)
	tu.Assert(t, err == nil)

	tx, err := db.Begin()
	tu.Assert(t, err == nil)

	questions := ed.ExerciceQuestions{
		{IdExercice: ex.Id, IdQuestion: qu.Id, Index: 0, Bareme: 2},
		{IdExercice: ex.Id, IdQuestion: qu.Id, Index: 1, Bareme: 1},
		{IdExercice: ex.Id, IdQuestion: qu.Id, Index: 2, Bareme: 3},
	}

	err = ed.InsertManyExerciceQuestions(tx, questions...)
	tu.Assert(t, err == nil)

	err = tx.Commit()
	tu.Assert(t, err == nil)

	return ex, questions
}

func TestEvaluateExercice(t *testing.T) {
	db := testutils.NewTestDB(t, "../sql/teacher/gen_create.sql", "../sql/editor/gen_create.sql", "../sql/tasks/gen_create.sql")
	defer db.Remove()

	tc, err := teacher.Teacher{IsAdmin: true}.Insert(db)
	tu.Assert(t, err == nil)

	ex, questions := createEx(t, db.DB, tc.Id)

	progExt := ProgressionExt{
		Questions: make([]ta.QuestionHistory, len(questions)),
	}

	// no error since the exercice is parallel
	_, err = EvaluateWorkIn{
		ID:          newWorkIDFromEx(ex.Id),
		Progression: progExt,
		Answers:     map[int]Answer{},
	}.Evaluate(db)
	tu.Assert(t, err == nil)

	out, err := EvaluateWorkIn{
		ID:          newWorkIDFromEx(ex.Id),
		Progression: progExt,
		Answers: map[int]Answer{
			0: {Answer: client.QuestionAnswersIn{Data: client.Answers{0: client.NumberAnswer{Value: 22}}}},
		},
	}.Evaluate(db)
	tu.Assert(t, err == nil)
	tu.Assert(t, out.Progression.NextQuestion == 0) // wrong answer

	out, err = EvaluateWorkIn{
		ID:          newWorkIDFromEx(ex.Id),
		Progression: progExt,
		Answers: map[int]Answer{
			0: {Answer: client.QuestionAnswersIn{Data: client.Answers{0: client.NumberAnswer{Value: 1}}}},
		},
	}.Evaluate(db)
	tu.Assert(t, err == nil)
	tu.Assert(t, out.Progression.NextQuestion == 1) // correct answer
}

func TestProgression(t *testing.T) {
	db := tu.NewTestDB(t, "../sql/teacher/gen_create.sql", "../sql/editor/gen_create.sql", "../sql/tasks/gen_create.sql")
	defer db.Remove()

	tc, err := teacher.Teacher{IsAdmin: true}.Insert(db)
	tu.Assert(t, err == nil)

	cl, err := teacher.Classroom{IdTeacher: tc.Id}.Insert(db)
	tu.Assert(t, err == nil)

	student, err := teacher.Student{IdClassroom: cl.Id}.Insert(db)
	tu.Assert(t, err == nil)

	// test with exercice
	ex, questions := createEx(t, db.DB, tc.Id)

	task, err := ta.Task{IdExercice: ex.Id.AsOptional()}.Insert(db)
	tu.Assert(t, err == nil)

	prog, err := ta.Progression{IdTask: task.Id, IdStudent: student.Id}.Insert(db)
	tu.Assert(t, err == nil)

	err = updateProgression(db.DB, prog, []ta.QuestionHistory{
		{false, true},
		{},
	})
	tu.Assert(t, err != nil) // invalid number of questions
	err = updateProgression(db.DB, prog, []ta.QuestionHistory{
		{false, true},
		{},
		{},
	})
	tu.Assert(t, err == nil)

	out, err := loadProgressions(db, ta.Progressions{prog.Id: prog})
	tu.Assert(t, err == nil)
	tu.Assert(t, out[prog.Id].NextQuestion == 1)

	// test with monoquestion
	monoquestion, err := ta.Monoquestion{IdQuestion: questions[0].IdQuestion, NbRepeat: 3}.Insert(db)
	tu.Assert(t, err == nil)

	task, err = ta.Task{IdMonoquestion: monoquestion.Id.AsOptional()}.Insert(db)
	tu.Assert(t, err == nil)

	prog, err = loadOrCreateProgressionFor(db, task.Id, student.Id)
	tu.Assert(t, err == nil)

	err = updateProgression(db.DB, prog, []ta.QuestionHistory{
		{false, true},
		{},
	})
	tu.Assert(t, err != nil) // invalid number of questions
	err = updateProgression(db.DB, prog, []ta.QuestionHistory{
		{false, true},
		{},
		{},
	})
	tu.Assert(t, err == nil)

	out, err = loadProgressions(db, ta.Progressions{prog.Id: prog})
	tu.Assert(t, err == nil)
	tu.Assert(t, out[prog.Id].NextQuestion == 1)
}
