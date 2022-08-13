package trivial

import (
	"github.com/benoitkugler/maths-online/prof/teacher"
	"github.com/benoitkugler/maths-online/trivial"
)

// QuestionCriterion is an union of intersection of tags.
type QuestionCriterion [][]string

// CategoriesQuestions defines a union of intersection of tags,
// for every category.
type CategoriesQuestions [trivial.NbCategories]QuestionCriterion

type TrivialExt struct {
	Config Trivial
	Origin teacher.Origin

	NbQuestionsByCategories [trivial.NbCategories]int
}

func (tc Trivial) withDetails(db DB, origin teacher.Origin, userID uID) (TrivialExt, error) {
	out := TrivialExt{Config: tc, Origin: origin}
	for i, cat := range tc.Questions {
		questions, err := cat.selectQuestions(db, userID)
		if err != nil {
			return out, err
		}
		out.NbQuestionsByCategories[i] = len(questions)
	}
	return out, nil
}

type LaunchSessionIn struct {
	IdConfig IdTrivial
	// Size of tje groups are fixed by the teacher.
	Groups []int
}

type LaunchSessionOut struct {
	GameIDs []string
}
