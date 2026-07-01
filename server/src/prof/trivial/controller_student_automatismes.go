package trivial

import (
	"time"

	"github.com/benoitkugler/maths-online/server/src/pass"
	"github.com/benoitkugler/maths-online/server/src/sql/events"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/benoitkugler/maths-online/server/src/sql/trivial"
	tv "github.com/benoitkugler/maths-online/server/src/trivial"
	"github.com/benoitkugler/maths-online/server/src/utils"
	"github.com/labstack/echo/v4"
)

// ----------------------- Student API -----------------------

// StudentLaunchTrivialAutomatisme creates a game for the given config,
// and returns the public game code, which may be joined with the regular API.
func (ct *Controller) StudentLaunchTrivialAutomatisme(c echo.Context) error {
	idC := pass.EncryptedID(c.QueryParam("client-id"))
	var idStudent int64
	if idC != "" {
		var err error
		idStudent, err = ct.studentKey.DecryptID(idC)
		if err != nil {
			return err
		}
	}
	id, err := utils.QueryParamInt64(c, "trivial-id")
	if err != nil {
		return err
	}

	out, err := ct.launchTrivialAutomatisme(trivial.IdTrivial(id), teacher.IdStudent(idStudent))
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

// idStudent may be null
func (ct *Controller) launchTrivialAutomatisme(idTrivial trivial.IdTrivial, idStudent teacher.IdStudent) (LaunchSelfaccessOut, error) {
	config, err := trivial.SelectTrivial(ct.db, idTrivial)
	if err != nil {
		return LaunchSelfaccessOut{}, utils.SQLError(err)
	}

	// select the questions
	questionPool, err := selectQuestions(ct.db, config.Questions, config.IdTeacher, true)
	if err != nil {
		return LaunchSelfaccessOut{}, err
	}

	options := tv.Options{
		Launch:          tv.LaunchStrategy{Manual: true},
		QuestionTimeout: time.Second * time.Duration(config.QuestionTimeout),
		ShowDecrassage:  config.ShowDecrassage,
		Questions:       questionPool,
	}

	gameID := ct.store.newSelfaccessGameID()

	ProgressLogger.Printf("Creating game for config %d", config.Id)

	ct.store.createGame(createGame{ID: gameID, Options: options})

	// update the success, for registered accounts
	var notif events.EventNotification
	if idStudent != 0 {
		notif, err = events.RegisterEvents(ct.db, idStudent, events.E_IsyTriv_Create)
		if err != nil {
			return LaunchSelfaccessOut{}, err
		}
		notif.HideIfNoPoints()
	}

	return LaunchSelfaccessOut{GameID: tv.RoomID(gameID.String()), Notification: notif}, nil
}
