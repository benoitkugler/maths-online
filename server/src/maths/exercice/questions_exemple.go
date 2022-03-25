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

func mustEvaluate(s string, vars ...expression.Variables) float64 {
	e := mustParse(s)

	var resolver expression.ValueResolver
	if len(vars) > 0 {
		resolver = vars[0]
	}

	out, err := e.Evaluate(resolver)
	if err != nil {
		panic(err)
	}
	return out
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
		Segments: []repere.Segment{
			{LabelName: "", From: "A", To: "B", LabelPos: repere.Left},
			{LabelName: "", From: "B", To: "C", LabelPos: repere.Bottom},
			{LabelName: "", From: "C", To: "D", LabelPos: repere.Right},
			{LabelName: "", From: "D", To: "A", LabelPos: repere.Top},

			// diagonal
			{LabelName: "", From: "A", To: "C", LabelPos: repere.Bottom},

			{LabelName: "", From: "B", To: "H", LabelPos: repere.Bottom},
			{LabelName: "", From: "D", To: "K", LabelPos: repere.Top},
		},
	}
)

var (
	__A = repere.Coord{X: -1, Y: -1}
	__D = repere.Coord{X: -1, Y: 1}
	__B = repere.Coord{X: 1, Y: 1}
	__J = repere.Coord{X: -3, Y: 0}
	__C = repere.Coord{X: -5, Y: 4}

	figure2 = repere.Figure{
		Origin: repere.Coord{X: 6, Y: 6},
		Width:  12,
		Height: 12,
		Points: map[string]repere.LabeledPoint{
			"A": {Point: __A, Pos: repere.TopLeft},
			"B": {Point: __B, Pos: repere.BottomLeft},
			"C": {Point: __C, Pos: repere.BottomRight},
			"D": {Point: __D, Pos: repere.TopRight},
			"J": {Point: __J, Pos: repere.TopRight},
		},
		Segments: []repere.Segment{},
		ShowGrid: true,
	}

	figure3 = repere.Figure{
		Origin: repere.Coord{X: 6, Y: 6},
		Width:  8,
		Height: 8,
		Points: map[string]repere.LabeledPoint{
			"A": {Point: __A, Pos: repere.TopLeft},
			"B": {Point: __B, Pos: repere.BottomLeft},
			"D": {Point: __D, Pos: repere.TopRight},
		},
		Segments: []repere.Segment{},
		ShowGrid: true,
	}

	_A4     = repere.Coord{X: 4, Y: 9}
	_B4     = repere.Coord{X: 6, Y: 6}
	_C4     = repere.Coord{X: 5, Y: 12}
	_D4     = repere.Coord{X: 2, Y: 0}
	figure4 = repere.Figure{
		Width:  7,
		Height: 13,
		Points: map[string]repere.LabeledPoint{
			"A": {Point: _A4, Pos: repere.TopLeft},
			"B": {Point: _B4, Pos: repere.BottomLeft},
			"C": {Point: _C4, Pos: repere.TopRight},
			"D": {Point: _D4, Pos: repere.TopRight},
		},
		Segments: []repere.Segment{},
	}

	_A5     = repere.Coord{X: 1, Y: 1}
	_B5     = repere.Coord{X: -1, Y: 5}
	_J5     = repere.Coord{X: -1, Y: 3}
	_H5     = repere.Coord{X: -2, Y: 4}
	_F5     = repere.Coord{X: -3, Y: 4}
	_G5     = repere.Coord{X: -2, Y: 2}
	figure5 = repere.Figure{
		Origin:   repere.Coord{X: 4, Y: 0},
		ShowGrid: true,
		Width:    7,
		Height:   8,
		Points: map[string]repere.LabeledPoint{
			"A": {Point: _A5, Pos: repere.BottomLeft},
			"B": {Point: _B5, Pos: repere.BottomLeft},
			"J": {Point: _J5, Pos: repere.TopRight},
			"H": {Point: _H5, Pos: repere.TopRight},
			"F": {Point: _F5, Pos: repere.TopRight},
			"G": {Point: _G5, Pos: repere.TopRight},
		},
		Segments: []repere.Segment{
			{From: "A", To: "B", AsVector: true},
			{From: "F", To: "G", AsVector: true},
			{From: "J", To: "H", AsVector: true},
		},
	}

	_line   = repere.Line{A: 3. / 2, B: 1, Label: "(d)"}
	figure6 = repere.Figure{
		Origin:   repere.Coord{X: 2, Y: 1},
		Width:    6,
		Height:   8,
		ShowGrid: true,
		Lines: []repere.Line{
			_line,
		},
	}
)

