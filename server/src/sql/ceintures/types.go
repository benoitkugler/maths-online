package ceintures

const (
	NbDomains = Matrices + 1 // gomacro:no-enum
	NbRanks   = Noire + 1    // gomacro:no-enum
)

type Level uint8

const (
	Seconde   Level = iota // Seconde
	Premiere               // Première
	Terminale              // Terminale
	PostBac                // Post-bac
)

// Domain is a sub topic, like "Calcul littéral",
// "Calcul mental".
//
// For now, we only support one scheme in maths.
type Domain uint8

const (
	CalculMental   Domain = iota // Calcul mental
	Puissances                   // Puissances et racines
	Fractions                    // Fractions
	Reduction                    // Réduction
	Factorisation                // Factorisation
	Developpement                // Développement
	IsolerVariable               // Isoler une variable
	Equations                    // Équations
	Inequations                  // Inéquations
	Derivation                   // Dérivation
	Matrices                     // Matrices et systèmes
)

// IsFor returns true if a student with [level] is
// qualified for the [Domain]
func (d Domain) IsFor(level Level) bool {
	switch d {
	case Derivation:
		return level >= Terminale
	case Matrices:
		return level >= PostBac
	default:
		return true
	}
}

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
	Rouge                 // Rouge
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
