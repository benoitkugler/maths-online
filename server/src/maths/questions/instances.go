package questions

import (
	"fmt"
	"log"
	"strings"

	"github.com/benoitkugler/maths-online/server/src/maths/expression"
	"github.com/benoitkugler/maths-online/server/src/maths/functiongrapher"
	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
)

type instance interface {
	toClient() client.Block

	// returns LaTeX code describing the block as close as possible
	toLatex() string
}

var (
	_ instance = TextInstance{}
	_ instance = FormulaDisplayInstance{}
	_ instance = VariationTableInstance{}
	_ instance = SignTableInstance{}
	_ instance = FigureInstance{}
	_ instance = TableInstance{}
	_ instance = FunctionsGraphInstance{}
	_ instance = ImageInstance{}
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
	Enonce     EnonceInstance
	Correction EnonceInstance
}

func (qu EnonceInstance) fields() map[int]fieldInstance {
	out := make(map[int]fieldInstance)
	for _, block := range qu {
		if field, isField := block.(fieldInstance); isField {
			out[field.fieldID()] = field
		}
	}
	return out
}

// CheckSyntaxe returns an error message if the syntaxe is not
func (qu EnonceInstance) CheckSyntaxe(answer client.QuestionSyntaxCheckIn) client.QuestionSyntaxCheckOut {
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
// An empty [answers] is supported, corresponding to the case where the student
// has left the question.
func (qu EnonceInstance) EvaluateAnswer(answers client.QuestionAnswersIn) client.QuestionAnswersOut {
	fields := qu.fields()

	out := client.QuestionAnswersOut{
		Results:         make(map[int]bool, len(fields)),
		ExpectedAnswers: make(map[int]client.Answer, len(fields)),
	}

	for id, reference := range fields {
		out.ExpectedAnswers[id] = reference.correctAnswer()
		out.Results[id] = false

		answer := answers.Data[id]
		if answer == nil {
			// the field was not provided, skip verification
			continue
		}

		if err := reference.validateAnswerSyntax(answer); err != nil {
			log.Printf("internal error: invalid field syntax for %T: %s", reference, err)
			continue
		}

		out.Results[id] = reference.evaluateAnswer(answer)
	}

	return out
}

// CorrectAnswer returns the expected answer for the question.
func (qu EnonceInstance) CorrectAnswer() (out client.QuestionAnswersIn) {
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
		Enonce:     make(client.Enonce, len(qi.Enonce)),
		Correction: make(client.Enonce, len(qi.Correction)),
	}
	for i, c := range qi.Enonce {
		out.Enonce[i] = c.toClient()
	}
	for i, c := range qi.Correction {
		out.Correction[i] = c.toClient()
	}
	return out
}

type EnonceInstance []instance

// StringOrExpression is either an expression or a static string,
// usually rendered as LaTeX, in text mode.
type StringOrExpression struct {
	Expression *expression.Expr
	String     string // LaTeX code, rendered in math mode
}

// IsEmpty returns `true` is the struct is the zero value.
func (se StringOrExpression) IsEmpty() bool {
	return se.Expression == nil && se.String == ""
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
	Expr  *expression.Expr
	Value float64
}

// subsitute variables
func newEvaluatedExpression(s string, params expression.Vars) (evaluatedExpression, error) {
	e, err := expression.Parse(s)
	if err != nil {
		return evaluatedExpression{}, err
	}
	e.Substitute(params)
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

// inferNumberAlignment return the alignment of the number at index i
// the arrow at index i is has then the opposite direction
func (vt VariationTableInstance) inferNumberAlignment(i int) (isUp bool) {
	if i == len(vt.Xs)-1 { // compute isUp from previous
		return vt.Fxs[i-1].Value < vt.Fxs[i].Value
	}
	// else, i < len(vt.Xs)-1, compute from following
	return vt.Fxs[i].Value > vt.Fxs[i+1].Value
}

// assume at least two columns
func (vt VariationTableInstance) toClient() client.Block {
	out := client.VariationTableBlock{
		Label: vt.Label,
	}
	for i := range vt.Xs {
		numberIsUp := vt.inferNumberAlignment(i)
		// add the number column
		out.Columns = append(out.Columns, client.VariationColumnNumber{
			X:    vt.Xs[i].Expr.AsLaTeX(),
			Y:    vt.Fxs[i].Expr.AsLaTeX(),
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
	Xs        []*expression.Expr
	Functions []client.FunctionSign
}

func (st SignTableInstance) toClient() client.Block {
	out := client.SignTableBlock{
		Xs:        make([]string, len(st.Xs)),
		Functions: append([]client.FunctionSign(nil), st.Functions...),
	}
	for i, x := range st.Xs {
		out.Xs[i] = x.AsLaTeX()
	}
	return out
}

type FigureInstance client.FigureBlock

func (f FigureInstance) toClient() client.Block { return client.FigureBlock(f) }

type FunctionsGraphInstance struct {
	Functions []functiongrapher.FunctionGraph
	Sequences []functiongrapher.SequenceGraph
	Areas     []client.FunctionArea
	Points    []client.FunctionPoint
}

func (fg FunctionsGraphInstance) toClient() client.Block { return fg.toClientG() }

func (fg FunctionsGraphInstance) toClientG() client.FunctionsGraphBlock {
	// compute the required bounding box for all drawings
	var allSegments []functiongrapher.BezierCurve
	for _, fn := range fg.Functions {
		allSegments = append(allSegments, fn.Segments...)
	}
	for _, seq := range fg.Sequences {
		for _, point := range seq.Points {
			allSegments = append(allSegments, functiongrapher.BezierCurve{P0: point, P1: point, P2: point})
		}
	}

	return client.FunctionsGraphBlock{
		Functions: fg.Functions,
		Sequences: fg.Sequences,
		Areas:     fg.Areas,
		Points:    fg.Points,
		Bounds:    functiongrapher.BoundingBox(allSegments),
	}
}

type TableInstance client.TableBlock

func (ti TableInstance) toClient() client.Block {
	return client.TableBlock(ti)
}

type TreeNodeInstance struct {
	Children      []TreeNodeInstance
	Probabilities []*expression.Expr // expression for edges, same length as Children
	Value         int                // index into the proposals, ignored for the root
}

func (node TreeNodeInstance) toClient() client.TreeNodeAnswer {
	out := client.TreeNodeAnswer{
		Value:         node.Value,
		Probabilities: make([]string, len(node.Probabilities)),
		Children:      make([]client.TreeNodeAnswer, len(node.Children)),
	}
	for i, expr := range node.Probabilities {
		out.Probabilities[i] = expr.String() // the client does not support latex in edges
	}
	for i, child := range node.Children { // recurse
		out.Children[i] = child.toClient()
	}
	return out
}

type TreeInstance struct {
	EventsProposals []client.TextLine
	AnswerRoot      TreeNodeInstance
}

func (ti TreeInstance) toClient() client.Block {
	return client.TreeBlock{
		EventsProposals: ti.EventsProposals,
		Root:            ti.AnswerRoot.toClient(),
	}
}

type ImageInstance ImageBlock

func (img ImageInstance) toClient() client.Block { return client.ImageBlock(img) }
