package expression

import (
	"reflect"
	"testing"
)

func TestExpr_SyntaxHints(t *testing.T) {
	tests := []struct {
		expr string
		want SyntaxHints
	}{
		{"2 + 5 * 4 ", SyntaxHints{}},
		{"-inf", SyntaxHints{shInf: true}},
		{"+inf", SyntaxHints{shInf: true}},
		{"3 / 8", SyntaxHints{shFractionSimple: true}},
		{"(3+x) / 8", SyntaxHints{shFractionComplex: true}},
		{"(x+1)^5", SyntaxHints{shPower: true}},
		{"sqrt(2x) - x^3", SyntaxHints{shPower: true, shSqrt: true}},
		{"sqrt(x^4) - x^3", SyntaxHints{shPower: true, shSqrt: true}},
	}
	for _, tt := range tests {
		e := mustParse(t, tt.expr)
		if got := e.SyntaxHints(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Expr.SyntaxHints() = %v, want %v", got, tt.want)
		}
	}
}
