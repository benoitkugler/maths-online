package examples

// import (
// 	_ "embed"
// 	"math"

// 	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
// 	"github.com/benoitkugler/maths-online/server/src/maths/expression"
// 	"github.com/benoitkugler/maths-online/server/src/maths/functiongrapher"
// 	"github.com/benoitkugler/maths-online/server/src/maths/repere"
// )

// func mustEvaluate(s string, vars ...expression.Variables) float64 {
// 	if len(vars) > 0 {
// 		return expression.mustEvaluate(s, vars[0])
// 	}
// 	return expression.mustEvaluate(s, nil)
// }

// func text(s string) client.TextOrMath {
// 	return client.TextOrMath{Text: s}
// }

// func staticMath(s string) client.TextOrMath {
// 	return client.TextOrMath{Text: s, IsMath: true}
// }

// func expr(s string) client.TextOrMath {
// 	return client.TextOrMath{Text: expression.MustParse(s).AsLaTeX(), IsMath: true}
// }

// var (
// 	_A = repere.Coord{X: 5, Y: 25}
// 	_B = repere.Coord{X: 5, Y: 5}
// 	_C = repere.Coord{X: 45, Y: 5}
// 	_D = repere.Coord{X: 45, Y: 25}
// 	_K = repere.OrthogonalProjection(_D, _A, _C)

// 	figure1 = repere.Figure{
// 		Bounds: repere.RepereBounds{
// 			Width:  50,
// 			Height: 30,
// 		},
// 		Drawings: repere.Drawings{
// 			Points: map[string]repere.LabeledPoint{
// 				"A": {Point: _A, Pos: repere.TopLeft},
// 				"B": {Point: _B, Pos: repere.BottomLeft},
// 				"C": {Point: _C, Pos: repere.BottomRight},
// 				"D": {Point: _D, Pos: repere.TopRight},
// 				"H": {Point: repere.OrthogonalProjection(_B, _A, _C), Pos: repere.Top},
// 				"K": {Point: _K, Pos: repere.Top},
// 			},
// 			Segments: []repere.Segment{
// 				{LabelName: "", From: "A", To: "B", LabelPos: repere.Left},
// 				{LabelName: "", From: "B", To: "C", LabelPos: repere.Bottom},
// 				{LabelName: "", From: "C", To: "D", LabelPos: repere.Right},
// 				{LabelName: "", From: "D", To: "A", LabelPos: repere.Top},

// 				// diagonal
// 				{LabelName: "", From: "A", To: "C", LabelPos: repere.Bottom},

// 				{LabelName: "", From: "B", To: "H", LabelPos: repere.Bottom},
// 				{LabelName: "", From: "D", To: "K", LabelPos: repere.Top},
// 			},
// 		},
// 	}
// )

// var (
// 	__A = repere.Coord{X: -1, Y: -1}
// 	__D = repere.Coord{X: -1, Y: 1}
// 	__B = repere.Coord{X: 1, Y: 1}
// 	__J = repere.Coord{X: -3, Y: 0}
// 	__C = repere.Coord{X: -5, Y: 4}

// 	figure2 = repere.Figure{
// 		Bounds: repere.RepereBounds{
// 			Origin: repere.Coord{X: 6, Y: 6},
// 			Width:  12,
// 			Height: 12,
// 		},
// 		Drawings: repere.Drawings{
// 			Points: map[string]repere.LabeledPoint{
// 				"A": {Point: __A, Pos: repere.TopLeft},
// 				"B": {Point: __B, Pos: repere.BottomLeft},
// 				"C": {Point: __C, Pos: repere.BottomRight},
// 				"D": {Point: __D, Pos: repere.TopRight},
// 				"J": {Point: __J, Pos: repere.TopRight},
// 			},
// 			Segments: []repere.Segment{},
// 		},
// 		ShowGrid: true,
// 	}

// 	figure3 = repere.Figure{
// 		Bounds: repere.RepereBounds{
// 			Origin: repere.Coord{X: 6, Y: 6},
// 			Width:  8,
// 			Height: 8,
// 		},
// 		Drawings: repere.Drawings{
// 			Points: map[string]repere.LabeledPoint{
// 				"A": {Point: __A, Pos: repere.TopLeft},
// 				"B": {Point: __B, Pos: repere.BottomLeft},
// 				"D": {Point: __D, Pos: repere.TopRight},
// 			},
// 			Segments: []repere.Segment{},
// 		},
// 		ShowGrid: true,
// 	}

// 	_A4     = repere.Coord{X: 4, Y: 9}
// 	_B4     = repere.Coord{X: 6, Y: 6}
// 	_C4     = repere.Coord{X: 5, Y: 12}
// 	_D4     = repere.Coord{X: 2, Y: 0}
// 	figure4 = repere.Figure{
// 		Bounds: repere.RepereBounds{
// 			Width:  7,
// 			Height: 13,
// 		},
// 		Drawings: repere.Drawings{
// 			Points: map[string]repere.LabeledPoint{
// 				"A": {Point: _A4, Pos: repere.TopLeft},
// 				"B": {Point: _B4, Pos: repere.BottomLeft},
// 				"C": {Point: _C4, Pos: repere.TopRight},
// 				"D": {Point: _D4, Pos: repere.TopRight},
// 			},
// 			Segments: []repere.Segment{},
// 		},
// 	}

