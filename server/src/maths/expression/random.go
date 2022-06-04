package expression

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

// ErrInvalidRandomParameters is returned when instantiating
// invalid parameter definitions
type ErrInvalidRandomParameters struct {
	Detail string
	Cause  Variable
}

func (irv ErrInvalidRandomParameters) Error() string {
	return fmt.Sprintf("Paramètre aléatoire %s invalide : %s", irv.Cause, irv.Detail)
}

// RandomParameters stores a set of random parameters definitions,
// which may be related, but cannot contain cycles.
type RandomParameters map[Variable]*Expression

// addAnonymousParam register the given expression under a new variable, not used yet,
// chosen among a private range
func (rp RandomParameters) addAnonymousParam(expr *Expression) Variable {
	for ru := firstPrivateVariable; ru > 0; ru++ {
		v := Variable{Name: ru}
		if _, has := rp[v]; !has {
			rp[v] = expr
			return v
		}
	}
	panic("implementation limit reached")
}

var _ ValueResolver = (*randomVarResolver)(nil)

type randomVarResolver struct {
	defs RandomParameters

	seen    map[Variable]bool // variable that we are currently resolving
	results map[Variable]rat  // resulting values

	err error
}

func (rvv *randomVarResolver) resolve(v Variable) (*Expression, bool) {
	if rvv.err != nil { // skip
		return nil, false
	}

	// first, check if it has already been resolved by side effect
	if nb, has := rvv.results[v]; has {
		return nb.toExpr(), true
	}

	if rvv.seen[v] {
		rvv.err = ErrInvalidRandomParameters{
			Cause:  v,
			Detail: fmt.Sprintf("%s est présente dans un cycle et ne peut donc pas être calculée", v),
		}
		return nil, false
	}

	// start the resolution : to detect invalid cycles,
	// register the variable
	rvv.seen[v] = true

	expr, ok := rvv.defs[v]
	if !ok {
		rvv.err = ErrInvalidRandomParameters{
			Cause:  v,
			Detail: fmt.Sprintf("%s n'est pas définie", v),
		}
		return nil, false
	}

	// recurse
	value, _ := expr.evalRat(rvv)

	// register the result
	rvv.results[v] = value

	return value.toExpr(), true
}

// Instantiate generate a random version of the
// variables, resolving possible dependencies.
// It returns an `ErrInvalidRandomParameters` error for invalid cycles, like a = a +1
// or a = b + 1; b = a.
// By design, a set of random parameters is either always valid, or always invalid,
// meaning this function may be used once as validation step.
func (rv RandomParameters) Instantiate() (Variables, error) {
	resolver := randomVarResolver{
		defs:    rv,
		seen:    make(map[Variable]bool),
		results: make(map[Variable]rat),
	}

	out := make(Variables, len(rv))
	for v, expr := range rv {
		// special case for randVariable
		if randV, isRandVariable := expr.atom.(randVariable); isRandVariable {
			resolver.seen[v] = true
			out[v] = &Expression{atom: randV.choice()}
			continue
		}

		value, _ := resolver.resolve(v) // this triggers the evaluation of the expression

		if resolver.err != nil {
			return nil, resolver.err
		}

		out[v] = value
	}

	return out, nil
}

func (rv randVariable) choice() Variable {
	index := rand.Intn(len(rv))
	return rv[index]
}

// return list of primes between min and max (inclusive)
func sieveOfEratosthenes(min, max int) (primes []int) {
	b := make([]bool, max+1)
	for i := 2; i <= max; i++ {
		if b[i] == true {
			continue
		}
		if i >= min {
			primes = append(primes, i)
		}
		for k := i * i; k <= max; k += i {
			b[k] = true
		}
	}
	return
}

// generateRandPrime panics if no prime is between min and max
func generateRandPrime(min, max int) int {
	choices := sieveOfEratosthenes(min, max)
	L := len(choices)
	index := rand.Intn(L)
	return choices[index]
}

func generateDivisors(n, threshold int) (out []int) {
	max := int(math.Sqrt(float64(n)))
	for i := 1; i <= max; i++ {
		if n%i == 0 {
			if i <= threshold {
				out = append(out, i)
			}
			if n/i <= threshold {
				out = append(out, n/i)
			}
		}
	}
	return out
}

