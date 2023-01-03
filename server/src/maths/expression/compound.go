package expression

import (
	"fmt"
	"strings"
)

// Compound is a sum type for a complex math object, built on [Expr]s,
// such as sets, intervals, or vectors.
// For compatibility, [*Expr] are also implementing [Compound].
// Note that nested compound objects, such as sets of sets, are not supported.
type Compound interface {
	isCompound()

	// String returns a human readable form of the expression.
	// The expression is prettified, meaning the structure of the
	// returned expression may slightly differ, but is garanteed
	// to be mathematically equivalent.
	// See `AsLaTeX` for a better display format.
	String() string

	// AsLaTeX returns a valid LaTeX code displaying the expression.
	AsLaTeX() string

	// Substitute replaces variables contained in `vars`, updating the math object in place.
	Substitute(vars Vars)
}

// ParseCompound accepts an expression describing a complex math object,
// extending [Parse].
// If invalid, [ErrInvalidExpr] is returned.
func ParseCompound(expr string) (out Compound, err error) {
	bytes := []byte(expr)
	pr := newParser(bytes)
	// look for a starting delimiter
	switch tok := pr.tk.Peek().data; tok {
	case openBracket, closeBracket: // Interval
		out, err = pr.parseInterval()
	case openCurly: // Set
		out, err = pr.parseSet()
	case openPar:
		// either an expression or a vector
		// to distinguish, try to parse as an expression,
		// if it fails, try as vector
		var parsedExpr *Expr
		parsedExpr, _, err = parseBytes(bytes)
		if err == nil { // simple expression
			out = parsedExpr
		} else { // it's a vector
			out, err = pr.parseVector()
		}
	default: // simple expression
		out, _, err = parseBytes(bytes)
	}

	if err != nil {
		errV := err.(ErrInvalidExpr)
		errV.Input = expr
		return nil, errV
	}

	return out, nil
}

func (Vector) isCompound()   {}
func (Set) isCompound()      {}
func (Interval) isCompound() {}
func (Expr) isCompound()     {}

// Vector is a n-uplet of expressions, such as (1;2;3;4;5)
type Vector []*Expr

func (pr *parser) parseVector() (Vector, error) {
	_ = pr.tk.Next() // consume the start
	args, _, err := pr.parseExpressionList(func(s symbol) bool { return s == closePar })
	if err != nil {
		return nil, err
	}
	return args, nil
}

func (vec Vector) String() string {
	chunks := make([]string, len(vec))
	for i, expr := range vec {
		chunks[i] = expr.String()
	}
	return fmt.Sprintf("( %s )", strings.Join(chunks, " ; "))
}

func (vec Vector) AsLaTeX() string {
	chunks := make([]string, len(vec))
	for i, expr := range vec {
		chunks[i] = expr.AsLaTeX()
	}
	return fmt.Sprintf("\\left( %s \\right)", strings.Join(chunks, " ; "))
}

func (vec Vector) Substitute(vars Vars) {
	for _, expr := range vec {
		expr.Substitute(vars)
	}
}

func areVectorsEquivalent(vec1, vec2 Vector, level ComparisonLevel) bool {
	if len(vec1) != len(vec2) {
		return false
	}
	for i, e1 := range vec1 {
		e2 := vec2[i]
		if !AreExpressionsEquivalent(e1, e2, level) {
			return false
		}
	}
	return true
}

// Set is a set of expression, such as {1; x^2; x^4}
type Set []*Expr

func (pr *parser) parseSet() (Set, error) {
	_ = pr.tk.Next() // consume the start
	args, _, err := pr.parseExpressionList(func(s symbol) bool { return s == closeCurly })
	if err != nil {
		return nil, err
	}
	return args, nil
}

func (set Set) String() string {
	chunks := make([]string, len(set))
	for i, expr := range set {
		chunks[i] = expr.String()
	}
	return fmt.Sprintf("{ %s }", strings.Join(chunks, " ; "))
}

func (set Set) AsLaTeX() string {
	chunks := make([]string, len(set))
	for i, expr := range set {
		chunks[i] = expr.AsLaTeX()
	}
	return fmt.Sprintf("\\left\\{ %s \\right\\}", strings.Join(chunks, " ; "))
}

