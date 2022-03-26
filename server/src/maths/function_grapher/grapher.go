// Package functiongrapher provides a way to convert an
// arbitrary function expression into a list of quadratic bezier curve.
package functiongrapher

import (
	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/repere"
)

type BezierCurve struct {
	P0, P1, P2 repere.Coord `dart-extern:"repere.gen.dart"`
}

func (seg segment) toCurve() BezierCurve {
	// special case when df1 = df2
	if seg.dFrom == seg.dTo {
		return BezierCurve{
			P0: seg.from,
			P1: repere.Coord{X: (seg.from.X + seg.to.X) / 2, Y: (seg.from.Y + seg.to.Y) / 2},
			P2: seg.to,
		}
	}

	p1 := controlFromDerivatives(seg.from, seg.to, seg.dFrom, seg.dTo)
	return BezierCurve{
		P0: seg.from,
		P1: p1,
		P2: seg.to,
	}
}

// compute the control point matching the given derivatives,
// which is the intersection between
// the tangents at from and to
func controlFromDerivatives(from, to repere.Coord, dFrom, dTo float64) repere.Coord {
	xIntersec := (to.Y - from.Y + dFrom*from.X - dTo*to.X) / (dFrom - dTo)
	yIntersec := dFrom*(xIntersec-from.X) + from.Y
	return repere.Coord{X: xIntersec, Y: yIntersec}
}

type segment struct {
	from, to   repere.Coord
	dFrom, dTo float64
}

// expr must be an expression containing only the variable `variable`
func newSegment(expr *expression.Expression, variable expression.Variable, from, to float64) segment {
	f := func(x float64) float64 {
		out, _ := expr.Evaluate(expression.Variables{variable: x})
		return out
	}

	yFrom := f(from)
	yTo := f(to)

	// compute derivative with finite differences
	epsilon := (to - from) / 100_000
	dFrom := (f(from+epsilon) - f(from)) / epsilon
	dTo := (f(to) - f(to-epsilon)) / epsilon

	return segment{
		from:  repere.Coord{X: from, Y: yFrom},
		to:    repere.Coord{X: to, Y: yTo},
		dFrom: dFrom,
		dTo:   dTo,
	}
}

type FunctionGraph struct {
	Segments []BezierCurve
	Bounds   repere.RepereBounds `dart-extern:"repere.gen.dart"`
}

// nbStep is the number of segments used when converting
// a function curve to Bezier curves.
const nbStep = 100

// Graph splits the curve of `expr(vari)` on [from, to] in small chunks on which it is approximated
// by Bezier curves.
func NewFunctionGraph(expr *expression.Expression, vari expression.Variable, from, to float64) FunctionGraph {
	step := (to - from) / nbStep
	curves := make([]BezierCurve, nbStep)
	for i := range curves {
		seg := newSegment(expr, vari, from+float64(i)*step, from+float64(i+1)*step)
		curves[i] = seg.toCurve()
	}

	return FunctionGraph{
		Segments: curves,
		Bounds:   boundingBox(curves),
	}
}

func controlFromPoints(from, to repere.Coord, firstHalf bool) repere.Coord {
	var dFrom, dTo float64
	if from.Y > to.Y { // decreasing
		if firstHalf {
			dFrom, dTo = 0, -4
		} else {
			dFrom, dTo = -4, 0
		}
	} else { // increasing
		if firstHalf {
			dFrom, dTo = 0, +4
		} else {
			dFrom, dTo = +4, 0
		}
	}
	return controlFromDerivatives(from, to, dFrom, dTo)
}

// GraphFromVariations builds one possible representation for the given variations
func GraphFromVariations(xs []expression.Number, ys []expression.Number) FunctionGraph {
	// each segment is drawn with two quadratic curves
	var curves []BezierCurve
	for i := range xs {
		if i >= len(xs)-1 {
			break
		}
		x1, x2 := float64(xs[i]), float64(xs[i+1])
		y1, y2 := float64(ys[i]), float64(ys[i+1])
		xMiddle, yMiddle := (x1+x2)/2, (y1+y2)/2

		// first half
		p0 := repere.Coord{X: x1, Y: y1}
		p2 := repere.Coord{X: xMiddle, Y: yMiddle}
		p1 := controlFromPoints(p0, p2, true)
		curves = append(curves, BezierCurve{p0, p1, p2})

		// second half
		p0 = repere.Coord{X: xMiddle, Y: yMiddle}
		p2 = repere.Coord{X: x2, Y: y2}
		p1 = controlFromPoints(p0, p2, false)
		curves = append(curves, BezierCurve{p0, p1, p2})
	}

	return FunctionGraph{
		Segments: curves,
		Bounds:   boundingBox(curves),
	}
}
