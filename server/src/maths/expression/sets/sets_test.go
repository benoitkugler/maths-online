package sets

import (
	"reflect"
	"testing"

	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func leaf(l SetID) SetExpr {
	return SetExpr{Op: SLeaf, Leaf: l}
}

func inter(args ...SetExpr) SetExpr {
	return SetExpr{Op: SInter, Args: args}
}

func union(args ...SetExpr) SetExpr {
	return SetExpr{Op: SUnion, Args: args}
}

func comp(arg SetExpr) SetExpr {
	return SetExpr{Op: SComplement, Args: []SetExpr{arg}}
}

func TestNewSet(t *testing.T) {
	for _, test := range []struct {
		in   RawSetExpr
		want Set
	}{
		{
			nil,
			Set{},
		},
		{
			Leaf("A"),
			Set{Expr: leaf(0), Leaves: []Leaf{"A"}},
		},
		{
			RNode{Left: Leaf("A"), Op: SInter, Right: Leaf("B")},
			Set{Expr: inter(leaf(0), leaf(1)), Leaves: []Leaf{"A", "B"}},
		},
		{
			RNode{Op: SComplement, Right: Leaf("A")},
			Set{Expr: comp(leaf(0)), Leaves: []Leaf{"A"}},
		},
		{
			RNode{Left: Leaf("A"), Op: SUnion, Right: Leaf("A")},
			Set{Expr: union(leaf(0), leaf(0)), Leaves: []Leaf{"A"}},
		},
		{
			RNode{
				Left:  RNode{Left: Leaf("A"), Op: SInter, Right: Leaf("B")},
				Op:    SInter,
				Right: RNode{Left: Leaf("B"), Op: SInter, Right: Leaf("C")},
			},
			Set{Expr: inter(leaf(0), leaf(1), leaf(1), leaf(2)), Leaves: []Leaf{"A", "B", "C"}},
		},
	} {
		got, err := NewSet(test.in)
		tu.Assert(t, err != nil == (len(test.want.Leaves) == 0))
		tu.Assert(t, reflect.DeepEqual(got, test.want))
	}
}
