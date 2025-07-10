package homework

import (
	"errors"

	"github.com/benoitkugler/maths-online/server/src/sql/tasks"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/benoitkugler/maths-online/server/src/utils"
)

type OptionalIdTravail struct {
	Valid bool
	ID    IdTravail
}

func (id IdTravail) AsOptional() OptionalIdTravail {
	return OptionalIdTravail{ID: id, Valid: true}
}

func LoadMonoquestionSheet(db DB, idMono tasks.IdMonoquestion) (tasks.IdTask, IdSheet, error) {
	ts, err := tasks.SelectTasksByIdMonoquestions(db, idMono)
	if err != nil {
		return 0, 0, utils.SQLError(err)
	}
	if len(ts) != 1 {
		return 0, 0, errors.New("internal error: expected one task for a monoquestion")
	}
	link, found, err := SelectSheetTaskByIdTask(db, ts.IDs()[0])
	if err != nil {
		return 0, 0, err
	}
	if !found {
		return 0, 0, errors.New("internal error: task without sheet")
	}
	return link.IdTask, link.IdSheet, nil
}

func LoadRandomMonoquestionSheet(db DB, idMono tasks.IdRandomMonoquestion) (tasks.IdTask, IdSheet, error) {
	ts, err := tasks.SelectTasksByIdRandomMonoquestions(db, idMono)
	if err != nil {
		return 0, 0, utils.SQLError(err)
	}
	if len(ts) != 1 {
		return 0, 0, errors.New("internal error: expected one task for a monoquestion")
	}
	link, found, err := SelectSheetTaskByIdTask(db, ts.IDs()[0])
	if err != nil {
		return 0, 0, err
	}
	if !found {
		return 0, 0, errors.New("internal error: task without sheet")
	}
	return link.IdTask, link.IdSheet, nil
}

// IsVisibleBy returns `true` if the Sheet is public or
// owned by `userID`
func (qu Sheet) IsVisibleBy(userID teacher.IdTeacher) bool {
	return qu.Public || qu.IdTeacher == userID
}

// RestrictVisible remove the sheets not visible by `userID`
func (qus Sheets) RestrictVisible(userID teacher.IdTeacher) {
	for id, qu := range qus {
		if !qu.IsVisibleBy(userID) {
			delete(qus, id)
		}
	}
}

type QuestionRepeat uint8

const (
	Unlimited QuestionRepeat = iota // Illimité
	OneTry                          // Un seul
)
