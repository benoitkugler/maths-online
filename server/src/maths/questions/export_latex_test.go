package questions

import (
	"testing"

	"github.com/benoitkugler/maths-online/server/src/maths/repere"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func tt(s string) TextPart {
	return TextPart{Content: s, Kind: Text}
}

func tm(s string) TextPart {
	return TextPart{Content: s, Kind: StaticMath}
}

func TestExportLatex(t *testing.T) {
	dummyCoord := CoordExpression{X: "1", Y: "1"}
	enonce := Enonce{
		TextBlock{Parts: "Consigne de la première question : que vaut N ?\n\n\n", Bold: true},
		TextBlock{Parts: "Conseil : considérer $f(x) = x - 8$", Italic: true, Smaller: true},
		TextBlock{Parts: "Compléter : N = "},
		NumberFieldBlock{"9"},
		FormulaBlock{`f(x) = \sqrt{x + 8}`},
		TextBlock{Parts: "Enoncer le lemme des poignées de main"},
		ExpressionFieldBlock{Label: "A = ", Expression: "2x + 7"},
		TextBlock{Parts: "Expression sans label"},
		ExpressionFieldBlock{Expression: "2x + 7"},
		// should be on one line
		TextBlock{Parts: "A = ("},
		NumberFieldBlock{"1"},
		TextBlock{Parts: ";"},
		NumberFieldBlock{"2"},
		TextBlock{Parts: ")"},
		// new line
		TextBlock{Parts: "\nQCM", Bold: true},
		RadioFieldBlock{Answer: "1", Proposals: []Interpolated{
			"Réponse 1",
			"$x + 4y = 5$",
			"Réponse D !",
			"Double line \n HAHAH",
		}, AsDropDown: false},
		TextBlock{Parts: "La suite est :"},
		RadioFieldBlock{Answer: "1", Proposals: []Interpolated{
			"Réponse 1",
			"$x + 4y = 5$",
			"Réponse D !",
			"Double line \n HAHAH",
		}, AsDropDown: true},
		TextBlock{Parts: "Quel est la négation de f croissante ?"},
		OrderedListFieldBlock{Label: ``, Answer: []Interpolated{"A", "$x+2$", "B"}, AdditionalProposals: []Interpolated{"C", `$ \forall x \in \R,$`}},
		OrderedListFieldBlock{Label: `$x \in $`, Answer: []Interpolated{"A", "$x+2$", "B"}, AdditionalProposals: []Interpolated{"C"}},
		VectorFieldBlock{DisplayColumn: true, Answer: dummyCoord},
		VectorFieldBlock{DisplayColumn: false, Answer: dummyCoord},
		// tables
		TableBlock{ // no headers
			Values: [][]TextPart{
				{tt("skjkdj"), tm("2x + 8"), tt("AA")},
				{tt("skjkdj"), tm("2x + 8"), tt("AA")},
				{tt("skjkdj"), tm("2x + 8"), tt("AA")},
				{tt("skjkdj"), tm("2x + 8"), tt("AA")},
			},
		},
		TableBlock{ // horizontal and vertical headers
			HorizontalHeaders: []TextPart{tt("H1"), tm("H2"), tt("H3")},
			VerticalHeaders:   []TextPart{tt("V1"), tm("V2")},
			Values: [][]TextPart{
				{tt("skjkdj"), tm("2x + 8"), tt("AA")},
				{tt("skjkdj"), tm("2x + 8"), tt("AA")},
			},
		},
		FigureBlock{
			ShowGrid:   true,
			ShowOrigin: true,
			Bounds:     repere.RepereBounds{Width: 50, Height: 30, Origin: repere.Coord{X: 10, Y: 5}},
			Drawings: repere.RandomDrawings{
				Points: []repere.NamedRandomLabeledPoint{
					{"A", repere.RandomLabeledPoint{"#FF0000", repere.RandomCoord{"1", "2"}, repere.Bottom}},
					{"B_xy", repere.RandomLabeledPoint{"#FF00FF", repere.RandomCoord{"1", "8"}, repere.Top}},
					{"C_1", repere.RandomLabeledPoint{"#FF00FF", repere.RandomCoord{"10", "8"}, repere.Top}},
				},
				Segments: []repere.RandomSegment{
					{"", "A", "B_xy", "#FF0FFF", 0, repere.SKSegment},
					{"D", "A", "C_1", "#F00FFF", 0, repere.SKLine},
					{"", "B_xy", "C_1", "#6F00FFFF", 0, repere.SKVector},
				},
				Lines: []repere.RandomLine{
					{"C_f", "0.1", "1", "#F0FF1020"},
					{"C_g", "inf", "4", "#F0FF1020"},
				},
				Circles: []repere.RandomCircle{
					{repere.RandomCoord{"2", "4"}, "3", "#00FF00", "#8800FFFF", "C"},
				},
				Areas: []repere.RandomArea{
					{Color: "#905F1820", Points: []string{"A", "B_xy", "C_1"}},
				},
			},
		},
		FigureBlock{
			ShowGrid:   true,
			ShowOrigin: false,
			Bounds:     repere.RepereBounds{Width: 30, Height: 30, Origin: repere.Coord{X: 10, Y: 5}},
		},
	}

	qu, err := enonce.InstantiateWith(nil)
	tu.AssertNoErr(t, err)

	latexQu := qu.ToLatex()

	tu.GenerateLatex(t, latexQu, "export-question.pdf")

	ques := []QuestionInstance{qu, qu, qu}

	latexEx := InstancesToLatex(ques)
	tu.GenerateLatex(t, latexEx, "export-exercice.pdf")
}
