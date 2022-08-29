package trivial

import (
	"reflect"
	"testing"

	ed "github.com/benoitkugler/maths-online/sql/editor"
	"github.com/benoitkugler/maths-online/sql/teacher"
	tr "github.com/benoitkugler/maths-online/sql/trivial"
	tu "github.com/benoitkugler/maths-online/utils/testutils"
)

func qu(group ed.IdQuestiongroup) ed.Question {
	return ed.Question{IdGroup: group.AsOptional()}
}

func quD(id ed.IdQuestion, group ed.IdQuestiongroup, diff ed.DifficultyTag) ed.Question {
	return ed.Question{Id: id, IdGroup: group.AsOptional(), Difficulty: diff}
}

func Test_weightQuestions(t *testing.T) {
	tests := []struct {
		args ed.Questions
		want []float64
	}{
		{nil, []float64{}},
		{ed.Questions{1: {}}, []float64{1}},
		{
			ed.Questions{
				1: quD(1, 1, ""),
				2: quD(2, 2, ""),
				3: quD(3, 3, ""),
			},
			[]float64{1. / 3, 1. / 3, 1. / 3},
		},
		{
			ed.Questions{
				0: quD(0, 1, ""),
				1: quD(1, 1, ""),
				2: quD(2, 2, ""),
				3: quD(3, 3, ""),
			},
			[]float64{1. / 6, 1. / 6, 1. / 3, 1. / 3},
		},
		{
			ed.Questions{
				0: quD(0, 1, ""),
				1: quD(1, 1, "_"),
				2: quD(2, 1, "_"),
				3: quD(3, 2, ""),
				4: quD(4, 3, ""),
			},
			[]float64{1. / 6, 1. / 12, 1. / 12, 1. / 3, 1. / 3},
		},
		{
			ed.Questions{
				0: quD(0, 1, ""),
				1: quD(1, 1, "_"),
				2: quD(2, 1, "_"),
				3: quD(3, 2, ""),
				4: quD(4, 2, ""),
				5: quD(5, 3, ""),
			},
			[]float64{1. / 6, 1. / 12, 1. / 12, 1. / 6, 1. / 6, 1. / 3},
		},
	}
	for _, tt := range tests {
		if got := weightQuestions(tt.args).Weights; !reflect.DeepEqual(got, tt.want) {
			t.Errorf("weightQuestions() = %v, want %v", got, tt.want)
		}
	}
}

