package trivial

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/benoitkugler/maths-online/server/src/pass"
	tcAPI "github.com/benoitkugler/maths-online/server/src/prof/teacher"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/benoitkugler/maths-online/server/src/sql/trivial"
	tv "github.com/benoitkugler/maths-online/server/src/trivial"
	"github.com/benoitkugler/maths-online/server/src/utils"
	"github.com/labstack/echo/v4"
)

// Student self access

type TrivialSelfaccess struct {
	Classrooms []teacher.Classroom
	Actives    []teacher.IdClassroom
}

// [TrivialGetSelfaccess] returns the classroom which may launch
// the given trivial
func (ct *Controller) TrivialGetSelfaccess(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id-trivial")
	if err != nil {
		return err
	}

	out, err := ct.selfaccess(trivial.IdTrivial(id), userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) selfaccess(id trivial.IdTrivial, userID uID) (TrivialSelfaccess, error) {
	// load all classrooms for the user
	cls_, err := teacher.SelectClassroomsByIdTeachers(ct.db, userID)
	if err != nil {
		return TrivialSelfaccess{}, utils.SQLError(err)
	}
	cls := make([]teacher.Classroom, 0, len(cls_))
	for _, cl := range cls_ {
		cls = append(cls, cl)
	}
	sort.Slice(cls, func(i, j int) bool { return cls[i].Name < cls[j].Name })

	links, err := trivial.SelectSelfaccessTrivialsByIdTrivials(ct.db, id)
	if err != nil {
		return TrivialSelfaccess{}, utils.SQLError(err)
	}

	actives := links.ByIdTeacher()[userID].IdClassrooms()
	return TrivialSelfaccess{Classrooms: cls, Actives: actives}, nil
}

type UpdateSelfaccessIn struct {
	IdTrivial    trivial.IdTrivial
	IdClassrooms []teacher.IdClassroom
}

// [TrivialUpdateSelfaccess] sets the classroom which may launch
// the given trivial
func (ct *Controller) TrivialUpdateSelfaccess(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var in UpdateSelfaccessIn
	if err := c.Bind(&in); err != nil {
		return fmt.Errorf("invalid parameters format: %s", err)
	}

	err := ct.updateSelfaccess(in, userID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) updateSelfaccess(args UpdateSelfaccessIn, userID uID) error {
	tx, err := ct.db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}

	_, err = trivial.DeleteSelfaccessTrivialsByIdTrivialAndIdTeacher(tx, args.IdTrivial, userID)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	newItems := make(trivial.SelfaccessTrivials, len(args.IdClassrooms))
	for i, id := range args.IdClassrooms {
		newItems[i] = trivial.SelfaccessTrivial{IdClassroom: id, IdTrivial: args.IdTrivial, IdTeacher: userID}
	}
	err = trivial.InsertManySelfaccessTrivials(tx, newItems...)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	err = tx.Commit()
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}

// ----------------------- Student API -----------------------

func (ct *Controller) StudentGetSelfaccess(c echo.Context) error {
	idC := pass.EncryptedID(c.QueryParam("client-id"))
	idStudent, err := ct.studentKey.DecryptID(idC)
	if err != nil {
		return err
	}
	out, err := ct.studentGetSelfaccess(teacher.IdStudent(idStudent))
	if err != nil {
		return err
	}
	return c.JSON(200, GetSelfaccessOut{out})
}

func (ct *Controller) studentGetSelfaccess(idStudent teacher.IdStudent) ([]trivial.Trivial, error) {
	student, err := teacher.SelectStudent(ct.db, idStudent)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	links, err := trivial.SelectSelfaccessTrivialsByIdClassrooms(ct.db, student.IdClassroom)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	trivsM, err := trivial.SelectTrivials(ct.db, links.IdTrivials()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	trivs := make([]trivial.Trivial, 0, len(trivsM))
	for _, triv := range trivsM {
		trivs = append(trivs, triv)
	}
	sort.Slice(trivs, func(i, j int) bool { return trivs[i].Name < trivs[j].Name })

	return trivs, nil
}

// StudentLaunchSelfaccess creates a game for the given config,
// and returns the public game code, which may be joined with the regular API.
func (ct *Controller) StudentLaunchSelfaccess(c echo.Context) error {
	idC := pass.EncryptedID(c.QueryParam("client-id"))
	idStudent, err := ct.studentKey.DecryptID(idC)
	if err != nil {
		return err
	}
	id, err := utils.QueryParamInt64(c, "trivial-id")
	if err != nil {
		return err
	}

	out, err := ct.launchSelfaccess(trivial.IdTrivial(id), teacher.IdStudent(idStudent))
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) launchSelfaccess(idTrivial trivial.IdTrivial, idStudent teacher.IdStudent) (LaunchSelfaccessOut, error) {
	student, err := teacher.SelectStudent(ct.db, idStudent)
	if err != nil {
		return LaunchSelfaccessOut{}, utils.SQLError(err)
	}
	classroom, err := teacher.SelectClassroom(ct.db, student.IdClassroom)
	if err != nil {
		return LaunchSelfaccessOut{}, utils.SQLError(err)
	}

	config, err := trivial.SelectTrivial(ct.db, idTrivial)
	if err != nil {
		return LaunchSelfaccessOut{}, utils.SQLError(err)
	}

	userID := classroom.IdTeacher
	// admin config may be launched, since it is a readonly operation
	if config.IdTeacher != userID && config.IdTeacher != ct.admin.Id {
		return LaunchSelfaccessOut{}, errAccessForbidden
	}

	// select the questions
	questionPool, err := selectQuestions(ct.db, config.Questions, userID)
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
	ct.store.createGame(createGame{ID: gameID, Options: options})

	return LaunchSelfaccessOut{gameID.roomID()}, nil
}

// StartSelfaccess starts a game previously created by [StudentLaunchSelfaccess]
// TODO: secure this access
func (ct *Controller) StudentStartSelfaccess(c echo.Context) error {
	gameID := c.QueryParam("game-id")
	parsed, err := ct.store.parseCode(gameID)
	if err != nil {
		return err
	}
	selfID, ok := parsed.(selfaccessCode)
	if !ok {
		return errors.New("internal error: invalid game code")
	}

	err = ct.store.startGame(selfID)
	if err != nil {
		return err
	}

	return nil
}
