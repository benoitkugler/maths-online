package expression

import (
	"math"
)

type ValueResolver interface {
	Resolve(v Variable) (float64, bool)
}

var _ ValueResolver = variables{}

type variables map[Variable]float64

func (vrs variables) Resolve(v Variable) (float64, bool) {
	value, ok := vrs[v]
	return value, ok
}

// Evaluate uses the given variables values to evaluate the formula.
func (expr *Expression) Evaluate(bindings ValueResolver) float64 {
	var left, right float64 // 0 is a valid default value
	if expr.left != nil {
		left = expr.left.Evaluate(bindings)
	}
	if expr.right != nil {
		right = expr.right.Evaluate(bindings)
	}
	return expr.atom.eval(left, right, bindings)
}

func (op operator) eval(left, right float64, _ ValueResolver) float64 {
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
		panic(exhaustiveOperatorSwitch)
	}
}

func (c constant) eval(_, _ float64, _ ValueResolver) float64 {
	switch c {
	case piConstant:
		return math.Pi
	case eConstant:
		return math.E
	default:
		panic(exhaustiveConstantSwitch)
	}
}

func (v number) eval(_, _ float64, _ ValueResolver) float64 { return float64(v) }

func (va Variable) eval(_, _ float64, b ValueResolver) float64 {
	out, _ := b.Resolve(va)
	return out
}

func (fn function) eval(_, right float64, _ ValueResolver) float64 {
	arg := right
	switch fn {
	case logFn:
		return math.Log(arg)
	case expFn:
		return math.Exp(arg)
	case sinFn:
		return math.Sin(arg)
	case cosFn:
		return math.Cos(arg)
	case absFn:
		return math.Abs(arg)
	default:
		panic(exhaustiveConstantSwitch)
	}
}

// partial evaluation a.k.a substitution

// Substitute replaces variables by the given values, that is
// the ones for which Resolve() returns `true`.
func (expr *Expression) Substitute(vars ValueResolver) {
	if expr == nil {
		return
	}
	expr.left.Substitute(vars)
	expr.right.Substitute(vars)

	if v, isVariable := expr.atom.(Variable); isVariable {
		value, has := vars.Resolve(v)
		if has {
			expr.atom = number(value)
		}
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
