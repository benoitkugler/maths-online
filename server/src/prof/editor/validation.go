package editor

import (
	"fmt"
	"strings"
)

// ValidateAllQuestions fetches all questions from the DB
// and calls Validate, returning all the errors encountered.
// It should be used as a maintenance helper when migrating the DB.
func ValidateAllQuestions(db DB) error {
	qu, err := SelectAllQuestions(db)
	if err != nil {
		return err
	}

	return validateAllQuestions(qu)
}

func validateAllQuestions(questions Questions) error {
	var errs []string
	for id, q := range questions {
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
