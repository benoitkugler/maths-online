package trivial

import (
	"github.com/benoitkugler/maths-online/server/src/sql/events"
	"github.com/benoitkugler/maths-online/server/src/sql/trivial"
	tv "github.com/benoitkugler/maths-online/server/src/trivial"
)

// file used to generate Dart types

type GetSelfaccessOut struct {
	Trivials []trivial.Trivial
}

type LaunchSelfaccessOut struct {
	GameID       tv.RoomID
	Notification events.EventNotification // new in version 1.7
}
