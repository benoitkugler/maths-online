package expression

import (
	"math"
	"reflect"
	"testing"

	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func TestRandomVariables_Instantiate(t *testing.T) {
	tests := []struct {
		rv      map[Variable]string
		want    Vars
		wantErr bool
	}{
		// cycle
		{
			map[Variable]string{NewVar('a'): "a +1"}, nil, true,
		},
		// cycle
		{
			map[Variable]string{NewVar('a'): "a + b + 1", NewVar('b'): "8"}, nil, true,
		},
		// cycle
		{
			map[Variable]string{NewVar('a'): "b + 1", NewVar('b'): "a+2"}, nil, true,
		},
		// b is a free variable
		{
			map[Variable]string{NewVar('a'): "b + 1"}, Vars{NewVar('a'): mustParse(t, "b+1")}, false,
		},
		// c is a free variable
		{
			map[Variable]string{NewVar('a'): "b + 1", NewVar('b'): "c+1"}, Vars{NewVar('a'): mustParse(t, "c+1+1"), NewVar('b'): mustParse(t, "c+1")}, false,
		},
		{
			map[Variable]string{NewVar('a'): "b + 1", NewVar('b'): " 2 * 3"}, Vars{NewVar('a'): NewNb(7), NewVar('b'): NewNb(6)}, false,
		},
		{
			map[Variable]string{NewVar('a'): "b + 1", NewVar('b'): " c+1", NewVar('c'): "8"}, Vars{NewVar('a'): NewNb(10), NewVar('b'): NewNb(9), NewVar('c'): NewNb(8)}, false,
		},
		{
			map[Variable]string{NewVar('a'): "0*randInt(1;3)"}, Vars{NewVar('a'): NewNb(0)}, false,
		},
		{
			map[Variable]string{NewVar('a'): "randInt(1;1)", NewVar('b'): "2*a"}, Vars{NewVar('a'): NewNb(1), NewVar('b'): NewNb(2)}, false,
		},
		{
			map[Variable]string{NewVar('a'): "randChoice(A)", NewVar('b'): "randInt(1;1)"},
			Vars{NewVar('a'): newVarExpr('A'), NewVar('b'): NewNb(1)},
			false,
		},
		{
			map[Variable]string{NewVar('a'): "choiceFrom(A; 2.1)"},
			nil,
			true,
		},
		{
			map[Variable]string{NewVar('a'): "choiceFrom(A; 2)"},
			nil,
			true,
		},
		{
			map[Variable]string{NewVar('a'): "choiceFrom(A; b)"},
			nil,
			true,
		},
		{
			map[Variable]string{NewVar('a'): "choiceFrom(0.5; 3; b)", NewVar('b'): "2"},
			Vars{NewVar('a'): newNb(3), NewVar('b'): newNb(2)},
			false,
		},
		{
			map[Variable]string{NewVar('a'): "choiceFrom(A; B; b)", NewVar('b'): "1+1"},
			Vars{NewVar('a'): newVarExpr('B'), NewVar('b'): newNb(2)},
			false,
		},
		{
			map[Variable]string{NewVar('a'): "7/3"},
			Vars{NewVar('a'): &Expr{atom: div, left: newNb(7), right: newNb(3)}},
			false,
		},
		{
			map[Variable]string{NewVar('a'): "randInt(b; 3)", NewVar('b'): "2.8"},
			nil,
			true,
		},
		// randChoice with arbitrary expression
		{
			map[Variable]string{NewVar('f'): "randChoice(i)"},
			Vars{NewVar('f'): newVarExpr('i')},
			false,
		},
		// randChoice with arbitrary expression and partial evaluation
		{
			map[Variable]string{NewVar('f'): "randChoice(i * (2 == 4))"},
			Vars{NewVar('f'): mustParse(t, "0")},
			false,
		},
		{
			map[Variable]string{NewVar('f'): "randChoice(i^2)"},
			Vars{NewVar('f'): mustParse(t, "i^2")},
			false,
		},
		// random matrix
		{
			map[Variable]string{NewVar('A'): "randMatrix(1; -1; 1; 10)"},
			nil,
			true,
		},
		{
			map[Variable]string{NewVar('A'): "randMatrix(-1; 2; 1; 10)"},
			nil,
			true,
		},
		{
			map[Variable]string{NewVar('A'): "randMatrix(1; 2; x; 10)"},
			nil,
			true,
		},
		{
			map[Variable]string{NewVar('A'): "randMatrix(1; 2; 20; 10)"},
			nil,
			true,
		},
		{
			map[Variable]string{NewVar('A'): "randMatrix(1.5; 2; 1; 10)"},
			nil,
			true,
		},
		{
			map[Variable]string{NewVar('A'): "randMatrix(1; 2.2; 1; 10)"},
			nil,
			true,
		},
		{
			map[Variable]string{NewVar('A'): "randMatrix(2; 3; 5; 5)"},
			Vars{NewVar('A'): &Expr{atom: matrix{
				{newNb(5), newNb(5), newNb(5)},
				{newNb(5), newNb(5), newNb(5)},
			}}},
			false,
		},
		// substitution in min (or max) functions
		{
			map[Variable]string{NewVar('f'): "min(i^2; j)", NewVar('j'): "4"},
			Vars{NewVar('f'): mustParse(t, "min(i^2; 4 )"), NewVar('j'): mustParse(t, "4")},
			false,
		},
		// evaluation
		{
			map[Variable]string{NewVar('f'): "cos(j)", NewVar('j'): "0"},
			Vars{NewVar('f'): mustParse(t, "1"), NewVar('j'): mustParse(t, "0")},
			false,
		},
		// partial evaluation
		{
			map[Variable]string{NewVar('f'): "(k == 1) * g + (k == 2) * h", NewVar('k'): "2"},
			Vars{NewVar('f'): mustParse(t, "h"), NewVar('k'): mustParse(t, "2")},
			false,
		},
		// substitution in indices
		{
			map[Variable]string{NewVar('f'): "randChoice(F)", NewVar('u'): "f_{1+2}"},
			Vars{NewVar('f'): mustParse(t, "F"), NewVar('u'): mustParse(t, "F_{3}")},
			false,
		},
		// zero length cycle through randChoice or choiceFrom is accepted
		{
			map[Variable]string{
				NewVar('a'): "randChoice(a)",
				NewVar('b'): "choiceFrom(b; 1)",
			},
			Vars{
				NewVar('a'): newVarExpr('a'),
				NewVar('b'): newVarExpr('b'),
			},
			false,
		},
		// matrices
		{map[Variable]string{NewVar('A'): "trace(2)"}, nil, true},
		{map[Variable]string{NewVar('A'): "det([[1;2]])"}, nil, true},
		{map[Variable]string{NewVar('A'): "det(2)"}, nil, true},
		{map[Variable]string{NewVar('A'): "trans(2)"}, nil, true},
		{map[Variable]string{NewVar('A'): "inv(2)"}, nil, true},
		{map[Variable]string{NewVar('A'): "set([[1]];2;2;1)"}, nil, true},
		{map[Variable]string{NewVar('A'): "set(0;2;2;1)"}, nil, true},
		{map[Variable]string{NewVar('A'): "set([[1; 2];[3;4]];2;2;10)"}, Vars{NewVar('A'): mustParse(t, "[[1; 2];[3;10]]")}, false},
		{map[Variable]string{NewVar('A'): "coeff([[1; 2];[3;4]]; 2; 2)"}, Vars{NewVar('A'): newNb(4)}, false},
		{map[Variable]string{
			NewVar('A'): "coeff([[1; 2];[3;4]]; 2; j)",
			NewVar('j'): "1",
		}, Vars{NewVar('A'): newNb(3), NewVar('j'): newNb(1)}, false},
		{
			map[Variable]string{
				NewVar('A'): "[[1;2]; [3;4]]",
				NewVar('b'): "trace(A)",
				NewVar('c'): "det(A)",
				NewVar('d'): "trace(trans(A))",
				NewVar('e'): "det(trans(A))",
			},
			Vars{
				NewVar('A'): mustParse(t, "[[1;2]; [3;4]]"),
				NewVar('b'): newNb(5),
				NewVar('c'): NewNb(-2),
				NewVar('d'): newNb(5),
				NewVar('e'): NewNb(-2),
			},
			false,
		},
		{
			map[Variable]string{
				NewVar('A'): "[[1;2]; [3;x]]",
				NewVar('b'): "trace(A)",
				NewVar('B'): "trans(A)",
			},
			Vars{
				NewVar('A'): mustParse(t, "[[1;2]; [3;x]]"),
				NewVar('b'): mustParse(t, "1+x"),
				NewVar('B'): mustParse(t, "[[1;3]; [2;x]]"),
			},
			false,
		},
		// sequence
		{
			map[Variable]string{
				NewVar('A'): "randint(2;2)",
				NewVar('b'): "prod(k; 1; A; n+k)",
				NewVar('c'): `prod(k; 1; A; n+k; "expand")`,
				NewVar('d'): `prod(k; 1; A; n+k; "expand-eval")`,
			},
			Vars{
				NewVar('A'): mustParse(t, "2"),
				NewVar('b'): mustParse(t, "prod(k; 1; 2; n+k)"),
				NewVar('c'): mustParse(t, "(n+1)(n+2)"),
				NewVar('d'): mustParse(t, "(n+1)(n+2)"),
			},
			false,
		},
	}
	for _, tt := range tests {
		rv := NewRandomParameters()
		for v, e := range tt.rv {
			rv.defs[v] = mustParse(t, e)
		}

		got, err := rv.Instantiate()
		if err != nil {
			err, ok := err.(ErrInvalidRandomParameters)
			if !ok {
				t.Fatal("invalid err type")
			}
			_ = err.Error()
		}
		if (err != nil) != tt.wantErr {
			t.Fatalf("RandomVariables.Instantiate(%s) error = %v, wantErr %v", rv, err, tt.wantErr)
		}
		if err := rv.Validate(); (err != nil) != tt.wantErr {
			t.Fatalf("RandomVariables.Validate(%s) error = %v, wantErr %v", rv, err, tt.wantErr)
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("RandomVariables.Instantiate() = %v, want %v", got, tt.want)
		}
	}
}

func Test_randChoice(t *testing.T) {
	for range [10]int{} {
		rv := RandomParameters{defs: map[Variable]*Expr{NewVar('P'): mustParse(t, "randChoice(A;B;C)")}}
		vars, err := rv.Instantiate()
		if err != nil {
			t.Fatal(err)
		}
		resolved := vars[NewVar('P')]
		asVar, isVar := resolved.atom.(Variable)
		if !isVar {
			t.Fatal(resolved)
		}
		if n := asVar.Name; !(n == 'A' || n == 'B' || n == 'C') {
			t.Fatal(n)
		}
	}
}

func TestCycle(t *testing.T) {
	// test that zero length cycles are properly accepted
	params := RandomParameters{defs: map[Variable]*Expr{
		NewVar('a'): mustParse(t, "randChoice(a;b)"),
		NewVar('c'): mustParse(t, "choiceFrom(c;d; randint(1;2))"),
	}}
	err := params.Validate()
	tu.AssertNoErr(t, err)

	params = RandomParameters{defs: map[Variable]*Expr{
		NewVar('a'): mustParse(t, "a"),
		NewVar('b'): mustParse(t, "b"),
	}}
	err = params.Validate()
	tu.AssertNoErr(t, err)

	// test that non trivial cycle are not accepted
	params = RandomParameters{defs: map[Variable]*Expr{
		NewVar('a'): mustParse(t, "randChoice(a;b)"),
		NewVar('b'): mustParse(t, "randChoice(a;b)"),
	}}
	err = params.Validate()
	tu.Assert(t, err != nil)
}

func TestInstantiateMinMax(t *testing.T) {
	pr := RandomParameters{defs: map[Variable]*Expr{
		NewVar('a'):       mustParse(t, "randChoice(1;2;3)"),
		NewVar('b'):       mustParse(t, "randChoice(8;9;10)"),
		NewVarI('s', "1"): mustParse(t, "min(a; b)"),
		NewVarI('s', "2"): mustParse(t, "max(a; b)"),
	}}
	vars, err := pr.Instantiate()
	tu.AssertNoErr(t, err)

	s1 := vars[NewVarI('s', "1")]
	s2 := vars[NewVarI('s', "2")]

	_, isNumber1 := s1.atom.(Number)
	_, isNumber2 := s2.atom.(Number)
	tu.Assert(t, isNumber1)
	tu.Assert(t, isNumber2)
}

func TestCycleFunc(t *testing.T) {
	params := RandomParameters{defs: map[Variable]*Expr{
		NewVar('a'): mustParse(t, "randChoice(a;b)"),
		NewVar('c'): mustParse(t, "min(a; n)"),
		NewVar('d'): mustParse(t, "binom(a; n)"),
	}}
	err := params.Validate()
	tu.AssertNoErr(t, err)
}

func TestInstantiateWithIntrinsics(t *testing.T) {
	params := NewRandomParameters()
	err := params.ParseIntrinsic("a, b = number_pair_sum(1)")
	tu.AssertNoErr(t, err)
	err = params.ParseVariable(NewVar('A'), "[[a;b]]")
	tu.AssertNoErr(t, err)
	v, err := params.Instantiate()
	tu.AssertNoErr(t, err)
	A := v[NewVar('A')].atom.(matrix)
	a, b := A[0][0], A[0][1]
	an, ok := a.isConstantTerm()
	tu.Assert(t, ok)
	bn, ok := b.isConstantTerm()
	tu.Assert(t, ok)
	_, ok = IsInt(an)
	tu.Assert(t, ok)
	_, ok = IsInt(bn)
	tu.Assert(t, ok)
}

func TestCos(t *testing.T) {
	e := mustParse(t, "sqrt(16+9-2*4*3*cos(27*pi/180))")
	v, err := e.Evaluate(Vars{})
	tu.AssertNoErr(t, err)
	tu.Assert(t, !math.IsNaN(v))

	e = mustParse(t, "sqrt(b*b+c*c-2*b*c*cos(d*pi/180))")
	v, err = e.Evaluate(Vars{NewVar('b'): newNb(4), NewVar('c'): newNb(3), NewVar('d'): newNb(27)})
	tu.AssertNoErr(t, err)
	tu.Assert(t, !math.IsNaN(v))
}

func TestBug353(t *testing.T) {
	// https://github.com/benoitkugler/maths-online/issues/353
	params := NewRandomParameters()
	err := params.ParseVariable(NewVar('r'), "sqrt(a*c)*x+b*sqrt(c)")
	tu.AssertNoErr(t, err)
	err = params.ParseVariable(NewVar('a'), "3")
	tu.AssertNoErr(t, err)
	err = params.ParseVariable(NewVar('b'), "4")
	tu.AssertNoErr(t, err)
	err = params.ParseVariable(NewVar('c'), "7")
	tu.AssertNoErr(t, err)

	vars, err := params.Instantiate()
	tu.AssertNoErr(t, err)

	got := vars[NewVar('r')]
	exp := mustParse(t, "sqrt(21)*x+4*sqrt(7)")

	// 'got' is evaluated, not 'exp'
	tu.Assert(t, !AreExpressionsEquivalent(exp, got, SimpleSubstitutions))
}
