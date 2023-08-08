package questions

import (
	"errors"
	"math/rand"
	"sort"

	"github.com/benoitkugler/maths-online/server/src/maths/expression"
	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	"github.com/benoitkugler/maths-online/server/src/maths/repere"
)

var (
	_ Block = NumberFieldBlock{}
	_ Block = ExpressionFieldBlock{}
	_ Block = RadioFieldBlock{}
	_ Block = OrderedListFieldBlock{}
	_ Block = FigurePointFieldBlock{}
	_ Block = FigureVectorFieldBlock{}
	_ Block = VariationTableFieldBlock{}
	_ Block = SignTableFieldBlock{}
	_ Block = FunctionPointsFieldBlock{}
	_ Block = FigureVectorPairFieldBlock{}
	_ Block = FigureAffineLineFieldBlock{}
	_ Block = TreeFieldBlock{}
	_ Block = TableFieldBlock{}
	_ Block = VectorFieldBlock{}
	_ Block = ProofFieldBlock{}
)

type NumberFieldBlock struct {
	// a valid expression, in the format used by expression.Expression
	// which is only parametrized by the random parameters
	Expression string
}

func (n NumberFieldBlock) instantiate(params expression.Vars, ID int) (instance, error) {
	answer, err := evaluateExpr(n.Expression, params)
	return NumberFieldInstance{ID: ID, Answer: answer}, err
}

func (n NumberFieldBlock) setupValidator(params expression.RandomParameters) (validator, error) {
	expr, err := expression.Parse(n.Expression)
	if err != nil {
		return nil, err
	}
	return numberValidator{expr: expr}, nil
}

type ExpressionFieldBlock struct {
	// A valid expression, in the format used by expression.Expression or expression.Compound
	Expression       string
	Label            Interpolated // optional
	ComparisonLevel  ComparisonLevel
	ShowFractionHelp bool // if true an hint for fraction is displayed when applicable
}

func (f ExpressionFieldBlock) SyntaxHint(params Parameters) (TextBlock, error) {
	answer, err := expression.ParseCompound(f.Expression)
	if err != nil {
		return TextBlock{}, err
	}
	allExprs := answer.Expressions()

	if err := params.Validate(); err != nil {
		return TextBlock{}, err
	}
	m := params.ToMap()

	const nbRepeat = 100
	allHints := make(expression.SyntaxHints)
	for i := 0; i < nbRepeat; i++ {
		params, err := m.Instantiate()
		if err != nil {
			return TextBlock{}, err
		}
		for _, expr := range allExprs {
			expr = expr.Copy()
			expr.Substitute(params)
			allHints.Append(expr.SyntaxHints())
		}
	}
	out := TextBlock{
		Italic:  true,
		Smaller: true,
		Parts:   Interpolated(allHints.Text()),
	}
	return out, nil
}

func (f ExpressionFieldBlock) instantiate(params expression.Vars, ID int) (instance, error) {
	answer, err := expression.ParseCompound(f.Expression)
	if err != nil {
		return nil, err
	}
	answer.Substitute(params)

	var showFractionHelp bool
	answerExpr, ok := answer.(*expression.Expr)
	if ok {
		answerExpr.DefaultSimplify() // needed for better [IsFraction] result
		showFractionHelp = f.ShowFractionHelp && answerExpr.IsFraction()
	}

	label, err := f.Label.instantiateAndMerge(params)
	if err != nil {
		return nil, err
	}

	return ExpressionFieldInstance{
		LabelLaTeX:       label,
		Answer:           answer,
		ComparisonLevel:  f.ComparisonLevel,
		ShowFractionHelp: showFractionHelp,
		ID:               ID,
	}, nil
}

func (f ExpressionFieldBlock) setupValidator(expression.RandomParameters) (validator, error) {
	expr, err := expression.ParseCompound(f.Expression)
	if err != nil {
		return nil, err
	}
	_, err = f.Label.parse()
	if err != nil {
		return nil, err
	}

	asExpr, isExpr := expr.(*expression.Expr)

	if f.ShowFractionHelp && !isExpr {
		return nil, errors.New("L'aide aux fractions n'est utilisable que pour une expression simple.")
	}

	switch f.ComparisonLevel {
	case AsLinearEquation:
		if !isExpr {
			return nil, errors.New("Une expression simple est attendue pour une équation cartésienne.")
		}
		return linearEquationValidator{expr: asExpr}, nil
	default:
		return noOpValidator{}, nil
	}
}

type RadioFieldBlock struct {
	Answer     string         // must satisfy expression.IsValidIndex
	Proposals  []Interpolated // slice of text parts
	AsDropDown bool
}

