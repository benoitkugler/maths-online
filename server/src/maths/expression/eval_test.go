package expression

import (
	"math"
	"reflect"
	"testing"
)

func Test_Expression_eval(t *testing.T) {
	tests := []struct {
		expr     string
		bindings valueResolver
		want     float64
	}{
		{
			"3 + 2", nil, 5,
		},
		{
			"3 + exp(0)", nil, 4,
		},
		{
			"sin(0)", nil, 0,
		},
		{
			"cos(0)", nil, 1,
		},
		{
			"abs(-3)", nil, 3,
		},
		{
			"ln(e)", nil, 1,
		},
		{
			"4/2", nil, 2,
		},
		{
			"4/3", nil, 4. / 3,
		},
		{
			"4/3", nil, 4. / 3,
		},
		{
			"4 * 3", nil, 12,
		},
		{
			"4 ^ 3", nil, 64,
		},
		{
			"\u03C0 / 2", nil, math.Pi / 2,
		},
		{
			"1 + 2 * (3 + 2)", nil, 11,
		},
		{
			"1 + 1 * 3 ^ 3 * 2 - 1", nil, 54,
		},
		{
			"x + 2", Variables{'x': 4}, 6,
		},
		{
			"2 + 0 * randInt(1,3)", nil, 2,
		},
		{
			"4 * sgn(-1)", nil, -4,
		},
		{
			"sqrt(16) * sqrt(9)", nil, 4 * 3,
		},
		{
			"2 * sqrt(16) * sqrt(9) * sqrt(25)", nil, 2 * 4 * 3 * 5,
		},
		{
			"4 * sgn(-1) * sgn(1) * sgn(0)", nil, 0,
		},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)

		if got := expr.Evaluate(tt.bindings); got != tt.want {
			t.Errorf("node.eval() = %v, want %v", got, tt.want)
		}
	}
}

func Test_Expression_simplifyNumbers(t *testing.T) {
	tests := []struct {
		expr string
		want string
	}{
		{"3 + x", "3 + x"}, // no op
		{"3 + 2", "5"},
		{"3 - 2", "1"},
		{"x - 0", "x"},
		{"x + 0", "x"},
		{"0+x", "x"},
		{"x * 0", "0"},
		{"0*x", "0"},
		{"0/x", "0"},
		{"x * 1", "x"},
		{"x / 1", "x"},
		{"x ^ 1", "x"},
		{"1 ^ x", "1"},
		{"- 2", "-2"},
		{"3 / 4", "0.75"},
		{"1 + 2*(5 - 3 + 4)", "13"},
		{"1 + x + 2", "1 + x + 2"}, // need commutativity, not handled by simplifyNumbers
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		expr.simplifyNumbers()

		want := mustParse(t, tt.want)
		want.simplifyNumbers()
		if !reflect.DeepEqual(expr, want) {
			t.Errorf("node.simplifyNumbers() = %v, want %v", expr, tt.want)
		}
	}
}

func TestExpression_Substitute(t *testing.T) {
	tests := []struct {
		expr string
		vars valueResolver
		want string
	}{
		{"a + b", Variables{}, "a+b"},
		{"a + b", Variables{'a': 4}, "4+b"},
		{"a + b / 2*a", Variables{'a': 4}, "4+b/2*4"},
		{"a + b", Variables{'a': 4, 'b': 5}, "4+5"},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		expr.Substitute(tt.vars)

		want := mustParse(t, tt.want)
		if !reflect.DeepEqual(expr, want) {
			t.Errorf("Substitute(%s) = %v, want %v", tt.expr, expr, tt.want)
		}
	}
}
