package client

import (
	"github.com/benoitkugler/maths-online/server/src/maths/expression/sets"
	"github.com/benoitkugler/maths-online/server/src/maths/functiongrapher"
	"github.com/benoitkugler/maths-online/server/src/maths/repere"
)

type Question struct {
	Enonce     Enonce
	Correction Enonce
}

type Enonce []Block

type Block interface {
	isBlock()
}

func (TextBlock) isBlock()           {}
func (FormulaBlock) isBlock()        {}
func (VariationTableBlock) isBlock() {}
func (SignTableBlock) isBlock()      {}
func (FigureBlock) isBlock()         {}
func (FunctionsGraphBlock) isBlock() {}
func (TableBlock) isBlock()          {}
func (TreeBlock) isBlock()           {}

func (NumberFieldBlock) isBlock()                {}
func (ExpressionFieldBlock) isBlock()            {}
func (RadioFieldBlock) isBlock()                 {}
func (DropDownFieldBlock) isBlock()              {}
func (OrderedListFieldBlock) isBlock()           {}
func (GeometricConstructionFieldBlock) isBlock() {}
func (VariationTableFieldBlock) isBlock()        {}
func (SignTableFieldBlock) isBlock()             {}
func (FunctionPointsFieldBlock) isBlock()        {}
func (TreeFieldBlock) isBlock()                  {}
func (TableFieldBlock) isBlock()                 {}
func (VectorFieldBlock) isBlock()                {}
func (ProofFieldBlock) isBlock()                 {}
func (SetFieldBlock) isBlock()                   {}

// TextOrMath is a part of a text line, rendered
// either as plain text or using LaTeX in text mode.
type TextOrMath struct {
	Text   string
	IsMath bool
}

type TextBlock struct {
	Parts   TextLine
	Bold    bool
	Italic  bool
	Smaller bool
}

// FormulaBlock is whole line, rendered as LaTeX in display mode
type FormulaBlock struct {
	Formula string // as latex
}

// VariationColumnNumber is a column in a variation table
// displaying (x, f(x)) values
type VariationColumnNumber struct {
	X, Y string // LaTeX
	IsUp bool   // to adjust the vertical alignment
}

type VariationTableBlock struct {
	Label   string // LaTeX
	Columns []VariationColumnNumber
	// Arrows displays the arrows between two local extrema,
	// with the convention that `true` means `isUp`.
	Arrows []bool
}

type FunctionSign struct {
	Label     string       // printed in math mode
	FxSymbols []SignSymbol // one for each X, alternate with [Signs]
	Signs     []bool       // is positive, with length len(Xs) - 1
}

type SignTableBlock struct {
	Xs        []string // as LaTeX code, includes empty cells
	Functions []FunctionSign
}

type FigureBlock struct {
	Figure repere.Figure
}

type FunctionArea struct {
	Color repere.ColorHex
	Path  []functiongrapher.BezierCurve
}

type FunctionPoint struct {
	Color  repere.ColorHex
	Legend string // LaTeX code
	Coord  repere.Coord
}
type FunctionsGraphBlock struct {
	Functions []functiongrapher.FunctionGraph
	Sequences []functiongrapher.SequenceGraph
	Areas     []FunctionArea
	Points    []FunctionPoint
	Bounds    repere.RepereBounds
}

type TableBlock struct {
	HorizontalHeaders []TextOrMath // optional
	VerticalHeaders   []TextOrMath // optional
	Values            [][]TextOrMath
}

// SizeHint is the length of the expected answer,
// in runes. It may be used by the client to adjust the field width.
type SizeHint = int

// NumberFieldBlock is an answer field where only
// numbers are allowed
// answers are compared as float values
type NumberFieldBlock struct {
	ID int
	// Typical values range from 1 to 15
	SizeHint SizeHint
}
type ExpressionFieldBlock struct {
	Label  string // as LaTeX, optional
	Suffix string // as LaTeX, optional

	// Typical values range from 1 to 30
	SizeHint SizeHint

	// If true, the field is diplayed with two subfields
	ShowFractionHelp bool

	ID int
}

