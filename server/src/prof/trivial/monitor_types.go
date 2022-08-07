package trivial

import (
	tv "github.com/benoitkugler/maths-online/trivial"
)

type gamePlayers struct {
	Player    string
	Successes tv.Success
}

type gameSummary struct {
	GameID        tv.RoomID
	CurrentPlayer string // empty when no one is playing
	Players       []gamePlayers
	RoomSize      int
}

type teacherSocketData struct {
	Games []gameSummary
}
