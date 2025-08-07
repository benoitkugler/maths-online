package questions

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	ex "github.com/benoitkugler/maths-online/server/src/maths/expression"
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
			EnonceInstance{NumberFieldInstance{ID: 0, Answer: 3.4}, ExpressionFieldInstance{ID: 1, Answer: ex.MustParse("x+2"), ComparisonLevel: SimpleSubstitutions}},
			map[int]client.Answer{0: client.NumberAnswer{Value: 3.4}, 1: client.ExpressionAnswer{Expression: "x+3"}},
			map[int]bool{0: true, 1: false},
		},
		{
			EnonceInstance{NumberFieldInstance{ID: 0, Answer: 3.4}, ExpressionFieldInstance{ID: 1, Answer: ex.MustParse("x+2"), ComparisonLevel: SimpleSubstitutions}},
			map[int]client.Answer{0: client.NumberAnswer{Value: 3.4}, 1: client.ExpressionAnswer{Expression: "x+ 2"}},
			map[int]bool{0: true, 1: true},
		},
	}
	for _, tt := range tests {
		qu := tt.fields
		if got := qu.EvaluateAnswer(client.QuestionAnswersIn{Data: tt.args}); !reflect.DeepEqual(got.Results, tt.want) {
			t.Errorf("QuestionInstance.CompareAnswer(%v) = %v, want %v", tt.args, got.Results, tt.want)
		}
	}
}

func TestBug144(t *testing.T) {
	jsonInput := `[{"variable": {"Name": 107, "Indice": ""}, "expression": "randint(0;1)"}, {"variable": {"Name": 107, "Indice": "1"}, "expression": "k==0"}, {"variable": {"Name": 107, "Indice": "2"}, "expression": "k==1"}, {"variable": {"Name": 97, "Indice": ""}, "expression": "k_1*(randChoice(-10;-9;-8;-7;-6;-5;-4;-3;-2;-1;-1;-1;-1;-1;-1;-1;-1))+k_2*(randChoice(-10;-8;-5;-4;-2;-1;-1;-1;-1))"}, {"variable": {"Name": 115, "Indice": ""}, "expression": "k_2*(1/a)+k_1*0"}, {"variable": {"Name": 98, "Indice": ""}, "expression": "k_1*(1/a)+k_2*0"}]`
	var prs []Rp
	err := json.Unmarshal([]byte(jsonInput), &prs)
	tu.AssertNoErr(t, err)

	asParams := make(Parameters, len(prs))
	for i, p := range prs {
		asParams[i] = p
	}
	out := asParams.ToMap()

	for range [100]int{} {
		rp, err := out.Instantiate()
		tu.AssertNoErr(t, err)

		k1, _ := rp[ex.NewVarI('k', "1")].Evaluate(nil)
		if k1 == 0 {
			tu.Assert(t, rp[ex.NewVar('b')].String() == "0")
		} else {
			tu.Assert(t, k1 == 1)
		}
	}
}

func TestInstantiateMinMax(t *testing.T) {
	jsonInput := `{"enonce": [{"Data": {"Bold": false, "Parts": "Résoudre dans $\\R$ l'équation :", "Italic": false, "Smaller": false}, "Kind": "TextBlock"}, {"Data": {"Parts": "&k*(x+a)*(x+b)&=&0&"}, "Kind": "FormulaBlock"}, {"Data": {"Label": "", "Answer": ["$S$", "$=$", "$\\{$", "&s_1&", "$;$", "&s_2&", "$\\}$"], "AdditionalProposals": ["$;$", "$:$", "$,$", "$.$", "$($", "$)$", "$[$", "$[$", "$]$", "$]$", "&-s_1&", "&-s_2&", "&k&", "&-k&", "$0$"]}, "Kind": "OrderedListFieldBlock"}, {"Data": {"Bold": false, "Parts": "On rangera les solutions par ordre croissant.", "Italic": true, "Smaller": true}, "Kind": "TextBlock"}], "parameters": {"Variables": [{"variable": {"Name": 113, "Indice": "1"}, "expression": "randint(1;4)*randChoice(-1;1)"}, {"variable": {"Name": 113, "Indice": "2"}, "expression": "randint(1;4)*randChoice(-1;1)"}, {"variable": {"Name": 107, "Indice": ""}, "expression": "randint(2;10)*randChoice(-1;1)"}, {"variable": {"Name": 97, "Indice": ""}, "expression": "2*q_1"}, {"variable": {"Name": 98, "Indice": ""}, "expression": "2*q_2+1"}, {"variable": {"Name": 115, "Indice": "1"}, "expression": "min(-a;-b)"}, {"variable": {"Name": 115, "Indice": "2"}, "expression": "max(-a;-b)"}], "Intrinsics": null}}`
	var page struct {
		Enonce     Enonce                   `json:"enonce"`
		Parameters struct{ Variables []Rp } `json:"parameters"`
	}
	err := json.Unmarshal([]byte(jsonInput), &page)
	tu.AssertNoErr(t, err)

	asParams := make(Parameters, len(page.Parameters.Variables))
	for i, p := range page.Parameters.Variables {
		asParams[i] = p
	}
	vars, err := asParams.ToMap().Instantiate()
	tu.AssertNoErr(t, err)

	// test that the min and max function are properly evaluated
	tu.Assert(t, !strings.Contains(vars[ex.NewVarI('s', "1")].String(), "min"))
	tu.Assert(t, !strings.Contains(vars[ex.NewVarI('s', "2")].String(), "max"))
}

