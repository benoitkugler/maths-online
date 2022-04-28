package trivialpoursuit

import (
	"sort"

	tv "github.com/benoitkugler/maths-online/trivial-poursuit"
	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
)

//go:generate ../../../../../structgen/structgen -source=teacher_types.go -mode=ts:../../../../prof/src/controller/trivial_config_socket_gen.ts

type gamePlayers struct {
	Player    string
	Successes game.Success
}

type gameSummary struct {
	GameID        GameID
	CurrentPlayer string // empty when no one is playing
	Players       []gamePlayers
}

func newGameSummary(s tv.GameSummary) (out gameSummary) {
	out.GameID = s.ID
	for p, su := range s.Successes {
		out.Players = append(out.Players, gamePlayers{
			Player:    p.Name,
			Successes: su,
		})
	}

	sort.Slice(out.Players, func(i, j int) bool { return out.Players[i].Player < out.Players[j].Player })

	if s.PlayerTurn != nil {
		out.CurrentPlayer = s.PlayerTurn.Name
	}
	return out
}

type teacherSocketData struct {
	Games []gameSummary
}
