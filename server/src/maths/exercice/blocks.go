package exercice

import (
	"strings"

	"github.com/benoitkugler/maths-online/maths/exercice/client"
	"github.com/benoitkugler/maths-online/maths/expression"
	functiongrapher "github.com/benoitkugler/maths-online/maths/function_grapher"
	"github.com/benoitkugler/maths-online/maths/repere"
)

const ExhaustiveTextKind = "exhaustiveTextKind"

var (
	_ Block = TextBlock{}
	_ Block = FormulaBlock{}
	_ Block = VariationTableBlock{}
	_ Block = SignTableBlock{}
	_ Block = FigureBlock{}
	_ Block = FunctionGraphBlock{}
	_ Block = FunctionVariationGraphBlock{}
	_ Block = TableBlock{}
)

type Enonce []Block

// Block form the actual content of a question
// it is stored in a DB in generic form, but may be instantiated
// against random parameter values
type Block interface {
	// ID is only used by answer fields
	instantiate(params expression.Variables, ID int) instance
}

type Parameters struct {
	Variables  randomParameters
	Intrinsics []string // validated by exercice.ParseIntrinsic
}

// ToMap may only be used after `Validate`
func (pr Parameters) ToMap() expression.RandomParameters {
	// start with basic variables
	out := pr.Variables.toMap()

	// add intrinsics
	for _, intrinsic := range pr.Intrinsics {
		it, _ := expression.ParseIntrinsic(intrinsic)
		_ = it.MergeTo(out)
	}

	return out
}

// randomParameters is a serialized form of expression.RandomParameters
type randomParameters []randomParameter

type randomParameter struct {
	Expression string              `json:"expression"` // as typed by the user, but validated
	Variable   expression.Variable `json:"variable"`
}

