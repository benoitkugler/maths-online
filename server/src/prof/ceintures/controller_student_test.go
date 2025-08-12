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
	db := tu.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../sql/ceintures/gen_create.sql")
	defer db.Remove()

	cl, _ := teacher.Classroom{}.Insert(db)
	tc, _ := teacher.Teacher{FavoriteMatiere: teacher.Francais}.Insert(db)
	_ = teacher.TeacherClassroom{IdTeacher: tc.Id, IdClassroom: cl.Id}.Insert(db)
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

	l, err := ct.getStudentsAdvance(cl.Id, tc.Id)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(l) == 1)
}

// populate each stage
func createQuestions(t *testing.T, db *sql.DB) (questionsNumber [ce.NbDomains][ce.NbRanks]int) {
	tx, err := db.Begin()
	tu.AssertNoErr(t, err)

	for d := ce.Domain(0); d < ce.NbDomains; d++ {
		for r := ce.Rank(0); r < ce.NbRanks; r++ {
			i := 1 + rand.Intn(3)
			questionsNumber[d][r] = i
			for n := 0; n < i; n++ {
				_, err = ce.Beltquestion{Domain: d, Rank: r, Repeat: 1, Enonce: questions.Enonce{
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

	return questionsNumber
}

func TestInitDevQuestions(t *testing.T) {
	t.Skip()

	db, err := tu.DB.ConnectPostgres()
	tu.AssertNoErr(t, err)

	createQuestions(t, db)
}

func TestSelectEvaluateQuestions(t *testing.T) {
	db := tu.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../sql/editor/gen_create.sql", "../../sql/ceintures/gen_create.sql")
	defer db.Remove()
	questionsNumber := createQuestions(t, db.DB)

	ct := NewController(db.DB, teacher.Teacher{Id: 1}, pass.Encrypter{})

	out, err := ct.selectQuestions(SelectQuestionsIn{Stage: Stage{ce.CalculMentalI, ce.Blanche}})
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(out.Questions) == questionsNumber[ce.CalculMentalI][ce.Blanche])

	for d := ce.Domain(0); d < ce.NbDomains; d++ {
		stage := Stage{d, ce.Rouge}
		out, err = ct.selectQuestions(SelectQuestionsIn{Stage: stage})
		tu.AssertNoErr(t, err)
		tu.Assert(t, len(out.Questions) == questionsNumber[stage.Domain][stage.Rank])
	}

	ev, err := ct.createEvolution(CreateEvolutionIn{Level: ce.Seconde})
	tu.AssertNoErr(t, err)
	stage := Stage{ce.CalculMentalI, ce.Blanche}
	out, err = ct.selectQuestions(SelectQuestionsIn{Stage: stage})
	tu.AssertNoErr(t, err)

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