// 	_A5     = repere.Coord{X: 1, Y: 1}
// 	_B5     = repere.Coord{X: -1, Y: 5}
// 	_J5     = repere.Coord{X: -1, Y: 3}
// 	_H5     = repere.Coord{X: -2, Y: 4}
// 	_F5     = repere.Coord{X: -3, Y: 4}
// 	_G5     = repere.Coord{X: -2, Y: 2}
// 	figure5 = repere.Figure{
// 		Bounds: repere.RepereBounds{
// 			Origin: repere.Coord{X: 4, Y: 0},
// 			Width:  7,
// 			Height: 8,
// 		},
// 		ShowGrid: true,
// 		Drawings: repere.Drawings{
// 			Points: map[string]repere.LabeledPoint{
// 				"A": {Point: _A5, Pos: repere.BottomLeft},
// 				"B": {Point: _B5, Pos: repere.BottomLeft},
// 				"J": {Point: _J5, Pos: repere.TopRight},
// 				"H": {Point: _H5, Pos: repere.TopRight},
// 				"F": {Point: _F5, Pos: repere.TopRight},
// 				"G": {Point: _G5, Pos: repere.TopRight},
// 			},
// 			Segments: []repere.Segment{
// 				{From: "A", To: "B", AsVector: true},
// 				{From: "F", To: "G", AsVector: true},
// 				{From: "J", To: "H", AsVector: true},
// 			},
// 		},
// 	}

// 	_line   = repere.Line{A: 3. / 2, B: 1, Label: "(d)"}
// 	figure6 = repere.Figure{
// 		Bounds: repere.RepereBounds{
// 			Origin: repere.Coord{X: 2, Y: 1},
// 			Width:  6,
// 			Height: 8,
// 		},
// 		ShowGrid: true,
// 		Drawings: repere.Drawings{
// 			Lines: []repere.Line{
// 				_line,
// 			},
// 		},
// 	}
// )

// // pythagorian triplet
// var (
// 	pythagorians = expression.PythagorianTriplet{
// 		A: expression.NewVar('a'), B: expression.NewVar('b'), C: expression.NewVar('c'),
// 		Bound: 20,
// 	}
// 	distanceParams = expression.RandomParameters{
// 		expression.NewVar('X'): expression.MustParse("-randInt(100;200)"),
// 		expression.NewVar('Y'): expression.MustParse("300"),
// 		expression.NewVar('A'): expression.MustParse("a + X"),
// 		expression.NewVar('B'): expression.MustParse("b + Y"),
// 	}
// 	distanceSample expression.Variables
// )

// func init() {
// 	pythagorians.MergeTo(distanceParams)

// 	var err error
// 	distanceSample, err = distanceParams.Instantiate()
// 	if err != nil {
// 		panic(err)
// 	}

// 	allFields, err := loadAllFieldsQuestion()
// 	if err != nil {
// 		panic(err)
// 	}

// 	PredefinedQuestions = append([]QuestionInstance{allFields}, PredefinedQuestions...)

// 	PredefinedQuestions = append(PredefinedQuestions, QuestionInstance{
// 		Title: "Repérage dans le plan", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Soient les points F("),
// 				staticMath(expression.Number(distanceSample[expression.NewVar('A')].N).String()),
// 				text(";"),
// 				staticMath(expression.Number(distanceSample[expression.NewVar('B')].N).String()),
// 				text(") et G("),
// 				staticMath(expression.Number(distanceSample[expression.NewVar('X')].N).String()),
// 				text(";"),
// 				staticMath(expression.Number(distanceSample[expression.NewVar('Y')].N).String()),
// 				text("). Calculer FG."),
// 			}},
// 			// TextInstance{
// 			// 	IsHint: true,
// 			// 	Parts: []client.TextOrMath{
// 			// 		text("(On utilisera sqrt(10) pour "),
// 			// 		staticMath(`\sqrt{10}`),
// 			// 		text(")."),
// 			// 	},
// 			// },
// 			// ExpressionFieldInstance{
// 			// 	ID:              0,
// 			// 	Label:           StringOrExpression{String: "FG = "},
// 			// 	Answer:          expression.MustParse("sqrt(1262900)"),
// 			// 	ComparisonLevel: expression.SimpleSubstitutions,
// 			// },
// 			TextInstance{Parts: []client.TextOrMath{
// 				staticMath("FG = "),
// 			}},
// 			NumberFieldInstance{ID: 0, Answer: distanceSample[expression.NewVar('c')].N},
// 		},
// 	},
// 	)
// }

// func valuesToExpr(vs ...float64) []evaluatedExpression {
// 	out := make([]evaluatedExpression, len(vs))
// 	for i, v := range vs {
// 		out[i] = evaluatedExpression{Value: v, Expr: expression.NewNb(v)}
// 	}
// 	return out
// }

