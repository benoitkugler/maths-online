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
	CalculMentalI  Domain = iota // Calcul mental I
	CalculMentalII               // Calcul mental II
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

// Rank is the belt color, that is the level
// of progression in one [Domain].
type Rank uint8

const (
	StartRank Rank = iota // Départ
	Blanche               // Blanche
	Jaune                 // Jaune
	Orange                // Orange
	VerteI                // Verte clair
	VerteII               // Verte foncée
	Bleue                 // Bleue
	Violet                // Violette
	Rouge                 // Rouge
	Marron                // Marron
	Noire                 // Noire
)

type Stat struct {
	Success uint16 // number of questions answered with success
	Failure uint16 // number of questions answered with failure
}

func (s *Stat) Add(other Stat) {
	s.Success += other.Success
	s.Failure += other.Failure
}

// Advance stores, for each [Domain], the rank
// the student has checked.
type Advance [NbDomains]Rank

type Stats [NbDomains][NbRanks]Stat // by domain and rank
