package homework

import (
	"github.com/benoitkugler/maths-online/server/src/pass"
	"github.com/benoitkugler/maths-online/server/src/sql/events"
	ho "github.com/benoitkugler/maths-online/server/src/sql/homework"
	sql "github.com/benoitkugler/maths-online/server/src/sql/tasks"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	taAPI "github.com/benoitkugler/maths-online/server/src/tasks"
)

// used to generate Dart code

// SheetProgression is the summary of the progression
// of one student for one sheet

type Sheet struct {
	Id            ho.IdSheet
	Title         string
	Noted         bool // new in version 1.5
	Deadline      ho.Time
	IgnoreForMark bool // new in version 1.6.4

	Matiere teacher.MatiereTag // new in version 1.9

	QuestionRepeat    ho.QuestionRepeat // new in version 1.9
	QuestionTimeLimit int               // new in version 1.9
}

type SheetProgression struct {
	IdTravail ho.IdTravail // new in version 1.5
	Sheet     Sheet
	Tasks     []taAPI.TaskProgressionHeader
}

type StudentSheets []SheetProgression

type StudentEvaluateTaskIn struct {
	StudentID pass.EncryptedID
	IdTask    sql.IdTask
	Ex        taAPI.EvaluateWorkIn
	IdTravail ho.IdTravail
}

type StudentEvaluateTaskOut struct {
	Ex   taAPI.EvaluateWorkOut
	Mark int // updated mark
	// WasProgressionRegistred is true if the server has updated the DB
	// It should be used to decide whether or not to update the sheet list.
	WasProgressionRegistred bool                     // new in v1.6.8
	Advance                 events.EventNotification // new in v1.7
}

type StudentResetTaskIn struct {
	StudentID pass.EncryptedID
	IdTravail ho.IdTravail
	IdTask    sql.IdTask
}
