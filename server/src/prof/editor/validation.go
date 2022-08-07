package editor

import (
	"fmt"
	"strings"

	"github.com/benoitkugler/maths-online/utils"
)

// ValidateAllQuestions fetches all questions from the DB
// and calls Validate, returning all the errors encountered.
// It should be used as a maintenance helper when migrating the DB.
func ValidateAllQuestions(db DB) error {
	qus, err := SelectAllQuestions(db)
	if err != nil {
		return utils.SQLError(err)
	}

	exs := make(IdExerciceSet)
	for _, question := range qus {
		if question.NeedExercice.Valid {
			exs.Add(IdExercice(question.NeedExercice.Int64))
		}
	}

	exercices, err := SelectExercices(db, exs.Keys()...)
	if err != nil {
		return utils.SQLError(err)
	}

	return validateAllQuestions(qus, exercices)
}

func validateAllQuestions(questions Questions, exercices Exercices) error {
	var errs []string
	for id, q := range questions {
		if q.NeedExercice.Valid {
			ex := exercices[IdExercice(q.NeedExercice.Int64)]
			q.Page.Parameters = q.Page.Parameters.Append(ex.Parameters)
		}

		err := q.Page.Validate()
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s (ID: %d) -> %s", q.Page.Title, id, err))
		}
	}
	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("inconsistent table questions: %s", strings.Join(errs, "\n"))
}
