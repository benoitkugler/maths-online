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
			map[Variable]string{NewVariable('a'): "a +1"}, nil, true,
		},
		{
			map[Variable]string{NewVariable('a'): "a + b + 1", NewVariable('b'): "8"}, nil, true,
		},
		{
			map[Variable]string{NewVariable('a'): "b + 1", NewVariable('b'): "a+2"}, nil, true,
		},
		{
			map[Variable]string{NewVariable('a'): "b + 1"}, nil, true,
		},
		{
			map[Variable]string{NewVariable('a'): "b + 1", NewVariable('b'): " 2 * 3"}, Variables{NewVariable('a'): NewRN(7), NewVariable('b'): NewRN(6)}, false,
		},
		{
			map[Variable]string{NewVariable('a'): "b + 1", NewVariable('b'): " c+1", NewVariable('c'): "8"}, Variables{NewVariable('a'): NewRN(10), NewVariable('b'): NewRN(9), NewVariable('c'): NewRN(8)}, false,
		},
		{
			map[Variable]string{NewVariable('a'): "0*randInt(1;3)"}, Variables{NewVariable('a'): NewRN(0)}, false,
		},
		{
			map[Variable]string{NewVariable('a'): "randInt(1;1)", NewVariable('b'): "2*a"}, Variables{NewVariable('a'): NewRN(1), NewVariable('b'): NewRN(2)}, false,
		},
		{
			map[Variable]string{NewVariable('a'): "randLetter(A)", NewVariable('b'): "randInt(1;1)"},
			Variables{NewVariable('a'): NewRV(NewVariable('A')), NewVariable('b'): NewRN(1)},
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
		rv := RandomParameters{NewVariable('P'): mustParse(t, "randLetter(A;B;C)")}
		vars, err := rv.Instantiate()
		if err != nil {
			t.Fatal(err)
		}
		resolved := vars[NewVariable('P')]
		if !resolved.IsVariable {
			t.Fatal(resolved)
		}
		if n := resolved.V.Name; !(n == 'A' || n == 'B' || n == 'C') {
			t.Fatal(n)
		}
	}
}

