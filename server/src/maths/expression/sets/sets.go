// Package sets implement basic math sets theory,
// usefull to evaluate an answer.
// For instance, the two expression
// A  and  A ∩ (B U neg(B)) are equal
//
// Note that this package treats sets as opaque, so that
// simplification using relation such as A = B U C are not supported.
package sets

import "errors"

// RawSetExpr is a temporary representation
// of a set expression.
type RawSetExpr interface {
	isSetExpr()
}

func (Leaf) isSetExpr()  {}
func (RNode) isSetExpr() {}

// Leaf is the LaTeX code of a single set,
// like "A", "B_i" or "B_{i+1}"
type Leaf string

// RNode is an expression using only sets and sets operator
type RNode struct {
	Left, Right RawSetExpr
	Op          SetOp // SLeaf is invalid here
}

type SetOp uint8

const (
	SLeaf = iota
	SUnion
	SInter
	SComplement
)

// SetID is the identifier of a set leaf.
// Currently, the number of different sets must be less than 16.
type SetID uint8

type SetExpr struct {
	// with length 0 if Op == SLeaf, one if Op == SComplement, > 1 otherwise
	Args []SetExpr
	Op   SetOp
	Leaf SetID // valid only if Op == SLeaf
}

type Set struct {
	Leaves []Leaf // indexed by SetID
	Expr   SetExpr
}

func NewSet(raw RawSetExpr) (Set, error) {
	const maxNbLeaves = 16
	leaves := map[Leaf]int{}
	lastID := -1

	var aux func(node RawSetExpr) (SetExpr, error)
	aux = func(node RawSetExpr) (SetExpr, error) {
		switch node := node.(type) {
		case nil:
			return SetExpr{}, errors.New("Ensemble non défini.")
		case Leaf:
			// register or use the leaf
			if _, has := leaves[node]; !has {
				leaves[node] = lastID + 1
				lastID++
			}
			id := leaves[node]
			if id > maxNbLeaves {
				return SetExpr{}, errors.New("Le nombre d'ensembles dépasse la limite supportée.")
			}
			return SetExpr{Op: SLeaf, Leaf: SetID(id)}, nil
		case RNode:
			switch node.Op {
			case SUnion, SInter:
				// group all same ops into one node
				left, err := aux(node.Left)
				if err != nil {
					return SetExpr{}, err
				}
				right, err := aux(node.Right)
				if err != nil {
					return SetExpr{}, err
				}
				var args []SetExpr
				if left.Op == node.Op { // move up
					args = append(args, left.Args...)
				} else {
					args = append(args, left)
				}
				if right.Op == node.Op {
					args = append(args, right.Args...)
				} else {
					args = append(args, right)
				}
				return SetExpr{Op: node.Op, Args: args}, nil
			case SComplement:
				right, err := aux(node.Right)
				if err != nil {
					return SetExpr{}, err
				}
				return SetExpr{Op: node.Op, Args: []SetExpr{right}}, nil
			default:
				panic("exhaustive switch")
			}
		default:
			panic("exhaustive switch")
		}
	}

	root, err := aux(raw)
	if err != nil {
		return Set{}, err
	}

	out := Set{Expr: root, Leaves: make([]Leaf, len(leaves))}
	for l, id := range leaves {
		out.Leaves[id] = l
	}

	return out, nil
}
