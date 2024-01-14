package ceintures

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	ce "github.com/benoitkugler/maths-online/server/src/sql/ceintures"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func (sh Scheme) nextPrerequisite() (out [ce.NbDomains][ce.NbRanks]ce.Rank) {
	for d := ce.Domain(0); d < ce.NbDomains; d++ {
		var hasP [ce.NbRanks]bool
		for _, pr := range sh {
			if pr.Pending.Domain == d {
				hasP[pr.Pending.Rank] = true
			}
		}
		for r := ce.Rank(0); r < ce.NbRanks; r++ {
			cursor := r + 1
			for ; cursor < ce.NbRanks; cursor++ {
				if hasP[cursor] {
					break
				}
			}

			out[d][r] = cursor - 1 // we have to stay at r and check the prerequisites
		}
	}
	return out
}

func (sh Scheme) isWellFormed(level ce.Level) bool {
	nextPs := sh.nextPrerequisite()
	var aux func(adv ce.Advance) bool
	aux = func(adv ce.Advance) bool {
		// jump accross the "obvious" steps
		for d, r := range adv {
			adv[d] = nextPs[d][r]
		}

		if adv == allNoire { // done !
			return true
		}

		nextSteps := sh.Pending(adv, level)
		if len(nextSteps) == 0 { // trap !
			return false
		}

		ok := true
		for _, ns := range nextSteps {
			adv[ns.Domain] = ns.Rank
			ok = ok && aux(adv) // recurse
		}

		return ok
	}

	return aux(ce.Advance{})
}

func TestWellFormed(t *testing.T) {
	brokenScheme := Scheme{
		{Need: Stage{ce.CalculMental, ce.Blanche}, Pending: Stage{ce.Developpement, ce.Jaune}},
		{Need: Stage{ce.Developpement, ce.Jaune}, Pending: Stage{ce.CalculMental, ce.Blanche}},
	}
	tu.Assert(t, !brokenScheme.isWellFormed(ce.Seconde))

	// test that, for every possible advance,
	// new locations are available
	tu.Assert(t, mathScheme.isWellFormed(ce.Seconde))
	tu.Assert(t, mathScheme.isWellFormed(ce.Premiere))
	tu.Assert(t, mathScheme.isWellFormed(ce.Terminale))
	tu.Assert(t, mathScheme.isWellFormed(ce.PostBac))
}

func pow(x, a int) int64 {
	return int64(math.Pow(float64(x), float64(a)))
}

func TestPrintNbAdvances(t *testing.T) {
	fmt.Println(pow(int(ce.NbRanks), int(ce.NbDomains))/1000000, "millions")
}

func TestPending(t *testing.T) {
	for _, test := range []struct {
		adv   ce.Advance
		level ce.Level
		got   []Stage
	}{
		{ // start
			ce.Advance{}, ce.PostBac, []Stage{
				{0, ce.Blanche},
				{1, ce.Blanche},
				{2, ce.Blanche},
				{3, ce.Blanche},
				{4, ce.Blanche},
				{5, ce.Blanche},
				{6, ce.Blanche},
				{7, ce.Blanche},
				{8, ce.Blanche},
				{9, ce.Blanche},
				{ce.Matrices, ce.Blanche},
			},
		},
		{
			ce.Advance{}, ce.Seconde, []Stage{
				{0, ce.Blanche},
				{1, ce.Blanche},
				{2, ce.Blanche},
				{3, ce.Blanche},
				{4, ce.Blanche},
				{5, ce.Blanche},
				{6, ce.Blanche},
				{7, ce.Blanche},
				{8, ce.Blanche},
				// {9, ce.Blanche},
				// {ce.Matrices, ce.Blanche},
			},
		},
		// TODO:
	} {
		tu.Assert(t, reflect.DeepEqual(mathScheme.Pending(test.adv, test.level), test.got))
	}
}

var allNoire ce.Advance = [...]ce.Rank{
	ce.Noire, ce.Noire, ce.Noire, ce.Noire, ce.Noire, ce.Noire,
	ce.Noire, ce.Noire, ce.Noire, ce.Noire, ce.Noire,
}
