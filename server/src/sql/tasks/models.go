package tasks

import (
	"github.com/benoitkugler/maths-online/sql/editor"
	"github.com/benoitkugler/maths-online/sql/teacher"
)

type (
	IdProgression  int64
	IdTask         int64
	IdMonoquestion int64
)

// Monoquestion is a shortcut for an exercice composed of only one question.
// It is used to avoid creating cumbersome exercice wrappers around questions.
// It may only be used for standalone questions, that is question in a group
type Monoquestion struct {
	Id         IdMonoquestion
	IdQuestion editor.IdQuestion
	NbRepeat   int
	Bareme     int // for one question
}

// Task is a pointer to an assignement to one exercice, either
// an `Exercice` or a `Monoquestion`
// gomacro:SQL ADD UNIQUE(Id, IdExercice)
// gomacro:SQL ADD CHECK(IdExercice IS NOT NULL OR IdMonoquestion IS NOT NULL)
// gomacro:SQL ADD CHECK(IdExercice IS NULL OR IdMonoquestion IS NULL)
type Task struct {
	Id             IdTask
	IdExercice     editor.OptionalIdExercice `gomacro-sql-foreign:"Exercice"`
	IdMonoquestion OptionalIdMonoquestion    `gomacro-sql-foreign:"Monoquestion"`
}

// Progression is the table storing the student progression
// for one exercice.
// gomacro:SQL ADD UNIQUE(IdStudent, IdTask)
type Progression struct {
	Id IdProgression

	IdStudent teacher.IdStudent `gomacro-sql-on-delete:"CASCADE"`

	IdTask IdTask `gomacro-sql-on-delete:"CASCADE"`
}

// ProgressionQuestion stores the result for one question,
// either in a Exercice or in a Monoquestion
type ProgressionQuestion struct {
	IdProgression IdProgression   `json:"id_progression" gomacro-sql-on-delete:"CASCADE"`
	Index         int             `json:"index"` // in the question list
	History       QuestionHistory `json:"history"`
}
