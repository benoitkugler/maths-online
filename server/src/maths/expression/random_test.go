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
			map[Variable]string{NewVariable('a'): "b + 1", NewVariable('b'): " 2 * 3"}, Variables{NewVariable('a'): 7, NewVariable('b'): 6}, false,
		},
		{
			map[Variable]string{NewVariable('a'): "b + 1", NewVariable('b'): " c+1", NewVariable('c'): "8"}, Variables{NewVariable('a'): 10, NewVariable('b'): 9, NewVariable('c'): 8}, false,
		},
		{
			map[Variable]string{NewVariable('a'): "0*randInt(1;3)"}, Variables{NewVariable('a'): 0}, false,
		},
		{
			map[Variable]string{NewVariable('a'): "randInt(1;1)", NewVariable('b'): "2*a"}, Variables{NewVariable('a'): 1, NewVariable('b'): 2}, false,
		},
	}
	for _, tt := range tests {
		rv := make(RandomParameters)
		for v, e := range tt.rv {
			rv[v] = mustParse(t, e)
		}

		got, err := rv.Instantiate()
		if err != nil {
			err, ok := err.(InvalidRandomVariable)
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
		if values[NewVariable('a')] != -values[NewVariable('b')] {
			t.Fatal(values)
		}

		if a := values[NewVariable('a')]; a < 3 || a > 30 {
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
		if a := values[NewVariable('a')]; a < 1 || a > 10 {
			t.Fatal(a)
		}
		if a, b := values[NewVariable('a')], values[NewVariable('b')]; math.Abs(a) != math.Abs(b) {
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
	expr, _, _ := Parse("1 * isZero(a-1) + 2 * isZero(a-2) + 2.5*isZero(a-3)")
	randExpr, _, _ := Parse("randInt(0;3)")
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
