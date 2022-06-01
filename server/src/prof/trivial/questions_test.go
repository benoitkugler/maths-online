package trivial

import (
	"reflect"
	"testing"

	ex "github.com/benoitkugler/maths-online/maths/exercice"
	"github.com/benoitkugler/maths-online/prof/editor"
	ed "github.com/benoitkugler/maths-online/prof/editor"
	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
	"github.com/benoitkugler/maths-online/utils/testutils"
)

func qu(title string) ed.Question {
	return ed.Question{Page: ex.QuestionPage{Title: title}}
}

func quD(title, diff string) questionDiff {
	return questionDiff{question: qu(title), diff: ed.DifficultyTag(diff)}
}

var dummyQuestions = game.QuestionPool{
	game.WeigthedQuestions{Questions: []ed.Question{qu("Qu1")}, Weights: []float64{1}},
	game.WeigthedQuestions{Questions: []ed.Question{qu("Qu1")}, Weights: []float64{1}},
	game.WeigthedQuestions{Questions: []ed.Question{qu("Qu1")}, Weights: []float64{1}},
	game.WeigthedQuestions{Questions: []ed.Question{qu("Qu1")}, Weights: []float64{1}},
	game.WeigthedQuestions{Questions: []ed.Question{qu("Qu1")}, Weights: []float64{1}},
}

func Test_weightQuestions(t *testing.T) {
	tests := []struct {
		args []questionDiff
		want []float64
	}{
		{nil, []float64{}},
		{[]questionDiff{{}}, []float64{1}},
		{
			[]questionDiff{
				quD("1", ""),
				quD("2", ""),
				quD("3", ""),
			},
			[]float64{1. / 3, 1. / 3, 1. / 3},
		},
		{
			[]questionDiff{
				quD("1", ""),
				quD("1", ""),
				quD("2", ""),
				quD("3", ""),
			},
			[]float64{1. / 6, 1. / 6, 1. / 3, 1. / 3},
		},
		{
			[]questionDiff{
				quD("1", ""),
				quD("1", "_"),
				quD("1", "_"),
				quD("2", ""),
				quD("3", ""),
			},
			[]float64{1. / 6, 1. / 12, 1. / 12, 1. / 3, 1. / 3},
		},
		{
			[]questionDiff{
				quD("1", ""),
				quD("1", "_"),
				quD("1", "_"),
				quD("2", ""),
				quD("2", ""),
				quD("3", ""),
			},
			[]float64{1. / 6, 1. / 12, 1. / 12, 1. / 6, 1. / 6, 1. / 3},
		},
	}
	for _, tt := range tests {
		if got := weightQuestions(tt.args); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("weightQuestions() = %v, want %v", got, tt.want)
		}
	}
}

func TestSelectQuestions(t *testing.T) {
	// create a DB shared by all tests
	db := testutils.CreateDBDev(t, "../editor/gen_create.sql")
	defer testutils.RemoveDBDev()
	defer db.Close()

	qu("title1").Insert(db)
	qu("title1").Insert(db)
	qu("title1").Insert(db)
	qu("title1").Insert(db)
	qu("title2").Insert(db)
	qu("title2").Insert(db)
	qu("title3").Insert(db)

	qu("title3").Insert(db)
	qu("title3").Insert(db)
	qu("title4").Insert(db)

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	err = ed.InsertManyQuestionTags(tx,
		ed.QuestionTag{IdQuestion: 1, Tag: string(ed.Diff1)},
		ed.QuestionTag{IdQuestion: 2, Tag: string(ed.Diff2)},
		ed.QuestionTag{IdQuestion: 3, Tag: string(ed.Diff3)},
		ed.QuestionTag{IdQuestion: 4, Tag: string(ed.Diff3)},
		// categorie tags
		ed.QuestionTag{IdQuestion: 1, Tag: "keep"},
		ed.QuestionTag{IdQuestion: 2, Tag: "keep"},
		ed.QuestionTag{IdQuestion: 3, Tag: "keep"},
		ed.QuestionTag{IdQuestion: 4, Tag: "keep"},
		ed.QuestionTag{IdQuestion: 5, Tag: "keep"},
		ed.QuestionTag{IdQuestion: 6, Tag: "keep"},
		ed.QuestionTag{IdQuestion: 7, Tag: "keep"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}

	criterion := QuestionCriterion{{"keep"}}
	cats := CategoriesQuestions{criterion, criterion, criterion, criterion, criterion}

	pool, err := cats.selectQuestions(db, 0)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(pool[0].Weights, []float64{
		1. / 9, 1. / 9, 1. / 18, 1. / 18, 1. / 6, 1. / 6, 1. / 3,
	}) {
		t.Fatal(pool[0].Weights)
	}
}

func TestQuestionCriterion_filter(t *testing.T) {
	tests := []struct {
		qc      QuestionCriterion
		args    editor.QuestionTags
		wantOut IDs
	}{
		{
			QuestionCriterion{},
			editor.QuestionTags{
				{IdQuestion: 1, Tag: ""},
				{IdQuestion: 2, Tag: ""},
			},
			nil,
		},
		{
			QuestionCriterion{
				{"TAG1"},
			},
			editor.QuestionTags{
				{IdQuestion: 1, Tag: "TAG1"},
				{IdQuestion: 1, Tag: "xx"},
				{IdQuestion: 2, Tag: "xx"},
			},
			IDs{1},
		},
		{
			QuestionCriterion{
				{"TAG1", "TAG2"},
			},
			editor.QuestionTags{
				{IdQuestion: 1, Tag: "TAG1"},
				{IdQuestion: 1, Tag: "xx"},
				{IdQuestion: 2, Tag: "xx"},
			},
			nil,
		},
		{
			QuestionCriterion{
				{"TAG1", "TAG2"},
			},
			editor.QuestionTags{
				{IdQuestion: 1, Tag: "TAG1"},
				{IdQuestion: 1, Tag: "TAG2"},
				{IdQuestion: 1, Tag: "TAG1"},
				{IdQuestion: 2, Tag: "xx"},
			},
			IDs{1},
		},
		{
			QuestionCriterion{
				{"TAG1", "TAG2"},
				{"TAG3"},
			},
			editor.QuestionTags{
				{IdQuestion: 1, Tag: "TAG1"},
				{IdQuestion: 1, Tag: "TAG2"},
				{IdQuestion: 1, Tag: "TAG1"},
				{IdQuestion: 2, Tag: "TAG3"},
			},
			IDs{1, 2},
		},
	}
	for _, tt := range tests {
		if gotOut := tt.qc.filter(tt.args.ByIdQuestion()); !reflect.DeepEqual(gotOut, tt.wantOut) {
			t.Errorf("QuestionCriterion.filter() = %v, want %v", gotOut, tt.wantOut)
		}
	}
}
