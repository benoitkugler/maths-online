package editor

import (
	"fmt"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/maths/questions"
	"github.com/benoitkugler/maths-online/prof/teacher"
	"github.com/benoitkugler/maths-online/utils/testutils"
)

func TestValidation(t *testing.T) {
	db, err := testutils.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", testutils.DB, err)
		return
	}

	qu, err := SelectAllQuestions(db)
	if err != nil {
		t.Fatal(err)
	}
	qu.RestrictNeedExercice()

	ti := time.Now()
	err = validateAllQuestions(qu, nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Validated in :", time.Since(ti), "average :", time.Since(ti)/time.Duration(len(qu)))
}

func BenchmarkValidation(b *testing.B) {
	db, err := testutils.DB.ConnectPostgres()
	if err != nil {
		b.Skipf("DB %v not available : %s", testutils.DB, err)
		return
	}

	qu, err := SelectAllQuestions(db)
	if err != nil {
		b.Fatal(err)
	}
	qu.RestrictNeedExercice()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		validateAllQuestions(qu, nil)
	}
}

func TestExerciceCRUD(t *testing.T) {
	db := testutils.CreateDBDev(t, "../teacher/gen_create.sql", "gen_create.sql")
	defer testutils.RemoveDBDev()
	defer db.Close()

	_, err := teacher.Teacher{IsAdmin: true}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	ct := NewController(db, teacher.Teacher{Id: 1})

	ex, err := ct.createExercice(1)
	if err != nil {
		t.Fatal(err)
	}

	l, err := ct.createQuestionEx(ExerciceCreateQuestionIn{IdExercice: ex.Exercice.Id}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(l.Questions) != 1 {
		t.Fatal(l)
	}

	qu, err := Question{IdTeacher: 1}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	l, err = ct.updateQuestionsEx(ExerciceUpdateQuestionsIn{
		IdExercice: ex.Exercice.Id,
		Questions: ExerciceQuestions{
			l.Questions[0],
			ExerciceQuestion{IdQuestion: qu.Id},
			ExerciceQuestion{IdQuestion: qu.Id},
		},
	}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(l.Questions) != 3 {
		t.Fatal(l)
	}

	updated := l.Questions[1]
	updated.Bareme = 5
	_, err = ct.updateQuestionsEx(ExerciceUpdateQuestionsIn{
		IdExercice: ex.Exercice.Id,
		Questions: ExerciceQuestions{
			l.Questions[0],
			updated,
		},
	}, 1)
	if err != nil {
		t.Fatal(err)
	}

	exe, err := ct.updateExercice(Exercice{Id: ex.Exercice.Id, Description: "test", Title: "test2", Flow: Sequencial}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if exe.Flow != Sequencial {
		t.Fatal(exe)
	}

	err = ct.deleteExercice(ex.Exercice.Id, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDB(t *testing.T) {
	db, err := testutils.DB.ConnectPostgres()
	if err != nil {
		t.Skip("DB not available")
	}

	ct := NewController(db, teacher.Teacher{Id: 1})
	_, err = ct.getExercices(1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGroupTagsEmpty(t *testing.T) {
	db := testutils.CreateDBDev(t, "../teacher/gen_create.sql", "gen_create.sql")
	defer testutils.RemoveDBDev()
	defer db.Close()

	_, err := teacher.Teacher{IsAdmin: true}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	ct := NewController(db, teacher.Teacher{Id: 1})

	// create an implicit groups with no tags
	_, err = Question{IdTeacher: 1, Page: questions.QuestionPage{Title: "test"}}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Question{IdTeacher: 1, Page: questions.QuestionPage{Title: "test"}}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	out, err := ct.updateGroupTags(UpdateGroupTagsIn{GroupTitle: "test", CommonTags: []string{"newtag1", "newtag2"}}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Tags) != 2 {
		t.Fatal(out)
	}
}

func TestProgression(t *testing.T) {
	db := testutils.CreateDBDev(t, "../teacher/gen_create.sql", "gen_create.sql")
	defer testutils.RemoveDBDev()
	defer db.Close()

	_, err := teacher.Teacher{IsAdmin: true}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	ct := NewController(db, teacher.Teacher{Id: 1})

	ex, err := ct.createExercice(1)
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

	tx, err := ct.db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	prog, err := Progression{IdExercice: ex.Exercice.Id}.Insert(tx)
	if err != nil {
		t.Fatal(err)
	}

	err = InsertManyProgressionQuestions(tx,
		ProgressionQuestion{IdProgression: prog.Id, IdExercice: prog.IdExercice, Index: 0, History: QuestionHistory{false, true}},
		ProgressionQuestion{IdProgression: prog.Id, IdExercice: prog.IdExercice, Index: 2, History: randQuestionHistory()},
	)
	if err != nil {
		t.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		t.Fatal(err)
	}

	out, err := ct.fetchProgression(prog.Id)
	if err != nil {
		t.Fatal(err)
	}
	if out.NextQuestion != 1 {
		t.Fatal(out)
	}
}
