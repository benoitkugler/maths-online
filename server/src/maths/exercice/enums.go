package exercice

import "github.com/benoitkugler/maths-online/maths/expression"

type TextKind uint8

const (
	Text       TextKind = iota // Text simple
	StaticMath                 // Code LaTeX
	Expression                 // Expression
)

type SignSymbol uint8

const (
	Nothing        SignSymbol = iota //
	Zero                             // 0
	ForbiddenValue                   // ||
)

type ComparisonLevel = expression.ComparisonLevel

const (
	Strict                ComparisonLevel = expression.Strict                // Exacte
	SimpleSubstitutions   ComparisonLevel = expression.SimpleSubstitutions   // Simple
	ExpandedSubstitutions ComparisonLevel = expression.ExpandedSubstitutions // Complète
)

type VectorPairCriterion uint8

const (
	VectorEquals     VectorPairCriterion = iota // Vecteurs égaux
	VectorColinear                              // Vecteurs colinéaires
	VectorOrthogonal                            // Vecteurs orthogonaux
)

const (
	Diff1 DifficultyTag = "\u2605"             // 1 étoile
	Diff2 DifficultyTag = "\u2605\u2605"       // 2 étoiles
	Diff3 DifficultyTag = "\u2605\u2605\u2605" // 3 étoiles
)
