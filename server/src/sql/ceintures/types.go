package ceintures

const (
	NbDomains = Developpement + 1 // gomacro:no-enum
	NbRanks   = Noire + 1         // gomacro:no-enum
)

// Domain is a sub topic, like "Calcul littéral",
// "Calcul mental".
//
// For now, we only support one scheme in maths.
type Domain uint8

const (
	// TODO: precise
	CalculMental  Domain = iota // Calcul mental
	Fractions                   // Fractions
	Factorisation               // Factorisation
	Developpement               // Développement
)

// Rank is the belt color, that is the level
// of progression in one [Domain].
type Rank uint8

const (
	StartRank Rank = iota // Départ
	Blanche               // Blanche
	Jaune                 // Jaune
	Orange                // Orange
	Verte                 // Verte
	Bleue                 // Bleue
	Marron                // Marron
	Noire                 // Noire
)

type Stat struct {
	Success uint16 // number of questions answered with success
	Failure uint16 // number of questions answered with failure
}

// Advance stores, for each [Domain], the rank
// the student has checked.
type Advance [NbDomains]Rank

type Stats [NbDomains][NbRanks]Stat // by domain and rank
