package exercice

import "github.com/benoitkugler/maths-online/maths/expression"

// randomParameters is the serialized form of expression.RandomParameters
type randomParameters []randomParameter

type randomParameter struct {
	Expression string `json:"expression"` // serialized form
	Variable   rune   `json:"variable"`
}

// toMap assumes `rp` only contains valid expressions
func (rp randomParameters) toMap() expression.RandomParameters {
	out := make(expression.RandomParameters, len(rp))
	for _, item := range rp {
		out[expression.Variable(item.Variable)], _ = expression.Parse(item.Expression)
	}
	return out
}

// Exercice is a sequence of questions
type ExerciceQuestions struct {
	Exercice
	Questions Questions
}
