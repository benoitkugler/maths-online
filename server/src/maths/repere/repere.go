// Package repere implements a simple DSL for
// 2D mathematical figure drawing
package repere

import "math"

//go:generate ../../../../../structgen/structgen -source=repere.go -mode=dart:../../../../eleve/lib/exercices/repere.gen.dart

// Coord is a coordinate pair, in the usual mathematical plan,
// where X and Y must be between 0 and the dimension of the figure
type Coord struct {
	X, Y float64
}

type PointName = string

type LabeledPoint struct {
	Point Coord
	Pos   LabelPos
}

type Figure struct {
	Points   map[PointName]LabeledPoint
	Segments []Segment
	Lines    []Line

	// Width and Height defines the logical size of
	// the figure. Since points comparison as performed
	// by rounding to integers it also influences the
	// tolerance allowed.
	Width, Height int
	// Origin defines the visual position of the mathematical origin (0;0),
	// counting from the bottom left of the figure.
	// All points are expressed in mathematical coordinates, meaning a point (x;y)
	// will be visually placed at Origin + (x; y)
	Origin   Coord
	ShowGrid bool
}

// Segment is a segment from a defined point to another
type Segment struct {
	LabelName string // optional
	From, To  PointName
	LabelPos  LabelPos // used only if LabelName is not zero
	AsVector  bool     // if true, add an arrow
}

// Line is an infinite line, defined by an equation y = ax + b
type Line struct {
	Label string
	A, B  float64
}

// OrthogonalProjection compute the coordinates of the orthogonal
// projection of B on (AC).
func OrthogonalProjection(B, A, C Coord) Coord {
	u := C.X - A.X // AC
	v := C.Y - A.Y // AC
	// det(AB, AC)
	abX := B.X - A.X
	abY := B.Y - A.Y
	d := abX*v - abY*u
	// solve for BH = (x, y)
	// xu + yv = 0
	// xv - yu = -d
	norm := u*u + v*v
	x := math.Round(-d * v / norm)
	y := math.Round(d * u / norm)
	return Coord{X: x + B.X, Y: y + B.Y}
}
