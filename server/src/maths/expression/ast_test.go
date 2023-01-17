package expression

import (
	"testing"

	"github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func TestPanics(t *testing.T) {
	testutils.ShouldPanic(t, func() { (invalidFn).eval(rat{}, rat{}, nil) })
	testutils.ShouldPanic(t, func() { (invalidConstant).eval(rat{}, rat{}, nil) })
	testutils.ShouldPanic(t, func() { (invalidOperator).eval(rat{}, rat{}, nil) })
	testutils.ShouldPanic(t, func() { (specialFunction{kind: invalidSpecialFunction}).eval(rat{}, rat{}, nil) })

	testutils.ShouldPanic(t, func() { _ = (invalidFn).String() })
	testutils.ShouldPanic(t, func() { _ = (invalidConstant).String() })
	testutils.ShouldPanic(t, func() { _ = (invalidOperator).String() })
	testutils.ShouldPanic(t, func() { _ = (specialFunction{kind: invalidSpecialFunction}).String() })

	testutils.ShouldPanic(t, func() { _ = (invalidFn).asLaTeX(nil, nil) })
	testutils.ShouldPanic(t, func() { _ = (invalidConstant).asLaTeX(nil, nil) })
	testutils.ShouldPanic(t, func() { _ = (invalidOperator).asLaTeX(nil, nil) })
	testutils.ShouldPanic(t, func() { _ = (invalidOperator).serialize(nil, nil) })

	testutils.ShouldPanic(t, func() { plus.needParenthesis(&Expr{}, false, true) })
	testutils.ShouldPanic(t, func() { shouldOmitTimes(false, &Expr{}) })
	testutils.ShouldPanic(t, func() {
		e := &Expr{atom: invalidOperator}
		e.simplifyNumbers()
	})

	testutils.ShouldPanic(t, func() {
		tk := tokenizer{currentToken: token{data: symbol(invalidSymbol)}}
		pr := parser{tk: &tk}
		pr.parseOneNode(true)
	})

	testutils.ShouldPanic(t, func() {
		tk := tokenizer{}
		pr := parser{tk: &tk}
		pr.parseOneNode(true)
	})

	testutils.ShouldPanic(t, func() {
		tk := newTokenizer([]byte{')'})
		pr := parser{tk: tk}
		pr.parseOneNode(true)
	})

	testutils.ShouldPanic(t, func() {
		mustEvaluate("x+2", nil)
	})
	testutils.ShouldPanic(t, func() {
		expr := MustParse("x+2")
		expr.mustEvaluate(nil)
	})

	testutils.ShouldPanic(t, func() {
		MustParse("x + ")
	})

	testutils.ShouldPanic(t, func() {
		(&Expr{}).isLinearTerm()
	})

	testutils.ShouldPanic(t, func() {
		specialFunction{kind: invalidSpecialFunction}.validate(0)
	})

	testutils.ShouldPanic(t, func() {
		AreCompoundsEquivalent(nil, &Expr{}, 0)
	})
}

func TestExpression_String(t *testing.T) {
	tests := []struct {
		expr string
		want string
	}{
		{"2 + x", "2 + x"},
		{"2 / x", "2 / x"},
		{"2 * x", "2x"},
		{"a * x", "ax"},
		{"2 * 3", "2 * 3"},
		{"2 - x", "2 - x"},
		{"2 ^ x", "2 ^ x"},
		{"2 ^ (x+1)", "2 ^ (x + 1)"},
		{"e * \u03C0", "e\u03C0"},
		{"\uE001", "\uE001"},

		{"exp(2)", "exp(2)"},
		{"sin(2)", "sin(2)"},
		{"cos(2)", "cos(2)"},
		{"abs(2)", "abs(2)"},
		{"sqrt(2)", "sqrt(2)"},
		{"2 + x + log(10)", "2 + x + log(10)"},
		{"- x + 3", "-x + 3"},
		{"1x", "x"},
		{"+x", "x"},
		{"min(2) + max(3)", "min(2) + max(3)"},
		{"-(-a)", "a"},
		{"floor(4)", "floor(4)"},
		{"x + (-4 + y)", "x - 4 + y"},
		{"(1<2)+(3>4)+(5<=6)+(7>=8)", "(1 < 2) + (3 > 4) + (5 <= 6) + (7 >= 8)"},
		{"1 + 2<3", "1 + 2 < 3"},
		{"-inf", "-inf"},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		got := expr.String()
		if got != tt.want {
			t.Errorf("Expression.String() = %v, want %v", got, tt.want)
		}

		expr2 := mustParse(t, got)
		if expr.String() != expr2.String() {
			t.Fatalf("inconsitent String() for %s", tt.expr)
		}
		if !AreExpressionsEquivalent(expr, expr2, SimpleSubstitutions) {
			t.Fatalf("inconsitent String() for %s:  %s", tt.expr, got)
		}
	}
}

func TestExpression_StringRoundtrip(t *testing.T) {
	for _, tt := range expressions {
		if tt.wantErr {
			continue
		}

		expr := mustParse(t, tt.expr)
		got := expr.String()
		expr2 := mustParse(t, got)
		if expr.String() != expr2.String() {
			t.Fatalf("inconsitent String() for %s", tt.expr)
		}
		if !AreExpressionsEquivalent(expr, expr2, SimpleSubstitutions) {
			t.Fatalf("inconsitent String() for %s:  %s", tt.expr, got)
		}
	}
}
