// Package examples provides a list of Question
// which may be used for demonstration or testing purposes.
package examples

import (
	"fmt"

	"github.com/benoitkugler/maths-online/maths/questions"
)

// Questions returns a version of the examples.
func Questions() (out []questions.QuestionInstance) {
	for _, qu := range questionsList {
		out = append(out, qu.Instantiate())
	}

	for _, block := range blockList {
		title := fmt.Sprintf("%T", block)
		enonce := questions.Enonce{
			questions.TextBlock{
				Parts: questions.Interpolated("Exemple du moule " + title),
				Bold:  true,
			},
			block,
		}
		qu := questions.QuestionPage{
			Title:  title,
			Enonce: enonce,
		}
		out = append(out, qu.Instantiate())
	}

	return out
}
