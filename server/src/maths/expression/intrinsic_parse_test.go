package expression

import (
	"testing"
)

var validPythagorians = []struct {
	args string
	want pythagorianTriplet
}{
	{
		"a,b,c = pythagorians()", pythagorianTriplet{
			NewVar('a'), NewVar('b'), NewVar('c'), 10,
		},
	},
	{
		"a,b, c = pythagorians( )", pythagorianTriplet{
			NewVar('a'), NewVar('b'), NewVar('c'), 10,
		},
	},
	{
		"a,b_21, c = pythagorians( )", pythagorianTriplet{
			NewVar('a'), Variable{Name: 'b', Indice: "21"}, NewVar('c'), 10,
		},
	},
	{
		"a,b,c = pythagorians(12)", pythagorianTriplet{
			NewVar('a'), NewVar('b'), NewVar('c'), 12,
		},
	},
}

var validOrthogonalProjection = []struct {
	args string
	want orthogonalProjection
}{
	{
		"H = projection(A, B, C)", orthogonalProjection{
			Variable{Name: 'x', Indice: "A"},
			Variable{Name: 'y', Indice: "A"},
			Variable{Name: 'x', Indice: "B"},
			Variable{Name: 'y', Indice: "B"},
			Variable{Name: 'x', Indice: "C"},
			Variable{Name: 'y', Indice: "C"},
			Variable{Name: 'x', Indice: "H"},
			Variable{Name: 'y', Indice: "H"},
		},
	},
	{
		"x, y = projection(A, B, C)", orthogonalProjection{
			Variable{Name: 'x', Indice: "A"},
			Variable{Name: 'y', Indice: "A"},
			Variable{Name: 'x', Indice: "B"},
			Variable{Name: 'y', Indice: "B"},
			Variable{Name: 'x', Indice: "C"},
			Variable{Name: 'y', Indice: "C"},
			Variable{Name: 'x'},
			Variable{Name: 'y'},
		},
	},
}

var validNumberPair = []struct {
	args string
	want numberPair
}{
	{
		"a, b = number_pair_sum(1)",
		numberPair{a: NewVar('a'), b: NewVar('b'), difficulty: 1, isMultiplicative: false},
	},
	{
		"a, b = number_pair_sum(2)",
		numberPair{a: NewVar('a'), b: NewVar('b'), difficulty: 2, isMultiplicative: false},
	},
	{
		"x_a, b = number_pair_sum(2)",
		numberPair{a: NewVarI('x', "a"), b: NewVar('b'), difficulty: 2, isMultiplicative: false},
	},
	{
		"a, b = number_pair_prod(2)",
		numberPair{a: NewVar('a'), b: NewVar('b'), difficulty: 2, isMultiplicative: true},
	},
}

func Test_parseIntrisic(t *testing.T) {
	type test struct {
		args    string
		wantErr bool
	}
	tests := []test{
		{"a, b = unknown()", true},
		{"a, b = = unknown()", true},
		{"a, b =  unknown)", true},
		{"a, b =  unknown)(", true},
		// pythagorians
		{"a, b = pythagorians()", true},
		{"a, b, c = pythagorians(10,21)", true},
		{"a, b, c = pythagorians(10.4)", true},
		// projection
		{"a, b, c = projection()", true},
		{"a, b, c = projection(a,b,c)", true},
		{"a, b = projection()", true},
		{"a, b = projection(a,b)", true},
		// number pair
		{"a = number_pair_sum(1)", true},
		{"a, b,c = number_pair_sum(1)", true},
		{"a, b = number_pair_sum(0)", true},
		{"a, b = number_pair_sum(0, 1)", true},
		{"a, b = number_pair_sum(1.2)", true},
		{"a, b = number_pair_sum(6)", true},
		{"a, b = number_pair_prod(6)", true},
	}
	for _, i := range validPythagorians {
		tests = append(tests, test{i.args, false})
	}
	for _, i := range validOrthogonalProjection {
		tests = append(tests, test{i.args, false})
	}
	for _, i := range validNumberPair {
		tests = append(tests, test{i.args, false})
	}
	for _, tt := range tests {
		p := NewRandomParameters()
		err := p.ParseIntrinsic(tt.args)
		if (err != nil) != tt.wantErr {
			t.Errorf("parseIntrisic() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
	}
}
