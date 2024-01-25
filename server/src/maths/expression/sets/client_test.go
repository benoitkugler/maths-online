package sets

import (
	"reflect"
	"testing"

	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func leaf(l Set) ListNode {
	return ListNode{Op: SLeaf, Leaf: l}
}

func inter(args ...ListNode) ListNode {
	return ListNode{Op: SInter, Args: args}
}

func union(args ...ListNode) ListNode {
	return ListNode{Op: SUnion, Args: args}
}

func comp(arg ListNode) ListNode {
	return ListNode{Op: SComplement, Args: []ListNode{arg}}
}

func TestNewSet(t *testing.T) {
	for _, test := range []struct {
		in   BinNode
		want ListNode
	}{
		{
			Set(0),
			leaf(0),
		},
		{
			Inter{Left: Set(0), Right: Set(1)},
			inter(leaf(0), leaf(1)),
		},
		{
			Complement{Right: Set(0)},
			comp(leaf(0)),
		},
		{
			Union{Left: Set(0), Right: Set(0)},
			union(leaf(0), leaf(0)),
		},
		{
			Inter{
				Left:  Inter{Left: Set(0), Right: Set(1)},
				Right: Inter{Left: Set(1), Right: Set(2)},
			},
			inter(leaf(0), leaf(1), leaf(1), leaf(2)),
		},
	} {
		got := BinarySet{Root: test.in}.ToList().Expr
		tu.Assert(t, reflect.DeepEqual(got, test.want))

		roundTrip := BinarySet{Root: got.ToBin()}.ToList().Expr
		tu.Assert(t, reflect.DeepEqual(roundTrip, test.want))
	}
}
