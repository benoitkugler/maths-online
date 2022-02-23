package exercice

// This file is used a source to auto generate dart bindings

// Exercice is a sequence of questions
type Exercice struct {
	Title       string            // name for the exercice
	Description string            // overall description for all questions (optional)
	Fields      map[string]string // TODO: the logic needed by each field
	Questions   []Question        // the actual content of the exercices
}

type Question struct {
	Content []block
}

type block interface {
	isBlock()
}

func (TextBlock) isBlock()    {}
func (Formula) isBlock()      {}
func (ListField) isBlock()    {}
func (FormulaField) isBlock() {}

type TextBlock struct {
	Text string
}

type Formula struct {
	Latex    string
	IsInline bool // else display
}

type ListField struct {
	Id      string
	Choices []string
}

type FormulaField struct {
	Id string
}
