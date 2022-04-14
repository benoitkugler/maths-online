package exercice

// This file is used as source to auto generate SQL statements

//go:generate ../../../../../structgen/structgen -source=models_sql.go -mode=sql:gen_scans.go -mode=sql_test:gen_scans_test.go  -mode=sql_composite:composites/auto.go  -mode=sql_gen:create_gen.sql  -mode=rand:gen_data_test.go -mode=itfs-json:gen_itfs.go -mode=ts:../../../../prof/src/controller/exercice_gen.ts

// Question is the fundamental object to build exercices.
// It is mainly consituted of a list of content blocks, which
// describes the question (description, question, field answer)
type Question struct {
	Id         int64      `json:"id"`
	Title      string     `json:"title"` // name of the question, optional
	Enonce     Enonce     `json:"enonce"`
	Parameters Parameters `json:"parameters"` // random parameters shared by the all the blocks
}

// sql: ADD UNIQUE(id_question, tag)
type QuestionTag struct {
	Tag        string `json:"tag"`
	IdQuestion int64  `sql_on_delete:"CASCADE" json:"id_question"`
}
