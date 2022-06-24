package editor

import (
	"testing"

	"github.com/benoitkugler/maths-online/maths/questions"
	"github.com/benoitkugler/maths-online/maths/questions/client"
	"github.com/benoitkugler/maths-online/prof/teacher"
	"github.com/benoitkugler/maths-online/utils/testutils"
)

func TestEvaluateExercice(t *testing.T) {
	db := testutils.CreateDBDev(t, "../teacher/gen_create.sql", "gen_create.sql")
	defer testutils.RemoveDBDev()
	defer db.Close()

	_, err := teacher.Teacher{IsAdmin: true}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	ct := NewController(db, teacher.Teacher{Id: 1})

	qu := Question{IdTeacher: 1, Page: questions.QuestionPage{Enonce: questions.Enonce{questions.NumberFieldBlock{Expression: "1"}}}}
	qu, err = qu.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	ex, err := ct.createExercice(1)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ct.updateQuestionsEx(ExerciceUpdateQuestionsIn{
		IdExercice: ex.Exercice.Id,
		Questions:  ExerciceQuestions{{IdQuestion: qu.Id}},
	}, 1)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ct.createQuestionEx(ExerciceCreateQuestionIn{IdExercice: ex.Exercice.Id}, 1)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ct.createQuestionEx(ExerciceCreateQuestionIn{IdExercice: ex.Exercice.Id}, 1)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ct.createQuestionEx(ExerciceCreateQuestionIn{IdExercice: ex.Exercice.Id}, 1)
	if err != nil {
		t.Fatal(err)
	}

	prog, err := Progression{IdExercice: ex.Exercice.Id}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	progExt, err := ct.fetchProgression(prog.Id)
	if err != nil {
		t.Fatal(err)
	}

	// no error since the exercice is parallel
	_, err = ct.evaluateExercice(EvaluateExerciceIn{IdExercice: ex.Exercice.Id, Progression: progExt, Answers: map[int]Answer{}})
	if err != nil {
		t.Fatal(err)
	}

	out, err := ct.evaluateExercice(EvaluateExerciceIn{IdExercice: ex.Exercice.Id, Progression: progExt, Answers: map[int]Answer{
		0: {Answer: client.QuestionAnswersIn{Data: client.Answers{0: client.NumberAnswer{Value: 22}}}},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if out.Progression.NextQuestion != 0 {
		t.Fatal(err)
	}

	out, err = ct.evaluateExercice(EvaluateExerciceIn{IdExercice: ex.Exercice.Id, Progression: progExt, Answers: map[int]Answer{
		0: {Answer: client.QuestionAnswersIn{Data: client.Answers{0: client.NumberAnswer{Value: 1}}}},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if out.Progression.NextQuestion != 1 {
		t.Fatal(err)
	}
}
