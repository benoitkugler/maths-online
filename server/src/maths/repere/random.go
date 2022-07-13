package repere

type RandomCoord struct {
	X, Y string // must be a valid expression.Expression
}

type RandomLabeledPoint struct {
	Color string // #FFFFFF format
	Coord RandomCoord
	Pos   LabelPos
}

type RandomSegment struct {
	LabelName string // optional
	From, To  string // expression.Expression resolving to point name
	Color     Color
	LabelPos  LabelPos    // used only if LabelName is not zero
	Kind      SegmentKind // what to actually draw
}

type RandomLine struct {
	Label string
	A, B  string // must be a valid expression.Expression
	Color string // #FFFFFF format
}

type RandomArea struct {
	Color  Color
	Points []string // expression.Expression for polyline point names
}

type NamedRandomLabeledPoint struct {
	Name  string // must be a valid expression.Expression
	Point RandomLabeledPoint
}

type RandomDrawings struct {
	Points   []NamedRandomLabeledPoint
	Segments []RandomSegment
	Lines    []RandomLine
	Areas    []RandomArea
}
