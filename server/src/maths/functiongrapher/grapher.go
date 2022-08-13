// Package functiongrapher provides a way to convert an
// arbitrary function expression into a list of quadratic bezier curve.
package functiongrapher

import (
	"math"

	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/repere"
)

// BezierCurve is a quadratic Bezier curve with the
// additional invariant that P0.X <= P2.X
type BezierCurve struct {
	P0, P1, P2 repere.Coord `gomacro-extern:"repere:dart:repere.gen.dart"`
}

func (seg segment) toCurve() BezierCurve {
	p1 := controlFromDerivatives(seg.from, seg.to, seg.dFrom, seg.dTo)
	return BezierCurve{
		P0: seg.from,
		P1: p1,
		P2: seg.to,
	}
}

func middle(from, to repere.Coord) repere.Coord {
	return repere.Coord{X: (from.X + to.X) / 2, Y: (from.Y + to.Y) / 2}
}

// compute the control point matching the given derivatives,
// which is the intersection between
// the tangents at from and to
func controlFromDerivatives(from, to repere.Coord, dFrom, dTo float64) repere.Coord {
	// special case when df1 = df2
	if math.Abs(dFrom-dTo) < 0.1 {
		return middle(from, to)
	}

	xIntersec := (to.Y - from.Y + dFrom*from.X - dTo*to.X) / (dFrom - dTo)
	yIntersec := dFrom*(xIntersec-from.X) + from.Y
	return repere.Coord{X: xIntersec, Y: yIntersec}
}

// compute derivative with finite differences
func computeDf(f func(float64) float64, x, epsilon float64) float64 {
	return (f(x+epsilon) - f(x)) / epsilon
}

// invariant: from.X <= to.X
type segment struct {
	from, to   repere.Coord
	dFrom, dTo float64
}

// expr must be an expression containing only the variable `variable`
// assume from <= to
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

func NewFunctionGraph(fn expression.FunctionDefinition) []BezierCurve {
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

// nbStep is the number of segments used when converting
// a function curve to Bezier curves.
const nbStep = 100

// choose appropriate derivatives and call controlFromDerivatives
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

// NewFunctionGraphFromVariations builds one possible representation for the given variations
func NewFunctionGraphFromVariations(xs []float64, ys []float64) []BezierCurve {
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

	return curves
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

// HorizontalAxis is a convinience function returning the
// curve defined by an horizontal axis at `y`, between x = `from` and x = `to`
func HorizontalAxis(from, to, y float64) []BezierCurve {
	p0 := repere.Coord{X: from, Y: y}
	p2 := repere.Coord{X: to, Y: y}
	return []BezierCurve{
		{P0: p0, P1: middle(p0, p2), P2: p2},
	}
}

func copy(l []BezierCurve) []BezierCurve {
	return append([]BezierCurve(nil), l...)
}

// NewAreaBetween returns a closed path delimiting an area between two curves
// `top` and `bottom` restricted in [`left`, `right`], which must be included in
// the two curves.
func NewAreaBetween(top, bottom []BezierCurve, left, right float64) []BezierCurve {
	top = restrictRight(restrictLeft(copy(top), left), right)
	bottom = restrictRight(restrictLeft(copy(bottom), left), right)
	// merge the two curves to construct a path top -> topToBottom -> bottom -> bottomToTop

	lastTop, lastBottom := top[len(top)-1], bottom[len(bottom)-1]
	out := append(top, BezierCurve{P0: lastTop.P2, P2: lastBottom.P2, P1: middle(lastTop.P2, lastBottom.P2)})
	// add the bottom part reversed
	for i := range bottom {
		c := bottom[len(bottom)-1-i]
		c = BezierCurve{P0: c.P2, P1: c.P1, P2: c.P0} // thanks to math property
		out = append(out, c)
	}
	// finaly close the path
	firstTop, firstBottom := top[0], bottom[0]
	out = append(out, BezierCurve{P0: firstBottom.P0, P2: firstTop.P0, P1: middle(firstBottom.P0, firstTop.P0)})
	return out
}

func restrictLeft(curves []BezierCurve, left float64) []BezierCurve {
	// find where to cut left
	for i, c := range curves {
		if c.P0.X == left {
			return curves[i:]
		} else if c.P2.X == left {
			return curves[i+1:]
		} else if c.P0.X <= left && left <= c.P2.X { // split to adjust
			_, c = c.splitByX(left)
			curves[i] = c
			return curves[i:]
		}
	}
	panic("invalid left cut")
}

func restrictRight(curves []BezierCurve, right float64) []BezierCurve {
	// find where to cut left
	for i, c := range curves {
		if c.P0.X == right {
			return curves[:i]
		} else if c.P2.X == right {
			return curves[:i+1]
		} else if c.P0.X <= right && right <= c.P2.X { // split to adjust
			c, _ = c.splitByX(right)
			curves[i] = c
			return curves[:i+1]
		}
	}
	panic("invalid left cut")
}

// assume x is inside the curve and split at `x`,
// returning two pieces whose concatenation is equivalent to the original curve
func (bc BezierCurve) splitByX(x float64) (left, right BezierCurve) {
	left, right = bc, bc

	ts := bc.invertX(x)
	xs, ys := bc.evaluateCurve(ts)
	left.P2 = repere.Coord{X: xs, Y: ys}
	right.P0 = repere.Coord{X: xs, Y: ys}

	// adjust the control points to respect the derivatives
	leftDFrom, leftDTo := bc.derivative(0), bc.derivative(ts)
	left.P1 = controlFromDerivatives(left.P0, left.P2, leftDFrom, leftDTo)

	rightDFrom, rightDTo := bc.derivative(ts), bc.derivative(1)
	right.P1 = controlFromDerivatives(right.P0, right.P2, rightDFrom, rightDTo)

	return left, right
}

// assume P0.X <= x <= P2.X and returbs the unique t in [0, 1]
// such that B(t).X = x
func (bc BezierCurve) invertX(x float64) (t float64) {
	// a quad bezier is also a quadratic polynomial
	// x = At^2 + Bt + C
	// where
	// A = p0 + p2 - 2p1
	// B = 2(p1 - p0)
	// C = p0
	p0, p1, p2 := bc.P0.X, bc.P1.X, bc.P2.X
	a, b, c := p0+p2-2*p1, 2*(p1-p0), p0-x

	if a == 0 {
		return -c / b
	}

	delta := math.Sqrt(b*b - 4*a*c)

	t1 := (-b + delta) / (2 * a)
	t2 := (-b - delta) / (2 * a)

	if t1 >= 0 && t1 <= 1 {
		return t1
	}
	return t2
}

// returns the derivative a point t in [0,1], that is y'(t) / x'(t)
func (bc BezierCurve) derivative(t float64) float64 {
	p0x, p1x, p2x := bc.P0.X, bc.P1.X, bc.P2.X
	p0y, p1y, p2y := bc.P0.Y, bc.P1.Y, bc.P2.Y

	dx := 2*(p2x-p1x-(p1x-p0x))*t + 2*(p1x-p0x)
	dy := 2*(p2y-p1y-(p1y-p0y))*t + 2*(p1y-p0y)
	return dy / dx
}
