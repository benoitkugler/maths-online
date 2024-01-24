package ceintures

import (
	"github.com/benoitkugler/maths-online/server/src/maths/questions"
	tc "github.com/benoitkugler/maths-online/server/src/sql/teacher"
)

type IdBeltquestion int64

// Beltevolution is the evolution of one student in a belt scheme.
//
// gomacro:SQL ADD UNIQUE(IdStudent)
type Beltevolution struct {
	IdStudent tc.IdStudent
	Level     Level
	Advance   Advance
	Stats     Stats
}

// Beltquestion is one question, contained in
// a domain and color.
//
// gomacro:SQL _SELECT KEY(Domain, Rank)
type Beltquestion struct {
	Id IdBeltquestion

	Domain Domain // Location
	Rank   Rank   // Location

	Parameters questions.Parameters
	Enonce     questions.Enonce
	// Correction an optional content describing the expected solution,
	// to be instantiated with the same parameters as [Enonce]
	Correction questions.Enonce
}

func (qu Beltquestion) Page() questions.QuestionPage {
	return questions.QuestionPage{Enonce: qu.Enonce, Correction: qu.Correction, Parameters: qu.Parameters}
}
