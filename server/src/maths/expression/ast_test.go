package expression

import (
	"testing"

	"github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func TestPanics(t *testing.T) {
	testutils.ShouldPanic(t, func() { (invalidFn).eval(real{}, real{}, nil) })
	testutils.ShouldPanic(t, func() { (invalidConstant).eval(real{}, real{}, nil) })
	testutils.ShouldPanic(t, func() { (invalidOperator).eval(real{}, real{}, nil) })
	testutils.ShouldPanic(t, func() { (specialFunction{kind: invalidSpecialFunction}).eval(real{}, real{}, nil) })

	testutils.ShouldPanic(t, func() { _ = (invalidFn).String() })
	testutils.ShouldPanic(t, func() { _ = (invalidConstant).String() })
	testutils.ShouldPanic(t, func() { _ = (invalidOperator).String() })
	testutils.ShouldPanic(t, func() { _ = (specialFunction{kind: invalidSpecialFunction}).String() })

	testutils.ShouldPanic(t, func() { _ = (invalidFn).asLaTeX(nil, nil) })
	testutils.ShouldPanic(t, func() { _ = (invalidConstant).asLaTeX(nil, nil) })
	testutils.ShouldPanic(t, func() { _ = (invalidOperator).asLaTeX(nil, nil) })
	testutils.ShouldPanic(t, func() { _ = (invalidOperator).serialize(nil, nil) })

	testutils.ShouldPanic(t, func() { plus.needParenthesis(&Expr{}, false, true) })
	testutils.ShouldPanic(t, func() { shouldOmitTimes(nil, false, &Expr{}) })
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
