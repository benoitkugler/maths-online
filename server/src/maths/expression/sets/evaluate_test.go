package sets

import (
	"fmt"
	"reflect"
	"testing"

	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

const (
	A = Set(iota)
	B
	C
	D
)

func bs(root BinNode, N int) BinarySet {
	return BinarySet{Sets: make([]string, N), Root: root}
}

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
		tu.Assert(t, reflect.DeepEqual(normalize(test.in), test.want))
	}
}

func TestToNormalized(t *testing.T) {
	for _, test := range []struct {
		in   BinarySet
		want normalized
	}{
		{
			bs(Inter{A, B}, 2),
			normalized{0b11},
		},
		{
			bs(Union{Inter{Complement{A}, B}, C}, 3),
			normalized{0b010, 0b100, 0b101, 0b110, 0b111},
		},
		{
			bs(Inter{Union{A, B}, Union{C, D}}, 4),
			normalized{0b0101, 0b0110, 0b0111, 0b1001, 0b1010, 0b1011, 0b1101, 0b1110, 0b1111},
		},
	} {
		tu.Assert(t, reflect.DeepEqual(test.in.toNormalized(), test.want))
	}
}

func TestIsEquivalent(t *testing.T) {
	for _, test := range []struct {
		s1         BinarySet
		s2         BinNode
		equivalent bool
	}{
		{
			bs(A, 1), A, true,
		},
		{
			bs(Union{A, A}, 1), A, true,
		},
		{
			bs(Inter{A, A}, 1), A, true,
		},
		{
			bs(A, 2), B, false,
		},
		{
			bs(Union{A, B}, 2), Union{B, A}, true,
		},
		{
			bs(Union{A, B}, 2), Union{B, A}, true,
		},
		{
			bs(Union{Inter{A, B}, Inter{A, Complement{B}}}, 2), A, true,
		},
		{
			bs(Union{A, B}, 2), Union{A, Complement{B}}, false,
		},
		{
			bs(A, 2), Union{A, Complement{B}}, false,
		},
		{
			bs(A, 10), Union{A, Complement{D}}, false,
		},
	} {
		_ = fmt.Sprintf("%s", test.s1)
		tu.Assert(t, test.s1.IsEquivalent(test.s2) == test.equivalent)
	}
}