func TestInstantiateTextFormula(t *testing.T) {
	b := TextBlock{Parts: `Soit 
		$$ f(x) = 3x + 5 $$
		Calculer f'
	`}
	blocks, err := Enonce{b}.InstantiateWith(nil)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(blocks) == 3)
}

func TestValidateImage(t *testing.T) {
	ok, _ := Enonce{ImageBlock{"ùmslùsmd", 20}}.validate(&ex.RandomParameters{})
	tu.Assert(t, !ok)
	ok, _ = Enonce{ImageBlock{"https://test.fr", 0}}.validate(&ex.RandomParameters{})
	tu.Assert(t, !ok)
	ok, _ = Enonce{ImageBlock{"https://test.fr", 100}}.validate(&ex.RandomParameters{})
	tu.Assert(t, ok)
}

func TestIssue373(t *testing.T) {
	jsonInputEnonce := `[
	{"Data": {"Bold": false, "Parts": "Dériver la fonction $&t&$ définie et dérivable sur $\\R$ par :", "Italic": false, "Smaller": false}, "Kind": "TextBlock"}, 
	{"Data": {"Parts": "&t&(x)=&f_1*e^f_2&"}, "Kind": "FormulaBlock"}, 
	{"Data": {"Label": "", "Expression": "f_3*e^f_2", "ComparisonLevel": 2, "ShowFractionHelp": false}, "Kind": "ExpressionFieldBlock"} 
	]`
	jsonInputParameters := `[{"Data": {"variable": {"Name": 119, "Indice": ""}, "expression": "randint(1;3)"}, "Kind": "Rp"}, {"Data": {"variable": {"Name": 116, "Indice": ""}, "expression": "choiceFrom(f;g;h;w)"}, "Kind": "Rp"}, {"Data": {"variable": {"Name": 84, "Indice": ""}, "expression": "choiceFrom(F;G;H;w)"}, "Kind": "Rp"}, {"Data": {"variable": {"Name": 97, "Indice": ""}, "expression": "0"}, "Kind": "Rp"}, {"Data": {"variable": {"Name": 98, "Indice": ""}, "expression": "1"}, "Kind": "Rp"}, {"Data": {"variable": {"Name": 109, "Indice": ""}, "expression": "-1"}, "Kind": "Rp"}, {"Data": {"variable": {"Name": 112, "Indice": ""}, "expression": "0"}, "Kind": "Rp"}, {"Data": {"variable": {"Name": 102, "Indice": "1"}, "expression": "a*x+b"}, "Kind": "Rp"}, {"Data": {"variable": {"Name": 102, "Indice": "2"}, "expression": "m*x+p"}, "Kind": "Rp"}, {"Data": {"variable": {"Name": 102, "Indice": "3"}, "expression": "m*f_1"}, "Kind": "Rp"}]`

	var enonce Enonce
	err := json.Unmarshal([]byte(jsonInputEnonce), &enonce)
	tu.AssertNoErr(t, err)

	var parameters Parameters
	err = json.Unmarshal([]byte(jsonInputParameters), &parameters)
	tu.AssertNoErr(t, err)

	vars, err := parameters.ToMap().Instantiate()
	tu.AssertNoErr(t, err)

	instance, err := enonce.InstantiateWith(vars)
	tu.AssertNoErr(t, err)

	out := instance.EvaluateAnswer(client.QuestionAnswersIn{
		Data: client.Answers{0: client.ExpressionAnswer{Expression: "-e^(-x)"}},
	})
	tu.Assert(t, out.IsCorrect())

	out = instance.EvaluateAnswer(client.QuestionAnswersIn{
		Data: client.Answers{0: client.ExpressionAnswer{Expression: "(-1)e^(-x)"}},
	})
	tu.Assert(t, out.IsCorrect())

	out = instance.EvaluateAnswer(client.QuestionAnswersIn{
		Data: client.Answers{0: client.ExpressionAnswer{Expression: "(-1)e^(-1x)"}},
	})
	tu.Assert(t, out.IsCorrect())
}