func (rd specialFunctionA) validateStartEnd(start, end float64, pos int) error {
	switch rd.kind {
	case randInt, randPrime:
		start, okStart := isInt(start)
		end, okEnd := isInt(end)
		if !(okStart && okEnd) {
			return ErrInvalidExpr{
				Reason: "randXXX attend deux entiers en paramètres",
				Pos:    pos,
			}
		}

		if start > end {
			return ErrInvalidExpr{
				Reason: "ordre invalide entre les arguments de randXXX",
				Pos:    pos,
			}
		}

		if rd.kind == randPrime && start < 0 {
			return ErrInvalidExpr{
				Reason: "randPrime n'accepte que des nombres positifs",
				Pos:    pos,
			}
		}

		if rd.kind == randPrime && len(sieveOfEratosthenes(start, end)) == 0 {
			return ErrInvalidExpr{
				Reason: fmt.Sprintf("aucun nombre premier n'existe entre %d et %d", start, end),
				Pos:    pos,
			}
		}
	}
	return nil
}

func (rd specialFunctionA) validate(pos int) error {
	switch rd.kind {
	case randInt, randPrime:
		if len(rd.args) < 2 {
			return ErrInvalidExpr{
				Reason: "randXXX attend deux paramètres",
				Pos:    pos,
			}
		}

		// eagerly try to eval start and end in case their are constant,
		// so that the error is detected during parameter setup
		start, end, err := rd.startEnd(nil)
		if err == nil {
			return rd.validateStartEnd(start, end, pos)
		}
	case randChoice:
		if len(rd.args) == 0 {
			return ErrInvalidExpr{
				Reason: "randChoice doit préciser au moins un argument",
				Pos:    pos,
			}
		}
	case randDenominator: // nothing to validate

	case minFn, maxFn:
		if len(rd.args) == 0 {
			return ErrInvalidExpr{
				Reason: "min et max requierent au moins un argument",
				Pos:    pos,
			}
		}
	default:
		panic(exhaustiveSpecialFunctionSwitch)
	}

	return nil
}

// ErrRandomTests is returned when a valid expression does always
// pass a given criteria
type ErrRandomTests struct {
	// frequency of successul tries, between 0 and 100
	SuccessFrequency int
}

func (e ErrRandomTests) Error() string {
	return fmt.Sprintf("success rate: %d", e.SuccessFrequency)
}

// IsValidNumber instantiates the expression using `parameters`, then evaluate the resulting
// expression, checking it is a valid finite number.
// If `checkPrecision` is true, it also checks that the numbers are not exceeding the
// float precision used in `AreFloatEqual`.
// `parameters` must be a valid set of parameters
func (expr *Expression) IsValidNumber(parameters RandomParameters, checkPrecision, rejectInfinite bool) error {
	const nbTries = 1000
	var nbSuccess int
	for i := 0; i < nbTries; i++ {
		ps, _ := parameters.Instantiate()
		value, err := expr.Evaluate(ps)
		if err != nil { // return early
			return err
		}

		isValid := !math.IsNaN(value)
		if rejectInfinite {
			isValid = isValid && !math.IsInf(value, 0)
		}
		if checkPrecision {
			isValid = isValid && !isFloatExceedingPrecision(value)
		}

		if isValid {
			nbSuccess++
		}
	}

	if nbSuccess != nbTries {
		return ErrRandomTests{nbSuccess * 100 / nbTries}
	}

	return nil
}

// IsValidProba is the same as IsValidNumber, but also checks the number
// are in [0;1]
func (expr *Expression) IsValidProba(parameters RandomParameters) (bool, int) {
	const checkPrecision = true
	const nbTries = 1000
	var nbSuccess int
	for i := 0; i < nbTries; i++ {
		ps, _ := parameters.Instantiate()
		value, err := expr.Evaluate(ps)

		isValid := err == nil && !(math.IsInf(value, 0) || math.IsNaN(value))

		if checkPrecision && isFloatExceedingPrecision(value) {
			continue
		}

		isValid = isValid && (0 <= value && value <= 1)

		if isValid {
			nbSuccess++
		}
	}
	return nbSuccess == nbTries, nbSuccess * 100 / nbTries
}

