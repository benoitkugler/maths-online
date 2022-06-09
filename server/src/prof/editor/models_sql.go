package editor

import "github.com/benoitkugler/maths-online/maths/exercice"

//go:generate ../../../../../structgen/structgen -source=models_sql.go -mode=sql:gen_scans.go -mode=sql_gen:gen_create.sql -mode=rand:gen_randdata_test.go -mode=ts:../../../../prof/src/controller/exercice_gen.ts

// Question is a standalone question, used for instance in games.
type Question struct {
	Id          int64                 `json:"id"`
	Page        exercice.QuestionPage `json:"page"`
	Public      bool                  `json:"public"` // in practice only true for admins
	IdTeacher   int64                 `json:"id_teacher"`
	Description string                `json:"description"`
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
