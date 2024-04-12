package trivial

import (
	"sort"

	"github.com/benoitkugler/maths-online/server/src/prof/teacher"
	"github.com/benoitkugler/maths-online/server/src/tasks"
	tv "github.com/benoitkugler/maths-online/server/src/trivial"
	"github.com/labstack/echo/v4"
)

type GamePlayers struct {
	Player    string
	Successes tv.Success
}

type GameSummary struct {
	GameID             tv.RoomID
	CurrentPlayer      string          // empty when no one is playing
	LatestQuestion     QuestionContent // empty before the first question
	Players            []GamePlayers
	RoomSize           tv.RoomSize
	InQuestionStudents []string
}

func newGameSummary(s tv.Summary) (out GameSummary) {
	out.GameID = s.ID
	out.RoomSize = s.RoomSize
	out.LatestQuestion = QuestionContent{
		Id:        s.LatestQuestion.ID,
		Categorie: s.LatestQuestion.Categorie,
		Question:  s.LatestQuestion.Question.ToClient(),
		Params:    tasks.NewParams(s.LatestQuestion.Vars),
	}
	out.InQuestionStudents = s.InQuestionStudents

	for p, su := range s.Successes {
		out.Players = append(out.Players, GamePlayers{
			Player:    p,
			Successes: su,
		})
	}

	sort.Slice(out.Players, func(i, j int) bool { return out.Players[i].Player < out.Players[j].Player })

	if s.PlayerTurn != nil {
		out.CurrentPlayer = s.PlayerTurn.Pseudo
	}
	return out
}

type MonitorOut struct {
	Games []GameSummary
}

func newMonitorOut(summaries map[tv.RoomID]tv.Summary) (out MonitorOut) {
	for _, su := range summaries {
		out.Games = append(out.Games, newGameSummary(su))
	}

	sort.Slice(out.Games, func(i, j int) bool {
		return out.Games[i].GameID < out.Games[j].GameID
	})

	return out
}

// lock and fetch summaries
func (gs *gameStore) collectSummaries(sessionID sessionID) map[tv.RoomID]tv.Summary {
	gs.lock.Lock()
	games := gs.getSession(sessionID)
	gs.lock.Unlock()

	out := make(map[tv.RoomID]tv.Summary)
	for _, ga := range games {
		su := ga.Summary()
		out[su.ID] = su
	}

	return out
}

// TrivialTeacherMonitor returns the summary of games currently playing
func (ct *Controller) TrivialTeacherMonitor(c echo.Context) error {
	userID := teacher.JWTTeacher(c)

	session := ct.store.getSessionID(userID)
	if session == "" {
		return c.JSON(200, MonitorOut{})
	}

	out := newMonitorOut(ct.store.collectSummaries(session))

	return c.JSON(200, out)
}
