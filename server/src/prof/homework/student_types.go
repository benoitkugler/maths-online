package homework

import (
	"github.com/benoitkugler/maths-online/pass"
	ho "github.com/benoitkugler/maths-online/sql/homework"
	sql "github.com/benoitkugler/maths-online/sql/tasks"
	"github.com/benoitkugler/maths-online/tasks"
	taAPI "github.com/benoitkugler/maths-online/tasks"
)

// used to generate Dart code

// SheetProgression is the summary of the progression
// of one student for one sheet
type SheetProgression struct {
	Sheet ho.Sheet
	Tasks []taAPI.TaskProgressionHeader
}

type StudentSheets = []SheetProgression

type StudentEvaluateTaskIn struct {
	StudentID pass.EncryptedID
	IdTask    sql.IdTask
	Ex        tasks.EvaluateWorkIn `gomacro-extern:"tasks#dart#package:eleve/shared_gen.dart"`
}

type StudentEvaluateTaskOut struct {
	Ex   tasks.EvaluateWorkOut `gomacro-extern:"tasks#dart#package:eleve/shared_gen.dart"`
	Mark int                   // updated mark
}
