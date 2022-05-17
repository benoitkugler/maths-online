package expression

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
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

type ErrMissingVariable struct {
	Missing Variable
}

func (mv ErrMissingVariable) Error() string {
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

// extrema returns an approximation of max |f(x)| on its definition domain.
// The approximation is exact for monotonous functions.
// `extrema` will panic if the expression if not a valid function.
// It returns -1 if one of the values is not a finite number.
func (f FunctionDefinition) extrema() float64 {
	const nbSteps = 100
	fn := f.Closure()
	step := (f.To - f.From) / nbSteps
	var max float64
	for i := 0; i <= nbSteps; i++ {
		fx := math.Abs(fn(f.From + float64(i)*step))
		if math.IsInf(fx, 0) || math.IsNaN(fx) {
			return -1
		}

		if fx > max {
			max = fx
		}
	}
	return max
}

const floatPrec = 1e-8

// AreFloatEqual returns `true` if v1 and v2 are equal up to
// a small threshold, so that floating point rouding errors are ignored
func AreFloatEqual(v1, v2 float64) bool {
	return math.Abs(v1-v2) < floatPrec
}

// RoundFloat returns `v` rounded to the precision used by `AreFloatEqual`.
// It should be used to avoid float imprecision when displaying numbers.
// It used internally when to display expressions.
func RoundFloat(v float64) float64 {
	// round to avoid errors in previous computation
	// and use FormatFloat to avoid imprecision if this computation
	s := strconv.FormatFloat(math.Round(v/floatPrec)*floatPrec, 'f', 8, 64)
	v, _ = strconv.ParseFloat(s, 64)
	return v
}

func isFloatExceedingPrecision(v float64) bool {
	// we rely on go format routine to avoid issue with floating
	// point computation
	s := fmt.Sprintf("%.9f", v) // 9 is one more than floatPrec exponent
	return s[len(s)-1] != '0'
}

// Evaluate uses the given variables values to evaluate the formula.
// If a variable is referenced in the expression but not in the bindings,
// `ErrMissingVariable` is returned.
// If the expression is not valid, like in randInt(2; -2), `ErrInvalidExpr` is returned
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
		return 0, ErrMissingVariable{Missing: va}
	}

	out, has := b.resolve(va)
	if !has {
		return 0, ErrMissingVariable{Missing: va}
	}
	return out, nil
}

func (rv randVariable) eval(_, _ float64, _ ValueResolver) (float64, error) {
	return 0, nil
}

func roundTo(v float64, digits int) float64 {
	exp := math.Pow10(digits)
	return math.Round(v*exp) / exp
}

func (round roundFn) eval(_, right float64, _ ValueResolver) (float64, error) {
	return roundTo(right, round.nbDigits), nil
}

func (fn function) eval(_, right float64, _ ValueResolver) (float64, error) {
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

func (r specialFunctionA) startEnd(res ValueResolver) (float64, float64, error) {
	start, err := r.args[0].Evaluate(res)
	if err != nil {
		return 0, 0, err
	}
	end, err := r.args[1].Evaluate(res)
	if err != nil {
		return 0, 0, err
	}

	return start, end, nil
}

// return a random number
func (r specialFunctionA) eval(_, _ float64, res ValueResolver) (float64, error) {
	switch r.kind {
	case randInt:
		start, end, err := r.startEnd(res)
		if err != nil {
			return 0, err
		}

		err = r.validateStartEnd(start, end, 0)
		if err != nil {
			return 0, err
		}
		return start + float64(rand.Intn(int(end-start)+1)), nil
	case randPrime:
		start, end, err := r.startEnd(res)
		if err != nil {
			return 0, err
		}

		err = r.validateStartEnd(start, end, 0)
		if err != nil {
			return 0, err
		}

		return float64(generateRandPrime(int(start), int(end))), nil
	case randChoice:
		index := rand.Intn(len(r.args))
		return r.args[index].Evaluate(res)
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
		left = NewNb(0)
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
