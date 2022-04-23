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
	_ Block = FormulaFieldBlock{}
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
)

type NumberFieldBlock struct {
	// a valid expression, in the format used by expression.Expression
	// which is only parametrized by the random parameters
	Expression string
}

func (n NumberFieldBlock) instantiate(params expression.Variables, ID int) instance {
	answer := expression.MustEvaluate(n.Expression, params)
	return NumberFieldInstance{ID: ID, Answer: answer}
}

func (n NumberFieldBlock) validate(params expression.RandomParameters) error {
	// note that we dont allow non decimal solutions, since it is confusing for the student.
	// they should rahter be handled with an expression field
	return validateNumberExpression(n.Expression, params, true)
}

type FormulaFieldBlock struct {
	Expression      string   // a valid expression, in the format used by expression.Expression
	Label           TextPart // optional
	ComparisonLevel ComparisonLevel
}

func (f FormulaFieldBlock) instantiate(params expression.Variables, ID int) instance {
	label := StringOrExpression{String: f.Label.Content}
	if f.Label.Kind == Expression {
		label = StringOrExpression{Expression: mustParse(f.Label.Content)}
		label.Expression.Substitute(params)
	}
	return ExpressionFieldInstance{
		Label:           label,
		Answer:          mustParse(f.Expression),
		ComparisonLevel: f.ComparisonLevel,
		ID:              ID,
	}
}

func (n FormulaFieldBlock) validate(params expression.RandomParameters) error {
	_, err := expression.Parse(n.Expression)
	return err
}

type RadioFieldBlock struct {
	Answer     string         // must satisfy expression.IsValidIndex
	Proposals  []Interpolated // slice of text parts
	AsDropDown bool
}

func (rf RadioFieldBlock) instantiate(params expression.Variables, ID int) instance {
	out := RadioFieldInstance{
		Proposals: make([]client.ListFieldProposal, len(rf.Proposals)),
		Answer:    int(expression.MustEvaluate(rf.Answer, params)),
		ID:        ID,
	}
	for i, p := range rf.Proposals {
		parts, _ := p.Parse()
		out.Proposals[i] = client.ListFieldProposal{Content: parts.instantiate(params)}
	}

	if rf.AsDropDown {
		return DropDownFieldInstance(out)
	}
	return out
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
	Label               string     // optionnal, LaTeX code displayed in front of the anwser field
	Answer              []TextPart // always math
	AdditionalProposals []TextPart // always math
}

func (ol OrderedListFieldBlock) instantiate(params expression.Variables, ID int) instance {
	out := OrderedListFieldInstance{
		Label:               ol.Label,
		Answer:              make([]string, len(ol.Answer)),
		AdditionalProposals: make([]string, len(ol.AdditionalProposals)),
		ID:                  ID,
	}

	for i, a := range ol.Answer {
		out.Answer[i] = a.instantiate(params).Text
	}

	for i, a := range ol.AdditionalProposals {
		out.AdditionalProposals[i] = a.instantiate(params).Text
	}

	return out
}

func (ol OrderedListFieldBlock) validate(params expression.RandomParameters) error {
	for _, a := range ol.Answer {
		if err := a.validate(); err != nil {
			return err
		}
	}

	for _, a := range ol.AdditionalProposals {
		if err := a.validate(); err != nil {
			return err
		}
	}

	return nil
}

// CoordExpression is a pair of valid expression.Expression
type CoordExpression struct {
	X, Y string
}

func (c CoordExpression) instantiate(params expression.Variables) repere.IntCoord {
	return repere.IntCoord{
		X: int(expression.MustEvaluate(c.X, params)),
		Y: int(expression.MustEvaluate(c.Y, params)),
	}
}

func (c CoordExpression) validate(params expression.RandomParameters) error {
	if err := validateNumberExpression(c.X, params, false); err != nil {
		return err
	}
	if err := validateNumberExpression(c.Y, params, false); err != nil {
		return err
	}
	return nil
}

type FigurePointFieldBlock struct {
	Answer CoordExpression
	Figure FigureBlock
}

func (fp FigurePointFieldBlock) instantiate(params expression.Variables, ID int) instance {
	return FigurePointFieldInstance{
		Figure: fp.Figure.instantiateF(params).Figure,
		Answer: fp.Answer.instantiate(params),
		ID:     ID,
	}
}

