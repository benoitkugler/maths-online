package events

import (
	"time"

	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
)

type Date time.Time

type Event struct {
	IdStudent teacher.IdStudent
	Event     EventK
	Date      Date
}

type EventK int16

const (
	E_IsyTriv_Create EventK = iota
	E_IsyTriv_Streak3
	E_IsyTriv_Win
	E_Homework_TaskDone
	E_Homework_TravailDone
	E_All_QuestionRight
	E_All_QuestionWrong
	E_Misc_SetPlaylist

	// these events are computed from the others and
	// not store in the DB, but are displayed like regular events
	E_ConnectStreak3
	E_ConnectStreak7
	E_ConnectStreak30
)
