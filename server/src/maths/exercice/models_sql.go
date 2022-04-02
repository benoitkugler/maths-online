package exercice

import "github.com/benoitkugler/maths-online/maths/expression"

// This file is used as source to auto generate SQL statements

//go:generate ../../../../../structgen/structgen -source=models_sql.go -mode=sql:gen_scans.go -mode=sql_test:gen_scans_test.go  -mode=sql_composite:composites/auto.go  -mode=sql_gen:create_gen.sql  -mode=rand:gen_data_test.go -mode=itfs-json:gen_itfs.go

type Exercice struct {
	Id               int64            `json:"id"`
	Title            string           `json:"title"`             // name for the exercice
	Description      string           `json:"description"`       // overall description for all questions (optional)
	RandomParameters randomParameters `json:"random_parameters"` // random parameters shared by the questions
}

// Question is the fundamental object to build exercices.
// It is mainly consituted of a list of content blocks, which
// describes the question (description, question, field answer)
type Question struct {
	IdExercice int64  `json:"id_exercice" sql_on_delete:"CASCADE"`
	Title      string `json:"title"` // theme of the question
	Enonce     Enonce `json:"enonce"`
}

type Enonce []block

// block form the actual content of a question
// it is stored in a DB in generic form, but may be instantiated
// against random parameter values
type block interface {
	// ID is only used by answer fields
	instantiate(params expression.Variables, ID int) instance
}

// TextBlock is a chunk of text
// which may contain maths
type TextBlock struct {
	Parts  TextParts
	IsHint bool
}

// FormulaBlock is a math formula, which should be display using
// a LaTeX renderer.
type FormulaBlock struct {
	Parts FormulaContent
}

type VariationTableBlock struct {
	Xs  []string // expressions
	Fxs []string // expressions
}

type SignTableBlock struct {
	Xs        FormulaContent
	FxSymbols []SignSymbol
	Signs     []bool // with length len(Xs) - 1
}

type NumberField struct {
	// a valid expression, in the format used by expression.Expression
	// which is only parametrized by the random parameters
	// TODO: carefully check that the prof expression is valid
	Expression string
}

type ListField struct {
	Choices []string
}

type FormulaField struct {
	Expression string // a valid expression, in the format used by expression.Expression
}