func (fp FigurePointFieldBlock) validate(params expression.RandomParameters) error {
	if err := fp.Figure.validate(params); err != nil {
		return err
	}
	if err := fp.Answer.validate(params); err != nil {
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

func (fv FigureVectorFieldBlock) instantiate(params expression.Variables, ID int) instance {
	out := FigureVectorFieldInstance{
		ID:             ID,
		Figure:         fv.Figure.instantiateF(params).Figure,
		Answer:         fv.Answer.instantiate(params),
		MustHaveOrigin: fv.MustHaveOrigin,
	}

	if fv.MustHaveOrigin {
		out.AnswerOrigin = fv.AnswerOrigin.instantiate(params)
	}
	return out
}

func (fp FigureVectorFieldBlock) validate(params expression.RandomParameters) error {
	if err := fp.Figure.validate(params); err != nil {
		return err
	}
	if err := fp.Answer.validate(params); err != nil {
		return err
	}
	if fp.MustHaveOrigin {
		if err := fp.AnswerOrigin.validate(params); err != nil {
			return err
		}
	}

	return nil
}

type VariationTableFieldBlock struct {
	Answer VariationTableBlock
}

func (vt VariationTableFieldBlock) instantiate(params expression.Variables, ID int) instance {
	return VariationTableFieldInstance{
		ID:     ID,
		Answer: vt.Answer.instantiateVT(params),
	}
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

func (fp FunctionPointsFieldBlock) instantiate(params expression.Variables, ID int) instance {
	return FunctionPointsFieldInstance{
		Function: expression.FunctionExpr{
			Function: expression.MustParse(fp.Function),
			Variable: fp.Variable,
		},
		ID:    ID,
		Label: fp.Label,
		XGrid: fp.XGrid,
	}
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
	return fn.validate(params)
}

type FigureVectorPairFieldBlock struct {
	Figure    FigureBlock
	Criterion VectorPairCriterion
}

func (fv FigureVectorPairFieldBlock) instantiate(params expression.Variables, ID int) instance {
	return FigureVectorPairFieldInstance{
		ID:        ID,
		Figure:    fv.Figure.instantiateF(params).Figure,
		Criterion: fv.Criterion,
	}
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

func (fa FigureAffineLineFieldBlock) instantiate(params expression.Variables, ID int) instance {
	return FigureAffineLineFieldInstance{
		ID:     ID,
		Label:  fa.Label,
		Figure: fa.Figure.instantiateF(params).Figure,
		Answer: [2]float64{
			expression.MustEvaluate(fa.A, params),
			expression.MustEvaluate(fa.B, params),
		},
	}
}

func (fa FigureAffineLineFieldBlock) validate(params expression.RandomParameters) error {
	if err := fa.Figure.validate(params); err != nil {
		return err
	}
	if err := validateNumberExpression(fa.A, params, false); err != nil {
		return err
	}
	if err := validateNumberExpression(fa.B, params, false); err != nil {
		return err
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

func (tf TreeFieldBlock) instantiate(params expression.Variables, ID int) instance {
	out := TreeFieldInstance{
		ID:              ID,
		EventsProposals: make([]client.TextOrMath, len(tf.EventsProposals)),
	}
	for i, p := range tf.EventsProposals {
		out.EventsProposals[i] = client.TextOrMath{Text: p}
	}

	var buildTree func(node TreeNodeAnswer) client.TreeNodeAnswer
	buildTree = func(node TreeNodeAnswer) client.TreeNodeAnswer {
		out := client.TreeNodeAnswer{
			Value:         node.Value,
			Probabilities: make([]float64, len(node.Probabilities)),
			Children:      make([]client.TreeNodeAnswer, len(node.Children)),
		}
		for i, c := range node.Probabilities {
			out.Probabilities[i] = expression.MustEvaluate(c, params)
		}
		for i, c := range node.Children {
			out.Children[i] = buildTree(c)
		}
		return out
	}

	out.Answer = client.TreeAnswer{Root: buildTree(tf.AnswerRoot)}
	return out
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

func (tf TableFieldBlock) instantiate(params expression.Variables, ID int) instance {
	out := TableFieldInstance{
		ID:                ID,
		HorizontalHeaders: make([]client.TextOrMath, len(tf.HorizontalHeaders)),
		VerticalHeaders:   make([]client.TextOrMath, len(tf.VerticalHeaders)),
		Answer:            client.TableAnswer{Rows: make([][]float64, len(tf.Answer))},
	}
	for i, cell := range tf.HorizontalHeaders {
		out.HorizontalHeaders[i] = cell.instantiate(params)
	}
	for i, cell := range tf.VerticalHeaders {
		out.VerticalHeaders[i] = cell.instantiate(params)
	}

	for i, row := range tf.Answer {
		rowInstance := make([]float64, len(row))
		for j, v := range row {
			rowInstance[j] = expression.MustEvaluate(v, params)
		}
		out.Answer.Rows[i] = rowInstance
	}
	return out
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
			if err := validateNumberExpression(cell, params, true); err != nil {
				return err
			}
		}
	}

	return nil
}
