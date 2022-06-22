package editor

import (
	"database/sql"

	"github.com/benoitkugler/maths-online/maths/questions"
	"github.com/benoitkugler/maths-online/utils"
)

//go:generate ../../../../../structgen/structgen -source=models_sql.go -mode=sql:gen_scans.go -mode=sql_gen:gen_create.sql -mode=rand:gen_randdata_test.go -mode=ts:../../../../prof/src/controller/exercice_gen.ts

// Question is a standalone question, used for instance in games.
type Question struct {
	Id          int64                  `json:"id"`
	Page        questions.QuestionPage `json:"page"`
	Public      bool                   `json:"public"` // in practice only true for admins
	IdTeacher   int64                  `json:"id_teacher"`
	Description string                 `json:"description"`
	// NeedExercice is `true` is the question cannot be instantiated (or edited)
	// on its own
	NeedExercice bool `json:"need_exercice"`
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

// LevelTag are special question tags used to indicate the
// level (class) for the question.
type LevelTag string

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

// IsVisibleBy returns `true` if the question is public or
// owned by `userID`
func (ex Exercice) IsVisibleBy(userID int64) bool {
	return ex.Public || ex.IdTeacher == userID
}

// TODO: check delete question API
// ExerciceQuestion models an ordered list of questions.
// All link items should be updated at once to preserve `Index` invariants
// sql: ADD UNIQUE(id_exercice, index)
type ExerciceQuestion struct {
	IdExercice int64 `json:"id_exercice" sql_on_delete:"CASCADE"`
	IdQuestion int64 `json:"id_question"`
	Bareme     int   `json:"bareme"`
	Index      int   `json:"-" sql:"index"`
}

// updateExerciceQuestionList set the questions for the given exercice,
// overidding `IdExercice` and `index` fields of the list items.
func updateExerciceQuestionList(db *sql.DB, idExercice int64, l ExerciceQuestions) ([]ExerciceQuestionExt, error) {
	// enforce fields
	for i := range l {
		l[i].Index = i
		l[i].IdExercice = idExercice
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, utils.SQLError(err)
	}
	_, err = DeleteExerciceQuestionsByIdExercices(db, idExercice)
	if err != nil {
		_ = tx.Rollback()
		return nil, utils.SQLError(err)
	}

	err = InsertManyExerciceQuestions(tx, l...)
	if err != nil {
		_ = tx.Rollback()
		return nil, utils.SQLError(err)
	}

	questions, err := SelectQuestions(tx, l.IdQuestions()...)
	if err != nil {
		_ = tx.Rollback()
		return nil, utils.SQLError(err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, utils.SQLError(err)
	}

	return fillQuestions(l, questions), nil
}
