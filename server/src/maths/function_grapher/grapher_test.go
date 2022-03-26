// Package functiongrapher provides a way to convert an
// arbitrary function expression into a list of quadratic bezier curve.
package functiongrapher

import (
	"math"
	"reflect"
	"testing"

	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/repere"
)

func assertApprox(t *testing.T, a, b float64) {
	if math.Abs(a-b) > 1e-5 {
		t.Fatalf("%g != %g", a, b)
	}
}

func assertSegmentApprox(t *testing.T, a, b segment) {
	assertApprox(t, a.dFrom, b.dFrom)
	assertApprox(t, a.dTo, b.dTo)
	assertApprox(t, a.from.X, b.from.X)
	assertApprox(t, a.from.Y, b.from.Y)
	assertApprox(t, a.to.X, b.to.X)
	assertApprox(t, a.to.Y, b.to.Y)
}

func Test_newSegment(t *testing.T) {
	tests := []struct {
		expr string
		from float64
		to   float64
		want segment
	}{
		{"7", 0, 1, segment{from: repere.Coord{X: 0, Y: 7}, to: repere.Coord{X: 1, Y: 7}, dFrom: 0, dTo: 0}},
		{"2x + 1", 0, 1, segment{from: repere.Coord{X: 0, Y: 1}, to: repere.Coord{X: 1, Y: 3}, dFrom: 2, dTo: 2}},
		{"sin(x)", 0, math.Pi, segment{from: repere.Coord{X: 0, Y: 0}, to: repere.Coord{X: math.Pi, Y: 0}, dFrom: 1, dTo: -1}},
		{"ln(x)", 1, 2, segment{from: repere.Coord{X: 1, Y: 0}, to: repere.Coord{X: 2, Y: math.Log(2)}, dFrom: 1, dTo: 1. / 2}},
	}
	for _, tt := range tests {
		expr, _, err := expression.Parse(tt.expr)
		if err != nil {
			t.Fatal(err)
		}
		got := newSegment(expr, 'x', tt.from, tt.to)
		assertSegmentApprox(t, got, tt.want)
	}
}

func Test_segment_toCurve(t *testing.T) {
	type fields struct {
		from  repere.Coord
		to    repere.Coord
		dFrom float64
		dTo   float64
	}
	tests := []struct {
		fields fields
		want   repere.Coord
	}{
		{
			fields: fields{
				repere.Coord{X: 0, Y: 1},
				repere.Coord{X: 2, Y: 3},
				2,
				2,
			},
			want: repere.Coord{X: 1, Y: 2},
		},
		{
			fields: fields{
				repere.Coord{X: -1, Y: 1},
				repere.Coord{X: 1, Y: 1},
				-1,
				1,
			},
			want: repere.Coord{X: 0, Y: 0},
		},
	}
	for _, tt := range tests {
		seg := segment{
			from:  tt.fields.from,
			to:    tt.fields.to,
			dFrom: tt.fields.dFrom,
			dTo:   tt.fields.dTo,
		}
		want := BezierCurve{P0: tt.fields.from, P1: tt.want, P2: tt.fields.to}
		if got := seg.toCurve(); !reflect.DeepEqual(got, want) {
			t.Errorf("segment.toCurve() = %v, want %v", got, want)
		}
	}
}

func TestGraph(t *testing.T) {
	tests := []struct {
		expr string
		vari expression.Variable
		from float64
		to   float64
	}{
		{"7x+1", 'x', 0, 10},
		{"ln(y)", 'y', 1, 10},
	}
	for _, tt := range tests {
		expr, _, _ := expression.Parse(tt.expr)

		got := NewFunctionGraph(expr, tt.vari, tt.from, tt.to)
		if len(got.Segments) != nbStep {
			t.Fatal()
		}

		if got.Bounds.Origin.X < 0 || got.Bounds.Origin.Y < 0 {
			t.Fatal(got.Bounds.Origin)
		}
	}
}
