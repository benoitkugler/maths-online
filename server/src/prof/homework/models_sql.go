package homework

import (
	"time"

	"github.com/benoitkugler/maths-online/prof/editor"
	"github.com/benoitkugler/maths-online/prof/teacher"
)

type IdSheet int64

// Notation is the kind of notation applied
// to a sheet, if any
type Notation uint8

const (
	NoNotation      Notation = iota // no notation
	SuccessNotation                 // a question gives point if it has been successfully completed (at least) once
)

// Time is an instant in a day.
type Time time.Time

func (d Time) MarshalJSON() ([]byte, error)     { return time.Time(d).MarshalJSON() }
func (d *Time) UnmarshalJSON(data []byte) error { return (*time.Time)(d).UnmarshalJSON(data) }

// Sheet is a list of exercices with
// a due date
type Sheet struct {
	Id          IdSheet
	IdClassroom teacher.IdClassroom `gomacro-sql-on-delete:"CASCADE"`
	Title       string
	Notation    Notation
	// If false, the sheet is in preparation, not shown to the student.
	Activated bool
	// Passed the Deadline, the sheet notations may not be modified anymore.
	Deadline Time
}

type SheetExercice struct {
	IdSheet    IdSheet `gomacro-sql-on-delete:"CASCADE"`
	IdExercice editor.IdExercice
	Index      int `json:"-"` // order in the list
}