// pythagorian triplet
var (
	pythagorians = expression.PythagorianTriplet{
		A: 'a', B: 'b', C: 'c',
		SeedStart: 2, SeedEnd: 20,
	}
	distanceParams = expression.RandomParameters{
		'X': mustParse("-randInt(100;200)"),
		'Y': mustParse("300"),
		'A': mustParse("a + X"),
		'B': mustParse("b + Y"),
	}
	distanceSample expression.Variables
)

func init() {
	pythagorians.MergeTo(distanceParams)

	var err error
	distanceSample, err = distanceParams.Instantiate()
	if err != nil {
		panic(err)
	}

	PredefinedQuestions = append(PredefinedQuestions, QuestionInstance{
		Title: "Repérage dans le plan", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Soient les points F("),
				staticMath(expression.Number(distanceSample['A']).String()),
				text(";"),
				staticMath(expression.Number(distanceSample['B']).String()),
				text(") et G("),
				staticMath(expression.Number(distanceSample['X']).String()),
				text(";"),
				staticMath(expression.Number(distanceSample['Y']).String()),
				text("). Calculer FG."),
			}},
			// TextInstance{
			// 	IsHint: true,
			// 	Parts: []TextOrMaths{
			// 		text("(On utilisera sqrt(10) pour "),
			// 		staticMath(`\sqrt{10}`),
			// 		text(")."),
			// 	},
			// },
			// ExpressionFieldInstance{
			// 	ID:              0,
			// 	Label:           StringOrExpression{String: "FG = "},
			// 	Answer:          mustParse("sqrt(1262900)"),
			// 	ComparisonLevel: expression.SimpleSubstitutions,
			// },
			TextInstance{Parts: []TextOrMaths{
				staticMath("FG = "),
			}},
			NumberFieldInstance{ID: 0, Answer: distanceSample['c']},
		},
	})
}

