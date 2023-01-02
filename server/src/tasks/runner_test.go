package tasks

import (
	"database/sql"
	"testing"

	"github.com/benoitkugler/maths-online/server/src/maths/questions"
	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	ed "github.com/benoitkugler/maths-online/server/src/sql/editor"
	ta "github.com/benoitkugler/maths-online/server/src/sql/tasks"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func TestInstantiateQuestions(t *testing.T) {
	db, err := tu.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", tu.DB, err)
		return
	}

	out, err := InstantiateQuestions(db, []ed.IdQuestion{24, 29, 37})
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(out) == 3)
	// s, _ := json.MarshalIndent(out, " ", " ")
	// fmt.Println(string(s)) // may be used as reference for client tests
}

func createEx(t *testing.T, db *sql.DB, idTeacher teacher.IdTeacher) (ed.Exercice, ed.ExerciceQuestions, ta.Monoquestion) {
	group, err := ed.Exercicegroup{IdTeacher: idTeacher}.Insert(db)
	tu.AssertNoErr(t, err)

	ex, err := ed.Exercice{IdGroup: group.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	qu1, err := ed.Question{
		NeedExercice: ex.Id.AsOptional(),
		Page: questions.QuestionPage{
			Enonce: questions.Enonce{
				questions.NumberFieldBlock{Expression: "1"},
			},
		},
	}.Insert(db)
	tu.AssertNoErr(t, err)

	qu2, err := qu1.Insert(db)
	tu.AssertNoErr(t, err)
	qu3, err := qu1.Insert(db)
	tu.AssertNoErr(t, err)

	tx, err := db.Begin()
	tu.AssertNoErr(t, err)

	qus := ed.ExerciceQuestions{
		{IdExercice: ex.Id, IdQuestion: qu1.Id, Index: 0, Bareme: 2},
		{IdExercice: ex.Id, IdQuestion: qu2.Id, Index: 1, Bareme: 1},
		{IdExercice: ex.Id, IdQuestion: qu3.Id, Index: 2, Bareme: 3},
	}

	err = ed.InsertManyExerciceQuestions(tx, qus...)
	tu.AssertNoErr(t, err)

	// monoquestion
	quGroup, err := ed.Questiongroup{IdTeacher: idTeacher}.Insert(tx)
	tu.AssertNoErr(t, err)
	qu4, err := ed.Question{
		IdGroup: quGroup.Id.AsOptional(),
		Page: questions.QuestionPage{
			Enonce: questions.Enonce{
				questions.NumberFieldBlock{Expression: "1"},
			},
		},
	}.Insert(tx)
	tu.AssertNoErr(t, err)

	mono, err := ta.Monoquestion{IdQuestion: qu4.Id, NbRepeat: 3, Bareme: 2}.Insert(tx)
	tu.AssertNoErr(t, err)

	err = tx.Commit()
	tu.AssertNoErr(t, err)

	return ex, qus, mono
}

func TestEvaluateExercice(t *testing.T) {
	db := tu.NewTestDB(t, "../sql/teacher/gen_create.sql", "../sql/editor/gen_create.sql", "../sql/tasks/gen_create.sql")
	defer db.Remove()

	tc, err := teacher.Teacher{IsAdmin: true}.Insert(db)
	tu.AssertNoErr(t, err)

	ex, questions, monoquestion := createEx(t, db.DB, tc.Id)

	progExt := ProgressionExt{
		Questions: make([]ta.QuestionHistory, len(questions)),
	}

	// no error since the exercice is parallel
	_, err = EvaluateWorkIn{
		ID:          newWorkIDFromMono(monoquestion.Id),
		Progression: progExt,
		Answers:     map[int]AnswerP{},
	}.Evaluate(db)
	tu.AssertNoErr(t, err)

	out, err := EvaluateWorkIn{
		ID:          newWorkIDFromEx(ex.Id),
		Progression: progExt,
		Answers: map[int]AnswerP{
			0: {Answer: client.QuestionAnswersIn{Data: client.Answers{0: client.NumberAnswer{Value: 22}}}},
		},
	}.Evaluate(db)
	tu.AssertNoErr(t, err)
	tu.Assert(t, out.Progression.NextQuestion == 0) // wrong answer

	out, err = EvaluateWorkIn{
		ID:          newWorkIDFromEx(ex.Id),
		Progression: progExt,
		Answers: map[int]AnswerP{
			0: {Answer: client.QuestionAnswersIn{Data: client.Answers{0: client.NumberAnswer{Value: 1}}}},
		},
	}.Evaluate(db)
	tu.AssertNoErr(t, err)
	tu.Assert(t, out.Progression.NextQuestion == 1) // correct answer
}

func TestProgression(t *testing.T) {
	db := tu.NewTestDB(t, "../sql/teacher/gen_create.sql", "../sql/editor/gen_create.sql", "../sql/tasks/gen_create.sql")
	defer db.Remove()

	tc, err := teacher.Teacher{IsAdmin: true}.Insert(db)
	tu.AssertNoErr(t, err)

	cl, err := teacher.Classroom{IdTeacher: tc.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	student, err := teacher.Student{IdClassroom: cl.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	// test with exercice
	ex, questions, _ := createEx(t, db.DB, tc.Id)

	task, err := ta.Task{IdExercice: ex.Id.AsOptional()}.Insert(db)
	tu.AssertNoErr(t, err)

	prog, err := ta.Progression{IdTask: task.Id, IdStudent: student.Id}.Insert(db)
	tu.AssertNoErr(t, err)

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
	tu.AssertNoErr(t, err)

	out, err := loadProgressions(db, ta.Progressions{prog.Id: prog})
	tu.AssertNoErr(t, err)
	tu.Assert(t, out[prog.Id].NextQuestion == 1)

	// test with monoquestion
	monoquestion, err := ta.Monoquestion{IdQuestion: questions[0].IdQuestion, NbRepeat: 3}.Insert(db)
	tu.AssertNoErr(t, err)

	task, err = ta.Task{IdMonoquestion: monoquestion.Id.AsOptional()}.Insert(db)
	tu.AssertNoErr(t, err)

	prog, err = loadOrCreateProgressionFor(db, task.Id, student.Id)
	tu.AssertNoErr(t, err)

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
	tu.AssertNoErr(t, err)

	out, err = loadProgressions(db, ta.Progressions{prog.Id: prog})
	tu.AssertNoErr(t, err)
	tu.Assert(t, out[prog.Id].NextQuestion == 1)
}
