// Package ceintures defines the graph structure
// used to organize questions into a progression,
// and modelize prerequisites of each task.
package ceintures

import (
	ce "github.com/benoitkugler/maths-online/server/src/sql/ceintures"
)

// Stage stores the position of a question in
// a scheme
type Stage struct {
	Domain ce.Domain
	Rank   ce.Rank
}

// returns true if [adv] is at least at [loc]
func (loc Stage) isReached(adv ce.Advance) bool {
	return adv[loc.Domain] >= loc.Rank
}

// Prerequisite modelizes the required location, [Need]
// to try a new one [For].
type Prerequisite struct {
	Need Stage
	For  Stage
}

type Scheme []Prerequisite

// we do not store the implicit links which are intern to one domain,
// such as (CalculMental, Blanche) -> (CalculMental, Jaune)
var mathScheme = Scheme{
	{Stage{ce.CalculMental, ce.Blanche}, Stage{ce.Factorisation, ce.Jaune}},
	{Stage{ce.Fractions, ce.Blanche}, Stage{ce.Factorisation, ce.Orange}},
}

// return, for each target, the list of prerequisite needed
func (sh Scheme) byTarget() map[Stage][]Stage {
	out := make(map[Stage][]Stage, len(sh))
	for _, pr := range sh {
		out[pr.For] = append(out[pr.For], pr.Need)
	}
	return out
}

// Pending returns the locations a student with [advance]
// may start.
// The current (or lower) positions are not included.
//
// An empty slice is returned if the scheme is complete.
// Otherwise, a well-formed scheme will always propose
// at least one new [Location].
func (sh Scheme) Pending(advance ce.Advance) (nexts []Stage) {
	byTarget := sh.byTarget()
	// try each next rank, and check if it is compatible
	// with prerequisites
	for domain, rank := range advance {
		if rank+1 == ce.NbRanks { // reached max rank
			continue
		}
		candidate := Stage{Domain: ce.Domain(domain), Rank: rank + 1}
		needed := byTarget[candidate]

		reached := true
		for _, need := range needed {
			if !need.isReached(advance) { // prerequisite not fullfiled
				reached = false
				break
			}
		}

		if reached {
			nexts = append(nexts, candidate)
		}
	}

	return nexts
}
