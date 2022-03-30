package exercice

import (
	"fmt"
	"math/rand"

	"github.com/benoitkugler/maths-online/maths/exercice/client"
	"github.com/benoitkugler/maths-online/maths/expression"
	functiongrapher "github.com/benoitkugler/maths-online/maths/function_grapher"
	"github.com/benoitkugler/maths-online/maths/repere"
)

// InvalidFieldAnswer is returned for syntactically incorrect answers
type InvalidFieldAnswer struct {
	ID     int
	Reason string
}

func (ifa InvalidFieldAnswer) Error() string {
	return fmt.Sprintf("field %d: %s", ifa.ID, ifa.Reason)
}

// fieldInstance is an answer field, identified with an integer ID
type fieldInstance interface {
	blockInstance

	fieldID() int

	// validateAnswerSyntax is called during editing for complex fields,
	// to catch syntax mistake before validating the answer
	// an error may also be returned against malicious query
	validateAnswerSyntax(answer client.Answer) error

	// evaluateAnswer evaluate the given answer against the reference
	// validateAnswerSyntax is assumed to have already been called on `answer`
	// so that is has a valid format.
	evaluateAnswer(answer client.Answer) (isCorrect bool)
}

var (
	_ fieldInstance = NumberFieldInstance{}
	_ fieldInstance = ExpressionFieldInstance{}
	_ fieldInstance = RadioFieldInstance{}
	_ fieldInstance = DropDownFieldInstance{}
	_ fieldInstance = OrderedListFieldInstance{}
	_ fieldInstance = FigurePointFieldInstance{}
	_ fieldInstance = FigureVectorFieldInstance{}
	_ fieldInstance = VariationTableFieldInstance{}
	_ fieldInstance = FunctionPointsFieldInstance{}
	// TODO: à tester
	_ fieldInstance = FigureVectorPairFieldInstance{}
	_ fieldInstance = FigureAffineLineFieldInstance{}
)

// NumberFieldInstance is an answer field where only
// numbers are allowed
// answers are compared as float values
type NumberFieldInstance struct {
	ID     int
	Answer float64 // expected answer
}

func (f NumberFieldInstance) fieldID() int { return f.ID }

func (f NumberFieldInstance) toClient() client.Block { return client.NumberFieldBlock{ID: f.ID} }

func (f NumberFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	_, ok := answer.(client.NumberAnswer)
	if !ok {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("expected NumberAnswer, got %T", answer),
		}
	}
	return nil
}

func (f NumberFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	return f.Answer == answer.(client.NumberAnswer).Value
}

// ExpressionFieldInstance is an answer field where a single mathematical expression
// if expected
type ExpressionFieldInstance struct {
	// if not empty, the field is displayed on a new line
	// expression are added an equal symbol : <expression> =
	Label StringOrExpression

	Answer          *expression.Expression
	ComparisonLevel expression.ComparisonLevel
	ID              int
}

func (f ExpressionFieldInstance) fieldID() int { return f.ID }

func (f ExpressionFieldInstance) toClient() client.Block {
	var label string
	if f.Label.Expression != nil {
		label = f.Label.Expression.AsLaTeX(nil) + " = "
	} else {
		label = f.Label.String
	}
	return client.ExpressionFieldBlock{
		ID:    f.ID,
		Label: label,
	}
}

func (f ExpressionFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	expr, ok := answer.(client.ExpressionAnswer)
	if !ok {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("expected ExpressionAnswer, got %T", answer),
		}
	}

	_, _, err := expression.Parse(expr.Expression)
	if err != nil {
		err := err.(expression.InvalidExpr)
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf(`Expression invalide : %s (à "%s")`, err.Reason, err.PortionOf(expr.Expression)),
		}
	}
	return nil
}

func (f ExpressionFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	expr, _, _ := expression.Parse(answer.(client.ExpressionAnswer).Expression)

	return expression.AreExpressionsEquivalent(f.Answer, expr, f.ComparisonLevel)
}

// type expressionOrText struct {
// 	Expression *expression.Expression
// 	Text       string
// 	IsMath     bool
// }

// func (e expressionOrText) instantiate() (out client.TextOrMath) {
// 	if e.Expression != nil {
// 		out.Content = e.Expression.AsLaTeX(nil)
// 		out.IsMath = true
// 	} else {
// 		out.Content = e.Text
// 		out.IsMath = e.IsMath
// 	}
// 	return out
// }

// type listFieldProposal struct {
// 	Content []expressionOrText
// }

// func (lf listFieldProposal) toClient() client.ListFieldProposal {
// 	out := client.ListFieldProposal{Content: make([]client.ListFieldProposalBlock, len(lf.Content))}
// 	for i, f := range lf.Content {
// 		out.Content[i] = f.toClient()
// 	}
// 	return out
// }

// RadioFieldInstance is an answer field where one choice
// is to be made against a fixed list
type RadioFieldInstance struct {
	Proposals []client.ListFieldProposal
	ID        int
	Answer    int // index into Proposals
}

