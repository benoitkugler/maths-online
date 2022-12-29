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
	tu.Assert(t, err == nil)
	tu.Assert(t, len(out) == 3)
	// s, _ := json.MarshalIndent(out, " ", " ")
	// fmt.Println(string(s)) // may be used as reference for client tests
}

func createEx(t *testing.T, db *sql.DB, idTeacher teacher.IdTeacher) (ed.Exercice, ed.ExerciceQuestions, ta.Monoquestion) {
	group, err := ed.Exercicegroup{IdTeacher: idTeacher}.Insert(db)
	tu.Assert(t, err == nil)

	ex, err := ed.Exercice{IdGroup: group.Id}.Insert(db)
	tu.Assert(t, err == nil)

	qu1, err := ed.Question{
		NeedExercice: ex.Id.AsOptional(),
		Page: questions.QuestionPage{
			Enonce: questions.Enonce{
				questions.NumberFieldBlock{Expression: "1"},
			},
		},
	}.Insert(db)
	tu.Assert(t, err == nil)

	qu2, err := qu1.Insert(db)
	tu.Assert(t, err == nil)
	qu3, err := qu1.Insert(db)
	tu.Assert(t, err == nil)

	tx, err := db.Begin()
	tu.Assert(t, err == nil)

	qus := ed.ExerciceQuestions{
		{IdExercice: ex.Id, IdQuestion: qu1.Id, Index: 0, Bareme: 2},
		{IdExercice: ex.Id, IdQuestion: qu2.Id, Index: 1, Bareme: 1},
		{IdExercice: ex.Id, IdQuestion: qu3.Id, Index: 2, Bareme: 3},
	}

	err = ed.InsertManyExerciceQuestions(tx, qus...)
	tu.Assert(t, err == nil)

	// monoquestion
	quGroup, err := ed.Questiongroup{IdTeacher: idTeacher}.Insert(tx)
	tu.Assert(t, err == nil)
	qu4, err := ed.Question{
		IdGroup: quGroup.Id.AsOptional(),
		Page: questions.QuestionPage{
			Enonce: questions.Enonce{
				questions.NumberFieldBlock{Expression: "1"},
			},
		},
	}.Insert(tx)
	tu.Assert(t, err == nil)

	mono, err := ta.Monoquestion{IdQuestion: qu4.Id, NbRepeat: 3, Bareme: 2}.Insert(tx)
	tu.Assert(t, err == nil)

	err = tx.Commit()
	tu.Assert(t, err == nil)

	return ex, qus, mono
}

func TestEvaluateExercice(t *testing.T) {
	db := tu.NewTestDB(t, "../sql/teacher/gen_create.sql", "../sql/editor/gen_create.sql", "../sql/tasks/gen_create.sql")
	defer db.Remove()

	tc, err := teacher.Teacher{IsAdmin: true}.Insert(db)
	tu.Assert(t, err == nil)

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
	tu.Assert(t, err == nil)

	out, err := EvaluateWorkIn{
		ID:          newWorkIDFromEx(ex.Id),
		Progression: progExt,
		Answers: map[int]AnswerP{
			0: {Answer: client.QuestionAnswersIn{Data: client.Answers{0: client.NumberAnswer{Value: 22}}}},
		},
	}.Evaluate(db)
	tu.Assert(t, err == nil)
	tu.Assert(t, out.Progression.NextQuestion == 0) // wrong answer

	out, err = EvaluateWorkIn{
		ID:          newWorkIDFromEx(ex.Id),
		Progression: progExt,
		Answers: map[int]AnswerP{
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
	ex, questions, _ := createEx(t, db.DB, tc.Id)

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
