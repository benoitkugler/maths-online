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
// to try a new one [Pending].
type Prerequisite struct {
	Need    Stage
	Pending Stage
}

type Scheme struct {
	Ps     []Prerequisite
	Levels [ce.NbDomains]ce.Level // level minimum
}

// we do not store the implicit links which are intern to one domain,
// such as (CalculMental, Blanche) -> (CalculMental, Jaune)
var mathScheme = Scheme{
	Ps: []Prerequisite{
		{Stage{ce.CalculMentalI, ce.VerteI}, Stage{ce.CalculMentalII, ce.Orange}},

		// TODO: this is not decided yet
		{Stage{ce.CalculMentalI, ce.Blanche}, Stage{ce.Equations, ce.Blanche}},
		{Stage{ce.Fractions, ce.Jaune}, Stage{ce.Equations, ce.Blanche}},
		{Stage{ce.Reduction, ce.Rouge}, Stage{ce.Equations, ce.Blanche}},
		{Stage{ce.CalculMentalII, ce.Blanche}, Stage{ce.Factorisation, ce.Jaune}},
		{Stage{ce.Fractions, ce.Blanche}, Stage{ce.Factorisation, ce.Orange}},
	},
	Levels: [ce.NbDomains]ce.Level{
		ce.Derivation: ce.Terminale,
		ce.Matrices:   ce.PostBac,
	},
}

// return, for each target, the list of prerequisite needed
func (sh Scheme) byTarget() map[Stage][]Stage {
	out := make(map[Stage][]Stage, len(sh.Ps))
	for _, pr := range sh.Ps {
		out[pr.Pending] = append(out[pr.Pending], pr.Need)
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
func (sh Scheme) Pending(advance ce.Advance, level ce.Level) (nexts []Stage) {
	byTarget := sh.byTarget()
	// try each next rank, and check if it is compatible
	// with prerequisites
	for domain_, rank := range advance {
		domain := ce.Domain(domain_)
		if rank+1 == ce.NbRanks || level < sh.Levels[domain_] { // reached max rank or wrong level
			continue
		}

		candidate := Stage{Domain: domain, Rank: rank + 1}
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

func (sh Scheme) suggestionIndex(nexts []Stage) int {
	if len(nexts) == 0 {
		return -1
	}
	// select the minimum rank, then minimum domain
	var (
		min   = nexts[0]
		index int
	)
	for i, stage := range nexts {
		if stage.Rank < min.Rank {
			min = stage
			index = i
		} else if stage.Rank == min.Rank && stage.Domain < min.Domain {
			min = stage
			index = i
		}
	}
	return index
}

func byStage(questions ce.Beltquestions) map[Stage][]ce.Beltquestion {
	out := map[Stage][]ce.Beltquestion{}
	for _, qu := range questions {
		key := Stage{qu.Domain, qu.Rank}
		out[key] = append(out[key], qu)
	}
	return out
}
