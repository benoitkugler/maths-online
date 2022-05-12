package exercice

import (
	"fmt"
	"log"
	"strings"

	"github.com/benoitkugler/maths-online/maths/exercice/client"
	"github.com/benoitkugler/maths-online/maths/expression"
	functiongrapher "github.com/benoitkugler/maths-online/maths/function_grapher"
)

type instance interface {
	toClient() client.Block
}

var (
	_ instance = TextInstance{}
	_ instance = FormulaDisplayInstance{}
	_ instance = VariationTableInstance{}
	_ instance = SignTableInstance{}
	_ instance = FigureInstance{}
	_ instance = FunctionVariationGraphInstance{}
	_ instance = TableInstance{}
	_ instance = FunctionGraphInstance{}
)

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
func (qu QuestionInstance) CheckSyntaxe(answer client.QuestionSyntaxCheckIn) client.QuestionSyntaxCheckOut {
	field, ok := qu.fields()[answer.ID]
	if !ok {
		return client.QuestionSyntaxCheckOut{
			ID:     answer.ID,
			Reason: fmt.Sprintf("champ %d inconnu", answer.ID),
		}
	}

	if err := field.validateAnswerSyntax(answer.Answer); err != nil {
		return client.QuestionSyntaxCheckOut{
			ID:     answer.ID,
			Reason: err.(InvalidFieldAnswer).Reason,
		}
	}

	return client.QuestionSyntaxCheckOut{
		ID:      answer.ID,
		IsValid: true,
	}
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
			log.Println("invalid id")
			continue
		}

		if err := reference.validateAnswerSyntax(answer); err != nil {
			log.Println("invalid field", err)
			out[id] = false
			continue
		}

		out[id] = reference.evaluateAnswer(answer)
	}

	return client.QuestionAnswersOut{Data: out}
}

// CorrectAnswer returns the expected answer for the question.
func (qu QuestionInstance) CorrectAnswer() (out client.QuestionAnswersIn) {
	fields := qu.fields()
	out.Data = make(map[int]client.Answer, len(fields))
	for k, v := range fields {
		out.Data[k] = v.correctAnswer()
	}
	return out
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

type EnonceInstance []instance

// StringOrExpression is either an expression or a static string,
// usually rendered as LaTeX, in text mode.
type StringOrExpression struct {
	Expression *expression.Expression
	String     string // LaTeX code, rendered in math mode
}

// IsEmpty returns `true` is the struct is the zero value.
func (se StringOrExpression) IsEmpty() bool {
	return se.Expression == nil && se.String == ""
}

func (fi StringOrExpression) asLaTeX() string {
	if fi.Expression != nil {
		return fi.Expression.AsLaTeX(nil)
	}
	return fi.String
}

// TextInstance is a paragraph of text, which may contain expression or
// math chunks, which is rendered on a single line (eventually wrapped).
type TextInstance client.TextBlock

func (t TextInstance) toClient() client.Block { return client.TextBlock(t) }

// FormulaDisplayInstance is rendered as LaTeX, in display mode.
type FormulaDisplayInstance []string

func (fi FormulaDisplayInstance) toClient() client.Block {
	return client.FormulaBlock{Formula: strings.Join(fi, " ")}
}

// evaluatedExpression groups an expression and its
// result. It is meant to handle cases where we want
// to display 1/3, but a numeric value is also needed.
type evaluatedExpression struct {
	Expr  *expression.Expression
	Value float64
}

func newEvaluatedExpression(s string, params expression.Variables) (evaluatedExpression, error) {
	e, err := expression.Parse(s)
	if err != nil {
		return evaluatedExpression{}, err
	}
	v, err := e.Evaluate(params)
	if err != nil {
		return evaluatedExpression{}, err
	}
	return evaluatedExpression{Expr: e, Value: v}, nil
}

type VariationTableInstance struct {
	Label string
	Xs    []evaluatedExpression // sorted expression values for x
	Fxs   []evaluatedExpression // corresponding values for f(x)
}

// inferAlignment return the alignment of the number at index i
func (vt VariationTableInstance) inferAlignment(i int) (isUp bool) {
	if i == len(vt.Xs)-1 { // compute isUp from previous
		return vt.Fxs[i-1].Value < vt.Fxs[i].Value
	}
	// else, i < len(vt.Xs)-1, compute from following
	return vt.Fxs[i].Value < vt.Fxs[i+1].Value
}

// assume at least two columns
func (vt VariationTableInstance) toClient() client.Block {
	out := client.VariationTableBlock{
		Label: vt.Label,
	}
	for i := range vt.Xs {
		numberIsUp := vt.inferAlignment(i)
		// add the number column
		out.Columns = append(out.Columns, client.VariationColumnNumber{
			X:    vt.Xs[i].Expr.AsLaTeX(nil),
			Y:    vt.Fxs[i].Expr.AsLaTeX(nil),
			IsUp: numberIsUp,
		})

		if i < len(vt.Xs)-1 {
			// add the arrow column
			out.Arrows = append(out.Arrows, !numberIsUp)
		}
	}
	return out
}

type SignTableInstance struct {
	Label     string
	Xs        []string
	FxSymbols []SignSymbol
	Signs     []bool // with length len(Xs) - 1
}

func (st SignTableInstance) toClient() client.Block {
	var columns []client.SignColumn
	for i, x := range st.Xs {
		col := client.SignColumn{
			X:      x,
			IsSign: false,
		}
		switch st.FxSymbols[i] {
		case Nothing:
		case Zero:
			col.IsPositive = true
		case ForbiddenValue:
			col.IsYForbiddenValue = true
		}
		columns = append(columns, col)

		if i != len(st.Xs)-1 {
			columns = append(columns, client.SignColumn{
				IsSign:     true,
				IsPositive: st.Signs[i],
			})
		}
	}
	return client.SignTableBlock{
		Label:   st.Label,
		Columns: columns,
	}
}

type FigureInstance client.FigureBlock

func (f FigureInstance) toClient() client.Block { return client.FigureBlock(f) }

type FunctionGraphInstance struct {
	Functions   []expression.FunctionDefinition
	Decorations []functiongrapher.FunctionDecoration
}

func (fg FunctionGraphInstance) toClient() client.Block {
	return client.FunctionGraphBlock{
		Graph: functiongrapher.NewFunctionGraph(fg.Functions, fg.Decorations),
	}
}

// FunctionVariationGraphInstance is the same as VariationTableInstance,
// but displays its content as a graph
type FunctionVariationGraphInstance VariationTableInstance

func (fg FunctionVariationGraphInstance) values() (xs, fxs []float64) {
	xs = make([]float64, len(fg.Xs))
	fxs = make([]float64, len(fg.Fxs))
	for i, v := range fg.Xs {
		xs[i] = v.Value
	}
	for i, v := range fg.Fxs {
		fxs[i] = v.Value
	}
	return
}

func (fg FunctionVariationGraphInstance) toClient() client.Block {
	xs, fxs := fg.values()
	return client.FunctionGraphBlock{
		Graph: functiongrapher.GraphFromVariations(functiongrapher.FunctionDecoration{Label: fg.Label}, xs, fxs),
	}
}

type TableInstance client.TableBlock

func (ti TableInstance) toClient() client.Block {
	return client.TableBlock(ti)
}
