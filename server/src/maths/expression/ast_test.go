package expression

import (
	"reflect"
	"testing"
)

func shouldPanic(t *testing.T, f func()) {
	t.Helper()

	defer func() { recover() }()
	f()
	t.Errorf("should have panicked")
}

func TestPanics(t *testing.T) {
	shouldPanic(t, func() { (invalidFn).eval(0, 0, nil) })
	shouldPanic(t, func() { (invalidConstant).eval(0, 0, nil) })
	shouldPanic(t, func() { (invalidOperator).eval(0, 0, nil) })

	shouldPanic(t, func() { _ = (invalidFn).String() })
	shouldPanic(t, func() { _ = (invalidConstant).String() })
	shouldPanic(t, func() { _ = (invalidOperator).String() })

	shouldPanic(t, func() { _ = (invalidFn).asLaTeX(nil, nil, nil) })
	shouldPanic(t, func() { _ = (invalidConstant).asLaTeX(nil, nil, nil) })
	shouldPanic(t, func() { _ = (invalidOperator).asLaTeX(nil, nil, nil) })

	shouldPanic(t, func() { Expression{}.needParenthesis(0, false) })
	shouldPanic(t, func() {
		e := &Expression{atom: invalidOperator}
		e.simplifyNumbers()
	})

	shouldPanic(t, func() {
		tk := tokenizer{nextToken: token{data: symbol(invalidSymbol)}}
		pr := parser{tk: &tk}
		pr.parseOneNode()
	})

	shouldPanic(t, func() {
		tk := tokenizer{}
		pr := parser{tk: &tk}
		pr.parseOneNode()
	})
}

func TestExpression_String(t *testing.T) {
	tests := []struct {
		expr string
		want string
	}{
		{"2 + x", "((2)+(x))"},
		{"2 / x", "((2)/(x))"},
		{"2 * x", "((2)*(x))"},
		{"2 - x", "((2)-(x))"},
		{"2 ^ x", "((2)^(x))"},
		{"e * \u03C0", "((e)*(\u03C0))"},
		{"\uE001", "(\uE001)"},

		{"exp(2)", "(exp(2))"},
		{"sin(2)", "(sin(2))"},
		{"cos(2)", "(cos(2))"},
		{"abs(2)", "(abs(2))"},
		{"sqrt(2)", "(sqrt(2))"},
		{"2 + x + log(10)", "(((2)+(x))+(log(10)))"},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		got := expr.String()
		if got != tt.want {
			t.Errorf("Expression.String() = %v, want %v", got, tt.want)
		}

		expr2 := mustParse(t, got)
		if !reflect.DeepEqual(expr, expr2) {
			t.Fatalf("inconsitent String() for %s", tt.expr)
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
		if !reflect.DeepEqual(expr, expr2) {
			t.Fatalf("inconsitent String() for %s", tt.expr)
		}
	}
}