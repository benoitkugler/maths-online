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
	_ Block = FunctionsGraphBlock{}
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

// Append appends `other` slices to `pr` slices
func (pr Parameters) Append(other Parameters) (out Parameters) {
	out.Intrinsics = append(pr.Intrinsics, other.Intrinsics...)
	out.Variables = append(pr.Variables, other.Variables...)
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
		return QuestionInstance{Enonce: EnonceInstance{
			TextInstance{
				Parts: []client.TextOrMath{
					{
						Text: fmt.Sprintf("Erreur inattendue : %v", err),
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
	return QuestionInstance{Enonce: enonce}, nil
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
	_, err := t.Parts.parse()
	return noOpValidator{}, err
}

// FormulaContent is a list of chunks, either
//   - static math symbols, such as f(x) =
//   - valid expression, such as a*x - b, which will be instantiated
//
// when rendering the question
//
// For instance, the formula "f(x) = a*(x + 2)"
// is represented by two FormulaPart elements:
//
//	{ f(x) = } and { a*(x + 2) }
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
	parts, err := f.Parts.parse()
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
	_, err := f.Parts.parse()
	return noOpValidator{}, err
}

func evaluateExpr(expr string, params expression.Vars) (float64, error) {
	e, err := expression.Parse(expr)
	if err != nil {
		return 0, err
	}
	return e.Evaluate(params)
}

// parse, substitute and return LaTeX format
func instantiateLaTeXExpr(expr string, params expression.Vars) (string, error) {
	e, err := expression.Parse(expr)
	if err != nil {
		return "", err
	}
	e.Substitute(params)
	return e.AsLaTeX(), nil
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

	var err error
	out.Label, err = vt.Label.instantiateAndMerge(params)
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

func (vt VariationTableBlock) setupValidatorVT() (variationTableValidator, error) {
	label, err := vt.Label.parse()
	if err != nil {
		return variationTableValidator{}, err
	}

	if len(vt.Xs) < 2 {
		return variationTableValidator{}, errors.New("Au moins deux colonnes sont attendues.")
	}

	if len(vt.Xs) != len(vt.Fxs) {
		return variationTableValidator{}, errors.New("internal error: expected same length for X and Fx")
	}

	xExprs := make([]*expression.Expr, len(vt.Xs))
	fxExprs := make([]*expression.Expr, len(vt.Fxs))
	for i, c := range vt.Xs {
		var err error
		xExprs[i], err = expression.Parse(c)
		if err != nil {
			return variationTableValidator{}, err
		}
		fxExprs[i], err = expression.Parse(vt.Fxs[i])
		if err != nil {
			return variationTableValidator{}, err
		}
	}

	return variationTableValidator{label: label, xs: xExprs, fxs: fxExprs}, nil
}

func (vt VariationTableBlock) setupValidator(expression.RandomParameters) (validator, error) {
	out, err := vt.setupValidatorVT()
	if err != nil {
		return nil, err
	}
	return out, nil
}

type SignTableBlock struct {
	Label     string
	FxSymbols []SignSymbol
	Xs        []string // valid expression
	Signs     []bool   // is positive, with length len(Xs) - 1
}

func (st SignTableBlock) instantiate(params expression.Vars, _ int) (instance, error) {
	return st.instantiateST(params)
}

func (st SignTableBlock) instantiateST(params expression.Vars) (SignTableInstance, error) {
	out := SignTableInstance{
		Label: st.Label,
		Xs:    make([]*expression.Expr, len(st.Xs)),
	}
	var err error
	for i, c := range st.Xs {
		out.Xs[i], err = expression.Parse(c)
		if err != nil {
			return out, err
		}
		out.Xs[i].Substitute(params)
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
		_, err := expression.Parse(c)
		if err != nil {
			return nil, err
		}
	}

	return noOpValidator{}, nil
}

type FigureBlock struct {
	Drawings   repere.RandomDrawings
	Bounds     repere.RepereBounds
	ShowGrid   bool
	ShowOrigin bool
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
				Circles:  make([]repere.Circle, len(f.Drawings.Circles)),
				Areas:    make([]repere.Area, len(f.Drawings.Areas)),
			},
			Bounds:     f.Bounds,
			ShowGrid:   f.ShowGrid,
			ShowOrigin: f.ShowOrigin,
		},
	}
	for _, v := range f.Drawings.Points {
		name, err := instantiateLaTeXExpr(v.Name, params)
		if err != nil {
			return out, err
		}

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

	var err error
	for i, s := range f.Drawings.Segments {
		instance := repere.Segment(s)

		instance.LabelName, err = Interpolated(s.LabelName).instantiateAndMerge(params)
		if err != nil {
			return out, err
		}

		instance.From, err = instantiateLaTeXExpr(s.From, params)
		if err != nil {
			return out, err
		}
		instance.To, err = instantiateLaTeXExpr(s.To, params)
		if err != nil {
			return out, err
		}

		out.Figure.Drawings.Segments[i] = instance
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

	for i, circle := range f.Drawings.Circles {
		legend, err := Interpolated(circle.Legend).instantiateAndMerge(params)
		if err != nil {
			return out, err
		}

		center, err := CoordExpression(circle.Center).instantiateToFloat(params)
		if err != nil {
			return out, err
		}

		radius, err := evaluateExpr(circle.Radius, params)
		if err != nil {
			return out, err
		}

		out.Figure.Drawings.Circles[i] = repere.Circle{
			Radius:    radius,
			Center:    center,
			LineColor: circle.LineColor,
			FillColor: circle.FillColor,
			Legend:    legend,
		}
	}

	for i, area := range f.Drawings.Areas {
		instance := repere.Area{
			Color:  area.Color,
			Points: make([]repere.PointName, len(area.Points)),
		}
		for j, p := range area.Points {
			instance.Points[j], err = instantiateLaTeXExpr(p, params)
			if err != nil {
				return out, err
			}
		}
		out.Figure.Drawings.Areas[i] = instance
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
		// validate the syntax for the name, which support interpolation
		_, err = Interpolated(seg.LabelName).parse()
		if err != nil {
			return nil, err
		}

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

	for _, area := range f.Drawings.Areas {
		if len(area.Points) < 3 {
			return nil, errors.New("Une surface requiert au moins 3 points.")
		}

		for _, point := range area.Points {
			e, err := expression.Parse(point)
			if err != nil {
				return nil, err
			}
			out.references = append(out.references, e)
		}
	}

	for _, circle := range f.Drawings.Circles {
		center, err := CoordExpression(circle.Center).parse()
		if err != nil {
			return nil, err
		}

		radius, err := expression.Parse(circle.Radius)
		if err != nil {
			return nil, err
		}

		_, err = Interpolated(circle.Legend).parse()
		if err != nil {
			return nil, err
		}

		out.circlesDims = append(out.circlesDims, center.X, center.Y, radius)
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

type FunctionArea struct {
	// reference to function [Label]s, with empty meaning
	// horizontal line
	Bottom, Top Interpolated
	Left, Right string // expression.Expression
	Color       repere.Color
}

// FunctionPoint draws a point at a given
// abscice for a given function
type FunctionPoint struct {
	Function Interpolated // reference to a function [Label]
	X        string       // expression.Expression
	Color    repere.Color
	Legend   Interpolated // legend
}

// FunctionsGraphBlock draws a figure with functions
// curves and colored areas
// Function are identifier by their [Label]
type FunctionsGraphBlock struct {
	FunctionExprs      []FunctionDefinition
	FunctionVariations []VariationTableBlock
	Areas              []FunctionArea
	Points             []FunctionPoint
}

func (fg FunctionsGraphBlock) setupValidator(params expression.RandomParameters) (validator, error) {
	out := functionsGraphValidator{
		functions:          make([]function, len(fg.FunctionExprs)),
		variationValidator: make([]variationTableValidator, len(fg.FunctionVariations)),
		areas:              make([]areaVData, len(fg.Areas)),
		points:             make([]functionPointVData, len(fg.Points)),
	}
	for i, f := range fg.FunctionExprs {
		var err error
		out.functions[i], err = newFunction(f, params)
		if err != nil {
			return nil, err
		}
	}
	for i, vt := range fg.FunctionVariations {
		var err error
		out.variationValidator[i], err = vt.setupValidatorVT()
		if err != nil {
			return nil, err
		}
	}
	for i, area := range fg.Areas {
		var err error
		out.areas[i].top, err = area.Top.parse()
		if err != nil {
			return nil, err
		}
		out.areas[i].bottom, err = area.Bottom.parse()
		if err != nil {
			return nil, err
		}
		out.areas[i].domain.From, err = expression.Parse(area.Left)
		if err != nil {
			return nil, err
		}
		out.areas[i].domain.To, err = expression.Parse(area.Right)
		if err != nil {
			return nil, err
		}
	}
	for i, point := range fg.Points {
		_, err := point.Legend.parse() // check the syntax for the legend
		if err != nil {
			return nil, err
		}

		out.points[i].fnLabel, err = point.Function.parse()
		if err != nil {
			return nil, err
		}
		out.points[i].x, err = expression.Parse(point.X)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func extractValues(vt VariationTableInstance) (xs, fxs []float64) {
	xs = make([]float64, len(vt.Xs))
	fxs = make([]float64, len(vt.Fxs))
	for i, v := range vt.Xs {
		xs[i] = v.Value
	}
	for i, v := range vt.Fxs {
		fxs[i] = v.Value
	}
	return
}

type domainCurves struct {
	curves []functiongrapher.BezierCurve
	domain [2]float64
}

func horizontalAxis(left, right float64) domainCurves {
	return domainCurves{curves: functiongrapher.HorizontalAxis(left, right, 0), domain: [2]float64{left, right}}
}

func selectByDomain(candidates []domainCurves, left, right float64) (domainCurves, error) {
	for _, c := range candidates {
		if c.domain[0] <= left && right <= c.domain[1] {
			return c, nil
		}
	}
	return domainCurves{}, fmt.Errorf("aucun domaine ne contient [%f, %f]", left, right)
}

func (fg FunctionsGraphBlock) instantiate(params expression.Vars, _ int) (instance, error) {
	out := FunctionsGraphInstance{}

	byNames := make(map[string][]domainCurves)

	// instantiate expression
	for _, f := range fg.FunctionExprs {
		fd, err := f.instantiate(params)
		if err != nil {
			return nil, err
		}
		fg := functiongrapher.FunctionGraph{
			Segments:   functiongrapher.NewFunctionGraph(fd),
			Decoration: f.Decoration,
		}
		out.Functions = append(out.Functions, fg)
		byNames[fg.Decoration.Label] = append(byNames[fg.Decoration.Label], domainCurves{
			curves: fg.Segments,
			domain: [2]float64{fd.From, fd.To},
		})
	}

	// instantiate variations
	for _, f := range fg.FunctionVariations {
		vt, err := f.instantiateVT(params)
		if err != nil {
			return nil, err
		}
		xs, fxs := extractValues(vt)
		fg := functiongrapher.FunctionGraph{
			Segments:   functiongrapher.NewFunctionGraphFromVariations(xs, fxs),
			Decoration: functiongrapher.FunctionDecoration{Label: vt.Label},
		}
		out.Functions = append(out.Functions, fg)
		byNames[fg.Decoration.Label] = append(byNames[fg.Decoration.Label], domainCurves{
			curves: fg.Segments,
			domain: [2]float64{vt.Xs[0].Value, vt.Xs[len(vt.Xs)-1].Value},
		})
	}

	// instantiate areas
	for _, area := range fg.Areas {
		topLabel, err := area.Top.instantiateAndMerge(params)
		if err != nil {
			return nil, err
		}
		bottomLabel, err := area.Bottom.instantiateAndMerge(params)
		if err != nil {
			return nil, err
		}

		left, err := evaluateExpr(area.Left, params)
		if err != nil {
			return nil, err
		}
		right, err := evaluateExpr(area.Right, params)
		if err != nil {
			return nil, err
		}

		// select the curve containing [left, right] (guarded by the validation)
		topCandidates := byNames[topLabel]
		if topLabel == "" { // use the abscisse axis
			topCandidates = []domainCurves{horizontalAxis(left, right)}
		}
		top, err := selectByDomain(topCandidates, left, right)
		if err != nil {
			return nil, err
		}
		bottomCandidates := byNames[bottomLabel]
		if bottomLabel == "" { // use the abscisse axis
			bottomCandidates = []domainCurves{horizontalAxis(left, right)}
		}
		bottom, err := selectByDomain(bottomCandidates, left, right)
		if err != nil {
			return nil, err
		}

		path := functiongrapher.NewAreaBetween(top.curves, bottom.curves, left, right)
		out.Areas = append(out.Areas, client.FunctionArea{
			Path:  path,
			Color: area.Color,
		})
	}

	// instantiate points
	for _, point := range fg.Points {
		legend, err := point.Legend.instantiateAndMerge(params)
		if err != nil {
			return nil, err
		}

		x, err := evaluateExpr(point.X, params)
		if err != nil {
			return nil, err
		}
		// select the curve containing x (guarded by the validation)
		fnLabel, err := point.Function.instantiateAndMerge(params)
		if err != nil {
			return nil, err
		}
		candidates := byNames[fnLabel]
		if fnLabel == "" { // use the abscisse axis
			candidates = []domainCurves{horizontalAxis(x, x)}
		}

		domain, err := selectByDomain(candidates, x, x)
		if err != nil {
			return nil, err
		}
		y, err := functiongrapher.OrdinateAt(domain.curves, x)
		if err != nil {
			return nil, err
		}
		out.Points = append(out.Points, client.FunctionPoint{
			Color:  point.Color,
			Legend: legend,
			Coord: repere.Coord{
				X: x,
				Y: y,
			},
		})
	}

	return out, nil
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
