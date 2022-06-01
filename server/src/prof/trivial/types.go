package trivial

import (
	"github.com/benoitkugler/maths-online/prof/editor"
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

func (tc TrivialConfig) withDetails(dict map[int64]editor.QuestionTags, origin teacher.Origin) TrivialConfigExt {
	out := TrivialConfigExt{Config: tc, Origin: origin}
	for i, cat := range tc.Questions {
		out.NbQuestionsByCategories[i] = len(cat.filter(dict))
	}
	return out
}

type LaunchSessionIn struct {
	IdConfig int64
	// Size of tje groups are fixed by the teacher.
	Groups []int
}

type LaunchSessionOut struct {
	GameIDs []string
}