func (rf RadioFieldInstance) fieldID() int {
	return rf.ID
}

func (rf RadioFieldInstance) toClient() client.Block {
	return client.RadioFieldBlock{
		ID:        rf.ID,
		Proposals: rf.Proposals,
	}
}

func (f RadioFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	_, ok := answer.(client.RadioAnswer)
	if !ok {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("expected RadioAnswer, got %T", answer),
		}
	}
	return nil
}

func (f RadioFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	return f.Answer == answer.(client.RadioAnswer).Index
}

type DropDownFieldInstance RadioFieldInstance

func (rf DropDownFieldInstance) fieldID() int { return rf.ID }

func (rf DropDownFieldInstance) toClient() client.Block {
	return client.DropDownFieldBlock{
		ID:        rf.ID,
		Proposals: rf.Proposals,
	}
}

func (f DropDownFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	return RadioFieldInstance(f).validateAnswerSyntax(answer)
}

func (f DropDownFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	return RadioFieldInstance(f).evaluateAnswer(answer)
}

// OrderedListFieldInstance asks the student to reorder part of the
// given symbols
type OrderedListFieldInstance struct {
	Label               string
	Answer              []StringOrExpression
	AdditionalProposals []StringOrExpression // added to Answer when displaying the field
	ID                  int
}

func (olf OrderedListFieldInstance) fieldID() int { return olf.ID }

// proposals groups Answer and AdditionalProposals and shuffle the list
// according to
func (olf OrderedListFieldInstance) proposals() (out []StringOrExpression) {
	out = append(append(out, olf.Answer...), olf.AdditionalProposals...)
	// shuffle in a deterministic way
	rd := rand.New(rand.NewSource(int64(len(olf.Answer)*1000 + len(olf.AdditionalProposals))))
	rd.Shuffle(len(out), func(i, j int) { out[i], out[j] = out[j], out[i] })
	return out
}

func (olf OrderedListFieldInstance) toClient() client.Block {
	out := client.OrderedListFieldBlock{
		ID:           olf.ID,
		Label:        olf.Label,
		AnswerLength: len(olf.Answer),
	}

	for _, v := range olf.proposals() {
		out.Proposals = append(out.Proposals, v.asLaTeX())
	}
	return out
}

func (olf OrderedListFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	list, ok := answer.(client.OrderedListAnswer)
	if !ok {
		return InvalidFieldAnswer{
			ID:     olf.ID,
			Reason: fmt.Sprintf("expected OrderedListAnswer, got %T", answer),
		}
	}

	if len(list.Indices) != len(olf.Answer) {
		return InvalidFieldAnswer{
			ID:     olf.ID,
			Reason: fmt.Sprintf("invalid answer length %d", len(list.Indices)),
		}
	}

	props := olf.proposals()
	for _, v := range list.Indices {
		if v >= len(props) {
			return InvalidFieldAnswer{
				ID:     olf.ID,
				Reason: fmt.Sprintf("invalid indice %d for length %d", v, len(props)),
			}
		}
	}

	return nil
}

func (olf OrderedListFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	list := answer.(client.OrderedListAnswer).Indices

	// reference and answer have the same length, as checked in `validateAnswerSyntax`
	proposals := olf.proposals()
	for i, ref := range olf.Answer {
		got := proposals[list[i]] // check in `validateAnswerSyntax`

		// we compare by value, not indices, since two different indices may have the same
		// value and then not be distinguable by the student,
		// and also, the indices has been shuffled
		if ref.asLaTeX() != got.asLaTeX() {
			return false
		}
	}

	return true
}

type FigurePointFieldInstance struct {
	Figure repere.Figure
	Answer repere.IntCoord
	ID     int
}

func (f FigurePointFieldInstance) fieldID() int { return f.ID }

func (f FigurePointFieldInstance) toClient() client.Block {
	return client.FigurePointFieldBlock{Figure: f.Figure, ID: f.ID}
}

func (f FigurePointFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	_, ok := answer.(client.PointAnswer)
	if !ok {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("expected PointAnswer, got %T", answer),
		}
	}
	return nil
}

func (f FigurePointFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	return f.Answer == answer.(client.PointAnswer).Point
}

type FigureVectorFieldInstance struct {
	Figure repere.Figure
	ID     int
	Answer repere.IntCoord
}

func (f FigureVectorFieldInstance) fieldID() int { return f.ID }

func (f FigureVectorFieldInstance) toClient() client.Block {
	return client.FigureVectorFieldBlock{Figure: f.Figure, ID: f.ID}
}

func (f FigureVectorFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	_, ok := answer.(client.DoublePointAnswer)
	if !ok {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("expected DoublePointAnswer, got %T", answer),
		}
	}
	return nil
}

func (f FigureVectorFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	ans := answer.(client.DoublePointAnswer)
	vector := repere.IntCoord{
		X: ans.To.X - ans.From.X,
		Y: ans.To.Y - ans.From.Y,
	}
	return f.Answer == vector
}

type FigureAffineLineFieldInstance struct {
	Label  string        // of the expect affine function
	Figure repere.Figure // usually empty, but set width and height
	ID     int
	Answer [2]float64 // a, b
}