// var PredefinedQuestions = []QuestionInstance{
// 	{
// 		Title: "Calcul littéral", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{text("Développer l’expression :")}},
// 			ExpressionFieldInstance{
// 				ID:              0,
// 				Label:           StringOrExpression{Expression: expression.MustParse("(x−6)*(4*x−3)")},
// 				ComparisonLevel: expression.SimpleSubstitutions,
// 				Answer:          expression.MustParse("4*x^2 - 27 *x + 18"),
// 			},
// 		},
// 	},
// 	{
// 		Title: "Calcul littéral", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{text("Écrire sous une seule fraction : ")}},
// 			FormulaDisplayInstance{
// 				expression.MustParse("(1/3)+(2/5)").AsLaTeX(),
// 				`= \frac{a}{b}`,
// 			},
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("avec "),
// 				staticMath("a = "),
// 			}},
// 			NumberFieldInstance{ID: 0, Answer: 1*5 + 2*3},
// 			TextInstance{Parts: []client.TextOrMath{
// 				text(" et "),
// 				staticMath("b = "),
// 			}},
// 			NumberFieldInstance{ID: 1, Answer: 3 * 5},
// 		},
// 	},
// 	{
// 		Title: "Calcul littéral", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Écrire sous la forme "),
// 				staticMath(`a\sqrt{b},`),
// 				text(" où "),
// 				staticMath("a"),
// 				text(" et "),
// 				staticMath("b"),
// 				text(" sont des entiers les plus petits possibles :"),
// 			}},
// 			FormulaDisplayInstance{
// 				expression.MustParse("sqrt(50)").AsLaTeX(),
// 				" = ",
// 				expression.MustParse("a * sqrt(b)").AsLaTeX(),
// 			},
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("avec "),
// 				staticMath("a = "),
// 			}},
// 			NumberFieldInstance{ID: 0, Answer: 5},
// 			TextInstance{Parts: []client.TextOrMath{
// 				text(" et "),
// 				staticMath("b = "),
// 			}},
// 			NumberFieldInstance{ID: 1, Answer: 2},
// 		},
// 	},
// 	{
// 		Title: "Nombres entiers", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				expr("89"),
// 				text(" est-il un nombre premier ?"),
// 			}},
// 			RadioFieldInstance{
// 				ID:     0,
// 				Answer: 0,
// 				Proposals: []client.ListFieldProposal{
// 					{Content: []client.TextOrMath{{Text: "Oui"}}},
// 					{Content: []client.TextOrMath{{Text: "Non"}}},
// 				},
// 			},
// 		},
// 	},
// 	{
// 		Title: "Nombres entiers", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				expr("987"),
// 				text(" est-il un mulitple de "),
// 				expr("3"),
// 			}},
// 			RadioFieldInstance{
// 				ID:     0,
// 				Answer: 0,
// 				Proposals: []client.ListFieldProposal{
// 					{Content: []client.TextOrMath{{Text: "Oui"}}},
// 					{Content: []client.TextOrMath{{Text: "Non"}}},
// 				},
// 			},
// 		},
// 	},
// 	{
// 		Title: "Nombres réels", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Traduire en intervalle : "),
// 				staticMath(`x \ge `),
// 				expr("12"),
// 			}},
// 			OrderedListFieldInstance{
// 				ID:    0,
// 				Label: `x \in`,
// 				Answer: []string{ // [12;+infty[
// 					"[",
// 					"12",
// 					";",
// 					"+",
// 					`\infty`,
// 					`[`,
// 				},
// 				AdditionalProposals: []string{
// 					"]", // some duplicats
// 					`\infty`,
// 					"11",
// 					"-",
// 				},
// 			},
// 		},
// 	},
// 	{
// 		Title: "Nombres réels", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("A quel plus petit ensemble appartient "),
// 				expr("4/7"),
// 			}},
// 			RadioFieldInstance{
// 				ID:     0,
// 				Answer: 4,
// 				Proposals: []client.ListFieldProposal{
// 					{Content: []client.TextOrMath{{Text: `\mathbb{N}`, IsMath: true}}},
// 					{Content: []client.TextOrMath{{Text: `\mathbb{Z}`, IsMath: true}}},
// 					{Content: []client.TextOrMath{{Text: `\mathbb{D}`, IsMath: true}}},
// 					{Content: []client.TextOrMath{{Text: `\mathbb{Q}`, IsMath: true}}},
// 					{Content: []client.TextOrMath{{Text: `\mathbb{R}`, IsMath: true}}},
// 				},
// 			},
// 		},
// 	},
// 	{
// 		Title: "Nombres réels", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Déterminer "),
// 				staticMath(`I \cap J`),
// 				text(" où "),
// 				staticMath(`I = [`),
// 				expr("3"),
// 				staticMath(`;`),
// 				expr("9"),
// 				staticMath(`]`),
// 				text(" et "),
// 				staticMath(`J = [`),
// 				expr("0"),
// 				staticMath(`;`),
// 				expr("6"),
// 				staticMath(`]`),
// 			}},
// 			OrderedListFieldInstance{
// 				ID:    0,
// 				Label: `I \cap J =`,
// 				Answer: []string{ // [12;+infty[
// 					"[",
// 					"3",
// 					";",
// 					"6",
// 					`]`,
// 				},
// 				AdditionalProposals: []string{
// 					"]", // some duplicates
// 					`\infty`,
// 					"0",
// 					"9",
// 					"-",
// 				},
// 			},
// 		},
// 	},
// 	{
// 		Title: "Equations et inéquations", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Résoudre dans "),
// 				staticMath(`\mathbb{R}`),
// 				text("l'équation :"),
// 			}},
// 			FormulaDisplayInstance{
// 				expression.MustParse("20*x - 7").AsLaTeX(),
// 				"=",
// 				expression.MustParse("34").AsLaTeX(),
// 			},
// 			ExpressionFieldInstance{
// 				ID:              0,
// 				Label:           StringOrExpression{String: "x ="},
// 				ComparisonLevel: expression.SimpleSubstitutions,
// 				Answer:          expression.MustParse("(34+7)/20"),
// 			},
// 		},
// 	},
// 	{
// 		Title: "Equations et inéquations", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Résoudre dans "),
// 				staticMath(`\mathbb{R}`),
// 				text("l'inéquation :"),
// 			}},
// 			FormulaDisplayInstance{
// 				expression.MustParse("-8*x +1").AsLaTeX(),
// 				"<",
// 				expression.MustParse("2*x - 5").AsLaTeX(),
// 			},
// 			TextInstance{Parts: []client.TextOrMath{
// 				staticMath("x "),
// 			}},
// 			DropDownFieldInstance{
// 				ID:     0,
// 				Answer: 1,
// 				Proposals: []client.ListFieldProposal{
// 					{Content: []client.TextOrMath{{Text: `<`, IsMath: true}}},
// 					{Content: []client.TextOrMath{{Text: `>`, IsMath: true}}},
// 					{Content: []client.TextOrMath{{Text: `=`, IsMath: true}}},
// 				},
// 			},
// 			NumberFieldInstance{ID: 1, Answer: 0.6},
// 		},
// 	},
// 	{
// 		Title: "Equations et inéquations", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Résoudre dans "),
// 				staticMath(`\mathbb{R}`),
// 				text("l'équation :"),
// 			}},
// 			FormulaDisplayInstance{
// 				expression.MustParse("(x-7)*(4*x - 9)").AsLaTeX(),
// 				"=",
// 				expression.MustParse("0").AsLaTeX(),
// 			},
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Solutions : "),
// 				staticMath("x = "),
// 			}},
// 			NumberFieldInstance{ID: 0, Answer: 7},
// 			TextInstance{Parts: []client.TextOrMath{
// 				text(" ou "),
// 				staticMath("x = "),
// 			}},
// 			NumberFieldInstance{ID: 1, Answer: 9. / 4},
// 		},
// 	},
// 	{
// 		Title: "Equations et inéquations", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Combien l’équation "),
// 				staticMath(`f(x) = `),
// 				expr("-1"),
// 				text(" admet-elle de solutions ? "),
// 			}},
// 			VariationTableInstance{
// 				Label: "f",
// 				Xs:    valuesToExpr(-20, -10, 0, 3, 18),
// 				Fxs:   valuesToExpr(-6, -2, -8, 0, -5),
// 			},
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Nombre de solutions : "),
// 			}},
// 			NumberFieldInstance{ID: 0, Answer: 2},
// 		},
// 	},
// 	{
// 		Title: "Equations et inéquations", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Résoudre l’inéquation "),
// 				staticMath(`f(x) \ge 0`),
// 			}},
// 			SignTableInstance{
// 				Xs: []string{
// 					`-\infty`,
// 					"-2",
// 					"0",
// 					"4",
// 					`+\infty`,
// 				},
// 				FxSymbols: []SignSymbol{Nothing, Zero, Zero, ForbiddenValue, Nothing},
// 				Signs:     []bool{false, true, false, true},
// 			},
// 			// SignTableInstance{Columns: []client.SignColumn{
// 			// 	{X: `-\infty`, IsYForbiddenValue: false, IsSign: false, IsPositive: false},
// 			// 	{X: "", IsYForbiddenValue: false, IsSign: true, IsPositive: false},
// 			// 	{X: "-2", IsYForbiddenValue: false, IsSign: false, IsPositive: true},
// 			// 	{X: "", IsYForbiddenValue: false, IsSign: true, IsPositive: true},
// 			// 	{X: "0", IsYForbiddenValue: false, IsSign: false, IsPositive: true},
// 			// 	{X: "", IsYForbiddenValue: false, IsSign: true, IsPositive: false},
// 			// 	{X: "4", IsYForbiddenValue: true, IsSign: false, IsPositive: true},
// 			// 	{X: "", IsYForbiddenValue: false, IsSign: true, IsPositive: true},
// 			// 	{X: `+\infty`, IsYForbiddenValue: false, IsSign: false, IsPositive: false},
// 			// }},
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Intervalle(s) solution(s) :"),
// 			}},
// 			OrderedListFieldInstance{
// 				ID: 0,
// 				Answer: []string{ // [12;+infty[
// 					"[",
// 					"-2",
// 					";",
// 					"0",
// 					`]`,
// 					`\cup`,
// 					"]",
// 					"4",
// 					";",
// 					`+\infty`,
// 					`[`,
// 				},
// 				AdditionalProposals: []string{
// 					"]", // some duplicates
// 					`-\infty`,
// 					"0",
// 					"-",
// 				},
// 			},
// 		},
// 	},
// 	// geometrie

