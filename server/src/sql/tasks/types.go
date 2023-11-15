// Package tasks exposes the data structure
// required to assign exercices during activities,
// and to track the progression of the students.
package tasks

import (
	"sort"

	"github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/utils"
)

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

// Stats returns the number of tries
func (qh QuestionHistory) Stats() (success, failure int) {
	for _, try := range qh {
		if try {
			success++
		} else {
			failure++
		}
	}
	return
}

// EnsureOrder must be call on the questions of one exercice,
// to make sure the order in the slice is consistent with the one
// indicated by `Index`
func (l RandomMonoquestionVariants) EnsureOrder() {
	sort.Slice(l, func(i, j int) bool { return l[i].Index < l[j].Index })
}

// EnsureOrder must be call on the questions of one exercice,
// to make sure the order in the slice is consistent with the one
// indicated by `Index`
func (l Progressions) EnsureOrder() {
	sort.Slice(l, func(i, j int) bool { return l[i].Index < l[j].Index })
}

// ResizeProgressions makes sure the given task has progression items
// in range [0; nbQuestions[
// This function should be called when updating a task content.
func ResizeProgressions(db DB, id IdTask, nbQuestions int) error {
	_, err := db.Exec("DELETE FROM progressions WHERE idTask = $1 AND Index >= $2", id, nbQuestions)
	if err != nil {
		return utils.SQLError(err)
	}

	task, err := SelectTask(db, id)
	if err != nil {
		return utils.SQLError(err)
	}

	if idM := task.IdRandomMonoquestion; idM.Valid {
		// also remove the variants
		_, err := db.Exec("DELETE FROM random_monoquestion_variants WHERE idRandomMonoquestion = $1 AND Index >= $2", idM.ID, nbQuestions)
		if err != nil {
			return utils.SQLError(err)
		}
	}

	return nil
}
