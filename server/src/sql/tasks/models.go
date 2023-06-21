package tasks

import (
	"github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
)

type (
	IdProgression        int64
	IdTask               int64
	IdMonoquestion       int64
	IdRandomMonoquestion int64
)

// Monoquestion is a shortcut for an exercice composed of only one question.
// It is used to avoid creating cumbersome exercice wrappers around questions.
type Monoquestion struct {
	Id         IdMonoquestion
	IdQuestion editor.IdQuestion
	NbRepeat   int
	Bareme     int // for one question
}

// RandomMonoquestion allows the teacher to specify a whole [Questiongroup],
// with questions chosen randomly for each student, according to an (optional)
// difficulty tag.
type RandomMonoquestion struct {
	Id              IdRandomMonoquestion
	IdQuestiongroup editor.IdQuestiongroup
	NbRepeat        int
	Bareme          int                  // for one question
	Difficulty      editor.DifficultyTag // optional, empty for all questions
}

// Task is a pointer to an assignement to one exercice, either
// an `Exercice`, a `Monoquestion` or a `RandomMonoquestion`
// gomacro:SQL ADD UNIQUE(Id, IdExercice)
// gomacro:SQL ADD CHECK((IdExercice IS NOT NULL)::int + (IdMonoquestion IS NOT NULL)::int + (IdRandomMonoquestion IS NOT NULL)::int = 1)
type Task struct {
	Id                   IdTask
	IdExercice           editor.OptionalIdExercice    `gomacro-sql-foreign:"Exercice"`
	IdMonoquestion       OptionalIdMonoquestion       `gomacro-sql-foreign:"Monoquestion"`
	IdRandomMonoquestion OptionalIdRandomMonoquestion `gomacro-sql-foreign:"RandomMonoquestion"`
}

// Progression is a link table storing the student progressions
// on tasks.
// gomacro:SQL ADD UNIQUE(IdStudent, IdTask, Index)
// gomacro:SQL ADD CHECK((IdQuestionVariant IS NOT NULL)::int + (IdExerciceVariant IS NOT NULL)::int = 1)
// gomacro:SQL _SELECT KEY (IdStudent, IdTask)
type Progression struct {
	IdStudent teacher.IdStudent `gomacro-sql-on-delete:"CASCADE"`
	IdTask    IdTask            `gomacro-sql-on-delete:"CASCADE"`

	// Index in the question list
	// For exercice, it is the question number
	// For monoquestion, it is the "repetion" number
	Index int `json:"index"`

	History QuestionHistory `json:"history"`

	// The following fields are mainly used when the task is random.
	// Otherwise, they are just copied from the task.

	IdQuestionVariant OptionalIdQuestion        `gomacro-sql-foreign:"Question"`
	IdExerciceVariant editor.OptionalIdExercice `gomacro-sql-foreign:"Exercice"`
}