// 	{
// 		Title: "Géométrie plane", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("ABC est un triangle rectangle en C avec "),
// 				staticMath(`AB = `),
// 				expr("sqrt(98)"),
// 				text(" et "),
// 				staticMath("BC = "),
// 				expr("7"),
// 				text(". Calculer, en degrés, "),
// 				staticMath(`\widehat{ABC}.`),
// 			}},
// 			TextInstance{Parts: []client.TextOrMath{
// 				staticMath(`\widehat{ABC} = `),
// 			}},
// 			NumberFieldInstance{ID: 0, Answer: mustEvaluate("acos(7/sqrt(98))") * 180 / math.Pi},
// 		},
// 	},
// 	{
// 		Title: "Géométrie plane", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("ABC est un triangle rectangle en C avec "),
// 				staticMath(`AB = `),
// 				expr("sqrt(98)"),
// 				text(" et "),
// 				staticMath("BC = "),
// 				expr("7"),
// 				text(". Calculer AC."),
// 			}},
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("AC = "),
// 			}},
// 			NumberFieldInstance{ID: 0, Answer: mustEvaluate("sqrt(sqrt(98)^2 - 7^2)")},
// 		},
// 	},
// 	{
// 		Title: "Géométrie plane", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Soient "),
// 				staticMath(`AB = `),
// 				expr("4"),
// 				staticMath(" ; AC = "),
// 				expr("12"),
// 				text(" et "),
// 				staticMath("BC = "),
// 				expr("8"),
// 				text(". Le triangle ABC est-il rectangle ? Si oui, en quoi ?"),
// 			}},
// 			RadioFieldInstance{
// 				ID: 0,
// 				Answer: int(mustEvaluate("1 * isZero(a^2 - b^2 - c^2) + 2*isZero(b^2 - a^2 - c^2) + 3*isZero(c^2 - a^2 - b^2)", expression.Variables{
// 					expression.NewVar('a'): expression.NewNb(8),  // BC
// 					expression.NewVar('b'): expression.NewNb(12), // AC
// 					expression.NewVar('c'): expression.NewNb(4),  // AB
// 				})),
// 				Proposals: []client.ListFieldProposal{
// 					{Content: []client.TextOrMath{{Text: `ABC n'est pas rectangle.`, IsMath: false}}},
// 					{Content: []client.TextOrMath{{Text: `ABC est rectangle en A.`, IsMath: false}}},
// 					{Content: []client.TextOrMath{{Text: `ABC est rectangle en B.`, IsMath: false}}},
// 					{Content: []client.TextOrMath{{Text: `ABC est rectangle en C.`, IsMath: false}}},
// 				},
// 			},
// 		},
// 	},

