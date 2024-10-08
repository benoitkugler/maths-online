package repere

type RandomCoord struct {
	X, Y string // must be a valid expression.Expression
}

type RandomLabeledPoint struct {
	Color ColorHex // #FFFFFF format
	Coord RandomCoord
	Pos   LabelPos
}

type RandomSegment struct {
	LabelName string // optional, support interpolation
	From, To  string // expression.Expression resolving to point name
	Color     ColorHex
	LabelPos  LabelPos    // used only if LabelName is not zero
	Kind      SegmentKind // what to actually draw
}

type RandomLine struct {
	Label string
	A, B  string // must be a valid expression.Expression
	Color ColorHex
}

type RandomArea struct {
	Color  ColorHex
	Points []string // expression.Expression for polyline point names
}

type RandomCircle struct {
	Center    RandomCoord
	Radius    string   // must be a valid expression.Expression
	LineColor ColorHex // optional, default to black
	FillColor ColorHex // optional, defaul to transparent
	Legend    string   // support interpolation
}

type NamedRandomLabeledPoint struct {
	Name  string // must be a valid expression.Expression
	Point RandomLabeledPoint
}

type RandomDrawings struct {
	Points   []NamedRandomLabeledPoint
	Segments []RandomSegment
	Lines    []RandomLine
	Circles  []RandomCircle
	Areas    []RandomArea
}
