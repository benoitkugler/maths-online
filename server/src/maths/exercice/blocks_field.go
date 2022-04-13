package exercice

import (
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

type FigureVectorFieldBlock struct {
	Answer CoordExpression

	AnswerOrigin CoordExpression // optionnal

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

type VariationTableFieldBlock struct {
	Answer VariationTableBlock
}

func (vt VariationTableFieldBlock) instantiate(params expression.Variables, ID int) instance {
	return VariationTableFieldInstance{
		ID:     ID,
		Answer: vt.Answer.instantiateVT(params),
	}
}

type FunctionPointsFieldBlock struct {
	Function string // valid expression.Expression
	Label    string
	Variable expression.Variable
	XGrid    []int
}

func (fp FunctionPointsFieldBlock) instantiate(params expression.Variables, ID int) instance {
	return FunctionPointsFieldInstance{
		ID:       ID,
		Function: expression.MustParse(fp.Function),
		Label:    fp.Label,
		XGrid:    fp.XGrid,
		Variable: fp.Variable,
	}
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
