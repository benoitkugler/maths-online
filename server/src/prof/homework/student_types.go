package homework

import (
	"github.com/benoitkugler/maths-online/server/src/pass"
	ho "github.com/benoitkugler/maths-online/server/src/sql/homework"
	sql "github.com/benoitkugler/maths-online/server/src/sql/tasks"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	taAPI "github.com/benoitkugler/maths-online/server/src/tasks"
)

// used to generate Dart code

// SheetProgression is the summary of the progression
// of one student for one sheet

type Sheet struct {
	Id       ho.IdSheet
	Title    string
	Noted    bool // new in version 1.5
	Deadline ho.Time

	// TODO: cleanup these unused fields

	Notation    int                 // Deprecated
	Activated   bool                // Deprecated
	IdClassroom teacher.IdClassroom // Deprecated
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
}

type StudentResetTaskIn struct {
	StudentID pass.EncryptedID
	IdTravail ho.IdTravail
	IdTask    sql.IdTask
}
