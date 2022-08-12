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

// shufflingMap returns the mapping from originalIndex -> shuffledIndex
// for the list [0, 1, ..., n-1]
// note that `shuffler` must be freshy created, or at least seeded properly
func shufflingMap(shuffler *rand.Rand, n int) []int {
	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}
	shuffler.Shuffle(len(indices), func(i, j int) { indices[i], indices[j] = indices[j], indices[i] })
	// build the reverse map (quadratic complexity)
	answer := make([]int, n)
	for i := range answer {
		// find the new index of i into indices
		rep := -1
		for r, val := range indices {
			if val == i {
				rep = r
				break
			}
		}
		answer[i] = rep
	}
	return answer
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
