// Package functiongrapher provides a way to convert an
// arbitrary function expression into a list of quadratic bezier curve.
package functiongrapher

import (
	"math"

	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/repere"
)

type BezierCurve struct {
	P0, P1, P2 repere.Coord `dart-extern:"repere:repere.gen.dart"`
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
	if math.Abs(dFrom-dTo) < 0.1 {
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
func newSegment(fn expression.FunctionExpr, from, to float64) segment {
	f := fn.Closure()

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
	Decoration FunctionDecoration
	Segments   []BezierCurve
}

func newFunctionGraph(fn expression.FunctionDefinition) []BezierCurve {
	step := (fn.To - fn.From) / nbStep
	curves := make([]BezierCurve, nbStep)
	for i := range curves {
		seg := newSegment(fn.FunctionExpr, fn.From+float64(i)*step, fn.From+float64(i+1)*step)
		curves[i] = seg.toCurve()
	}

	return curves
}

type FunctionDecoration struct {
	Label string
	Color string
}

type FunctionsGraph struct {
	Functions []FunctionGraph
	Bounds    repere.RepereBounds `dart-extern:"repere:repere.gen.dart"`
}

// nbStep is the number of segments used when converting
// a function curve to Bezier curves.
const nbStep = 100

// Graph splits the curve of each function on the definition domain in small chunks on which it is approximated
// by Bezier curves.
// It will panic if the expression are not valid or evaluable given their input variable.
func NewFunctionGraph(functions []expression.FunctionDefinition, decorations []FunctionDecoration) FunctionsGraph {
	var allCurves []BezierCurve

	out := FunctionsGraph{
		Functions: make([]FunctionGraph, len(functions)),
	}
	for i, fn := range functions {
		out.Functions[i].Segments = newFunctionGraph(fn)
		out.Functions[i].Decoration = decorations[i]
		allCurves = append(allCurves, out.Functions[i].Segments...)
	}

	out.Bounds = boundingBox(allCurves)

	return out
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
func GraphFromVariations(dec FunctionDecoration, xs []float64, ys []float64) FunctionsGraph {
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

	return FunctionsGraph{
		Functions: []FunctionGraph{{Decoration: dec, Segments: curves}},
		Bounds:    boundingBox(curves),
	}
}

// BoundsFromExpression returns the f(x) and f'(x) values for x in the grid
func BoundsFromExpression(fn expression.FunctionExpr, grid []int) (bounds repere.RepereBounds, fxs []int, dfxs []float64) {
	f := fn.Closure()

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