// toMap assumes `rp` only contains valid expressions,
// or it will panic
// It may only be used after `ValidateParameters`
func (rp randomParameters) toMap() expression.RandomParameters {
	out := make(expression.RandomParameters, len(rp))
	for _, item := range rp {
		out[expression.Variable(item.Variable)] = mustParse(item.Expression)
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
	rp := eq.Parameters.ToMap()
	// generate random params
	params, _ := rp.Instantiate()

	out := ExerciceInstance{
		Id:          eq.Id,
		Title:       eq.Title,
		Description: eq.Description,
	}
	out.Questions = make([]QuestionInstance, len(eq.Questions))

	for i, qu := range eq.Questions {
		out.Questions[i] = qu.instantiateWith(params)
	}

	return out
}

// Instantiate returns a deep copy of `qu`, where all random parameters
// have been resolved.
// It assumes that the expressions and random parameters definitions are valid.
func (qu Question) Instantiate() QuestionInstance {
	// generate random params
	rp, _ := qu.Parameters.ToMap().Instantiate()
	return qu.instantiateWith(rp)
}

func (qu Question) instantiateWith(params expression.Variables) QuestionInstance {
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

// TextPart is either a plain text, a LaTeX code or an expression
type TextPart struct {
	Content string
	Kind    TextKind
}

func (tp TextPart) instantiate(params expression.Variables) client.TextOrMath {
	switch tp.Kind {
	case Text:
		return client.TextOrMath{Text: tp.Content}
	case StaticMath:
		return client.TextOrMath{Text: tp.Content, IsMath: true}
	case Expression:
		expr, _, _ := expression.Parse(tp.Content)
		expr.Substitute(params)
		return client.TextOrMath{Text: expr.AsLaTeX(nil), IsMath: true}
	default:
		panic(ExhaustiveTextKind)
	}
}

func (tp TextPart) instantiateAndEvaluate(params expression.Variables) client.TextOrMath {
	if tp.Kind == Expression {
		expr, _, _ := expression.Parse(tp.Content)
		expr.Substitute(params)
		v, err := expr.Evaluate(nil)
		if err == nil {
			expr = expression.NewNumber(v)
		}
		return client.TextOrMath{Text: expr.AsLaTeX(nil), IsMath: true}
	}
	return tp.instantiate(params)
}

type TextParts []TextPart

func (tp TextParts) instantiate(params expression.Variables) []client.TextOrMath {
	parts := make([]client.TextOrMath, len(tp))
	for i, p := range tp {
		parts[i] = p.instantiate(params)
	}
	return parts
}

// assume all parts are either static math or expression.
// after instantiating, tries to evaluate the expression
// and returns the LaTeX concatenated code
func (tp TextParts) instantiateAndEvaluate(params expression.Variables) string {
	parts := make([]string, len(tp))
	for i, p := range tp {
		parts[i] = p.instantiateAndEvaluate(params).Text
	}
	return strings.Join(parts, "")
}

// TextBlock is a chunk of text
// which may contain maths
// It support basic interpolation syntax.
type TextBlock struct {
	Parts  Interpolated
	IsHint bool
}

func (t TextBlock) instantiate(params expression.Variables, _ int) instance {
	content, _ := t.Parts.Parse()
	return TextInstance{
		IsHint: t.IsHint,
		Parts:  content.instantiate(params),
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
	Parts Interpolated
}

func (f FormulaBlock) instantiate(params expression.Variables, _ int) instance {
	parts, _ := f.Parts.Parse()
	partsInstance := parts.instantiate(params)
	out := make(FormulaDisplayInstance, len(partsInstance))
	for i, c := range partsInstance {
		out[i] = c.Text
	}
	return out
}

type VariationTableBlock struct {
	Xs  []string // expressions
	Fxs []string // expressions
}

func (vt VariationTableBlock) instantiateVT(params expression.Variables) VariationTableInstance {
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

func (vt VariationTableBlock) instantiate(params expression.Variables, _ int) instance {
	return vt.instantiateVT(params)
}

type SignTableBlock struct {
	Xs        []Interpolated // always math content
	FxSymbols []SignSymbol
	Signs     []bool // with length len(Xs) - 1
}

func (st SignTableBlock) instantiate(params expression.Variables, _ int) instance {
	out := SignTableInstance{
		Xs: make([]string, len(st.Xs)),
	}
	for i, c := range st.Xs {
		parts, _ := c.Parse()
		out.Xs[i] = parts.instantiateAndEvaluate(params)
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
	return f.instantiateF(params)
}

func (f FigureBlock) instantiateF(params expression.Variables) FigureInstance {
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
	for _, v := range f.Drawings.Points {
		out.Figure.Drawings.Points[v.Name] = repere.LabeledPoint{
			Point: repere.Coord{
				X: mustEvaluate(v.Point.Coord.X, params),
				Y: mustEvaluate(v.Point.Coord.Y, params),
			},
			Pos: v.Point.Pos,
		}
	}

	for i, l := range f.Drawings.Lines {
		out.Figure.Drawings.Lines[i] = repere.Line{
			Label: l.Label,
			A:     mustEvaluate(l.A, params),
			B:     mustEvaluate(l.B, params),
			Color: l.Color,
		}
	}
	return out
}

type FunctionDefinition struct {
	Function   string // expression.Expression
	Decoration functiongrapher.FunctionDecoration
	Variable   expression.Variable // usually x
	Range      [2]float64          // definition domain
}

func (fg FunctionDefinition) instantiate(params expression.Variables) expression.FunctionDefinition {
	expr := mustParse(fg.Function)
	expr.Substitute(params)
	return expression.FunctionDefinition{
		FunctionExpr: expression.FunctionExpr{
			Function: expr,
			Variable: fg.Variable,
		},
		From: fg.Range[0],
		To:   fg.Range[1],
	}
}

// TODO: validate using extrema to avoid edge cases
type FunctionGraphBlock struct {
	Functions []FunctionDefinition
}

func (fg FunctionGraphBlock) instantiate(params expression.Variables, _ int) instance {
	out := FunctionGraphInstance{
		Functions:   make([]expression.FunctionDefinition, len(fg.Functions)),
		Decorations: make([]functiongrapher.FunctionDecoration, len(fg.Functions)),
	}
	for i, f := range fg.Functions {
		out.Functions[i] = f.instantiate(params)
		out.Decorations[i] = f.Decoration
	}
	return out
}

type FunctionVariationGraphBlock VariationTableBlock

func (f FunctionVariationGraphBlock) instantiate(params expression.Variables, _ int) instance {
	return FunctionVariationGraphInstance(VariationTableBlock(f).instantiateVT(params))
}

type TableBlock struct {
	HorizontalHeaders []TextPart
	VerticalHeaders   []TextPart
	Values            [][]TextPart
}

func (t TableBlock) instantiate(params expression.Variables, _ int) instance {
	out := TableInstance{
		HorizontalHeaders: make([]client.TextOrMath, len(t.HorizontalHeaders)),
		VerticalHeaders:   make([]client.TextOrMath, len(t.VerticalHeaders)),
		Values:            make([][]client.TextOrMath, len(t.Values)),
	}
	for i, cell := range t.HorizontalHeaders {
		out.HorizontalHeaders[i] = cell.instantiate(params)
	}
	for i, cell := range t.VerticalHeaders {
		out.VerticalHeaders[i] = cell.instantiate(params)
	}

	for i, row := range t.Values {
		rowInstance := make([]client.TextOrMath, len(row))
		for j, v := range row {
			rowInstance[j] = v.instantiate(params)
		}
		out.Values[i] = rowInstance
	}
	return out
}
