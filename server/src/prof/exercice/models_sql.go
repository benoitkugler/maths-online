package exercice

import "github.com/benoitkugler/maths-online/maths/questions"

//go:generate ../../../../../structgen/structgen -source=models_sql.go -mode=sql:gen_scans.go -mode=sql_gen:gen_create.sql -mode=rand:gen_randdata_test.go

type Exercice struct {
	Id int64
	// Parameters are parameters shared by all the questions,
	// which are added to the individual ones.
	// It will be empty for parallel exercices
	Parameters questions.Parameters
	Flow       Flow
}

// TODO: check delete question API
type ExerciceQuestion struct {
	IdExercice int64 `json:"id_exercice" sql_on_delete:"CASCADE"`
	IdQuestion int64 `json:"id_question"`
	Bareme     int   `json:"bareme"`
}
