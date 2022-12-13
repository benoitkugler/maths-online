// Package examples provides a list of Question
// which may be used for demonstration or testing purposes.
package examples

import (
	"fmt"

	que "github.com/benoitkugler/maths-online/server/src/maths/questions"
)

type LabeledQuestion struct {
	Title    string
	Question que.QuestionInstance
}

// Questions returns a version of the examples.
func Questions() (out []LabeledQuestion) {
	out = append(out, LabeledQuestion{
		Title: "Remplir un tableau de variation",
		Question: que.QuestionPage{
			Enonce: que.Enonce{
				que.TextBlock{
					Parts: "Réponse attendue : -2, 0, 2/3 \n 4/9, -2, 2/3",
				},
				que.VariationTableFieldBlock{
					Answer: que.VariationTableBlock{
						Label: "y = h(x)",
						Xs:    []string{"-2", "0", "2/3"},
						Fxs:   []string{"4/9", "-2", "2/3"},
					},
				},
			},
		}.Instantiate(),
	},
	)

	for _, block := range blockList {
		title := fmt.Sprintf("%T", block)
		enonce := que.Enonce{
			que.TextBlock{
				Parts: que.Interpolated("Exemple du moule " + title),
				Bold:  true,
			},
			block,
		}
		qu := que.QuestionPage{Enonce: enonce}.Instantiate()
		out = append(out, LabeledQuestion{
			Title:    title,
			Question: qu,
		})
	}

	return out
}
