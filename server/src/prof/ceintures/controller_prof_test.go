package ceintures

import (
	"testing"

	"github.com/benoitkugler/maths-online/server/src/maths/expression"
	"github.com/benoitkugler/maths-online/server/src/maths/questions"
	"github.com/benoitkugler/maths-online/server/src/pass"
	ce "github.com/benoitkugler/maths-online/server/src/sql/ceintures"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func TestCRUDQuestion(t *testing.T) {
	db := tu.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../sql/ceintures/gen_create.sql")
	defer db.Remove()

	ct := NewController(db.DB, teacher.Teacher{Id: 1}, pass.Encrypter{})

	stage := Stage{ce.Equations, ce.Blanche}
	l, err := ct.getQuestions(stage)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(l) == 0)

	qu, err := ct.createQuestion(stage)
	tu.AssertNoErr(t, err)

	qu.Parameters = questions.Parameters{
		questions.Co("Un commentaire"),
		questions.Rp{Expression: "2", Variable: expression.NewVar('a')},
		questions.Co("TODO !"),
	}
	qu.Enonce = questions.Enonce{questions.TextBlock{}}
	pr, err := ct.saveQuestion(SaveBeltQuestionIn{Question: qu})
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(pr.Preview.Questions) == 1)

	qu.Enonce = questions.Enonce{questions.TextBlock{}, questions.NumberFieldBlock{Expression: "x"}}
	pr, err = ct.saveQuestion(SaveBeltQuestionIn{Question: qu})
	tu.Assert(t, !pr.IsValid) // x is not defined

	dup, err := ct.duplicateQuestion(qu.Id)
	tu.AssertNoErr(t, err)
	tu.Assert(t, dup.Preview.QuestionIndex == 1)

	out, err := ct.getScheme(1)
	tu.AssertNoErr(t, err)
	tu.Assert(t, out.Stages[stage.Domain][stage.Rank].HasTODO)

	err = ct.deleteQuestion(qu.Id)
	tu.AssertNoErr(t, err)
	err = ct.deleteQuestion(dup.Question.Id)
	tu.AssertNoErr(t, err)

	l, err = ct.getQuestions(stage)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(l) == 0)
}
