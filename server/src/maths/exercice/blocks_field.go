package exercice

import (
	"errors"
	"fmt"
	"sort"

	"github.com/benoitkugler/maths-online/maths/exercice/client"
	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/repere"
)

var (
	_ Block = NumberFieldBlock{}
	_ Block = ExpressionFieldBlock{}
	_ Block = RadioFieldBlock{}
	_ Block = OrderedListFieldBlock{}
	_ Block = FigurePointFieldBlock{}
	_ Block = FigureVectorFieldBlock{}
	_ Block = VariationTableFieldBlock{}
	_ Block = FunctionPointsFieldBlock{}
	_ Block = FigureVectorPairFieldBlock{}
	_ Block = FigureAffineLineFieldBlock{}
	_ Block = TreeFieldBlock{}
	_ Block = TableFieldBlock{}
	_ Block = VectorFieldBlock{}
)

type NumberFieldBlock struct {
	// a valid expression, in the format used by expression.Expression
	// which is only parametrized by the random parameters
	Expression string
}

func (n NumberFieldBlock) instantiate(params expression.Variables, ID int) (instance, error) {
	answer, err := evaluateExpr(n.Expression, params)
	return NumberFieldInstance{ID: ID, Answer: answer}, err
}

func (n NumberFieldBlock) validate(params expression.RandomParameters) error {
	// note that we dont allow non decimal solutions, since it is confusing for the student.
	// they should rather be handled with an expression field, or rounded using the
	// builtin round() function
	return validateNumberExpression(n.Expression, params, true, true)
}

type ExpressionFieldBlock struct {
	Expression      string   // a valid expression, in the format used by expression.Expression
	Label           TextPart // optional
	ComparisonLevel ComparisonLevel
}

func (f ExpressionFieldBlock) instantiate(params expression.Variables, ID int) (instance, error) {
	label := StringOrExpression{String: f.Label.Content}
	if f.Label.Kind == Expression {
		e, err := expression.Parse(f.Label.Content)
		if err != nil {
			return nil, err
		}
		label = StringOrExpression{Expression: e}
		label.Expression.Substitute(params)
	}
	answer, err := expression.Parse(f.Expression)
	if err != nil {
		return nil, err
	}
	answer.Substitute(params)
	return ExpressionFieldInstance{
		Label:           label,
		Answer:          answer,
		ComparisonLevel: f.ComparisonLevel,
		ID:              ID,
	}, nil
}

func (f ExpressionFieldBlock) validate(params expression.RandomParameters) error {
	_, err := expression.Parse(f.Expression)
	return err
}

type RadioFieldBlock struct {
	Answer     string         // must satisfy expression.IsValidIndex
	Proposals  []Interpolated // slice of text parts
	AsDropDown bool
}