var PredefinedQuestions = []QuestionInstance{
	{
		Title: "Calcul littéral", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{text("Développer l’expression :")}},
			ExpressionFieldInstance{
				ID:              0,
				Label:           StringOrExpression{Expression: mustParse("(x−6)*(4*x−3)")},
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
				Label:           StringOrExpression{String: "x ="},
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
				staticMath("x "),
			}},
			DropDownFieldInstance{
				ID:     0,
				Answer: 1,
				Proposals: []client.ListFieldProposal{
					{Content: []client.TextOrMath{{Text: `<`, IsMath: true}}},
					{Content: []client.TextOrMath{{Text: `>`, IsMath: true}}},
					{Content: []client.TextOrMath{{Text: `=`, IsMath: true}}},
				},
			},
			NumberFieldInstance{ID: 1, Answer: 0.6},
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
			NumberFieldInstance{ID: 0, Answer: mustEvaluate("acos(7/sqrt(98))") * 180 / math.Pi},
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
			NumberFieldInstance{ID: 0, Answer: mustEvaluate("sqrt(sqrt(98)^2 - 7^2)")},
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
				Answer: int(mustEvaluate("1 * isZero(a^2 - b^2 - c^2) + 2*isZero(b^2 - a^2 - c^2) + 3*isZero(c^2 - a^2 - b^2)", expression.Variables{
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
			NumberFieldInstance{ID: 0, Answer: mustEvaluate("(8 + (-6))/2")},
			TextInstance{Parts: []TextOrMaths{
				staticMath(";"),
			}},
			NumberFieldInstance{ID: 1, Answer: mustEvaluate("(19 + 0)/2")},
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
				text("Lire les coordonnées de B."),
			}},
			FigureInstance{
				Figure: figure3,
			},
			TextInstance{Parts: []TextOrMaths{
				staticMath("B = ("),
			}},
			NumberFieldInstance{
				ID:     0,
				Answer: __B.X,
			},
			TextInstance{Parts: []TextOrMaths{
				staticMath("; "),
			}},
			NumberFieldInstance{
				ID:     1,
				Answer: __B.Y,
			},
			TextInstance{Parts: []TextOrMaths{
				staticMath(")"),
			}},
		},
	},
	{
		Title: "Repérage dans le plan", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Quelle est la nature de ABCD ?"),
			}},
			FigureInstance{
				Figure: figure4,
			},
			RadioFieldInstance{
				ID:     0,
				Answer: 0,
				Proposals: []client.ListFieldProposal{
					{Content: []client.TextOrMath{{Text: "Quadrilatère quelconque"}}},
					{Content: []client.TextOrMath{{Text: "Rectangle"}}},
					{Content: []client.TextOrMath{{Text: "Losange"}}},
					{Content: []client.TextOrMath{{Text: "Carré"}}},
				},
			},
		},
	},
	{
		Title: "Vecteurs", Enonce: EnonceInstance{
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
		Title: "Vecteurs", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Simplifier au maximum l'expression vectorielle"),
				staticMath(`\overrightarrow{AB} + \overrightarrow{EB} -  \overrightarrow{EG}`),
			}},
			OrderedListFieldInstance{
				Label: `\overrightarrow{AB} + \overrightarrow{EB} - \overrightarrow{EG} = `,
				Answer: []StringOrExpression{
					{String: "G"},
					{String: "A"},
				},
				AdditionalProposals: []StringOrExpression{
					{String: "E"},
					{String: "B"},
					{String: "F"},
					{String: "A"},
				},
			},
		},
	},
	{
		Title: "Vecteurs", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Construire deux vecteurs égaux."),
			}},
			FigureVectorPairFieldInstance{
				Figure:    figure1,
				Criterion: VectorEquals,
				ID:        0,
			},
		},
	},
	{
		Title: "Vecteurs", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Construire deux vecteurs colinéaires."),
			}},
			FigureVectorPairFieldInstance{
				Figure:    figure1,
				Criterion: VectorColinear,
				ID:        0,
			},
		},
	},
	{
		Title: "Vecteurs", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Lire les coordonnées de "),
				staticMath(`\overrightarrow{FG}.`),
			}},
			FigureInstance{
				Figure: figure5,
			},
			TextInstance{Parts: []TextOrMaths{
				staticMath(`\overrightarrow{FG} = (`),
			}},
			NumberFieldInstance{
				ID:     0,
				Answer: _G5.X - _F5.X,
			},
			TextInstance{Parts: []TextOrMaths{
				staticMath(`;`),
			}},
			NumberFieldInstance{
				ID:     1,
				Answer: _G5.Y - _F5.Y,
			},
			TextInstance{Parts: []TextOrMaths{
				staticMath(`)`),
			}},
		},
	},
	{
		Title: "Droites du plan", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Lire l'équation de la droite (d)."),
			}},
			FigureInstance{
				Figure: figure6,
			},
			FormulaDisplayInstance{Parts: []StringOrExpression{
				{String: "y = ax + b"},
			}},
			TextInstance{Parts: []TextOrMaths{
				text("avec "),
				staticMath("a = "),
			}},
			NumberFieldInstance{ID: 0, Answer: _line.A},
			TextInstance{Parts: []TextOrMaths{
				text(" et "),
				staticMath("b = "),
			}},
			NumberFieldInstance{ID: 1, Answer: _line.B},
		},
	},
	{
		Title: "Droites du plan", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Construire la droite (d') d'équation"),
			}},
			FormulaDisplayInstance{Parts: []StringOrExpression{
				{String: `y =`},
				{Expression: mustParse("((-1)/3)x + 2")},
			}},
			FigureAffineLineFieldInstance{
				Figure: repere.Figure{
					Origin:   repere.Coord{X: 2, Y: 2},
					Width:    7,
					Height:   7,
					ShowGrid: true,
				},
				Label:  "(d')",
				Answer: [2]float64{-1. / 3, 2},
			},
		},
	},
	{
		Title: "Droites du plan", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{
				text("Quelle est la position relative entre"),
			}},
			FormulaDisplayInstance{Parts: []StringOrExpression{
				{String: `(d) : `},
				{Expression: mustParse("3x - 7y + 1")},
				{String: `= 0`},
			}},
			FormulaDisplayInstance{Parts: []StringOrExpression{
				{String: `(d') : `},
				{Expression: mustParse("x + 2y - 7")},
				{String: `= 0`},
			}},
			RadioFieldInstance{
				ID: 0,
				Proposals: []client.ListFieldProposal{
					{Content: []client.TextOrMath{{Text: "(d) est au dessus de (d')"}}},
					{Content: []client.TextOrMath{{Text: "(d) est au dessous de (d')"}}},
					{Content: []client.TextOrMath{{Text: "Ni l'un, ni l'autre"}}},
				},
				Answer: 2,
			},
		},
	},
	{
		Title: "Droites du plan", Enonce: EnonceInstance{
			TextInstance{Parts: []TextOrMaths{text("Quelle est l'équation réduite de :")}},
			FormulaDisplayInstance{Parts: []StringOrExpression{
				{String: `(d) : `},
				{Expression: mustParse("3x - 7y + 1")},
				{String: `= 0`},
			}},
			ExpressionFieldInstance{
				ID:              0,
				Label:           StringOrExpression{String: "y = "},
				ComparisonLevel: expression.ExpandedSubstitutions,
				Answer:          mustParse("(3/7)x + 1"),
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
