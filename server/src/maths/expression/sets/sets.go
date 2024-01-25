// Package sets implement basic math sets theory,
// usefull to evaluate an answer.
// For instance, the two expression
// A  and  A ∩ (B U neg(B)) are equal
//
// Note that this package treats sets as opaque, so that
// simplification using relation such as A = B U C are not supported.
package sets

import (
	"fmt"
	"strconv"
)

type BinarySet struct {
	// Leaves stores the LaTeX code of each set,
	// like "A", "B_i" or "B_{i+1}", and is indexed by
	// a Leaf id.
	// A well formed set has no more than [MaxNumberOfLeaves] leaves.
	Leaves []string
	Root   BinNode
}

// BinNode is one node of a binary tree representation
// of a set expression.
//
// The nil value is invalid and does not appear in a well formed
// tree.
type BinNode interface {
	isBinNode()

	String() string
}

func (Set) isBinNode()        {}
func (Union) isBinNode()      {}
func (Inter) isBinNode()      {}
func (Complement) isBinNode() {}

// Set is an index into the [BinarySet.Leaves] slice.
type Set int

func (l Set) String() string { return strconv.Itoa(int(l)) }

type Union struct {
	Left, Right BinNode
}

func (u Union) String() string { return fmt.Sprintf("%s \u222A %s", u.Left, u.Right) }

type Inter struct {
	Left, Right BinNode
}

func (i Inter) String() string { return fmt.Sprintf("%s \u2229 %s", i.Left, i.Right) }

type Complement struct {
	Right BinNode
}

func (c Complement) String() string { return fmt.Sprintf("\u00AC %s", c.Right) }
