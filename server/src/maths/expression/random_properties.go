package expression

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
)

// RandomParameters stores a set of random parameters definitions,
// which may be related, but cannot contain cycles.
type RandomParameters map[Variable]*Expr

// addAnonymousParam register the given expression under a new variable, not used yet,
// chosen among a private range
func (rp RandomParameters) addAnonymousParam(expr *Expr) Variable {
	for ru := firstPrivateVariable; ru > 0; ru++ {
		v := Variable{Name: ru}
		if _, has := rp[v]; !has {
			rp[v] = expr
			return v
		}
	}
	panic("implementation limit reached")
}

// return list of primes between min and max (inclusive)
func sieveOfEratosthenes(min, max int) (primes []int) {
	sieve := make(map[int]bool)
	for i := 2; i <= max; i++ {
		if sieve[i] {
			continue
		}
		if i >= min {
			primes = append(primes, i)
		}
		for k := i * i; k <= max; k += i {
			sieve[k] = true
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

func (rd specialFunction) validateStartEnd(start, end float64, pos int) error {
	switch rd.kind {
	case randInt, randPrime:
		start, okStart := IsInt(start)
		end, okEnd := IsInt(end)
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

// IsValidNumber evaluates the expression using `vars`,
// checking if it is a valid number.
// If `checkPrecision` is true, it also checks that the number is not exceeding the
// float precision used in `AreFloatEqual`.
// If `rejectInfinite` is true, it returns an error for +/-Inf
func (expr *Expr) IsValidNumber(vars Vars, checkPrecision, rejectInfinite bool) error {
	value, err := expr.Evaluate(vars)
	if err != nil {
		return err
	}

	if math.IsNaN(value) {
		return fmt.Errorf("L'expression %s produit une valeur invalide (NaN).", expr)
	}

	if rejectInfinite && math.IsInf(value, 0) {
		return fmt.Errorf("L'expression %s produit une valeur infinie.", expr)
	}

	if checkPrecision && isFloatExceedingPrecision(value) {
		return fmt.Errorf("L'expression %s produit un nombre non décimal (%f).", expr, value)
	}

	return nil
}

// IsValidProba is the same as IsValidNumber, but also checks the number
// is in [0;1]
func (expr *Expr) IsValidProba(vars Vars) error {
	value, err := expr.Evaluate(vars)
	if err != nil {
		return err
	}

	if isFloatExceedingPrecision(value) {
		return fmt.Errorf("L'expression %s produit un nombre non décimal (%f).", expr, value)
	}

	if !(0 <= value && value <= 1) {
		return fmt.Errorf("L'expression %s ne produit pas un nombre entre 0 et 1 (%f).", expr, value)
	}

	return nil
}

// AreSortedNumbers evaluates the expressions using `vars`, checking if all are valid numbers and are sorted (in ascending order)
func AreSortedNumbers(exprs []*Expr, vars Vars) error {
	values := make([]float64, len(exprs))

	for j, expr := range exprs {
		var err error
		values[j], err = expr.Evaluate(vars)
		if err != nil {
			return err
		}
		if math.IsNaN(values[j]) {
			return fmt.Errorf("L'expression %s produit une valeur invalide (NaN).", expr)
		}
	}

	if !sort.Float64sAreSorted(values) {
		return fmt.Errorf("Les expressions produisent des valeurs <b>non croissantes</b>.")
	}

	return nil
}

// IsValidIndex evaluates the expression using `vars`,
// then checks if it is usable as input in a slice of length `length`.
// Note that we adopt the mathematical convention, with indices starting at 1. Thus the result is
// valid if it is an integer in [1, length]
func (expr *Expr) IsValidIndex(vars Vars, length int) error {
	value, err := expr.Evaluate(vars)
	if err != nil {
		return err
	}
	if index, ok := IsInt(value); ok && 1 <= index && index <= length {
		return nil
	}
	return fmt.Errorf("L'expression %s ne définit pas un index valide dans une liste de longueur %d.", expr, length)
}

// IsValid evaluates the function expression using `vars`, and checks
// if the (estimated) extrema of |f| is less than `bound`, returning an error if not.
func (fn FunctionExpr) IsValid(domain Domain, vars Vars, bound float64) error {
	fnExpr := fn.Function.Copy()
	fnExpr.Substitute(vars)

	fromV, toV, err := domain.eval(vars)
	if err != nil {
		return err
	}

	if fromV >= toV {
		return fmt.Errorf("Les expressions %s ne définissent pas un intervalle valide (%f, %f).", domain, fromV, toV)
	}

	def := FunctionDefinition{
		FunctionExpr: FunctionExpr{Function: fnExpr, Variable: fn.Variable},
		From:         fromV,
		To:           toV,
	}
	ext := def.extrema()
	if ext == -1 {
		return fmt.Errorf("L'expression %s ne définit pas une fonction valide.", fnExpr)
	} else if ext > bound {
		return fmt.Errorf("L'expression %s prend des valeurs trop importantes (%f)", fnExpr, ext)
	}

	return nil
}

// If one of the expression bound is nil, it is interpreted as Infinity (no constraint).

type Domain struct {
	From, To *Expr
}

func (d Domain) String() string {
	from := "-Inf"
	if d.From != nil {
		from = d.From.String()
	}
	to := "+Inf"
	if d.To != nil {
		to = d.To.String()
	}
	return fmt.Sprintf("[%s;%s]", from, to)
}

func (d Domain) eval(vars Vars) (from, to float64, err error) {
	if d.From == nil {
		from = math.Inf(-1)
	} else {
		from, err = d.From.Evaluate(vars)
		if err != nil {
			return 0, 0, err
		}
	}

	if d.To == nil {
		to = math.Inf(+1)
	} else {
		to, err = d.To.Evaluate(vars)
		if err != nil {
			return 0, 0, err
		}
	}
	return from, to, nil
}

// IsIncludedIntoOne returns an error if [expr], after evaluation does not
// belong to one of [domains].
func (expr *Expr) IsIncludedIntoOne(domains []Domain, vars Vars) error {
	x, err := expr.Evaluate(vars)
	if err != nil {
		return err
	}
	var ds []string
	for _, d := range domains {
		dFrom, dTo, err := d.eval(vars)
		if err != nil {
			return err
		}
		if dFrom <= x && x <= dTo { // found it
			return nil
		}
		ds = append(ds, d.String())
	}
	return fmt.Errorf("L'abscisse %s n'appartient à aucun des domaines %s.", expr, strings.Join(ds, ", "))
}

// IsIncludedIntoOne returns an error if [d] is not included in any [domains].
func (d Domain) IsIncludedIntoOne(domains []Domain, vars Vars) error {
	dFrom, dTo, err := d.eval(vars)
	if err != nil {
		return err
	}

	var ds []string
	for _, other := range domains {
		otherFrom, otherTo, err := other.eval(vars)
		if err != nil {
			return err
		}
		if otherFrom <= dFrom && dTo <= otherTo { // found it
			return nil
		}
		ds = append(ds, other.String())
	}

	return fmt.Errorf("L'intervalle %s n'est inclut dans aucun des domaines %s.", d, strings.Join(ds, ", "))
}

// AreDisjointsDomains returns an error if the given intervals [from, to] are not disjoints.
// Domains must be valid, as defined by `FunctionExpr.IsValid`.
func AreDisjointsDomains(domains []Domain, vars Vars) error {
	intervals := make([][2]float64, len(domains))
	for i, ds := range domains {
		fromV, toV, err := ds.eval(vars)
		if err != nil {
			return err
		}
		intervals[i] = [2]float64{fromV, toV}
	}

	if err := checkIntervalsDisjoints(intervals); err != nil {
		i1, i2 := domains[err.index1], domains[err.index2]
		return fmt.Errorf("les expressions %s et %s ne définissent pas des domains disjoints", i1, i2)
	}

	return nil
}

type jointIntervals struct {
	index1, index2 int
}

// returns a non nil value if the given intervals are not disjoints
func checkIntervalsDisjoints(intervals [][2]float64) *jointIntervals {
	// keep track of the original indices
	type indexedInterval struct {
		from, to float64
		index    int
	}

	tmp := make([]indexedInterval, len(intervals))
	for i, v := range intervals {
		tmp[i] = indexedInterval{from: v[0], to: v[1], index: i}
	}

	// sort by interval start
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].from < tmp[j].from
	})

	// now check that the end of one interval is less than the start of the next
	for i := range tmp {
		if i == len(tmp)-1 {
			break
		}

		i1, i2 := tmp[i], tmp[i+1]
		if i1.to > i2.from {
			return &jointIntervals{index1: i1.index, index2: i2.index}
		}
	}

	return nil
}

// IsValidLinearEquation checks that, once instantiated, [expr]
// is a linear equation (such as 2x - 3y + t/2 - 0.5)
func (expr *Expr) IsValidLinearEquation(vars Vars) error {
	expr = expr.Copy()
	expr.Substitute(vars)
	_, err := expr.isLinearEquation()
	return err
}