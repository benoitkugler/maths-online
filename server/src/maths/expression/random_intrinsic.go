package expression

import "fmt"

// Intrinsic maps some variables to expression,
// serving as shortcut for complex conditions
type Intrinsic interface {
	MergeTo(params RandomParameters)
}

// PythagorianTriplet are 3 positive integers (a,b,c) such that a^2 + b^2 = c^2
type PythagorianTriplet struct {
	A, B, C Variable
	// SeedStart and SeedEnd define the magnitude of a, b, c such that
	// 2*SeedStart^2 <= a <= 2*SeedEnd^2
	SeedStart, SeedEnd int
}

func (pt PythagorianTriplet) MergeTo(params RandomParameters) {
	expr, _, _ := Parse(fmt.Sprintf("2 * randInt(%d;%d)", pt.SeedStart, pt.SeedEnd))
	p := params.addAnonymousParam(expr)
	q := params.addAnonymousParam(&Expression{atom: random{start: pt.SeedStart, end: pt.SeedEnd}})
	a := &Expression{atom: mult, left: &Expression{atom: p}, right: &Expression{atom: q}} // p * q
	c, _, _ := Parse(fmt.Sprintf("(%s %s^2 + %s)  / 2", p, q, p))                         // (p q^2 + p)  / 2
	b := &Expression{atom: minus, left: c, right: &Expression{atom: p}}                   // c - p

	params[pt.A] = a
	params[pt.B] = b
	params[pt.C] = c
}
