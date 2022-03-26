package expression

import "fmt"

// Intrinsic maps some variables to expression,
// serving as shortcut for complex conditions
type Intrinsic interface {
	MergeTo(params RandomParameters)
}

// BuildParams is a convenience shortcut to build the
// parameters associated with the given intrinsic.
func BuildParams(it Intrinsic) RandomParameters {
	out := make(RandomParameters)
	it.MergeTo(out)
	return out
}

// PythagorianTriplet are 3 positive integers (a,b,c) such that a^2 + b^2 = c^2
type PythagorianTriplet struct {
	A, B, C Variable
	// SeedStart and SeedEnd define the magnitude of a, b, c such that
	// 2*SeedStart^2 <= a <= 2*SeedEnd^2
	SeedStart, SeedEnd int
}

func (pt PythagorianTriplet) MergeTo(params RandomParameters) {
	p := params.addAnonymousParam(mustParseE(fmt.Sprintf("2 * randInt(%d;%d)", pt.SeedStart, pt.SeedEnd)))
	q := params.addAnonymousParam(mustParseE(fmt.Sprintf("randInt(%d;%d)", pt.SeedStart, pt.SeedEnd)))
	a := mustParseE(fmt.Sprintf("%s*%s", p, q))                  // p * q
	c := mustParseE(fmt.Sprintf("(%s %s^2 + %s)  / 2", p, q, p)) // (p q^2 + p)  / 2
	b := mustParseE(fmt.Sprintf("%s-%s", pt.C, p))               // c - p

	params[pt.A] = a
	params[pt.B] = b
	params[pt.C] = c
}

// PolynomialCoeffs are the coefficient of
// P = 3/4 X^4 +bX^3 + cX^2 + dX, such that the roots of
// P' are integers in the given range, which must be at least with length 2.
type PolynomialCoeffs struct {
	B, C, D              Variable
	X1, X2, X3           Variable
	RootsStart, RootsEnd int
}

func (qp PolynomialCoeffs) MergeTo(params RandomParameters) {
	width := qp.RootsEnd - qp.RootsStart
	// ensure no solutions are too close
	// |--_|_-_|_--|
	x1 := mustParseE(fmt.Sprintf("randInt(%d;%d)", qp.RootsStart, qp.RootsStart+2*width/9))
	x2 := mustParseE(fmt.Sprintf("randInt(%d;%d)", qp.RootsStart+4*width/9, qp.RootsStart+5*width/9))
	x3 := mustParseE(fmt.Sprintf("randInt(%d;%d)", qp.RootsStart+7*width/9, qp.RootsEnd))

	// we solve
	// P' = 4a X^3 + 3bX^2 + 2cX + d = 4a(X - x1)(X - x2)(X - x3)
	// 3b = 4a * ( -x1 - x2 - x3  )
	// 2c = 4a * ( x1 x2 + x1 x3 + x2 x3 )
	// d = 4a * -(x1 x2 x3 )

	// for now we simply choose a = 3/4 and e = 0
	b := mustParseE(fmt.Sprintf("-(%s + %s + %s)", qp.X1, qp.X2, qp.X3))
	c := mustParseE(fmt.Sprintf("3 * (%s * %s + %s * %s + %s * %s) / 2", qp.X1, qp.X2, qp.X1, qp.X3, qp.X2, qp.X3))
	d := mustParseE(fmt.Sprintf("-3 * %s * %s * %s", qp.X1, qp.X2, qp.X3))

	params[qp.X1] = x1
	params[qp.X2] = x2
	params[qp.X3] = x3
	params[qp.B] = b
	params[qp.C] = c
	params[qp.D] = d
}
