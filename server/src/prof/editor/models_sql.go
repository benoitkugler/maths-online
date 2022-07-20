package editor

import (
	"database/sql"

	"github.com/benoitkugler/maths-online/maths/questions"
)

//go:generate ../../../../../structgen/structgen -source=models_sql.go -mode=sql:gen_scans.go -mode=sql_gen:gen_create.sql -mode=rand:gen_randdata_test.go

// Question is a standalone question, used for instance in games.
type Question struct {
	Id          int64                  `json:"id"`
	Page        questions.QuestionPage `json:"page"`
	Public      bool                   `json:"public"` // in practice only true for admins
	IdTeacher   int64                  `json:"id_teacher"`
	Description string                 `json:"description"`
	// NeedExercice is not null if the question cannot be instantiated (or edited)
	// on its own
	NeedExercice sql.NullInt64 `json:"need_exercice" sql_foreign_key:"exercice"`
}

// sql: ADD UNIQUE(id_question, tag)
type QuestionTag struct {
	Tag        string `json:"tag"`
	IdQuestion int64  `sql_on_delete:"CASCADE" json:"id_question"`
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
	Id          int64
	Title       string // displayed to the students
	Description string // used internally by the teachers
	// Parameters are parameters shared by all the questions,
	// which are added to the individual ones.
	// It will be empty for parallel exercices
	Parameters questions.Parameters
	Flow       Flow
	// IdTeacher is the owner of the exercice
	IdTeacher int64 `json:"id_teacher"`
	Public    bool
}

// TODO: check delete question API
// ExerciceQuestion models an ordered list of questions.
// All link items should be updated at once to preserve `Index` invariants
// sql: ADD PRIMARY KEY (id_exercice, index)
type ExerciceQuestion struct {
	IdExercice int64 `json:"id_exercice" sql_on_delete:"CASCADE"`
	IdQuestion int64 `json:"id_question"`
	Bareme     int   `json:"bareme"`
	Index      int   `json:"-" sql:"index"`
}

// Progression is the table storing the student progression
// for one exercice.
// Note that this data structure may also be used in memory,
// for instance for the editor loopback.
// sql: ADD UNIQUE(id, id_exercice)
type Progression struct {
	Id         int64
	IdExercice int64 `json:"id_exercice" sql_on_delete:"CASCADE"`
}

// We enforce consistency with the additional `id_exercice` field
// sql: ADD FOREIGN KEY (id_exercice, index) REFERENCES exercice_questions ON DELETE CASCADE
// sql: ADD FOREIGN KEY (id_progression, id_exercice) REFERENCES progressions (id, id_exercice) ON DELETE CASCADE
type ProgressionQuestion struct {
	IdProgression int64           `json:"id_progression" sql_on_delete:"CASCADE"`
	IdExercice    int64           `json:"id_exercice" sql_on_delete:"CASCADE"`
	Index         int             `json:"index"` // in the question list
	History       QuestionHistory `json:"history"`
}
