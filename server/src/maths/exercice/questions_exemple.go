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

var PredefinedQuestions = [...]QuestionInstance{
	{
		Title: "Calcul littéral", Enonce: EnonceInstance{
			TextBlock{"Développer l’expression :"},
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
			TextBlock{"Écrire sous une seule fraction : "},
			FormulaInstance{IsInline: false, Chunks: []FormulaPartInstance{
				{Expression: mustParse("(1/3)+(2/5)")},
				{StaticContent: `= \frac{a}{b}`},
			}},
			TextBlock{"avec "},
			FormulaInstance{IsInline: true, Chunks: []FormulaPartInstance{{StaticContent: "a = "}}},
			NumberFieldInstance{ID: 0, Answer: 1*5 + 2*3},
			TextBlock{" et "},
			FormulaInstance{IsInline: true, Chunks: []FormulaPartInstance{{StaticContent: "b = "}}},
			NumberFieldInstance{ID: 1, Answer: 3 * 5},
		},
	},
	{
		Title: "Calcul littéral", Enonce: EnonceInstance{
			TextBlock{"Écrire sous la forme "},
			FormulaInstance{IsInline: true, Chunks: []FormulaPartInstance{
				{StaticContent: `a\sqrt{b},`},
			}},
			TextBlock{" où "},
			FormulaInstance{IsInline: true, Chunks: []FormulaPartInstance{
				{StaticContent: `a`},
			}},
			TextBlock{" et "},
			FormulaInstance{IsInline: true, Chunks: []FormulaPartInstance{
				{StaticContent: `b`},
			}},
			TextBlock{" sont des entiers les plus petits possibles :"},
			FormulaInstance{IsInline: false, Chunks: []FormulaPartInstance{
				{Expression: mustParse("sqrt(50)")},
				{StaticContent: " = "},
				{Expression: mustParse("a * sqrt(b)")},
			}},
			TextBlock{"avec "},
			FormulaInstance{IsInline: true, Chunks: []FormulaPartInstance{{StaticContent: "a = "}}},
			NumberFieldInstance{ID: 0, Answer: 5},
			TextBlock{" et "},
			FormulaInstance{IsInline: true, Chunks: []FormulaPartInstance{{StaticContent: "b = "}}},
			NumberFieldInstance{ID: 1, Answer: 2},
		},
	},
	{
		Title: "Nombres entiers", Enonce: EnonceInstance{
			FormulaInstance{IsInline: true, Chunks: []FormulaPartInstance{
				{Expression: mustParse("89")},
			}},
			TextBlock{" est-il un nombre premier ?"},
			RadioFieldInstance{
				ID:     0,
				Answer: 1,
				Proposals: []client.ListFieldProposal{
					{Content: []client.ListFieldProposalPart{{Content: "Oui"}}},
					{Content: []client.ListFieldProposalPart{{Content: "Non"}}},
				},
			},
		},
	},
	{
		Title: "Très longue question horizontale", Enonce: EnonceInstance{
			TextBlock{"Écrire sous une seule fraction : "},
			FormulaInstance{IsInline: true, Chunks: []FormulaPartInstance{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextBlock{"Écrire sous une seule fraction : "},
			FormulaInstance{IsInline: true, Chunks: []FormulaPartInstance{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextBlock{"Écrire sous une seule fraction : "},
			FormulaInstance{IsInline: true, Chunks: []FormulaPartInstance{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextBlock{"Écrire sous une seule fraction : "},
			FormulaInstance{IsInline: true, Chunks: []FormulaPartInstance{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextBlock{"Écrire sous une seule fraction : "},
			FormulaInstance{IsInline: true, Chunks: []FormulaPartInstance{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextBlock{"Écrire sous une seule fraction : "},
			FormulaInstance{IsInline: true, Chunks: []FormulaPartInstance{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextBlock{"Écrire sous une seule fraction : "},
			FormulaInstance{IsInline: true, Chunks: []FormulaPartInstance{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
		},
	},
	{
		Title: "Très longue question verticale", Enonce: EnonceInstance{
			TextBlock{"Écrire sous une seule fraction : "},
			FormulaInstance{IsInline: false, Chunks: []FormulaPartInstance{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextBlock{"Écrire sous une seule fraction : "},
			FormulaInstance{IsInline: false, Chunks: []FormulaPartInstance{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextBlock{"Écrire sous une seule fraction : "},
			FormulaInstance{IsInline: false, Chunks: []FormulaPartInstance{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextBlock{"Écrire sous une seule fraction : "},
			FormulaInstance{IsInline: false, Chunks: []FormulaPartInstance{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextBlock{"Écrire sous une seule fraction : "},
			FormulaInstance{IsInline: false, Chunks: []FormulaPartInstance{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextBlock{"Écrire sous une seule fraction : "},
			FormulaInstance{IsInline: false, Chunks: []FormulaPartInstance{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextBlock{"Écrire sous une seule fraction : "},
			FormulaInstance{IsInline: false, Chunks: []FormulaPartInstance{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
		},
	},
}
