package homework

import (
	"github.com/benoitkugler/maths-online/server/src/pass"
	ho "github.com/benoitkugler/maths-online/server/src/sql/homework"
	sql "github.com/benoitkugler/maths-online/server/src/sql/tasks"
	"github.com/benoitkugler/maths-online/server/src/tasks"
	taAPI "github.com/benoitkugler/maths-online/server/src/tasks"
)

// used to generate Dart code

// SheetProgression is the summary of the progression
// of one student for one sheet
type SheetProgression struct {
	Sheet ho.Sheet
	Tasks []taAPI.TaskProgressionHeader
}

type StudentSheets []SheetProgression

type StudentEvaluateTaskIn struct {
	StudentID pass.EncryptedID
	IdTask    sql.IdTask
	Ex        tasks.EvaluateWorkIn
}

type StudentEvaluateTaskOut struct {
	Ex   tasks.EvaluateWorkOut
	Mark int // updated mark
}
