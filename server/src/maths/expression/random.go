package expression

import (
	"fmt"
	"math/rand"
)

// InvalidRandomVariable is returned when instantiating
// invalid parameter definitions
type InvalidRandomVariable struct {
	Detail string
	Cause  Variable
}

func (irv InvalidRandomVariable) Error() string {
	return fmt.Sprintf("paramètre aléatoire %s invalide : %s", irv.Cause, irv.Detail)
}

// RandomParameters stores a set of random parameters definitions,
// which may be related, but cannot contain cycles.
// RandomParameters may be stored in JSON format.
type RandomParameters map[Variable]*Expression

var _ valueResolver = &randomVarResolver{}

type randomVarResolver struct {
	defs RandomParameters

	seen    map[Variable]bool   // variable that we are currently resolving
	results map[Variable]number // resulting values

	err error
}

func (rvv *randomVarResolver) resolve(v Variable) (float64, bool) {
	if rvv.err != nil { // skip
		return 0, false
	}

	// first, check if it has already been resolved by side effect
	if nb, has := rvv.results[v]; has {
		return float64(nb), true
	}

	if rvv.seen[v] {
		rvv.err = InvalidRandomVariable{
			Cause:  v,
			Detail: fmt.Sprintf("%s est présente dans un cycle et ne peut donc pas être calculée", v),
		}
		return 0, false
	}

	// start the resolution : to detect invalid cycles,
	// register the variable
	rvv.seen[v] = true

	expr, ok := rvv.defs[v]
	if !ok {
		rvv.err = InvalidRandomVariable{
			Cause:  v,
			Detail: fmt.Sprintf("%s n'est pas définie", v),
		}
		return 0, false
	}

	// recurse
	value := expr.Evaluate(rvv)

	// register the result
	rvv.results[v] = number(value)

	return value, true
}

// Instantiate generate a random version of the
// variables, resolving possible dependencies.
// It returns an InvalidRandomVariable error for invalid cycles, like a = a +1
// or a = b + 1; b = a
func (rv RandomParameters) Instantiate() (Variables, error) {
	resolver := randomVarResolver{
		defs:    rv,
		seen:    make(map[Variable]bool),
		results: make(map[Variable]number),
	}

	out := make(Variables, len(rv))
	for v := range rv {
		value, _ := resolver.resolve(v) // this triggers the evaluation of the expression

		if resolver.err != nil {
			return nil, resolver.err
		}

		out[v] = value
	}

	return out, nil
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

func randPrime(min, max int) int {
	choices := sieveOfEratosthenes(min, max)
	L := len(choices)
	index := rand.Intn(L)
	return choices[index]
}

func (rd random) validate(pos int) error {
	if rd.start > rd.end {
		return InvalidExpr{
			Reason: "ordre invalide entre les arguments de randInt",
			Pos:    pos,
		}
	}

	if rd.isPrime && rd.start < 0 {
		return InvalidExpr{
			Reason: "randPrime n'accepte que des nombres positifs",
			Pos:    pos,
		}
	}

	if rd.isPrime && len(sieveOfEratosthenes(rd.start, rd.end)) == 0 {
		return InvalidExpr{
			Reason: fmt.Sprintf("aucun nombre premier n'existe entre %d et %d", rd.start, rd.end),
			Pos:    pos,
		}
	}

	return nil
}
