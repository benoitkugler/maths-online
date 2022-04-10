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

	exhaustiveIntrinsicSwitch = "intrinsic"
)

// Expression is a parsed mathematical expression
type Expression struct {
	left, right *Expression
	atom        atom
}

// String returns a human readable form of the expression.
// It is also suitable as a storage format, since
// it always produces a valid expression output, whose parsing
// yields a structurally equal expression.
// The returned error is always nil.
// See `AsLaTeX` for a better display format.
func (expr *Expression) String() string {
	return string(expr.serialize())
}

func (expr *Expression) serialize() []byte {
	out := []byte{'('}
	if expr.left != nil {
		out = append(out, expr.left.serialize()...)
	}
	out = append(out, fmt.Sprint(expr.atom)...)
	if expr.right != nil {
		out = append(out, expr.right.serialize()...)
	}
	out = append(out, ')')
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
	eval(left, right float64, context ValueResolver) (float64, error)
	asLaTeX(left, right *Expression, res LaTeXResolver) string
}

func (operator) lexicographicOrder() int { return 5 }
func (random) lexicographicOrder() int   { return 4 }
func (function) lexicographicOrder() int { return 3 }
func (Variable) lexicographicOrder() int { return 2 }
func (constant) lexicographicOrder() int { return 1 }
func (Number) lexicographicOrder() int   { return 0 }

type operator uint8

const (
	// the order is the precedence of operators
	// used during parsing
	plus operator = iota
	minus
	mult
	div
	mod // modulo(a, x) := a % x
	rem // remainder(a, x) := a // x
	pow // x^2

	invalidOperator
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
	case mod:
		return "%"
	case rem:
		return "//"
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
	tanFn
	acosFn // invert
	asinFn // invert
	atanFn // invert
	absFn
	sqrtFn
	sgnFn     // returns -1 0 or 1
	isZeroFn  // returns 1 is its argument is 0, 0 otherwise
	isPrimeFn // returns 0 or 1
	// round

	invalidFn
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
	case tanFn:
		return "tan"
	case acosFn:
		return "acos"
	case asinFn:
		return "asin"
	case atanFn:
		return "atan"
	case absFn:
		return "abs"
	case sqrtFn:
		return "sqrt"
	case sgnFn:
		return "sgn"
	case isZeroFn:
		return "isZero"
	case isPrimeFn:
		return "isPrime"
	default:
		panic(exhaustiveFunctionSwitch)
	}
}

// Variable is a (one letter) mathematical variable,
// such as 'a', 'b' in (a + b)^2 or 'x' in 2x + 3.
// Indices are also permitted, written in LaTeX format :
// x_A or x_AB
// Private Unicode points are also permitted, so that
// custom compounded symbols may be used.
type Variable struct {
	Indice string // optional
	Name   rune
}

const firstPrivateVariable rune = '\uE001'

func NewVariable(x rune) Variable { return Variable{Name: x} }

func newVarExpr(r rune) *Expression {
	return &Expression{atom: NewVariable(r)}
}

func (v Variable) String() string {
	// we have to output valid expression syntax
	out := string(v.Name)
	if v.Indice != "" {
		return out + "_" + v.Indice + " " // notice the white space to avoid x_Ay_A
	}
	return out
}

type constant uint8

const (
	piConstant constant = iota
	eConstant
	// i

	invalidConstant
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

type Number float64

// NewNumber returns the one element expression containing
// the given number
func NewNumber(v float64) *Expression {
	return &Expression{atom: Number(v)}
}

func (v Number) String() string {
	return strconv.FormatFloat(float64(v), 'f', -1, 64)
}

// random is an integer random parameter, used to create unique and distinct
// version of the same general formula
type random struct {
	start, end int  // inclusive, only accepts number as arguments (not expression)
	isPrime    bool // if true, only generate prime numbers
}

func (r random) String() string {
	if r.isPrime {
		return fmt.Sprintf("randPrime(%d; %d)", r.start, r.end)
	}
	return fmt.Sprintf("randInt(%d; %d)", r.start, r.end)
}
