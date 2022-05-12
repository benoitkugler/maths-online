// Package examples provides a list of Question
// which may be used for demonstration or testing purposes.
package examples

import (
	"fmt"

	"github.com/benoitkugler/maths-online/maths/exercice"
)

// Questions returns a version of the examples.
func Questions() (out []exercice.QuestionInstance) {
	for _, block := range blockList {
		title := fmt.Sprintf("%T", block)
		enonce := exercice.Enonce{
			exercice.TextBlock{
				Parts: exercice.Interpolated("Exemple du moule " + title),
				Bold:  true,
			},
			block,
		}
		qu := exercice.Question{
			Title:  title,
			Enonce: enonce,
		}
		out = append(out, qu.Instantiate())
	}
	return out
}
