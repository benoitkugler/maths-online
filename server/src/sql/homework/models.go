package homework

import (
	"sort"
	"time"

	"github.com/benoitkugler/maths-online/server/src/sql/tasks"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
)

type (
	IdSheet   int64
	IdTravail int64
)

// Time is an instant in a day.
type Time time.Time

func (d Time) MarshalJSON() ([]byte, error)     { return time.Time(d).MarshalJSON() }
func (d *Time) UnmarshalJSON(data []byte) error { return (*time.Time)(d).UnmarshalJSON(data) }

// Travail associates a [Sheet] to a classroom, with
// an optional deadline
type Travail struct {
	Id          IdTravail
	IdClassroom teacher.IdClassroom `gomacro-sql-on-delete:"CASCADE"`
	IdSheet     IdSheet             `gomacro-sql-on-delete:"CASCADE"`

	// When 'true', the [Sheet] is evaluated, and may only
	// be done once.
	// Notation : a question gives point if it has been successfully completed (at least) once
	// When 'false' the sheet is always available as free training.
	Noted bool

	// If [Noted] is true,
	// passed the Deadline, the sheet notations may not be modified anymore.
	// Else, this field is ignored
	Deadline Time
}

// Sheet is a list of exercices.
type Sheet struct {
	Id IdSheet

	Title string

	// IdTeacher is the creator of the [Sheet]
	IdTeacher teacher.IdTeacher `gomacro-sql-on-delete:"CASCADE"`

	Level string // tag to classify by expected level
}

// gomacro:SQL ADD PRIMARY KEY (IdSheet, Index)
// A task may only appear in one sheet
// gomacro:SQL ADD UNIQUE (IdTask)
type SheetTask struct {
	IdSheet IdSheet `gomacro-sql-on-delete:"CASCADE"`
	Index   int     `json:"-"` // order in the list
	IdTask  tasks.IdTask
}

// EnsureOrder enforce the slice order indicated by `Index`
func (l SheetTasks) EnsureOrder() {
	sort.Slice(l, func(i, j int) bool { return l[i].Index < l[j].Index })
}
