package exercice

import (
	"math"

	"github.com/benoitkugler/maths-online/maths/exercice/client"
	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/repere"
)

func mustParse(s string) *expression.Expression {
	e, _, err := expression.Parse(s)
	if err != nil {
		panic(err)
	}
	return e
}

func text(s string) TextOrMaths {
	return TextOrMaths{StringOrExpression: StringOrExpression{String: s}}
}

func staticMath(s string) TextOrMaths {
	return TextOrMaths{StringOrExpression: StringOrExpression{String: s}, IsMath: true}
}

func expr(s string) TextOrMaths {
	return TextOrMaths{StringOrExpression: StringOrExpression{Expression: mustParse(s)}, IsMath: true}
}

var (
	_A = repere.Coord{X: 5, Y: 25}
	_B = repere.Coord{X: 5, Y: 5}
	_C = repere.Coord{X: 45, Y: 5}
	_D = repere.Coord{X: 45, Y: 25}
	_K = repere.OrthogonalProjection(_D, _A, _C)

	figure1 = repere.Figure{
		Width:  50,
		Height: 30,
		Points: map[string]repere.LabeledPoint{
			"A": {Point: _A, Pos: repere.TopLeft},
			"B": {Point: _B, Pos: repere.BottomLeft},
			"C": {Point: _C, Pos: repere.BottomRight},
			"D": {Point: _D, Pos: repere.TopRight},
			"H": {Point: repere.OrthogonalProjection(_B, _A, _C), Pos: repere.Top},
			"K": {Point: _K, Pos: repere.Top},
		},
		Lines: []repere.Line{
			{LabelName: "a", From: "A", To: "B", LabelPos: repere.Left},
			{LabelName: "b", From: "B", To: "C", LabelPos: repere.Bottom},
			{LabelName: "c", From: "C", To: "D", LabelPos: repere.Right},
			{LabelName: "d", From: "D", To: "A", LabelPos: repere.Top},

			// diagonal
			{LabelName: "", From: "A", To: "C", LabelPos: repere.Bottom},

			{LabelName: "", From: "B", To: "H", LabelPos: repere.Bottom},
			{LabelName: "", From: "D", To: "K", LabelPos: repere.Top},
		},
	}
)

var (
	__O = repere.Coord{X: 6, Y: 6}
	__A = repere.Coord{X: 5, Y: 5}
	__D = repere.Coord{X: 5, Y: 7}
	__B = repere.Coord{X: 7, Y: 7}
	__J = repere.Coord{X: 3, Y: 6}
	__C = repere.Coord{X: 1, Y: 10}

	figure2 = repere.Figure{
		Width:  12,
		Height: 12,
		Points: map[string]repere.LabeledPoint{
			"A": {Point: __A, Pos: repere.TopLeft},
			"B": {Point: __B, Pos: repere.BottomLeft},
			"C": {Point: __C, Pos: repere.BottomRight},
			"D": {Point: __D, Pos: repere.TopRight},
			"J": {Point: __J, Pos: repere.TopRight},
			"O": {Point: __O, Pos: repere.TopRight},
		},
		Lines: []repere.Line{},
	}
)

