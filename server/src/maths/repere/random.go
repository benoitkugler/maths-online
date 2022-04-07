package repere

type RandomCoord struct {
	X, Y string // must be a valid expression.Expression
}

type RandomLabeledPoint struct {
	Coord RandomCoord
	Pos   LabelPos
}

type RandomLine struct {
	Label string
	A, B  string // must be a valid expression.Expression
}

type NamedRandomLabeledPoint struct {
	Name  PointName
	Point RandomLabeledPoint
}

type RandomDrawings struct {
	Points   []NamedRandomLabeledPoint
	Segments []Segment
	Lines    []RandomLine
}