// TextLine is the general form of a static chunk of text,
// alternating LaTeX or basic text
type TextLine []TextOrMath

type RadioFieldBlock struct {
	Proposals []TextLine
	ID        int
}

// DropDownFieldBlock is the same has RadioFieldBlock,
// but is displayed inline.
type DropDownFieldBlock struct {
	Proposals []TextLine
	ID        int
}

type OrderedListFieldBlock struct {
	Label string // as LaTeX, optional, displayed before the answer
	// Proposals is a shuffled version of the list
	Proposals    []TextLine
	AnswerLength int
	ID           int
}

type GeometricConstructionFieldBlock struct {
	ID         int
	Field      GeoField
	Background FigureOrGraph
}

type FigureOrGraph interface {
	isFigureOrGraph()
	FigBounds() repere.RepereBounds
}

func (FigureBlock) isFigureOrGraph()         {}
func (FunctionsGraphBlock) isFigureOrGraph() {}

func (fg FigureBlock) FigBounds() repere.RepereBounds         { return fg.Figure.Bounds }
func (fg FunctionsGraphBlock) FigBounds() repere.RepereBounds { return fg.Bounds }

type GeoField interface {
	isGF()
}

func (GFPoint) isGF()      {}
func (GFVector) isGF()     {}
func (GFVectorPair) isGF() {}

// FigurePointFieldBlock asks for one 2D point
type GFPoint struct{}

// FigureVectorFieldBlock asks for a vector,
// represented by start and end.
// It may be used for vectors and affine functions
type GFVector struct {
	LineLabel string // ignored if AsLine is false
	AsLine    bool
}

// FigureVectorPairFieldBlock asks for two vectors,
// represented by start and end, but evaluated
// as vector.
// The trivial case where the two pair of points are equals
// is not allowed
type GFVectorPair struct{}

// VariationTableFieldBlock asks to complete a
// variation table (with fixed length)
type VariationTableFieldBlock struct {
	Label           string // LaTeX code
	LengthProposals []int  // propositions of the number of arrows
	ID              int
}

// SignTableFieldBlock asks to complete a
// sign table (with fixed length)
type SignTableFieldBlock struct {
	LengthProposals []int    // propositions of the number of signs to fill
	Labels          []string // LaTeX code, for each function
	ID              int
}

// FunctionPointsFieldBlock asks to place points
// to draw the graph of a function
type FunctionPointsFieldBlock struct {
	IsDiscrete bool   // true for sequences, removing the curve between points
	Label      string // name of the function
	Xs         []int  // the grid
	// the derivatives of the function, to plot a nice curve
	// empty if [IsDiscrete] is true
	Dfxs   []float64
	Bounds repere.RepereBounds
	ID     int
}

type TreeBlock struct {
	EventsProposals []TextLine
	Root            TreeNodeAnswer
}

// TreeShape defines the shape of a "regular" tree,
// specifying the number of children for each level
type TreeShape []int

// TreeFieldBlock asks to choose the shape and complete a
// probability tree
type TreeFieldBlock struct {
	ShapeProposals  []TreeShape
	EventsProposals []TextLine
	ID              int
}

type TableFieldBlock struct {
	HorizontalHeaders []TextOrMath
	VerticalHeaders   []TextOrMath
	ID                int
}

type VectorFieldBlock struct {
	ID            int
	DisplayColumn bool
	// Typical values range from 1 to 15
	SizeHintX, SizeHintY SizeHint
}

// Statement is a basic statement
type Statement struct {
	Content TextLine
}

// Equality is an equality of the form A1 = A2 = A3
type Equality struct {
	Terms   []TextLine
	Def     TextLine // Optional avec x = 2k term
	WithDef bool
}

