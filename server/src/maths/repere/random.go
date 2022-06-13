package repere

type RandomCoord struct {
	X, Y string // must be a valid expression.Expression
}

type RandomLabeledPoint struct {
	Color string // #FFFFFF format
	Coord RandomCoord
	Pos   LabelPos
}

type RandomLine struct {
	Label string
	A, B  string // must be a valid expression.Expression
	Color string // #FFFFFF format
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
