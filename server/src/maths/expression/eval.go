package expression

import (
	"fmt"
	"math"
	"math/rand"
)

type ValueResolver interface {
	resolve(v Variable) (float64, bool)
}

var _ ValueResolver = Variables{}

// Variables maps variables to a chosen value.
type Variables map[Variable]ResolvedVariable

func (vrs Variables) resolve(v Variable) (float64, bool) {
	value, ok := vrs[v]
	if !ok || value.IsVariable {
		return 0, false
	}
	return value.N, ok
}

type MissingVariableErr struct {
	Missing Variable
}

func (mv MissingVariableErr) Error() string {
	return fmt.Sprintf("missing value for variable %s", mv.Missing)
}

// MustEvaluate panics if the expression is invalid or if
// a variable is missing from `vars`.
func MustEvaluate(expr string, vars Variables) float64 {
	e := MustParse(expr)
	return e.MustEvaluate(vars)
}

// MustEvaluate panics if a variable is missing.
func (expr *Expression) MustEvaluate(bindings ValueResolver) float64 {
	out, err := expr.Evaluate(bindings)
	if err != nil {
		panic(fmt.Sprintf("%s: %s", expr.String(), err))
	}
	return out
}

type singleVarResolver struct {
	v     Variable
	value float64
}

func (res singleVarResolver) resolve(v Variable) (float64, bool) {
	if res.v != v {
		return 0, false
	}
	return res.value, true
}

type FunctionExpr struct {
	Function *Expression
	Variable Variable // usually x
}

// FunctionDefinition interprets an expression as mathematical function
type FunctionDefinition struct {
	FunctionExpr
	From, To float64 // definition domain
}

// Closure returns a function computing f(x), where f is defined by the expression.
// The closure will panic if the expression depends on other variables
func (f FunctionExpr) Closure() func(float64) float64 {
	return func(xValue float64) float64 {
		return f.Function.MustEvaluate(singleVarResolver{f.Variable, xValue})
	}
}

// Extrema returns an approximation of max |f(x)| on its definition domain.
// The approximation is exact for monotonous functions.
// `Extrema` will panic if the expression if not a valid function.
func (f FunctionDefinition) Extrema() float64 {
	const nbSteps = 100
	fn := f.Closure()
	step := (f.To - f.From) / nbSteps
	var max float64
	for i := 0; i <= nbSteps; i++ {
		fx := math.Abs(fn(f.From + float64(i)*step))
		if fx > max {
			max = fx
		}
	}
	return max
}

// AreFloatEqual returns `true` if v1 and v2 are equal up to
// a small threshold, so that floating point rouding errors are ignored
func AreFloatEqual(v1, v2 float64) bool {
	const floatPrec = 1e-8
	return math.Abs(v1-v2) < floatPrec
}

// Evaluate uses the given variables values to evaluate the formula.
// If a variable is referenced in the expression but not in the bindings,
// `MissingVariableErr` is returned.
func (expr *Expression) Evaluate(bindings ValueResolver) (float64, error) {
	var (
		left, right float64 // 0 is a valid default value
		err         error
	)
	if expr.left != nil {
		left, err = expr.left.Evaluate(bindings)
		if err != nil {
			return 0, err
		}
	}
	if expr.right != nil {
		right, err = expr.right.Evaluate(bindings)
		if err != nil {
			return 0, err
		}
	}
	return expr.atom.eval(left, right, bindings)
}

func (op operator) eval(left, right float64, _ ValueResolver) (float64, error) {
	return op.evaluate(left, right), nil
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

func (c constant) eval(_, _ float64, _ ValueResolver) (float64, error) {
	switch c {
	case piConstant:
		return math.Pi, nil
	case eConstant:
		return math.E, nil
	default:
		panic(exhaustiveConstantSwitch)
	}
}

func (v Number) eval(_, _ float64, _ ValueResolver) (float64, error) { return float64(v), nil }

func (va Variable) eval(_, _ float64, b ValueResolver) (float64, error) {
	if b == nil {
		return 0, MissingVariableErr{Missing: va}
	}

	out, has := b.resolve(va)
	if !has {
		return 0, MissingVariableErr{Missing: va}
	}
	return out, nil
}

func (rv randVariable) eval(_, _ float64, _ ValueResolver) (float64, error) {
	return 0, nil
}

func (fn function) eval(left, right float64, _ ValueResolver) (float64, error) {
	arg := right
	switch fn {
	case logFn:
		return math.Log(arg), nil
	case expFn:
		return math.Exp(arg), nil
	case sinFn:
		return math.Sin(arg), nil
	case cosFn:
		return math.Cos(arg), nil
	case tanFn:
		return math.Tan(arg), nil
	case asinFn:
		return math.Asin(arg), nil
	case acosFn:
		return math.Acos(arg), nil
	case atanFn:
		return math.Atan(arg), nil
	case absFn:
		return math.Abs(arg), nil
	case sqrtFn:
		return math.Sqrt(arg), nil
	case sgnFn:
		if arg > 0 {
			return 1, nil
		} else if arg < 0 {
			return -1, nil
		}
		return 0, nil
	case isZeroFn:
		if arg == 0 {
			return 1, nil
		}
		return 0, nil
	case isPrimeFn:
		argInt, isInt := isInt(arg)
		if !isInt {
			return 0, nil
		}
		if isPrime(argInt) {
			return 1, nil
		}
		return 0, nil
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

const (
	maxDecDen       = 10_000
	thresholdDecDen = 100
)

var decimalDividors = generateDivisors(maxDecDen, thresholdDecDen)

// return a random number
func (r specialFunctionA) eval(_, _ float64, _ ValueResolver) (float64, error) {
	switch r.kind {
	case randInt:
		start, end := int(r.args[0]), int(r.args[1])
		return float64(start + rand.Intn(end-start+1)), nil
	case randPrime:
		start, end := int(r.args[0]), int(r.args[1])
		return float64(generateRandPrime(start, end)), nil
	case randChoice:
		index := rand.Intn(len(r.args))
		return float64(r.args[index]), nil
	case randDenominator:
		index := rand.Intn(len(decimalDividors))
		return float64(decimalDividors[index]), nil
	default:
		panic(exhaustiveSpecialFunctionSwitch)
	}
}

// --------------------------- numbers computations ---------------------------

// performs some basic simplifications to convert expressions to numbers
// examples :
//	2*3 -> 6
//  ln(1) -> 0
// due to the binary representation, some expressions cannot be simplified, such as
// (1 + x + 2)
func (expr *Expression) simplifyNumbers() {
	if expr == nil {
		return
	}

	expr.left.simplifyNumbers()
	expr.right.simplifyNumbers()

	op, ok := expr.atom.(operator)
	if !ok {
		return
	}

	left := expr.left
	if expr.left == nil { // 0 is a valid default value
		left = NewNumber(0)
	}
	right := expr.right

	// general case with two numbers
	leftNumber, leftOK := left.atom.(Number)
	rightNumber, rightOK := right.atom.(Number)
	if leftOK && rightOK {
		res := op.evaluate(float64(leftNumber), float64(rightNumber))
		*expr = Expression{atom: Number(res)}
	}
}
