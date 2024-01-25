package sets

import (
	"reflect"
	"testing"

	"github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

const (
	A = Set(iota)
	B
	C
	D
)

func TestNormalize(t *testing.T) {
	for _, test := range []struct {
		in   BinNode
		want BinNode
	}{
		{
			A, A,
		},
		{
			Complement{A}, Complement{A},
		},
		{
			Inter{A, B}, Inter{A, B},
		},
		{
			Union{Inter{Complement{A}, B}, C},
			Union{Inter{Complement{A}, B}, C},
		},
		{
			Complement{Complement{A}}, A,
		},
		{
			Complement{Inter{A, B}},
			Union{Complement{A}, Complement{B}},
		},
		{
			Inter{A, Union{B, C}},
			Union{Inter{A, B}, Inter{A, C}},
		},
		{
			Complement{Union{A, B}},
			Inter{Complement{A}, Complement{B}},
		},
		{
			Inter{Union{A, B}, C},
			Union{Inter{A, C}, Inter{B, C}},
		},
		{
			Inter{Union{A, B}, Union{C, D}},
			Union{Union{Inter{A, C}, Inter{A, D}}, Union{Inter{B, C}, Inter{B, D}}},
		},
	} {
		testutils.Assert(t, reflect.DeepEqual(normalize(test.in), test.want))
	}
}
