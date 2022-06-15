package examples

import (
	ex "github.com/benoitkugler/maths-online/maths/exercice"
	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/functiongrapher"
	"github.com/benoitkugler/maths-online/maths/repere"
)

var questions = [...]ex.QuestionPage{
	{
		Title: "Remplir un tableau de variation",
		Enonce: ex.Enonce{
			ex.TextBlock{
				Parts: "Réponse attendue : -2, 0, 2/3 \n 4/9, -2, 2/3",
			},
			ex.VariationTableFieldBlock{
				Answer: ex.VariationTableBlock{
					Label: "y = h(x)",
					Xs:    []string{"-2", "0", "2/3"},
					Fxs:   []string{"4/9", "-2", "2/3"},
				},
			},
		},
	},
}

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
				From:     "-3", To: "3",
			},
			{
				Function: "1/x",
				Decoration: functiongrapher.FunctionDecoration{
					Label: "y = 1/x",
					Color: "#FF00BB",
				},
				Variable: expression.NewVar('x'),
				From:     "-3", To: "-0.1",
			},
			{
				Function: "1/x",
				Decoration: functiongrapher.FunctionDecoration{
					Label: "y = 1/x",
					Color: "#FF00BB",
				},
				Variable: expression.NewVar('x'),
				From:     "0.1", To: "3",
			},
		},
	},
	ex.FunctionPointsFieldBlock{
		Function: "3x- 5",
		Label:    "h",
		Variable: expression.NewVar('x'),
		XGrid:    []string{"-3", "-2", "-1", "0", "1", "2", "3"},
	},

	ex.NumberFieldBlock{Expression: "1.2"},
	ex.OrderedListFieldBlock{
		Label: "x \\in",
		Answer: []ex.Interpolated{
			"[",
			"-2",
			";",
			"$+\\infty$",
			"[",
		},
		AdditionalProposals: []ex.Interpolated{
			"$\\{$",
			"$\\}$",
			"3",
			"Un long texte",
			"Un long texte",
			"Un long texte",
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
