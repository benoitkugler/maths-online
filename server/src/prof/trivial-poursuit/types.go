package trivialpoursuit

import (
	"encoding/json"
	"sort"

	"github.com/benoitkugler/maths-online/maths/exercice"
	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
)

//go:generate ../../../../../structgen/structgen -source=types.go -mode=itfs-json:gen_itfs.go

// QuestionCriterion is an union of intersection of tags.
type QuestionCriterion [][]string

// remove empty intersection and normalizes tags
func (qc QuestionCriterion) normalize() (out QuestionCriterion) {
	for _, q := range qc {
		for i, t := range q {
			q[i] = exercice.NormalizeTag(t)
		}

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

	Running                 LaunchSessionOut // empty when not running
	NbQuestionsByCategories [game.NbCategories]int
}

func (tc TrivialConfig) withDetails(dict map[int64]exercice.QuestionTags, sessions map[int64]LaunchSessionOut) TrivialConfigExt {
	out := TrivialConfigExt{Config: tc}
	out.Running = sessions[tc.Id]
	for i, cat := range tc.Questions {
		out.NbQuestionsByCategories[i] = len(cat.filter(dict))
	}
	return out
}

type LaunchSessionIn struct {
	IdConfig      int64
	GroupStrategy GroupStrategy
}

type LaunchSessionOut struct {
	SessionID         string
	GroupsID          []string // used for FixedGroupStrategy
	GroupStrategyKind int
}

func (ls *LaunchSessionIn) UnmarshalJSON(data []byte) error {
	var tmp struct {
		IdConfig      int64
		GroupStrategy GroupStrategyWrapper
	}
	err := json.Unmarshal(data, &tmp)
	ls.IdConfig = tmp.IdConfig
	ls.GroupStrategy = tmp.GroupStrategy.Data

	return err
}

var (
	_ GroupStrategy = RandomGroupStrategy{}
	_ GroupStrategy = FixedSizeGroupStrategy{}
)

// GroupStrategy defines how players are matched
type GroupStrategy interface {
	kind() int

	// initGames is called once when creating a new sesssion,
	initGames(*gameSession)

	// selectOrCreateGame return the game `studentID` should join,
	// or an empty string to start a new room with the returned number of players.
	// `studentID` is -1 for anonymous players
	// `selectOrCreateGame` must also validate `clientGameID` when needed
	selectOrCreateGame(clientGameID string, studentID int64, session *gameSession) (GameID, error)
}

// RandomGroupStrategy groups players at random,
// according to their connection time.
type RandomGroupStrategy struct {
	MaxPlayersPerGroup int
	TotalPlayersNumber int
}

// FixedSizeGroupStrategy is a broader version of `FixedGroupStrategy`,
// in which only the size of the groups are fixed by the teacher.
type FixedSizeGroupStrategy struct {
	Groups []int
}

// FixedGroupStrategy follow the groups precised by
// the teacher. Anonymous players are added randomly into existing groups.
type FixedGroupStrategy struct {
	Groups [][]int64 // student ids
}
