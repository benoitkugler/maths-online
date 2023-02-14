// Package expression provides support for parsing,
// evaluating and comparing simple mathematical expressions.
package expression

import (
	"fmt"
	"math"
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

	exhaustiveCompoundSwitch = "compound"
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

	lexicographicOrder() int // smaller is first; unique among concrete types

	// given the serialized form of the left and right terms,
	// serialize returns plain text, prettified, valid form
	// of the expression node
	serialize(left, right *Expr) string

	// return a value if possible as rational, so that
	// it may be simplified by subsequent operations
	eval(left, right rat, context varEvaluer) (rat, error)

	asLaTeX(left, right *Expr) string
}

func (Number) lexicographicOrder() int          { return 0 }
func (constant) lexicographicOrder() int        { return 1 }
func (operator) lexicographicOrder() int        { return 2 }
func (Variable) lexicographicOrder() int        { return 3 }
func (function) lexicographicOrder() int        { return 4 }
func (indice) lexicographicOrder() int          { return 5 }
func (specialFunction) lexicographicOrder() int { return 6 }
func (roundFn) lexicographicOrder() int         { return 7 }

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

type operator uint8

const (
	// the order is the precedence of operators
	// used during parsing
	equals operator = iota
	greater
	strictlyGreater
	lesser
	strictlyLesser
	plus
	minus
	mult
	div
	mod       // modulo(a, x) := a % x
	rem       // remainder(a, x) := a // x
	pow       // x^2
	factorial // n!

	invalidOperator
)

func (op operator) String() string {
	switch op {
	case equals:
		return "=="
	case greater:
		return ">="
	case strictlyGreater:
		return ">"
	case lesser:
		return "<="
	case strictlyLesser:
		return "<"
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
	case factorial:
		return "!"
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
// Variable is also used to store a custom symbol.
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
	case piConstant: // favor the nicer π instead of pi
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
	if math.IsInf(float64(v), 1) {
		// do not write the + to avoid -+inf
		return "inf"
	}
	const decimalSeparator = "," // prefer french notation
	out := strconv.FormatFloat(RoundFloat(float64(v)), 'f', -1, 64)
	return strings.ReplaceAll(out, ".", decimalSeparator)
}

func (v Number) serialize(_, _ *Expr) string { return v.String() }

// <left>_{<right>}
// the parser only accepts a name as left argument
type indice struct{}

func (indice) String() string { return "_" }

func (indice) serialize(left, right *Expr) string {
	return fmt.Sprintf("%s_{%s}", left.Serialize(), right.Serialize())
}

type specialFunction struct {
	kind specialFunctionKind
	args []*Expr // the correct length of args is check during parsing
}

func (sf specialFunction) String() string {
	name := sf.kind.String()

	var args []string
	for _, n := range sf.args {
		args = append(args, n.String())
	}

	return name + "(" + strings.Join(args, " ; ") + ")"
}

func (sf specialFunction) serialize(_, _ *Expr) string { return sf.String() }
