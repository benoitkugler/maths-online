package exercice

type TextKind uint8

const (
	Text       TextKind = iota // Text simple
	StaticMath                 // Code LaTeX
	Expression                 // Expression
)

type SignSymbol uint8

const (
	Nothing        SignSymbol = iota // Vide
	Zero                             // Zéro
	ForbiddenValue                   // Valeur interdite
)
