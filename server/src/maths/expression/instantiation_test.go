package expression

import (
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
	}
	for _, tt := range tests {
		rv := make(RandomParameters)
		for v, e := range tt.rv {
			rv[v] = mustParse(t, e)
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
		rv := RandomParameters{NewVar('P'): mustParse(t, "randChoice(A;B;C)")}
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
	params := RandomParameters{
		NewVar('a'): mustParse(t, "randChoice(a;b)"),
		NewVar('c'): mustParse(t, "choiceFrom(c;d; randint(1;2))"),
	}
	err := params.Validate()
	tu.AssertNoErr(t, err)

	// test that non trivial cycle are not accepted
	params = RandomParameters{
		NewVar('a'): mustParse(t, "randChoice(a;b)"),
		NewVar('b'): mustParse(t, "randChoice(a;b)"),
	}
	err = params.Validate()
	tu.Assert(t, err != nil)
}
