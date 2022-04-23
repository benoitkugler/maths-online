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
	Difficulty1 DifficultyTag = "__PRIVATE_1_STAR" // 1 étoile
	Difficulty2 DifficultyTag = "__PRIVATE_2_STAR" // 2 étoiles
	Difficulty3 DifficultyTag = "__PRIVATE_3_STAR" // 3 étoiles
)
