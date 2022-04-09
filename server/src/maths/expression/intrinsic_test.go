package expression

import (
	"fmt"
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
			A: 'a', B: 'b', C: 'c',
			Bound: tt.Bound,
		}
		out := BuildParams(pt)

		for range [100]int{} {
			vr, err := out.Instantiate()
			if err != nil {
				t.Fatal(err)
			}
			a, b, c := vr['a'], vr['b'], vr['c']
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
			B:          'b',
			C:          'c',
			D:          'd',
			X1:         'u',
			X2:         'v',
			X3:         'w',
			RootsStart: tt.RootsStart,
			RootsEnd:   tt.RootsEnd,
		}
		out := BuildParams(qp)

		for range [10]int{} {
			vs, err := out.Instantiate()
			if err != nil {
				t.Fatal(err)
			}
			expr := mustParse(t, "4*((3/4)X^4 + bX^3 + cX^2 + dX)")
			expr.Substitute(vs)

			if v := expr.MustEvaluate(Variables{'X': 0}); v != 0 {
				t.Fatal(v)
			}

			x1, x2, x3 := vs['u'], vs['v'], vs['w']

			derivative := mustParse(t, "3X^3 + 3bX^2 + 2cX + d")
			derivative.Substitute(vs)

			if v := derivative.MustEvaluate(Variables{'X': x1}); v != 0 {
				t.Fatalf("expected df(%v) = 0, got %v", x1, v)
			}
			if v := derivative.MustEvaluate(Variables{'X': x2}); v != 0 {
				t.Fatalf("expected df(%v) = 0, got %v", x2, v)
			}
			if v := derivative.MustEvaluate(Variables{'X': x3}); v != 0 {
				t.Fatalf("expected df(%v) = 0, got %v", x3, v)
			}

			fmt.Println(expr.MustEvaluate(Variables{'X': x1}))
			fmt.Println(expr.MustEvaluate(Variables{'X': x2}))
			fmt.Println(expr.MustEvaluate(Variables{'X': x3}))

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
