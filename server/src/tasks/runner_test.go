package tasks

import (
	"database/sql"
	"reflect"
	"sort"
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

func createEx(t *testing.T, db *sql.DB, idTeacher teacher.IdTeacher) (ed.Exercice, ed.ExerciceQuestions, ta.Monoquestion, ta.RandomMonoquestion) {
	group, err := ed.Exercicegroup{IdTeacher: idTeacher}.Insert(db)
	tu.AssertNoErr(t, err)

	ex, err := ed.Exercice{IdGroup: group.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	qu1, err := ed.Question{
		NeedExercice: ex.Id.AsOptional(),
		Enonce: questions.Enonce{
			questions.NumberFieldBlock{Expression: "1"},
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
		Enonce: questions.Enonce{
			questions.NumberFieldBlock{Expression: "1"},
		},
	}.Insert(tx)
	tu.AssertNoErr(t, err)

	_, err = ed.Question{IdGroup: quGroup.Id.AsOptional()}.Insert(tx)
	tu.AssertNoErr(t, err)
	_, err = ed.Question{IdGroup: quGroup.Id.AsOptional()}.Insert(tx)
	tu.AssertNoErr(t, err)
	_, err = ed.Question{IdGroup: quGroup.Id.AsOptional()}.Insert(tx)
	tu.AssertNoErr(t, err)
	_, err = ed.Question{IdGroup: quGroup.Id.AsOptional()}.Insert(tx)
	tu.AssertNoErr(t, err)
	_, err = ed.Question{IdGroup: quGroup.Id.AsOptional()}.Insert(tx)
	tu.AssertNoErr(t, err)

	mono, err := ta.Monoquestion{IdQuestion: qu4.Id, NbRepeat: 3, Bareme: 2}.Insert(tx)
	tu.AssertNoErr(t, err)

	// randommonoquestion
	randomMono, err := ta.RandomMonoquestion{IdQuestiongroup: quGroup.Id, NbRepeat: 4, Bareme: 2}.Insert(tx)
	tu.AssertNoErr(t, err)

	err = tx.Commit()
	tu.AssertNoErr(t, err)

	return ex, qus, mono, randomMono
}

func TestEvaluateExercice(t *testing.T) {
	db := tu.NewTestDB(t, "../sql/teacher/gen_create.sql", "../sql/editor/gen_create.sql", "../sql/tasks/gen_create.sql")
	defer db.Remove()

	tc, err := teacher.Teacher{IsAdmin: true, FavoriteMatiere: teacher.Mathematiques}.Insert(db)
	tu.AssertNoErr(t, err)

	ex, questions, monoquestion, _ := createEx(t, db.DB, tc.Id)

	prog := Progression(make([]ta.QuestionHistory, len(questions)))

	// no error since the exercice is parallel
	_, err = EvaluateWorkIn{
		ID:          newWorkIDFromMono(monoquestion.Id),
		Progression: prog,
	}.Evaluate(db, -1, false)
	tu.AssertNoErr(t, err)

	out, err := EvaluateWorkIn{
		ID:          newWorkIDFromEx(ex.Id),
		Progression: prog,
		AnswerIndex: 0,
		Answer:      AnswerP{Answer: client.QuestionAnswersIn{Data: client.Answers{0: client.NumberAnswer{Value: 22}}}},
	}.Evaluate(db, -1, false)
	tu.AssertNoErr(t, err)
	tu.Assert(t, out.Progression.NextQuestion == 0) // wrong answer

	out, err = EvaluateWorkIn{
		ID:          newWorkIDFromEx(ex.Id),
		Progression: prog,
		AnswerIndex: 0,
		Answer:      AnswerP{Answer: client.QuestionAnswersIn{Data: client.Answers{0: client.NumberAnswer{Value: 22}}}},
	}.Evaluate(db, -1, true)
	tu.AssertNoErr(t, err)
	tu.Assert(t, out.Progression.NextQuestion == 1) // wrong answer in one try mode

	out, err = EvaluateWorkIn{
		ID:          newWorkIDFromEx(ex.Id),
		Progression: prog,
		AnswerIndex: 0,
		Answer:      AnswerP{Answer: client.QuestionAnswersIn{Data: client.Answers{0: client.NumberAnswer{Value: 1}}}},
	}.Evaluate(db, -1, false)
	tu.AssertNoErr(t, err)
	tu.Assert(t, out.Progression.NextQuestion == 1) // correct answer
}

func Test_inferNextQuestion(t *testing.T) {
	for _, test := range []struct {
		progression Progression
		isOneTry    bool
		exp         int
	}{
		{Progression{{}}, false, 0},
		{Progression{{true}, {}}, false, 1},
		{Progression{{true}, {}, {true}}, false, 1},
		{Progression{{true}, {false}}, false, 1},
		{Progression{{true}, {false}}, true, -1},
		{Progression{{true}, {false}, {}}, true, 2},
		{Progression{{true}, {true}}, false, -1},
		{Progression{{true}, {true}}, true, -1},
	} {
		tu.Assert(t, test.progression.inferNextQuestion(test.isOneTry) == test.exp)
	}
}

func TestProgression(t *testing.T) {
	db := tu.NewTestDB(t, "../sql/teacher/gen_create.sql", "../sql/editor/gen_create.sql", "../sql/tasks/gen_create.sql")
	defer db.Remove()

	tc, err := teacher.Teacher{IsAdmin: true, FavoriteMatiere: teacher.Mathematiques}.Insert(db)
	tu.AssertNoErr(t, err)

	cl, err := teacher.Classroom{}.Insert(db)
	tu.AssertNoErr(t, err)

	student, err := teacher.Student{IdClassroom: cl.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	// test with exercice
	ex, _, monoquestion, randomMono := createEx(t, db.DB, tc.Id)

	task, err := ta.Task{IdExercice: ex.Id.AsOptional()}.Insert(db)
	tu.AssertNoErr(t, err)

	err = updateProgression(db.DB, student.Id, task.Id, []ta.QuestionHistory{
		{false, true},
		{},
	})
	tu.Assert(t, err != nil) // invalid number of questions

	err = updateProgression(db.DB, student.Id, task.Id, []ta.QuestionHistory{
		{false, true},
		{},
		{},
	})
	tu.AssertNoErr(t, err)

	out, err := LoadTasksProgression(db, student.Id, []ta.IdTask{task.Id})
	tu.AssertNoErr(t, err)
	tu.Assert(t, out[task.Id].HasProgression)

	// test with monoquestion
	task, err = ta.Task{IdMonoquestion: monoquestion.Id.AsOptional()}.Insert(db.DB)
	tu.AssertNoErr(t, err)

	err = updateProgression(db.DB, student.Id, task.Id, []ta.QuestionHistory{
		{false, true},
		{},
	})
	tu.Assert(t, err != nil) // invalid number of questions
	err = updateProgression(db.DB, student.Id, task.Id, []ta.QuestionHistory{
		{false, true},
		{},
		{},
	})
	tu.AssertNoErr(t, err)

	out, err = LoadTasksProgression(db, student.Id, []ta.IdTask{task.Id})
	tu.AssertNoErr(t, err)
	tu.Assert(t, out[task.Id].HasProgression)

	// test with random mono
	task, err = ta.Task{IdRandomMonoquestion: randomMono.Id.AsOptional()}.Insert(db.DB)
	tu.AssertNoErr(t, err)

	rd, err := newRandomMonoquestionData(db, randomMono.Id, student.Id)
	tu.AssertNoErr(t, err)
	rd, err = rd.selectQuestions(db.DB, student.Id)
	tu.AssertNoErr(t, err)

	out, err = LoadTasksProgression(db, student.Id, []ta.IdTask{task.Id})
	tu.AssertNoErr(t, err)
	tu.Assert(t, !out[task.Id].HasProgression)

	err = updateProgression(db.DB, student.Id, task.Id, []ta.QuestionHistory{
		{false, true},
		{},
		{},
		{false, true},
	})
	tu.AssertNoErr(t, err)

	out, err = LoadTasksProgression(db, student.Id, []ta.IdTask{task.Id})
	tu.AssertNoErr(t, err)
	tu.Assert(t, out[task.Id].HasProgression)
}

func Test_selectVariants(t *testing.T) {
	among := []ed.Question{
		{},
		{},
		{},
		{},
	}

	tests := []struct {
		nbToSelect int
		want       []ed.Question
	}{
		{0, []ed.Question{}},
		{2, []ed.Question{{}, {}}},
		{4, []ed.Question{{}, {}, {}, {}}},
		{5, []ed.Question{{}, {}, {}, {}, {}}},
	}
	for _, tt := range tests {
		if got := selectVariants(tt.nbToSelect, among); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("selectVariants() = %v, want %v", got, tt.want)
		}
	}

	tu.Assert(t, ed.Diff1 < ed.Diff2)
	got := selectVariants(4, []ed.Question{
		{Difficulty: ed.Diff1},
		{Difficulty: ed.Diff1},
		{Difficulty: ed.Diff2},
		{Difficulty: ed.Diff2},
		{Difficulty: ed.Diff3},
		{Difficulty: ed.Diff3},
		{Difficulty: ed.Diff3},
	})
	tu.Assert(t, sort.SliceIsSorted(got, func(i, j int) bool { return got[i].Difficulty < got[j].Difficulty }))
}

func TestRandomMonoquestion(t *testing.T) {
	db := tu.NewTestDB(t, "../sql/teacher/gen_create.sql", "../sql/editor/gen_create.sql", "../sql/tasks/gen_create.sql")
	defer db.Remove()

	tc, err := teacher.Teacher{IsAdmin: true, FavoriteMatiere: teacher.Mathematiques}.Insert(db)
	tu.AssertNoErr(t, err)

	cl, err := teacher.Classroom{}.Insert(db)
	tu.AssertNoErr(t, err)

	student, err := teacher.Student{IdClassroom: cl.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	_, _, _, randomMono := createEx(t, db.DB, tc.Id)

	task, err := ta.Task{IdRandomMonoquestion: randomMono.Id.AsOptional()}.Insert(db.DB)
	tu.AssertNoErr(t, err)

	_, err = InstantiateWork(db.DB, NewWorkID(task), student.Id)
	tu.AssertNoErr(t, err)

	out, err := newRandomMonoquestionData(db.DB, randomMono.Id, student.Id)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(out.selectedQuestions) == 4)
	selected := out.selectedQuestions

	// make sure InstantiateWork preserve already chosen questions
	_, err = InstantiateWork(db.DB, NewWorkID(task), student.Id)
	tu.AssertNoErr(t, err)

	out, err = newRandomMonoquestionData(db.DB, randomMono.Id, student.Id)
	tu.AssertNoErr(t, err)
	tu.Assert(t, reflect.DeepEqual(selected, out.selectedQuestions))
}
