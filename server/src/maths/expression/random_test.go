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
			map[Variable]string{'a': "a +1"}, nil, true,
		},
		{
			map[Variable]string{'a': "a + b + 1", 'b': "8"}, nil, true,
		},
		{
			map[Variable]string{'a': "b + 1", 'b': "a+2"}, nil, true,
		},
		{
			map[Variable]string{'a': "b + 1"}, nil, true,
		},
		{
			map[Variable]string{'a': "b + 1", 'b': " 2 * 3"}, Variables{'a': 7, 'b': 6}, false,
		},
		{
			map[Variable]string{'a': "b + 1", 'b': " c+1", 'c': "8"}, Variables{'a': 10, 'b': 9, 'c': 8}, false,
		},
		{
			map[Variable]string{'a': "0*randInt(1;3)"}, Variables{'a': 0}, false,
		},
		{
			map[Variable]string{'a': "randInt(1;1)", 'b': "2*a"}, Variables{'a': 1, 'b': 2}, false,
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
			'a': mustParse(t, "3*randInt(1; 10)"),
			'b': mustParse(t, "-a"),
		}
		values, err := rv.Instantiate()
		if err != nil {
			t.Fatal(err)
		}
		if values['a'] != -values['b'] {
			t.Fatal(values)
		}

		if a := values['a']; a < 3 || a > 30 {
			t.Fatal(a)
		}

		rv = RandomParameters{
			'a': mustParse(t, "randInt(1; 10)"),
			'b': mustParse(t, "sgn(2*randInt(0;1)-1) * a"),
		}
		values, err = rv.Instantiate()
		if err != nil {
			t.Fatal(err)
		}
		if a := values['a']; a < 1 || a > 10 {
			t.Fatal(a)
		}
		if a, b := values['a'], values['b']; math.Abs(a) != math.Abs(b) {
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
			"1 * isZero(a-1) + 2 * isZero(a-2) + 3*isZero(a-3)", RandomParameters{'a': mustParse(t, "2")}, 4, true,
		},
		{
			"1 * isZero(a-1) + 2 * isZero(a-2) + 3*isZero(a-3)", RandomParameters{'a': mustParse(t, "randInt(0;3)")}, 4, true,
		},
		{
			"1 * isZero(a-1) + 2 * isZero(a-2) + 2.5*isZero(a-3)", RandomParameters{'a': mustParse(t, "randInt(0;3)")}, 4, false,
		},
		{
			"1 * isZero(a-1) + 2 * isZero(a-2) + 4*isZero(a-3)", RandomParameters{'a': mustParse(t, "randInt(0;3)")}, 4, false,
		},
		{
			"1 * isZero(a^2 - b^2 - c^2) + 2*isZero(b^2 - a^2 - c^2) + 3*isZero(c^2 - a^2 - b^2)", RandomParameters{
				'a': mustParse(t, "randInt(3;12)"), // BC
				'b': mustParse(t, "randInt(3;12)"), // AC
				'c': mustParse(t, "randInt(3;12)"), // AB
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
	params := RandomParameters{'a': randExpr}

	_, freq := expr.IsValidIndex(params, 3)
	fmt.Println(freq) // approx 100 - 25 = 75
}