// AreSortedNumbers instantiates the expressions using `parameters`, then evaluate the resulting
// expression, checking if all are valid numbers and are sorted (in ascending order)
// It also returns the frequency of successul tries, in % (between 0 and 100)
// `parameters` must be a valid set of parameters
func AreSortedNumbers(exprs []*Expression, parameters RandomParameters) (bool, int) {
	const nbTries = 1000
	var nbSuccess int
	values := make([]float64, len(exprs))
outer:
	for i := 0; i < nbTries; i++ {
		ps, _ := parameters.Instantiate()

		for j, expr := range exprs {
			var err error
			values[j], err = expr.Evaluate(ps)
			if err != nil || math.IsNaN(values[j]) {
				continue outer
			}
		}

		if sort.Float64sAreSorted(values) {
			nbSuccess++
		}
	}
	return nbSuccess == nbTries, nbSuccess * 100 / nbTries
}

// IsValidIndex instantiates the expression using `parameters`, then evaluate the resulting
// expression and checks if it is usable as input in a slice of length `length`.
// Note that we adopt the mathematical convention, with indices starting at 1. Thus the result is
// valid if it is an integer in [1, length]
// It also returns the frequency of successul tries, in % (between 0 and 100)
// `parameters` must be a valid set of parameters
func (expr *Expression) IsValidIndex(parameters RandomParameters, length int) (bool, int) {
	const nbTries = 1000
	var nbSuccess int
	for i := 0; i < nbTries; i++ {
		ps, _ := parameters.Instantiate()
		value, _ := expr.Evaluate(ps)
		if index, ok := isInt(value); ok && 1 <= index && index <= length {
			nbSuccess++
		}
	}
	return nbSuccess == nbTries, nbSuccess * 100 / nbTries
}

// IsValidInteger instantiates the expression using `parameters`, then evaluate the resulting
// expression and checks if it yields an integer.
// It also returns the frequency of successul tries, in % (between 0 and 100)
// `parameters` must be a valid set of parameters
func (expr *Expression) IsValidInteger(parameters RandomParameters) (bool, int) {
	const nbTries = 1000
	var nbSuccess int
	for i := 0; i < nbTries; i++ {
		ps, _ := parameters.Instantiate()
		value, _ := expr.Evaluate(ps)
		if _, ok := isInt(value); ok {
			nbSuccess++
		}
	}
	return nbSuccess == nbTries, nbSuccess * 100 / nbTries
}

// AreFxsIntegers instantiates the expression using `parameters`, then evaluate the resulting
// expression replacing `x` by the values from `grid` and checks if the values are integers.
// It returns the frequency of successul tries, in % (between 0 and 100)
// `parameters` must be a valid set of parameters
func (expr *Expression) AreFxsIntegers(parameters RandomParameters, x Variable, grid []int) (bool, int) {
	const nbTries = 1000
	var nbSuccess int
	for i := 0; i < nbTries; i++ {
		ps, _ := parameters.Instantiate()
		// checks that all grid values are integers
		areIntegers := true
		for _, xValue := range grid {
			ps[x] = NewNb(float64(xValue))
			value, _ := expr.Evaluate(ps)
			_, ok := isInt(value)
			areIntegers = areIntegers && ok
		}

		if areIntegers {
			nbSuccess++
		}
	}
	return nbSuccess == nbTries, nbSuccess * 100 / nbTries
}

// IsValid instantiates the function expression using `parameters`, then checks
// if the (estimated) extrema of |f| is less than `bound`.
// It returns the frequency of successul tries, in % (between 0 and 100)
// `parameters` must be a valid set of parameters.
func (fn FunctionDefinition) IsValid(parameters RandomParameters, bound float64) (bool, int) {
	const nbTries = 1000
	var nbSuccess int
	for i := 0; i < nbTries; i++ {
		ps, _ := parameters.Instantiate()

		fnExpr := fn.Function.copy()
		fnExpr.Substitute(ps)

		// by design, it is enough to check against one value
		// to see if the expression is valid
		_, err := fnExpr.Evaluate(singleVarResolver{v: fn.Variable, value: 0})
		if err != nil {
			continue
		}

		def := fn
		def.FunctionExpr.Function = fnExpr
		if ext := def.extrema(); ext == -1 || ext > bound {
			continue
		}

		nbSuccess++
	}
	return nbSuccess == nbTries, nbSuccess * 100 / nbTries
}