func TestSelectQuestions(t *testing.T) {
	// create a DB shared by all tests
	db := tu.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../sql/editor/gen_create.sql")
	defer db.Remove()

	tc, err := teacher.Teacher{}.Insert(db)
	tu.Assert(t, err == nil)

	g1, err := ed.Questiongroup{IdTeacher: tc.Id}.Insert(db)
	tu.Assert(t, err == nil)
	g2, err := ed.Questiongroup{IdTeacher: tc.Id}.Insert(db)
	tu.Assert(t, err == nil)
	g3, err := ed.Questiongroup{IdTeacher: tc.Id}.Insert(db)
	tu.Assert(t, err == nil)
	g4, err := ed.Questiongroup{IdTeacher: tc.Id}.Insert(db)
	tu.Assert(t, err == nil)

	quD(0, g1.Id, ed.Diff1).Insert(db)
	quD(0, g1.Id, ed.Diff2).Insert(db)
	quD(0, g1.Id, ed.Diff3).Insert(db)
	quD(0, g1.Id, ed.Diff3).Insert(db)

	qu(g2.Id).Insert(db)
	qu(g2.Id).Insert(db)

	qu(g3.Id).Insert(db)
	qu(g3.Id).Insert(db)
	qu(g3.Id).Insert(db)

	qu(g4.Id).Insert(db)

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	err = ed.InsertManyQuestiongroupTags(tx,
		// categorie tags
		ed.QuestiongroupTag{IdQuestiongroup: 1, Tag: "KEEP"},
		ed.QuestiongroupTag{IdQuestiongroup: 2, Tag: "KEEP"},
		ed.QuestiongroupTag{IdQuestiongroup: 3, Tag: "KEEP"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}

	criterion := tr.QuestionCriterion{{"KEEP"}}
	cats := tr.CategoriesQuestions{
		Tags:         [...]tr.QuestionCriterion{criterion, criterion, criterion, criterion, criterion},
		Difficulties: nil,
	}

	pool, err := selectQuestions(db, cats, tc.Id)
	tu.Assert(t, err == nil)

	if !reflect.DeepEqual(pool[0].Weights, []float64{
		// group 1 : 4 questions -> 3 sub group
		1. / 9,
		1. / 9,
		1. / 18, 1. / 18,
		// group 2 : 2 questions
		1. / 6, 1. / 6,
		// group 3 : 3 questions
		1. / 9, 1. / 9, 1. / 9,
	}) {
		t.Fatal(pool[0].Weights)
	}
}

func TestQuestionCriterion_filter(t *testing.T) {
	tests := []struct {
		qc      tr.QuestionCriterion
		args    ed.QuestiongroupTags
		wantOut []ed.IdQuestiongroup
	}{
		{
			tr.QuestionCriterion{},
			ed.QuestiongroupTags{
				{IdQuestiongroup: 1, Tag: ""},
				{IdQuestiongroup: 2, Tag: ""},
			},
			nil,
		},
		{
			tr.QuestionCriterion{
				{"TAG1"},
			},
			ed.QuestiongroupTags{
				{IdQuestiongroup: 1, Tag: "TAG1"},
				{IdQuestiongroup: 1, Tag: "xx"},
				{IdQuestiongroup: 2, Tag: "xx"},
			},
			[]ed.IdQuestiongroup{1},
		},
		{
			tr.QuestionCriterion{
				{"TAG1", "TAG2"},
			},
			ed.QuestiongroupTags{
				{IdQuestiongroup: 1, Tag: "TAG1"},
				{IdQuestiongroup: 1, Tag: "xx"},
				{IdQuestiongroup: 2, Tag: "xx"},
			},
			nil,
		},
		{
			tr.QuestionCriterion{
				{"TAG1", "TAG2"},
			},
			ed.QuestiongroupTags{
				{IdQuestiongroup: 1, Tag: "TAG1"},
				{IdQuestiongroup: 1, Tag: "TAG2"},
				{IdQuestiongroup: 1, Tag: "TAG1"},
				{IdQuestiongroup: 2, Tag: "xx"},
			},
			[]ed.IdQuestiongroup{1},
		},
		{
			tr.QuestionCriterion{
				{"TAG1", "TAG2"},
				{"TAG3"},
			},
			ed.QuestiongroupTags{
				{IdQuestiongroup: 1, Tag: "TAG1"},
				{IdQuestiongroup: 1, Tag: "TAG2"},
				{IdQuestiongroup: 1, Tag: "TAG1"},
				{IdQuestiongroup: 2, Tag: "TAG3"},
			},
			[]ed.IdQuestiongroup{1, 2},
		},
	}
	for _, tt := range tests {
		if gotOut := filterTags(tt.qc, tt.args); !reflect.DeepEqual(gotOut, tt.wantOut) {
			t.Errorf("QuestionCriterion.filter() = %v, want %v", gotOut, tt.wantOut)
		}
	}
}

func BenchmarkQuestionSearch(b *testing.B) {
	b.StopTimer()

	db, err := tu.DB.ConnectPostgres()
	if err != nil {
		b.Skipf("DB %v not available : %s", tu.DB, err)
		return
	}

	sel, err := newQuestionSelector(db)
	if err != nil {
		b.Fatal(err)
	}

	tvs, err := tr.SelectAllTrivials(db)
	if err != nil {
		b.Fatal(err)
	}

	questions := tvs[60].Questions

	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sel.search(questions, 0)
	}
}
