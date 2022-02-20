// Package expression provides support for parsing,
// evaluating and comparing simple mathematical expressions.
package expression

import "fmt"

// Expression is a parsed mathematical expression
type Expression struct {
	left, right *Expression
	atom        atom
}

// TODO: polish
func (expr *Expression) String() string {
	var out string
	if expr.left != nil {
		out = expr.left.String()
	}
	out += fmt.Sprint(expr.atom)
	if expr.right != nil {
		out += expr.right.String()
	}
	if expr.left != nil && expr.right != nil {
		out = "(" + out + ")"
	}
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
	lexicographicOrder() int // smaller is first; unique among concrete types
	eval(left, right float64, context VariablesBinding) float64
}

func (operator) lexicographicOrder() int { return 4 }
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
		panic("invalid operator")
	}
}

type function uint8

const (
	log function = iota
	exp
	sin
	cos
	abs
	// round
)

// Variable is a (one letter) mathematical variable,
// such as 'a', 'b' in (a + b)^2 or 'x' in 2x + 3
type Variable rune

func newVariable(r rune) *Expression {
	return &Expression{atom: Variable(r)}
}

type constant uint8

const (
	numberPi constant = iota
	numberE
	// i
)

type number float64

func newNumber(v float64) *Expression {
	return &Expression{atom: number(v)}
}
