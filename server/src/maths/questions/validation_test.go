package questions

import (
	"testing"

	ex "github.com/benoitkugler/maths-online/maths/expression"
)

func mustParseMany(exprs []string) []*ex.Expr {
	out := make([]*ex.Expr, len(exprs))
	for i, s := range exprs {
		out[i] = ex.MustParse(s)
	}
	return out
}

func testAllValid(t *testing.T, parameters ex.RandomParameters, v validator, expected bool) {
	t.Helper()

	allValid := true
	for range [100]int{} {
		vars, err := parameters.Instantiate()
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
		xs         []string
		parameters ex.RandomParameters
		allValid   bool
	}{
		{
			[]string{"1", "2", "a"}, ex.RandomParameters{ex.NewVar('a'): ex.MustParse("randInt(2;5)")}, true,
		},
		{
			[]string{"1", "2", "a"}, ex.RandomParameters{ex.NewVar('a'): ex.MustParse("randInt(0;5)")}, false,
		},
		{
			[]string{"1", "2", "b"}, ex.RandomParameters{ex.NewVar('a'): ex.MustParse("randInt(0;5)")}, false,
		},
	}
	for _, tt := range tests {
		v := variationTableValidator{
			xs:  mustParseMany(tt.xs),
			fxs: nil, // ignored here
		}

		testAllValid(t, tt.parameters, v, tt.allValid)

	}
}

func Test_figureValidator_validate(t *testing.T) {
	tests := []struct {
		pointNames []string
		references []string
		parameters ex.RandomParameters
		wantErr    bool
	}{
		{[]string{"A", "B"}, nil, nil, false},
		{[]string{"A", "R"}, nil, ex.RandomParameters{ex.NewVar('R'): ex.MustParse("randLetter(U;V)")}, false},
		{[]string{"A", "R"}, nil, ex.RandomParameters{ex.NewVar('R'): ex.MustParse("randLetter(U;A)")}, true},
		{[]string{"A_1", "R"}, nil, ex.RandomParameters{ex.NewVar('R'): ex.MustParse("randLetter(U;A_1)")}, true},
		{[]string{"A_c2", "R"}, nil, ex.RandomParameters{ex.NewVar('R'): ex.MustParse("randLetter(U;A_c2)")}, true},
		{[]string{"A", "B"}, []string{"A", "A"}, nil, false},
		{[]string{"A", "B"}, []string{"A", "C"}, nil, true},
		{[]string{"A", "R"}, []string{"R"}, ex.RandomParameters{ex.NewVar('R'): ex.MustParse("randLetter(U;V)")}, false},
		{[]string{"A", "B", "C"}, []string{"R"}, ex.RandomParameters{ex.NewVar('R'): ex.MustParse("randLetter(A;B;C)")}, false},
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
		parameters ex.RandomParameters
		grid       []string
		want       bool
	}{
		{"2x + 1", nil, []string{"-2", "-1", "0", "4"}, true},
		{"ax^2 - 2x + c", ex.RandomParameters{ex.NewVar('a'): ex.MustParse("randInt(2;4)"), ex.NewVar('c'): ex.MustParse("7")}, []string{"-2", "-1", "0", "4"}, true},
		{"2x + 0.5", nil, []string{"-2", "-1", "0", "4"}, false},
		{"x", nil, []string{"1", "1"}, false},
		{"x", ex.RandomParameters{ex.NewVar('a'): ex.MustParse("randInt(2;4)")}, []string{"a", "2"}, false},
		{"x", ex.RandomParameters{ex.NewVar('a'): ex.MustParse("randInt(2;4)")}, []string{"a", "5"}, true},
	}
	for _, tt := range tests {
		expr := ex.MustParse(tt.expr)
		v := functionPointsValidator{
			function: function{
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
