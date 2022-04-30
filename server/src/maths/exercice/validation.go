package exercice

import (
	"fmt"
	"strings"

	"github.com/benoitkugler/maths-online/maths/expression"
)

// maxFunctionBound is the maximum value a function
// may reached. Higher values are either a bug, or won't be properly
// displayed on the student client
const maxFunctionBound = 100

type ErrParameters struct {
	Origin  string
	Details string
}

func (err ErrParameters) Error() string {
	return "invalid random parameters"
}

// Validate ensure the given `Parameters` are sound,
// by parsing the expression, checking for duplicate parameters,
// and detecting definition cycles.
// It the error is not nil, it will be of type `ErrParameters`.
// Once called without error, `ToMap` may be safely used.
func (pr Parameters) Validate() error {
	params := make(expression.RandomParameters)
	for _, def := range pr.Variables {
		if _, has := params[def.Variable]; has {
			return ErrParameters{
				Origin:  def.Expression,
				Details: expression.ErrDuplicateParameter{Duplicate: def.Variable}.Error(),
			}
		}

		expr, err := expression.Parse(def.Expression)
		if err != nil {
			return ErrParameters{
				Origin:  def.Expression,
				Details: err.Error(),
			}
		}

		params[def.Variable] = expr
	}

	for _, it := range pr.Intrinsics {
		parsed, err := expression.ParseIntrinsic(it)
		if err != nil {
			return ErrParameters{
				Origin:  it,
				Details: err.Error(),
			}
		}

		err = parsed.MergeTo(params)
		if err != nil {
			return ErrParameters{
				Origin:  it,
				Details: err.Error(),
			}
		}
	}

	for v := range params {
		if v.Name == 'e' {
			return ErrParameters{
				Origin:  v.String(),
				Details: "La variable e n'est pas autorisée (car utilisée pour exp).",
			}
		}
	}

	_, err := params.Instantiate()
	if err != nil {
		return ErrParameters{
			Origin:  "Liste des paramètres",
			Details: err.Error(),
		}
	}

	return nil
}

type errEnonce struct {
	Error string // detailed error
	Block int    // index of the invalid block
}

// ErrQuestionInvalid is returned by  Question.Validate()
// It is either an error about the random parameters, or the blocks content.
type ErrQuestionInvalid struct {
	ErrParameters     ErrParameters
	ErrEnonce         errEnonce
	ParametersInvalid bool
}

func (e ErrQuestionInvalid) Error() string {
	return "invalid question content"
}

// Validate ensure the enonce blocks are sound.
// If not, an `ErrQuestionInvalid` is returned.
func (qu Question) Validate() error {
	// the client validate the random parameters on the fly,
	// so they should be valid here
	// err on the side of caution though
	if err := qu.Parameters.Validate(); err != nil {
		return ErrQuestionInvalid{ParametersInvalid: true, ErrParameters: err.(ErrParameters)}
	}

	params := qu.Parameters.ToMap()
	for i, block := range qu.Enonce {
		err := block.validate(params)
		if err != nil {
			return ErrQuestionInvalid{ErrEnonce: errEnonce{Block: i, Error: err.Error()}}
		}
	}

	return nil
}

// ValidateAllQuestions fetches all questions from the DB
// and calls Validate, returning all the errors encountered.
// It should be used as a maintenance helper when migrating the DB.
func ValidateAllQuestions(db DB) error {
	qu, err := SelectAllQuestions(db)
	if err != nil {
		return err
	}

	var errs []string
	for id, q := range qu {
		err := q.Validate()
		if err != nil {
			errs = append(errs, fmt.Sprintf("ID: %d -> %s", id, err))
		}
	}
	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("inconsistent table questions: %s", strings.Join(errs, "\n"))
}
