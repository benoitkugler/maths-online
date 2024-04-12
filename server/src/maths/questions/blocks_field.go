package questions

import (
	"errors"
	"math/rand"
	"sort"

	ex "github.com/benoitkugler/maths-online/server/src/maths/expression"
	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	"github.com/benoitkugler/maths-online/server/src/maths/repere"
)

var (
	_ Block = NumberFieldBlock{}
	_ Block = ExpressionFieldBlock{}
	_ Block = RadioFieldBlock{}
	_ Block = OrderedListFieldBlock{}
	_ Block = GeometricConstructionFieldBlock{}
	_ Block = VariationTableFieldBlock{}
	_ Block = SignTableFieldBlock{}
	_ Block = FunctionPointsFieldBlock{}
	_ Block = TreeFieldBlock{}
	_ Block = TableFieldBlock{}
	_ Block = VectorFieldBlock{}
	_ Block = ProofFieldBlock{}
	_ Block = SetFieldBlock{}
)

type NumberFieldBlock struct {
	// a valid expression, in the format used by expression.Expression
	// which is only parametrized by the random parameters
	Expression string
}

func (n NumberFieldBlock) instantiate(params ex.Vars, ID int) (instance, error) {
	answer, err := evaluateExpr(n.Expression, params)
	return NumberFieldInstance{ID: ID, Answer: answer}, err
}

