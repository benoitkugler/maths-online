package proof

import (
	"fmt"
	"testing"
)

// ExampleParityProof shows how the following
// property proof may be modelled :
// If m and n are even, then m+n is even
func Test(t *testing.T) {
	n := Node{
		Op:    And,
		Left:  ProofPart{Statement("m est pair"), Equality{"m", "2k"}},
		Right: ProofPart{Statement("n est pair"), Equality{"n", "2k'"}},
	}

	proof := Proof{Root: ProofPart{
		n,
		Equality{"m+n", "2k+2k'", "2(k+k')"},
		Equality{"m+n", "2k''"},
		Statement("m+n est pair"),
	}}
	fmt.Println(proof)
}
