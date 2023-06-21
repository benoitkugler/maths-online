// Package tasks exposes the data structure
// required to assign exercices during activities,
// and to track the progression of the students.
package tasks

import "github.com/benoitkugler/maths-online/server/src/sql/editor"

type OptionalIdQuestion struct {
	Valid bool
	ID    editor.IdQuestion
}

type OptionalIdMonoquestion struct {
	Valid bool
	ID    IdMonoquestion
}

func (id IdMonoquestion) AsOptional() OptionalIdMonoquestion {
	return OptionalIdMonoquestion{ID: id, Valid: true}
}

type OptionalIdRandomMonoquestion struct {
	Valid bool
	ID    IdRandomMonoquestion
}

func (id IdRandomMonoquestion) AsOptional() OptionalIdRandomMonoquestion {
	return OptionalIdRandomMonoquestion{ID: id, Valid: true}
}

// QuestionHistory stores the successes for one question,
// in chronological order.
// For instance, [true, false, true] means : first try: correct, second: wrong answer,third: correct
type QuestionHistory []bool

// Success return true if at least one try is sucessful
func (qh QuestionHistory) Success() bool {
	for _, try := range qh {
		if try {
			return true
		}
	}
	return false
}
