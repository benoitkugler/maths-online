package exercice

import (
	"fmt"
	"math/rand"

	"github.com/benoitkugler/maths-online/maths/exercice/client"
	"github.com/benoitkugler/maths-online/maths/expression"
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
	_ fieldInstance = OrderedListFieldInstance{}
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
	Label *expression.Expression // if not nil, the field is displayed on a new line, prefixed by <expression> =

	Answer          *expression.Expression
	ComparisonLevel expression.ComparisonLevel
	ID              int
}

func (f ExpressionFieldInstance) fieldID() int { return f.ID }

func (f ExpressionFieldInstance) toClient() client.Block {
	var label string
	if f.Label != nil {
		label = f.Label.AsLaTeX(nil)
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

// OrderedListFieldInstance asks the student to reorder part of the
// given symbols
type OrderedListFieldInstance struct {
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
