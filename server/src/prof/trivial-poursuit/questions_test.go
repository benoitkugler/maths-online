package trivialpoursuit

import (
	"reflect"
	"testing"

	ex "github.com/benoitkugler/maths-online/maths/exercice"
	"github.com/benoitkugler/maths-online/utils/testutils"
)

func qu(title, diff string) questionDiff {
	return questionDiff{question: ex.Question{Title: title}, diff: ex.DifficultyTag(diff)}
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
				qu("1", ""),
				qu("2", ""),
				qu("3", ""),
			},
			[]float64{1. / 3, 1. / 3, 1. / 3},
		},
		{
			[]questionDiff{
				qu("1", ""),
				qu("1", ""),
				qu("2", ""),
				qu("3", ""),
			},
			[]float64{1. / 6, 1. / 6, 1. / 3, 1. / 3},
		},
		{
			[]questionDiff{
				qu("1", ""),
				qu("1", "_"),
				qu("1", "_"),
				qu("2", ""),
				qu("3", ""),
			},
			[]float64{1. / 6, 1. / 12, 1. / 12, 1. / 3, 1. / 3},
		},
		{
			[]questionDiff{
				qu("1", ""),
				qu("1", "_"),
				qu("1", "_"),
				qu("2", ""),
				qu("2", ""),
				qu("3", ""),
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
	db := testutils.CreateDBDev(t, "../../maths/exercice/create_gen.sql")
	defer testutils.RemoveDBDev()
	defer db.Close()

	ex.Question{Title: "title1"}.Insert(db)
	ex.Question{Title: "title1"}.Insert(db)
	ex.Question{Title: "title1"}.Insert(db)
	ex.Question{Title: "title1"}.Insert(db)
	ex.Question{Title: "title2"}.Insert(db)
	ex.Question{Title: "title2"}.Insert(db)
	ex.Question{Title: "title3"}.Insert(db)

	ex.Question{Title: "title3"}.Insert(db)
	ex.Question{Title: "title3"}.Insert(db)
	ex.Question{Title: "title4"}.Insert(db)

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	err = ex.InsertManyQuestionTags(tx,
		ex.QuestionTag{IdQuestion: 1, Tag: string(ex.Diff1)},
		ex.QuestionTag{IdQuestion: 2, Tag: string(ex.Diff2)},
		ex.QuestionTag{IdQuestion: 3, Tag: string(ex.Diff3)},
		ex.QuestionTag{IdQuestion: 4, Tag: string(ex.Diff3)},
		// categorie tags
		ex.QuestionTag{IdQuestion: 1, Tag: "keep"},
		ex.QuestionTag{IdQuestion: 2, Tag: "keep"},
		ex.QuestionTag{IdQuestion: 3, Tag: "keep"},
		ex.QuestionTag{IdQuestion: 4, Tag: "keep"},
		ex.QuestionTag{IdQuestion: 5, Tag: "keep"},
		ex.QuestionTag{IdQuestion: 6, Tag: "keep"},
		ex.QuestionTag{IdQuestion: 7, Tag: "keep"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}

	criterion := QuestionCriterion{{"keep"}}
	cats := CategoriesQuestions{criterion, criterion, criterion, criterion, criterion}

	pool, err := cats.selectQuestions(db)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(pool[0].Weights, []float64{
		1. / 9, 1. / 9, 1. / 18, 1. / 18, 1. / 6, 1. / 6, 1. / 3,
	}) {
		t.Fatal(pool[0].Weights)
	}
}
