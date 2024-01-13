package utils

import (
	"math/rand"
	"testing"
)

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

func TestShuffler(t *testing.T) {
	sh := NewDeterministicShuffler([]byte{1, 2, 4, 6, 8}, 4)

	input := []uint8{1, 2, 3, 4}
	output := make([]uint8, 4)
	sh.Shuffle(func(dst, src int) { output[dst] = input[src] })

	shMap := sh.OriginalToShuffled()
	for original, value := range input {
		shuffledIndex := shMap[original]
		if output[shuffledIndex] != value {
			t.Fatal()
		}
	}
}

func TestRandomID(t *testing.T) {
	rand.Seed(1)
	i1 := RandomID(false, 16, func(s string) bool { return false })
	rand.Seed(1)
	i2 := RandomID(false, 16, func(s string) bool { return s == i1 })
	if i1 == i2 {
		t.Fatal("duplicate!")
	}
}
