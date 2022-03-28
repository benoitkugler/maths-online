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
	// special case when df1 = df2
	if dFrom == dTo {
		return repere.Coord{X: (from.X + to.X) / 2, Y: (from.Y + to.Y) / 2}
	}

	xIntersec := (to.Y - from.Y + dFrom*from.X - dTo*to.X) / (dFrom - dTo)
	yIntersec := dFrom*(xIntersec-from.X) + from.Y
	return repere.Coord{X: xIntersec, Y: yIntersec}
}

// compute derivative with finite differences
func computeDf(f func(float64) float64, x, epsilon float64) float64 {
	return (f(x+epsilon) - f(x)) / epsilon
}

type segment struct {
	from, to   repere.Coord
	dFrom, dTo float64
}

// expr must be an expression containing only the variable `variable`
func newSegment(expr *expression.Expression, variable expression.Variable, from, to float64) segment {
	f := expr.Closure(variable)

	yFrom := f(from)
	yTo := f(to)

	// compute derivative with finite differences
	epsilon := (to - from) / 100_000
	dFrom := computeDf(f, from, epsilon)
	dTo := computeDf(f, to, epsilon)

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
// It will panic if `expr` is not valid or evaluable given `vari`.
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
func GraphFromVariations(xs []float64, ys []float64) FunctionGraph {
	// each segment is drawn with two quadratic curves
	var curves []BezierCurve
	for i := range xs {
		if i >= len(xs)-1 {
			break
		}
		x1, x2 := xs[i], xs[i+1]
		y1, y2 := ys[i], ys[i+1]
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

// BoundsFromExpression returns the f(x) and f'(x) values for x in the grid
func BoundsFromExpression(expr *expression.Expression, vari expression.Variable, grid []int) (bounds repere.RepereBounds, fxs []int, dfxs []float64) {
	f := expr.Closure(vari)

	fxs = make([]int, len(grid))
	dfxs = make([]float64, len(grid))
	// always add the origin
	minX, maxX, minY, maxY := -1., 1., -1., 1.
	for i, xValue := range grid {
		x := float64(xValue)
		y := f(x)
		fxs[i] = int(y)
		dfxs[i] = computeDf(f, x, 1e-5)

		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}

	bounds = boundsFromBoudingBox(minX, minY, maxX, maxY)
	return
}
