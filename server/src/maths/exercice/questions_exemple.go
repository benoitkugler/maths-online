package exercice

import (
	"github.com/benoitkugler/maths-online/maths/exercice/client"
	"github.com/benoitkugler/maths-online/maths/expression"
)

func mustParse(s string) *expression.Expression {
	e, _, err := expression.Parse(s)
	if err != nil {
		panic(err)
	}
	return e
}

func text(s string) TextOrMaths {
	return TextOrMaths{MathOrExpression: MathOrExpression{StaticContent: s}}
}

func staticMath(s string) TextOrMaths {
	return TextOrMaths{MathOrExpression: MathOrExpression{StaticContent: s}, IsMath: true}
}

func expr(s string) TextOrMaths {
	return TextOrMaths{MathOrExpression: MathOrExpression{Expression: mustParse(s)}, IsMath: true}
}

var PredefinedQuestions = [...]QuestionInstance{
	{
		Title: "Calcul littéral", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{text("Développer l’expression :")}},
			ExpressionFieldInstance{
				ID:              0,
				Label:           mustParse("(x−6)*(4*x−3)"),
				ComparisonLevel: expression.SimpleSubstitutions,
				Answer:          mustParse("24*x^2 - 27 *x + 18"),
			},
		},
	},
	{
		Title: "Calcul littéral", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{text("Écrire sous une seule fraction : ")}},
			FormulaDisplayInstance{Parts: []MathOrExpression{
				{Expression: mustParse("(1/3)+(2/5)")},
				{StaticContent: `= \frac{a}{b}`},
			}},
			TextInstance{Parts: []TextOrMaths{
				text("avec "),
				staticMath("a = "),
			}},
			NumberFieldInstance{ID: 0, Answer: 1*5 + 2*3},
			TextInstance{Parts: []TextOrMaths{
				text(" et "),
				staticMath("b = "),
			}},
			NumberFieldInstance{ID: 1, Answer: 3 * 5},
		},
	},
	{
		Title: "Calcul littéral", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Écrire sous la forme "),
				staticMath(`a\sqrt{b},`),
				text(" où "),
				staticMath("a"),
				text(" et "),
				staticMath("b"),
				text(" sont des entiers les plus petits possibles :"),
			}},
			FormulaDisplayInstance{Parts: []MathOrExpression{
				{Expression: mustParse("sqrt(50)")},
				{StaticContent: " = "},
				{Expression: mustParse("a * sqrt(b)")},
			}},
			TextInstance{Parts: []TextOrMaths{
				text("avec "),
				staticMath("a = "),
			}},
			NumberFieldInstance{ID: 0, Answer: 5},
			TextInstance{Parts: []TextOrMaths{
				text(" et "),
				staticMath("b = "),
			}},
			NumberFieldInstance{ID: 1, Answer: 2},
		},
	},
	{
		Title: "Nombres entiers", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				expr("89"),
				text(" est-il un nombre premier ?"),
			}},
			RadioFieldInstance{
				ID:     0,
				Answer: 1,
				Proposals: []client.ListFieldProposal{
					{Content: []client.TextOrMath{{Text: "Oui"}}},
					{Content: []client.TextOrMath{{Text: "Non"}}},
				},
			},
		},
	},
	{
		Title: "Nombres entiers", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				expr("987"),
				text(" est-il un mulitple de "),
				expr("3"),
			}},
			RadioFieldInstance{
				ID:     0,
				Answer: 0,
				Proposals: []client.ListFieldProposal{
					{Content: []client.TextOrMath{{Text: "Oui"}}},
					{Content: []client.TextOrMath{{Text: "Non"}}},
				},
			},
		},
	},
	{
		Title: "Nombres réels", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("A quel plus petit ensemble appartient "),
				expr("4/7"),
			}},
			RadioFieldInstance{
				ID:     0,
				Answer: 4,
				Proposals: []client.ListFieldProposal{
					{Content: []client.TextOrMath{{Text: `\mathbb{N}`, IsMath: true}}},
					{Content: []client.TextOrMath{{Text: `\mathbb{Z}`, IsMath: true}}},
					{Content: []client.TextOrMath{{Text: `\mathbb{D}`, IsMath: true}}},
					{Content: []client.TextOrMath{{Text: `\mathbb{Q}`, IsMath: true}}},
					{Content: []client.TextOrMath{{Text: `\mathbb{R}`, IsMath: true}}},
				},
			},
		},
	},
	{
		Title: "Très longue question horizontale", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Écrire sous une seule fraction : "),
				expr("(1/3)+(2/5)"),
				text("Écrire sous une seule fraction : "),
				expr("(1/3)+(2/5)"),
				text("Écrire sous une seule fraction : "),
				expr("(1/3)+(2/5)"),
				text("Écrire sous une seule fraction : "),
				expr("(1/3)+(2/5)"),
				text("Écrire sous une seule fraction : "),
				expr("(1/3)+(2/5)"),
				text("Écrire sous une seule fraction : "),
				expr("(1/3)+(2/5)"),
				text("Écrire sous une seule fraction : "),
				expr("(1/3)+(2/5)"),
			}},
		},
	},
	{
		Title: "Très longue question verticale", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{text("Écrire sous une seule fraction : ")}},
			FormulaDisplayInstance{Parts: []MathOrExpression{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextInstance{Parts: []TextOrMaths{text("Écrire sous une seule fraction : ")}},
			FormulaDisplayInstance{Parts: []MathOrExpression{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextInstance{Parts: []TextOrMaths{text("Écrire sous une seule fraction : ")}},
			FormulaDisplayInstance{Parts: []MathOrExpression{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextInstance{Parts: []TextOrMaths{text("Écrire sous une seule fraction : ")}},
			FormulaDisplayInstance{Parts: []MathOrExpression{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextInstance{Parts: []TextOrMaths{text("Écrire sous une seule fraction : ")}},
			FormulaDisplayInstance{Parts: []MathOrExpression{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextInstance{Parts: []TextOrMaths{text("Écrire sous une seule fraction : ")}},
			FormulaDisplayInstance{Parts: []MathOrExpression{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
		},
	},
}
