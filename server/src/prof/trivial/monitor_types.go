package trivial

import (
	tv "github.com/benoitkugler/maths-online/trivial"
)

//go:generate ../../../../../structgen/structgen -source=monitor_types.go -mode=ts:../../../../prof/src/controller/trivial_config_socket_gen.ts

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
