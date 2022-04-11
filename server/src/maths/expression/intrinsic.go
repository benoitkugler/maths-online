package expression

import (
	"fmt"
)

type ErrDuplicateParameter struct {
	Duplicate Variable
}

func (err ErrDuplicateParameter) Error() string {
	return fmt.Sprintf("Paramètre %s défini deux fois", err.Duplicate)
}

func checkDuplicates(params RandomParameters, vars ...Variable) error {
	for _, v := range vars {
		if _, has := params[v]; has {
			return ErrDuplicateParameter{Duplicate: v}
		}
	}
	return nil
}

// Intrinsic maps some variables to expression,
// serving as shortcut for complex conditions
type Intrinsic interface {
	// MergeTo returns an `ErrDuplicateParameter` error if parameters are already defined
	MergeTo(params RandomParameters) error
}

// buildParams is a convenience shortcut to build the
// parameters associated with the given intrinsic.
func buildParams(it Intrinsic) RandomParameters {
	out := make(RandomParameters)
	_ = it.MergeTo(out)
	return out
}

// PythagorianTriplet are 3 positive integers (a,b,c) such that a^2 + b^2 = c^2
type PythagorianTriplet struct {
	A, B, C Variable
	// Bound roughly defines the magnitude of a, b, c, such that
	// a <= 2*Bound^2.
	// It must be >= 2
	Bound int
}

func (pt PythagorianTriplet) MergeTo(params RandomParameters) error {
	if err := checkDuplicates(params, pt.A, pt.B, pt.C); err != nil {
		return err
	}

	const seedStart = 1
	p := params.addAnonymousParam(MustParse(fmt.Sprintf("2 * randInt(%d;%d)", seedStart, pt.Bound)))
	// q = 1 yield b = 0, avoid this edge case
	q := params.addAnonymousParam(MustParse(fmt.Sprintf("randInt(%d;%d)", seedStart+1, pt.Bound)))
	a := MustParse(fmt.Sprintf("%s*%s", p, q))                  // p * q
	c := MustParse(fmt.Sprintf("(%s %s^2 + %s)  / 2", p, q, p)) // (p q^2 + p)  / 2
	b := MustParse(fmt.Sprintf("%s-%s", pt.C, p))               // c - p

	params[pt.A] = a
	params[pt.B] = b
	params[pt.C] = c

	return nil
}

// PolynomialCoeffs are the coefficient of
// P = 3/4 X^4 +bX^3 + cX^2 + dX, such that the roots of
// P' are integers in the given range, which must be at least with length 2.
type PolynomialCoeffs struct {
	B, C, D              Variable
	X1, X2, X3           Variable
	RootsStart, RootsEnd int
}

func (qp PolynomialCoeffs) MergeTo(params RandomParameters) error {
	if err := checkDuplicates(params, qp.B, qp.C, qp.D, qp.X1, qp.X2, qp.X3); err != nil {
		return err
	}

	width := qp.RootsEnd - qp.RootsStart
	// ensure no solutions are too close
	// |--_|_-_|_--|
	x1 := MustParse(fmt.Sprintf("randInt(%d;%d)", qp.RootsStart, qp.RootsStart+2*width/9))
	x2 := MustParse(fmt.Sprintf("randInt(%d;%d)", qp.RootsStart+4*width/9, qp.RootsStart+5*width/9))
	x3 := MustParse(fmt.Sprintf("randInt(%d;%d)", qp.RootsStart+7*width/9, qp.RootsEnd))

	// we solve
	// P' = 4a X^3 + 3bX^2 + 2cX + d = 4a(X - x1)(X - x2)(X - x3)
	// 3b = 4a * ( -x1 - x2 - x3  )
	// 2c = 4a * ( x1 x2 + x1 x3 + x2 x3 )
	// d = 4a * -(x1 x2 x3 )

	// for now we simply choose a = 3/4 and e = 0
	b := MustParse(fmt.Sprintf("-(%s + %s + %s)", qp.X1, qp.X2, qp.X3))
	c := MustParse(fmt.Sprintf("3 * (%s * %s + %s * %s + %s * %s) / 2", qp.X1, qp.X2, qp.X1, qp.X3, qp.X2, qp.X3))
	d := MustParse(fmt.Sprintf("-3 * %s * %s * %s", qp.X1, qp.X2, qp.X3))

	params[qp.X1] = x1
	params[qp.X2] = x2
	params[qp.X3] = x3
	params[qp.B] = b
	params[qp.C] = c
	params[qp.D] = d

	return nil
}

// OrthogonalProjection compute the coordinates of H, the orthogonal
// projection of A on (BC).
type OrthogonalProjection struct {
	Ax, Ay, Bx, By, Cx, Cy Variable // arguments

	Hx, Hy Variable // output
}

func (op OrthogonalProjection) MergeTo(params RandomParameters) error {
	if err := checkDuplicates(params, op.Hx, op.Hy); err != nil {
		return err
	}

	// BC
	u := params.addAnonymousParam(MustParse(fmt.Sprintf("%s - %s", op.Cx, op.Bx))) // Cx - Bx
	v := params.addAnonymousParam(MustParse(fmt.Sprintf("%s - %s", op.Cy, op.By))) // Cy - By

	// det(BA,BC)
	d := params.addAnonymousParam(MustParse(fmt.Sprintf("(%s - %s)%s - (%s - %s)%s", op.Ax, op.Bx, v, op.Ay, op.By, u)))

	// solve for AH = (x, y)
	// xu + yv = 0
	// xv - yu = -d
	norm := params.addAnonymousParam(MustParse(fmt.Sprintf("%s^2 + %s^2", u, v)))

	params[op.Hx] = MustParse(fmt.Sprintf("(-%s * %s / %s) + %s", d, v, norm, op.Ax))
	params[op.Hy] = MustParse(fmt.Sprintf("(%s * %s / %s) + %s", d, u, norm, op.Ay))

	return nil
}