// 	{
// 		Title: "Repérage dans le plan", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Soient les points "),
// 				staticMath(`A(`),
// 				expr("8"),
// 				staticMath(";"),
// 				expr("19"),
// 				staticMath(")"),
// 				text(" et "),
// 				staticMath(`B(`),
// 				expr("-6"),
// 				staticMath(";"),
// 				expr("0"),
// 				staticMath(")."),
// 				text("Quelles sont les coordonnées de M, milieu de [AB] ?"),
// 			}},
// 			TextInstance{Parts: []client.TextOrMath{
// 				staticMath("M = ("),
// 			}},
// 			NumberFieldInstance{ID: 0, Answer: mustEvaluate("(8 + (-6))/2")},
// 			TextInstance{Parts: []client.TextOrMath{
// 				staticMath(";"),
// 			}},
// 			NumberFieldInstance{ID: 1, Answer: mustEvaluate("(19 + 0)/2")},
// 			TextInstance{Parts: []client.TextOrMath{
// 				staticMath(")"),
// 			}},
// 		},
// 	},

// 	{
// 		Title: "Repérage dans le plan", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Quel point est le projeté orthogonal de D sur (AH) ?"),
// 			}},
// 			FigureInstance{Figure: figure1},
// 			RadioFieldInstance{
// 				ID:     0,
// 				Answer: 2,
// 				Proposals: []client.ListFieldProposal{
// 					{Content: []client.TextOrMath{{Text: `A`, IsMath: false}}},
// 					{Content: []client.TextOrMath{{Text: `B`, IsMath: false}}},
// 					{Content: []client.TextOrMath{{Text: `K`, IsMath: false}}},
// 					{Content: []client.TextOrMath{{Text: `H`, IsMath: false}}},
// 					{Content: []client.TextOrMath{{Text: `D`, IsMath: false}}},
// 				},
// 			},
// 		},
// 	},
// 	{
// 		Title: "Repérage dans le plan", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Construire le point M, projeté orthogonal de K sur (BC)."),
// 			}},
// 			FigurePointFieldInstance{
// 				Figure: figure1,
// 				Answer: repere.OrthogonalProjection(_K, _B, _C).Round(),
// 				ID:     0,
// 			},
// 		},
// 	},

