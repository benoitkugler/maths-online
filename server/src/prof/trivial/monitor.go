package trivial

import (
	"fmt"
	"sort"

	"github.com/benoitkugler/maths-online/prof/teacher"
	tv "github.com/benoitkugler/maths-online/trivial"
	"github.com/labstack/echo/v4"
)

type GamePlayers struct {
	Player    string
	Successes tv.Success
}

type GameSummary struct {
	GameID        tv.RoomID
	CurrentPlayer string // empty when no one is playing
	Players       []GamePlayers
	RoomSize      int
}

type MonitorOut struct {
	Games []GameSummary
}

func newGameSummary(s tv.Summary) (out GameSummary) {
	out.GameID = s.ID
	out.RoomSize = s.RoomSize
	for p, su := range s.Successes {
		out.Players = append(out.Players, GamePlayers{
			Player:    p.Pseudo,
			Successes: su,
		})
	}

	sort.Slice(out.Players, func(i, j int) bool { return out.Players[i].Player < out.Players[j].Player })

	if s.PlayerTurn != nil {
		out.CurrentPlayer = s.PlayerTurn.Pseudo
	}
	return out
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
func (gs *gameSession) collectSummaries() map[tv.RoomID]tv.Summary {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	out := make(map[tv.RoomID]tv.Summary)
	for _, ga := range gs.games {
		su := ga.Summary()
		out[su.ID] = su
	}

	return out
}

// TrivialTeacherMonitor returns the summary of games currently playing
func (ct *Controller) TrivialTeacherMonitor(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	session := ct.getSession(user.Id)
	if session == nil {
		return fmt.Errorf("internal error: no running session for %d", user.Id)
	}

	out := newMonitorOut(session.collectSummaries())

	return c.JSON(200, out)
}