func (rf RadioFieldBlock) instantiate(params expression.Vars, ID int) (instance, error) {
	ans, err := evaluateExpr(rf.Answer, params)
	if err != nil {
		return nil, err
	}
	out := RadioFieldInstance{
		Proposals: make([]client.TextLine, len(rf.Proposals)),
		Answer:    int(ans),
		ID:        ID,
	}
	for i, p := range rf.Proposals {
		props, err := p.instantiate(params)
		if err != nil {
			return nil, err
		}
		out.Proposals[i] = props
	}

	if rf.AsDropDown {
		return DropDownFieldInstance(out), nil
	}
	return out, nil
}

func (rf RadioFieldBlock) setupValidator(params expression.RandomParameters) (validator, error) {
	for _, p := range rf.Proposals {
		_, err := p.parse()
		if err != nil {
			return nil, err
		}
	}

	expr, err := expression.Parse(rf.Answer)
	if err != nil {
		return nil, err
	}

	return radioValidator{expr: expr, proposalsLength: len(rf.Proposals)}, nil
}

type OrderedListFieldBlock struct {
	Label               Interpolated
	Answer              []Interpolated // the order matters
	AdditionalProposals []Interpolated
}

func (ol OrderedListFieldBlock) instantiate(params expression.Vars, ID int) (instance, error) {
	out := OrderedListFieldInstance{
		Answer:              make([]client.TextLine, len(ol.Answer)),
		AdditionalProposals: make([]client.TextLine, len(ol.AdditionalProposals)),
		ID:                  ID,
	}

	var err error
	out.Label, err = ol.Label.instantiateAndMerge(params)
	if err != nil {
		return nil, err
	}

	for i, a := range ol.Answer {
		out.Answer[i], err = a.instantiate(params)
		if err != nil {
			return nil, err
		}
	}

	for i, a := range ol.AdditionalProposals {
		out.AdditionalProposals[i], err = a.instantiate(params)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

func (ol OrderedListFieldBlock) setupValidator(expression.RandomParameters) (validator, error) {
	_, err := ol.Label.parse()
	if err != nil {
		return nil, err
	}

	for _, a := range ol.Answer {
		if _, err := a.parse(); err != nil {
			return nil, err
		}
	}

	for _, a := range ol.AdditionalProposals {
		if _, err := a.parse(); err != nil {
			return nil, err
		}
	}

	return noOpValidator{}, nil
}

// CoordExpression is a pair of valid expression.Expression
type CoordExpression struct {
	X, Y string
}

func (c CoordExpression) instantiateToFloat(params expression.Vars) (repere.Coord, error) {
	x, err := evaluateExpr(c.X, params)
	if err != nil {
		return repere.Coord{}, err
	}
	y, err := evaluateExpr(c.Y, params)
	if err != nil {
		return repere.Coord{}, err
	}
	return repere.Coord{
		X: x,
		Y: y,
	}, nil
}

func (c CoordExpression) instantiate(params expression.Vars) (repere.IntCoord, error) {
	out, err := c.instantiateToFloat(params)
	return out.Round(), err
}

func (c CoordExpression) parse() (out parsedCoord, err error) {
	out.X, err = expression.Parse(c.X)
	if err != nil {
		return out, err
	}
	out.Y, err = expression.Parse(c.Y)
	if err != nil {
		return out, err
	}
	return out, nil
}

type FigurePointFieldBlock struct {
	Answer CoordExpression
	Figure FigureBlock
}

func (fp FigurePointFieldBlock) instantiate(params expression.Vars, ID int) (instance, error) {
	fig, err := fp.Figure.instantiateF(params)
	if err != nil {
		return nil, err
	}
	ans, err := fp.Answer.instantiate(params)
	if err != nil {
		return nil, err
	}
	return FigurePointFieldInstance{
		Figure: fig.Figure,
		Answer: ans,
		ID:     ID,
	}, nil
}

func (fp FigurePointFieldBlock) setupValidator(params expression.RandomParameters) (validator, error) {
	figure, err := fp.Figure.setupValidator(params)
	if err != nil {
		return nil, err
	}

	answer, err := fp.Answer.parse()
	if err != nil {
		return nil, err
	}

	return figurePointValidator{figure: figure, answer: answer}, nil
}

type FigureVectorFieldBlock struct {
	Answer CoordExpression

	AnswerOrigin CoordExpression // optionnal, used when MustHaveOrigin is true

	Figure FigureBlock

	MustHaveOrigin bool
}

func (fv FigureVectorFieldBlock) instantiate(params expression.Vars, ID int) (instance, error) {
	fig, err := fv.Figure.instantiateF(params)
	if err != nil {
		return nil, err
	}
	ans, err := fv.Answer.instantiate(params)
	if err != nil {
		return nil, err
	}

	out := FigureVectorFieldInstance{
		ID:             ID,
		Figure:         fig.Figure,
		Answer:         ans,
		MustHaveOrigin: fv.MustHaveOrigin,
	}

	if fv.MustHaveOrigin {
		out.AnswerOrigin, err = fv.AnswerOrigin.instantiate(params)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (fp FigureVectorFieldBlock) setupValidator(params expression.RandomParameters) (validator, error) {
	figure, err := fp.Figure.setupValidator(params)
	if err != nil {
		return nil, err
	}

	answer, err := fp.Answer.parse()
	if err != nil {
		return nil, err
	}

	out := figureVectorValidator{figure: figure, answer: answer}

	if fp.MustHaveOrigin {
		origin, err := fp.AnswerOrigin.parse()
		if err != nil {
			return nil, err
		}
		out.answerOrigin = &origin
	}

	return out, nil
}

type VariationTableFieldBlock struct {
	Answer VariationTableBlock
}

func (vt VariationTableFieldBlock) instantiate(params expression.Vars, ID int) (instance, error) {
	ans, err := vt.Answer.instantiateVT(params)
	return VariationTableFieldInstance{
		ID:     ID,
		Answer: ans,
	}, err
}

func (fp VariationTableFieldBlock) setupValidator(params expression.RandomParameters) (validator, error) {
	return fp.Answer.setupValidator(params)
}

type SignTableFieldBlock struct {
	Answer SignTableBlock
}

func (vt SignTableFieldBlock) instantiate(params expression.Vars, ID int) (instance, error) {
	ans, err := vt.Answer.instantiateST(params)
	return SignTableFieldInstance{
		ID:     ID,
		Answer: ans,
	}, err
}

func (fp SignTableFieldBlock) setupValidator(params expression.RandomParameters) (validator, error) {
	return fp.Answer.setupValidator(params)
}

type FunctionPointsFieldBlock struct {
	IsDiscrete bool
	Function   string // function or sequence, valid expression.Expression
	Label      string
	Variable   expression.Variable
	XGrid      []string // valid expression.Expression
}

func (fp FunctionPointsFieldBlock) instantiate(params expression.Vars, ID int) (instance, error) {
	fn, err := expression.Parse(fp.Function)
	if err != nil {
		return nil, err
	}
	fn.Substitute(params)

	xGrid := make([]int, len(fp.XGrid))
	for i, x := range fp.XGrid {
		v, err := evaluateExpr(x, params)
		if err != nil {
			return nil, err
		}
		xGrid[i] = int(v)
	}
	sort.Ints(xGrid)

	return FunctionPointsFieldInstance{
		IsDiscrete: fp.IsDiscrete,
		Function: expression.FunctionExpr{
			Function: fn,
			Variable: fp.Variable,
		},
		ID:           ID,
		Label:        fp.Label,
		XGrid:        xGrid,
		offsetHeight: rand.Intn(3),
	}, nil
}

func (fp FunctionPointsFieldBlock) setupValidator(params expression.RandomParameters) (validator, error) {
	if len(fp.XGrid) < 2 {
		return nil, errors.New("Au moins deux valeurs pour x doivent être précisées.")
	}

	var (
		out functionPointsValidator
		err error
	)

	out.xGrid = make([]*expression.Expr, len(fp.XGrid))
	for i, x := range fp.XGrid {
		out.xGrid[i], err = expression.Parse(x)
		if err != nil {
			return nil, err
		}
	}

	fn := FunctionDefinition{
		Function: fp.Function,
		Variable: fp.Variable,
		From:     fp.XGrid[0],
		To:       fp.XGrid[len(fp.XGrid)-1],
	}
	out.function, err = newFunction(fn, params)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type FigureVectorPairFieldBlock struct {
	Figure    FigureBlock
	Criterion VectorPairCriterion
}

func (fv FigureVectorPairFieldBlock) instantiate(params expression.Vars, ID int) (instance, error) {
	fig, err := fv.Figure.instantiateF(params)
	return FigureVectorPairFieldInstance{
		ID:        ID,
		Figure:    fig.Figure,
		Criterion: fv.Criterion,
	}, err
}

func (fp FigureVectorPairFieldBlock) setupValidator(params expression.RandomParameters) (validator, error) {
	return fp.Figure.setupValidator(params)
}

type FigureAffineLineFieldBlock struct {
	Label  string
	A      string // valid expression.Expression
	B      string // valid expression.Expression
	Figure FigureBlock
}

func (fa FigureAffineLineFieldBlock) instantiate(params expression.Vars, ID int) (instance, error) {
	fig, err := fa.Figure.instantiateF(params)
	if err != nil {
		return nil, err
	}
	ansA, err := evaluateExpr(fa.A, params)
	if err != nil {
		return nil, err
	}
	ansB, err := evaluateExpr(fa.B, params)
	if err != nil {
		return nil, err
	}
	return FigureAffineLineFieldInstance{
		ID:      ID,
		Label:   fa.Label,
		Figure:  fig.Figure,
		AnswerA: ansA,
		AnswerB: int(ansB),
	}, nil
}

func (fa FigureAffineLineFieldBlock) setupValidator(params expression.RandomParameters) (validator, error) {
	figure, err := fa.Figure.setupValidator(params)
	if err != nil {
		return nil, err
	}

	a, err := expression.Parse(fa.A)
	if err != nil {
		return nil, err
	}
	b, err := expression.Parse(fa.B)
	if err != nil {
		return nil, err
	}
	return figureAffineLineValidator{figure: figure, a: a, b: b}, nil
}

type TreeFieldBlock struct {
	Answer TreeBlock
}

func (tf TreeFieldBlock) instantiate(params expression.Vars, ID int) (instance, error) {
	answer, err := tf.Answer.instantiateT(params)
	return TreeFieldInstance{
		ID:     ID,
		Answer: answer,
	}, err
}

func (tf TreeFieldBlock) setupValidator(params expression.RandomParameters) (validator, error) {
	return tf.Answer.setupValidator(params)
}

type TableFieldBlock struct {
	HorizontalHeaders []TextPart
	VerticalHeaders   []TextPart
	Answer            [][]string // valid expression.Expression
}

func (tf TableFieldBlock) instantiate(params expression.Vars, ID int) (instance, error) {
	out := TableFieldInstance{
		ID:                ID,
		HorizontalHeaders: make([]client.TextOrMath, len(tf.HorizontalHeaders)),
		VerticalHeaders:   make([]client.TextOrMath, len(tf.VerticalHeaders)),
		Answer:            client.TableAnswer{Rows: make([][]float64, len(tf.Answer))},
	}
	var err error
	for i, cell := range tf.HorizontalHeaders {
		out.HorizontalHeaders[i], err = cell.instantiate(params)
		if err != nil {
			return nil, err
		}
	}
	for i, cell := range tf.VerticalHeaders {
		out.VerticalHeaders[i], err = cell.instantiate(params)
		if err != nil {
			return nil, err
		}
	}

	for i, row := range tf.Answer {
		rowInstance := make([]float64, len(row))
		for j, v := range row {
			rowInstance[j], err = evaluateExpr(v, params)
			if err != nil {
				return nil, err
			}
		}
		out.Answer.Rows[i] = rowInstance
	}
	return out, nil
}

func (tf TableFieldBlock) setupValidator(params expression.RandomParameters) (validator, error) {
	for _, cell := range tf.HorizontalHeaders {
		if err := cell.validate(); err != nil {
			return nil, err
		}
	}
	for _, cell := range tf.VerticalHeaders {
		if err := cell.validate(); err != nil {
			return nil, err
		}
	}

	out := tableValidator{answer: make([][]*expression.Expr, len(tf.Answer))}
	var err error
	for i, row := range tf.Answer {
		rowExpr := make([]*expression.Expr, len(row))
		for j, cell := range row {
			rowExpr[j], err = expression.Parse(cell)
			if err != nil {
				return nil, err
			}
		}
		out.answer[i] = rowExpr
	}

	return out, nil
}

// VectorFieldBlock is a two-number field, with
// option to interpret the answer up to colinearity
type VectorFieldBlock struct {
	Answer         CoordExpression
	AcceptColinear bool // if true, all vectors colinears to `Answer` are accepted
	DisplayColumn  bool // if true, the field are displayed in column, instead of being on the same line
}

func (v VectorFieldBlock) instantiate(params expression.Vars, ID int) (instance, error) {
	ans, err := v.Answer.instantiateToFloat(params)
	if err != nil {
		return nil, err
	}

	out := VectorFieldInstance{
		ID:             ID,
		Answer:         ans,
		AcceptColinear: v.AcceptColinear,
		DisplayColumn:  v.DisplayColumn,
	}
	return out, nil
}

func (v VectorFieldBlock) setupValidator(expression.RandomParameters) (validator, error) {
	answer, err := v.Answer.parse()
	if err != nil {
		return nil, err
	}

	return vectorValidator{answer: answer}, nil
}
