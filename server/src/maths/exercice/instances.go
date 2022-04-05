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
type FormulaDisplayInstance struct {
	Parts []StringOrExpression
}

func (fi FormulaDisplayInstance) toClient() client.Block {
	chunks := make([]string, len(fi.Parts))
	for i, c := range fi.Parts {
		chunks[i] = c.asLaTeX()
	}

	return client.FormulaBlock{Formula: strings.Join(chunks, " ")}
}

type VariationTableInstance struct {
	Xs  []float64 // sorted values for x
	Fxs []float64 // corresponding values for f(x)
}

// inferAlignment return the alignment of the number at index i
func (vt VariationTableInstance) inferAlignment(i int) (isUp bool) {
	if i == len(vt.Xs)-1 { // compute isUp from previous
		return vt.Fxs[i-1] < vt.Fxs[i]
	}
	// else continue from following
	return vt.Fxs[i] > vt.Fxs[i+1]
}

// assume at least two columns
func (vt VariationTableInstance) toClient() client.Block {
	out := client.VariationTableBlock{}
	for i := range vt.Xs {
		numberIsUp := vt.inferAlignment(i)
		// add the number column
		out.Columns = append(out.Columns, client.VariationColumnNumber{
			X:    vt.Xs[i],
			Y:    vt.Fxs[i],
			IsUp: numberIsUp,
		})

		if i <= len(vt.Xs)-1 {
			// add the arrow column
			out.Arrows = append(out.Arrows, !numberIsUp)
		}
	}
	return out
}

type SignTableInstance struct {
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
	return client.SignTableBlock{Columns: columns}
}

type FigureInstance client.FigureBlock

func (f FigureInstance) toClient() client.Block { return client.FigureBlock(f) }

type FunctionGraphInstance struct {
	Function *expression.Expression
	Label    string
	Variable expression.Variable // usually x
	Range    [2]float64          // definition domain
}

func (fg FunctionGraphInstance) toClient() client.Block {
	return client.FunctionGraphBlock{
		Graph: functiongrapher.NewFunctionGraph(fg.Function, fg.Variable, fg.Range[0], fg.Range[1]),
		Label: fg.Label,
	}
}

// FunctionVariationGraphInstance is the same as VariationTableInstance,
// but displays its content as a graph
type FunctionVariationGraphInstance VariationTableInstance

func (fg FunctionVariationGraphInstance) toClient() client.Block {
	return client.FunctionGraphBlock{
		Graph: functiongrapher.GraphFromVariations(fg.Xs, fg.Fxs),
		Label: "y = f(x)",
	}
}

type TableInstance client.TableBlock

func (ti TableInstance) toClient() client.Block {
	return client.TableBlock(ti)
}
