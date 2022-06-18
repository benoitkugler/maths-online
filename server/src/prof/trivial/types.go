package trivial

import (
	"github.com/benoitkugler/maths-online/prof/teacher"
	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
)

// QuestionCriterion is an union of intersection of tags.
type QuestionCriterion [][]string

// CategoriesQuestions defines a union of intersection of tags,
// for every category.
type CategoriesQuestions [game.NbCategories]QuestionCriterion

type TrivialConfigExt struct {
	Config TrivialConfig
	Origin teacher.Origin

	NbQuestionsByCategories [game.NbCategories]int
}

func (tc TrivialConfig) withDetails(db DB, origin teacher.Origin, userID int64) (TrivialConfigExt, error) {
	out := TrivialConfigExt{Config: tc, Origin: origin}
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
	IdConfig int64
	// Size of tje groups are fixed by the teacher.
	Groups []int
}

type LaunchSessionOut struct {
	GameIDs []string
}
