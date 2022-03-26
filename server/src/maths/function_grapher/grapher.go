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

	// compute the control point, which is the intersection between
	// the tangents in seg.from.x and seg.to.x
	xIntersec := (seg.to.Y - seg.from.Y + seg.dFrom*seg.from.X - seg.dTo*seg.to.X) / (seg.dFrom - seg.dTo)
	yIntersec := seg.dFrom*(xIntersec-seg.from.X) + seg.from.Y
	p1 := repere.Coord{X: xIntersec, Y: yIntersec}
	return BezierCurve{
		P0: seg.from,
		P1: p1,
		P2: seg.to,
	}
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
	out := make([]BezierCurve, nbStep)
	for i := range out {
		seg := newSegment(expr, vari, from+float64(i)*step, from+float64(i+1)*step)
		out[i] = seg.toCurve()
	}

	fn := FunctionGraph{Segments: out}
	fn.Bounds = boundingBox(out)
	return fn
}
