package questions

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/benoitkugler/maths-online/server/src/maths/expression"
	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
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
			EnonceInstance{NumberFieldInstance{ID: 0, Answer: 3.4}, ExpressionFieldInstance{ID: 1, Answer: expression.MustParse("x+2"), ComparisonLevel: SimpleSubstitutions}},
			map[int]client.Answer{0: client.NumberAnswer{Value: 3.4}, 1: client.ExpressionAnswer{Expression: "x+3"}},
			map[int]bool{0: true, 1: false},
		},
		{
			EnonceInstance{NumberFieldInstance{ID: 0, Answer: 3.4}, ExpressionFieldInstance{ID: 1, Answer: expression.MustParse("x+2"), ComparisonLevel: SimpleSubstitutions}},
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

func TestBug144(t *testing.T) {
	jsonInput := `[{"variable": {"Name": 107, "Indice": ""}, "expression": "randint(0;1)"}, {"variable": {"Name": 107, "Indice": "1"}, "expression": "k==0"}, {"variable": {"Name": 107, "Indice": "2"}, "expression": "k==1"}, {"variable": {"Name": 97, "Indice": ""}, "expression": "k_1*(randChoice(-10;-9;-8;-7;-6;-5;-4;-3;-2;-1;-1;-1;-1;-1;-1;-1;-1))+k_2*(randChoice(-10;-8;-5;-4;-2;-1;-1;-1;-1))"}, {"variable": {"Name": 115, "Indice": ""}, "expression": "k_2*(1/a)+k_1*0"}, {"variable": {"Name": 98, "Indice": ""}, "expression": "k_1*(1/a)+k_2*0"}]`
	var prs RandomParameters
	err := json.Unmarshal([]byte(jsonInput), &prs)
	tu.Assert(t, err == nil)

	out := prs.toMap()

	rp, err := out.Instantiate()
	tu.Assert(t, err == nil)

	tu.Assert(t, rp[expression.NewVar('b')].String() == "0")
}
