package expression

import (
	"math"
	"sort"
	"strings"
)

// SyntaxHints is a set of [SyntaxHint]
type SyntaxHints map[SyntaxHint]bool

// Append add [other] to the set
func (sh SyntaxHints) Append(other SyntaxHints) {
	for k, v := range other {
		sh[k] = v
	}
}

// Text returns a french string using interpolation syntax $ $ for LaTeX
func (sh SyntaxHints) Text() string {
	var chunks []string
	for hint := range sh {
		chunks = append(chunks, hint.text())
	}
	sort.Strings(chunks)

	if len(chunks) == 1 {
		return "\n\nNotation : " + chunks[0]
	}
	return "\n\nNotation : " + strings.Join(chunks[:len(chunks)-1], " ; ") + " et " + chunks[len(chunks)-1]
}

// SyntaxHint refers to a hint for a peculiar expression,
// among :
//
//	inf : \infty s'écrit inf
//	-inf: -\infty s'écrit -inf
//	a/b (simple): \frac{3}{5} s'écrit 3 / 5
//	a/b (complex): \frac{2x+1}{-3x+7} s'écrit (2x + 1) / (-3x + 7)
//	x^n (n >= 2) : x^2 s'écrit x^2
//	sqrt(3) : \sqrt{3} s'écrit sqrt(3)
type SyntaxHint uint8

const (
	_ SyntaxHint = iota
	shInf
	shFractionSimple
	shFractionComplex
	shPower
	shSqrt
)

// text returns a string using interpolation syntax $ $ for LaTeX
func (hint SyntaxHint) text() string {
	switch hint {
	case shInf:
		return `$\infty$ s'écrit inf`
	case shFractionSimple:
		return `$\frac{3}{5}$ s'écrit 3 / 5`
	case shFractionComplex:
		return `$\frac{2x+1}{-3x+7}$ s'écrit (2x + 1) / (-3x + 7)`
	case shPower:
		return `$x^2$ s'écrit x^2`
	case shSqrt:
		return `$\sqrt{3}$ s'écrit sqrt(3)`
	default:
		panic("exhaustive SyntaxHint")
	}
}

// SyntaxHints returns a set of hints needed to type [e]
func (e *Expr) SyntaxHints() SyntaxHints {
	if e == nil {
		return nil
	}

	hints := make(SyntaxHints)
	// study the root
	switch a := e.atom.(type) {
	case Number:
		if math.IsInf(float64(a), 0) {
			hints[shInf] = true
		}
	case operator:
		if a == pow {
			hints[shPower] = true
		} else if a == div {
			_, isNumNumber := e.left.isConstantTerm()
			_, isDemNumber := e.right.isConstantTerm()
			if isNumNumber && isDemNumber {
				hints[shFractionSimple] = true
			} else {
				hints[shFractionComplex] = true
			}
		}
	case function:
		if a == sqrtFn {
			hints[shSqrt] = true
		}
	}

	// recurse on children
	hints.Append(e.left.SyntaxHints())
	hints.Append(e.right.SyntaxHints())
	return hints
}
