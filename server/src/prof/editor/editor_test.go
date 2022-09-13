package editor

import (
	"fmt"
	"testing"
	"time"

	ed "github.com/benoitkugler/maths-online/sql/editor"
	"github.com/benoitkugler/maths-online/sql/teacher"
	tu "github.com/benoitkugler/maths-online/utils/testutils"
)

func TestValidation(t *testing.T) {
	db, err := tu.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", tu.DB, err)
		return
	}

	qu, err := ed.SelectAllQuestions(db)
	tu.Assert(t, err == nil)
	qu.RestrictNeedExercice()

	ti := time.Now()
	err = validateAllQuestions(qu, nil)
	tu.Assert(t, err == nil)
	fmt.Println("Validated in :", time.Since(ti), "average :", time.Since(ti)/time.Duration(len(qu)))
}

func BenchmarkValidation(b *testing.B) {
	db, err := tu.DB.ConnectPostgres()
	if err != nil {
		b.Skipf("DB %v not available : %s", tu.DB, err)
		return
	}

	qu, err := ed.SelectAllQuestions(db)
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
	db := tu.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../sql/editor/gen_create.sql")
	defer db.Remove()

	_, err := teacher.Teacher{IsAdmin: true}.Insert(db)
	tu.Assert(t, err == nil)

	ct := NewController(db.DB, teacher.Teacher{Id: 1})

	group, err := ct.createExercice(1)
	tu.Assert(t, err == nil)
	ex := group.Variants[0]

	l, err := ct.createQuestionEx(ExerciceCreateQuestionIn{IdExercice: ex.Id}, 1)
	tu.Assert(t, err == nil)
	if len(l.Questions) != 1 {
		t.Fatal(l)
	}

	qu, err := ed.Question{NeedExercice: ex.Id.AsOptional()}.Insert(db)
	tu.Assert(t, err == nil)

	l, err = ct.updateQuestionsEx(ExerciceUpdateQuestionsIn{
		IdExercice: ex.Id,
		Questions: ed.ExerciceQuestions{
			ed.ExerciceQuestion{IdQuestion: l.Questions[0].Question.Id},
			ed.ExerciceQuestion{IdQuestion: qu.Id},
			ed.ExerciceQuestion{IdQuestion: qu.Id},
		},
	}, 1)
	tu.Assert(t, err == nil)
	if len(l.Questions) != 3 {
		t.Fatal(l)
	}

	updated := l.Questions[1]
	updated.Bareme = 5
	_, err = ct.updateQuestionsEx(ExerciceUpdateQuestionsIn{
		IdExercice: ex.Id,
		Questions: ed.ExerciceQuestions{
			ed.ExerciceQuestion{IdQuestion: l.Questions[0].Question.Id},
			ed.ExerciceQuestion{IdQuestion: updated.Question.Id},
		},
	}, 1)
	tu.Assert(t, err == nil)

	exe, err := ct.updateExercice(ed.Exercice{Id: ex.Id, Description: "test", Subtitle: "test2"}, 1)
	tu.Assert(t, err == nil)
	if exe.Description != "test" {
		t.Fatal(exe)
	}

	err = ct.deleteExercice(ex.Id, 1)
	tu.Assert(t, err == nil)
}

func TestDB(t *testing.T) {
	db, err := tu.DB.ConnectPostgres()
	if err != nil {
		t.Skip("DB not available")
	}

	ct := NewController(db, teacher.Teacher{Id: 1})
	_, err = ct.searchExercices(Query{}, 1)
	tu.Assert(t, err == nil)
}

func TestGroupTagsEmpty(t *testing.T) {
	db := tu.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../sql/editor/gen_create.sql")
	defer db.Remove()

	_, err := teacher.Teacher{IsAdmin: true}.Insert(db)
	tu.Assert(t, err == nil)

	ct := NewController(db.DB, teacher.Teacher{Id: 1})

	// create a group with no tags
	group, err := ed.Questiongroup{IdTeacher: 1}.Insert(db)
	tu.Assert(t, err == nil)

	_, err = ed.Question{IdGroup: group.Id.AsOptional()}.Insert(db)
	tu.Assert(t, err == nil)
	_, err = ed.Question{IdGroup: group.Id.AsOptional()}.Insert(db)
	tu.Assert(t, err == nil)

	err = ct.updateTags(UpdateQuestiongroupTagsIn{Id: group.Id, Tags: []string{"newtag1", "newtag2"}}, 1)
	tu.Assert(t, err == nil)
}
