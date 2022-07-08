package questions

import (
	"reflect"
	"testing"

	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/questions/client"
)

func TestQuestionInstance_CompareAnswer(t *testing.T) {
	tests := []struct {
		fields EnonceInstance
		args   map[int]client.Answer
		want   map[int]bool
	}{
		{
			EnonceInstance{NumberFieldInstance{ID: 0, Answer: 3.4}, NumberFieldInstance{ID: 1, Answer: -0.2}},
			map[int]client.Answer{},
			map[int]bool{0: false, 1: false},
		},
		{
			EnonceInstance{NumberFieldInstance{ID: 0, Answer: 3.4}, NumberFieldInstance{ID: 1, Answer: -0.2}},
			map[int]client.Answer{0: client.NumberAnswer{Value: 3.4}},
			map[int]bool{0: true, 1: false},
		},
		{
			EnonceInstance{NumberFieldInstance{ID: 0, Answer: 3.4}, NumberFieldInstance{ID: 1, Answer: -0.2}},
			map[int]client.Answer{0: client.NumberAnswer{Value: 3.4}, 1: client.RadioAnswer{}},
			map[int]bool{0: true, 1: false},
		},
		{
			EnonceInstance{NumberFieldInstance{ID: 0, Answer: 3.4}, NumberFieldInstance{ID: 1, Answer: -0.2}},
			map[int]client.Answer{0: client.NumberAnswer{Value: 3.4}, 1: client.NumberAnswer{Value: 0}},
			map[int]bool{0: true, 1: false},
		},
		{
			EnonceInstance{NumberFieldInstance{ID: 0, Answer: 3.4}, NumberFieldInstance{ID: 1, Answer: -0.2}},
			map[int]client.Answer{0: client.NumberAnswer{Value: 3.4}, 1: client.NumberAnswer{Value: -0.2}},
			map[int]bool{0: true, 1: true},
		},
		{
			EnonceInstance{NumberFieldInstance{ID: 0, Answer: 3.4}, ExpressionFieldInstance{ID: 1, Answer: expression.MustParse("x+2"), ComparisonLevel: expression.SimpleSubstitutions}},
			map[int]client.Answer{0: client.NumberAnswer{Value: 3.4}, 1: client.ExpressionAnswer{Expression: "x+3"}},
			map[int]bool{0: true, 1: false},
		},
		{
			EnonceInstance{NumberFieldInstance{ID: 0, Answer: 3.4}, ExpressionFieldInstance{ID: 1, Answer: expression.MustParse("x+2"), ComparisonLevel: expression.SimpleSubstitutions}},
			map[int]client.Answer{0: client.NumberAnswer{Value: 3.4}, 1: client.ExpressionAnswer{Expression: "x+ 2"}},
			map[int]bool{0: true, 1: true},
		},
	}
	for _, tt := range tests {
		qu := QuestionInstance{
			Enonce: tt.fields,
		}
		if got := qu.EvaluateAnswer(client.QuestionAnswersIn{Data: tt.args}); !reflect.DeepEqual(got.Results, tt.want) {
			t.Errorf("QuestionInstance.CompareAnswer(%v) = %v, want %v", tt.args, got.Results, tt.want)
		}
	}
}