// 	{
// 		Title: "Repérage dans le plan", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Lire les coordonnées de B."),
// 			}},
// 			FigureInstance{
// 				Figure: figure3,
// 			},
// 			TextInstance{Parts: []client.TextOrMath{
// 				staticMath("B = ("),
// 			}},
// 			NumberFieldInstance{
// 				ID:     0,
// 				Answer: __B.X,
// 			},
// 			TextInstance{Parts: []client.TextOrMath{
// 				staticMath("; "),
// 			}},
// 			NumberFieldInstance{
// 				ID:     1,
// 				Answer: __B.Y,
// 			},
// 			TextInstance{Parts: []client.TextOrMath{
// 				staticMath(")"),
// 			}},
// 		},
// 	},
// 	{
// 		Title: "Repérage dans le plan", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Quelle est la nature de ABCD ?"),
// 			}},
// 			FigureInstance{
// 				Figure: figure4,
// 			},
// 			RadioFieldInstance{
// 				ID:     0,
// 				Answer: 0,
// 				Proposals: []client.ListFieldProposal{
// 					{Content: []client.TextOrMath{{Text: "Quadrilatère quelconque"}}},
// 					{Content: []client.TextOrMath{{Text: "Rectangle"}}},
// 					{Content: []client.TextOrMath{{Text: "Losange"}}},
// 					{Content: []client.TextOrMath{{Text: "Carré"}}},
// 				},
// 			},
// 		},
// 	},
// 	{
// 		Title: "Vecteurs", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Construire le vecteur "),
// 				staticMath(`\overrightarrow{AB} + \overrightarrow{CD}.`),
// 			}},
// 			FigureVectorFieldInstance{
// 				Figure: figure2,
// 				Answer: repere.IntCoord{X: 6, Y: -1},
// 				ID:     0,
// 			},
// 		},
// 	},
// 	{
// 		Title: "Vecteurs", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Simplifier au maximum l'expression vectorielle"),
// 				staticMath(`\overrightarrow{AB} + \overrightarrow{EB} -  \overrightarrow{EG}`),
// 			}},
// 			OrderedListFieldInstance{
// 				Label: `\overrightarrow{AB} + \overrightarrow{EB} - \overrightarrow{EG} = `,
// 				Answer: []string{
// 					"G",
// 					"A",
// 				},
// 				AdditionalProposals: []string{
// 					"E",
// 					"B",
// 					"F",
// 					"A",
// 				},
// 			},
// 		},
// 	},
// 	{
// 		Title: "Vecteurs", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Construire deux vecteurs égaux."),
// 			}},
// 			FigureVectorPairFieldInstance{
// 				Figure:    figure1,
// 				Criterion: VectorEquals,
// 				ID:        0,
// 			},
// 		},
// 	},
// 	{
// 		Title: "Vecteurs", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Construire deux vecteurs colinéaires."),
// 			}},
// 			FigureVectorPairFieldInstance{
// 				Figure:    figure1,
// 				Criterion: VectorColinear,
// 				ID:        0,
// 			},
// 		},
// 	},
// 	{
// 		Title: "Vecteurs", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Lire les coordonnées de "),
// 				staticMath(`\overrightarrow{FG}.`),
// 			}},
// 			FigureInstance{
// 				Figure: figure5,
// 			},
// 			TextInstance{Parts: []client.TextOrMath{
// 				staticMath(`\overrightarrow{FG} = (`),
// 			}},
// 			NumberFieldInstance{
// 				ID:     0,
// 				Answer: _G5.X - _F5.X,
// 			},
// 			TextInstance{Parts: []client.TextOrMath{
// 				staticMath(`;`),
// 			}},
// 			NumberFieldInstance{
// 				ID:     1,
// 				Answer: _G5.Y - _F5.Y,
// 			},
// 			TextInstance{Parts: []client.TextOrMath{
// 				staticMath(`)`),
// 			}},
// 		},
// 	},
// 	{
// 		Title: "Droites du plan", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Lire l'équation de la droite (d)."),
// 			}},
// 			FigureInstance{
// 				Figure: figure6,
// 			},
// 			FormulaDisplayInstance{
// 				"y = ax + b",
// 			},
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("avec "),
// 				staticMath("a = "),
// 			}},
// 			NumberFieldInstance{ID: 0, Answer: _line.A},
// 			TextInstance{Parts: []client.TextOrMath{
// 				text(" et "),
// 				staticMath("b = "),
// 			}},
// 			NumberFieldInstance{ID: 1, Answer: _line.B},
// 		},
// 	},
// 	{
// 		Title: "Droites du plan", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Construire la droite (d') d'équation"),
// 			}},
// 			FormulaDisplayInstance{
// 				`y =`,
// 				expression.MustParse("((-1)/3)x + 2").AsLaTeX(),
// 			},
// 			FigureAffineLineFieldInstance{
// 				Figure: repere.Figure{
// 					Bounds: repere.RepereBounds{
// 						Origin: repere.Coord{X: 2, Y: 2},
// 						Width:  7,
// 						Height: 7,
// 					},
// 					ShowGrid: true,
// 				},
// 				Label:   "(d')",
// 				AnswerA: -1 / 3,
// 				AnswerB: 2,
// 			},
// 		},
// 	},
// 	{
// 		Title: "Droites du plan", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Quelle est la position relative entre"),
// 			}},
// 			FormulaDisplayInstance{
// 				`(d) : `,
// 				expression.MustParse("3x - 7y + 1").AsLaTeX(),
// 				`= 0`,
// 			},
// 			FormulaDisplayInstance{
// 				`(d') : `,
// 				expression.MustParse("x + 2y - 7").AsLaTeX(),
// 				`= 0`,
// 			},
// 			RadioFieldInstance{
// 				ID: 0,
// 				Proposals: []client.ListFieldProposal{
// 					{Content: []client.TextOrMath{{Text: "(d) est au dessus de (d')"}}},
// 					{Content: []client.TextOrMath{{Text: "(d) est au dessous de (d')"}}},
// 					{Content: []client.TextOrMath{{Text: "Ni l'un, ni l'autre"}}},
// 				},
// 				Answer: 2,
// 			},
// 		},
// 	},
// 	{
// 		Title: "Droites du plan", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{text("Quelle est l'équation réduite de :")}},
// 			FormulaDisplayInstance{
// 				`(d) : `,
// 				expression.MustParse("3x - 7y + 1").AsLaTeX(),
// 				`= 0`,
// 			},
// 			ExpressionFieldInstance{
// 				ID:              0,
// 				Label:           StringOrExpression{String: "y = "},
// 				ComparisonLevel: expression.ExpandedSubstitutions,
// 				Answer:          expression.MustParse("(3/7)x + 1"),
// 			},
// 		},
// 	},

