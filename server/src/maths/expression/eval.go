package expression

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

type ValueResolver interface {
	resolve(v Variable) (*Expression, bool)
}

var _ ValueResolver = Variables{}

// Variables maps variables to a chosen value.
type Variables map[Variable]*Expression

func (vrs Variables) resolve(v Variable) (*Expression, bool) {
	value, ok := vrs[v]
	return value, ok
}

type ErrMissingVariable struct {
	Missing Variable
}

func (mv ErrMissingVariable) Error() string {
	return fmt.Sprintf("La variable %s n'est pas définie.", mv.Missing)
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

func (res singleVarResolver) resolve(v Variable) (*Expression, bool) {
	if res.v != v {
		return nil, false
	}
	return NewNb(res.value), true
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

const floatPrec = 1e-10

// AreFloatEqual returns `true` if v1 and v2 are equal up to
// a small threshold, so that floating point rouding errors are ignored
func AreFloatEqual(v1, v2 float64) bool {
	return math.Abs(v1-v2) < floatPrec
}

// RoundFloat returns `v` rounded to the precision used by `AreFloatEqual`.
// It should be used to avoid float imprecision when displaying numbers.
// It used internally when displaying expressions.
func RoundFloat(v float64) float64 {
	// round to avoid errors in previous computation
	// and use FormatFloat to avoid imprecision in this computation
	s := strconv.FormatFloat(math.Round(v/floatPrec)*floatPrec, 'f', 11, 64)
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
	r, err := expr.evalRat(bindings)
	return r.eval(), err
}

func (expr *Expression) evalRat(bindings ValueResolver) (rat, error) {
	var (
		left, right = newRat(0), newRat(0) // 0 is a valid default value
		err         error
	)
	if expr.left != nil {
		left, err = expr.left.evalRat(bindings)
		if err != nil {
			return rat{}, err
		}
	}
	if expr.right != nil {
		right, err = expr.right.evalRat(bindings)
		if err != nil {
			return rat{}, err
		}
	}
	return expr.atom.eval(left, right, bindings)
}

func (op operator) eval(left, right rat, _ ValueResolver) (rat, error) {
	return op.evaluate(left, right), nil
}

func (op operator) evaluate(left, right rat) rat {
	// 0 is fine as default value for + and -
	// the other have mandatory left operands
	switch op {
	case plus:
		return sumRat(left, right)
	case minus:
		return minusRat(left, right)
	case mult:
		return multRat(left, right)
	case div:
		return divRat(left, right)
	case mod:
		leftInt, leftIsInt := isInt(left.eval())
		rightInt, rightIsInt := isInt(right.eval())
		if !(leftIsInt && rightIsInt) {
			return newRat(0)
		}
		return newRat(float64(leftInt % rightInt))
	case rem:
		leftInt, leftIsInt := isInt(left.eval())
		rightInt, rightIsInt := isInt(right.eval())
		if !(leftIsInt && rightIsInt) {
			return newRat(0)
		}
		return newRat(float64(leftInt / rightInt))
	case pow:
		return powRat(left, right.eval())
	default:
		panic(exhaustiveOperatorSwitch)
	}
}

func (c constant) eval(_, _ rat, _ ValueResolver) (rat, error) {
	switch c {
	case piConstant:
		return newRat(math.Pi), nil
	case eConstant:
		return newRat(math.E), nil
	default:
		panic(exhaustiveConstantSwitch)
	}
}

func (v Number) eval(_, _ rat, _ ValueResolver) (rat, error) { return newRat(float64(v)), nil }

func (va Variable) eval(_, _ rat, b ValueResolver) (rat, error) {
	if b == nil {
		return rat{}, ErrMissingVariable{Missing: va}
	}

	out, has := b.resolve(va)
	if !has {
		return rat{}, ErrMissingVariable{Missing: va}
	}
	return out.evalRat(b)
}

func (rv randVariable) eval(_, _ rat, _ ValueResolver) (rat, error) {
	return newRat(0), nil
}

func roundTo(v float64, digits int) float64 {
	exp := math.Pow10(digits)
	return math.Round(v*exp) / exp
}

func (round roundFn) eval(_, right rat, _ ValueResolver) (rat, error) {
	return newRat(roundTo(right.eval(), round.nbDigits)), nil
}

func (fn function) eval(_, right rat, _ ValueResolver) (rat, error) {
	arg := right.eval()
	switch fn {
	case logFn:
		return newRat(math.Log(arg)), nil
	case expFn:
		return newRat(math.Exp(arg)), nil
	case sinFn:
		return newRat(math.Sin(arg)), nil
	case cosFn:
		return newRat(math.Cos(arg)), nil
	case tanFn:
		return newRat(math.Tan(arg)), nil
	case asinFn:
		return newRat(math.Asin(arg)), nil
	case acosFn:
		return newRat(math.Acos(arg)), nil
	case atanFn:
		return newRat(math.Atan(arg)), nil
	case absFn:
		return newRat(math.Abs(arg)), nil
	case floorFn:
		return newRat(math.Floor(arg)), nil
	case sqrtFn:
		return newRat(math.Sqrt(arg)), nil
	case sgnFn:
		if arg > 0 {
			return newRat(1), nil
		} else if arg < 0 {
			return newRat(-1), nil
		}
		return newRat(0), nil
	case isZeroFn:
		if arg == 0 {
			return newRat(1), nil
		}
		return newRat(0), nil
	case isPrimeFn:
		argInt, isInt := isInt(arg)
		if !isInt {
			return newRat(0), nil
		}
		if isPrime(argInt) {
			return newRat(1), nil
		}
		return newRat(0), nil
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

func minMax(args []*Expression, res ValueResolver) (float64, float64, error) {
	if len(args) == 0 {
		return 0, 0, ErrInvalidExpr{
			Reason: "min et max requierent au moins un argument",
		}
	}
	min, err := args[0].Evaluate(res)
	if err != nil {
		return 0, 0, err
	}
	max := min
	for _, arg := range args[1:] {
		v, err := arg.Evaluate(res)
		if err != nil {
			return 0, 0, err
		}
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	return min, max, nil
}

// return a random number
func (r specialFunctionA) eval(_, _ rat, res ValueResolver) (rat, error) {
	switch r.kind {
	case randInt:
		start, end, err := r.startEnd(res)
		if err != nil {
			return rat{}, err
		}

		err = r.validateStartEnd(start, end, 0)
		if err != nil {
			return rat{}, err
		}
		return newRat(start + float64(rand.Intn(int(end-start)+1))), nil
	case randPrime:
		start, end, err := r.startEnd(res)
		if err != nil {
			return newRat(0), err
		}

		err = r.validateStartEnd(start, end, 0)
		if err != nil {
			return newRat(0), err
		}

		return newRat(float64(generateRandPrime(int(start), int(end)))), nil
	case randChoice:
		index := rand.Intn(len(r.args))
		return r.args[index].evalRat(res)
	case randDenominator:
		index := rand.Intn(len(decimalDividors))
		return newRat(float64(decimalDividors[index])), nil
	case minFn:
		min, _, err := minMax(r.args, res)
		return newRat(min), err
	case maxFn:
		_, max, err := minMax(r.args, res)
		return newRat(max), err
	default:
		panic(exhaustiveSpecialFunctionSwitch)
	}
}

// --------------------------- numbers computations ---------------------------

// greatest common divisor (GCD) via Euclidean algorithm
func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

// rat has the form of a rational number p/q,
// but no assumption is actually made on the nature of p and q
// a real x (non rational number) is represented by rat{p:x, q:1}
type rat struct {
	p float64
	q float64
}

func newRat(v float64) rat { return rat{p: v, q: 1} }

func (r rat) eval() float64 { return r.p / r.q }

// return the better representation for the "rational", after reducing
// 8 / 1 -> 8
// 4 / 3 -> 4/3
// 3/4 -> 0.75
// 2.4 / 2 -> 1.2
func (r rat) toExpr() *Expression {
	r.reduce()

	if r.q == 1 || r.p == 0 {
		return newNb(r.p)
	}

	// test if the evaluation is a decimal number
	val := r.eval()
	if !isFloatExceedingPrecision(val) {
		return newNb(val)
	}

	// else for integers, return a fraction
	_, ok1 := isInt(r.p)
	_, ok2 := isInt(r.q)
	if ok1 && ok2 {
		return &Expression{atom: div, left: newNb(r.p), right: newNb(r.q)}
	}

	// general case : evaluate
	return newNb(val)
}

// for integers number, update `r` to be in irreductible form
func (r *rat) reduce() {
	num, ok1 := isInt(r.p)
	den, ok2 := isInt(r.q)
	if ok1 && ok2 {
		// simplify integer denominators
		// commonDen = den1 * den2 / gcd()
		g := float64(gcd(num, den))
		r.p /= g
		r.q /= g
	}

	// simplify the minus
	if r.q < 0 {
		r.p = -r.p
		r.q = -r.q
	}
}

func sumRat(r1, r2 rat) rat {
	den1, ok1 := isInt(r1.q)
	den2, ok2 := isInt(r2.q)
	if ok1 && ok2 {
		// simplify integer denominators
		// commonDen = den1 * den2 / gcd()
		commonDen := float64(lcm(den1, den2))
		factor1 := commonDen / r1.q
		factor2 := commonDen / r2.q
		return rat{p: r1.p*factor1 + r2.p*factor2, q: commonDen}
	}
	// general case: do not simplify
	return rat{p: r1.p*r2.q + r2.p*r1.q, q: r1.q * r2.q}
}

// return r1 - r2
func minusRat(r1, r2 rat) rat {
	return sumRat(r1, rat{p: -r2.p, q: r2.q})
}

func multRat(r1, r2 rat) rat {
	return rat{p: r1.p * r2.p, q: r1.q * r2.q}
}

// return r1 / r2
func divRat(r1, r2 rat) rat {
	return rat{p: r1.p * r2.q, q: r1.q * r2.p}
}

func powRat(r rat, pow float64) rat {
	return rat{p: math.Pow(r.p, pow), q: math.Pow(r.q, pow)}
}

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
		res := op.evaluate(newRat(float64(leftNumber)), newRat(float64(rightNumber)))
		*expr = *res.toExpr()
	}
}
