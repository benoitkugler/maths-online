package proof

import (
	"encoding/json"
	"fmt"
	"testing"
)

// TestParityProof shows how the following
// property proof may be modelled :
// If m and n are even, then m+n is even
func TestParityProof(t *testing.T) {
	n := Node{
		Op:    And,
		Left:  NewSequence(Statement{"m est pair"}, NewEquality("m", "2k")),
		Right: NewSequence(Statement{"n est pair"}, NewEquality("n", "2k'")),
	}

	proof := Proof{NewSequence(
		n,
		NewEquality("m+n", "2k+2k'", "2(k+k')"),
		NewEquality("m+n", "2k''"),
		Statement{"m+n est pair"},
	)}
	fmt.Println(proof)

	b, err := json.Marshal(proof)
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Println(string(b))

	var proof2 Proof
	err = json.Unmarshal(b, &proof2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestProof_IsEquivalent(t *testing.T) {
	s1, s2 := Statement{"n is odd"}, Statement{"n is even"}
	tests := []struct {
		pr1  Sequence
		pr2  Sequence
		want bool
	}{
		{
			NewSequence(), NewSequence(), true,
		},
		{
			NewSequence(Node{Left: s1, Right: s2, Op: And}), NewSequence(Node{Left: s2, Right: s1, Op: And}), true,
		},
		{
			NewSequence(Node{Left: s1, Right: s2, Op: Or}), NewSequence(Node{Left: s1, Right: s2, Op: And}), false,
		},
		{
			NewSequence(NewEquality("a", "b")), NewSequence(NewEquality("b", "a")), false,
		},
	}
	for _, tt := range tests {
		if got := (Proof{tt.pr1}.IsEquivalent(Proof{tt.pr2})); got != tt.want {
			t.Errorf("Proof.IsEquivalent() = %v, want %v", got, tt.want)
		}
	}
}
