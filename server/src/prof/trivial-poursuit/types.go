package trivialpoursuit

import (
	"sort"

	"github.com/benoitkugler/maths-online/maths/exercice"
	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
)

// QuestionCriterion is an union of intersection of tags.
type QuestionCriterion [][]string

// remove empty intersection
func (qc QuestionCriterion) normalize() (out QuestionCriterion) {
	for _, q := range qc {
		if len(q) != 0 {
			out = append(out, q)
		}
	}
	return out
}

func (qc QuestionCriterion) filter(dict map[int64]exercice.QuestionTags) (out IDs) {
	qc = qc.normalize()

	// an empty criterion is interpreted as an invalid criterion,
	// since it is never something you want in practice (at least the class level should be specified)
	if len(qc) == 0 {
		return nil
	}

	for idQuestion, questions := range dict {
		questionTags := questions.Crible()
		for _, union := range qc { // at least one union must match
			if questionTags.HasAll(union) {
				out = append(out, idQuestion)
				break // no need to check the other unions
			}
		}
	}

	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] }) // deterministic order

	return out
}

// CategoriesQuestions defines a union of intersection of tags,
// for every category.
type CategoriesQuestions [game.NbCategories]QuestionCriterion

type TrivialConfigExt struct {
	Config TrivialConfig

	NbQuestionsByCategories [game.NbCategories]int
}

func (tc TrivialConfig) withQuestionsNumber(dict map[int64]exercice.QuestionTags) TrivialConfigExt {
	out := TrivialConfigExt{Config: tc}
	for i, cat := range tc.Questions {
		out.NbQuestionsByCategories[i] = len(cat.filter(dict))
	}
	return out
}

type LaunchSessionIn struct {
	IdConfig      int64
	GroupStrategy GroupStrategy
}

var _ GroupStrategy = RandomGroupStrategy{}

// GroupStrategy defines how players are matched
type GroupStrategy interface {
	// initGames is called once when creating a new sesssion
	initGames(*gameSession)

	// selectGame return the game `player` should join,
	// or an empty string to start a new room with the returned number of players
	// `studentID` is -1 for anonymous players
	selectGame(studentID int64, session *gameSession) (GameID, int)
}

// RandomGroupStrategy groups players at random,
// according to their connection time.
type RandomGroupStrategy struct {
	MaxPlayersPerGroup int
	TotalPlayersNumber int
}

// FixedGroupStrategy follow the groups precised by
// the teacher. Anonymous players are added randomly into existing groups.
type FixedGroupStrategy struct {
	Groups [][]int64 // student ids
}
