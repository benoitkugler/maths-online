// Package expression provides support for parsing,
// evaluating and comparing simple mathematical expressions.
package expression

import (
	"fmt"
	"strconv"
)

const (
	exhaustiveFunctionSwitch = "function"
	exhaustiveOperatorSwitch = "operator"
	exhaustiveConstantSwitch = "constant"
	exhaustiveSymbolSwitch   = "symbol"
	exhaustiveTokenSwitch    = "token"
	exhaustiveAtomSwitch     = "atom"
)

// Expression is a parsed mathematical expression
type Expression struct {
	left, right *Expression
	atom        atom
}

// String returns a human readable form of the expression,
// mainly used for debug.
// It is also suitable as a storage format, since
// it always produces a valid expression output, whose parsing
// yields a structurally equal expression.
// See `AsLaTeX` for a better display format.
func (expr *Expression) String() string {
	var out string
	if expr.left != nil {
		out = expr.left.String()
	}
	out += fmt.Sprint(expr.atom)
	if expr.right != nil {
		out += expr.right.String()
	}
	out = "(" + out + ")"
	return out
}

// returns a deep copy
func (expr *Expression) copy() *Expression {
	if expr == nil {
		return nil
	}

	out := *expr
	out.left = expr.left.copy()
	out.right = expr.right.copy()
	return &out
}

// returns `true` if both expression are the same (structurally, not mathematicaly)
func (expr *Expression) equals(other *Expression) bool {
	if expr == other {
		return true
	}

	if expr == nil && other != nil || expr != nil && other == nil {
		return false
	}

	if expr.atom != other.atom {
		return false
	}

	return expr.left.equals(other.left) && expr.right.equals(other.right)
}

// atom is either an operator, a function,
// a variable, a predefined constant or a numerical value
type atom interface {
	fmt.Stringer

	lexicographicOrder() int // smaller is first; unique among concrete types
	eval(left, right float64, context ValueResolver) float64
	asLaTeX(left, right *Expression, res LaTeXResolver) string
}

func (operator) lexicographicOrder() int { return 5 }
func (random) lexicographicOrder() int   { return 4 }
func (function) lexicographicOrder() int { return 3 }
func (Variable) lexicographicOrder() int { return 2 }
func (constant) lexicographicOrder() int { return 1 }
func (number) lexicographicOrder() int   { return 0 }

type operator uint8

const (
	// the order is the precedence of operators
	// used during parsing
	plus operator = iota
	minus
	mult
	div
	pow // x^2
)

func (op operator) String() string {
	switch op {
	case plus:
		return "+"
	case minus:
		return "-"
	case mult:
		return "*"
	case div:
		return "/"
	case pow:
		return "^"
	default:
		panic(exhaustiveOperatorSwitch)
	}
}

type function uint8

const (
	logFn function = iota
	expFn
	sinFn
	cosFn
	absFn
	sqrtFn
	// sgnFn // return -1 0 or 1 // TODO:
	// round
)

func (fn function) String() string {
	switch fn {
	case logFn:
		return "log"
	case expFn:
		return "exp"
	case sinFn:
		return "sin"
	case cosFn:
		return "cos"
	case absFn:
		return "abs"
	case sqrtFn:
		return "sqrt"
	default:
		panic(exhaustiveFunctionSwitch)
	}
}

// Variable is a (one letter) mathematical variable,
// such as 'a', 'b' in (a + b)^2 or 'x' in 2x + 3
// Private Unicode points are also permitted, so that
// custom compounded symbols may be used.
type Variable rune

func newVariable(r rune) *Expression {
	return &Expression{atom: Variable(r)}
}

func (v Variable) String() string {
	// we have to output valid expression syntax
	return string(v)
}

type constant uint8

const (
	piConstant constant = iota
	eConstant
	// i
)

func (c constant) String() string {
	switch c {
	case piConstant:
		return string(piRune)
	case eConstant:
		return "e"
	default:
		panic(exhaustiveConstantSwitch)
	}
}

type number float64

func newNumber(v float64) *Expression {
	return &Expression{atom: number(v)}
}

func (v number) String() string {
	return strconv.FormatFloat(float64(v), 'f', -1, 64)
}

// random is a random parameter, used to create unique and distinct
// version of the same general formula
type random struct {
	start, end int // inclusive, only accepts number as arguments (not expression)
}

func (r random) String() string {
	return fmt.Sprintf("randInt(%d, %d)", r.start, r.end)
}
