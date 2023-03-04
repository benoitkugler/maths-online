// Package repere implements a simple DSL for
// 2D mathematical figure drawing
package repere

import (
	"math"
)

// Coord is a coordinate pair, in the usual mathematical plan,
// where X and Y must be between 0 and the dimension of the figure
type Coord struct {
	X, Y float64
}

type PointName = string

type RepereBounds struct {
	// Width and Height defines the logical size of
	// the figure. Since points comparison are performed
	// by rounding to integers it also influences the
	// tolerance allowed.
	Width, Height int
	// Origin defines the visual position of the mathematical origin (0;0),
	// counting from the bottom left of the figure.
	// All points are expressed in mathematical coordinates, meaning a point (x;y)
	// will be visually placed at Origin + (x;y)
	Origin Coord
}

type Area struct {
	Color  ColorHex
	Points []PointName // polyline, refering to the Drawings.Points
}

type Circle struct {
	Center    Coord
	Radius    float64
	LineColor ColorHex // optional, default to black
	FillColor ColorHex // optional, defaul to transparent
	Legend    string   // LaTeX
}

type Drawings struct {
	Points   map[PointName]LabeledPoint
	Segments []Segment
	Lines    []Line
	Circles  []Circle
	Areas    []Area
}

type Figure struct {
	Drawings Drawings

	Bounds RepereBounds

	ShowGrid   bool
	ShowOrigin bool
}

type PosPoint struct {
	Point Coord
	Pos   LabelPos
}

type LabeledPoint struct {
	Color ColorHex
	Point PosPoint
}

// Segment is a segment from a defined point to another
type Segment struct {
	LabelName string // optional
	From, To  PointName
	Color     ColorHex
	LabelPos  LabelPos    // used only if LabelName is not zero
	Kind      SegmentKind // what to actually draw
}

// InferLine returns the (infinite) affine line passing by the segment [from, to]
func InferLine(from, to Coord) (a, b float64) {
	if to.X == from.X {
		a = math.Inf(1)
		b = to.X
	} else {
		a = (to.Y - from.Y) / (to.X - from.X)
		b = from.Y - a*from.X
	}
	return a, b
}

// Line is an infinite line, defined by an equation y = ax + b
type Line struct {
	Label string
	Color ColorHex
	A, B  float64
}

// Bounds return the extremal points which should be seen on the given repere.
func (li Line) Bounds(rep RepereBounds) (start, end Coord) {
	fw, fh := float64(rep.Width), float64(rep.Height)
	or := rep.Origin
	if math.IsInf(li.A, 0) {
		// start point
		start = Coord{li.B, -or.Y}

		// end point
		end = Coord{li.B, fh - or.Y}
	} else {
		// start point
		start = Coord{-or.X, li.A*(-or.X) + li.B}

		// end point
		end = Coord{fw - or.X, li.A*(fw-or.X) + li.B}
	}

	start, end, _ = clipLine(-or.X, fw-or.X, -or.Y, fh-or.Y, start, end)

	return
}

// Define the x/y clipping values for the border.
// Define the start and end points of the line.
func clipLine(edgeLeft, edgeRight, edgeBottom, edgeTop float64,
	src0, src1 Coord,
) (clipStart, clipEnd Coord, draw bool) {
	// Liang-Barsky function by Daniel White @ http://www.skytopia.com/project/articles/compsci/clipping.html

	t0, t1 := 0., 1.
	xdelta := src1.X - src0.X
	ydelta := src1.Y - src0.Y
	var p, q, r float64

	for edge := 0; edge < 4; edge++ { // Traverse through left, right, bottom, top edges.
		if edge == 0 {
			p = -xdelta
			q = -(edgeLeft - src0.X)
		}
		if edge == 1 {
			p = xdelta
			q = (edgeRight - src0.X)
		}
		if edge == 2 {
			p = -ydelta
			q = -(edgeBottom - src0.Y)
		}
		if edge == 3 {
			p = ydelta
			q = (edgeTop - src0.Y)
		}
		r = q / p
		if p == 0 && q < 0 {
			return Coord{}, Coord{}, false
		} // Don't draw line at all. (parallel line outside)

		if p < 0 {
			if r > t1 {
				return Coord{}, Coord{}, false
				// Don't draw line at all.
			} else if r > t0 {
				t0 = r
			} // Line is clipped!
		} else if p > 0 {
			if r < t0 {
				return Coord{}, Coord{}, false // Don't draw line at all.
			} else if r < t1 {
				t1 = r
			} // Line is clipped!
		}
	}

	clipStart = Coord{X: src0.X + t0*xdelta, Y: src0.Y + t0*ydelta}
	clipEnd = Coord{X: src0.X + t1*xdelta, Y: src0.Y + t1*ydelta}
	draw = true // (clipped) line is drawn
	return
}
