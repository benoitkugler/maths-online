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

type RandomDrawings struct {
	Points   map[PointName]RandomLabeledPoint
	Segments []Segment
	Lines    []RandomLine
}
