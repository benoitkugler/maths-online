package questions

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/benoitkugler/maths-online/maths/questions/client"
	cl "github.com/benoitkugler/maths-online/maths/questions/client"
)

func tl(text string) cl.TextLine {
	return cl.TextLine{{Text: text}}
}

func ml(text string) cl.TextLine {
	return cl.TextLine{{Text: text, IsMath: true}}
}

// TestParityProof shows how the following
// property proof may be modelled :
// If m and n are even, then m+n is even
func TestParityProof(t *testing.T) {
	n := proofNodeIns{
		Op:    cl.And,
		Left:  proofSequenceIns{proofStatementIns(tl("m est pair")), proofStatementIns(ml("m = 2k"))},
		Right: proofSequenceIns{proofStatementIns(tl("n est pair")), proofStatementIns(ml("n = 2k'"))},
	}

	proof := ProofFieldInstance{Answer: proofSequenceIns{
		n,
		proofEqualityIns{ml("m+n"), ml("2k+2k'"), ml("2(k+k')")},
		proofEqualityIns{ml("m+n"), ml("2k''")},
		proofStatementIns(tl("m+n est pair")),
	}}

	block := proof.toClient()

	b, err := json.MarshalIndent(block, " ", " ")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(b))

	var block2 client.ProofFieldBlock
	err = json.Unmarshal(b, &block2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestProof_IsEquivalent(t *testing.T) {
	s1, s2 := proofStatementIns(tl("n is odd")), proofStatementIns(tl("n is even"))
	tests := []struct {
		pr1  proofSequenceIns
		pr2  proofSequenceIns
		want bool
	}{
		{
			proofSequenceIns{}, proofSequenceIns{}, true,
		},
		{
			proofSequenceIns{proofNodeIns{Left: s1, Right: s2, Op: cl.And}}, proofSequenceIns{proofNodeIns{Left: s2, Right: s1, Op: cl.And}}, true,
		},
		{
			proofSequenceIns{proofNodeIns{Left: s1, Right: s2, Op: cl.Or}}, proofSequenceIns{proofNodeIns{Left: s1, Right: s2, Op: cl.And}}, false,
		},
		{
			proofSequenceIns{proofEqualityIns{tl("a"), tl("b")}}, proofSequenceIns{proofEqualityIns{tl("b"), tl("a")}}, false,
		},
	}
	for _, tt := range tests {
		if got := tt.pr1.isEquivalent(tt.pr2.toClient()); got != tt.want {
			t.Errorf("Proof.IsEquivalent() = %v, want %v", got, tt.want)
		}
	}
}