func (n NumberFieldBlock) setupValidator(params *ex.RandomParameters) (validator, error) {
	expr, err := ex.Parse(n.Expression)
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
	answer, err := ex.ParseCompound(f.Expression)
	if err != nil {
		return TextBlock{}, err
	}
	allExprs := answer.Expressions()

	if err := params.Validate(); err != nil {
		return TextBlock{}, err
	}
	m := params.ToMap()

	const nbRepeat = 100
	allHints := make(ex.SyntaxHints)
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

func (f ExpressionFieldBlock) instantiate(params ex.Vars, ID int) (instance, error) {
	answer, err := ex.ParseCompound(f.Expression)
	if err != nil {
		return nil, err
	}
	answer.Substitute(params)

	var showFractionHelp bool
	answerExpr, ok := answer.(*ex.Expr)
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

func (f ExpressionFieldBlock) setupValidator(*ex.RandomParameters) (validator, error) {
	expr, err := ex.ParseCompound(f.Expression)
	if err != nil {
		return nil, err
	}
	_, err = f.Label.parse()
	if err != nil {
		return nil, err
	}

	asExpr, isExpr := expr.(*ex.Expr)

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

func (rf RadioFieldBlock) instantiate(params ex.Vars, ID int) (instance, error) {
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

func (rf RadioFieldBlock) setupValidator(params *ex.RandomParameters) (validator, error) {
	for _, p := range rf.Proposals {
		_, err := p.parse()
		if err != nil {
			return nil, err
		}
	}

	expr, err := ex.Parse(rf.Answer)
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

func (ol OrderedListFieldBlock) instantiate(params ex.Vars, ID int) (instance, error) {
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

func (ol OrderedListFieldBlock) setupValidator(*ex.RandomParameters) (validator, error) {
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

func (c CoordExpression) instantiateToFloat(params ex.Vars) (repere.Coord, error) {
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

func (c CoordExpression) instantiate(params ex.Vars) (repere.IntCoord, error) {
	out, err := c.instantiateToFloat(params)
	return out.Round(), err
}

func (c CoordExpression) parse() (out parsedCoord, err error) {
	out.X, err = ex.Parse(c.X)
	if err != nil {
		return out, err
	}
	out.Y, err = ex.Parse(c.Y)
	if err != nil {
		return out, err
	}
	return out, nil
}

func (fv GeometricConstructionFieldBlock) setupValidator(params *ex.RandomParameters) (validator, error) {
	background, err := fv.Background.setupValidator(params)
	if err != nil {
		return nil, err
	}

	field, err := fv.Field.setupValidator(params)
	if err != nil {
		return nil, err
	}

	return geometricConstructionValidator{field: field, background: background}, nil
}

func (fv GeometricConstructionFieldBlock) instantiate(params ex.Vars, ID int) (instance, error) {
	var (
		out GeometricConstructionFieldInstance
		err error
	)

	out.ID = ID

	out.Field, err = fv.Field.instantiate(params)
	if err != nil {
		return nil, err
	}

	out.Background, err = fv.Background.instantiateFG(params)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type GFPoint struct {
	Answer CoordExpression
}

type GFVector struct {
	Answer         CoordExpression
	AnswerOrigin   CoordExpression // optionnal, used when MustHaveOrigin is true
	MustHaveOrigin bool
}

type GFVectorPair struct {
	Criterion VectorPairCriterion
}

type GFAffineLine struct {
	Label string
	A     string // valid expression.Expression
	B     string // valid expression.Expression
}

func (f FigureBlock) instantiateFG(params ex.Vars) (client.FigureOrGraph, error) {
	fig, err := f.instantiateF(params)
	return client.FigureBlock(fig), err
}

func (f FunctionsGraphBlock) instantiateFG(params ex.Vars) (client.FigureOrGraph, error) {
	graphs, err := f.instantiateG(params)
	return graphs.toClientG(), err
}

func (fp GFPoint) instantiate(params ex.Vars) (geoFieldInstance, error) {
	ans, err := fp.Answer.instantiate(params)
	if err != nil {
		return nil, err
	}
	return gfPoint(ans), nil
}

func (fp GFPoint) setupValidator(params *ex.RandomParameters) (validator, error) {
	answer, err := fp.Answer.parse()
	if err != nil {
		return nil, err
	}

	return gfPointValidator(answer), nil
}

func (fv GFVector) instantiate(params ex.Vars) (geoFieldInstance, error) {
	ans, err := fv.Answer.instantiate(params)
	if err != nil {
		return nil, err
	}

	out := gfVector{
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

func (fp GFVector) setupValidator(params *ex.RandomParameters) (validator, error) {
	answer, err := fp.Answer.parse()
	if err != nil {
		return nil, err
	}

	out := gfVectorValidator{answer: answer}

	if fp.MustHaveOrigin {
		origin, err := fp.AnswerOrigin.parse()
		if err != nil {
			return nil, err
		}
		out.answerOrigin = &origin
	}

	return out, nil
}

func (fv GFVectorPair) instantiate(params ex.Vars) (geoFieldInstance, error) {
	return gfVectorPair(fv.Criterion), nil
}

func (fp GFVectorPair) setupValidator(params *ex.RandomParameters) (validator, error) {
	return noOpValidator{}, nil
}

func (fa GFAffineLine) instantiate(params ex.Vars) (geoFieldInstance, error) {
	ansA, err := evaluateExpr(fa.A, params)
	if err != nil {
		return nil, err
	}
	ansB, err := evaluateExpr(fa.B, params)
	if err != nil {
		return nil, err
	}
	return gfAffineLine{
		Label:   fa.Label,
		AnswerA: ansA,
		AnswerB: int(ansB),
	}, nil
}

func (fa GFAffineLine) setupValidator(params *ex.RandomParameters) (validator, error) {
	a, err := ex.Parse(fa.A)
	if err != nil {
		return nil, err
	}
	b, err := ex.Parse(fa.B)
	if err != nil {
		return nil, err
	}
	return gfAffineLineValidator{a: a, b: b}, nil
}

type VariationTableFieldBlock struct {
	Answer VariationTableBlock
}

func (vt VariationTableFieldBlock) instantiate(params ex.Vars, ID int) (instance, error) {
	ans, err := vt.Answer.instantiateVT(params)
	return VariationTableFieldInstance{
		ID:     ID,
		Answer: ans,
	}, err
}

func (fp VariationTableFieldBlock) setupValidator(params *ex.RandomParameters) (validator, error) {
	return fp.Answer.setupValidator(params)
}

type SignTableFieldBlock struct {
	Answer SignTableBlock
}

func (vt SignTableFieldBlock) instantiate(params ex.Vars, ID int) (instance, error) {
	ans, err := vt.Answer.instantiateST(params)
	return SignTableFieldInstance{
		ID:     ID,
		Answer: ans,
	}, err
}

func (fp SignTableFieldBlock) setupValidator(params *ex.RandomParameters) (validator, error) {
	return fp.Answer.setupValidator(params)
}

type FunctionPointsFieldBlock struct {
	IsDiscrete bool
	Function   string // function or sequence, valid expression.Expression
	Label      string
	Variable   ex.Variable
	XGrid      []string // valid expression.Expression
}

func (fp FunctionPointsFieldBlock) instantiate(params ex.Vars, ID int) (instance, error) {
	fn, err := ex.Parse(fp.Function)
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
		Function: ex.FunctionExpr{
			Function: fn,
			Variable: fp.Variable,
		},
		ID:           ID,
		Label:        fp.Label,
		XGrid:        xGrid,
		offsetHeight: rand.Intn(3),
	}, nil
}

func (fp FunctionPointsFieldBlock) setupValidator(params *ex.RandomParameters) (validator, error) {
	if len(fp.XGrid) < 2 {
		return nil, errors.New("Au moins deux valeurs pour x doivent être précisées.")
	}

	var (
		out functionPointsValidator
		err error
	)

	out.xGrid = make([]*ex.Expr, len(fp.XGrid))
	for i, x := range fp.XGrid {
		out.xGrid[i], err = ex.Parse(x)
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
	out.function, err = newFunctionValidator(fn, params)
	if err != nil {
		return nil, err
	}

	return out, nil
}

type TreeFieldBlock struct {
	Answer TreeBlock
}

func (tf TreeFieldBlock) instantiate(params ex.Vars, ID int) (instance, error) {
	answer, err := tf.Answer.instantiateT(params)
	return TreeFieldInstance{
		ID:     ID,
		Answer: answer,
	}, err
}

func (tf TreeFieldBlock) setupValidator(params *ex.RandomParameters) (validator, error) {
	return tf.Answer.setupValidator(params)
}

type TableFieldBlock struct {
	HorizontalHeaders []TextPart
	VerticalHeaders   []TextPart
	Answer            [][]string // valid expression.Expression
}

func (tf TableFieldBlock) instantiate(params ex.Vars, ID int) (instance, error) {
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

func (tf TableFieldBlock) setupValidator(params *ex.RandomParameters) (validator, error) {
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

	out := tableValidator{answer: make([][]*ex.Expr, len(tf.Answer))}
	var err error
	for i, row := range tf.Answer {
		rowExpr := make([]*ex.Expr, len(row))
		for j, cell := range row {
			rowExpr[j], err = ex.Parse(cell)
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

func (v VectorFieldBlock) instantiate(params ex.Vars, ID int) (instance, error) {
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

func (v VectorFieldBlock) setupValidator(*ex.RandomParameters) (validator, error) {
	answer, err := v.Answer.parse()
	if err != nil {
		return nil, err
	}

	return vectorValidator{answer: answer}, nil
}

type SetFieldBlock struct {
	Answer         string // expression
	AdditionalSets []Interpolated
}

func (v SetFieldBlock) instantiate(params ex.Vars, ID int) (instance, error) {
	answer, err := ex.Parse(v.Answer)
	if err != nil {
		return nil, err
	}
	answer.Substitute(params)

	setExpr, err := answer.ToBinarySet()
	if err != nil {
		return nil, err
	}
	// add the sets to the answer
	for _, s := range v.AdditionalSets {
		setLatex, err := s.instantiateAndMerge(params)
		if err != nil {
			return nil, err
		}
		setExpr.Sets = append(setExpr.Sets, setLatex)
	}

	// the order of the sets presented is shuffled on the client

	out := SetFieldInstance{
		ID:     ID,
		Answer: setExpr,
	}
	return out, nil
}

func (v SetFieldBlock) setupValidator(*ex.RandomParameters) (validator, error) {
	answer, err := ex.Parse(v.Answer)
	if err != nil {
		return nil, err
	}

	for _, set := range v.AdditionalSets {
		_, err := set.parse()
		if err != nil {
			return nil, err
		}
	}

	return setValidator{answer: answer}, nil
}
