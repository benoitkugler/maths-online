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
			A: NewVar('a'), B: NewVar('b'), C: NewVar('c'),
			Bound: tt.Bound,
		}
		out := buildParams(pt)

		for range [100]int{} {
			vr, err := out.Instantiate()
			if err != nil {
				t.Fatal(err)
			}
			a, b, c := vr[NewVar('a')].mustEvaluate(nil), vr[NewVar('b')].mustEvaluate(nil), vr[NewVar('c')].mustEvaluate(nil)
			if _, ok := IsInt(a); !ok {
				t.Fatal()
			}
			if _, ok := IsInt(b); !ok {
				t.Fatal()
			}
			if _, ok := IsInt(c); !ok {
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
			B:          NewVar('b'),
			C:          NewVar('c'),
			D:          NewVar('d'),
			X1:         NewVar('u'),
			X2:         NewVar('v'),
			X3:         NewVar('w'),
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

			if v := expr.mustEvaluate(Vars{NewVar('X'): NewNb(0)}); v != 0 {
				t.Fatal(v)
			}

			x1, x2, x3 := vs[NewVar('u')].mustEvaluate(nil), vs[NewVar('v')].mustEvaluate(nil), vs[NewVar('w')].mustEvaluate(nil)

			derivative := mustParse(t, "3X^3 + 3bX^2 + 2cX + d")
			derivative.Substitute(vs)

			if v := derivative.mustEvaluate(Vars{NewVar('X'): NewNb(x1)}); v != 0 {
				t.Fatalf("expected df(%v) = 0, got %v", x1, v)
			}
			if v := derivative.mustEvaluate(Vars{NewVar('X'): NewNb(x2)}); v != 0 {
				t.Fatalf("expected df(%v) = 0, got %v", x2, v)
			}
			if v := derivative.mustEvaluate(Vars{NewVar('X'): NewNb(x3)}); v != 0 {
				t.Fatalf("expected df(%v) = 0, got %v", x3, v)
			}

			// fmt.Println(expr.mustEvaluate(Variables{NewVariable('X'): x1}))
			// fmt.Println(expr.mustEvaluate(Variables{NewVariable('X'): x2}))
			// fmt.Println(expr.mustEvaluate(Variables{NewVariable('X'): x3}))

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
				op.Ax: NewNb(0), op.Ay: NewNb(2),
				op.Bx: NewNb(-1), op.By: NewNb(0),
				op.Cx: NewNb(4), op.Cy: NewNb(0),
			},
			0, 0,
		},
		{
			RandomParameters{
				op.Ax: NewNb(-1), op.Ay: NewNb(1),
				op.Bx: NewNb(-1), op.By: NewNb(-1),
				op.Cx: NewNb(1), op.Cy: NewNb(1),
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
		if v[op.Hx].mustEvaluate(nil) != tt.expectedX || v[op.Hy].mustEvaluate(nil) != tt.expectedY {
			t.Fatal()
		}
	}
}
