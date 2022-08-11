package tasks

import (
	"github.com/benoitkugler/maths-online/prof/editor"
	"github.com/benoitkugler/maths-online/prof/teacher"
)

type (
	IdProgression int64
	IdTask        int64
)

// Task is a pointer to an assignement to one exercice
// gomacro:SQL ADD UNIQUE(Id, IdExercice)
type Task struct {
	Id         IdTask
	IdExercice editor.IdExercice
}

// Progression is the table storing the student progression
// for one exercice.
// gomacro:SQL ADD UNIQUE(IdStudent, IdTask)
// gomacro:SQL ADD UNIQUE(Id, IdExercice)
// gomacro:SQL ADD FOREIGN KEY (IdTask, IdExercice) REFERENCES Task (Id, IdExercice)
type Progression struct {
	Id IdProgression

	IdStudent teacher.IdStudent `gomacro-sql-on-delete:"CASCADE"`

	IdTask IdTask `gomacro-sql-on-delete:"CASCADE"`
	// IdExercice is used for consistency with ProgressionQuestion
	IdExercice editor.IdExercice `gomacro-sql-on-delete:"CASCADE"`
}

// We enforce consistency with the additional `IdExercice` field
// gomacro:SQL ADD FOREIGN KEY (IdExercice, Index) REFERENCES exercice_questions ON DELETE CASCADE
// gomacro:SQL ADD FOREIGN KEY (IdProgression, IdExercice) REFERENCES Progression (Id, IdExercice) ON DELETE CASCADE
type ProgressionQuestion struct {
	IdProgression IdProgression     `json:"id_progression" gomacro-sql-on-delete:"CASCADE"`
	IdExercice    editor.IdExercice `json:"id_exercice" gomacro-sql-on-delete:"CASCADE"`
	Index         int               `json:"index"` // in the question list
	History       QuestionHistory   `json:"history"`
}
