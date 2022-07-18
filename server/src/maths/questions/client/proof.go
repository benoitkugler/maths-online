package client

// Binary represents a logical connector between two propositions
type Binary uint8

type Assertions []Assertion

// Assertion is the general container for
// an element of the proof
type Assertion interface {
	isAssertion()
}

func (st Statement) isAssertion() {}
func (eq Equality) isAssertion()  {}
func (nd Node) isAssertion()      {}
func (pp Sequence) isAssertion()  {}

type Proof struct {
	Root Sequence
}
