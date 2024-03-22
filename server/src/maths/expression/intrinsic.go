package expression

import (
	"fmt"
	"math/rand"
)

type ErrDuplicateParameter struct {
	Duplicate Variable
}

func (err ErrDuplicateParameter) Error() string {
	return fmt.Sprintf("Paramètre %s défini deux fois", err.Duplicate)
}

func checkDuplicates(params map[Variable]*Expr, vars ...Variable) error {
	for _, v := range vars {
		if _, has := params[v]; has {
			return ErrDuplicateParameter{Duplicate: v}
		}
	}
	return nil
}

type intrinsic interface {
	instantiateTo(target Vars) error
	isDef(v Variable) bool // returns true if v is an OUTPUT variable
	vars() []Variable
}

// pythagorianTriplet are 3 positive integers (a,b,c) such that a^2 + b^2 = c^2
type pythagorianTriplet struct {
	a, b, c Variable
	// bound roughly defines the magnitude of a, b, c, such that
	// a <= 2*bound^2.
	// It must be >= 2
	bound int
}

func (pt pythagorianTriplet) isDef(v Variable) bool {
	return v == pt.a || v == pt.b || v == pt.c
}

func (pt pythagorianTriplet) vars() []Variable {
	return []Variable{pt.a, pt.b, pt.c}
}

func (pt pythagorianTriplet) instantiateTo(target Vars) error {
	if err := checkDuplicates(target, pt.a, pt.b, pt.c); err != nil {
		return err
	}

	const seedStart = 1

	p := 2 * randomInt(seedStart, pt.bound)
	// q = 1 yield b = 0, avoid this edge case
	q := randomInt(seedStart+1, pt.bound)
	a := p * q
	c := (p*q*q + p) / 2
	b := c - p

	target[pt.a] = rat{a, 1}.toExpr()
	target[pt.b] = rat{b, 1}.toExpr()
	target[pt.c] = rat{c, 1}.toExpr()

	return nil
}

// orthogonalProjection compute the coordinates of H, the orthogonal
// projection of A on (BC).
type orthogonalProjection struct {
	Ax, Ay, Bx, By, Cx, Cy Variable // arguments

	Hx, Hy Variable // output
}

func (op orthogonalProjection) mergeTo(params *RandomParameters) error {
	if err := checkDuplicates(params.defs, op.Hx, op.Hy); err != nil {
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

	params.defs[op.Hx] = MustParse(fmt.Sprintf("(-%s * %s / %s) + %s", d, v, norm, op.Ax))
	params.defs[op.Hy] = MustParse(fmt.Sprintf("(%s * %s / %s) + %s", d, u, norm, op.Ay))

	return nil
}

// numberPair generates random number suitable
// to be added (or multiplied), according to a difficulty level (from 1 to 5)
// Example :
// a,b = number_pair_add(1) # difficulty one
type numberPair struct {
	a, b             Variable // output
	difficulty       uint8    // in [1, 5]
	isMultiplicative bool
}

func (np numberPair) instantiateTo(target Vars) error {
	if err := checkDuplicates(target, np.a, np.b); err != nil {
		return err
	}

	ranges := addPairTable[np.difficulty]
	if np.isMultiplicative {
		ranges = multPairTable[np.difficulty]
	}
	selectedRange := ranges[rand.Intn(len(ranges))]
	a := randomInt(selectedRange.a[0], selectedRange.a[1])
	b := randomInt(selectedRange.b[0], selectedRange.b[1])

	target[np.a] = rat{a, 1}.toExpr()
	target[np.b] = rat{b, 1}.toExpr()

	return nil
}

func (np numberPair) isDef(v Variable) bool {
	return v == np.a || v == np.b
}

func (np numberPair) vars() []Variable {
	return []Variable{np.a, np.b}
}

var (
	addPairTable = [...][]pairRange{
		1: {
			{p{1, 8}, p{1, 1}},
			{p{3, 7}, p{1, 2}},
			{p{4, 6}, p{1, 3}},
			{p{5, 5}, p{1, 4}},
		},
		2: {
			{p{1, 4}, p{1, 5}},
			{p{1, 3}, p{1, 6}},
			{p{1, 2}, p{1, 7}},
			{p{1, 1}, p{1, 8}},
		},
		3: {
			{p{1, 10}, p{1, 10}},
		},
		4: {
			{p{-10, -5}, p{-10, -5}},
			{p{-10, -5}, p{5, 10}},
			{p{5, 10}, p{-10, -5}},
			{p{5, 10}, p{5, 10}},
		},
		5: {
			{p{-20, -10}, p{-20, -10}},
			{p{-20, -10}, p{10, 20}},
			{p{10, 20}, p{-20, -10}},
			{p{10, 20}, p{10, 20}},
		},
	}
	multPairTable = [...][]pairRange{
		1: {
			{p{1, 5}, p{1, 5}},
		},
		2: {
			{p{1, 7}, p{1, 7}},
		},
		3: {
			{p{1, 7}, p{1, 7}},
			{p{-7, -1}, p{-7, -1}},
			{p{1, 7}, p{-7, -1}},
		},
		4: {
			{p{4, 10}, p{4, 10}},
			{p{-10, -4}, p{-10, -4}},
			{p{4, 10}, p{-10, -4}},
		},
		5: {
			{p{5, 12}, p{5, 12}},
			{p{-12, -5}, p{-12, -5}},
			{p{5, 12}, p{-12, -5}},
		},
	}
)

type p = [2]int

type pairRange struct {
	a p
	b p
}
