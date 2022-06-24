// Package proof implements a field
// asking for a (structured) proof.
package proof

import "fmt"

// Binary represents a logical connector between two propositions
type Binary uint8

const (
	And Binary = iota
	Or
)

// Statement is a basic statement
type Statement string

// Equality is an equality of the form A1 = A2 = A3
// TODO: should the order matter ?
type Equality []string

// Node is an higher level assertion, such as
// m is even AND n is odd
type Node struct {
	Left, Right Assertion
	Op          Binary
}

// Assertion is the general container for
// an element of the proof
type Assertion interface {
	isAssertion()
	fmt.Stringer
}

func (Statement) isAssertion() {}
func (Equality) isAssertion()  {}
func (Node) isAssertion()      {}
func (ProofPart) isAssertion() {}

// ProofPart is a list of elementary steps needed
// to write a mathematical proof, where each step are
// implicitely connected by a "So" (Donc) connector.
type ProofPart []Assertion

type Proof struct {
	Root ProofPart
}
