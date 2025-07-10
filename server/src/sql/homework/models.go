package homework

import (
	"database/sql"
	"sort"
	"time"

	"github.com/benoitkugler/maths-online/server/src/sql/tasks"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
)

//go:generate ../../../../../gomacro/cmd/gomacro models.go sql:gen_create.sql go/sqlcrud:gen_scans.go go/randdata:gen_randdata_test.go

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
// gomacro:SQL ADD UNIQUE(Id, IdSheet)
type Travail struct {
	Id          IdTravail
	IdClassroom teacher.IdClassroom `gomacro-sql-on-delete:"CASCADE"`
	IdSheet     IdSheet             `gomacro-sql-on-delete:"CASCADE"`

	// When 'true', the [Sheet] is evaluated, and may only
	// be done until the [Deadline].
	// Notation : a question gives point if it has been successfully completed (at least) once.
	//
	// When 'false' the sheet is always available as free training.
	Noted bool

	// If [Noted] is true,
	// passed the Deadline, the sheet notations may not be modified anymore.
	// Else, this field is ignored
	Deadline Time

	// Pospone the access for students to this work
	ShowAfter Time

	QuestionRepeat QuestionRepeat

	// When not zero, every question is time limited.
	// (in seconds, zero means no limit)
	QuestionTimeLimit int
}

// Sheet is a list of exercices.
// gomacro:SQL ADD FOREIGN KEY (Id, Anonymous) REFERENCES Travail(IdSheet, Id) ON DELETE CASCADE
type Sheet struct {
	Id IdSheet

	Title string

	// IdTeacher is the creator of the [Sheet]
	IdTeacher teacher.IdTeacher `gomacro-sql-on-delete:"CASCADE"`

	Level string // tag to classify by expected level, ignored on anonymous sheets

	// Anonymous is not null when the sheet is only
	// link to one [Travail].
	// Anonymous [Sheet]s are deleted when the [Travail] is,
	// and are not shown in the favorites sheets panel.
	Anonymous OptionalIdTravail `gomacro-sql-on-delete:"CASCADE" gomacro-sql-foreign:"Travail"`

	Public bool // only true for admin account

	Matiere teacher.MatiereTag // tag to classify by expected topic, ignored on anonymous sheets
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

// TravailException is a link table storing per student
// settings for a [Travail]
//
// gomacro:SQL ADD UNIQUE(IdStudent, IdTravail)
// gomacro:SQL _SELECT KEY(IdStudent, IdTravail)
type TravailException struct {
	IdStudent teacher.IdStudent `gomacro-sql-on-delete:"CASCADE"`
	IdTravail IdTravail         `gomacro-sql-on-delete:"CASCADE"`

	// [Deadline] is an optionnal deadline overriding the one
	// setup in the related [Travail].
	Deadline sql.NullTime

	// [IgnoreForMark] may be set to true to ignore this mark
	// when displaying the average.
	IgnoreForMark bool
}
