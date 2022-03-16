package exercice

import (
	"fmt"
	"strings"

	"github.com/benoitkugler/maths-online/maths/exercice/client"
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
	StaticContent string // LaTeX code, rendered in math mode
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
		out.Questions[i] = qu.instantiate(params)
	}

	return out
}

func (qu Question) instantiate(params expression.Variables) QuestionInstance {
	enonce := make(EnonceInstance, len(qu.Enonce))
	var currentID int
	for j, bl := range qu.Enonce {
		enonce[j] = bl.instantiate(params, currentID)
		if _, isField := enonce[j].(fieldInstance); isField {
			currentID++
		}
	}
	return QuestionInstance{Title: qu.Title, Enonce: enonce}
}

func (t TextBlock) instantiate(params expression.Variables, _ int) blockInstance {
	return t
}

func (f Formula) instantiate(params expression.Variables, _ int) blockInstance {
	out := FormulaInstance{IsInline: f.IsInline}
	out.Chunks = make([]FormulaPartInstance, len(f.Chunks))
	for i, c := range f.Chunks {
		out.Chunks[i] = c.instantiate(params)
	}
	return out
}

func (n NumberField) instantiate(params expression.Variables, ID int) blockInstance {
	expr, _, _ := expression.Parse(n.Expression)
	answer := expr.Evaluate(params)
	return NumberFieldInstance{ID: ID, Answer: answer}
}

// TODO
func (f FormulaField) instantiate(params expression.Variables, ID int) blockInstance {
	return ExpressionFieldInstance{}
}

// TODO
func (l ListField) instantiate(params expression.Variables, ID int) blockInstance {
	return RadioFieldInstance{}
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
	Title  string
	Enonce EnonceInstance
}

func (qu QuestionInstance) fields() map[int]fieldInstance {
	out := make(map[int]fieldInstance)
	for _, block := range qu.Enonce {
		if field, isField := block.(fieldInstance); isField {
			out[field.fieldID()] = field
		}
	}
	return out
}

// CheckSyntaxe returns an error message if the syntaxe is not
func (qu QuestionInstance) CheckSyntaxe(answer client.QuestionSyntaxCheckIn) error {
	field, ok := qu.fields()[answer.ID]
	if !ok {
		return InvalidFieldAnswer{
			ID:     answer.ID,
			Reason: fmt.Sprintf("champ %d inconnu", answer.ID),
		}
	}

	return field.validateAnswerSyntax(answer.Answer)
}

// EvaluateAnswer check if the given answers are correct, and complete.
// TODO: provide detailled feedback
func (qu QuestionInstance) EvaluateAnswer(answers client.QuestionAnswersIn) client.QuestionAnswersOut {
	fields := qu.fields()

	out := make(map[int]bool, len(fields))
	for id, reference := range fields {
		answer, ok := answers.Data[id]
		if !ok { // should not happen since the client forces the user to fill all fields
			out[id] = false
			continue
		}

		if err := reference.validateAnswerSyntax(answer); err != nil {
			out[id] = false // should not happen since a pre-validation step has been done previously
			continue
		}

		out[id] = reference.evaluateAnswer(answer)
	}

	return client.QuestionAnswersOut{Data: out}
}

// ToClient convert the question to a client version, stripping
// expected answers and converting expressions to LaTeX strings.
func (qi QuestionInstance) ToClient() client.Question {
	out := client.Question{
		Title:  qi.Title,
		Enonce: make(client.Enonce, len(qi.Enonce)),
	}
	for i, c := range qi.Enonce {
		out.Enonce[i] = c.toClient()
	}
	return out
}

type EnonceInstance []blockInstance

type blockInstance interface {
	toClient() client.Block
}

func (t TextBlock) toClient() client.Block { return client.TextBlock(t) }

type FormulaInstance struct {
	Chunks   []FormulaPartInstance
	IsInline bool
}

func (fi FormulaInstance) toClient() client.Block {
	chunks := make([]string, len(fi.Chunks))
	for i, c := range fi.Chunks {
		chunks[i] = c.asLaTeX()
	}

	return client.FormulaBlock{Content: strings.Join(chunks, " "), IsInline: fi.IsInline}
}
