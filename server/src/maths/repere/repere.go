// Package repere implements a simple DSL for
// 2D mathematical figure drawing
package repere

//go:generate ../../../../../structgen/structgen -source=repere.go -mode=dart:../../../../eleve/lib/questions/repere.gen.dart

// Coord is a coordinate pair, in the usual mathematical plan,
// where X and Y must be between 0 and the dimension of the figure
type Coord struct {
	X, Y float64
}

// Color is an hex or ahex color string, with
// #FFFFFF or #AAFFFFFF format
type Color = string

type PointName = string

type RepereBounds struct {
	// Width and Height defines the logical size of
	// the figure. Since points comparison as performed
	// by rounding to integers it also influences the
	// tolerance allowed.
	Width, Height int
	// Origin defines the visual position of the mathematical origin (0;0),
	// counting from the bottom left of the figure.
	// All points are expressed in mathematical coordinates, meaning a point (x;y)
	// will be visually placed at Origin + (x; y)
	Origin Coord
}

type Area struct {
	Color  Color
	Points []PointName // polyline
}
type Drawings struct {
	Points   map[PointName]LabeledPoint
	Segments []Segment
	Lines    []Line
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
	Color Color
	Point PosPoint
}

// Segment is a segment from a defined point to another
type Segment struct {
	LabelName string // optional
	From, To  PointName
	Color     Color
	LabelPos  LabelPos    // used only if LabelName is not zero
	Kind      SegmentKind // what to actually draw
}

// Line is an infinite line, defined by an equation y = ax + b
type Line struct {
	Label string
	Color Color
	A, B  float64
}
