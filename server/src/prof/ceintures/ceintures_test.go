package ceintures

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	ce "github.com/benoitkugler/maths-online/server/src/sql/ceintures"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

// return NbColors ^ NbDomains elements
func allAdvances() []ce.Advance {
	fillAt := func(tmp []ce.Advance, i ce.Domain) []ce.Advance {
		var out []ce.Advance
		for _, partial := range tmp {
			for c := ce.Rank(0); c < ce.NbRanks; c++ {
				extended := partial // copy
				extended[i] = c
				out = append(out, extended)
			}
		}
		return out
	}
	l := []ce.Advance{{}}
	for i := ce.Domain(0); i < ce.NbDomains; i++ {
		l = fillAt(l, i)
	}
	return l
}

func pow(x, a int) int {
	return int(math.Pow(float64(x), float64(a)))
}

func TestAllAdvances(t *testing.T) {
	advances := allAdvances()
	tu.Assert(t, len(advances) == pow(int(ce.NbRanks), int(ce.NbDomains)))
	fmt.Println(advances[:100])
}

func TestPending(t *testing.T) {
	for _, test := range []struct {
		adv ce.Advance
		got []Stage
	}{
		{ // start
			ce.Advance{}, []Stage{{0, ce.Blanche}, {1, ce.Blanche}, {2, ce.Blanche}, {3, ce.Blanche}},
		},
		// TODO:
	} {
		tu.Assert(t, reflect.DeepEqual(mathScheme.Pending(test.adv), test.got))
	}
}

var allNoire ce.Advance = [...]ce.Rank{
	ce.Noire, ce.Noire, ce.Noire, ce.Noire,
}

func TestWellFormed(t *testing.T) {
	// test that, for every possible advance,
	// new locations are available
	for _, adv := range allAdvances() {
		if adv == allNoire {
			tu.Assert(t, len(mathScheme.Pending(adv)) == 0)
		} else {
			tu.Assert(t, len(mathScheme.Pending(adv)) > 0)
		}
	}
}
