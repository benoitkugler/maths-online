package expression

import (
	"reflect"
	"testing"

	"github.com/benoitkugler/maths-online/server/src/maths/expression/sets"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func set(root sets.BinNode, leaves ...string) sets.BinarySet {
	return sets.BinarySet{Root: root, Sets: leaves}
}

func TestToSetExpr(t *testing.T) {
	for _, test := range []struct {
		expr    string
		want    sets.BinarySet
		wantErr bool
	}{
		{"", sets.BinarySet{}, true},
		{"A", set(sets.Set(0), "A"), false},
		{"A_5", set(sets.Set(0), "A_{5}"), false},
		{"A_{i+j}", set(sets.Set(0), "A_{i + j}"), false},
		{"A \u222A B", set(sets.Union{Left: sets.Set(0), Right: sets.Set(1)}, "A", "B"), false},
		{"A \u2229 B", set(sets.Inter{Left: sets.Set(0), Right: sets.Set(1)}, "A", "B"), false},
		{"\u00AC B", set(sets.Complement{Right: sets.Set(0)}, "B"), false},
		{"A - B", set(sets.Inter{
			Left:  sets.Set(0),
			Right: sets.Complement{Right: sets.Set(1)},
		}, "A", "B"), false},
		{"A + 5", sets.BinarySet{}, true},
		{"(A + 5) \u222A B", sets.BinarySet{}, true},
		{"A * B", sets.BinarySet{}, true},
	} {
		e := mustParse(t, test.expr)
		set, err := e.ToBinarySet()
		tu.Assert(t, err != nil == test.wantErr)
		tu.Assert(t, reflect.DeepEqual(set, test.want))
	}
}
