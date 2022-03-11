package expression

import (
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
			map[Variable]string{'a': "0*randInt(1,3)"}, Variables{'a': 0}, false,
		},
		{
			map[Variable]string{'a': "randInt(1,1)", 'b': "2*a"}, Variables{'a': 1, 'b': 2}, false,
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
			'a': mustParse(t, "3*randInt(1, 10)"),
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
			'a': mustParse(t, "randInt(1, 10)"),
			'b': mustParse(t, "sgn(2*randInt(0,1)-1) * a"),
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