// Node is an higher level assertion, such as
// (m is even) AND (n is odd)
type Node struct {
	Left, Right Assertion
	Op          Binary
}

// Sequence is a list of elementary steps needed
// to write a mathematical proof, where each step are
// implicitely connected by a "So" (Donc) connector.
type Sequence struct {
	Parts Assertions
}

type ProofFieldBlock struct {
	Shape         Proof
	TermProposals []TextLine
	ID            int
}

// SetFieldBlock asks the student to build a set expression,
// using the given [Sets] and math set operators
type SetFieldBlock struct {
	Sets []string
	ID   int
}

// Answer is a sum type for the possible answers
// of question fields
type Answer interface {
	isAnswer()
}

func (NumberAnswer) isAnswer()          {}
func (RadioAnswer) isAnswer()           {}
func (ExpressionAnswer) isAnswer()      {}
func (OrderedListAnswer) isAnswer()     {}
func (PointAnswer) isAnswer()           {}
func (DoublePointAnswer) isAnswer()     {}
func (DoublePointPairAnswer) isAnswer() {}
func (VariationTableAnswer) isAnswer()  {}
func (SignTableAnswer) isAnswer()       {}
func (FunctionPointsAnswer) isAnswer()  {}
func (TreeAnswer) isAnswer()            {}
func (TableAnswer) isAnswer()           {}
func (VectorNumberAnswer) isAnswer()    {}
func (ProofAnswer) isAnswer()           {}
func (SetAnswer) isAnswer()             {}

// NumberAnswer is compared with float equality, with a fixed
// precision of 8 digits
type NumberAnswer struct {
	Value float64
}

type ExpressionAnswer struct {
	Expression string
}

// RadioAnswer is compared against a reference index
// It is shared by Radio and DropDown fields.
type RadioAnswer struct {
	Index int
}

type OrderedListAnswer struct {
	Indices []int // indices into the question field proposals
}

// PointAnswer is a 2D point, whoose coordinates
// are rounded before begin compared
type PointAnswer struct {
	Point repere.IntCoord
}

type DoublePointAnswer struct {
	From repere.IntCoord
	To   repere.IntCoord
}

type DoublePointPairAnswer struct {
	From1 repere.IntCoord
	To1   repere.IntCoord
	From2 repere.IntCoord
	To2   repere.IntCoord
}

type VariationTableAnswer struct {
	Xs     []string // expressions
	Fxs    []string // expressions
	Arrows []bool   // isUp
}

type SignTableAnswer struct {
	Xs        []string       // expressions
	Functions []FunctionSign // each label is ignored
}

type FunctionPointsAnswer struct {
	Fxs []int
}

type TreeNodeAnswer struct {
	Children      []TreeNodeAnswer
	Probabilities []string // expression for edges, same length as Children
	Value         int      // index into the proposals, ignored for the root
}

type TreeAnswer struct {
	Root TreeNodeAnswer
}

type TableAnswer struct {
	Rows [][]float64
}

type VectorNumberAnswer struct {
	X, Y float64
}

type ProofAnswer struct {
	Proof Proof
}

type SetAnswer struct {
	// Sets are relative to the question [Sets] list.
	Root sets.ListNode
}

// QuestionAnswersIn map the field ids to their answer
type QuestionAnswersIn struct {
	Data Answers
}

type Answers map[int]Answer

type QuestionAnswersOut struct {
	Results         map[int]bool
	ExpectedAnswers Answers
}

// IsCorrect returns `true` if all the fields are correct.
func (qu QuestionAnswersOut) IsCorrect() bool {
	for _, v := range qu.Results {
		if !v {
			return false
		}
	}
	return true
}

// QuestionSyntaxCheckIn is emitted by the client
// to perform a preliminary check of the syntax,
// without validating the answer
type QuestionSyntaxCheckIn struct {
	Answer Answer
	ID     int
}

type QuestionSyntaxCheckOut struct {
	Reason  string
	ID      int
	IsValid bool
}
