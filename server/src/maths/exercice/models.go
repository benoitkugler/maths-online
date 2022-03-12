package exercice

import (
	"strings"

	"github.com/benoitkugler/maths-online/maths/expression"
)

// randomParameters is a serialized form of expression.RandomParameters
type randomParameters []randomParameter

type randomParameter struct {
	Expression string `json:"expression"` // as typed by the user, but validated
	Variable   rune   `json:"variable"`
}

// toMap assumes `rp` only contains valid expressions
func (rp randomParameters) toMap() expression.RandomParameters {
	out := make(expression.RandomParameters, len(rp))
	for _, item := range rp {
		out[expression.Variable(item.Variable)], _, _ = expression.Parse(item.Expression)
	}
	return out
}

// FormulaContent is a list of chunks, either
//	- static math symbols, such as f(x) =
//	- valid expression, such as a*x - b, which will be instantiated when rendering the question
//
// For instance, the formula f(x) = a*(x + 2)
// is represented by two FormulaPart elements:
// 	{ f(x) = } and { a*(x + 2) }
type FormulaContent []FormulaPart

// FormulaPart forms a logic chunk of a formula.
type FormulaPart struct {
	Content      string
	IsExpression bool // when true, Content is interpreted as an expression.Expression
}

// assume the expression is valid
func (fp FormulaPart) instantiate(params expression.Variables) FormulaPartInstance {
	if !fp.IsExpression { // nothing to do
		return FormulaPartInstance{StaticContent: fp.Content}
	}

	expr, _, _ := expression.Parse(fp.Content)
	expr.Substitute(params)
	return FormulaPartInstance{Expression: expr}
}

// FormulaPartInstance is either an expression or a static string
type FormulaPartInstance struct {
	Expression    *expression.Expression
	StaticContent string
}

func (fi FormulaPartInstance) asLaTeX() string {
	if fi.Expression != nil {
		return fi.Expression.AsLaTeX(nil)
	}
	return fi.StaticContent
}

// Exercice is a sequence of questions
type ExerciceQuestions struct {
	Exercice
	Questions Questions
}

// instantiate returns a deep copy of `eq`, where all random parameters
// have been resolved
// It assumes that the expressions and random parameters definitions are valid.
func (eq ExerciceQuestions) instantiate() ExerciceInstance {
	rp := eq.RandomParameters.toMap()
	// generate random params
	params, _ := rp.Instantiate()

	out := ExerciceInstance{
		Id:          eq.Id,
		Title:       eq.Title,
		Description: eq.Description,
	}
	out.Questions = make([]QuestionInstance, len(eq.Questions))
	for i, qu := range eq.Questions {
		content := make(ContentInstance, len(qu.Content))
		for j, bl := range qu.Content {
			content[j] = bl.instantiate(params)
		}

		out.Questions[i].Content = content
	}

	return out
}

func (t TextBlock) instantiate(params expression.Variables) blockInstance {
	return t
}

func (f Formula) instantiate(params expression.Variables) blockInstance {
	out := FormulaInstance{IsInline: f.IsInline}
	out.Chunks = make([]FormulaPartInstance, len(f.Chunks))
	for i, c := range f.Chunks {
		out.Chunks[i] = c.instantiate(params)
	}
	return out
}

// TODO
func (l ListField) instantiate(params expression.Variables) blockInstance {
	return ListFieldInstance{}
}

// TODO
func (f FormulaField) instantiate(params expression.Variables) blockInstance {
	return FormulaFieldInstance{}
}

// ExerciceInstance is an in memory version of an Exercice,
// where all random parameters have been generated and substituted
type ExerciceInstance struct {
	Title       string
	Description string
	Questions   []QuestionInstance
	Id          int64
}

type QuestionInstance struct {
	Title   string
	Content ContentInstance
}

func (qi QuestionInstance) toClient() ClientQuestion {
	out := ClientQuestion{
		Title:   qi.Title,
		Content: make(ClientContent, len(qi.Content)),
	}
	for i, c := range qi.Content {
		out.Content[i] = c.toClient()
	}
	return out
}

type ContentInstance []blockInstance

type blockInstance interface {
	toClient() clientBlock
}

func (t TextBlock) toClient() clientBlock { return textBlock(t) }

type FormulaInstance struct {
	Chunks   []FormulaPartInstance
	IsInline bool
}

func (fi FormulaInstance) toClient() clientBlock {
	chunks := make([]string, len(fi.Chunks))
	for i, c := range fi.Chunks {
		chunks[i] = c.asLaTeX()
	}

	return formulaBlock{Content: strings.Join(chunks, " "), IsInline: fi.IsInline}
}

type ListFieldInstance struct{}

// TODO:
func (ListFieldInstance) toClient() clientBlock { return clientListFieldBlock{} }

type FormulaFieldInstance struct{}

// TODO:
func (FormulaFieldInstance) toClient() clientBlock { return clientFormulaFieldBlock{} }
