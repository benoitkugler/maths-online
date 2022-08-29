package trivial

import (
	"github.com/benoitkugler/maths-online/prof/teacher"
	tcAPI "github.com/benoitkugler/maths-online/prof/teacher"
	tr "github.com/benoitkugler/maths-online/sql/trivial"
	"github.com/benoitkugler/maths-online/trivial"
)

type TrivialExt struct {
	Config tr.Trivial
	Origin teacher.Origin

	NbQuestionsByCategories [trivial.NbCategories]int
}

func newTrivialExt(sel questionSelector, config tr.Trivial, userID, adminID uID) (TrivialExt, error) {
	vis := tcAPI.NewVisibility(config.IdTeacher, userID, adminID, config.Public)
	origin := tcAPI.Origin{
		AllowPublish: userID == adminID,
		Visibility:   vis,
		IsPublic:     config.Public,
	}

	out := TrivialExt{Config: config, Origin: origin}
	questions, err := sel.search(config.Questions, userID)
	if err != nil {
		return out, err
	}
	for i, cat := range questions {
		out.NbQuestionsByCategories[i] = len(cat.Questions)
	}
	return out, nil
}

type LaunchSessionIn struct {
	IdConfig tr.IdTrivial
	// Size of the groups as chosen by the teacher.
	Groups []int
}

type LaunchSessionOut struct {
	GameIDs []string
}
