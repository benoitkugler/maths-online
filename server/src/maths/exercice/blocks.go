package exercice

import (
	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/repere"
)

const ExhaustiveTextKind = "exhaustiveTextKind"

var (
	_ Block = TextBlock{}
	_ Block = FormulaBlock{}
	_ Block = VariationTableBlock{}
	_ Block = SignTableBlock{}
	_ Block = FigureBlock{}
)

type Enonce []Block

// Block form the actual content of a question
// it is stored in a DB in generic form, but may be instantiated
// against random parameter values
type Block interface {
	// ID is only used by answer fields
	instantiate(params expression.Variables, ID int) instance
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

type TextPart struct {
	Content string
	Kind    TextKind
}

func (tp TextPart) instantiate(params expression.Variables) TextOrMaths {
	switch tp.Kind {
	case Text:
		return TextOrMaths{StringOrExpression: StringOrExpression{String: tp.Content}}
	case StaticMath:
		return TextOrMaths{StringOrExpression: StringOrExpression{String: tp.Content}, IsMath: true}
	case Expression:
		expr, _, _ := expression.Parse(tp.Content)
		expr.Substitute(params)
		return TextOrMaths{StringOrExpression: StringOrExpression{Expression: expr}, IsMath: true}
	default:
		panic(ExhaustiveTextKind)
	}
}

type TextParts []TextPart

// TextBlock is a chunk of text
// which may contain maths
// It support basic interpolation syntax.
type TextBlock struct {
	Parts  string
	IsHint bool
}

func (t TextBlock) instantiate(params expression.Variables, _ int) instance {
	content, _ := ParseInterpolatedString(t.Parts)
	parts := make([]TextOrMaths, len(content))
	for i, p := range content {
		parts[i] = p.instantiate(params)
	}
	return TextInstance{
		IsHint: t.IsHint,
		Parts:  parts,
	}
}

// FormulaContent is a list of chunks, either
//	- static math symbols, such as f(x) =
//	- valid expression, such as a*x - b, which will be instantiated
// when rendering the question
//
// For instance, the formula "f(x) = a*(x + 2)"
// is represented by two FormulaPart elements:
// 	{ f(x) = } and { a*(x + 2) }
type FormulaContent []FormulaPart

// FormulaPart forms a logic chunk of a formula.
type FormulaPart struct {
	Content      string
	IsExpression bool // when true, Content is interpreted as an expression.Expression
}

// assume the expression is valid
func (fp FormulaPart) instantiate(params expression.Variables) StringOrExpression {
	if !fp.IsExpression { // nothing to do
		return StringOrExpression{String: fp.Content}
	}

	expr, _, _ := expression.Parse(fp.Content)
	expr.Substitute(params)
	return StringOrExpression{Expression: expr}
}

// FormulaBlock is a math formula, which should be display using
// a LaTeX renderer.
type FormulaBlock struct {
	Parts FormulaContent
}

func (f FormulaBlock) instantiate(params expression.Variables, _ int) instance {
	out := FormulaDisplayInstance{}
	out.Parts = make([]StringOrExpression, len(f.Parts))
	for i, c := range f.Parts {
		out.Parts[i] = c.instantiate(params)
	}
	return FormulaDisplayInstance{}
}

type VariationTableBlock struct {
	Xs  []string // expressions
	Fxs []string // expressions
}

func (vt VariationTableBlock) instantiate(params expression.Variables, _ int) instance {
	out := VariationTableInstance{
		Xs:  make([]float64, len(vt.Xs)),
		Fxs: make([]float64, len(vt.Fxs)),
	}
	for i, c := range vt.Xs {
		out.Xs[i] = mustEvaluate(c, params)
	}
	for i, c := range vt.Fxs {
		out.Fxs[i] = mustEvaluate(c, params)
	}
	return out
}

type SignTableBlock struct {
	Xs        FormulaContent
	FxSymbols []SignSymbol
	Signs     []bool // with length len(Xs) - 1
}

func (st SignTableBlock) instantiate(params expression.Variables, _ int) instance {
	out := SignTableInstance{
		Xs: make([]string, len(st.Xs)),
	}
	for i, c := range st.Xs {
		out.Xs[i] = c.instantiate(params).asLaTeX()
	}
	out.FxSymbols = append([]SignSymbol(nil), st.FxSymbols...)
	out.Signs = append([]bool(nil), st.Signs...)
	return out
}

type FigureBlock struct {
	Drawings repere.RandomDrawings
	Bounds   repere.RepereBounds
	ShowGrid bool
}

func (f FigureBlock) instantiate(params expression.Variables, _ int) instance {
	out := FigureInstance{
		Figure: repere.Figure{
			Drawings: repere.Drawings{
				Segments: f.Drawings.Segments,
				Points:   make(map[string]repere.LabeledPoint),
				Lines:    make([]repere.Line, len(f.Drawings.Lines)),
			},
			Bounds:   f.Bounds,
			ShowGrid: f.ShowGrid,
		},
	}
	for k, v := range f.Drawings.Points {
		out.Figure.Drawings.Points[k] = repere.LabeledPoint{
			Point: repere.Coord{
				X: mustEvaluate(v.Coord.X, params),
				Y: mustEvaluate(v.Coord.Y, params),
			},
			Pos: v.Pos,
		}
	}

	for i, l := range f.Drawings.Lines {
		out.Figure.Drawings.Lines[i] = repere.Line{
			Label: l.Label,
			A:     mustEvaluate(l.A, params),
			B:     mustEvaluate(l.B, params),
		}
	}
	return out
}

func (n NumberField) instantiate(params expression.Variables, ID int) instance {
	expr, _, _ := expression.Parse(n.Expression)
	answer, _ := expr.Evaluate(params)
	return NumberFieldInstance{ID: ID, Answer: answer}
}

// TODO
func (f FormulaField) instantiate(params expression.Variables, ID int) instance {
	return ExpressionFieldInstance{}
}

// TODO
func (l ListField) instantiate(params expression.Variables, ID int) instance {
	return RadioFieldInstance{}
}
