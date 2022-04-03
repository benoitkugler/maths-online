package editor

import "github.com/benoitkugler/maths-online/maths/exercice"

//go:generate ../../../../../structgen/structgen -source=types.go -mode=ts:test.ts

type Block = exercice.Block

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
