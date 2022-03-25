// Package functiongrapher provides a way to convert an
// arbitrary function expression into a list of quadratic bezier curve.
package functiongrapher

import (
	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/repere"
)

type BezierCurve struct {
	P0, P1, P2 repere.Coord
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

// TODO: tests

// expr must be an expression containing only the variable `variable`
func newSegment(expr *expression.Expression, variable expression.Variable, from, to float64) segment {
	f := func(x float64) float64 {
		out, _ := expr.Evaluate(expression.Variables{variable: x})
		return out
	}

	yFrom := f(from)
	yTo := f(to)

	// compute derivative with finite differences
	epsilon := (to - from) / 100
	dFrom := (f(from-epsilon) - f(from+epsilon)) / epsilon
	dTo := (f(to-epsilon) - f(to+epsilon)) / epsilon

	return segment{
		from:  repere.Coord{X: from, Y: yFrom},
		to:    repere.Coord{X: to, Y: yTo},
		dFrom: dFrom,
		dTo:   dTo,
	}
}
