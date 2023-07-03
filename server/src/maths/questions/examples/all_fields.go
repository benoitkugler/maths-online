package examples

import (
	"github.com/benoitkugler/maths-online/server/src/maths/expression"
	"github.com/benoitkugler/maths-online/server/src/maths/functiongrapher"
	que "github.com/benoitkugler/maths-online/server/src/maths/questions"
	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	"github.com/benoitkugler/maths-online/server/src/maths/repere"
)

var BlockList = [...]que.Block{
	que.ExpressionFieldBlock{
		Expression:      "x^2 + 2x + 1",
		Label:           "&(x+1)^2& =",
		ComparisonLevel: que.SimpleSubstitutions,
	},
	que.FigureAffineLineFieldBlock{
		Label: "f",
		A:     "1",
		B:     "3",
		Figure: que.FigureBlock{
			Drawings: repere.RandomDrawings{},
			Bounds: repere.RepereBounds{
				Width:  10,
				Height: 10,
				Origin: repere.Coord{
					X: 3,
					Y: 3,
				},
			},
			ShowGrid:   true,
			ShowOrigin: true,
		},
	},
	que.FigureBlock{
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
			Segments: []repere.RandomSegment{
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
		ShowGrid:   true,
		ShowOrigin: true,
	},
	que.FigurePointFieldBlock{
		Figure: que.FigureBlock{
			Bounds: repere.RepereBounds{
				Width:  10,
				Height: 10,
				Origin: repere.Coord{
					X: 3,
					Y: 3,
				},
			},
		},
		Answer: que.CoordExpression{X: "2", Y: "-1"},
	},
	que.FigureVectorFieldBlock{
		Figure: que.FigureBlock{
			Bounds: repere.RepereBounds{
				Width:  10,
				Height: 10,
				Origin: repere.Coord{
					X: 3,
					Y: 3,
				},
			},
		},
		Answer:         que.CoordExpression{X: "2", Y: "-1"},
		AnswerOrigin:   que.CoordExpression{X: "2", Y: "-1"},
		MustHaveOrigin: true,
	},
	que.FigureVectorPairFieldBlock{
		Figure: que.FigureBlock{
			Bounds: repere.RepereBounds{
				Width:  10,
				Height: 10,
				Origin: repere.Coord{
					X: 3,
					Y: 3,
				},
			},
		},
		Criterion: que.VectorColinear,
	},
	que.FormulaBlock{
		Parts: "f(x) = &2x + 1&",
	},
	que.FunctionsGraphBlock{
		FunctionExprs: []que.FunctionDefinition{
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
	que.FunctionPointsFieldBlock{
		Function: "3x- 5",
		Label:    "h",
		Variable: expression.NewVar('x'),
		XGrid:    []string{"-3", "-2", "-1", "0", "1", "2", "3"},
	},
	que.NumberFieldBlock{Expression: "1.2"},
	que.OrderedListFieldBlock{
		Label: "x \\in",
		Answer: []que.Interpolated{
			"[",
			"-2",
			";",
			"$+\\infty$",
			"[",
		},
		AdditionalProposals: []que.Interpolated{
			"$\\{$",
			"$\\}$",
			"3",
			"Un long texte",
			"Un long texte",
			"Un long texte",
		},
	},
	que.RadioFieldBlock{
		Proposals: []que.Interpolated{
			"Oui",
			"Non",
			"x = 4",
		},
		Answer:     "3",
		AsDropDown: false,
	},
	que.SignTableBlock{
		Xs: []string{
			"-inf",
			"1/2",
			"3",
		},
		Functions: []client.FunctionSign{
			{
				Label: "g",
				FxSymbols: []client.SignSymbol{
					client.Nothing,
					client.Zero,
					client.ForbiddenValue,
				},
				Signs: []bool{true, false},
			},
			{
				Label: "h",
				FxSymbols: []client.SignSymbol{
					client.Nothing,
					client.Zero,
					client.ForbiddenValue,
				},
				Signs: []bool{false, true},
			},
		},
	},
	que.TableBlock{
		VerticalHeaders: []que.TextPart{
			que.NewPText("Issues"),
			que.NewPText("Probabilités"),
		},
		HorizontalHeaders: nil,
		Values: [][]que.TextPart{
			{
				que.NewPText("Rouge"), que.NewPExpr("Vert"), que.NewPExpr("Bleu"),
			},
			{
				que.NewPExpr("1/3"), que.NewPExpr("1/3"), que.NewPExpr("1/3"),
			},
		},
	},
	que.TableFieldBlock{
		HorizontalHeaders: []que.TextPart{
			que.NewPText("Homme"),
			que.NewPText("Femme"),
		},
		VerticalHeaders: []que.TextPart{
			que.NewPText("Salarié"),
			que.NewPText("Chomeur"),
		},
		Answer: [][]string{
			{"899", "253"},
			{"520", "50"},
		},
	},
	que.TextBlock{
		Parts:   "Soit $f(x) = &2x + 1&$. Calculer $f'$.",
		Bold:    true,
		Italic:  true,
		Smaller: true,
	},
	que.TreeFieldBlock{
		EventsProposals: []string{"P", "F", "?"},
		AnswerRoot: que.TreeNodeAnswer{
			Probabilities: []string{"1/5", "4/5"},
			Children: []que.TreeNodeAnswer{
				{
					Value: 0,
				},
				{
					Value: 1,
				},
			},
		},
	},
	que.VariationTableBlock{
		Label: "g(x)",
		Xs:    []string{"-5", "0", "2/3"},
		Fxs:   []string{"4.5", "2/9", "-12"},
	},
	que.VariationTableFieldBlock{
		Answer: que.VariationTableBlock{
			Label: "g(x)",
			Xs:    []string{"-5", "0", "2/3"},
			Fxs:   []string{"4.5", "2/9", "-12"},
		},
	},
	que.VectorFieldBlock{
		Answer: que.CoordExpression{
			X: "-2", Y: "3",
		},
		AcceptColinear: true,
		DisplayColumn:  true,
	},
}
