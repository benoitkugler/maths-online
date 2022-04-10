package exercice

import (
	"github.com/benoitkugler/maths-online/maths/expression"
)

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

		expr, _, err := expression.Parse(def.Expression)
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

	_, err := params.Instantiate()
	if err != nil {
		return ErrParameters{
			Origin:  "Liste des paramètres",
			Details: err.Error(),
		}
	}

	return nil
}
