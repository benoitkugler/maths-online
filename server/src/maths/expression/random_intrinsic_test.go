package expression

import (
	"testing"
)

func TestPythagorianTriplet_mergeTo(t *testing.T) {
	tests := []struct {
		RangeStart int
		RangeEnd   int
	}{
		{2, 10},
		{3, 40},
	}
	for _, tt := range tests {
		pt := PythagorianTriplet{
			A: 'a', B: 'b', C: 'c',
			SeedStart: tt.RangeStart,
			SeedEnd:   tt.RangeEnd,
		}
		out := make(RandomParameters)
		pt.MergeTo(out)

		for range [100]int{} {
			vr, err := out.Instantiate()
			if err != nil {
				t.Fatal(err)
			}
			a, b, c := vr['a'], vr['b'], vr['c']
			if _, ok := isInt(a); !ok {
				t.Fatal()
			}
			if _, ok := isInt(b); !ok {
				t.Fatal()
			}
			if _, ok := isInt(c); !ok {
				t.Fatal()
			}
			if a*a+b*b != c*c {
				t.Fatal()
			}
		}
	}
}
