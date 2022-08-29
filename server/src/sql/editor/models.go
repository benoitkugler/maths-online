// Package editor provides the data structures used to
// create questions and exercices
package editor

import (
	"github.com/benoitkugler/maths-online/maths/questions"
	"github.com/benoitkugler/maths-online/sql/teacher"
)

type (
	IdQuestion      int64
	IdExercice      int64
	IdQuestiongroup int64
	IdExercicegroup int64
)

// Question is a standalone question version, used for instance in games
// or as part of a full exercice.
// gomacro:SQL ADD CHECK(NeedExercice IS NOT NULL OR IdGroup IS NOT NULL)
type Question struct {
	Id          IdQuestion
	Page        questions.QuestionPage
	Description string

	Difficulty DifficultyTag

	// NeedExercice is not null if the question is part
	// of an exercice and requires its parameters to be
	// instantiated
	NeedExercice OptionalIdExercice `gomacro-sql-foreign:"Exercice"`

	// IdGroup is not null for the standalone questions, accessed by a group
	IdGroup OptionalIdQuestiongroup `gomacro-sql-foreign:"Questiongroup" gomacro-sql-on-delete:"CASCADE"`
}

// Questiongroup groups several variant of the same question
type Questiongroup struct {
	Id        IdQuestiongroup
	Title     string
	Public    bool // in practice only true for admins
	IdTeacher teacher.IdTeacher
}

// gomacro:SQL ADD UNIQUE(IdQuestiongroup, Tag)
type QuestiongroupTag struct {
	Tag             string
	IdQuestiongroup IdQuestiongroup `gomacro-sql-on-delete:"CASCADE"`
}

// Exercicegroup groups the variant of the same exercice
type Exercicegroup struct {
	Id        IdExercicegroup
	Title     string            // title shown to the student
	Public    bool              // in practice only true for admins
	IdTeacher teacher.IdTeacher // IdTeacher is the owner of the exercice
}

// Exercice is the data structure for a full exercice, composed of a list of questions.
// There are two kinds of exercice :
//	- parallel : all the questions are independant
//	- progression : the questions are linked together by a shared Parameters set
type Exercice struct {
	Id      IdExercice
	IdGroup IdExercicegroup

	Title       string // subtitle, only shown to the teacher
	Description string // used internally by the teachers
	// Parameters are parameters shared by all the questions,
	// which are added to the individual ones.
	// It will be empty for parallel exercices
	Parameters questions.Parameters
	Flow       Flow
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

// DifficultyTag are special question tags used to indicate the
// difficulty of one question.
// It is used to select question among implicit groups
type DifficultyTag string

const (
	DiffEmpty DifficultyTag = ""
	Diff1     DifficultyTag = "\u2605"             // 1 étoile
	Diff2     DifficultyTag = "\u2605\u2605"       // 2 étoiles
	Diff3     DifficultyTag = "\u2605\u2605\u2605" // 3 étoiles
)

// LevelTag are special question tags used to indicate the
// level (class) for the question.
type LevelTag string

const (
	Seconde   LevelTag = "2NDE" // Seconde
	Premiere  LevelTag = "1ERE" // Première
	Terminale LevelTag = "TERM" // Terminale
)

type Flow uint8

const (
	Parallel   Flow = iota // Questions indépendantes
	Sequencial             // Questions liées
)