func (set Set) Substitute(vars Vars) {
	for _, expr := range set {
		expr.Substitute(vars)
	}
}

// returns true if set1 is included in set2
func isSetIncluded(set1, set2 Set, level ComparisonLevel) bool {
	for _, e1 := range set1 {
		// does e1 belong to set2 ?
		inSet2 := false
		for _, e2 := range set2 {
			if AreExpressionsEquivalent(e1, e2, level) { // e1 is in set2
				inSet2 = true
				break
			}
		}
		if !inSet2 {
			return false
		}
	}
	return true
}

// double inclusion
func areSetsEquivalent(set1, set2 Set, level ComparisonLevel) bool {
	// we can't rely on length, since they may be some repetitions,
	// just check manually the reverse inclusion
	return isSetIncluded(set1, set2, level) && isSetIncluded(set2, set1, level)
}

// Interval is a math interval (range) such as [2;4]
type Interval struct {
	Left, Right *Expr
	// Open means reject the boundary
	LeftOpen, RightOpen bool
}

func (pr *parser) parseInterval() (Interval, error) {
	tok := pr.tk.Next() // consume the start
	out := Interval{
		LeftOpen: tok.data == closeBracket, // ]
	}
	args, closer, err := pr.parseExpressionList(func(s symbol) bool { return s == openBracket || s == closeBracket })
	if err != nil {
		return Interval{}, err
	}
	if len(args) != 2 {
		return Interval{}, ErrInvalidExpr{
			Reason: "Deux expressions sont attendues pour un intervalle.",
			Pos:    0,
		}
	}
	out.Left, out.Right = args[0], args[1]
	out.RightOpen = closer == openBracket // [
	return out, nil
}

func (inter Interval) String() string {
	leftSymbol := "["
	if inter.LeftOpen {
		leftSymbol = "]"
	}
	rightSymbol := "]"
	if inter.RightOpen {
		rightSymbol = "["
	}
	return fmt.Sprintf("%s %s ; %s %s", leftSymbol, inter.Left.String(), inter.Right.String(), rightSymbol)
}

func (inter Interval) AsLaTeX() string {
	leftSymbol := "["
	if inter.LeftOpen {
		leftSymbol = "]"
	}
	rightSymbol := "]"
	if inter.RightOpen {
		rightSymbol = "["
	}
	return fmt.Sprintf("\\left%s %s ; %s \\right%s", leftSymbol, inter.Left.AsLaTeX(), inter.Right.AsLaTeX(), rightSymbol)
}

func (inter Interval) Substitute(vars Vars) {
	inter.Left.Substitute(vars)
	inter.Right.Substitute(vars)
}

func areIntervalsEquivalent(inter1, inter2 Interval, level ComparisonLevel) bool {
	if inter1.LeftOpen != inter2.LeftOpen || inter1.RightOpen != inter2.RightOpen {
		return false
	}
	return AreExpressionsEquivalent(inter1.Left, inter2.Left, level) && AreExpressionsEquivalent(inter1.Right, inter2.Right, level)
}

// AreCompoundsEquivalent compares the two objects using
// mathematical knowledge.
// Note that objects with different types are never equivalent,
// so that (2;) != 2 != {2}
//
// See AreExpressionsEquivalent for the meaning of [level].
func AreCompoundsEquivalent(e1, e2 Compound, level ComparisonLevel) bool {
	switch e1 := e1.(type) {
	case Vector:
		e2, ok := e2.(Vector)
		if !ok {
			return false
		}
		return areVectorsEquivalent(e1, e2, level)
	case Set:
		e2, ok := e2.(Set)
		if !ok {
			return false
		}
		return areSetsEquivalent(e1, e2, level)
	case Interval:
		e2, ok := e2.(Interval)
		if !ok {
			return false
		}
		return areIntervalsEquivalent(e1, e2, level)
	case *Expr:
		e2, ok := e2.(*Expr)
		if !ok {
			return false
		}
		return AreExpressionsEquivalent(e1, e2, level)
	default:
		panic(exhaustiveCompoundSwitch)
	}
}
