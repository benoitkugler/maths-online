// Package proof implements a field
// asking for a (structured) proof.
package questions

import (
	"fmt"
	"strings"

	cl "github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	"github.com/benoitkugler/maths-online/server/src/utils"
)

type proofStatementIns cl.TextLine

type proofEqualityIns cl.Equality

type proofNodeIns struct {
	Op          cl.Binary
	Left, Right proofAssertionIns
}

type proofSequenceIns []proofAssertionIns

// proofAssertionIns is the general container for
// an element of the proof
type proofAssertionIns interface {
	toClient() cl.Assertion

	// shape returns an Assertion with the same concrete type
	// and shape, but with empty terms
	shape() proofAssertionIns

	// terms returns one or more parts constituing this assertion
	terms() []cl.TextLine

	// isEquivalent returns `true` if `pr1` and `pr2` are logically
	// similar, in the sense of begin both acceptables.
	// It will usually be called on similarly shaped proofs, but does not assume so.
	isEquivalent(other cl.Assertion) bool
}

func (v proofStatementIns) toClient() cl.Assertion { return cl.Statement{Content: cl.TextLine(v)} }
func (v proofEqualityIns) toClient() cl.Assertion  { return cl.Equality(v) }

func (v proofNodeIns) toClient() cl.Assertion {
	return cl.Node{Op: v.Op, Left: v.Left.toClient(), Right: v.Right.toClient()}
}

func (v proofSequenceIns) toClientS() cl.Sequence {
	out := make(cl.Assertions, len(v))
	for i, k := range v {
		out[i] = k.toClient()
	}
	return cl.Sequence{Parts: out}
}
func (v proofSequenceIns) toClient() cl.Assertion { return v.toClientS() }

func (st proofStatementIns) shape() proofAssertionIns { return proofStatementIns{} }

func (eq proofEqualityIns) shape() proofAssertionIns {
	return proofEqualityIns{Terms: make([]cl.TextLine, len(eq.Terms)), WithDef: eq.WithDef}
}

func (nd proofNodeIns) shape() proofAssertionIns {
	return proofNodeIns{Left: nd.Left.shape(), Right: nd.Right.shape(), Op: 0}
}

func (pp proofSequenceIns) shape() proofAssertionIns {
	out := make(proofSequenceIns, len(pp))
	for i, a := range pp {
		out[i] = a.shape()
	}
	return out
}

func (st proofStatementIns) terms() []cl.TextLine { return []cl.TextLine{cl.TextLine(st)} }

func (eq proofEqualityIns) terms() []cl.TextLine {
	out := eq.Terms
	if len(eq.Def) != 0 {
		out = append(out, eq.Def)
	}
	return out
}

func (nd proofNodeIns) terms() []cl.TextLine {
	return append(nd.Left.terms(), nd.Right.terms()...)
}

func (pp proofSequenceIns) terms() []cl.TextLine {
	var out []cl.TextLine
	for _, c := range pp {
		out = append(out, c.terms()...)
	}
	return out
}

func (st proofStatementIns) isEquivalent(other cl.Assertion) bool {
	otherSt, ok := other.(cl.Statement)
	return ok && areLineEquals(cl.TextLine(st), otherSt.Content)
}

func areManyLinesEquals(s1, s2 []cl.TextLine) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, s := range s1 {
		if !areLineEquals(s, s2[i]) {
			return false
		}
	}
	return true
}

func (eq proofEqualityIns) isEquivalent(other cl.Assertion) bool {
	otherEq, ok := other.(cl.Equality)
	return ok && areManyLinesEquals(eq.Terms, otherEq.Terms) && areLineEquals(eq.Def, otherEq.Def)
}

func (nd proofNodeIns) isEquivalent(other cl.Assertion) bool {
	otherNode, ok := other.(cl.Node)
	if !ok {
		return false
	}

	if nd.Op != otherNode.Op {
		return false
	}

	// AND and OR are commutative, so check for the possible combination
	c1 := nd.Left.isEquivalent(otherNode.Left) && nd.Right.isEquivalent(otherNode.Right)
	c2 := nd.Left.isEquivalent(otherNode.Right) && nd.Right.isEquivalent(otherNode.Left)
	return c1 || c2
}

func (seq proofSequenceIns) isEquivalent(other cl.Assertion) bool {
	otherSeq, ok := other.(cl.Sequence)
	if !ok {
		return false
	}
	if len(seq) != len(otherSeq.Parts) {
		return false
	}

	// the order matters since the logical So are not always commutative
	// (for n natural) n is even => 2n is even (but not 2n is even => n is even)
	for i := range seq {
		a1, a2 := seq[i], otherSeq.Parts[i]
		if !a1.isEquivalent(a2) {
			return false
		}
	}
	return true
}

type ProofFieldInstance struct {
	Answer proofSequenceIns
	ID     int
}

// termProposals returns all the constituent of the various parts
// of the answer, in randomized order.
func (pr ProofFieldInstance) termProposals() []cl.TextLine {
	tmp := pr.Answer.terms()
	contents := make([]string, len(tmp))
	for i, t := range tmp {
		contents[i] = textLineToString(t)
	}
	shuffler := utils.NewDeterministicShuffler([]byte(strings.Join(contents, "")), len(tmp))
	out := make([]cl.TextLine, len(tmp))
	shuffler.Shuffle(func(dst, src int) { out[dst] = tmp[src] })
	return tmp
}

func (pr ProofFieldInstance) shape() cl.Proof {
	seq := pr.Answer.shape().toClient().(cl.Sequence)
	return cl.Proof{Root: seq}
}

func (pr ProofFieldInstance) fieldID() int { return pr.ID }

func (pr ProofFieldInstance) toClient() cl.Block {
	return cl.ProofFieldBlock{
		Shape:         pr.shape(),
		TermProposals: pr.termProposals(),
		ID:            pr.ID,
	}
}

func (pr ProofFieldInstance) validateAnswerSyntax(answer cl.Answer) error {
	_, ok := answer.(cl.ProofAnswer)
	if !ok {
		return InvalidFieldAnswer{
			ID:     pr.ID,
			Reason: fmt.Sprintf("expected ProofAnswer, got %T", answer),
		}
	}
	return nil
}

func (pr ProofFieldInstance) evaluateAnswer(answer cl.Answer) (isCorrect bool) {
	return pr.Answer.isEquivalent(answer.(cl.ProofAnswer).Proof.Root)
}

func (pr ProofFieldInstance) correctAnswer() cl.Answer {
	return cl.ProofAnswer{Proof: cl.Proof{Root: pr.Answer.toClientS()}}
}
