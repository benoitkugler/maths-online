package questions

import (
	"encoding/json"

	"github.com/benoitkugler/maths-online/maths/expression"
)

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

// By default a slice of SignSymbol is marshalled as string
// by Go, which is not recognized by the PSQL JSON constraints
func (s SignSymbol) MarshalJSON() ([]byte, error) { return json.Marshal(uint8(s)) }

type ComparisonLevel uint8

const (
	Strict                                = ComparisonLevel(expression.Strict)                // Exacte
	SimpleSubstitutions                   = ComparisonLevel(expression.SimpleSubstitutions)   // Simple
	ExpandedSubstitutions                 = ComparisonLevel(expression.ExpandedSubstitutions) // Complète
	AsLinearEquation      ComparisonLevel = ExpandedSubstitutions + 100
)

type VectorPairCriterion uint8

const (
	VectorEquals     VectorPairCriterion = iota // Vecteurs égaux
	VectorColinear                              // Vecteurs colinéaires
	VectorOrthogonal                            // Vecteurs orthogonaux
)