func (f FigureAffineLineFieldInstance) fieldID() int { return f.ID }

func (f FigureAffineLineFieldInstance) toClient() client.Block {
	return client.FigureVectorFieldBlock{Figure: f.Figure, ID: f.ID, AsLine: true, LineLabel: f.Label}
}

func (f FigureAffineLineFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	ans, ok := answer.(client.DoublePointAnswer)
	if !ok {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("expected DoublePointAnswer, got %T", answer),
		}
	}

	if ans.To.X-ans.From.X == 0 {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: "invalid 0 x increment",
		}
	}
	return nil
}

func (f FigureAffineLineFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	ans := answer.(client.DoublePointAnswer)
	a := float64(ans.To.Y-ans.From.Y) / float64(ans.To.X-ans.From.X)
	b := float64(ans.From.Y) - a*float64(ans.From.X)
	return f.Answer == [2]float64{a, b}
}

type VectorPairCriterion uint8

const (
	VectorEquals VectorPairCriterion = iota
	VectorColinear
	VectorOrthogonal
)

type FigureVectorPairFieldInstance struct {
	Figure    repere.Figure
	ID        int
	Criterion VectorPairCriterion
}

func (f FigureVectorPairFieldInstance) fieldID() int { return f.ID }

func (f FigureVectorPairFieldInstance) toClient() client.Block {
	return client.FigureVectorPairFieldBlock{Figure: f.Figure, ID: f.ID}
}

func (f FigureVectorPairFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	_, ok := answer.(client.DoublePointPairAnswer)
	if !ok {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("expected DoublePointPairAnswer, got %T", answer),
		}
	}
	return nil
}

func (f FigureVectorPairFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	ans := answer.(client.DoublePointPairAnswer)
	vector1 := repere.IntCoord{
		X: ans.To1.X - ans.From1.X,
		Y: ans.To1.Y - ans.From1.Y,
	}
	vector2 := repere.IntCoord{
		X: ans.To2.X - ans.From2.X,
		Y: ans.To2.Y - ans.From2.Y,
	}
	switch f.Criterion {
	case VectorEquals:
		return vector1 == vector2
	case VectorColinear: // check if det(v1, v2) = 0
		return vector1.X*vector2.Y-vector1.Y*vector2.X == 0
	case VectorOrthogonal: // check if v1.v2 = 0
		return vector1.X*vector2.X+vector1.Y*vector2.Y == 0
	default:
		panic("exhaustive switch")
	}
}

type VariationTableFieldInstance struct {
	Answer VariationTableInstance
	ID     int
}

func (f VariationTableFieldInstance) fieldID() int { return f.ID }

func (f VariationTableFieldInstance) toClient() client.Block {
	return client.VariationTableFieldBlock{Length: len(f.Answer.Xs), ID: f.ID}
}

func (f VariationTableFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	ans, ok := answer.(client.VariationTableAnswer)
	if !ok {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("expected DoublePointPairAnswer, got %T", answer),
		}
	}

	if L := len(f.Answer.Xs); len(ans.Xs) != L || len(ans.Fxs) != L || len(ans.Arrows) != L-1 {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("invalid lengths %d %d %d", len(ans.Xs), len(ans.Fxs), len(ans.Arrows)),
		}
	}

	return nil
}

func (f VariationTableFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	ans := answer.(client.VariationTableAnswer)
	for i := range f.Answer.Xs {
		if ans.Xs[i] != f.Answer.Xs[i] || ans.Fxs[i] != f.Answer.Fxs[i] {
			return false
		}
		if i < len(f.Answer.Xs)-1 && !f.Answer.inferAlignment(i) != ans.Arrows[i] {
			return false
		}
	}

	return true
}

type FunctionPointsFieldInstance struct {
	Function *expression.Expression
	Label    string
	XGrid    []int
	Variable expression.Variable
	ID       int
}

func (f FunctionPointsFieldInstance) fieldID() int { return f.ID }

func (f FunctionPointsFieldInstance) toClient() client.Block {
	bounds, _, dfxs := functiongrapher.BoundsFromExpression(f.Function, f.Variable, f.XGrid)
	return client.FunctionPointsFieldBlock{
		Label: f.Label,
		Xs:    f.XGrid, ID: f.ID,
		Bounds: bounds,
		Dfxs:   dfxs,
	}
}

func (f FunctionPointsFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	ans, ok := answer.(client.FunctionPointsAnswer)
	if !ok {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("expected DoublePointPairAnswer, got %T", answer),
		}
	}

	if L := len(ans.Fxs); L != len(f.XGrid) {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("invalid length %d", L),
		}
	}

	return nil
}

func (f FunctionPointsFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	ans := answer.(client.FunctionPointsAnswer).Fxs
	_, ys, _ := functiongrapher.BoundsFromExpression(f.Function, f.Variable, f.XGrid)
	for i := range ys {
		if ans[i] != ys[i] {
			return false
		}
	}
	return true
}
