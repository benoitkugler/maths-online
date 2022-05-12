package examples

import (
	ex "github.com/benoitkugler/maths-online/maths/exercice"
	"github.com/benoitkugler/maths-online/maths/expression"
	functiongrapher "github.com/benoitkugler/maths-online/maths/function_grapher"
	"github.com/benoitkugler/maths-online/maths/repere"
)

var blockList = [...]ex.Block{
	ex.ExpressionFieldBlock{
		Expression:      "x^2 + 2x + 1",
		Label:           ex.NewPExpr("(x+1)^2"),
		ComparisonLevel: ex.SimpleSubstitutions,
	},
	ex.FigureAffineLineFieldBlock{
		Label: "f",
		A:     "1",
		B:     "3",
		Figure: ex.FigureBlock{
			Drawings: repere.RandomDrawings{},
			Bounds: repere.RepereBounds{
				Width:  10,
				Height: 10,
				Origin: repere.Coord{
					X: 3,
					Y: 3,
				},
			},
			ShowGrid: true,
		},
	},
	ex.FigureBlock{
		Drawings: repere.RandomDrawings{
			Points: []repere.NamedRandomLabeledPoint{
				{
					Name: "A", Point: repere.RandomLabeledPoint{
						Coord: repere.RandomCoord{X: "1", Y: "2"},
					},
				},
				{
					Name: "B", Point: repere.RandomLabeledPoint{
						Coord: repere.RandomCoord{X: "1", Y: "4"},
					},
				},
				{
					Name: "C", Point: repere.RandomLabeledPoint{
						Coord: repere.RandomCoord{X: "3", Y: "2"},
					},
				},
			},
			Segments: []repere.Segment{
				{From: "A", To: "B"},
				{From: "B", To: "C"},
				{From: "C", To: "A"},
			},
			Lines: []repere.RandomLine{
				{Label: "g", A: "3", B: "-1", Color: "#FF0022"},
			},
		},
		Bounds: repere.RepereBounds{
			Width:  20,
			Height: 20,
			Origin: repere.Coord{X: 5, Y: 5},
		},
		ShowGrid: true,
	},
	ex.FigurePointFieldBlock{
		Answer: ex.CoordExpression{X: "2", Y: "-1"},
	},
	ex.FigureVectorFieldBlock{
		Answer:         ex.CoordExpression{X: "2", Y: "-1"},
		AnswerOrigin:   ex.CoordExpression{X: "2", Y: "-1"},
		MustHaveOrigin: true,
	},
	ex.FigureVectorPairFieldBlock{
		Criterion: ex.VectorColinear,
	},
	ex.FormulaBlock{
		Parts: "Soit $f(x) = &2x + 1&$. Quelle est la dérivée de f ?",
	},
	ex.FunctionGraphBlock{
		Functions: []ex.FunctionDefinition{
			{
				Function: "x^2 - 5",
				Decoration: functiongrapher.FunctionDecoration{
					Label: "C_g",
					Color: "#FF0000",
				},
				Variable: expression.NewVar('x'),
				Range:    [2]float64{-3, 3},
			},
			{
				Function: "1/x",
				Decoration: functiongrapher.FunctionDecoration{
					Label: "y = 1/x",
					Color: "#FF00BB",
				},
				Variable: expression.NewVar('x'),
				Range:    [2]float64{-3, -0.1},
			},
			{
				Function: "1/x",
				Decoration: functiongrapher.FunctionDecoration{
					Label: "y = 1/x",
					Color: "#FF00BB",
				},
				Variable: expression.NewVar('x'),
				Range:    [2]float64{0.1, 3},
			},
		},
	},
	ex.FunctionPointsFieldBlock{
		Function: "3x- 5",
		Label:    "h",
		Variable: expression.NewVar('x'),
		XGrid:    []int{-3, -2, -1, 0, 1, 2, 3},
	},
	ex.FunctionVariationGraphBlock{
		Label: "y = h(x)",
		Xs:    []string{"-2", "0", "1/3"},
		Fxs:   []string{"4/9", "-2", "2/3"},
	},
	ex.NumberFieldBlock{Expression: "1.2"},
	ex.OrderedListFieldBlock{
		Label: "x \\in",
		Answer: []ex.TextPart{
			ex.NewPMath("["),
			ex.NewPMath("-2"),
			ex.NewPMath(";"),
			ex.NewPMath("+\\infty"),
			ex.NewPMath("["),
		},
		AdditionalProposals: []ex.TextPart{
			ex.NewPMath("\\{"),
			ex.NewPMath("\\}"),
			ex.NewPMath("3"),
		},
	},
	ex.RadioFieldBlock{
		Proposals: []ex.Interpolated{
			"Oui",
			"Non",
			"x = 4",
		},
		Answer:     "3",
		AsDropDown: false,
	},
	ex.SignTableBlock{
		Label: "g",
		Xs: []ex.Interpolated{
			"$-\\infty$",
			"&1/3&",
			"$3$",
		},
		FxSymbols: []ex.SignSymbol{
			ex.Nothing,
			ex.Zero,
			ex.ForbiddenValue,
		},
		Signs: []bool{true, false},
	},
	ex.TableBlock{
		VerticalHeaders: []ex.TextPart{
			ex.NewPText("Issues"),
			ex.NewPText("Probabilités"),
		},
		HorizontalHeaders: nil,
		Values: [][]ex.TextPart{
			{
				ex.NewPText("Rouge"), ex.NewPExpr("Vert"), ex.NewPExpr("Bleu"),
			},
			{
				ex.NewPExpr("1/3"), ex.NewPExpr("1/3"), ex.NewPExpr("1/3"),
			},
		},
	},
	ex.TableFieldBlock{
		HorizontalHeaders: []ex.TextPart{
			ex.NewPText("Homme"),
			ex.NewPText("Femme"),
		},
		VerticalHeaders: []ex.TextPart{
			ex.NewPText("Salarié"),
			ex.NewPText("Chomeur"),
		},
		Answer: [][]string{
			{"899", "253"},
			{"520", "50"},
		},
	},
	ex.TextBlock{
		Parts:   "Soit $f(x) = &2x + 1&$. Calculer $f'$.",
		Bold:    true,
		Italic:  true,
		Smaller: true,
	},
	ex.TreeFieldBlock{
		EventsProposals: []string{"P", "F", "?"},
		AnswerRoot: ex.TreeNodeAnswer{
			Probabilities: []string{"1/3", "2/3"},
			Children: []ex.TreeNodeAnswer{
				{
					Value: 0,
				},
				{
					Value: 1,
				},
			},
		},
	},
	ex.VariationTableBlock{
		Label: "g(x)",
		Xs:    []string{"-5", "0", "2/3"},
		Fxs:   []string{"4.5", "2/9", "-12"},
	},
	ex.VariationTableFieldBlock{
		Answer: ex.VariationTableBlock{
			Label: "g(x)",
			Xs:    []string{"-5", "0", "2/3"},
			Fxs:   []string{"4.5", "2/9", "-12"},
		},
	},
}
