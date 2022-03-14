package exercice

import (
	"fmt"

	"github.com/benoitkugler/maths-online/maths/exercice/client"
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
	validateAnswerSyntax(answer client.Answer) error

	// evaluateAnswer evaluate the given answer against the reference
	// an error may be returned, for instance against malicious query
	// validateAnswerSyntax is assumed to have already been called on `answer`
	evaluateAnswer(answer client.Answer) (isCorrect bool, err error)
}

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

func (f NumberFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool, err error) {
	return f.Answer == answer.(client.NumberAnswer).Value, nil
}

type ListFieldInstance struct{}

// TODO:
func (ListFieldInstance) toClient() client.Block { return client.ListFieldBlock{} }

type FormulaFieldInstance struct{}

// TODO:
func (FormulaFieldInstance) toClient() client.Block { return client.FormulaFieldBlock{} }
