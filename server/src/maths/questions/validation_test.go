package questions

import (
	"testing"

	ex "github.com/benoitkugler/maths-online/server/src/maths/expression"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func mustParseMany(exprs []string) []*ex.Expr {
	out := make([]*ex.Expr, len(exprs))
	for i, s := range exprs {
		out[i] = ex.MustParse(s)
	}
	return out
}

func testAllValid(t *testing.T, parameters []Rp, v validator, expected bool) {
	t.Helper()

	ps := ex.NewRandomParameters()
	for _, p := range parameters {
		err := ps.ParseVariable(p.Variable, p.Expression)
		tu.AssertNoErr(t, err)
	}

	allValid := true
	for range [100]int{} {
		vars, err := ps.Instantiate()
		if err != nil {
			t.Fatal(err)
		}
		err = v.validate(vars)
		allValid = allValid && err == nil
	}
	if allValid != expected {
		t.Fatalf("expected %v, got %v", expected, allValid)
	}
}

func Test_variationTableValidator_validate(t *testing.T) {
	tests := []struct {
		xs        []string
		parameter Rp
		allValid  bool
	}{
		{
			[]string{"1", "2", "a"}, Rp{Variable: ex.NewVar('a'), Expression: "randInt(2;5)"}, true,
		},
		{
			[]string{"1", "2", "a"}, Rp{Variable: ex.NewVar('a'), Expression: "randInt(0;5)"}, false,
		},
		{
			[]string{"1", "2", "b"}, Rp{Variable: ex.NewVar('a'), Expression: "randInt(0;5)"}, false,
		},
	}
	for _, tt := range tests {
		v := variationTableValidator{
			xs:  mustParseMany(tt.xs),
			fxs: nil, // ignored here
		}
		testAllValid(t, []Rp{tt.parameter}, v, tt.allValid)
	}

	testAllValid(t, nil, variationTableValidator{
		fxs: mustParseMany([]string{"8", "inf", "-inf"}),
	}, true)
}

func Test_figureValidator_validate(t *testing.T) {
	tests := []struct {
		pointNames []string
		references []string
		parameters []Rp
		wantErr    bool
	}{
		{[]string{"A", "B"}, nil, nil, false},
		{[]string{"A", "R"}, nil, []Rp{{Variable: ex.NewVar('R'), Expression: "randChoice(U;V)"}}, false},
		{[]string{"A", "R"}, nil, []Rp{{Variable: ex.NewVar('R'), Expression: "randChoice(U;A)"}}, true},
		{[]string{"A_1", "R"}, nil, []Rp{{Variable: ex.NewVar('R'), Expression: "randChoice(U;A_1)"}}, true},
		{[]string{"A_c2", "R"}, nil, []Rp{{Variable: ex.NewVar('R'), Expression: "randChoice(U;A_c2)"}}, true},
		{[]string{"A", "B"}, []string{"A", "A"}, nil, false},
		{[]string{"A", "B"}, []string{"A", "C"}, nil, true},
		{[]string{"A", "R"}, []string{"R"}, []Rp{{Variable: ex.NewVar('R'), Expression: "randChoice(U;V)"}}, false},
		{[]string{"A", "B", "C"}, []string{"R"}, []Rp{{Variable: ex.NewVar('R'), Expression: "randChoice(A;B;C)"}}, false},
	}
	for _, tt := range tests {
		v := figureValidator{
			pointNames: mustParseMany(tt.pointNames),
			references: mustParseMany(tt.references),
		}

		testAllValid(t, tt.parameters, v, !tt.wantErr)
	}
}

func TestExpression_AreFxsIntegers(t *testing.T) {
	tests := []struct {
		expr       string
		parameters []Rp
		grid       []string
		want       bool
	}{
		{"2x + 1", nil, []string{"-2", "-1", "0", "4"}, true},
		{"ax^2 - 2x + c", []Rp{
			{Variable: ex.NewVar('a'), Expression: "randInt(2;4)"},
			{Variable: ex.NewVar('c'), Expression: "7"},
		}, []string{"-2", "-1", "0", "4"}, true},
		{"2x + 0.5", nil, []string{"-2", "-1", "0", "4"}, false},
		{"x", nil, []string{"1", "1"}, false},
		{"x", []Rp{{Variable: ex.NewVar('a'), Expression: "randInt(2;4)"}}, []string{"a", "2"}, false},
		{"x", []Rp{{Variable: ex.NewVar('a'), Expression: "randInt(2;4)"}}, []string{"a", "5"}, true},
	}
	for _, tt := range tests {
		expr := ex.MustParse(tt.expr)
		v := functionPointsValidator{
			function: functionValidator{
				FunctionExpr: ex.FunctionExpr{Function: expr, Variable: ex.NewVar('x')},
				domain: ex.Domain{
					From: ex.MustParse("-10"),
					To:   ex.MustParse("-10"),
				},
			},
			xGrid: mustParseMany(tt.grid),
		}
		testAllValid(t, tt.parameters, v, tt.want)
	}
}

func TestSequences(t *testing.T) {
	bl := FunctionsGraphBlock{
		SequenceExprs: []FunctionDefinition{
			{Function: "(-1)^n + (n/2)^2", Variable: ex.NewVar('n'), From: "-4", To: "4"},
		},
	}
	v, err := bl.setupValidator(&ex.RandomParameters{})
	tu.AssertNoErr(t, err)
	err = v.validate(nil)
	tu.AssertNoErr(t, err)
}
