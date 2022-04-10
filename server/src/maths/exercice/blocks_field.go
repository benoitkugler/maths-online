package exercice

import (
	"github.com/benoitkugler/maths-online/maths/exercice/client"
	"github.com/benoitkugler/maths-online/maths/expression"
)

var (
	_ Block = NumberFieldBlock{}
	_ Block = FormulaFieldBlock{}
	_ Block = RadioFieldBlock{}
	_ Block = OrderedListFieldBlock{}
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
