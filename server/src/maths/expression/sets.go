package expression

import (
	"errors"
	"fmt"

	"github.com/benoitkugler/maths-online/server/src/maths/expression/sets"
)

// // ToSetExpr converts the given expression to a specialized set format,
// // or returns an error if [e] is not a valid set expression.
// func (e *Expr) ToSetExpr() (sets.Set, error) {
// 	tmp, err := e.toBinarySet()
// 	if err != nil {
// 		return sets.Set{}, err
// 	}
// 	return sets.NewSet(tmp)
// }

func (expr *Expr) toBinarySet() (sets.BinarySet, error) {
	leaves := map[string]sets.Set{}
	lastID := sets.Set(-1)

	var aux func(e *Expr) (sets.BinNode, error)
	aux = func(e *Expr) (sets.BinNode, error) {
		if e == nil {
			return nil, errors.New("Une expression vide ne définit pas un ensemble.")
		}

		switch atom := e.atom.(type) {
		case indice, Variable:
			leafText := e.AsLaTeX()
			// register or use the leaf
			if _, has := leaves[leafText]; !has {
				leaves[leafText] = lastID + 1
				lastID++
			}
			leaf := leaves[leafText]
			return leaf, nil
		case operator:
			// recurse
			var err error
			var left sets.BinNode
			if atom != complement { // complement as a nil left kid
				left, err = aux(e.left)
				if err != nil {
					return nil, err
				}
			}
			right, err := aux(e.right)
			if err != nil {
				return nil, err
			}
			switch atom {
			case minus:
				// rewrite A - B to A ∩ comp(B)
				return sets.Inter{
					Left:  left,
					Right: sets.Complement{Right: right},
				}, nil
			case union:
				return sets.Union{Left: left, Right: right}, nil
			case intersection:
				return sets.Inter{Left: left, Right: right}, nil
			case complement:
				return sets.Complement{Right: right}, nil
			default:
				// no other operator are supported
				return nil, fmt.Errorf("L'opérateur %s n'est pas défini sur les ensembles.", atom)
			}
		default:
			return nil, fmt.Errorf("L'expression %s ne définit pas un ensemble.", e.String())
		}
	}

	root, err := aux(expr)
	if err != nil {
		return sets.BinarySet{}, err
	}

	if len(leaves) > sets.MaxNumberOfLeaves {
		return sets.BinarySet{}, errors.New("Le nombre d'ensembles dépasse la limite supportée.")
	}

	out := sets.BinarySet{Root: root, Leaves: make([]string, len(leaves))}
	for l, id := range leaves {
		out.Leaves[id] = l
	}
	return out, nil
}