// 	{
// 		Title: "Généralités sur les fonctions", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Construire la courbe représentative de "),
// 				staticMath("g(x) = "),
// 				expr("x^2 − 3x + 1"),
// 				text("."),
// 			}},
// 			FunctionPointsFieldInstance{
// 				ID: 0,
// 				Function: expression.FunctionExpr{
// 					Function: expression.MustParse("x^2 − 3x + 1"),
// 					Variable: expression.NewVar('x'),
// 				},
// 				Label: "g(x)",
// 				XGrid: []int{-2, -1, 0, 1, 2, 3, 4, 5},
// 			},
// 		},
// 	},
// 	{
// 		Title: "Variations", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{text("Quel est le maximum de h ?")}},
// 			FunctionGraphInstance{
// 				Functions: []expression.FunctionDefinition{{
// 					FunctionExpr: expression.FunctionExpr{
// 						Function: expression.MustParse("10*(1.1(x/10)^4 - 2.2(x/10)^3+(x/10)) + 0.047"),
// 						Variable: expression.NewVar('x'),
// 					},
// 					From: -5,
// 					To:   16,
// 				}},
// 				Decorations: []functiongrapher.FunctionDecoration{{Label: `y = h(x)`}},
// 			},
// 			TextInstance{Parts: []client.TextOrMath{text("Max : ")}},
// 			NumberFieldInstance{
// 				ID:     0,
// 				Answer: 3,
// 			},
// 		},
// 	},

// 	{
// 		Title: "Variations", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{text("Compléter le tableau de variations de h.")}},
// 			FunctionVariationGraphInstance{
// 				Xs:  valuesToExpr(0, 2, 6, 10, 12),
// 				Fxs: valuesToExpr(3, 0, 6, 1, 4),
// 			},
// 			VariationTableFieldInstance{
// 				ID: 0,
// 				Answer: VariationTableInstance{
// 					Xs:  valuesToExpr(0, 2, 6, 10, 12),
// 					Fxs: valuesToExpr(3, 0, 6, 1, 4),
// 				},
// 			},
// 		},
// 	},
// 	{
// 		Title: "Probabilités", Enonce: EnonceInstance{
// 			TableInstance{
// 				VerticalHeaders: []client.TextOrMath{
// 					{Text: "Issue"}, {Text: "Probabilité"},
// 				},
// 				Values: [][]client.TextOrMath{
// 					{
// 						{Text: "Pique"}, {Text: "Trèfle"}, {Text: "Coeur"}, {Text: "Carreau"}, {Text: "TEST"}, {Text: "TEST"}, {Text: "TEST"},
// 					},
// 					{
// 						{Text: "0.25", IsMath: true}, {Text: "x", IsMath: true}, {Text: "3x", IsMath: true}, {Text: "0.1", IsMath: true}, {}, {}, {},
// 					},
// 				},
// 			},
// 			TextInstance{Parts: []client.TextOrMath{text("Quelle est la probabilité de tirer un coeur ?")}},
// 			NumberFieldInstance{
// 				ID:     0,
// 				Answer: mustEvaluate("3 * (1 -(0.25 + 0.1)) / 4"),
// 			},
// 		},
// 	},

// 	{
// 		Title: "Probabilités", Enonce: EnonceInstance{
// 			TableInstance{
// 				HorizontalHeaders: []client.TextOrMath{
// 					{Text: "Femmes"}, {Text: "Hommes"}, {Text: "Total"},
// 				},
// 				VerticalHeaders: []client.TextOrMath{
// 					{Text: "Inexpérimentés"}, {Text: "Expérimentés"}, {Text: "Total"},
// 				},
// 				Values: [][]client.TextOrMath{
// 					{
// 						{Text: "49", IsMath: true}, {Text: "48", IsMath: true}, {Text: "97", IsMath: true},
// 					},
// 					{
// 						{Text: "191", IsMath: true}, {Text: "912", IsMath: true}, {Text: "1103", IsMath: true},
// 					},
// 					{
// 						{Text: "240", IsMath: true}, {Text: "960", IsMath: true}, {Text: "1200", IsMath: true},
// 					},
// 				},
// 			},
// 			TextInstance{Parts: []client.TextOrMath{
// 				text(`On choisit au hasard un employé de l’entreprise.
// Quel est la probabilité que ce soit un employé inexpérimenté ?`),
// 			}},
// 			NumberFieldInstance{
// 				ID:     0,
// 				Answer: mustEvaluate("97/1200"),
// 			},
// 			TextInstance{Parts: []client.TextOrMath{text("\n\nUn expérimenté si on sait que c’est un homme ?")}},
// 			NumberFieldInstance{
// 				ID:     1,
// 				Answer: mustEvaluate("48/960"),
// 			},
// 			TextInstance{Parts: []client.TextOrMath{text("\n\nQuelle est la probabilité que ce soit une femme ou un employé expérimenté ?")}},
// 			NumberFieldInstance{
// 				ID:     2,
// 				Answer: mustEvaluate("240/1200 + 912/1200"),
// 			},
// 		},
// 	},

