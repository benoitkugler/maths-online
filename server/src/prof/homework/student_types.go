package homework

import (
	"github.com/benoitkugler/maths-online/pass"
	"github.com/benoitkugler/maths-online/prof/editor"
	"github.com/benoitkugler/maths-online/tasks"
)

// used to generate Dart code

type StudentSheets = []SheetProgression

type StudentEvaluateExerciceIn struct {
	StudentID pass.EncryptedID
	IdTask    tasks.IdTask
	Ex        editor.EvaluateExerciceIn `gomacro-extern:"editor:dart:../shared_gen.dart"`
}

type StudentEvaluateExerciceOut struct {
	Ex   editor.EvaluateExerciceOut `gomacro-extern:"editor:dart:../shared_gen.dart"`
	Mark int                        // updated mark
}
