// Package expression provides support for parsing,
// evaluating and comparing simple mathematical expressions.
package expression

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	exhaustiveSpecialFunctionSwitch = "specialFunction"
	exhaustiveFunctionSwitch        = "function"
	exhaustiveOperatorSwitch        = "operator"
	exhaustiveConstantSwitch        = "constant"
	exhaustiveSymbolSwitch          = "symbol"
	exhaustiveTokenSwitch           = "token"
	exhaustiveAtomSwitch            = "atom"

	exhaustiveIntrinsicSwitch = "intrinsic"
)

// Expr is a parsed mathematical expression
type Expr struct {
	left, right *Expr
	atom        atom
}

// Serialize returns the expression as text.
// It is meant to be used for internal exchange; see
// String() and AsLaTex() for display.
func (expr *Expr) Serialize() string {
	if expr == nil {
		return ""
	}
	return expr.atom.serialize(expr.left, expr.right)
}

// Copy returns a deep copy of the expression.
func (expr *Expr) Copy() *Expr {
	if expr == nil {
		return nil
	}

	out := *expr
	out.left = expr.left.Copy()
	out.right = expr.right.Copy()
	return &out
}

// returns `true` if both expression are the same (structurally, not mathematicaly)
func (expr *Expr) equals(other *Expr) bool {
	if expr == other {
		return true
	}

	if expr == nil && other != nil || expr != nil && other == nil {
		return false
	}

	if expr.atom.String() != other.atom.String() {
		return false
	}

	return expr.left.equals(other.left) && expr.right.equals(other.right)
}

// atom is either an operator, a function,
// a variable, a predefined constant or a numerical value
type atom interface {
	fmt.Stringer // used to compare atom in equals

	// given the serialized form of the left and right terms,
	// serialize returns plain text, prettified, valid form
	// of the expression node
	serialize(left, right *Expr) string

	lexicographicOrder() int // smaller is first; unique among concrete types

	// return a value is possible as rational, so that
	// it may be simplified by subsequent operations
	eval(left, right rat, context ValueResolver) (rat, error)

	asLaTeX(left, right *Expr) string
}

func (Number) lexicographicOrder() int           { return 0 }
func (constant) lexicographicOrder() int         { return 1 }
func (operator) lexicographicOrder() int         { return 2 }
func (Variable) lexicographicOrder() int         { return 3 }
func (function) lexicographicOrder() int         { return 4 }
func (specialFunctionA) lexicographicOrder() int { return 5 }
func (roundFn) lexicographicOrder() int          { return 6 }
func (randVariable) lexicographicOrder() int     { return 7 }

// roundFn act as a function, but takes an integer parameter
// in addition to its regular parameter
type roundFn struct {
	nbDigits int
}

func (rd roundFn) String() string {
	return rd.serialize(nil, nil)
}

func (rd roundFn) serialize(_, right *Expr) string {
	return fmt.Sprintf("round(%s ; %d)", right.Serialize(), rd.nbDigits)
}

// randVariable is a special atom, randomly choosing a variable
// among the propositions
// it is only useful when used in random parameters definition,
// and is treated as zero elsewhere
type randVariable struct {
	choices []Variable
	// nil means choose uniformly in `choices`
	// if not nil, is is expected to return an indice into [1, len(choices)]
	selector *Expr
}

// String returns one of the two forms
// 			randSymbol(A;B;C)
//			choiceSymbol((A;B;C); randInt(1;4))
func (rv randVariable) String() string {
	var args []string
	for _, v := range rv.choices {
		args = append(args, v.String())
	}
	argS := strings.Join(args, ";")
	if rv.selector == nil { // shortcut form
		return fmt.Sprintf("randSymbol(%s)", argS)
	}
	return fmt.Sprintf("choiceSymbol((%s); %s)", argS, rv.selector.String())
}

func (rv randVariable) serialize(_, _ *Expr) string { return rv.String() }

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
	floorFn // floor (partie entière)
	sqrtFn
	sgnFn     // returns -1 0 or 1
	isZeroFn  // returns 1 is its argument is 0, 0 otherwise
	isPrimeFn // returns 0 or 1
	round     // round(<expr>; <digits>) : round(1.1256, 2) = 1.123

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
	case floorFn:
		return "floor"
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

func (fn function) serialize(_, right *Expr) string {
	return fn.String() + "(" + right.Serialize() + ")"
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

// NewVar is a convenience constructor for a simple variable.
func NewVar(x rune) Variable { return Variable{Name: x} }

// NewVarI is a convenience constructor supporting indices.
func NewVarI(x rune, indice string) Variable { return Variable{Name: x, Indice: indice} }

// NewVarExpr is a convenience constructor converting a `Variable` to an expression.
func NewVarExpr(v Variable) *Expr { return &Expr{atom: v} }

func newVarExpr(r rune) *Expr { return NewVarExpr(NewVar(r)) }

func (v Variable) serialize(_, _ *Expr) string { return v.String() }

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

func (c constant) serialize(_, _ *Expr) string { return c.String() }

type Number float64

func newNb(v float64) *Expr { return &Expr{atom: Number(v)} }

// NewNb returns the one element expression containing
// the given number.
// For consistency with the parser, negative numbers are actually
// returned as -(...)
func NewNb(v float64) *Expr {
	if v < 0 {
		return &Expr{atom: minus, right: newNb(-v)}
	}
	return newNb(v)
}

func (v Number) String() string {
	const decimalSeparator = ","
	out := strconv.FormatFloat(RoundFloat(float64(v)), 'f', -1, 64)
	return strings.ReplaceAll(out, ".", decimalSeparator)
}

func (v Number) serialize(_, _ *Expr) string { return v.String() }

type specialFunctionA struct {
	kind specialFunction
	args []*Expr // the correct length of args is check during parsing
}

func (sf specialFunctionA) String() string {
	var name string
	switch sf.kind {
	case randInt:
		name = "randInt"
	case randPrime:
		name = "randPrime"
	case randChoice:
		name = "randChoice"
	case randDenominator:
		name = "randDecDen"
	case minFn:
		name = "min"
	case maxFn:
		name = "max"
	default:
		panic(exhaustiveSpecialFunctionSwitch)
	}

	var args []string
	for _, n := range sf.args {
		args = append(args, n.String())
	}

	return name + "(" + strings.Join(args, " ; ") + ")"
}

func (sf specialFunctionA) serialize(_, _ *Expr) string { return sf.String() }
