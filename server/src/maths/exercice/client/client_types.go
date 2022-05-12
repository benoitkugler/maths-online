package client

import (
	"encoding/json"

	functiongrapher "github.com/benoitkugler/maths-online/maths/function_grapher"
	"github.com/benoitkugler/maths-online/maths/repere"
)

//go:generate ../../../../../../structgen/structgen -source=client_types.go -mode=dart:../../../../../eleve/lib/exercices/types.gen.dart  -mode=itfs-json:gen_itfs_client.go

type Question struct {
	Title  string
	Enonce Enonce
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
func (FunctionGraphBlock) isBlock()  {}
func (TableBlock) isBlock()          {}

func (NumberFieldBlock) isBlock()           {}
func (ExpressionFieldBlock) isBlock()       {}
func (RadioFieldBlock) isBlock()            {}
func (DropDownFieldBlock) isBlock()         {}
func (OrderedListFieldBlock) isBlock()      {}
func (FigurePointFieldBlock) isBlock()      {}
func (FigureVectorFieldBlock) isBlock()     {}
func (FigureVectorPairFieldBlock) isBlock() {}
func (VariationTableFieldBlock) isBlock()   {}
func (FunctionPointsFieldBlock) isBlock()   {}
func (TreeFieldBlock) isBlock()             {}
func (TableFieldBlock) isBlock()            {}

// TextOrMath is a part of a text line, rendered
// either as plain text or using LaTeX in text mode.
type TextOrMath struct {
	Text   string
	IsMath bool
}

type TextBlock struct {
	Parts   []TextOrMath
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
	Label   string
	Columns []VariationColumnNumber
	// Arrows displays the arrows between two local extrema.
	Arrows []bool
}

// SignColumn is a column in a sign table.
type SignColumn struct {
	X                 string // as LaTeX code
	IsYForbiddenValue bool   // if true, a double bar is displayed
	IsSign            bool
	IsPositive        bool // for signs, displays a +, for numbers displays a 0 (else nothing)
}

type SignTableBlock struct {
	Label   string
	Columns []SignColumn
}

type FigureBlock struct {
	Figure repere.Figure `dart-extern:"repere.gen.dart"`
}

type FunctionGraphBlock struct {
	Graph functiongrapher.FunctionsGraph
}

type TableBlock struct {
	HorizontalHeaders []TextOrMath // optional
	VerticalHeaders   []TextOrMath // optional
	Values            [][]TextOrMath
}

// NumberFieldBlock is an answer field where only
// numbers are allowed
// answers are compared as float values
type NumberFieldBlock struct {
	ID int
}
type ExpressionFieldBlock struct {
	Label string // as LaTeX, optional
	ID    int
}

type ListFieldProposal struct {
	Content []TextOrMath
}

type RadioFieldBlock struct {
	Proposals []ListFieldProposal
	ID        int
}

// DropDownFieldBlock is the same has RadioFieldBlock,
// but is displayed inline.
type DropDownFieldBlock struct {
	Proposals []ListFieldProposal
	ID        int
}

type OrderedListFieldBlock struct {
	Label string // as LaTeX, optional, displayed before the answer
	// Proposals is a shuffled version of the list,
	// displayed as math text
	Proposals    []string
	AnswerLength int
	ID           int
}

// FigurePointFieldBlock asks for one 2D point
type FigurePointFieldBlock struct {
	Figure repere.Figure `dart-extern:"repere.gen.dart"`
	ID     int
}

// FigureVectorFieldBlock asks for a vector,
// represented by start and end.
// It may be used for vectors and affine functions

type FigureVectorFieldBlock struct {
	LineLabel string        // ignored if AsLine is false
	Figure    repere.Figure `dart-extern:"repere.gen.dart"`
	ID        int
	AsLine    bool
}

// FigureVectorPairFieldBlock asks for two vectors,
// represented by start and end, but evaluated
// as vector.
// The trivial case where the two pair of points are equals
// is not allowed

type FigureVectorPairFieldBlock struct {
	Figure repere.Figure `dart-extern:"repere.gen.dart"`
	ID     int
}

// VariationTableFieldBlock asks to complete a
// variation table (with fixed length)
type VariationTableFieldBlock struct {
	Label           string
	LengthProposals []int // propositions of the number of arrows
	ID              int
}

// FunctionPointsFieldBlock asks to place points
// to draw the graph of a function
type FunctionPointsFieldBlock struct {
	Label  string              // name of the function
	Xs     []int               // the grid
	Dfxs   []float64           // the derivatives of the function, to plot a nice curve
	Bounds repere.RepereBounds `dart-extern:"repere.gen.dart"`
	ID     int
}

// TreeShape defines the shape of a "regular" tree,
// specifying the number of children for each level
type TreeShape []int

// TreeFieldBlock asks to choose the shape and complete a
// probability tree
type TreeFieldBlock struct {
	ShapeProposals  []TreeShape
	EventsProposals []TextOrMath
	ID              int
}

type TableFieldBlock struct {
	HorizontalHeaders []TextOrMath
	VerticalHeaders   []TextOrMath
	ID                int
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
func (FunctionPointsAnswer) isAnswer()  {}
func (TreeAnswer) isAnswer()            {}
func (TableAnswer) isAnswer()           {}

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

type FunctionPointsAnswer struct {
	Fxs []int
}

type TreeNodeAnswer struct {
	Children      []TreeNodeAnswer
	Probabilities []float64 // edges, same length as Children
	Value         int       // index into the proposals
}

type TreeAnswer struct {
	Root TreeNodeAnswer
}

type TableAnswer struct {
	Rows [][]float64
}

// QuestionAnswersIn map the field ids to their answer
type QuestionAnswersIn struct {
	Data map[int]Answer
}

func (out *QuestionAnswersIn) UnmarshalJSON(src []byte) error {
	var wr struct {
		Data map[int]AnswerWrapper
	}

	err := json.Unmarshal(src, &wr)
	out.Data = make(map[int]Answer)
	for i, v := range wr.Data {
		out.Data[i] = v.Data
	}

	return err
}

func (out QuestionAnswersIn) MarshalJSON() ([]byte, error) {
	var tmp struct {
		Data map[int]AnswerWrapper
	}
	tmp.Data = make(map[int]AnswerWrapper)
	for k, v := range out.Data {
		tmp.Data[k] = AnswerWrapper{v}
	}
	return json.Marshal(tmp)
}

type QuestionAnswersOut struct {
	Data map[int]bool
}

// IsCorrect returns `true` if all the fields are correct.
func (qu QuestionAnswersOut) IsCorrect() bool {
	for _, v := range qu.Data {
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

func (out *QuestionSyntaxCheckIn) UnmarshalJSON(src []byte) error {
	var wr struct {
		Answer AnswerWrapper
		ID     int
	}
	err := json.Unmarshal(src, &wr)
	out.Answer = wr.Answer.Data
	out.ID = wr.ID

	return err
}

func (out QuestionSyntaxCheckIn) MarshalJSON() ([]byte, error) {
	wr := struct {
		Answer AnswerWrapper
		ID     int
	}{
		Answer: AnswerWrapper{out.Answer},
		ID:     out.ID,
	}
	return json.Marshal(wr)
}

type QuestionSyntaxCheckOut struct {
	Reason  string
	ID      int
	IsValid bool
}
