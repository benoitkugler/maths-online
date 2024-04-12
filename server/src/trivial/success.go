package trivial

import "github.com/benoitkugler/maths-online/server/src/sql/events"

// SuccessHandler reacts to students acheiving
// global success.
//
// It will typically be implemented with a [*sql.DB]
type SuccessHandler interface {
	OnQuestion(player PlayerID, correct, hasStreak3 bool) events.EventNotification
	OnWin(player PlayerID) events.EventNotification
}
