package exercice

// This file is used as source to auto generate SQL statements

//go:generate ../../../../../structgen/structgen -source=sql_models.go -mode=sql:gen_scans.go -mode=sql_test:gen_scans_test.go  -mode=sql_composite:composites/auto.go  -mode=sql_gen:create_gen.sql  -mode=rand:gen_data_test.go -mode=itfs-json:gen_itfs.go

// Exercice is a sequence of questions
type Exercice struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`       // name for the exercice
	Description string `json:"description"` // overall description for all questions (optional)
}

// Question is the fundamental object to build exercices.
// It is mainly consituted of a list of content blocks, which
// describes the question (description, question, field answer)
type Question struct {
	IdExercice int64   `json:"id_exercice" sql_on_delete:"CASCADE"`
	Content    Content `json:"content"`
}

type Content []block

type block interface {
	isBlock()
}

func (TextBlock) isBlock()    {}
func (Formula) isBlock()      {}
func (ListField) isBlock()    {}
func (FormulaField) isBlock() {}

// TextBlock is a regular chunk of text
type TextBlock struct {
	Text string
}

// Formula is a math formula, which should be display using
// a LaTeX renderer
type Formula struct {
	Latex    string
	IsInline bool // else display
}

type ListField struct {
	// Id      string
	Choices []string
}

type FormulaField struct {
	// Id string
	Content string
}
