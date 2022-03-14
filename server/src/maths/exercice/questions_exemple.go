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
		Title: "Calcul littéral", Enonce: EnonceInstance{
			TextBlock{"Développer l’expression : "},
			FormulaInstance{IsInline: true, Chunks: []FormulaPartInstance{
				{Expression: mustParse("(x−6)*(4*x−3)")},
			}},
		},
	},
	{
		Title: "Calcul littéral", Enonce: EnonceInstance{
			TextBlock{"Écrire sous une seule fraction : "},
			FormulaInstance{IsInline: true, Chunks: []FormulaPartInstance{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
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
