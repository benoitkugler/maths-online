package expression

import (
	"reflect"
	"testing"

	"github.com/benoitkugler/maths-online/server/src/maths/expression/sets"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func TestToSetExpr(t *testing.T) {
	for _, test := range []struct {
		expr    string
		want    sets.RawSetExpr
		wantErr bool
	}{
		{"", nil, false},
		{"A", sets.Leaf("A"), false},
		{"A_5", sets.Leaf("A_{5}"), false},
		{"A_{i+j}", sets.Leaf("A_{i + j}"), false},
		{"A \u222A B", sets.RNode{Left: sets.Leaf("A"), Op: sets.SUnion, Right: sets.Leaf("B")}, false},
		{"A \u2229 B", sets.RNode{Left: sets.Leaf("A"), Op: sets.SInter, Right: sets.Leaf("B")}, false},
		{"\u00AC B", sets.RNode{Op: sets.SComplement, Right: sets.Leaf("B")}, false},
		{"A - B", sets.RNode{
			Left:  sets.Leaf("A"),
			Op:    sets.SInter,
			Right: sets.RNode{Op: sets.SComplement, Right: sets.Leaf("B")},
		}, false},
		{"A + 5", nil, true},
		{"(A + 5) \u222A B", nil, true},
		{"A * B", nil, true},
	} {
		e := mustParse(t, test.expr)
		set, err := e.toRawSetExpr()
		tu.Assert(t, err != nil == test.wantErr)
		tu.Assert(t, reflect.DeepEqual(set, test.want))
	}
}
