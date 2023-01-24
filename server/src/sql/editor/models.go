// Package editor provides the data structures used to
// create questions and exercices
package editor

import (
	"github.com/benoitkugler/maths-online/server/src/maths/questions"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
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
// gomacro:SQL ADD UNIQUE(Id, NeedExercice)
// TODO: add constraint for empty Parameters for questions in exercices
type Question struct {
	Id IdQuestion

	// only used for question in groups (not in exercices)
	Subtitle   string // used to differentiate questions inside a group
	Difficulty DifficultyTag

	// NeedExercice is not null if the question is part
	// of an exercice and requires its parameters to be
	// instantiated
	NeedExercice OptionalIdExercice `gomacro-sql-foreign:"Exercice"`

	// IdGroup is not null for the standalone questions, accessed by a group
	IdGroup OptionalIdQuestiongroup `gomacro-sql-foreign:"Questiongroup" gomacro-sql-on-delete:"CASCADE"`

	Enonce     questions.Enonce
	Parameters questions.Parameters
}

// Questiongroup groups several variant of the same question
type Questiongroup struct {
	Id        IdQuestiongroup
	Title     string
	Public    bool // in practice only true for admins
	IdTeacher teacher.IdTeacher
}

// gomacro:SQL ADD UNIQUE(IdQuestiongroup, Tag)
// gomacro:SQL ADD CHECK(Tag = upper(Tag))
// gomacro:SQL CREATE UNIQUE INDEX QuestiongroupTag_level ON QuestiongroupTag (IdQuestiongroup) WHERE Section = #[Section.Level]
// gomacro:SQL CREATE UNIQUE INDEX QuestiongroupTag_chapter ON QuestiongroupTag (IdQuestiongroup) WHERE Section = #[Section.Chapter]
type QuestiongroupTag struct {
	Tag             string
	IdQuestiongroup IdQuestiongroup `gomacro-sql-on-delete:"CASCADE"`
	Section         Section
}

// Exercicegroup groups the variant of the same exercice
type Exercicegroup struct {
	Id        IdExercicegroup
	Title     string            // title shown to the student
	Public    bool              // in practice only true for admins
	IdTeacher teacher.IdTeacher // IdTeacher is the owner of the exercice
}

// gomacro:SQL ADD UNIQUE(IdExercicegroup, Tag)
// gomacro:SQL ADD CHECK(Tag = upper(Tag))
// gomacro:SQL CREATE UNIQUE INDEX ExercicegroupTag_level ON ExercicegroupTag (IdExercicegroup) WHERE Section = #[Section.Level]
// gomacro:SQL CREATE UNIQUE INDEX ExercicegroupTag_chapter ON ExercicegroupTag (IdExercicegroup) WHERE Section = #[Section.Chapter]
type ExercicegroupTag struct {
	Tag             string
	IdExercicegroup IdExercicegroup `gomacro-sql-on-delete:"CASCADE"`
	Section         Section
}

// Exercice is the data structure for a full exercice, composed of a list of questions.
// The questions are linked together by a shared Parameters set, and are to be tried sequentially.
type Exercice struct {
	Id      IdExercice
	IdGroup IdExercicegroup

	Subtitle string // subtitle, only shown to the teacher

	// Parameters are parameters shared by all the questions.
	// It is used instead of the Question.Parameters field, which is empty
	Parameters questions.Parameters

	Difficulty DifficultyTag
}

// TODO: check delete question API
// ExerciceQuestion models an ordered list of questions.
// All link items should be updated at once to preserve `Index` invariants
// gomacro:SQL ADD PRIMARY KEY (IdExercice, Index)
// gomacro:SQL ADD FOREIGN KEY (IdExercice, IdQuestion) REFERENCES Questions (NeedExercice, Id)
// gomacro:SQL ADD UNIQUE (IdQuestion)
type ExerciceQuestion struct {
	IdExercice IdExercice `json:"id_exercice" gomacro-sql-on-delete:"CASCADE"`
	IdQuestion IdQuestion `json:"id_question"`
	Bareme     int        `json:"bareme"`
	Index      int        `json:"-"`
}

// Section defines one kind of tag.
type Section uint8

const (
	_        Section = iota
	Level            // Niveau
	Chapter          // Chapitre
	TrivMath         // Triv'Math
)

// DifficultyTag are special question tags used to indicate the
// difficulty of one question.
// It is used to select question among question groups
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
	CPGE      LevelTag = "CPGE" // CPGE
)

type Flow uint8

const (
	Parallel   Flow = iota // Questions indépendantes
	Sequencial             // Questions liées
)
