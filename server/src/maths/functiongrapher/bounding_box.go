package functiongrapher

import (
	"math"

	"github.com/benoitkugler/maths-online/server/src/maths/repere"
)

// quadratic polinomial
// x = At^2 + Bt + C
// where
// A = p0 + p2 - 2p1
// B = 2(p1 - p0)
// C = p0
func bezierQuad(p0, p1, p2, t float64) float64 {
	return (p0+p2-2*p1)*t*t + 2*(p1-p0)*t + p0
}

// derivative as at + b where a,b :
func quadraticDerivative(p0, p1, p2 float64) (a, b float64) {
	return 2 * (p2 - p1 - (p1 - p0)), 2 * (p1 - p0)
}

// handle the case where a = 0
func linearRoots(a, b float64) []float64 {
	if a == 0 {
		return nil
	}
	return []float64{-b / a}
}

func (cu BezierCurve) criticalPoints() (tX, tY []float64) {
	p0x, p0y := cu.P0.X, cu.P0.Y
	p1x, p1y := cu.P1.X, cu.P1.Y
	p2x, p2y := cu.P2.X, cu.P2.Y

	aX, bX := quadraticDerivative(p0x, p1x, p2x)
	aY, bY := quadraticDerivative(p0y, p1y, p2y)

	return linearRoots(aX, bX), linearRoots(aY, bY)
}

func (cu BezierCurve) evaluateCurve(t float64) (x, y float64) {
	p0x, p0y := cu.P0.X, cu.P0.Y
	p1x, p1y := cu.P1.X, cu.P1.Y
	p2x, p2y := cu.P2.X, cu.P2.Y
	return bezierQuad(p0x, p1x, p2x, t), bezierQuad(p0y, p1y, p2y, t)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func (cu BezierCurve) boundingBox() (minX, maxX, minY, maxY float64) {
	resX, resY := cu.criticalPoints()

	// draw min and max
	var bbox []repere.Coord

	// add begin and end point
	for _, t := range append(append(resX, 0, 1), resY...) {
		// filter invalid value
		if !(0 <= t && t <= 1) {
			continue
		}
		x, y := cu.evaluateCurve(t)

		bbox = append(bbox, repere.Coord{X: x, Y: y})
	}

	// bbox is never empty since it always contains the begin and end point
	minX = bbox[0].X
	minY = bbox[0].Y
	maxX, maxY = minX, minY

	for _, e := range bbox[1:] {
		minX = min(e.X, minX)
		minY = min(e.Y, minY)
		maxX = max(e.X, maxX)
		maxY = max(e.Y, maxY)
	}
	return
}

// BoundingBox compute the bounds enclosing the given segments,
// always adding the origin and some padding.
func BoundingBox(curves []BezierCurve) repere.RepereBounds {
	minX, maxX, minY, maxY := -1., 1., -1., 1.
	for _, curve := range curves {
		minX2, maxX2, minY2, maxY2 := curve.boundingBox()
		minX = min(minX, minX2)
		minY = min(minY, minY2)
		maxX = max(maxX, maxX2)
		maxY = max(maxY, maxY2)
	}

	return boundsFromBoudingBox(minX, minY, maxX, maxY)
}

// the client display exactly the width and height asked for:
// add a bit of padding
func boundsFromBoudingBox(minX, minY, maxX, maxY float64) repere.RepereBounds {
	const paddingX = 1
	const paddingY = 1
	return repere.RepereBounds{
		Width:  int(math.Ceil(maxX-minX)) + 2*paddingX,
		Height: int(math.Ceil(maxY-minY)) + 2*paddingY,
		Origin: repere.Coord{
			X: -minX + paddingX,
			Y: -minY + paddingY,
		},
	}
}
