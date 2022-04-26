package expression

import (
	"reflect"
	"testing"
)

func Test_parseIntrisic(t *testing.T) {
	tests := []struct {
		args    string
		want    Intrinsic
		wantErr bool
	}{
		{"a, b = unknown()", nil, true},
		{"a, b = = unknown()", nil, true},
		{"a, b =  unknown)", nil, true},
		{"a, b =  unknown)(", nil, true},
		// pythagorians
		{
			"a,b,c = pythagorians()", PythagorianTriplet{
				NewVar('a'), NewVar('b'), NewVar('c'), 10,
			}, false,
		},
		{
			"a,b, c = pythagorians( )", PythagorianTriplet{
				NewVar('a'), NewVar('b'), NewVar('c'), 10,
			}, false,
		},
		{
			"a,b_21, c = pythagorians( )", PythagorianTriplet{
				NewVar('a'), Variable{Name: 'b', Indice: "21"}, NewVar('c'), 10,
			}, false,
		},
		{
			"a,b,c = pythagorians(12)", PythagorianTriplet{
				NewVar('a'), NewVar('b'), NewVar('c'), 12,
			}, false,
		},
		{"a, b = pythagorians()", PythagorianTriplet{}, true},
		{"a, b, c = pythagorians(10,21)", PythagorianTriplet{}, true},
		{"a, b, c = pythagorians(10.4)", PythagorianTriplet{}, true},
		// projection
		{
			"H = projection(A, B, C)", OrthogonalProjection{
				Variable{Name: 'x', Indice: "A"},
				Variable{Name: 'y', Indice: "A"},
				Variable{Name: 'x', Indice: "B"},
				Variable{Name: 'y', Indice: "B"},
				Variable{Name: 'x', Indice: "C"},
				Variable{Name: 'y', Indice: "C"},
				Variable{Name: 'x', Indice: "H"},
				Variable{Name: 'y', Indice: "H"},
			}, false,
		},
		{
			"x, y = projection(A, B, C)", OrthogonalProjection{
				Variable{Name: 'x', Indice: "A"},
				Variable{Name: 'y', Indice: "A"},
				Variable{Name: 'x', Indice: "B"},
				Variable{Name: 'y', Indice: "B"},
				Variable{Name: 'x', Indice: "C"},
				Variable{Name: 'y', Indice: "C"},
				Variable{Name: 'x'},
				Variable{Name: 'y'},
			}, false,
		},
		{"a, b, c = projection()", OrthogonalProjection{}, true},
		{"a, b, c = projection(a,b,c)", OrthogonalProjection{}, true},
		{"a, b = projection()", OrthogonalProjection{}, true},
		{"a, b = projection(a,b)", OrthogonalProjection{}, true},
	}
	for _, tt := range tests {
		got, err := ParseIntrinsic(tt.args)
		if (err != nil) != tt.wantErr {
			t.Errorf("parseIntrisic() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("parseIntrisic() = %v, want %v", got, tt.want)
		}
	}
}
