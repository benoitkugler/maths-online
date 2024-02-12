package expression

import (
	"testing"
)

func TestPythagorianTriplet_instantiateTo(t *testing.T) {
	tests := []struct {
		Bound int
	}{
		{10},
		{40},
	}
	for _, tt := range tests {
		pt := pythagorianTriplet{
			a: NewVar('a'), b: NewVar('b'), c: NewVar('c'),
			bound: tt.Bound,
		}
		out := RandomParameters{specials: []intrinsic{pt}}

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

func TestOrthogonalProjection_mergeTo(t *testing.T) {
	op := orthogonalProjection{
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
			RandomParameters{defs: map[Variable]*Expr{
				op.Ax: NewNb(0), op.Ay: NewNb(2),
				op.Bx: NewNb(-1), op.By: NewNb(0),
				op.Cx: NewNb(4), op.Cy: NewNb(0),
			}},
			0, 0,
		},
		{
			RandomParameters{defs: map[Variable]*Expr{
				op.Ax: NewNb(-1), op.Ay: NewNb(1),
				op.Bx: NewNb(-1), op.By: NewNb(-1),
				op.Cx: NewNb(1), op.Cy: NewNb(1),
			}},
			0, 0,
		},
	}
	for _, tt := range tests {
		op.mergeTo(&tt.args)
		v, err := tt.args.Instantiate()
		if err != nil {
			t.Fatal(err)
		}
		if v[op.Hx].mustEvaluate(nil) != tt.expectedX || v[op.Hy].mustEvaluate(nil) != tt.expectedY {
			t.Fatal()
		}
	}
}
