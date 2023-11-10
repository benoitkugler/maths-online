package expression

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

type varEvaluer interface {
	resolve(v Variable) (real, error)
}

var _ varEvaluer = (*evalResolver)(nil)

// Vars maps variables to a chosen value.
type Vars map[Variable]*Expr

// CompleteFrom adds entries in [other] not defined in [vs]
func (vs Vars) CompleteFrom(other Vars) {
	for c, v := range other {
		if _, has := vs[c]; !has {
			vs[c] = v
		}
	}
}

type evalResolver struct {
	defs Vars // source of the variables expressions

	seen    map[Variable]bool // variable that we are currently resolving
	results map[Variable]real // resulting values
}

// handle cycle
func (vrs Vars) resolver() *evalResolver {
	return &evalResolver{
		defs:    vrs,
		seen:    make(map[Variable]bool),
		results: make(map[Variable]real),
	}
}

func (vrs *evalResolver) resolve(v Variable) (real, error) {
	if value, has := vrs.results[v]; has {
		return value, nil
	}

	if vrs.seen[v] {
		return real{}, ErrCycleVariable{v}
	}

	expr, ok := vrs.defs[v]
	if !ok {
		return real{}, ErrMissingVariable{v}
	}

	vrs.seen[v] = true

	// recurse
	value, err := expr.evalReal(vrs)
	if err != nil {
		return real{}, err
	}

	vrs.results[v] = value

	return value, nil
}

// Evaluate uses the given variables values to evaluate the formula.
// If a variable is referenced in the expression but not in the bindings,
// `ErrMissingVariable` is returned.
// If a variable is in a cycle and can't be resolved, ErrCycleVariable` is returned.
// If the expression is not valid, like in randInt(2; -2), `ErrInvalidExpr` is returned
func (expr *Expr) Evaluate(vars Vars) (float64, error) {
	resolver := vars.resolver()
	return expr.evalFloat(resolver)
}

type ErrMissingVariable struct {
	Missing Variable
}

func (mv ErrMissingVariable) Error() string {
	return fmt.Sprintf("La variable %s n'est pas définie.", mv.Missing)
}

type ErrCycleVariable struct {
	InCycle Variable
}

func (mv ErrCycleVariable) Error() string {
	return fmt.Sprintf("La variable %s est présente dans un cycle.", mv.InCycle)
}

// mustEvaluate panics if the expression is invalid or if
// a variable is missing from `vars`.
func mustEvaluate(expr string, vars Vars) float64 {
	e := MustParse(expr)
	return e.mustEvaluate(vars)
}

// mustEvaluate panics if a variable is missing.
func (expr *Expr) mustEvaluate(vars Vars) float64 {
	out, err := expr.Evaluate(vars)
	if err != nil {
		panic(fmt.Sprintf("%s: %s", expr.String(), err))
	}
	return out
}

type singleVarResolver struct {
	v     Variable
	value float64
}

func (res singleVarResolver) resolve(v Variable) (real, error) {
	if res.v != v {
		return real{}, ErrMissingVariable{v}
	}
	return newReal(res.value), nil
}

type FunctionExpr struct {
	Function *Expr
	Variable Variable // usually x
}

// FunctionDefinition interprets an expression as mathematical function,
// where random parameters have been resolved
type FunctionDefinition struct {
	FunctionExpr         // instantiated version
	From, To     float64 // definition domain, with From <= To
}

// Closure returns a function computing f(x), where f is defined by the expression.
// The closure will silently return NaN if the expression is invalid.
func (f FunctionExpr) Closure() func(float64) float64 {
	return func(xValue float64) float64 {
		out, err := f.Function.evalFloat(singleVarResolver{f.Variable, xValue})
		if err != nil {
			return math.NaN()
		}
		return out
	}
}

