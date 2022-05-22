package trivialpoursuit

import (
	"encoding/json"

	"github.com/benoitkugler/maths-online/prof/editor"
	"github.com/benoitkugler/maths-online/prof/teacher"
	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
)

//go:generate ../../../../../structgen/structgen -source=types.go -mode=itfs-json:gen_itfs.go

// QuestionCriterion is an union of intersection of tags.
type QuestionCriterion [][]string

// CategoriesQuestions defines a union of intersection of tags,
// for every category.
type CategoriesQuestions [game.NbCategories]QuestionCriterion

type TrivialConfigExt struct {
	Config     TrivialConfig
	Visibility teacher.Visibility

	Running                 LaunchSessionOut // empty when not running
	NbQuestionsByCategories [game.NbCategories]int
}

func (tc TrivialConfig) withDetails(dict map[int64]editor.QuestionTags, sessions map[int64]LaunchSessionOut) TrivialConfigExt {
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