var PredefinedQuestions = [...]QuestionInstance{
	{
		Title: "Calcul littéral", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{text("Développer l’expression :")}},
			ExpressionFieldInstance{
				ID:              0,
				Label:           mustParse("(x−6)*(4*x−3)"),
				ComparisonLevel: expression.SimpleSubstitutions,
				Answer:          mustParse("4*x^2 - 27 *x + 18"),
			},
		},
	},
	{
		Title: "Calcul littéral", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{text("Écrire sous une seule fraction : ")}},
			FormulaDisplayInstance{Parts: []StringOrExpression{
				{Expression: mustParse("(1/3)+(2/5)")},
				{String: `= \frac{a}{b}`},
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
			FormulaDisplayInstance{Parts: []StringOrExpression{
				{Expression: mustParse("sqrt(50)")},
				{String: " = "},
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
				Answer: 0,
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
				text("Traduire en intervalle : "),
				staticMath(`x \ge `),
				expr("12"),
			}},
			OrderedListFieldInstance{
				ID:    0,
				Label: `x \in`,
				Answer: []StringOrExpression{ // [12;+infty[
					{String: "["},
					{Expression: mustParse("12")},
					{String: ";"},
					{String: "+"},
					{String: `\infty`},
					{String: `[`},
				},
				AdditionalProposals: []StringOrExpression{
					{String: "]"}, // some duplicates
					{String: `\infty`},
					{Expression: mustParse("11")},
					{String: "-"},
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
		Title: "Nombres réels", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Déterminer "),
				staticMath(`I \cap J`),
				text(" où "),
				staticMath(`I = [`),
				expr("3"),
				staticMath(`;`),
				expr("9"),
				staticMath(`]`),
				text(" et "),
				staticMath(`J = [`),
				expr("0"),
				staticMath(`;`),
				expr("6"),
				staticMath(`]`),
			}},
			OrderedListFieldInstance{
				ID:    0,
				Label: `I \cap J =`,
				Answer: []StringOrExpression{ // [12;+infty[
					{String: "["},
					{Expression: mustParse("3")},
					{String: ";"},
					{Expression: mustParse("6")},
					{String: `]`},
				},
				AdditionalProposals: []StringOrExpression{
					{String: "]"}, // some duplicates
					{String: `\infty`},
					{Expression: mustParse("0")},
					{Expression: mustParse("9")},
					{String: "-"},
				},
			},
		},
	},
	{
		Title: "Equations et inéquations", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Résoudre dans "),
				staticMath(`\mathbb{R}`),
				text("l'équation :"),
			}},
			FormulaDisplayInstance{Parts: []StringOrExpression{
				{Expression: mustParse("20*x - 7")},
				{String: "="},
				{Expression: mustParse("34")},
			}},
			ExpressionFieldInstance{
				ID:              0,
				Label:           mustParse("x"),
				ComparisonLevel: expression.SimpleSubstitutions,
				Answer:          mustParse("(34+7)/20"),
			},
		},
	},
	{
		Title: "Equations et inéquations", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Résoudre dans "),
				staticMath(`\mathbb{R}`),
				text("l'inéquation :"),
			}},
			FormulaDisplayInstance{Parts: []StringOrExpression{
				{Expression: mustParse("-8*x +1")},
				{String: "<"},
				{Expression: mustParse("2*x - 5")},
			}},
			TextInstance{Parts: []TextOrMaths{
				staticMath("x > "),
			}},
			NumberFieldInstance{ID: 0, Answer: 0.6},
		},
	},
	{
		Title: "Equations et inéquations", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Résoudre dans "),
				staticMath(`\mathbb{R}`),
				text("l'équation :"),
			}},
			FormulaDisplayInstance{Parts: []StringOrExpression{
				{Expression: mustParse("(x-7)*(4*x - 9)")},
				{String: "="},
				{Expression: mustParse("0")},
			}},
			TextInstance{Parts: []TextOrMaths{
				text("Solutions : "),
				staticMath("x = "),
			}},
			NumberFieldInstance{ID: 0, Answer: 7},
			TextInstance{Parts: []TextOrMaths{
				text(" ou "),
				staticMath("x = "),
			}},
			NumberFieldInstance{ID: 1, Answer: 9. / 4},
		},
	},
	{
		Title: "Equations et inéquations", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Combien l’équation "),
				staticMath(`f(x) = `),
				expr("-1"),
				text(" admet-elle de solutions ? "),
			}},
			VariationTableInstance{
				Xs:  []expression.Number{-20, -10, 0, 3, 18},
				Fxs: []expression.Number{-6, -2, -8, 0, -5},
			},
			TextInstance{Parts: []TextOrMaths{
				text("Nombre de solutions : "),
			}},
			NumberFieldInstance{ID: 0, Answer: 2},
		},
	},
	{
		Title: "Equations et inéquations", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Résoudre l’inéquation "),
				staticMath(`f(x) \ge 0`),
			}},
			SignTableInstance{Columns: []client.SignColumn{
				{X: `-\infty`, IsYForbiddenValue: false, IsSign: false, IsPositive: false},
				{X: "", IsYForbiddenValue: false, IsSign: true, IsPositive: false},
				{X: "-2", IsYForbiddenValue: false, IsSign: false, IsPositive: true},
				{X: "", IsYForbiddenValue: false, IsSign: true, IsPositive: true},
				{X: "0", IsYForbiddenValue: false, IsSign: false, IsPositive: true},
				{X: "", IsYForbiddenValue: false, IsSign: true, IsPositive: false},
				{X: "4", IsYForbiddenValue: true, IsSign: false, IsPositive: true},
				{X: "", IsYForbiddenValue: false, IsSign: true, IsPositive: true},
				{X: `+\infty`, IsYForbiddenValue: false, IsSign: false, IsPositive: false},
			}},
			TextInstance{Parts: []TextOrMaths{
				text("Intervalle(s) solution(s) :"),
			}},
			OrderedListFieldInstance{
				ID: 0,
				Answer: []StringOrExpression{ // [12;+infty[
					{String: "["},
					{Expression: mustParse("-2")},
					{String: ";"},
					{Expression: mustParse("0")},
					{String: `]`},
					{String: `\cup`},
					{String: "]"},
					{Expression: mustParse("4")},
					{String: ";"},
					{String: `+\infty`},
					{String: `[`},
				},
				AdditionalProposals: []StringOrExpression{
					{String: "]"}, // some duplicates
					{String: `-\infty`},
					{Expression: mustParse("0")},
					{String: "-"},
				},
			},
		},
	},
	// geometrie

	{
		Title: "Géométrie plane", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("ABC est un triangle rectangle en C avec "),
				staticMath(`AB = `),
				expr("sqrt(98)"),
				text(" et "),
				staticMath("BC = "),
				expr("7"),
				text(". Calculer, en degrés, "),
				staticMath(`\widehat{ABC}.`),
			}},
			TextInstance{Parts: []TextOrMaths{
				staticMath(`\widehat{ABC} = `),
			}},
			NumberFieldInstance{ID: 0, Answer: mustParse("acos(7/sqrt(98))").Evaluate(nil) * 180 / math.Pi},
		},
	},
	{
		Title: "Géométrie plane", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("ABC est un triangle rectangle en C avec "),
				staticMath(`AB = `),
				expr("sqrt(98)"),
				text(" et "),
				staticMath("BC = "),
				expr("7"),
				text(". Calculer AC."),
			}},
			TextInstance{Parts: []TextOrMaths{
				text("AC = "),
			}},
			NumberFieldInstance{ID: 0, Answer: mustParse("sqrt(sqrt(98)^2 - 7^2)").Evaluate(nil)},
		},
	},
	{
		Title: "Géométrie plane", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Soient "),
				staticMath(`AB = `),
				expr("4"),
				staticMath(" ; AC = "),
				expr("12"),
				text(" et "),
				staticMath("BC = "),
				expr("8"),
				text(". Le triangle ABC est-il rectangle ? Si oui, en quoi ?"),
			}},
			RadioFieldInstance{
				ID: 0,
				Answer: int(mustParse("1 * isZero(a^2 - b^2 - c^2) + 2*isZero(b^2 - a^2 - c^2) + 3*isZero(c^2 - a^2 - b^2)").Evaluate(expression.Variables{
					'a': 8,  // BC
					'b': 12, // AC
					'c': 4,  // AB
				})),
				Proposals: []client.ListFieldProposal{
					{Content: []client.TextOrMath{{Text: `ABC n'est pas rectangle.`, IsMath: false}}},
					{Content: []client.TextOrMath{{Text: `ABC est rectangle en A.`, IsMath: false}}},
					{Content: []client.TextOrMath{{Text: `ABC est rectangle en B.`, IsMath: false}}},
					{Content: []client.TextOrMath{{Text: `ABC est rectangle en C.`, IsMath: false}}},
				},
			},
		},
	},

	{
		Title: "Repérage dans le plan", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Soient les points "),
				staticMath(`A(`),
				expr("8"),
				staticMath(";"),
				expr("19"),
				staticMath(")"),
				text(" et "),
				staticMath(`B(`),
				expr("-6"),
				staticMath(";"),
				expr("0"),
				staticMath(")."),
				text("Quelles sont les coordonnées de M, milieu de [AB] ?"),
			}},
			TextInstance{Parts: []TextOrMaths{
				staticMath("M = ("),
			}},
			NumberFieldInstance{ID: 0, Answer: mustParse("(8 + (-6))/2").Evaluate(nil)},
			TextInstance{Parts: []TextOrMaths{
				staticMath(";"),
			}},
			NumberFieldInstance{ID: 1, Answer: mustParse("(19 + 0)/2").Evaluate(nil)},
			TextInstance{Parts: []TextOrMaths{
				staticMath(")"),
			}},
		},
	},

	{
		Title: "Repérage dans le plan", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Quel point est le projeté orthogonal de D sur (AH) ?"),
			}},
			FigureInstance{Figure: figure1},
			RadioFieldInstance{
				ID:     0,
				Answer: 2,
				Proposals: []client.ListFieldProposal{
					{Content: []client.TextOrMath{{Text: `A`, IsMath: false}}},
					{Content: []client.TextOrMath{{Text: `B`, IsMath: false}}},
					{Content: []client.TextOrMath{{Text: `K`, IsMath: false}}},
					{Content: []client.TextOrMath{{Text: `H`, IsMath: false}}},
					{Content: []client.TextOrMath{{Text: `D`, IsMath: false}}},
				},
			},
		},
	},
	{
		Title: "Repérage dans le plan", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Construire le point M, projeté orthogonal de K sur (BC)."),
			}},
			FigurePointFieldInstance{
				Figure: figure1,
				Answer: repere.OrthogonalProjection(_K, _B, _C).Round(),
				ID:     0,
			},
		},
	},
	{
		Title: "Repérage dans le plan", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Construire le vecteur "),
				staticMath(`\overrightarrow{AB} + \overrightarrow{CD}.`),
			}},
			FigureVectorFieldInstance{
				Figure: figure2,
				Answer: repere.IntCoord{X: 6, Y: -1},
				ID:     0,
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
			FormulaDisplayInstance{Parts: []StringOrExpression{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextInstance{Parts: []TextOrMaths{text("Écrire sous une seule fraction : ")}},
			FormulaDisplayInstance{Parts: []StringOrExpression{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextInstance{Parts: []TextOrMaths{text("Écrire sous une seule fraction : ")}},
			FormulaDisplayInstance{Parts: []StringOrExpression{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextInstance{Parts: []TextOrMaths{text("Écrire sous une seule fraction : ")}},
			FormulaDisplayInstance{Parts: []StringOrExpression{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextInstance{Parts: []TextOrMaths{text("Écrire sous une seule fraction : ")}},
			FormulaDisplayInstance{Parts: []StringOrExpression{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
			TextInstance{Parts: []TextOrMaths{text("Écrire sous une seule fraction : ")}},
			FormulaDisplayInstance{Parts: []StringOrExpression{
				{Expression: mustParse("(1/3)+(2/5)")},
			}},
		},
	},
}
