package trivialpoursuit

import (
	"reflect"
	"testing"

	"github.com/benoitkugler/maths-online/maths/exercice"
)

func TestQuestionCriterion_filter(t *testing.T) {
	tests := []struct {
		qc      QuestionCriterion
		args    exercice.QuestionTags
		wantOut IDs
	}{
		{
			QuestionCriterion{},
			exercice.QuestionTags{
				{IdQuestion: 1, Tag: ""},
				{IdQuestion: 2, Tag: ""},
			},
			nil,
		},
		{
			QuestionCriterion{
				{"TAG1"},
			},
			exercice.QuestionTags{
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
			exercice.QuestionTags{
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
			exercice.QuestionTags{
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
			exercice.QuestionTags{
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
