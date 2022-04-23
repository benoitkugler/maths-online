package utils

import "testing"

func TestSampleIndex(t *testing.T) {
	tests := []struct {
		args []float64
		want int
	}{
		{
			[]float64{0, 0, 0, 1}, 3,
		},
		{
			[]float64{0, 1, 0, 0}, 1,
		},
	}
	for _, tt := range tests {
		if got := SampleIndex(tt.args); got != tt.want {
			t.Errorf("SampleIndex() = %v, want %v", got, tt.want)
		}
	}
}
