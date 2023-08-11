package questions

import (
	"errors"
	"fmt"
	"strings"

	ex "github.com/benoitkugler/maths-online/server/src/maths/expression"
	"github.com/benoitkugler/maths-online/server/src/maths/functiongrapher"
	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	"github.com/benoitkugler/maths-online/server/src/maths/repere"
)

const exhaustiveTextKind = "exhaustiveTextKind"

var (
	_ Block = TextBlock{}
	_ Block = FormulaBlock{}
	_ Block = VariationTableBlock{}
	_ Block = SignTableBlock{}
	_ Block = FigureBlock{}
	_ Block = FunctionsGraphBlock{}
	_ Block = TableBlock{}
	_ Block = TreeBlock{}
)

// Block form the actual content of a question
// it is stored in a DB in generic form, but may be instantiated
// against random parameter values
type Block interface {
	// ID is only used by answer fields
	instantiate(params ex.Vars, ID int) (instance, error)

	// setupValidator is called on teacher input
	// it must :
	//	- performs validation not depending on instantiated parameters
	//  - delegates to `validator` for the other
	// is is also meant to ensure that only valid content is persisted on DB
	setupValidator(ex.RandomParameters) (validator, error)
}

// ParameterEntry is either a single variable definition,
// a special function or a (possibly multiline) comment.
type ParameterEntry interface {
	// Return a user friendly description
	String() string

	// mergeTo returns an `ErrDuplicateParameter` error if parameters are already defined
	mergeTo(vars ex.RandomParameters) error
}

func (rp Rp) String() string {
	return fmt.Sprintf("%s = %s", rp.Variable, rp.Expression)
}
func (it In) String() string { return string(it) }
func (cm Co) String() string { return string(cm) }

func (rp Rp) mergeTo(vars ex.RandomParameters) error {
	expr, err := ex.Parse(rp.Expression)
	if err != nil {
		return err
	}
	if _, has := vars[rp.Variable]; has {
		return ex.ErrDuplicateParameter{Duplicate: rp.Variable}
	}
	vars[rp.Variable] = expr
	return nil
}

func (it In) mergeTo(vars ex.RandomParameters) error {
	intr, err := ex.ParseIntrinsic(string(it))
	if err != nil {
		return err
	}
	return intr.MergeTo(vars)
}

// Comment are ignored
func (Co) mergeTo(vars ex.RandomParameters) error { return nil }

// ToMap may only be used after `Validate`
func (pr Parameters) ToMap() ex.RandomParameters {
	out := make(ex.RandomParameters)
	for _, entry := range pr {
		_ = entry.mergeTo(out) // error is check in Validate
	}
	return out
}

type Rp struct {
	Expression string      `json:"expression"` // as typed by the user, but validated
	Variable   ex.Variable `json:"variable"`
}

// String form on an intrinsic call, validated by expression.ParseIntrinsic
type In string

type Co string

// QuestionPage is the fundamental object to build exercices.
// It is mainly consituted of a list of content blocks, which
// describes the question (description, question, field answer),
// and are parametrized by random values.
type QuestionPage struct {
	Enonce     Enonce     `json:"enonce" gomacro-opaque:"dart"`
	Parameters Parameters `json:"parameters" gomacro-opaque:"dart"` // random parameters shared by the all the blocks
	Correction Enonce     `json:"correction" gomacro-opaque:"dart"`
}

// Instantiate returns a deep copy of `qu`, where all random parameters
// have been resolved.
// It assumes that the expressions and random parameters definitions are valid :
// if an error is encountered, it is returned as a TextInstance displaying the error.
func (qu QuestionPage) Instantiate() (out QuestionInstance, vars ex.Vars) {
	out, vars, err := qu.InstantiateErr()
	if err != nil {
		out = QuestionInstance{Enonce: EnonceInstance{
			TextInstance{
				Parts: []client.TextOrMath{
					{
						Text: fmt.Sprintf("Erreur inattendue : %v", err),
					},
				},
			},
		}}
	}
	return out, vars
}

