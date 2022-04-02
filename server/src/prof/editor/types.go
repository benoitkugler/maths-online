package editor

import "github.com/benoitkugler/maths-online/maths/exercice"

// See exercice.TextBlock
type TextBlock struct {
	Parts  string
	IsHint bool
}

type (
	FormulaBlock        = exercice.FormulaBlock
	VariationTableBlock = exercice.VariationTableBlock
	SignTableBlock      = exercice.SignTableBlock
)
