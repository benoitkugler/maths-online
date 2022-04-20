package trivialpoursuit

import (
	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
)

// CategoriesQuestions defines a union of intersection of tags,
// for every category.
type CategoriesQuestions [game.NbCategories][][]string

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
