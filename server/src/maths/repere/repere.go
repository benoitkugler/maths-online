// Package repere implements a simple DSL for
// 2D mathematical figure drawing
package repere

//go:generate ../../../../../structgen/structgen -source=repere.go -mode=dart:../../../../eleve/lib/exercices/repere.gen.dart

// Coord is a coordinate pair, in the usual mathematical plan,
// where X and Y are expressed in a fraction (1/1000) of the total square.
type Coord struct {
	X, Y int
}

type PointName = string

type LabeledPoint struct {
	Point Coord
	Pos   LabelPos
}

type Figure struct {
	Points map[PointName]LabeledPoint
	Lines  []Line
	// Width and Height are values between 0 and 1000
	// restricting the displayed area
	Width, Height int
}

type Line struct {
	LabelName string // optional
	From, To  PointName
	LabelPos  LabelPos // used only if LabelName is not zero
}

// ProjeteOrtho compute the coordinates of the orthogonal
// projection of B on (AC).
func ProjeteOrtho(B, A, C Coord) Coord {
	u := C.X - A.X // AC
	v := C.Y - A.Y // AC
	// det(AB, AC)
	abX := B.X - A.X
	abY := B.Y - A.Y
	d := (abX*v - abY*u)
	// solve for BH = (x, y)
	// xu + yv = 0
	// xv - yu = -d
	x := -d * v / (u*u + v*v)
	y := d * u / (u*u + v*v)
	return Coord{X: (x + B.X), Y: (y + B.Y)}
}
