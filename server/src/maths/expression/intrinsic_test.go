package expression

import (
	"testing"
)

func TestPythagorianTriplet_mergeTo(t *testing.T) {
	tests := []struct {
		Bound int
	}{
		{10},
		{40},
	}
	for _, tt := range tests {
		pt := PythagorianTriplet{
			A: NewVariable('a'), B: NewVariable('b'), C: NewVariable('c'),
			Bound: tt.Bound,
		}
		out := buildParams(pt)

		for range [100]int{} {
			vr, err := out.Instantiate()
			if err != nil {
				t.Fatal(err)
			}
			a, b, c := vr[NewVariable('a')], vr[NewVariable('b')], vr[NewVariable('c')]
			if _, ok := isInt(a); !ok {
				t.Fatal()
			}
			if _, ok := isInt(b); !ok {
				t.Fatal()
			}
			if _, ok := isInt(c); !ok {
				t.Fatal()
			}
			if a*a+b*b != c*c {
				t.Fatal()
			}
		}
	}
}

func TestQuadraticPolynomialCoeffs_MergeTo(t *testing.T) {
	tests := []struct {
		RootsStart int
		RootsEnd   int
	}{
		{-3, 3},
	}
	for _, tt := range tests {
		qp := PolynomialCoeffs{
			B:          NewVariable('b'),
			C:          NewVariable('c'),
			D:          NewVariable('d'),
			X1:         NewVariable('u'),
			X2:         NewVariable('v'),
			X3:         NewVariable('w'),
			RootsStart: tt.RootsStart,
			RootsEnd:   tt.RootsEnd,
		}
		out := buildParams(qp)

		for range [10]int{} {
			vs, err := out.Instantiate()
			if err != nil {
				t.Fatal(err)
			}
			expr := mustParse(t, "4*((3/4)X^4 + bX^3 + cX^2 + dX)")
			expr.Substitute(vs)

			if v := expr.MustEvaluate(Variables{NewVariable('X'): 0}); v != 0 {
				t.Fatal(v)
			}

			x1, x2, x3 := vs[NewVariable('u')], vs[NewVariable('v')], vs[NewVariable('w')]

			derivative := mustParse(t, "3X^3 + 3bX^2 + 2cX + d")
			derivative.Substitute(vs)

			if v := derivative.MustEvaluate(Variables{NewVariable('X'): x1}); v != 0 {
				t.Fatalf("expected df(%v) = 0, got %v", x1, v)
			}
			if v := derivative.MustEvaluate(Variables{NewVariable('X'): x2}); v != 0 {
				t.Fatalf("expected df(%v) = 0, got %v", x2, v)
			}
			if v := derivative.MustEvaluate(Variables{NewVariable('X'): x3}); v != 0 {
				t.Fatalf("expected df(%v) = 0, got %v", x3, v)
			}

			// fmt.Println(expr.MustEvaluate(Variables{NewVariable('X'): x1}))
			// fmt.Println(expr.MustEvaluate(Variables{NewVariable('X'): x2}))
			// fmt.Println(expr.MustEvaluate(Variables{NewVariable('X'): x3}))

			width := qp.RootsEnd - qp.RootsStart
			minDist := float64(2*width)/9 - 1 // -1 to account for rouding error
			if x1 > x2 || x2 > x3 {
				t.Fatal(x1, x2, x3)
			}
			if x2-x1 < minDist || x3-x2 < minDist {
				t.Fatal(width, minDist, x1, x2, x3)
			}
		}
	}
}

func TestOrthogonalProjection_MergeTo(t *testing.T) {
	op := OrthogonalProjection{
		Ax: Variable{Name: 'x', Indice: "A"},
		Ay: Variable{Name: 'y', Indice: "A"},
		Bx: Variable{Name: 'x', Indice: "B"},
		By: Variable{Name: 'y', Indice: "B"},
		Cx: Variable{Name: 'x', Indice: "C"},
		Cy: Variable{Name: 'y', Indice: "C"},
		Hx: Variable{Name: 'x', Indice: "H"},
		Hy: Variable{Name: 'y', Indice: "H"},
	}
	tests := []struct {
		args                 RandomParameters
		expectedX, expectedY float64
	}{
		{
			RandomParameters{
				op.Ax: NewNumber(0), op.Ay: NewNumber(2),
				op.Bx: NewNumber(-1), op.By: NewNumber(0),
				op.Cx: NewNumber(4), op.Cy: NewNumber(0),
			},
			0, 0,
		},
		{
			RandomParameters{
				op.Ax: NewNumber(-1), op.Ay: NewNumber(1),
				op.Bx: NewNumber(-1), op.By: NewNumber(-1),
				op.Cx: NewNumber(1), op.Cy: NewNumber(1),
			},
			0, 0,
		},
	}
	for _, tt := range tests {
		op.MergeTo(tt.args)
		v, err := tt.args.Instantiate()
		if err != nil {
			t.Fatal(err)
		}
		if v[op.Hx] != tt.expectedX || v[op.Hy] != tt.expectedY {
			t.Fatal()
		}
	}
}
