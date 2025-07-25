package events

import (
	"time"

	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
)

type Date time.Time

func NewDate(year int, month time.Month, day int) Date {
	return Date(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
}

func NewDateFrom(ti time.Time) Date {
	return NewDate(ti.Year(), ti.Month(), ti.Day())
}

func (d Date) Time() time.Time {
	ti := time.Time(d)
	return time.Date(ti.Year(), ti.Month(), ti.Day(), 0, 0, 0, 0, time.UTC)
}

type Event struct {
	IdStudent teacher.IdStudent `gomacro-sql-on-delete:"CASCADE"`
	Event     EventK
	Date      Date
}

type EventK int16

const (
	E_IsyTriv_Create       EventK = iota // Créer une partie d'IsyTriv
	E_IsyTriv_Streak3                    // Réussir trois questions IsyTriv d'affilée
	E_IsyTriv_Win                        // Remporter une partie IsyTriv
	E_Homework_TaskDone                  // Terminer un exercice
	E_Homework_TravailDone               // Terminer une feuille d'exercices
	E_All_QuestionRight                  // Répondre correctement à une question
	E_All_QuestionWrong                  // Répondre incorrectement à une question
	E_Misc_SetPlaylist                   // Modifier sa playlist

	// these events are computed from the others and
	// not store in the DB, but are displayed like regular events
	E_ConnectStreak3  // Se connecter 3 jours de suite
	E_ConnectStreak7  // Se connecter 7 jours de suite
	E_ConnectStreak30 // Se connecter 30 jours de suite
)

const NbEvents = 11