func (rf RadioFieldBlock) instantiate(params expression.Variables, ID int) (instance, error) {
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

func (rf RadioFieldBlock) validate(params expression.RandomParameters) error {
	for _, p := range rf.Proposals {
		_, err := p.Parse()
		if err != nil {
			return err
		}
	}

	expr, err := expression.Parse(rf.Answer)
	if err != nil {
		return err
	}
	if ok, freq := expr.IsValidIndex(params, len(rf.Proposals)); !ok {
		return fmt.Errorf("L'expression %s ne définit pas un index valide dans la liste des propositions (%d %% des tests ont échoué)", rf.Answer, 100-freq)
	}

	return nil
}

type OrderedListFieldBlock struct {
	Label               string         // optionnal, LaTeX code displayed in front of the anwser field
	Answer              []Interpolated // the order matters
	AdditionalProposals []Interpolated
}

func (ol OrderedListFieldBlock) instantiate(params expression.Variables, ID int) (instance, error) {
	out := OrderedListFieldInstance{
		Label:               ol.Label,
		Answer:              make([]client.TextLine, len(ol.Answer)),
		AdditionalProposals: make([]client.TextLine, len(ol.AdditionalProposals)),
		ID:                  ID,
	}

	var err error
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

func (ol OrderedListFieldBlock) validate(params expression.RandomParameters) error {
	for _, a := range ol.Answer {
		if _, err := a.Parse(); err != nil {
			return err
		}
	}

	for _, a := range ol.AdditionalProposals {
		if _, err := a.Parse(); err != nil {
			return err
		}
	}

	return nil
}

// CoordExpression is a pair of valid expression.Expression
type CoordExpression struct {
	X, Y string
}

func (c CoordExpression) instantiateToFloat(params expression.Variables) (repere.Coord, error) {
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

func (c CoordExpression) instantiate(params expression.Variables) (repere.IntCoord, error) {
	out, err := c.instantiateToFloat(params)
	return out.Round(), err
}

func (c CoordExpression) validate(params expression.RandomParameters, checkPrecision bool) error {
	if err := validateNumberExpression(c.X, params, checkPrecision, true); err != nil {
		return err
	}
	if err := validateNumberExpression(c.Y, params, checkPrecision, true); err != nil {
		return err
	}
	return nil
}

type FigurePointFieldBlock struct {
	Answer CoordExpression
	Figure FigureBlock
}

func (fp FigurePointFieldBlock) instantiate(params expression.Variables, ID int) (instance, error) {
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

func (fp FigurePointFieldBlock) validate(params expression.RandomParameters) error {
	if err := fp.Figure.validate(params); err != nil {
		return err
	}
	if err := fp.Answer.validate(params, false); err != nil {
		return err
	}
	return nil
}

type FigureVectorFieldBlock struct {
	Answer CoordExpression

	AnswerOrigin CoordExpression // optionnal, used when MustHaveOrigin is true

	Figure FigureBlock

	MustHaveOrigin bool
}

func (fv FigureVectorFieldBlock) instantiate(params expression.Variables, ID int) (instance, error) {
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

func (fp FigureVectorFieldBlock) validate(params expression.RandomParameters) error {
	if err := fp.Figure.validate(params); err != nil {
		return err
	}
	if err := fp.Answer.validate(params, false); err != nil {
		return err
	}
	if fp.MustHaveOrigin {
		if err := fp.AnswerOrigin.validate(params, false); err != nil {
			return err
		}
	}

	return nil
}

type VariationTableFieldBlock struct {
	Answer VariationTableBlock
}

func (vt VariationTableFieldBlock) instantiate(params expression.Variables, ID int) (instance, error) {
	ans, err := vt.Answer.instantiateVT(params)
	return VariationTableFieldInstance{
		ID:     ID,
		Answer: ans,
	}, err
}

func (fp VariationTableFieldBlock) validate(params expression.RandomParameters) error {
	return fp.Answer.validate(params)
}

type FunctionPointsFieldBlock struct {
	Function string // valid expression.Expression
	Label    string
	Variable expression.Variable
	XGrid    []int
}

func (fp FunctionPointsFieldBlock) instantiate(params expression.Variables, ID int) (instance, error) {
	fn, err := expression.Parse(fp.Function)
	if err != nil {
		return nil, err
	}
	fn.Substitute(params)
	return FunctionPointsFieldInstance{
		Function: expression.FunctionExpr{
			Function: fn,
			Variable: fp.Variable,
		},
		ID:    ID,
		Label: fp.Label,
		XGrid: fp.XGrid,
	}, nil
}

func (fp FunctionPointsFieldBlock) validate(params expression.RandomParameters) error {
	if !sort.IntsAreSorted(fp.XGrid) {
		return errors.New("Les valeurs x doivent être en ordre croissant.")
	}

	if len(fp.XGrid) < 2 {
		return errors.New("Au moins deux valeurs pour x doivent être précisées.")
	}

	fn := FunctionDefinition{
		Function: fp.Function,
		Variable: fp.Variable,
		Range: [2]float64{
			float64(fp.XGrid[0]),
			float64(fp.XGrid[len(fp.XGrid)-1]),
		},
	}
	if err := fn.validate(params); err != nil {
		return err
	}

	// check that every point is an integer (otherwise is can't be selected on the client)
	fnExpr, _ := fn.parse() // guarded by `validate`
	if ok, freq := fnExpr.AreFxsIntegers(params, fp.XGrid); !ok {
		return fmt.Errorf("Les valeurs de la fonction ne sont pas des nombres entiers (%d %% des tests ont échoué).", 100-freq)
	}

	return nil
}

type FigureVectorPairFieldBlock struct {
	Figure    FigureBlock
	Criterion VectorPairCriterion
}

func (fv FigureVectorPairFieldBlock) instantiate(params expression.Variables, ID int) (instance, error) {
	fig, err := fv.Figure.instantiateF(params)
	return FigureVectorPairFieldInstance{
		ID:        ID,
		Figure:    fig.Figure,
		Criterion: fv.Criterion,
	}, err
}

func (fp FigureVectorPairFieldBlock) validate(params expression.RandomParameters) error {
	return fp.Figure.validate(params)
}

type FigureAffineLineFieldBlock struct {
	Label  string
	A      string // valid expression.Expression
	B      string // valid expression.Expression
	Figure FigureBlock
}

func (fa FigureAffineLineFieldBlock) instantiate(params expression.Variables, ID int) (instance, error) {
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

func (fa FigureAffineLineFieldBlock) validate(params expression.RandomParameters) error {
	if err := fa.Figure.validate(params); err != nil {
		return err
	}
	if err := validateNumberExpression(fa.A, params, false, false); err != nil {
		return err
	}
	if err := validateNumberExpression(fa.B, params, false, true); err != nil {
		return err
	}

	bExpr := expression.MustParse(fa.B) // guarded by `validateNumberExpression`
	if ok, freq := bExpr.IsValidInteger(params); !ok {
		return fmt.Errorf("L'expression de B n'est pas un nombre entier (%d %% des tests ont échoué).", 100-freq)
	}
	return nil
}

type TreeNodeAnswer struct {
	Children      []TreeNodeAnswer `structgen-data:"ignore"`
	Probabilities []string         // edges, same length as Children, valid expression.Expression
	Value         int              // index into the proposals, 0 for the root
}

type TreeFieldBlock struct {
	EventsProposals []string
	AnswerRoot      TreeNodeAnswer
}

func (tf TreeFieldBlock) instantiate(params expression.Variables, ID int) (instance, error) {
	out := TreeFieldInstance{
		ID:              ID,
		EventsProposals: make([]client.TextOrMath, len(tf.EventsProposals)),
	}
	for i, p := range tf.EventsProposals {
		out.EventsProposals[i] = client.TextOrMath{Text: p}
	}

	var buildTree func(node TreeNodeAnswer) (client.TreeNodeAnswer, error)
	buildTree = func(node TreeNodeAnswer) (client.TreeNodeAnswer, error) {
		out := client.TreeNodeAnswer{
			Value:         node.Value,
			Probabilities: make([]float64, len(node.Probabilities)),
			Children:      make([]client.TreeNodeAnswer, len(node.Children)),
		}
		var err error
		for i, c := range node.Probabilities {
			out.Probabilities[i], err = evaluateExpr(c, params)
			if err != nil {
				return out, err
			}
		}
		for i, c := range node.Children {
			out.Children[i], err = buildTree(c)
			if err != nil {
				return out, err
			}
		}
		return out, nil
	}

	root, err := buildTree(tf.AnswerRoot)
	out.Answer = client.TreeAnswer{Root: root}
	return out, err
}

func (tf TreeFieldBlock) validate(params expression.RandomParameters) error {
	var checkTree func(node TreeNodeAnswer) error
	checkTree = func(node TreeNodeAnswer) error {
		if node.Value < 0 || node.Value >= len(tf.EventsProposals) {
			return fmt.Errorf("L'index %d n'est pas compatible avec le nombre de propositions.", node.Value)
		}

		for _, c := range node.Probabilities {
			expr, err := expression.Parse(c)
			if err != nil {
				return err
			}

			if ok, freq := expr.IsValidProba(params); !ok {
				return fmt.Errorf("L'expression %s ne définit pas une probabilité valide. (%d %% des tests ont échoué)", c, 100-freq)
			}
		}
		for _, c := range node.Children {
			if err := checkTree(c); err != nil {
				return err
			}
		}
		return nil
	}
	return checkTree(tf.AnswerRoot)
}

type TableFieldBlock struct {
	HorizontalHeaders []TextPart
	VerticalHeaders   []TextPart
	Answer            [][]string // valid expression.Expression
}

func (tf TableFieldBlock) instantiate(params expression.Variables, ID int) (instance, error) {
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

func (tf TableFieldBlock) validate(params expression.RandomParameters) error {
	for _, cell := range tf.HorizontalHeaders {
		if err := cell.validate(); err != nil {
			return err
		}
	}
	for _, cell := range tf.VerticalHeaders {
		if err := cell.validate(); err != nil {
			return err
		}
	}
	for _, row := range tf.Answer {
		for _, cell := range row {
			if err := validateNumberExpression(cell, params, true, true); err != nil {
				return err
			}
		}
	}

	return nil
}

// VectorFieldBlock is a two-number field, with
// option to interpret the answer up to colinearity
type VectorFieldBlock struct {
	Answer         CoordExpression
	AcceptColinear bool // if true, all vectors colinears to `Answer` are accepted
	DisplayColumn  bool // if true, the field are displayed in column, instead of being on the same line
}

func (v VectorFieldBlock) instantiate(params expression.Variables, ID int) (instance, error) {
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

func (v VectorFieldBlock) validate(params expression.RandomParameters) error {
	if err := v.Answer.validate(params, true); err != nil {
		return err
	}

	return nil
}
