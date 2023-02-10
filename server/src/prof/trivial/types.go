package trivial

import (
	"fmt"

	tcAPI "github.com/benoitkugler/maths-online/server/src/prof/teacher"
	tr "github.com/benoitkugler/maths-online/server/src/sql/trivial"
	"github.com/benoitkugler/maths-online/server/src/trivial"
)

// LoadQuestionNumbers returns the number of questions available for
// each categories, as defined by [config.Questions]
func LoadQuestionNumbers(db tr.DB, config tr.Trivial, userID uID) (out [trivial.NbCategories]int, err error) {
	qus, err := selectQuestions(db, config.Questions, userID)
	if err != nil {
		return out, err
	}
	for i, cat := range qus {
		out[i] = len(cat.Questions)
	}
	return out, nil
}

type TrivialExt struct {
	Config tr.Trivial
	Origin tcAPI.Origin

	NbQuestionsByCategories [trivial.NbCategories]int
}

func newTrivialExt(sel questionSelector, config tr.Trivial, inReview tcAPI.OptionalIdReview, userID, adminID uID) (TrivialExt, error) {
	vis := tcAPI.NewVisibility(config.IdTeacher, userID, adminID, config.Public)
	origin := tcAPI.Origin{
		AllowPublish: userID == adminID,
		Visibility:   vis,
		IsPublic:     config.Public,
		IsInReview:   inReview,
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

	Groups GroupsStrategy
}

type LaunchSessionOut struct {
	GameIDs []trivial.RoomID
}

type GroupsStrategy interface {
	// return an error if the sizes are invalid
	groups() ([]trivial.LaunchStrategy, error)
}

// GroupsStrategyManual does not fix the number of players
type GroupsStrategyManual struct {
	NbGroups int
}

func (gsm GroupsStrategyManual) groups() ([]trivial.LaunchStrategy, error) {
	out := make([]trivial.LaunchStrategy, gsm.NbGroups)
	for i := range out {
		out[i] = trivial.LaunchStrategy{Manual: true}
	}
	return out, nil
}

// GroupsStrategyAuto starts the game when full
type GroupsStrategyAuto struct {
	Groups []int
}

func (gsa GroupsStrategyAuto) groups() ([]trivial.LaunchStrategy, error) {
	out := make([]trivial.LaunchStrategy, len(gsa.Groups))
	for i, max := range gsa.Groups {
		if max <= 0 {
			return nil, fmt.Errorf("internal error: invalid room size %d", max)
		}
		out[i] = trivial.LaunchStrategy{Manual: false, Max: max}
	}
	return out, nil
}