// extrema returns an approximation of max |f(x)| on its definition domain.
// The approximation is exact for monotonous functions.
// `extrema` returns -1 if one of the values is not a finite number, or
// if the expression is invalid
func (f FunctionDefinition) extrema(isDiscrete bool) float64 {
	const nbSteps = 100
	fn := f.Closure()
	var max float64
	if isDiscrete {
		for x := math.Ceil(f.From); x <= math.Floor(f.To); x++ {
			fx := math.Abs(fn(x))
			if math.IsInf(fx, 0) || math.IsNaN(fx) {
				return -1
			}

			if fx > max {
				max = fx
			}
		}
	} else {
		step := (f.To - f.From) / nbSteps
		for i := 0; i <= nbSteps; i++ {
			x := f.From + float64(i)*step
			fx := math.Abs(fn(x))
			if math.IsInf(fx, 0) || math.IsNaN(fx) {
				return -1
			}

			if fx > max {
				max = fx
			}
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

func (expr *Expr) evalFloat(bindings varEvaluer) (float64, error) {
	r, err := expr.evalReal(bindings)
	return r.eval(), err
}

func (expr *Expr) evalReal(bindings varEvaluer) (real, error) {
	var (
		left, right = newRealInt(0), newRealInt(0) // 0 is a valid default value
		err         error
	)
	if expr.left != nil {
		left, err = expr.left.evalReal(bindings)
		if err != nil {
			return real{}, err
		}
	}
	if expr.right != nil {
		right, err = expr.right.evalReal(bindings)
		if err != nil {
			return real{}, err
		}
	}
	return expr.atom.eval(left, right, bindings)
}

func (op operator) eval(left, right real, _ varEvaluer) (real, error) {
	return op.evaluate(left, right), nil
}

// returns 1 if b is true
func evalBool(b bool) real {
	if b {
		return newRealInt(1)
	}
	return newRealInt(0)
}

func (op operator) evaluate(left, right real) real {
	// 0 is fine as default value for + and -
	// the other have mandatory left operands
	switch op {
	case equals:
		return evalBool(AreFloatEqual(left.eval(), right.eval()))
	case greater:
		return evalBool(left.eval() >= right.eval())
	case strictlyGreater:
		return evalBool(left.eval() > right.eval())
	case lesser:
		return evalBool(left.eval() <= right.eval())
	case strictlyLesser:
		return evalBool(left.eval() < right.eval())
	case plus:
		return sumReal(left, right)
	case minus:
		return minusReal(left, right)
	case mult:
		return multReal(left, right)
	case div:
		return divReal(left, right)
	case mod:
		leftInt, leftIsInt := IsInt(left.eval())
		rightInt, rightIsInt := IsInt(right.eval())
		if !(leftIsInt && rightIsInt) {
			return newRealInt(0)
		}
		return newRealInt(leftInt % rightInt)
	case rem:
		leftInt, leftIsInt := IsInt(left.eval())
		rightInt, rightIsInt := IsInt(right.eval())
		if !(leftIsInt && rightIsInt) {
			return newRealInt(0)
		}
		return newRealInt(leftInt / rightInt)
	case pow:
		return powReal(left, right.eval())
	case factorial:
		argInt, argIsInt := IsInt(left.eval())
		if !argIsInt {
			return newRealInt(0)
		}
		f := evalFactorial(argInt)
		return newRealInt(f)
	default:
		panic(exhaustiveOperatorSwitch)
	}
}

func (c constant) evalRat() real {
	switch c {
	case piConstant:
		return newReal(math.Pi)
	case eConstant:
		return newReal(math.E)
	default:
		panic(exhaustiveConstantSwitch)
	}
}

func (c constant) eval(_, _ real, _ varEvaluer) (real, error) {
	return c.evalRat(), nil
}

func (v Number) eval(_, _ real, _ varEvaluer) (real, error) { return newReal(float64(v)), nil }

func (indice) eval(_, _ real, _ varEvaluer) (real, error) {
	return real{}, errors.New("Une expression indicée ne peut pas être évaluée.")
}

func (matrix) eval(_, _ real, _ varEvaluer) (real, error) {
	return real{}, errors.New("Une matrice ne peut pas être évaluée.")
}

func (va Variable) eval(_, _ real, b varEvaluer) (real, error) {
	if b == nil {
		return real{}, ErrMissingVariable{Missing: va}
	}

	return b.resolve(va)
}

func roundTo(v float64, digits int) float64 {
	exp := math.Pow10(digits)
	return math.Round(v*exp) / exp
}

func (round roundFunc) eval(_, right real, _ varEvaluer) (real, error) {
	return newReal(roundTo(right.eval(), round.nbDigits)), nil
}

func (fn function) eval(_, right real, _ varEvaluer) (real, error) {
	arg := right.eval()
	switch fn {
	case logFn:
		return newReal(math.Log(arg)), nil
	case expFn:
		return newReal(math.Exp(arg)), nil
	case sinFn:
		return newReal(math.Sin(arg)), nil
	case cosFn:
		return newReal(math.Cos(arg)), nil
	case tanFn:
		return newReal(math.Tan(arg)), nil
	case asinFn:
		return newReal(math.Asin(arg)), nil
	case acosFn:
		return newReal(math.Acos(arg)), nil
	case atanFn:
		return newReal(math.Atan(arg)), nil
	case absFn:
		return newReal(math.Abs(arg)), nil
	case floorFn:
		return newRealInt(int(math.Floor(arg))), nil
	case sqrtFn:
		return newReal(math.Sqrt(arg)), nil
	case sgnFn:
		if arg > 0 {
			return newRealInt(1), nil
		} else if arg < 0 {
			return newRealInt(-1), nil
		}
		return newRealInt(0), nil
	case isPrimeFn:
		argInt, isInt := IsInt(arg)
		if !isInt {
			return newRealInt(0), nil
		}
		if isPrime(argInt) {
			return newRealInt(1), nil
		}
		return newRealInt(0), nil
	case forceDecimalFn:
		return real{val: right.eval(), isRational: false}, nil
	case detFn, traceFn, transposeFn, invertFn:
		return real{}, errors.New("internal error: matrice functions are not evaluable")
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

func startEnd(startE, endE *Expr, res varEvaluer) (float64, float64, error) {
	start, err := startE.evalFloat(res)
	if err != nil {
		return 0, 0, err
	}
	end, err := endE.evalFloat(res)
	if err != nil {
		return 0, 0, err
	}

	return start, end, nil
}

func minMax(args []*Expr, res varEvaluer) (float64, float64, error) {
	if len(args) == 0 {
		return 0, 0, ErrInvalidExpr{
			Reason: "min et max requierent au moins un argument",
		}
	}
	min, err := args[0].evalFloat(res)
	if err != nil {
		return 0, 0, err
	}
	max := min
	for _, arg := range args[1:] {
		v, err := arg.evalFloat(res)
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
func (r specialFunction) evalRat(res varEvaluer) (real, error) {
	switch r.kind {
	case randInt:
		start, end, err := startEnd(r.args[0], r.args[1], res)
		if err != nil {
			return real{}, err
		}

		err = r.validateStartEnd(start, end, 0)
		if err != nil {
			return real{}, err
		}
		return newRealInt(int(start) + rand.Intn(int(end-start)+1)), nil
	case randPrime:
		start, end, err := startEnd(r.args[0], r.args[1], res)
		if err != nil {
			return real{}, err
		}

		err = r.validateStartEnd(start, end, 0)
		if err != nil {
			return real{}, err
		}

		return newRealInt(generateRandPrime(int(start), int(end))), nil
	case randChoice:
		index := rand.Intn(len(r.args))
		return r.args[index].evalReal(res)
	case choiceFrom:
		// the parsing step ensure len(r.args) >= 2
		choice, err := choiceFromSelect(r.args, res)
		if err != nil {
			return real{}, err
		}
		return choice.evalReal(res)
	case randDenominator:
		index := rand.Intn(len(decimalDividors))
		return newRealInt(decimalDividors[index]), nil
	case minFn:
		min, _, err := minMax(r.args, res)
		return newReal(min), err
	case maxFn:
		_, max, err := minMax(r.args, res)
		return newReal(max), err
	case sumFn:
		start, end, err := startEnd(r.args[1], r.args[2], res)
		if err != nil {
			return real{}, err
		}
		err = r.validateStartEnd(start, end, 0)
		if err != nil {
			return real{}, err
		}

		// extract the variable
		indice, ok := r.args[0].atom.(Variable)
		if !ok {
			return real{}, errors.New("Le premier argument de sum() doit être une variable.")
		}
		expr := r.args[3]
		if start > end { //  ensure start <= end
			start, end = end, start
		}
		sum := newRealInt(0)
		for indiceVal := int(start); indiceVal <= int(end); indiceVal++ {
			out, err := expr.evalReal(singleVarResolver{indice, float64(indiceVal)})
			if err != nil {
				return real{}, fmt.Errorf("Impossible d'évaluer le terme d'indice %s = %d : %s", indice, indiceVal, err)
			}
			sum = sumReal(sum, out)
		}
		return sum, nil
	case matCoeff:
		mat, i, j := r.args[0], r.args[1], r.args[2]
		if mat, ok := mat.atom.(matrix); ok {
			n, m := mat.dims()
			i, err := evalIntInRange(i, res, 1, n)
			if err != nil {
				return real{}, fmt.Errorf("Le deuxième argument de coeff() doit être un indice de ligne : %s", err)
			}
			j, err := evalIntInRange(j, res, 1, m)
			if err != nil {
				return real{}, fmt.Errorf("Le troisième argument de coeff() doit être un indice de colonne : %s", err)
			}
			// human -> computer convention
			i--
			j--
			return mat[i][j].evalReal(res)
		} else {
			return real{}, fmt.Errorf("Le premier argument de coeff() doit être une matrice.")
		}
	case binomial:
		// the parsing step ensure len(r.args) == 2
		k, err := evalInt(r.args[0], res)
		if err != nil {
			return real{}, fmt.Errorf("Le premier argument de binom() doit être un entier (%s).", err)
		}
		n, err := evalInt(r.args[1], res)
		if err != nil {
			return real{}, fmt.Errorf("Le second argument de binom() doit être un entier (%s).", err)
		}
		return newRealInt(binomialCoefficient(k, n)), nil
	default:
		panic(exhaustiveSpecialFunctionSwitch)
	}
}

func evalInt(arg *Expr, res varEvaluer) (int, error) {
	i, err := arg.evalFloat(res)
	if err != nil {
		return 0, err
	}
	i_, ok := IsInt(i)
	if !ok {
		return 0, fmt.Errorf("valeur %g non entière", i)
	}
	return i_, nil
}

func evalIntInRange(arg *Expr, res varEvaluer, min, max int) (int, error) {
	i, err := evalInt(arg, res)
	if err != nil {
		return 0, err
	}
	if i < min || i > max {
		return 0, fmt.Errorf("valeur %d en dehors de [%d;%d]", i, min, max)
	}
	return i, nil
}

// evaluate the selector and return the expression at the index
// args must have length >= 2
func choiceFromSelect(args []*Expr, res varEvaluer) (choice *Expr, err error) {
	choices, selector := args[:len(args)-1], args[len(args)-1]
	index, err := evalIntInRange(selector, res, 1, len(choices))
	if err != nil {
		return nil, fmt.Errorf("Le dernier argument de la fonction choiceFrom est invalide : %s", err)
	}
	index -= 1 // using "human" convention
	return choices[index], nil
}

// return a random number
func (r specialFunction) eval(_, _ real, res varEvaluer) (real, error) {
	return r.evalRat(res)
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

// rat stores a rational number p/q, with q != 0
type rat struct {
	p, q int
}

func (r rat) eval() float64 { return float64(r.p) / float64(r.q) }

// update `r` to be in irreductible form
func (r *rat) reduce() {
	// special case for 0 / -5 : avoid spurious -0
	if r.p == 0 {
		r.p = 0
		r.q = 1
		return
	}

	// simplify integer denominators
	// commonDen = den1 * den2 / gcd()
	g := gcd(r.p, r.q)
	r.p /= g
	r.q /= g

	// simplify the minus
	if r.q < 0 {
		r.p = -r.p
		r.q = -r.q
	}
}

func (r rat) toExpr() *Expr {
	r.reduce()

	// avoid useless 4 / 1 or 0 / 1 fractions
	if r.q == 1 || r.p == 0 {
		return newNb(float64(r.p))
	}

	return &Expr{atom: div, left: newNb(float64(r.p)), right: newNb(float64(r.q))}
}

func sumRat(r1, r2 rat) rat {
	// den1, ok1 := IsInt(r1.q)
	// den2, ok2 := IsInt(r2.q)
	// if ok1 && ok2 {
	// 	// simplify integer denominators
	// 	// commonDen = den1 * den2 / gcd()
	// 	commonDen := float64(lcm(den1, den2))
	// 	factor1 := commonDen / r1.q
	// 	factor2 := commonDen / r2.q
	// 	return real{p: r1.p*factor1 + r2.p*factor2, q: commonDen}
	// }
	// // general case: do not simplify
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

func powRat(r rat, pow int) rat {
	powF := float64(pow)
	pF, qF := float64(r.p), float64(r.q)
	var pPow, qPow float64
	if powF < 0 { // invert the fraction
		pPow, qPow = math.Pow(qF, powF), math.Pow(pF, powF)
	} else {
		pPow, qPow = math.Pow(pF, powF), math.Pow(qF, powF)
	}

	return rat{p: int(pPow), q: int(qPow)}
}

// real store a real number, which may be represented as
// a rational, or not, depending on the flag isRational
type real struct {
	isRational bool
	rat        rat     // meaningful only if isRational is true
	val        float64 // meaningful only if isRational is false
}

func newRealInt(p int) real {
	return real{isRational: true, rat: rat{p: p, q: 1}}
}

// returns a rational number if v is an integer
func newReal(v float64) real {
	if vInt, isInt := IsInt(v); isInt {
		return newRealInt(vInt)
	}
	return real{val: v, isRational: false}
}

func (r real) eval() float64 {
	if r.isRational {
		return r.rat.eval()
	}
	return r.val
}

// return the better representation for the "rational", after reducing
// 8 / 1 -> 8
// 4 / 3 -> 4/3
// 3/4 -> 0.75
// 2.4 / 2 -> 1.2
func (r real) toExpr() *Expr {
	if r.isRational {
		return r.rat.toExpr()
	}

	// general case for real numbers : just returns the value
	return newNb(r.val)
}

func sumReal(r1, r2 real) real {
	if r1.isRational && r2.isRational {
		return real{isRational: true, rat: sumRat(r1.rat, r2.rat)}
	}
	// use eval to handle the case where r1 or r2 is rational
	return real{isRational: false, val: r1.eval() + r2.eval()}
}

// return r1 - r2
func minusReal(r1, r2 real) real {
	if r1.isRational && r2.isRational {
		return real{isRational: true, rat: minusRat(r1.rat, r2.rat)}
	}
	// use eval to handle the case where r1 or r2 is rational
	return real{isRational: false, val: r1.eval() - r2.eval()}
}

func multReal(r1, r2 real) real {
	if r1.isRational && r2.isRational {
		return real{isRational: true, rat: multRat(r1.rat, r2.rat)}
	}
	// use eval to handle the case where r1 or r2 is rational
	return real{isRational: false, val: r1.eval() * r2.eval()}
}

// return r1 / r2
func divReal(r1, r2 real) real {
	if r1.isRational && r2.isRational {
		return real{isRational: true, rat: divRat(r1.rat, r2.rat)}
	}
	// use eval to handle the case where r1 or r2 is rational
	return real{isRational: false, val: r1.eval() / r2.eval()}
}

func powReal(r real, pow float64) real {
	if powInt, isPowInt := IsInt(pow); r.isRational && isPowInt {
		return real{isRational: true, rat: powRat(r.rat, powInt)}
	}
	// use eval to handle the case where r is rational
	return real{isRational: false, val: math.Pow(r.eval(), pow)}
}

// performs some basic simplifications to convert expressions to numbers
// examples :
//
//	2*3 -> 6
//	ln(1) -> 0
//	3/2 -> 1.5
//
// due to the binary representation, some expressions cannot be simplified, such as
// (1 + x + 2)
func (expr *Expr) simplifyNumbers() {
	if expr == nil {
		return
	}

	expr.left.simplifyNumbers()
	expr.right.simplifyNumbers()

	op, ok := expr.atom.(operator)
	if !ok {
		return
	}

	// general case with two numbers
	var (
		leftNumber, rightNumber Number
		leftOK, rightOK         = true, true
	)
	if expr.left != nil {
		leftNumber, leftOK = expr.left.atom.(Number)
	}
	if expr.right != nil {
		rightNumber, rightOK = expr.right.atom.(Number)
	}

	if leftOK && rightOK {
		res := op.evaluate(newReal(float64(leftNumber)), newReal(float64(rightNumber)))
		// ensure fractions are converted
		*expr = *newNb(res.eval())
	}
}

// IsFraction returns true for expression of the form (...) / (...)
func (expr *Expr) IsFraction() bool {
	return expr.atom == div
}

func evalFactorial(n int) int {
	f := 1
	for k := 2; k <= n; k++ {
		f *= k
	}
	return f
}
