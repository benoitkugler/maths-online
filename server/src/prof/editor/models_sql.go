package editor

import (
	"database/sql"

	"github.com/benoitkugler/maths-online/maths/questions"
	"github.com/benoitkugler/maths-online/prof/teacher"
)

type (
	IdQuestion    int64
	IdExercice    int64
	IdProgression int64
)

// Question is a standalone question, used for instance in games.
type Question struct {
	Id          IdQuestion             `json:"id"`
	Page        questions.QuestionPage `json:"page"`
	Public      bool                   `json:"public"` // in practice only true for admins
	IdTeacher   teacher.IdTeacher      `json:"id_teacher"`
	Description string                 `json:"description"`
	// NeedExercice is not null if the question cannot be instantiated (or edited)
	// on its own
	NeedExercice sql.NullInt64 `json:"need_exercice" gomacro-sql-foreign:"Exercice"`
}

// gomacro:SQL ADD UNIQUE(IdQuestion, Tag)
type QuestionTag struct {
	Tag        string     `json:"tag"`
	IdQuestion IdQuestion `gomacro-sql-on-delete:"CASCADE" json:"id_question"`
}

// DifficultyTag are special question tags used to indicate the
// difficulty of one question.
// It is used to select question among implicit groups
type DifficultyTag string

// Exercice is the data structure for a full exercice, composed of a list of questions.
// There are two kinds of exercice :
//	- parallel : all the questions are independant
//	- progression : the questions are linked together by a shared Parameters set
type Exercice struct {
	Id          IdExercice
	Title       string // displayed to the students
	Description string // used internally by the teachers
	// Parameters are parameters shared by all the questions,
	// which are added to the individual ones.
	// It will be empty for parallel exercices
	Parameters questions.Parameters
	Flow       Flow
	// IdTeacher is the owner of the exercice
	IdTeacher teacher.IdTeacher
	Public    bool
}

// TODO: check delete question API
// ExerciceQuestion models an ordered list of questions.
// All link items should be updated at once to preserve `Index` invariants
// gomacro:SQL ADD PRIMARY KEY (IdExercice, Index)
type ExerciceQuestion struct {
	IdExercice IdExercice `json:"id_exercice" gomacro-sql-on-delete:"CASCADE"`
	IdQuestion IdQuestion `json:"id_question"`
	Bareme     int        `json:"bareme"`
	Index      int        `json:"-"`
}

// Progression is the table storing the student progression
// for one exercice.
// Note that this data structure may also be used in memory,
// for instance for the editor loopback.
// gomacro:SQL ADD UNIQUE(Id, IdExercice)
type Progression struct {
	Id         IdProgression
	IdExercice IdExercice `gomacro-sql-on-delete:"CASCADE"`
}

// We enforce consistency with the additional `IdExercice` field
// gomacro:SQL ADD FOREIGN KEY (IdExercice, Index) REFERENCES ExerciceQuestion ON DELETE CASCADE
// gomacro:SQL ADD FOREIGN KEY (IdProgression, IdExercice) REFERENCES Progression (Id, IdExercice) ON DELETE CASCADE
type ProgressionQuestion struct {
	IdProgression IdProgression   `json:"id_progression" gomacro-sql-on-delete:"CASCADE"`
	IdExercice    IdExercice      `json:"id_exercice" gomacro-sql-on-delete:"CASCADE"`
	Index         int             `json:"index"` // in the question list
	History       QuestionHistory `json:"history"`
}
