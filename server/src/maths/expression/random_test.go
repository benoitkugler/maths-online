package expression

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestRandomVariables_instantiate(t *testing.T) {
	tests := []struct {
		rv      map[Variable]string
		want    Variables
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
			map[Variable]string{NewVar('a'): "b + 1", NewVar('b'): " 2 * 3"}, Variables{NewVar('a'): NewNb(7), NewVar('b'): NewNb(6)}, false,
		},
		{
			map[Variable]string{NewVar('a'): "b + 1", NewVar('b'): " c+1", NewVar('c'): "8"}, Variables{NewVar('a'): NewNb(10), NewVar('b'): NewNb(9), NewVar('c'): NewNb(8)}, false,
		},
		{
			map[Variable]string{NewVar('a'): "0*randInt(1;3)"}, Variables{NewVar('a'): NewNb(0)}, false,
		},
		{
			map[Variable]string{NewVar('a'): "randInt(1;1)", NewVar('b'): "2*a"}, Variables{NewVar('a'): NewNb(1), NewVar('b'): NewNb(2)}, false,
		},
		{
			map[Variable]string{NewVar('a'): "randLetter(A)", NewVar('b'): "randInt(1;1)"},
			Variables{NewVar('a'): newVarExpr('A'), NewVar('b'): NewNb(1)},
			false,
		},
		{
			map[Variable]string{NewVar('a'): "7/3"},
			Variables{NewVar('a'): &Expression{atom: div, left: newNb(7), right: newNb(3)}},
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
			t.Errorf("RandomVariables.instantiate() error = %v, wantErr %v", err, tt.wantErr)
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
		if values[NewVar('a')].MustEvaluate(nil) != -values[NewVar('b')].MustEvaluate(nil) {
			t.Fatal(values)
		}

		if a := values[NewVar('a')].MustEvaluate(nil); a < 3 || a > 30 {
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
		if a := values[NewVar('a')].MustEvaluate(nil); a < 1 || a > 10 {
			t.Fatal(a)
		}
		if a, b := values[NewVar('a')].MustEvaluate(nil), values[NewVar('b')].MustEvaluate(nil); math.Abs(a) != math.Abs(b) {
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
		err := expr.IsValidNumber(tt.parameters, tt.checkPrecision, true)
		if (err != nil) != tt.wantErr {
			t.Errorf("Expression.IsValidNumber(%s) got = %v", tt.expr, err)
		}
	}
}

func TestExpression_IsValidProba(t *testing.T) {
	tests := []struct {
		expr       string
		parameters RandomParameters
		want       bool
	}{
		{
			"1.1", RandomParameters{}, false,
		},
		{
			"1/a", RandomParameters{NewVar('a'): mustParse(t, "randInt(4;5)")}, true,
		},
		{
			"1/a", RandomParameters{NewVar('a'): mustParse(t, "randDecDen()")}, true,
		},
		{
			"0.2 + 1/a", RandomParameters{NewVar('a'): mustParse(t, "randInt(1;4)")}, false,
		},
		{
			"round(1/3; 3)", nil, true,
		},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		got, _ := expr.IsValidProba(tt.parameters)
		if got != tt.want {
			t.Errorf("Expression.IsValidProba(%s) got = %v, want %v", tt.expr, got, tt.want)
		}
	}
}

func mustParseMany(t *testing.T, exprs []string) []*Expression {
	out := make([]*Expression, len(exprs))
	for i, s := range exprs {
		out[i] = mustParse(t, s)
	}
	return out
}

func TestExpression_AreSortedNumbers(t *testing.T) {
	tests := []struct {
		exprs      []string
		parameters RandomParameters
		want       bool
	}{
		{
			[]string{"1", "2", "a"}, RandomParameters{NewVar('a'): mustParse(t, "randInt(2;5)")}, true,
		},
		{
			[]string{"1", "2", "a"}, RandomParameters{NewVar('a'): mustParse(t, "randInt(0;5)")}, false,
		},
		{
			[]string{"1", "2", "b"}, RandomParameters{NewVar('a'): mustParse(t, "randInt(0;5)")}, false,
		},
	}
	for _, tt := range tests {
		exprs := mustParseMany(t, tt.exprs)
		got, _ := AreSortedNumbers(exprs, tt.parameters)
		if got != tt.want {
			t.Errorf("Expression.IsValidIndex() got = %v, want %v", got, tt.want)
		}
	}
}

func TestExpression_IsValidIndex(t *testing.T) {
	tests := []struct {
		expr       string
		parameters RandomParameters
		length     int
		want       bool
	}{
		{
			"+1 + 1 * isZero(a-1) + 2 * isZero(a-2) + 3*isZero(a-3)", RandomParameters{NewVar('a'): mustParse(t, "2")}, 4, true,
		},
		{
			"+1 + 1 * isZero(a-1) + 2 * isZero(a-2) + 3*isZero(a-3)", RandomParameters{NewVar('a'): mustParse(t, "randInt(0;3)")}, 4, true,
		},
		{
			"+1 + 1 * isZero(a-1) + 2 * isZero(a-2) + 2.5*isZero(a-3)", RandomParameters{NewVar('a'): mustParse(t, "randInt(0;3)")}, 4, false,
		},
		{
			"+1 + 1 * isZero(a-1) + 2 * isZero(a-2) + 4*isZero(a-3)", RandomParameters{NewVar('a'): mustParse(t, "randInt(0;3)")}, 4, false,
		},
		{
			"+1 + 1 * isZero(a^2 - b^2 - c^2) + 2*isZero(b^2 - a^2 - c^2) + 3*isZero(c^2 - a^2 - b^2)", RandomParameters{
				NewVar('a'): mustParse(t, "randInt(3;12)"), // BC
				NewVar('b'): mustParse(t, "randInt(3;12)"), // AC
				NewVar('c'): mustParse(t, "randInt(3;12)"), // AB
			}, 4, true,
		},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		got, _ := expr.IsValidIndex(tt.parameters, tt.length)
		if got != tt.want {
			t.Errorf("Expression.IsValidIndex() got = %v, want %v", got, tt.want)
		}
	}
}

func TestExpression_IsValidInteger(t *testing.T) {
	tests := []struct {
		expr       string
		parameters RandomParameters
		want       bool
	}{
		{
			"1 * a - 1", RandomParameters{NewVar('a'): mustParse(t, "2")}, true,
		},
		{
			"1 / a - 1", RandomParameters{NewVar('a'): mustParse(t, "2")}, false,
		},
		{
			"2.5", nil, false,
		},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		got, _ := expr.IsValidInteger(tt.parameters)
		if got != tt.want {
			t.Errorf("Expression.IsValidInteger() got = %v, want %v", got, tt.want)
		}
	}
}

func TestExample(t *testing.T) {
	expr, _ := Parse("1 + 1 * isZero(a-1) + 2 * isZero(a-2) + 2.5*isZero(a-3)")
	randExpr, _ := Parse("randInt(0;3)")
	params := RandomParameters{NewVar('a'): randExpr}

	_, freq := expr.IsValidIndex(params, 3)
	fmt.Println(freq) // approx 100 - 25 = 75
}

func TestExpression_AreFxsIntegers(t *testing.T) {
	tests := []struct {
		expr       string
		parameters RandomParameters
		grid       []int
		want       bool
	}{
		{"2x + 1", nil, []int{-2, -1, 0, 4}, true},
		{"ax^2 - 2x + c", RandomParameters{NewVar('a'): mustParse(t, "randInt(2;4)"), NewVar('c'): mustParse(t, "7")}, []int{-2, -1, 0, 4}, true},
		{"2x + 0.5", nil, []int{-2, -1, 0, 4}, false},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		got, freq := FunctionExpr{Function: expr, Variable: NewVar('x')}.AreFxsIntegers(tt.parameters, tt.grid)
		if got != tt.want {
			t.Errorf("Expression.AreFxsIntegers() got = %v (%d), want %v", got, freq, tt.want)
		}
	}
}

func TestFunctionDefinition_IsValid(t *testing.T) {
	tests := []struct {
		expr       string
		variable   rune
		parameters RandomParameters
		bound      float64 // extrema on [-10;10]
		want       bool
	}{
		{"2x + 1", 'x', nil, 25, true},
		{"2x + 1", 'x', nil, 10, false},
		{"2x + a", 'x', nil, 10, false},
		{"1/x", 'x', nil, 100, false},
		{"exp(x)", 'x', nil, 100, false},
		{"ax + b", 'x', RandomParameters{
			NewVar('a'): mustParse(t, "randInt(2;4)"),
			NewVar('b'): mustParse(t, "7"),
		}, 100, true},
		{"ax + b", 'x', RandomParameters{
			NewVar('a'): mustParse(t, "randInt(2;100)"),
			NewVar('b'): mustParse(t, "7"),
		}, 100, false},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		fn := FunctionDefinition{
			FunctionExpr: FunctionExpr{
				Function: expr,
				Variable: NewVar(tt.variable),
			},
			From: -10,
			To:   10,
		}
		got, freq := fn.IsValid(tt.parameters, tt.bound)
		if got != tt.want {
			t.Errorf("Expression.AreFxsIntegers() got = %v (%d), want %v", got, freq, tt.want)
		}
	}
}

func TestAreExprsValidDict(t *testing.T) {
	tests := []struct {
		exprs      []string
		refs       []string
		parameters RandomParameters
		wantErr1   bool
		wantErr2   bool
	}{
		{[]string{"A", "B"}, nil, nil, false, false},
		{[]string{"A", "R"}, nil, RandomParameters{NewVar('R'): mustParse(t, "randLetter(U;V)")}, false, false},
		{[]string{"A", "R"}, nil, RandomParameters{NewVar('R'): mustParse(t, "randLetter(U;A)")}, true, false},
		{[]string{"A_1", "R"}, nil, RandomParameters{NewVar('R'): mustParse(t, "randLetter(U;A_1)")}, true, false},
		{[]string{"A_c2", "R"}, nil, RandomParameters{NewVar('R'): mustParse(t, "randLetter(U;A_c2)")}, true, false},
		{[]string{"A", "B"}, []string{"A", "A"}, nil, false, false},
		{[]string{"A", "B"}, []string{"A", "C"}, nil, false, true},
		{[]string{"A", "R"}, []string{"R"}, RandomParameters{NewVar('R'): mustParse(t, "randLetter(U;V)")}, false, false},
		{[]string{"A", "B", "C"}, []string{"R"}, RandomParameters{NewVar('R'): mustParse(t, "randLetter(A;B;C)")}, false, false},
	}
	for _, tt := range tests {
		exprs := mustParseMany(t, tt.exprs)
		refs := mustParseMany(t, tt.refs)
		if err := AreExprsDistincsNames(exprs, tt.parameters); (err != nil) != tt.wantErr1 {
			t.Errorf("AreExprsDistinctsNames() error = %v, wantErr1 %v", err, tt.wantErr1)
		}
		if err := AreExprsValidRefs(exprs, refs, tt.parameters); (err != nil) != tt.wantErr2 {
			t.Errorf("AreExprsDistinctsNames() error = %v, wantErr2 %v", err, tt.wantErr2)
		}
	}
}
