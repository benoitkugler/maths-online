package questions

import (
	"testing"

	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func TestExportLatex(t *testing.T) {
	enonce := Enonce{
		TextBlock{Parts: "Consigne de la première question : que vaut N ?", Bold: true},
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
		OrderedListFieldBlock{Label: ``, Answer: []Interpolated{"A", "$x+2$", "B"}, AdditionalProposals: []Interpolated{"C"}},
		OrderedListFieldBlock{Label: `$x \in $`, Answer: []Interpolated{"A", "$x+2$", "B"}, AdditionalProposals: []Interpolated{"C"}},
	}

	qu, err := enonce.InstantiateWith(nil)
	tu.AssertNoErr(t, err)

	latexQu := qu.ToLatex()

	tu.GenerateLatex(t, latexQu, "export-question.pdf")

	ques := []QuestionInstance{qu, qu, qu}

	latexEx := InstancesToLatex(ques)
	tu.GenerateLatex(t, latexEx, "export-exercice.pdf")
}
