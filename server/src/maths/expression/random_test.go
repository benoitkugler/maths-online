package expression

import (
	"math"
	"reflect"
	"testing"
)

func TestRandomVariables_instantiate(t *testing.T) {
	tests := []struct {
		rv      map[Variable]string
		want    Vars
		wantErr bool
	}{
		{
			map[Variable]string{NewVar('a'): "a +1"}, nil, true,
		},
		{
			map[Variable]string{NewVar('a'): "a + b + 1", NewVar('b'): "8"}, nil, true,
		},
		{
			map[Variable]string{NewVar('a'): "b + 1", NewVar('b'): "a+2"}, nil, true,
		},
		{
			map[Variable]string{NewVar('a'): "b + 1"}, nil, true,
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
			map[Variable]string{NewVar('a'): "randLetter(A)", NewVar('b'): "randInt(1;1)"},
			Vars{NewVar('a'): newVarExpr('A'), NewVar('b'): NewNb(1)},
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
			t.Errorf("RandomVariables.instantiate() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if err := rv.Validate(); (err != nil) != tt.wantErr {
			t.Errorf("RandomVariables.Validate() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("RandomVariables.instantiate() = %v, want %v", got, tt.want)
		}
	}
}

func TestRandLetter(t *testing.T) {
	for range [10]int{} {
		rv := RandomParameters{NewVar('P'): mustParse(t, "randLetter(A;B;C)")}
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

func TestRandomVariables_range(t *testing.T) {
	for range [10]int{} {
		rv := RandomParameters{
			NewVar('a'): mustParse(t, "3*randInt(1; 10)"),
			NewVar('b'): mustParse(t, "-a"),
		}
		values, err := rv.Instantiate()
		if err != nil {
			t.Fatal(err)
		}
		if values[NewVar('a')].mustEvaluate(nil) != -values[NewVar('b')].mustEvaluate(nil) {
			t.Fatal(values)
		}

		if a := values[NewVar('a')].mustEvaluate(nil); a < 3 || a > 30 {
			t.Fatal(a)
		}

		rv = RandomParameters{
			NewVar('a'): mustParse(t, "randInt(1; 10)"),
			NewVar('b'): mustParse(t, "sgn(2*randInt(0;1)-1) * a"),
		}
		values, err = rv.Instantiate()
		if err != nil {
			t.Fatal(err)
		}
		if a := values[NewVar('a')].mustEvaluate(nil); a < 1 || a > 10 {
			t.Fatal(a)
		}
		if a, b := values[NewVar('a')].mustEvaluate(nil), values[NewVar('b')].mustEvaluate(nil); math.Abs(a) != math.Abs(b) {
			t.Fatal(a, b)
		}
	}
}

func Test_sieveOfEratosthenes(t *testing.T) {
	tests := []struct {
		min, max   int
		wantPrimes []int
	}{
		{4, 4, nil},
		{0, 10, []int{2, 3, 5, 7}},
		{0, 11, []int{2, 3, 5, 7, 11}},
		{3, 10, []int{3, 5, 7}},
		{4, 11, []int{5, 7, 11}},
	}
	for _, tt := range tests {
		if gotPrimes := sieveOfEratosthenes(tt.min, tt.max); !reflect.DeepEqual(gotPrimes, tt.wantPrimes) {
			t.Errorf("sieveOfEratosthenes() = %v, want %v", gotPrimes, tt.wantPrimes)
		}
	}
}

func TestExpression_IsValidNumber(t *testing.T) {
	tests := []struct {
		expr           string
		parameters     RandomParameters
		checkPrecision bool
		wantErr        bool
	}{
		{
			"2a - sin(a) + exp(1 + a)", RandomParameters{NewVar('a'): mustParse(t, "2")}, false, false,
		},
		{
			"2a + b", RandomParameters{NewVar('a'): mustParse(t, "2")}, false, true,
		},
		{
			"1/0", RandomParameters{}, false, true,
		},
		{
			"1/a", RandomParameters{NewVar('a'): mustParse(t, "randInt(0;4)")}, false, true,
		},
		{
			"1/a", RandomParameters{NewVar('a'): mustParse(t, "randInt(1;4)")}, false, false,
		},
		{
			"1/a", RandomParameters{NewVar('a'): mustParse(t, "randDecDen()")}, true, false,
		},
		{
			"(v_f - v_i) / v_i", RandomParameters{NewVarI('v', "f"): mustParse(t, "randint(1;10)"), NewVarI('v', "i"): mustParse(t, "randDecDen()")}, true, false,
		},
		{
			"round(1/3; 3)", nil, true, false,
		},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)

		allValid := true
		for range [100]int{} {
			vars, err := tt.parameters.Instantiate()
			if err != nil {
				t.Fatal(err)
			}

			err = expr.IsValidNumber(vars, tt.checkPrecision, true)
			allValid = allValid && err == nil
		}
		if !allValid != tt.wantErr {
			t.Errorf("Expression.IsValidNumber(%s) got = %v", tt.expr, !allValid)
		}
	}
}

func TestExpression_IsValidProba(t *testing.T) {
	tests := []struct {
		expr string
		vars Vars
		want bool
	}{
		{
			"1.1", nil, false,
		},
		{
			"1/a", Vars{NewVar('a'): mustParse(t, "4")}, true,
		},
		{
			"1/a", Vars{NewVar('a'): mustParse(t, "randDecDen()")}, true,
		},
		{
			"0.2 + 1/a", Vars{NewVar('a'): mustParse(t, "3")}, false,
		},
		{
			"round(1/3; 3)", nil, true,
		},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		err := expr.IsValidProba(tt.vars)
		if (err == nil) != tt.want {
			t.Errorf("Expression.IsValidProba(%s) got = %v, want %v", tt.expr, err, tt.want)
		}
	}
}

func mustParseMany(t *testing.T, exprs []string) []*Expr {
	out := make([]*Expr, len(exprs))
	for i, s := range exprs {
		out[i] = mustParse(t, s)
	}
	return out
}

func TestExpression_AreSortedNumbers(t *testing.T) {
	tests := []struct {
		exprs []string
		vars  Vars
		want  bool
	}{
		{
			[]string{"1", "2", "a"}, Vars{NewVar('a'): mustParse(t, "3")}, true,
		},
		{
			[]string{"1", "2", "a"}, Vars{NewVar('a'): mustParse(t, "1")}, false,
		},
		{
			[]string{"1", "2", "b"}, Vars{NewVar('a'): mustParse(t, "3")}, false,
		},
	}
	for _, tt := range tests {
		exprs := mustParseMany(t, tt.exprs)

		err := AreSortedNumbers(exprs, tt.vars)
		if (err == nil) != tt.want {
			t.Errorf("AreSortedNumbers(%s) got = %v, want %v", tt.exprs, err, tt.want)
		}
	}
}

func TestExpression_IsValidIndex(t *testing.T) {
	tests := []struct {
		expr   string
		vars   Vars
		length int
		want   bool
	}{
		{
			"+1 + 1 * isZero(a-1) + 2 * isZero(a-2) + 3*isZero(a-3)", Vars{NewVar('a'): mustParse(t, "2")}, 4, true,
		},
		{
			"+1 + 1 * isZero(a-1) + 2 * isZero(a-2) + 3*isZero(a-3)", Vars{NewVar('a'): mustParse(t, "3")}, 4, true,
		},
		{
			"+1 + 1 * isZero(a-1) + 2 * isZero(a-2) + 2.5*isZero(a-3)", Vars{NewVar('a'): mustParse(t, "3")}, 4, false,
		},
		{
			"+1 + 1 * isZero(a-1) + 2 * isZero(a-2) + 4*isZero(a-3)", Vars{NewVar('a'): mustParse(t, "3")}, 4, false,
		},
		{
			"+1 + 1 * isZero(a^2 - b^2 - c^2) + 2*isZero(b^2 - a^2 - c^2) + 3*isZero(c^2 - a^2 - b^2)", Vars{
				NewVar('a'): mustParse(t, "4"), // BC
				NewVar('b'): mustParse(t, "5"), // AC
				NewVar('c'): mustParse(t, "6"), // AB
			}, 4, true,
		},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		err := expr.IsValidIndex(tt.vars, tt.length)
		if (err == nil) != tt.want {
			t.Errorf("Expression.IsValidIndex() got = %v, want %v", err, tt.want)
		}
	}
}

func TestFunctionDefinition_IsValid(t *testing.T) {
	tests := []struct {
		expr     string
		variable rune
		vars     Vars
		from, to string
		bound    float64 // expected extrema
		want     bool
	}{
		{"2x + 1", 'x', nil, "-10", "10", 25, true},
		{"2x + 1", 'x', nil, "-10", "10", 10, false},
		{"2x + 1", 'x', nil, "2", "2", 10, false},
		{"2x + a", 'x', nil, "-10", "10", 10, false},
		{"1/x", 'x', nil, "-10", "10", 100, false},
		{"exp(x)", 'x', nil, "-10", "10", 100, false},
		{"ax + b", 'x', Vars{
			NewVar('a'): mustParse(t, "3"),
			NewVar('b'): mustParse(t, "7"),
		}, "-10", "10", 100, true},
		{"ax + b", 'x', Vars{
			NewVar('a'): mustParse(t, "90"),
			NewVar('b'): mustParse(t, "7"),
		}, "-10", "10", 100, false},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		fn := FunctionExpr{
			Function: expr,
			Variable: NewVar(tt.variable),
		}
		err := fn.IsValid(Domain{mustParse(t, tt.from), mustParse(t, tt.to)}, tt.vars, tt.bound)
		if (err == nil) != tt.want {
			t.Errorf("Expression.AreFxsIntegers() got = %v, want %v", err, tt.want)
		}
	}
}

func TestAreDisjointsDomains(t *testing.T) {
	tests := []struct {
		domains [][2]string
		vars    Vars
		wantErr bool
	}{
		{
			[][2]string{{"0", "1"}}, nil, false,
		},
		{
			[][2]string{{"0", "1"}, {"1", "2"}}, nil, false,
		},
		{
			[][2]string{{"0", "1"}, {"0", "2"}}, nil, true,
		},
		{
			[][2]string{{"0", "1"}, {"x", "0.5"}}, Vars{NewVar('x'): newNb(0)}, true,
		},
	}
	for _, tt := range tests {
		domains := make([]Domain, len(tt.domains))
		for i, d := range tt.domains {
			domains[i] = Domain{mustParse(t, d[0]), mustParse(t, d[1])}
		}
		if err := AreDisjointsDomains(domains, tt.vars); (err != nil) != tt.wantErr {
			t.Errorf("AreDisjointsDomains() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}