func TestRandomVariables_range(t *testing.T) {
	for range [10]int{} {
		rv := RandomParameters{
			NewVariable('a'): mustParse(t, "3*randInt(1; 10)"),
			NewVariable('b'): mustParse(t, "-a"),
		}
		values, err := rv.Instantiate()
		if err != nil {
			t.Fatal(err)
		}
		if values[NewVariable('a')].N != -values[NewVariable('b')].N {
			t.Fatal(values)
		}

		if a := values[NewVariable('a')].N; a < 3 || a > 30 {
			t.Fatal(a)
		}

		rv = RandomParameters{
			NewVariable('a'): mustParse(t, "randInt(1; 10)"),
			NewVariable('b'): mustParse(t, "sgn(2*randInt(0;1)-1) * a"),
		}
		values, err = rv.Instantiate()
		if err != nil {
			t.Fatal(err)
		}
		if a := values[NewVariable('a')].N; a < 1 || a > 10 {
			t.Fatal(a)
		}
		if a, b := values[NewVariable('a')].N, values[NewVariable('b')].N; math.Abs(a) != math.Abs(b) {
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
		want           bool
	}{
		{
			"2a - sin(a) + exp(1 + a)", RandomParameters{NewVariable('a'): mustParse(t, "2")}, false, true,
		},
		{
			"2a + b", RandomParameters{NewVariable('a'): mustParse(t, "2")}, false, false,
		},
		{
			"1/0", RandomParameters{}, false, false,
		},
		{
			"1/a", RandomParameters{NewVariable('a'): mustParse(t, "randInt(0;4)")}, false, false,
		},
		{
			"1/a", RandomParameters{NewVariable('a'): mustParse(t, "randInt(1;4)")}, false, true,
		},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		got, _ := expr.IsValidNumber(tt.parameters, tt.checkPrecision)
		if got != tt.want {
			t.Errorf("Expression.IsValidIndex() got = %v, want %v", got, tt.want)
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
			"1/a", RandomParameters{NewVariable('a'): mustParse(t, "randInt(4;5)")}, true,
		},
		{
			"1/a", RandomParameters{NewVariable('a'): mustParse(t, "randDecDen()")}, true,
		},
		{
			"0.2 + 1/a", RandomParameters{NewVariable('a'): mustParse(t, "randInt(1;4)")}, false,
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

func TestExpression_AreSortedNumbers(t *testing.T) {
	tests := []struct {
		exprs      []string
		parameters RandomParameters
		want       bool
	}{
		{
			[]string{"1", "2", "a"}, RandomParameters{NewVariable('a'): mustParse(t, "randInt(2;5)")}, true,
		},
		{
			[]string{"1", "2", "a"}, RandomParameters{NewVariable('a'): mustParse(t, "randInt(0;5)")}, false,
		},
		{
			[]string{"1", "2", "b"}, RandomParameters{NewVariable('a'): mustParse(t, "randInt(0;5)")}, false,
		},
	}
	for _, tt := range tests {
		exprs := make([]*Expression, len(tt.exprs))
		for i, s := range tt.exprs {
			exprs[i] = mustParse(t, s)
		}
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
			"1 * isZero(a-1) + 2 * isZero(a-2) + 3*isZero(a-3)", RandomParameters{NewVariable('a'): mustParse(t, "2")}, 4, true,
		},
		{
			"1 * isZero(a-1) + 2 * isZero(a-2) + 3*isZero(a-3)", RandomParameters{NewVariable('a'): mustParse(t, "randInt(0;3)")}, 4, true,
		},
		{
			"1 * isZero(a-1) + 2 * isZero(a-2) + 2.5*isZero(a-3)", RandomParameters{NewVariable('a'): mustParse(t, "randInt(0;3)")}, 4, false,
		},
		{
			"1 * isZero(a-1) + 2 * isZero(a-2) + 4*isZero(a-3)", RandomParameters{NewVariable('a'): mustParse(t, "randInt(0;3)")}, 4, false,
		},
		{
			"1 * isZero(a^2 - b^2 - c^2) + 2*isZero(b^2 - a^2 - c^2) + 3*isZero(c^2 - a^2 - b^2)", RandomParameters{
				NewVariable('a'): mustParse(t, "randInt(3;12)"), // BC
				NewVariable('b'): mustParse(t, "randInt(3;12)"), // AC
				NewVariable('c'): mustParse(t, "randInt(3;12)"), // AB
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

func TestExample(t *testing.T) {
	expr, _ := Parse("1 * isZero(a-1) + 2 * isZero(a-2) + 2.5*isZero(a-3)")
	randExpr, _ := Parse("randInt(0;3)")
	params := RandomParameters{NewVariable('a'): randExpr}

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
		{"ax^2 - 2x + c", RandomParameters{NewVariable('a'): mustParse(t, "randInt(2;4)"), NewVariable('c'): mustParse(t, "7")}, []int{-2, -1, 0, 4}, true},
		{"2x + 0.5", nil, []int{-2, -1, 0, 4}, false},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		got, freq := expr.AreFxsIntegers(tt.parameters, NewVariable('x'), tt.grid)
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
			NewVariable('a'): mustParse(t, "randInt(2;4)"),
			NewVariable('b'): mustParse(t, "7"),
		}, 100, true},
		{"ax + b", 'x', RandomParameters{
			NewVariable('a'): mustParse(t, "randInt(2;100)"),
			NewVariable('b'): mustParse(t, "7"),
		}, 100, false},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		fn := FunctionDefinition{
			FunctionExpr: FunctionExpr{
				Function: expr,
				Variable: NewVariable(tt.variable),
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