func (qu QuestionPage) InstantiateWith(params ex.Vars) (QuestionInstance, error) {
	enonce, err := qu.Enonce.InstantiateWith(params)
	if err != nil {
		return QuestionInstance{}, err
	}
	correction, err := qu.Correction.InstantiateWith(params)
	if err != nil {
		return QuestionInstance{}, err
	}
	return QuestionInstance{enonce, correction}, err
}

// InstantiateErr is a shortcut to :
// - instantiate [Parameters]
// - instantiate [Enonce] with these parameters
// - instantiate [Correction] with these parameters
func (qu QuestionPage) InstantiateErr() (QuestionInstance, ex.Vars, error) {
	// generate random params
	rp, err := qu.Parameters.ToMap().Instantiate()
	if err != nil {
		return QuestionInstance{}, nil, err
	}
	instance, err := qu.InstantiateWith(rp)
	return instance, rp, err
}

// InstantiateWith uses the given values to instantiate the general question
func (qu Enonce) InstantiateWith(params ex.Vars) (EnonceInstance, error) {
	enonce := make(EnonceInstance, len(qu))
	var currentID int
	for j, bl := range qu {
		var err error
		enonce[j], err = bl.instantiate(params, currentID)
		if err != nil {
			return nil, err
		}
		if _, isField := enonce[j].(fieldInstance); isField {
			currentID++
		}
	}
	return enonce, nil
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

func (tp TextPart) instantiate(params ex.Vars) (client.TextOrMath, error) {
	switch tp.Kind {
	case Text:
		return client.TextOrMath{Text: tp.Content}, nil
	case StaticMath:
		return client.TextOrMath{Text: tp.Content, IsMath: true}, nil
	case Expression:
		expr, err := ex.ParseCompound(tp.Content)
		if err != nil {
			return client.TextOrMath{}, err
		}
		expr.Substitute(params)
		return client.TextOrMath{Text: expr.AsLaTeX(), IsMath: true}, nil
	default:
		panic(exhaustiveTextKind)
	}
}

func (tp TextPart) validate() error {
	switch tp.Kind {
	case Text, StaticMath:
		return nil // nothing to do
	case Expression:
		_, err := ex.ParseCompound(tp.Content)
		return err
	default:
		panic(exhaustiveTextKind)
	}
}

type TextParts []TextPart

// instantiate merges adjacent math chunks so that latex expression are not split up
// and may be successfully parsed by the client
func (tp TextParts) instantiate(params ex.Vars) (client.TextLine, error) {
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
func (tp TextParts) instantiateAndMerge(params ex.Vars) (string, error) {
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

// TextBlock is a chunk of text
// which may contain maths
// It support basic interpolation syntax.
type TextBlock struct {
	Parts   Interpolated
	Bold    bool
	Italic  bool
	Smaller bool
}

func (t TextBlock) instantiate(params ex.Vars, _ int) (instance, error) {
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

func (t TextBlock) setupValidator(ex.RandomParameters) (validator, error) {
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

// FormulaBlock is a math formula, which should be display using
// a LaTeX renderer.
type FormulaBlock struct {
	Parts Interpolated
}

func (f FormulaBlock) instantiate(params ex.Vars, _ int) (instance, error) {
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

func (f FormulaBlock) setupValidator(ex.RandomParameters) (validator, error) {
	_, err := f.Parts.parse()
	return noOpValidator{}, err
}

func evaluateExpr(expr string, params ex.Vars) (float64, error) {
	e, err := ex.Parse(expr)
	if err != nil {
		return 0, err
	}
	return e.Evaluate(params)
}

// parse, substitute and return LaTeX format
func instantiateLaTeXExpr(expr string, params ex.Vars) (string, error) {
	e, err := ex.Parse(expr)
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

func (vt VariationTableBlock) instantiateVT(params ex.Vars) (VariationTableInstance, error) {
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

func (vt VariationTableBlock) instantiate(params ex.Vars, _ int) (instance, error) {
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

	xExprs := make([]*ex.Expr, len(vt.Xs))
	fxExprs := make([]*ex.Expr, len(vt.Fxs))
	for i, c := range vt.Xs {
		var err error
		xExprs[i], err = ex.Parse(c)
		if err != nil {
			return variationTableValidator{}, err
		}
		fxExprs[i], err = ex.Parse(vt.Fxs[i])
		if err != nil {
			return variationTableValidator{}, err
		}
	}

	return variationTableValidator{label: label, xs: xExprs, fxs: fxExprs}, nil
}

func (vt VariationTableBlock) setupValidator(ex.RandomParameters) (validator, error) {
	out, err := vt.setupValidatorVT()
	if err != nil {
		return nil, err
	}
	return out, nil
}

type SignTableBlock struct {
	Xs        []string // valid expression
	Functions []client.FunctionSign
}

func (st SignTableBlock) instantiate(params ex.Vars, _ int) (instance, error) {
	return st.instantiateST(params)
}

func (st SignTableBlock) instantiateST(params ex.Vars) (SignTableInstance, error) {
	out := SignTableInstance{
		Xs:        make([]*ex.Expr, len(st.Xs)),
		Functions: append([]client.FunctionSign(nil), st.Functions...),
	}
	var err error
	for i, c := range st.Xs {
		out.Xs[i], err = ex.Parse(c)
		if err != nil {
			return out, err
		}
		out.Xs[i].Substitute(params)
	}
	return out, nil
}

func (st SignTableBlock) setupValidator(ex.RandomParameters) (validator, error) {
	if len(st.Xs) < 2 {
		return nil, errors.New("Au moins deux colonnes sont attendues.")
	}

	for _, function := range st.Functions {
		if len(st.Xs) != len(function.FxSymbols) || len(function.Signs) != len(st.Xs)-1 {
			return nil, errors.New("internal error: unexpected length for X and Fx")
		}
	}

	for _, c := range st.Xs {
		_, err := ex.Parse(c)
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

func (f FigureBlock) instantiate(params ex.Vars, _ int) (instance, error) {
	return f.instantiateF(params)
}

func (f FigureBlock) instantiateF(params ex.Vars) (FigureInstance, error) {
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

func (f FigureBlock) setupValidator(ex.RandomParameters) (validator, error) {
	var (
		out figureValidator
		err error
	)

	if f.Bounds.Height <= 0 || f.Bounds.Width <= 0 {
		return nil, errors.New("Les dimensions de la figure sont invalides.")
	}

	out.pointNames = make([]*ex.Expr, len(f.Drawings.Points))
	out.points = make([]*ex.Expr, 0, 2*len(f.Drawings.Points))
	for i, v := range f.Drawings.Points {
		out.pointNames[i], err = ex.Parse(v.Name)
		if err != nil {
			return nil, err
		}
		ptX, err := ex.Parse(v.Point.Coord.X)
		if err != nil {
			return nil, err
		}
		ptY, err := ex.Parse(v.Point.Coord.Y)
		if err != nil {
			return nil, err
		}
		out.points = append(out.points, ptX, ptY)
	}

	// ... and undefined points
	out.references = make([]*ex.Expr, 0, 2*len(f.Drawings.Segments))
	for _, seg := range f.Drawings.Segments {
		// validate the syntax for the name, which support interpolation
		_, err = Interpolated(seg.LabelName).parse()
		if err != nil {
			return nil, err
		}

		from, err := ex.Parse(seg.From)
		if err != nil {
			return nil, err
		}
		to, err := ex.Parse(seg.To)
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
			e, err := ex.Parse(point)
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

		radius, err := ex.Parse(circle.Radius)
		if err != nil {
			return nil, err
		}

		_, err = Interpolated(circle.Legend).parse()
		if err != nil {
			return nil, err
		}

		out.circlesDims = append(out.circlesDims, center.X, center.Y, radius)
	}

	out.lines = make([][2]*ex.Expr, len(f.Drawings.Lines))
	for i, l := range f.Drawings.Lines {
		out.lines[i][0], err = ex.Parse(l.A)
		if err != nil {
			return nil, err
		}
		out.lines[i][1], err = ex.Parse(l.B)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

type FunctionDefinition struct {
	Function   string // expression.Expression
	Decoration functiongrapher.FunctionDecoration
	Variable   ex.Variable // usually x
	From, To   string      // definition domain, expression.Expression
}

func (fg FunctionDefinition) parse() (fn ex.FunctionExpr, from, to *ex.Expr, err error) {
	expr, err := ex.Parse(fg.Function)
	if err != nil {
		return fn, from, to, err
	}
	from, err = ex.Parse(fg.From)
	if err != nil {
		return fn, from, to, err
	}
	to, err = ex.Parse(fg.To)
	if err != nil {
		return fn, from, to, err
	}

	return ex.FunctionExpr{
		Function: expr,
		Variable: fg.Variable,
	}, from, to, nil
}

func (fg FunctionDefinition) instantiate(params ex.Vars) (ex.FunctionDefinition, error) {
	fnExpr, from, to, err := fg.parse()
	if err != nil {
		return ex.FunctionDefinition{}, err
	}
	fnExpr.Function.Substitute(params)

	fromV, err := from.Evaluate(params)
	if err != nil {
		return ex.FunctionDefinition{}, err
	}
	toV, err := to.Evaluate(params)
	if err != nil {
		return ex.FunctionDefinition{}, err
	}

	return ex.FunctionDefinition{
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
	Color       repere.ColorHex
}

// FunctionPoint draws a point at a given
// abscice for a given function
type FunctionPoint struct {
	Function Interpolated // reference to a function [Label]
	X        string       // expression.Expression
	Color    repere.ColorHex
	Legend   Interpolated // legend
}

// FunctionsGraphBlock draws a figure with functions
// curves and colored areas
// Function are identifier by their [Label]
type FunctionsGraphBlock struct {
	FunctionExprs      []FunctionDefinition
	FunctionVariations []VariationTableBlock
	SequenceExprs      []FunctionDefinition // displayed as discrete sequences
	Areas              []FunctionArea
	Points             []FunctionPoint
}

func (fg FunctionsGraphBlock) setupValidator(params ex.RandomParameters) (validator, error) {
	out := functionsGraphValidator{
		functions:          make([]function, len(fg.FunctionExprs)),
		variationValidator: make([]variationTableValidator, len(fg.FunctionVariations)),
		sequences:          make([]function, len(fg.SequenceExprs)),
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
	for i, f := range fg.SequenceExprs {
		var err error
		out.sequences[i], err = newFunction(f, params)
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
		out.areas[i].domain.From, err = ex.Parse(area.Left)
		if err != nil {
			return nil, err
		}
		out.areas[i].domain.To, err = ex.Parse(area.Right)
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
		out.points[i].x, err = ex.Parse(point.X)
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

func (fg FunctionsGraphBlock) instantiate(params ex.Vars, _ int) (instance, error) {
	return fg.instantiateG(params)
}

func (fg FunctionsGraphBlock) instantiateG(params ex.Vars) (FunctionsGraphInstance, error) {
	out := FunctionsGraphInstance{}

	byNames := make(map[string][]domainCurves)

	// instantiate expression
	for _, f := range fg.FunctionExprs {
		fd, err := f.instantiate(params)
		if err != nil {
			return out, err
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
	for _, f := range fg.SequenceExprs {
		fd, err := f.instantiate(params)
		if err != nil {
			return out, err
		}
		fg := functiongrapher.SequenceGraph{
			Points:     functiongrapher.NewSequenceGraph(fd),
			Decoration: f.Decoration,
		}
		out.Sequences = append(out.Sequences, fg)
	}

	// instantiate variations
	for _, f := range fg.FunctionVariations {
		vt, err := f.instantiateVT(params)
		if err != nil {
			return out, err
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
			return out, err
		}
		bottomLabel, err := area.Bottom.instantiateAndMerge(params)
		if err != nil {
			return out, err
		}

		left, err := evaluateExpr(area.Left, params)
		if err != nil {
			return out, err
		}
		right, err := evaluateExpr(area.Right, params)
		if err != nil {
			return out, err
		}

		// select the curve containing [left, right] (guarded by the validation)
		topCandidates := byNames[topLabel]
		if topLabel == "" { // use the abscisse axis
			topCandidates = []domainCurves{horizontalAxis(left, right)}
		}
		top, err := selectByDomain(topCandidates, left, right)
		if err != nil {
			return out, err
		}
		bottomCandidates := byNames[bottomLabel]
		if bottomLabel == "" { // use the abscisse axis
			bottomCandidates = []domainCurves{horizontalAxis(left, right)}
		}
		bottom, err := selectByDomain(bottomCandidates, left, right)
		if err != nil {
			return out, err
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
			return out, err
		}

		x, err := evaluateExpr(point.X, params)
		if err != nil {
			return out, err
		}
		// select the curve containing x (guarded by the validation)
		fnLabel, err := point.Function.instantiateAndMerge(params)
		if err != nil {
			return out, err
		}
		candidates := byNames[fnLabel]
		if fnLabel == "" { // use the abscisse axis
			candidates = []domainCurves{horizontalAxis(x, x)}
		}

		domain, err := selectByDomain(candidates, x, x)
		if err != nil {
			return out, err
		}
		y, err := functiongrapher.OrdinateAt(domain.curves, x)
		if err != nil {
			return out, err
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

func (t TableBlock) instantiate(params ex.Vars, _ int) (instance, error) {
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

func (t TableBlock) setupValidator(ex.RandomParameters) (validator, error) {
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

// TreeNodeAnswer is an event, with (optional) children
type TreeNodeAnswer struct {
	Children      []TreeNodeAnswer `gomacro-data:"ignore"`
	Probabilities []string         // edges, same length as Children, valid expression.Expression
	Value         int              // index into the proposals, 0 for the root
}

type TreeBlock struct {
	EventsProposals []Interpolated
	AnswerRoot      TreeNodeAnswer
}

func (tf TreeBlock) instantiate(params ex.Vars, ID int) (instance, error) {
	return tf.instantiateT(params)
}

func (tf TreeBlock) instantiateT(params ex.Vars) (TreeInstance, error) {
	out := TreeInstance{
		EventsProposals: make([]client.TextLine, len(tf.EventsProposals)),
	}
	for i, p := range tf.EventsProposals {
		var err error
		out.EventsProposals[i], err = p.instantiate(params)
		if err != nil {
			return TreeInstance{}, err
		}
	}

	var buildTree func(node TreeNodeAnswer) (TreeNodeInstance, error)
	buildTree = func(node TreeNodeAnswer) (TreeNodeInstance, error) {
		out := TreeNodeInstance{
			Value:         node.Value,
			Probabilities: make([]*ex.Expr, len(node.Probabilities)),
			Children:      make([]TreeNodeInstance, len(node.Children)),
		}
		for i, c := range node.Probabilities {
			expr, err := ex.Parse(c)
			if err != nil {
				return out, err
			}
			expr.Substitute(params)
			out.Probabilities[i] = expr
		}
		for i, c := range node.Children {
			var err error
			out.Children[i], err = buildTree(c)
			if err != nil {
				return out, err
			}
		}
		return out, nil
	}

	root, err := buildTree(tf.AnswerRoot)
	out.AnswerRoot = root
	return out, err
}

func (tf TreeBlock) setupValidator(params ex.RandomParameters) (validator, error) {
	// check events syntax
	for _, event := range tf.EventsProposals {
		_, err := event.parse()
		if err != nil {
			return nil, err
		}
	}
	return treeValidator{data: tf}, nil
}
