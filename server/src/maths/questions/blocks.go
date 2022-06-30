package questions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/functiongrapher"
	"github.com/benoitkugler/maths-online/maths/questions/client"
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

// Block form the actual content of a question
// it is stored in a DB in generic form, but may be instantiated
// against random parameter values
type Block interface {
	// ID is only used by answer fields
	instantiate(params expression.Vars, ID int) (instance, error)

	// setupValidator is called on teacher input
	// it must :
	//	- performs validation not depending on instantiated parameters
	//  - delegates to `validator` for the other
	// is is also meant to ensure that only valid content is persisted on DB
	setupValidator(expression.RandomParameters) (validator, error)
}

type Parameters struct {
	Variables  RandomParameters
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

// RandomParameters is a serialized form of expression.RandomParameters
type RandomParameters []RandomParameter

type RandomParameter struct {
	Expression string              `json:"expression"` // as typed by the user, but validated
	Variable   expression.Variable `json:"variable"`
}

// toMap assumes `rp` only contains valid expressions,
// or it will panic
// It may only be used after `ValidateParameters`
func (rp RandomParameters) toMap() expression.RandomParameters {
	out := make(expression.RandomParameters, len(rp))
	for _, item := range rp {
		out[expression.Variable(item.Variable)] = expression.MustParse(item.Expression)
	}
	return out
}

// QuestionPage is the fundamental object to build exercices.
// It is mainly consituted of a list of content blocks, which
// describes the question (description, question, field answer),
// and are parametrized by random values.
type QuestionPage struct {
	Title      string     `json:"title"` // name of the question, optional
	Enonce     Enonce     `json:"enonce"`
	Parameters Parameters `json:"parameters"` // random parameters shared by the all the blocks
}

// Instantiate returns a deep copy of `qu`, where all random parameters
// have been resolved.
// It assumes that the expressions and random parameters definitions are valid :
// if an error is encountered, it is returned as a TextInstance displaying the error.
func (qu QuestionPage) Instantiate() (out QuestionInstance) {
	out, err := qu.instantiate()
	if err != nil {
		return QuestionInstance{Title: "Erreur", Enonce: EnonceInstance{
			TextInstance{
				Parts: []client.TextOrMath{
					{
						Text: fmt.Sprintf("Erreur inatendue : %v", err),
					},
				},
			},
		}}
	}
	return out
}

func (qu QuestionPage) instantiate() (QuestionInstance, error) {
	// generate random params
	rp, err := qu.Parameters.ToMap().Instantiate()
	if err != nil {
		return QuestionInstance{}, err
	}
	return qu.InstantiateWith(rp)
}

// InstantiateWith uses the given values to instantiate the general question
func (qu QuestionPage) InstantiateWith(params expression.Vars) (QuestionInstance, error) {
	enonce := make(EnonceInstance, len(qu.Enonce))
	var currentID int
	for j, bl := range qu.Enonce {
		var err error
		enonce[j], err = bl.instantiate(params, currentID)
		if err != nil {
			return QuestionInstance{}, err
		}
		if _, isField := enonce[j].(fieldInstance); isField {
			currentID++
		}
	}
	return QuestionInstance{Title: qu.Title, Enonce: enonce}, nil
}

// TextPart is either a plain text, a LaTeX code or an expression
type TextPart struct {
	Content string
	Kind    TextKind
}

func NewPText(content string) TextPart {
	return TextPart{Content: content, Kind: Text}
}

func NewPMath(content string) TextPart {
	return TextPart{Content: content, Kind: StaticMath}
}

func NewPExpr(content string) TextPart {
	return TextPart{Content: content, Kind: Expression}
}

func (tp TextPart) instantiate(params expression.Vars) (client.TextOrMath, error) {
	switch tp.Kind {
	case Text:
		return client.TextOrMath{Text: tp.Content}, nil
	case StaticMath:
		return client.TextOrMath{Text: tp.Content, IsMath: true}, nil
	case Expression:
		expr, err := expression.Parse(tp.Content)
		if err != nil {
			return client.TextOrMath{}, err
		}
		expr.Substitute(params)
		return client.TextOrMath{Text: expr.AsLaTeX(), IsMath: true}, nil
	default:
		panic(ExhaustiveTextKind)
	}
}

func (tp TextPart) validate() error {
	switch tp.Kind {
	case Text, StaticMath:
		return nil // nothing to do
	case Expression:
		_, err := expression.Parse(tp.Content)
		return err
	default:
		panic(ExhaustiveTextKind)
	}
}

type TextParts []TextPart

// instantiate merges adjacent math chunks so that latex expression are not split up
// and may be successfully parsed by the client
func (tp TextParts) instantiate(params expression.Vars) (client.TextLine, error) {
	var parts client.TextLine
	for _, p := range tp {
		sample, err := p.instantiate(params)
		if err != nil {
			return nil, err
		}
		L := len(parts)
		if L == 0 {
			parts = append(parts, sample)
			continue
		}

		// check if the previous chunk as same type
		if parts[L-1].IsMath == sample.IsMath {
			// simply merge the contents
			parts[L-1].Text = parts[L-1].Text + sample.Text
		} else { // start a new chunk
			parts = append(parts, sample)
		}
	}
	return parts, nil
}

// assume all parts are either static math or expression.
func (tp TextParts) instantiateAndMerge(params expression.Vars) (string, error) {
	parts, err := tp.instantiate(params)
	if err != nil {
		return "", err
	}
	chunks := make([]string, len(parts))
	for i, p := range parts {
		chunks[i] = p.Text
	}
	return strings.Join(chunks, ""), nil
}

func (tp TextParts) validate() error {
	for _, text := range tp {
		if err := text.validate(); err != nil {
			return err
		}
	}
	return nil
}

// TextBlock is a chunk of text
// which may contain maths
// It support basic interpolation syntax.
type TextBlock struct {
	Parts   Interpolated
	Bold    bool
	Italic  bool
	Smaller bool
}

func (t TextBlock) instantiate(params expression.Vars, _ int) (instance, error) {
	parts, err := t.Parts.instantiate(params)
	if err != nil {
		return nil, err
	}
	return TextInstance{
		Parts:   parts,
		Bold:    t.Bold,
		Italic:  t.Italic,
		Smaller: t.Smaller,
	}, nil
}

func (t TextBlock) setupValidator(expression.RandomParameters) (validator, error) {
	_, err := t.Parts.Parse()
	return noOpValidator{}, err
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
func (fp FormulaPart) instantiate(params expression.Vars) (StringOrExpression, error) {
	if !fp.IsExpression { // nothing to do
		return StringOrExpression{String: fp.Content}, nil
	}

	expr, err := expression.Parse(fp.Content)
	if err != nil {
		return StringOrExpression{}, err
	}
	expr.Substitute(params)
	return StringOrExpression{Expression: expr}, nil
}

// FormulaBlock is a math formula, which should be display using
// a LaTeX renderer.
type FormulaBlock struct {
	Parts Interpolated
}

func (f FormulaBlock) instantiate(params expression.Vars, _ int) (instance, error) {
	parts, err := f.Parts.Parse()
	if err != nil {
		return nil, err
	}
	partsInstance, err := parts.instantiate(params)
	if err != nil {
		return nil, err
	}
	out := make(FormulaDisplayInstance, len(partsInstance))
	for i, c := range partsInstance {
		out[i] = c.Text
	}
	return out, nil
}

func (f FormulaBlock) setupValidator(expression.RandomParameters) (validator, error) {
	_, err := f.Parts.Parse()
	return noOpValidator{}, err
}

func evaluateExpr(expr string, params expression.Vars) (float64, error) {
	e, err := expression.Parse(expr)
	if err != nil {
		return 0, err
	}
	return e.Evaluate(params)
}

type VariationTableBlock struct {
	Label Interpolated
	Xs    []string // expressions
	Fxs   []string // expressions
}

func (vt VariationTableBlock) instantiateVT(params expression.Vars) (VariationTableInstance, error) {
	out := VariationTableInstance{
		Xs:  make([]evaluatedExpression, len(vt.Xs)),
		Fxs: make([]evaluatedExpression, len(vt.Fxs)),
	}

	parts, err := vt.Label.Parse()
	if err != nil {
		return out, err
	}
	out.Label, err = parts.instantiateAndMerge(params)
	if err != nil {
		return out, err
	}

	for i, c := range vt.Xs {
		out.Xs[i], err = newEvaluatedExpression(c, params)
		if err != nil {
			return out, err
		}
	}
	for i, c := range vt.Fxs {
		out.Fxs[i], err = newEvaluatedExpression(c, params)
		if err != nil {
			return out, err
		}
	}

	return out, nil
}

func (vt VariationTableBlock) instantiate(params expression.Vars, _ int) (instance, error) {
	return vt.instantiateVT(params)
}

func (vt VariationTableBlock) setupValidator(expression.RandomParameters) (validator, error) {
	_, err := vt.Label.Parse()
	if err != nil {
		return nil, err
	}

	if len(vt.Xs) < 2 {
		return nil, errors.New("Au moins deux colonnes sont attendues.")
	}

	if len(vt.Xs) != len(vt.Fxs) {
		return nil, errors.New("internal error: expected same length for X and Fx")
	}

	xExprs := make([]*expression.Expr, len(vt.Xs))
	fxExprs := make([]*expression.Expr, len(vt.Fxs))
	for i, c := range vt.Xs {
		var err error
		xExprs[i], err = expression.Parse(c)
		if err != nil {
			return nil, err
		}
		fxExprs[i], err = expression.Parse(vt.Fxs[i])
		if err != nil {
			return nil, err
		}
	}

	return variationTableValidator{xs: xExprs, fxs: fxExprs}, nil
}

type SignTableBlock struct {
	Label     string
	FxSymbols []SignSymbol
	Xs        []Interpolated // always math content
	Signs     []bool         // with length len(Xs) - 1
}

func (st SignTableBlock) instantiate(params expression.Vars, _ int) (instance, error) {
	out := SignTableInstance{
		Label: st.Label,
		Xs:    make([]string, len(st.Xs)),
	}
	for i, c := range st.Xs {
		parts, err := c.Parse()
		if err != nil {
			return nil, err
		}
		out.Xs[i], err = parts.instantiateAndMerge(params)
		if err != nil {
			return nil, err
		}
	}
	out.FxSymbols = append([]SignSymbol(nil), st.FxSymbols...)
	out.Signs = append([]bool(nil), st.Signs...)
	return out, nil
}

func (st SignTableBlock) setupValidator(expression.RandomParameters) (validator, error) {
	if len(st.Xs) < 2 {
		return nil, errors.New("Au moins deux colonnes sont attendues.")
	}

	if len(st.Xs) != len(st.FxSymbols) || len(st.Signs) != len(st.Xs)-1 {
		return nil, errors.New("internal error: unexpected length for X and Fx")
	}

	for _, c := range st.Xs {
		_, err := c.Parse()
		if err != nil {
			return nil, err
		}
	}

	return noOpValidator{}, nil
}

type FigureBlock struct {
	Drawings repere.RandomDrawings
	Bounds   repere.RepereBounds
	ShowGrid bool
}

func (f FigureBlock) instantiate(params expression.Vars, _ int) (instance, error) {
	return f.instantiateF(params)
}

func (f FigureBlock) instantiateF(params expression.Vars) (FigureInstance, error) {
	out := FigureInstance{
		Figure: repere.Figure{
			Drawings: repere.Drawings{
				Segments: make([]repere.Segment, len(f.Drawings.Segments)),
				Points:   make(map[string]repere.LabeledPoint),
				Lines:    make([]repere.Line, len(f.Drawings.Lines)),
			},
			Bounds:   f.Bounds,
			ShowGrid: f.ShowGrid,
		},
	}
	for _, v := range f.Drawings.Points {
		nameExpr, err := expression.Parse(v.Name)
		if err != nil {
			return out, err
		}
		nameExpr.Substitute(params)
		name := nameExpr.AsLaTeX()

		x, err := evaluateExpr(v.Point.Coord.X, params)
		if err != nil {
			return out, err
		}
		y, err := evaluateExpr(v.Point.Coord.Y, params)
		if err != nil {
			return out, err
		}
		out.Figure.Drawings.Points[name] = repere.LabeledPoint{
			Point: repere.PosPoint{
				Point: repere.Coord{
					X: x,
					Y: y,
				},
				Pos: v.Point.Pos,
			},
			Color: v.Point.Color,
		}
	}

	for i, s := range f.Drawings.Segments {
		fromExpr, err := expression.Parse(s.From)
		if err != nil {
			return out, err
		}
		fromExpr.Substitute(params)
		s.From = fromExpr.AsLaTeX()

		toExpr, err := expression.Parse(s.To)
		if err != nil {
			return out, err
		}
		toExpr.Substitute(params)
		s.To = toExpr.AsLaTeX()

		out.Figure.Drawings.Segments[i] = s
	}

	for i, l := range f.Drawings.Lines {
		a, err := evaluateExpr(l.A, params)
		if err != nil {
			return out, err
		}
		b, err := evaluateExpr(l.B, params)
		if err != nil {
			return out, err
		}
		out.Figure.Drawings.Lines[i] = repere.Line{
			Label: l.Label,
			A:     a,
			B:     b,
			Color: l.Color,
		}
	}
	return out, nil
}

func (f FigureBlock) setupValidator(expression.RandomParameters) (validator, error) {
	var (
		out figureValidator
		err error
	)

	if f.Bounds.Height <= 0 || f.Bounds.Width <= 0 {
		return nil, errors.New("Les dimensions de la figure sont invalides.")
	}

	out.pointNames = make([]*expression.Expr, len(f.Drawings.Points))
	out.points = make([]*expression.Expr, 0, 2*len(f.Drawings.Points))
	for i, v := range f.Drawings.Points {
		out.pointNames[i], err = expression.Parse(v.Name)
		if err != nil {
			return nil, err
		}
		ptX, err := expression.Parse(v.Point.Coord.X)
		if err != nil {
			return nil, err
		}
		ptY, err := expression.Parse(v.Point.Coord.Y)
		if err != nil {
			return nil, err
		}
		out.points = append(out.points, ptX, ptY)
	}

	// ... and undefined points
	out.references = make([]*expression.Expr, 0, 2*len(f.Drawings.Segments))
	for _, seg := range f.Drawings.Segments {
		from, err := expression.Parse(seg.From)
		if err != nil {
			return nil, err
		}
		to, err := expression.Parse(seg.To)
		if err != nil {
			return nil, err
		}
		out.references = append(out.references, from, to)
	}

	out.lines = make([][2]*expression.Expr, len(f.Drawings.Lines))
	for i, l := range f.Drawings.Lines {
		out.lines[i][0], err = expression.Parse(l.A)
		if err != nil {
			return nil, err
		}
		out.lines[i][1], err = expression.Parse(l.B)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

type FunctionDefinition struct {
	Function   string // expression.Expression
	Decoration functiongrapher.FunctionDecoration
	Variable   expression.Variable // usually x
	From, To   string              // definition domain, expression.Expression
}

func (fg FunctionDefinition) parse() (fn expression.FunctionExpr, from, to *expression.Expr, err error) {
	expr, err := expression.Parse(fg.Function)
	if err != nil {
		return fn, from, to, err
	}
	from, err = expression.Parse(fg.From)
	if err != nil {
		return fn, from, to, err
	}
	to, err = expression.Parse(fg.To)
	if err != nil {
		return fn, from, to, err
	}

	return expression.FunctionExpr{
		Function: expr,
		Variable: fg.Variable,
	}, from, to, nil
}

func (fg FunctionDefinition) instantiate(params expression.Vars) (expression.FunctionDefinition, error) {
	fnExpr, from, to, err := fg.parse()
	if err != nil {
		return expression.FunctionDefinition{}, err
	}
	fnExpr.Function.Substitute(params)

	fromV, err := from.Evaluate(params)
	if err != nil {
		return expression.FunctionDefinition{}, err
	}
	toV, err := to.Evaluate(params)
	if err != nil {
		return expression.FunctionDefinition{}, err
	}

	return expression.FunctionDefinition{
		FunctionExpr: fnExpr,
		From:         fromV,
		To:           toV,
	}, nil
}

type FunctionGraphBlock struct {
	Functions []FunctionDefinition
}

func (fg FunctionGraphBlock) instantiate(params expression.Vars, _ int) (instance, error) {
	out := FunctionGraphInstance{
		Functions:   make([]expression.FunctionDefinition, len(fg.Functions)),
		Decorations: make([]functiongrapher.FunctionDecoration, len(fg.Functions)),
	}
	for i, f := range fg.Functions {
		var err error
		out.Functions[i], err = f.instantiate(params)
		if err != nil {
			return nil, err
		}
		out.Decorations[i] = f.Decoration
	}
	return out, nil
}

func (fg FunctionGraphBlock) setupValidator(params expression.RandomParameters) (validator, error) {
	out := functionGraphValidator{functions: make([]function, len(fg.Functions))}
	for i, f := range fg.Functions {
		var err error
		out.functions[i], err = newFunction(f, params)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

type FunctionVariationGraphBlock VariationTableBlock

func (f FunctionVariationGraphBlock) instantiate(params expression.Vars, _ int) (instance, error) {
	out, err := VariationTableBlock(f).instantiateVT(params)
	return FunctionVariationGraphInstance(out), err
}

func (f FunctionVariationGraphBlock) setupValidator(params expression.RandomParameters) (validator, error) {
	return VariationTableBlock(f).setupValidator(params)
}

type TableBlock struct {
	HorizontalHeaders []TextPart
	VerticalHeaders   []TextPart
	Values            [][]TextPart
}

func (t TableBlock) instantiate(params expression.Vars, _ int) (instance, error) {
	out := TableInstance{
		HorizontalHeaders: make([]client.TextOrMath, len(t.HorizontalHeaders)),
		VerticalHeaders:   make([]client.TextOrMath, len(t.VerticalHeaders)),
		Values:            make([][]client.TextOrMath, len(t.Values)),
	}
	var err error
	for i, cell := range t.HorizontalHeaders {
		out.HorizontalHeaders[i], err = cell.instantiate(params)
		if err != nil {
			return nil, err
		}
	}
	for i, cell := range t.VerticalHeaders {
		out.VerticalHeaders[i], err = cell.instantiate(params)
		if err != nil {
			return nil, err
		}
	}

	for i, row := range t.Values {
		rowInstance := make([]client.TextOrMath, len(row))
		for j, cell := range row {
			rowInstance[j], err = cell.instantiate(params)
			if err != nil {
				return nil, err
			}
		}
		out.Values[i] = rowInstance
	}
	return out, nil
}

func (t TableBlock) setupValidator(expression.RandomParameters) (validator, error) {
	for _, cell := range t.HorizontalHeaders {
		if err := cell.validate(); err != nil {
			return nil, err
		}
	}
	for _, cell := range t.VerticalHeaders {
		if err := cell.validate(); err != nil {
			return nil, err
		}
	}
	for _, row := range t.Values {
		for _, cell := range row {
			if err := cell.validate(); err != nil {
				return nil, err
			}
		}
	}

	return noOpValidator{}, nil
}
