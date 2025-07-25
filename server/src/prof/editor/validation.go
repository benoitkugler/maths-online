package editor

import (
	"fmt"
	"strings"

	ed "github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/utils"
)

// ValidateAllQuestions fetches all questions from the DB
// and calls Validate, returning all the errors encountered.
// It should be used as a maintenance helper when migrating the DB.
func ValidateAllQuestions(db ed.DB) error {
	qus, err := ed.SelectAllQuestions(db)
	if err != nil {
		return utils.SQLError(err)
	}

	exs := utils.Set[ed.IdExercice]{}
	for _, question := range qus {
		if question.NeedExercice.Valid {
			exs.Add(question.NeedExercice.ID)
		}
	}

	exercices, err := ed.SelectExercices(db, exs.Keys()...)
	if err != nil {
		return utils.SQLError(err)
	}

	return validateAllQuestions(qus, exercices)
}

func validateAllQuestions(qus ed.Questions, exercices ed.Exercices) error {
	var errs []string
	for id, q := range qus {
		page := q.Page()
		if q.NeedExercice.Valid {
			ex := exercices[q.NeedExercice.ID]
			page.Parameters = append(page.Parameters, ex.Parameters...)
		}

		err := page.Validate()
		if err != nil {
			errs = append(errs, fmt.Sprintf("ID: %d (%s) -> %s", id, q.Subtitle, err))
		}
	}
	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("%d invalid table questions: %s", len(errs), strings.Join(errs, "\n"))
}