// 	{
// 		Title: "Probabilités", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text(`Compléter ce tableau sachant que : l’entreprise comporte 1200 salariés dont 960 hommes.
// Parmi les hommes, 48 sont inexpérimentés et parmi les femmes, 191 sont expérimentées.`),
// 			}},
// 			TableFieldInstance{
// 				HorizontalHeaders: []client.TextOrMath{
// 					{Text: "Femmes"}, {Text: "Hommes"}, {Text: "Total"},
// 				},
// 				VerticalHeaders: []client.TextOrMath{
// 					{Text: "Inexpérimentés"}, {Text: "Expérimentés"}, {Text: "Total"},
// 				},
// 				Answer: client.TableAnswer{
// 					Rows: [][]float64{
// 						{
// 							49, 48, 97,
// 						},
// 						{
// 							191, 912, 1103,
// 						},
// 						{
// 							240, 960, 1200,
// 						},
// 					},
// 				},
// 			},
// 		},
// 	},

// 	{
// 		Title: "Probabilités", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text(`Construire l’arbre pondéré associé à l’expérience aléatoire suivante :
// on lance une pièce de monnaie truquée deux fois d’affilé, la probabilité d’obtenir Pile (P) est de 0,3.`),
// 			}},
// 			TextInstance{Bold: true, Parts: []client.TextOrMath{
// 				text(`On ordonnera les issues par ordre alphabétique.`),
// 			}},
// 			TreeFieldInstance{
// 				ID: 0,
// 				EventsProposals: []client.TextOrMath{
// 					{Text: "F"},
// 					{Text: "P"},
// 					{Text: "?"},
// 				},
// 				Answer: client.TreeAnswer{
// 					Root: client.TreeNodeAnswer{
// 						Probabilities: []float64{0.7, 0.3},
// 						Children: []client.TreeNodeAnswer{
// 							{
// 								Value:         0,
// 								Probabilities: []float64{0.7, 0.3},
// 								Children: []client.TreeNodeAnswer{
// 									{Value: 0},
// 									{Value: 1},
// 								},
// 							},
// 							{
// 								Value:         1,
// 								Probabilities: []float64{0.7, 0.3},
// 								Children: []client.TreeNodeAnswer{
// 									{Value: 0},
// 									{Value: 1},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	},

// 	{
// 		Title: "Très longue question horizontale", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{
// 				text("Écrire sous une seule fraction : "),
// 				expr("(1/3)+(2/5)"),
// 				text("Écrire sous une seule fraction : "),
// 				expr("(1/3)+(2/5)"),
// 				text("Écrire sous une seule fraction : "),
// 				expr("(1/3)+(2/5)"),
// 				text("Écrire sous une seule fraction : "),
// 				expr("(1/3)+(2/5)"),
// 				text("Écrire sous une seule fraction : "),
// 				expr("(1/3)+(2/5)"),
// 				text("Écrire sous une seule fraction : "),
// 				expr("(1/3)+(2/5)"),
// 				text("Écrire sous une seule fraction : "),
// 				expr("(1/3)+(2/5)"),
// 			}},
// 		},
// 	},
// 	{
// 		Title: "Très longue question verticale", Enonce: EnonceInstance{
// 			TextInstance{Parts: []client.TextOrMath{text("Écrire sous une seule fraction : ")}},
// 			FormulaDisplayInstance{
// 				expression.MustParse("(1/3)+(2/5)").AsLaTeX(),
// 			},
// 			TextInstance{Parts: []client.TextOrMath{text("Écrire sous une seule fraction : ")}},
// 			FormulaDisplayInstance{
// 				expression.MustParse("(1/3)+(2/5)").AsLaTeX(),
// 			},
// 			TextInstance{Parts: []client.TextOrMath{text("Écrire sous une seule fraction : ")}},
// 			FormulaDisplayInstance{
// 				expression.MustParse("(1/3)+(2/5)").AsLaTeX(),
// 			},
// 			TextInstance{Parts: []client.TextOrMath{text("Écrire sous une seule fraction : ")}},
// 			FormulaDisplayInstance{
// 				expression.MustParse("(1/3)+(2/5)").AsLaTeX(),
// 			},
// 			TextInstance{Parts: []client.TextOrMath{text("Écrire sous une seule fraction : ")}},
// 			FormulaDisplayInstance{
// 				expression.MustParse("(1/3)+(2/5)").AsLaTeX(),
// 			},
// 		},
// 	},
// }
