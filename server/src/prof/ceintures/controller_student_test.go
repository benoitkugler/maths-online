package ceintures

import (
	"database/sql"
	"math/rand"
	"testing"

	"github.com/benoitkugler/maths-online/server/src/maths/questions"
	"github.com/benoitkugler/maths-online/server/src/pass"
	ce "github.com/benoitkugler/maths-online/server/src/sql/ceintures"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/benoitkugler/maths-online/server/src/tasks"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func TestGetEvolution(t *testing.T) {
	db := tu.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../sql/editor/gen_create.sql", "../../sql/ceintures/gen_create.sql")
	defer db.Remove()

	tc, _ := teacher.Teacher{FavoriteMatiere: teacher.Francais}.Insert(db)
	cl, _ := teacher.Classroom{IdTeacher: tc.Id}.Insert(db)
	student, err := teacher.Student{IdClassroom: cl.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	ct := NewController(db.DB, teacher.Teacher{Id: 1}, pass.Encrypter{})
	// anonymous connection
	_, has, err := ct.getEvolution(StudentTokens{})
	tu.AssertNoErr(t, err)
	tu.Assert(t, !has)

	out, err := ct.createEvolution(CreateEvolutionIn{Level: ce.PostBac})
	tu.AssertNoErr(t, err)
	tu.Assert(t, out.AnonymousID != "")

	get, has, err := ct.getEvolution(StudentTokens{AnonymousID: out.AnonymousID})
	tu.AssertNoErr(t, err)
	tu.Assert(t, has)
	tu.Assert(t, len(get.Pending) != 0)
	tu.Assert(t, get.SuggestionIndex != -1)

	// registred connection
	id := ct.studentKey.EncryptID(int64(student.Id))
	_, has, err = ct.getEvolution(StudentTokens{ClientID: id})
	tu.AssertNoErr(t, err)
	tu.Assert(t, !has)

	out, err = ct.createEvolution(CreateEvolutionIn{ClientID: id, Level: ce.PostBac})
	tu.AssertNoErr(t, err)
	tu.Assert(t, out.AnonymousID == "")

	get, has, err = ct.getEvolution(StudentTokens{ClientID: id})
	tu.AssertNoErr(t, err)
	tu.Assert(t, has)
	tu.Assert(t, len(get.Pending) != 0)
	tu.Assert(t, get.SuggestionIndex != -1)
}

// populate each stage
func createQuestions(t *testing.T, db *sql.DB) {
	tx, err := db.Begin()
	tu.AssertNoErr(t, err)

	for d := ce.Domain(0); d < ce.NbDomains; d++ {
		for r := ce.Rank(0); r < ce.NbRanks; r++ {
			i := 1 + rand.Intn(3)
			for n := 0; n < i; n++ {
				_, err = ce.Beltquestion{Domain: d, Rank: r, Enonce: questions.Enonce{
					questions.TextBlock{Parts: "1+1="},
					questions.RadioFieldBlock{
						Answer: "1",
						Proposals: []questions.Interpolated{
							"La bonne réponse !",
							"La mauvaise..",
						},
						AsDropDown: false,
					},
				}}.Insert(tx)
				tu.AssertNoErr(t, err)
			}
		}
	}

	err = tx.Commit()
	tu.AssertNoErr(t, err)
}

func TestInitDevQuestions(t *testing.T) {
	// t.Skip()

	db, err := tu.DB.ConnectPostgres()
	tu.AssertNoErr(t, err)

	createQuestions(t, db)
}

func TestSelectEvaluateQuestions(t *testing.T) {
	db := tu.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../sql/editor/gen_create.sql", "../../sql/ceintures/gen_create.sql")
	defer db.Remove()
	createQuestions(t, db.DB)

	ct := NewController(db.DB, teacher.Teacher{Id: 1}, pass.Encrypter{})

	out, err := ct.selectQuestions(SelectQuestionsIn{Stage: Stage{ce.CalculMental, ce.Blanche}})
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(out.Questions) == 3)

	for d := ce.Domain(0); d < ce.NbDomains; d++ {
		out, err = ct.selectQuestions(SelectQuestionsIn{Stage: Stage{d, ce.Rouge}})
		tu.AssertNoErr(t, err)
		tu.Assert(t, len(out.Questions) == 4)
	}

	ev, _ := ct.createEvolution(CreateEvolutionIn{Level: ce.Seconde})
	stage := Stage{ce.CalculMental, ce.Blanche}
	out, _ = ct.selectQuestions(SelectQuestionsIn{Stage: stage})

	var ids []ce.IdBeltquestion
	for _, qu := range out.Questions {
		ids = append(ids, qu.Id)
	}
	res, err := ct.evaluateAnswers(EvaluateAnswersIn{
		Tokens:    StudentTokens{AnonymousID: ev.AnonymousID},
		Stage:     stage,
		Questions: ids,
		Answers:   make([]tasks.AnswerP, len(ids)),
	})
	tu.AssertNoErr(t, err)
	tu.Assert(t, res.Evolution.Advance == ce.Advance{})                                                                     // incorrect answer
	tu.Assert(t, res.Evolution.Stats[stage.Domain][stage.Rank] == ce.Stat{Success: 0, Failure: uint16(len(out.Questions))}) // incorrect answer
}
