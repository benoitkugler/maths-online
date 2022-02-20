package expression

import (
	"math"
)

type VariablesBinding interface {
	Resolve(v Variable) float64
}

type variables map[Variable]float64

func (vrs variables) Resolve(v Variable) float64 { return vrs[v] }

// Evaluate uses the given variables values to evaluate the formula.
func (expr *Expression) Evaluate(bindings VariablesBinding) float64 {
	var left, right float64 // 0 is a valid default value
	if expr.left != nil {
		left = expr.left.Evaluate(bindings)
	}
	if expr.right != nil {
		right = expr.right.Evaluate(bindings)
	}
	return expr.atom.eval(left, right, bindings)
}

func (op operator) eval(left, right float64, _ VariablesBinding) float64 {
	return op.evaluate(left, right)
}

func (op operator) evaluate(left, right float64) float64 {
	// 0 is fine as default value for + and -
	// the other have mandatory left operands
	switch op {
	case plus:
		return left + right
	case minus:
		return left - right
	case mult:
		return left * right
	case div:
		return left / right
	case pow:
		return math.Pow(left, right)
	default:
		panic("unknown operator")
	}
}

func (c constant) eval(_, _ float64, _ VariablesBinding) float64 {
	switch c {
	case numberPi:
		return math.Pi
	case numberE:
		return math.E
	default:
		panic("unknown constant")
	}
}

func (v number) eval(_, _ float64, _ VariablesBinding) float64 { return float64(v) }

func (va Variable) eval(_, _ float64, b VariablesBinding) float64 { return b.Resolve(va) }

func (fn function) eval(_, right float64, _ VariablesBinding) float64 {
	arg := right
	switch fn {
	case log:
		return math.Log(arg)
	case exp:
		return math.Exp(arg)
	case sin:
		return math.Sin(arg)
	case cos:
		return math.Cos(arg)
	case abs:
		return math.Abs(arg)
	default:
		panic("unknown constant")
	}
}

// --------------------------- numbers computations ---------------------------

// performs some basic simplifications to convert expressions to numbers
// examples : 2*3 -> 6
// examples : ln(1) -> 0
// due to the binary representation, some expressions cannot be simplified, such as
// (1 + x + 2)
func (expr *Expression) simplifyNumbers() {
	if expr == nil {
		return
	}

	expr.left.simplifyNumbers()
	expr.right.simplifyNumbers()

	// we only simplify operators for now
	op, ok := expr.atom.(operator)
	if !ok {
		return
	}

	var (
		leftNumber number
		leftOK     bool
	)
	if expr.left == nil {
		leftOK = true // 0 is a valid default value
	} else {
		leftNumber, leftOK = expr.left.atom.(number)
	}
	rightNumber, rightOK := expr.right.atom.(number)

	if leftOK && rightOK {
		res := op.evaluate(float64(leftNumber), float64(rightNumber))
		*expr = Expression{atom: number(res)}
	}
}
