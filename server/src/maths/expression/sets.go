package expression

import (
	"fmt"

	"github.com/benoitkugler/maths-online/server/src/maths/expression/sets"
)

// ToSetExpr converts the given expression to a specialized set format,
// or returns an error if [e] is not a valid set expression.
func (e *Expr) ToSetExpr() (sets.Set, error) {
	tmp, err := e.toRawSetExpr()
	if err != nil {
		return sets.Set{}, err
	}
	return sets.NewSet(tmp)
}

func (e *Expr) toRawSetExpr() (sets.RawSetExpr, error) {
	if e == nil {
		return nil, nil
	}

	switch atom := e.atom.(type) {
	case indice, Variable:
		return sets.Leaf(e.AsLaTeX()), nil
	case operator:
		// recurse
		left, err := e.left.toRawSetExpr()
		if err != nil {
			return nil, err
		}
		right, err := e.right.toRawSetExpr()
		if err != nil {
			return nil, err
		}
		out := sets.RNode{Left: left, Right: right}
		switch atom {
		case minus:
			// rewrite A - B to A ∩ comp(B)
			out = sets.RNode{
				Left:  left,
				Op:    sets.SInter,
				Right: sets.RNode{Op: sets.SComplement, Right: right},
			}
		case union:
			out.Op = sets.SUnion
		case intersection:
			out.Op = sets.SInter
		case complement:
			out.Op = sets.SComplement
		default:
			// no other operator are supported
			return nil, fmt.Errorf("L'opérateur %s n'est pas défini sur les ensembles.", atom)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("L'expression %s ne définit pas un ensemble.", e.String())
	}
}
