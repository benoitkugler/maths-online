package editor

import (
	"fmt"
	"testing"
	"time"

	ed "github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func TestValidation(t *testing.T) {
	db, err := tu.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", tu.DB, err)
		return
	}

	qu, err := ed.SelectAllQuestions(db)
	tu.AssertNoErr(t, err)

	exs, err := ed.SelectAllExercices(db)
	tu.AssertNoErr(t, err)

	ti := time.Now()
	err = validateAllQuestions(qu, exs)
	tu.AssertNoErr(t, err)
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
	db := tu.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../sql/editor/gen_create.sql", "../../sql/tasks/gen_create.sql",
		"../../sql/homework/gen_create.sql", "../../sql/reviews/gen_create.sql")
	defer db.Remove()

	_, err := teacher.Teacher{IsAdmin: true}.Insert(db)
	tu.AssertNoErr(t, err)

	ct := NewController(db.DB, teacher.Teacher{Id: 1})

	group, err := ct.createExercice(1)
	tu.AssertNoErr(t, err)
	ex := group.Variants[0]

	l, err := ct.createQuestionEx(ExerciceCreateQuestionIn{IdExercice: ex.Id}, 1)
	tu.AssertNoErr(t, err)
	if len(l.Questions) != 1 {
		t.Fatal(l)
	}

	qu, err := ed.Question{NeedExercice: ex.Id.AsOptional()}.Insert(db)
	tu.AssertNoErr(t, err)

	l, err = ct.updateQuestionsEx(ExerciceUpdateQuestionsIn{
		IdExercice: ex.Id,
		Questions: ed.ExerciceQuestions{
			ed.ExerciceQuestion{IdQuestion: l.Questions[0].Question.Id},
			ed.ExerciceQuestion{IdQuestion: qu.Id},
		},
	}, 1)
	tu.AssertNoErr(t, err)
	if len(l.Questions) != 2 {
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
	tu.AssertNoErr(t, err)

	exe, err := ct.updateExercice(ExerciceHeader{Id: ex.Id, Difficulty: ed.Diff3, Subtitle: "test2"}, 1)
	tu.AssertNoErr(t, err)
	if exe.Difficulty != ed.Diff3 {
		t.Fatal(exe)
	}

	err = ct.duplicateExercicegroup(group.Group.Id, 1)
	tu.AssertNoErr(t, err)

	_, err = ct.deleteExercice(ex.Id, 1)
	tu.AssertNoErr(t, err)
}

func TestDB(t *testing.T) {
	db, err := tu.DB.ConnectPostgres()
	if err != nil {
		t.Skip("DB not available")
	}

	ct := NewController(db, teacher.Teacher{Id: 1})
	_, err = ct.searchExercices(Query{}, 1)
	tu.AssertNoErr(t, err)
}

func TestGroupTagsEmpty(t *testing.T) {
	db := tu.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../sql/editor/gen_create.sql")
	defer db.Remove()

	_, err := teacher.Teacher{IsAdmin: true}.Insert(db)
	tu.AssertNoErr(t, err)

	ct := NewController(db.DB, teacher.Teacher{Id: 1})

	// create a group with no tags
	group, err := ed.Questiongroup{IdTeacher: 1}.Insert(db)
	tu.AssertNoErr(t, err)

	_, err = ed.Question{IdGroup: group.Id.AsOptional()}.Insert(db)
	tu.AssertNoErr(t, err)
	_, err = ed.Question{IdGroup: group.Id.AsOptional()}.Insert(db)
	tu.AssertNoErr(t, err)

	err = ct.updateQuestionTags(UpdateQuestiongroupTagsIn{Id: group.Id, Tags: ed.Tags{
		{Tag: "newtag1", Section: ed.Level},
		{Tag: "newtag2", Section: ed.Chapter},
		{Tag: "newtag3", Section: ed.TrivMath},
		{Tag: "newtag4", Section: ed.TrivMath},
	}}, 1)
	tu.AssertNoErr(t, err)

	// same for exercice
	group2, err := ed.Exercicegroup{IdTeacher: 1}.Insert(db)
	tu.AssertNoErr(t, err)

	err = ct.updateExerciceTags(UpdateExercicegroupTagsIn{Id: group2.Id, Tags: ed.Tags{
		{Tag: "newtag1", Section: ed.Level},
		{Tag: "newtag2", Section: ed.Chapter},
		{Tag: "newtag3", Section: ed.TrivMath},
		{Tag: "newtag4", Section: ed.TrivMath},
	}}, 1)
	tu.AssertNoErr(t, err)
}

func TestLoadTags(t *testing.T) {
	db, err := tu.DB.ConnectPostgres()
	if err != nil {
		t.Skip("DB not available")
	}

	ct := NewController(db, teacher.Teacher{Id: 1})
	tags, err := LoadTags(ct.db, ct.admin.Id)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(tags.Levels) >= 1)
	tu.Assert(t, len(tags.ChaptersByLevel) >= 1)
	// there might be some question with Level but no chapter
	tu.Assert(t, len(tags.ChaptersByLevel) <= len(tags.Levels))
}

func TestLoadIndex(t *testing.T) {
	db, err := tu.DB.ConnectPostgres()
	if err != nil {
		t.Skip("DB not available")
	}

	ct := NewController(db, teacher.Teacher{Id: 1})
	index, err := ct.loadQuestionsIndex(1)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(index) >= 3)

	index, err = ct.loadExercicesIndex(1)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(index) >= 1)
}
