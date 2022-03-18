package expression

import (
	"math"
	"math/rand"
)

type valueResolver interface {
	resolve(v Variable) (float64, bool)
}

var _ valueResolver = Variables{}

// Variables maps variables to a chosen value.
type Variables map[Variable]float64

func (vrs Variables) resolve(v Variable) (float64, bool) {
	value, ok := vrs[v]
	return value, ok
}

// Evaluate uses the given variables values to evaluate the formula.
func (expr *Expression) Evaluate(bindings valueResolver) float64 {
	var left, right float64 // 0 is a valid default value
	if expr.left != nil {
		left = expr.left.Evaluate(bindings)
	}
	if expr.right != nil {
		right = expr.right.Evaluate(bindings)
	}
	return expr.atom.eval(left, right, bindings)
}

func (op operator) eval(left, right float64, _ valueResolver) float64 {
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
	case mod:
		leftInt, leftIsInt := isInt(left)
		rightInt, rightIsInt := isInt(right)
		if !(leftIsInt && rightIsInt) {
			return 0
		}
		return float64(leftInt % rightInt)
	case rem:
		leftInt, leftIsInt := isInt(left)
		rightInt, rightIsInt := isInt(right)
		if !(leftIsInt && rightIsInt) {
			return 0
		}
		return float64(leftInt / rightInt)
	case pow:
		return math.Pow(left, right)
	default:
		panic(exhaustiveOperatorSwitch)
	}
}

func (c constant) eval(_, _ float64, _ valueResolver) float64 {
	switch c {
	case piConstant:
		return math.Pi
	case eConstant:
		return math.E
	default:
		panic(exhaustiveConstantSwitch)
	}
}

func (v Number) eval(_, _ float64, _ valueResolver) float64 { return float64(v) }

func (va Variable) eval(_, _ float64, b valueResolver) float64 {
	out, _ := b.resolve(va)
	return out
}

func (fn function) eval(left, right float64, _ valueResolver) float64 {
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
	case sqrtFn:
		return math.Sqrt(arg)
	case sgnFn:
		if arg > 0 {
			return 1
		} else if arg < 0 {
			return -1
		}
		return 0
	case isPrimeFn:
		argInt, isInt := isInt(arg)
		if !isInt {
			return 0
		}
		if isPrime(argInt) {
			return 1
		}
		return 0
	default:
		panic(exhaustiveFunctionSwitch)
	}
}

func isPrime(n int) bool {
	if n < 0 {
		n = -n
	}
	for i := 2; i <= int(math.Floor(math.Sqrt(float64(n)))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return n > 1
}

// return a random number
func (r random) eval(_, _ float64, _ valueResolver) float64 {
	if r.isPrime {
		return float64(randPrime(r.start, r.end))
	}
	return float64(r.start + rand.Intn(r.end-r.start+1))
}

// partial evaluation a.k.a substitution

// Substitute replaces variables contained in `vars`.
func (expr *Expression) Substitute(vars Variables) {
	if expr == nil {
		return
	}
	expr.left.Substitute(vars)
	expr.right.Substitute(vars)

	if v, isVariable := expr.atom.(Variable); isVariable {
		value, has := vars[v]
		if has {
			expr.atom = Number(value)
		}
	}
}

// --------------------------- numbers computations ---------------------------

// performs some basic simplifications to convert expressions to numbers
// examples :
//	2*3 -> 6
//  ln(1) -> 0
// 	1 * x -> x
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

	left := expr.left
	if expr.left == nil { // 0 is a valid default value
		left = newNumber(0)
	}
	right := expr.right

	// multiplying or dividing by 1;
	// adding or substracting 0 are no-ops
	switch op {
	case plus:
		if left.atom == Number(0) { // 0 + x = x
			*expr = *expr.right
			return
		} else if right.atom == Number(0) { // x + 0 = x
			*expr = *expr.left
			return
		}
	case minus:
		if right.atom == Number(0) { // x - 0 = x
			*expr = *expr.left
			return
		}
	case mult:
		if left.atom == Number(1) { // 1 * x = x
			*expr = *expr.right
			return
		} else if right.atom == Number(1) { // x * 1 = x
			*expr = *expr.left
			return
		} else if left.atom == Number(0) { // 0 * x = 0
			*expr = Expression{atom: Number(0)}
			return
		} else if right.atom == Number(0) {
			*expr = Expression{atom: Number(0)}
			return
		}
	case div:
		if right.atom == Number(1) { // x / 1 = x
			*expr = *expr.left
			return
		} else if left.atom == Number(0) && right.atom != Number(0) { // 0 / x = 0
			*expr = Expression{atom: Number(0)}
			return
		}
	case pow:
		if right.atom == Number(1) { // x ^ 1 = x
			*expr = *expr.left
			return
		} else if left.atom == Number(1) { // 1 ^ x = 1
			*expr = Expression{atom: Number(1)}
			return
		}
	default:
		panic(exhaustiveOperatorSwitch)
	}

	// general case with two numbers
	leftNumber, leftOK := left.atom.(Number)
	rightNumber, rightOK := right.atom.(Number)
	if leftOK && rightOK {
		res := op.evaluate(float64(leftNumber), float64(rightNumber))
		*expr = Expression{atom: Number(res)}
	}
}
