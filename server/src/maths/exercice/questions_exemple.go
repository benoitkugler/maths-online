package exercice

import "github.com/benoitkugler/maths-online/maths/expression"

func mustParse(s string) *expression.Expression {
	e, _, err := expression.Parse(s)
	if err != nil {
		panic(err)
	}
	return e
}

var questions = [...]QuestionInstance{
	{
		Title: "Calcul littéral", Content: ContentInstance{
			TextBlock{"Développer l’expression : "},
			FormulaInstance{IsInline: true, Chunks: []FormulaPartInstance{
				{Expression: mustParse("(x−6)*(4*x−3)")},
			}},
		},
	},
}
