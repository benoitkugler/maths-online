// Package proof implements a field
// asking for a (structured) proof.
package proof

import (
	"fmt"
	"strings"

	"github.com/benoitkugler/maths-online/utils"
)

//go:generate ../../../../../structgen/structgen -source=proof.go -mode=itfs-json:gen_itfs.go -mode=dart:../../../../eleve/lib/questions/proof.gen.dart

// Binary represents a logical connector between two propositions
type Binary uint8

// Statement is a basic statement
type Statement struct {
	Content string
}

// Equality is an equality of the form A1 = A2 = A3
// TODO: should the order matter ? if not so, adjust isEquivalent
type Equality struct {
	Terms []string
}

func NewEquality(terms ...string) Equality { return Equality{Terms: terms} }

// Node is an higher level assertion, such as
// (m is even) AND (n is odd)
type Node struct {
	Left, Right Assertion
	Op          Binary
}

type Assertions []Assertion

// Sequence is a list of elementary steps needed
// to write a mathematical proof, where each step are
// implicitely connected by a "So" (Donc) connector.
type Sequence struct {
	Parts Assertions
}

func NewSequence(parts ...Assertion) Sequence { return Sequence{Parts: parts} }

// Assertion is the general container for
// an element of the proof
type Assertion interface {
	// shape returns an Assertion with the same concrete type
	// and shape, but with empty terms
	shape() Assertion

	// terms returns one or more parts constituing this assertion
	terms() []string

	// returns true if other is logically accceptable
	isEquivalent(other Assertion) bool

	fmt.Stringer
}

func (st Statement) shape() Assertion { return Statement{} }

func (eq Equality) shape() Assertion { return Equality{Terms: make([]string, len(eq.Terms))} }

func (nd Node) shape() Assertion {
	return Node{Left: nd.Left.shape(), Right: nd.Right.shape(), Op: 0}
}

func (pp Sequence) shape() Assertion {
	out := make([]Assertion, len(pp.Parts))
	for i, a := range pp.Parts {
		out[i] = a.shape()
	}
	return Sequence{Parts: out}
}

func (st Statement) terms() []string { return []string{st.Content} }

func (eq Equality) terms() []string { return eq.Terms }

func (nd Node) terms() []string {
	return append(nd.Left.terms(), nd.Right.terms()...)
}

func (pp Sequence) terms() []string {
	var out []string
	for _, c := range pp.Parts {
		out = append(out, c.terms()...)
	}
	return out
}

func (st Statement) isEquivalent(other Assertion) bool {
	otherSt, ok := other.(Statement)
	return ok && otherSt.Content == st.Content
}

func (eq Equality) isEquivalent(other Assertion) bool {
	otherEq, ok := other.(Equality)
	return ok && utils.StringSlicesEqual(eq.Terms, otherEq.Terms)
}

func (nd Node) isEquivalent(other Assertion) bool {
	otherNode, ok := other.(Node)
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

func (seq Sequence) isEquivalent(other Assertion) bool {
	otherSeq, ok := other.(Sequence)
	if !ok {
		return false
	}
	if len(seq.Parts) != len(otherSeq.Parts) {
		return false
	}

	// the order matters since the logical So are not always commutative
	// (for n natural) n is even => 2n is even (but not 2n is even => n is even)
	for i := range seq.Parts {
		a1, a2 := seq.Parts[i], otherSeq.Parts[i]
		if !a1.isEquivalent(a2) {
			return false
		}
	}
	return true
}

type Proof struct {
	Root Sequence
}

// proposals returns all the constituent of the various parts
// of the answer, in randomized order
func (pr Proof) proposals() []string {
	out := pr.Root.terms()
	shuffler := utils.NewDeterministicRand([]byte(strings.Join(out, "")))
	shuffler.Shuffle(len(out), func(i, j int) { out[i], out[j] = out[j], out[i] })
	return out
}

// IsEquivalent returns `true` if `pr1` and `pr2` are logically
// similar, in the sense of begin both acceptables.
// It will usually be called on similarly shaped proofs, but does not assume so.
func (pr1 Proof) IsEquivalent(pr2 Proof) bool {
	return pr1.Root.isEquivalent(pr2.Root)
}

// Shape returns a `Proof` with the same shape, but with fields content stripped out.
func (pr Proof) Shape() Proof {
	return Proof{pr.Root.shape().(Sequence)}
}
