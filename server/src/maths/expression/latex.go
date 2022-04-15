package expression

import (
	"fmt"
)

//go:generate go run unicode-latex/gen.go

// LaTeXResolver returns the latex code for the given variable `v`.
type LaTeXResolver = func(v Variable) string

// DefaultLatexResolver maps from unicode values to LaTeX commands
// for usual maths symbols.
func DefaultLatexResolver(v Variable) string { return unicodeToLaTeX[v.Name] }

// AsLaTeX returns a valid LaTeX code displaying the expression.
// `res` is an optional mapping from variables to latex symbols
func (expr *Expression) AsLaTeX(res LaTeXResolver) string {
	if expr == nil {
		return ""
	}

	if res == nil {
		res = DefaultLatexResolver
	}

	expr.contractPlusMinus()
	expr.contractMinusMinus()
	expr.simplify0And1()

	return expr.atom.asLaTeX(expr.left, expr.right, res)
}

func addParenthesis(s string) string {
	return `\left(` + s + `\right)`
}

func (op operator) asLaTeX(left, right *Expression, res LaTeXResolver) string {
	leftCode, rightCode := left.AsLaTeX(res), right.AsLaTeX(res)
	switch op {
	case plus:
		if leftCode == "" {
			return rightCode // plus is optional
		}
		return fmt.Sprintf("%s + %s", leftCode, rightCode)
	case minus:
		if right.needParenthesis(op, false) {
			rightCode = addParenthesis(rightCode)
		}
		return fmt.Sprintf("%s - %s", leftCode, rightCode)
	case mult:
		if left.needParenthesis(op, true) {
			leftCode = addParenthesis(leftCode)
		}
		if right.needParenthesis(op, false) {
			rightCode = addParenthesis(rightCode)
		}

		// check for implicit multiplication
		// the only case where x is required is for two numbers
		if _, isRightNumber := right.atom.(Number); isRightNumber {
			return fmt.Sprintf(`%s \times %s`, leftCode, rightCode)
		}
		return fmt.Sprintf(`%s%s`, leftCode, rightCode)
	case div:
		return fmt.Sprintf(`\frac{%s}{%s}`, leftCode, rightCode)
	case mod:
		return fmt.Sprintf(`\text{mod}\left(%s, %s\right)`, leftCode, rightCode)
	case rem:
		return fmt.Sprintf(`\text{rem}\left(%s, %s\right)`, leftCode, rightCode)
	case pow:
		if left.needParenthesis(op, true) {
			leftCode = addParenthesis(leftCode)
		}
		if right.needParenthesis(op, false) {
			rightCode = addParenthesis(rightCode)
		}
		return fmt.Sprintf(`{%s}^{%s}`, leftCode, rightCode)
	default:
		panic(exhaustiveOperatorSwitch)
	}
}

func (fn function) asLaTeX(left, right *Expression, res LaTeXResolver) string {
	arg := right.AsLaTeX(res)
	switch fn {
	case logFn:
		return fmt.Sprintf(`\log\left(%s\right)`, arg)
	case expFn:
		return fmt.Sprintf(`\exp\left(%s\right)`, arg)
	case sinFn:
		return fmt.Sprintf(`\sin\left(%s\right)`, arg)
	case cosFn:
		return fmt.Sprintf(`\cos\left(%s\right)`, arg)
	case tanFn:
		return fmt.Sprintf(`\tan\left(%s\right)`, arg)
	case asinFn:
		return fmt.Sprintf(`\arcsin\left(%s\right)`, arg)
	case acosFn:
		return fmt.Sprintf(`\arccos\left(%s\right)`, arg)
	case atanFn:
		return fmt.Sprintf(`\arctan\left(%s\right)`, arg)
	case absFn:
		return fmt.Sprintf(`\left|%s\right|`, arg)
	case sqrtFn:
		return fmt.Sprintf(`\sqrt{%s}`, arg)
	case sgnFn:
		return fmt.Sprintf(`\text{sgn}\left(%s\right)`, arg)
	case isZeroFn:
		return fmt.Sprintf(`\text{isZero}\left(%s\right)`, arg)
	case isPrimeFn:
		return fmt.Sprintf(`\text{isPrime}\left(%s\right)`, arg)
	default:
		panic(exhaustiveFunctionSwitch)
	}
}

func (r specialFunctionA) asLaTeX(_, _ *Expression, _ LaTeXResolver) string {
	return fmt.Sprintf(`\text{%s}`, r.String())
}

func (v Variable) asLaTeX(_, _ *Expression, res LaTeXResolver) string {
	name := res(v)
	if v.Indice != "" {
		name += "_{" + v.Indice + "}"
	}
	return name
}

func (c constant) asLaTeX(_, _ *Expression, _ LaTeXResolver) string {
	switch c {
	case piConstant:
		return "\\pi"
	case eConstant:
		return "e"
	default:
		panic(exhaustiveConstantSwitch)
	}
}

func (v Number) asLaTeX(_, _ *Expression, _ LaTeXResolver) string { return v.String() }

// returns `true` is the expression is compound and requires parenthesis
// when used with `op`
// if `isLeft` is true, this is expr op ...
// else this is ... op expr
func (expr Expression) needParenthesis(op operator, isLeft bool) bool {
	switch atom := expr.atom.(type) {
	case Number, constant, function, Variable:
		return false
	case operator:
		switch op {
		case minus:
			return atom <= op // - or +
		case pow:
			if atom == pow {
				return isLeft // (2^3)^5, but 2^(3^5) is usually written 2^3^5
			}
			return atom < op
		default:
			return atom < op
		}
	default:
		panic(exhaustiveAtomSwitch)
	}
}
