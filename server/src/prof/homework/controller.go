package homework

import (
	"database/sql"
	"errors"
	"sort"
	"time"

	"github.com/benoitkugler/maths-online/prof/editor"
	"github.com/benoitkugler/maths-online/prof/teacher"
	"github.com/benoitkugler/maths-online/utils"
	"github.com/labstack/echo/v4"
)

type uID = teacher.IdTeacher

var accessForbidden = errors.New("resource access forbidden")

type Controller struct {
	db    *sql.DB
	admin teacher.Teacher
}

func NewController(db *sql.DB, admin teacher.Teacher) *Controller {
	return &Controller{
		db:    db,
		admin: admin,
	}
}

func (ct *Controller) HomeworkGetSheets(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	out, err := ct.getSheets(user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) getSheets(userID uID) (out []ClassroomSheets, err error) {
	// load the classrooms
	classrooms, err := teacher.SelectClassroomsByIdTeachers(ct.db, userID)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	// load all the sheets required
	sheetsDict, err := SelectSheetsByIdClassrooms(ct.db, classrooms.IDs()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	// load all the exercices required
	links, err := SelectSheetExercicesByIdSheets(ct.db, sheetsDict.IDs()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	exes, err := editor.SelectExercices(ct.db, links.IdExercices()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	// finally agregate the results
	sheets := buildSheetExts(sheetsDict, links, exes)
	for _, class := range classrooms {
		out = append(out, newClassroomSheets(class, sheets))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Classroom.Id < out[j].Classroom.Id })
	return out, nil
}

type CreateSheetIn struct {
	IdClassroom teacher.IdClassroom
}

func (ct *Controller) HomeworkCreateSheet(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	var args CreateSheetIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	sheet, err := ct.createSheet(args.IdClassroom, user.Id)
	if err != nil {
		return err
	}

	out := SheetExt{Sheet: sheet}
	return c.JSON(200, out)
}

func (ct *Controller) createSheet(idClassroom teacher.IdClassroom, userID uID) (Sheet, error) {
	class, err := teacher.SelectClassroom(ct.db, idClassroom)
	if err != nil {
		return Sheet{}, utils.SQLError(err)
	}

	if class.IdTeacher != userID {
		return Sheet{}, accessForbidden
	}

	sheet, err := Sheet{
		IdClassroom: class.Id,
		Title:       "Feuille d'exercices",
		Notation:    SuccessNotation,
		Activated:   false,
		Deadline:    Time(time.Now().Add(time.Hour * 24 * 24)),
	}.Insert(ct.db)
	if err != nil {
		return Sheet{}, utils.SQLError(err)
	}

	return sheet, nil
}

type UpdateSheetIn struct {
	Sheet     Sheet
	Exercices []editor.IdExercice
}

func (ct *Controller) HomeworkUpdateSheet(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	var args UpdateSheetIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	err := ct.updateSheet(args, user.Id)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) updateSheet(args UpdateSheetIn, userID uID) error {
	// check the classroom is owned by the user
	class, err := teacher.SelectClassroom(ct.db, args.Sheet.IdClassroom)
	if err != nil {
		return utils.SQLError(err)
	}

	if class.IdTeacher != userID {
		return accessForbidden
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}

	_, err = args.Sheet.Update(tx)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	// update the links
	err = updateSheetExercices(tx, args.Sheet.Id, args.Exercices)
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

func (ct *Controller) HomeworkDeleteSheet(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	idSheet, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	err = ct.deleteSheet(IdSheet(idSheet), user.Id)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) deleteSheet(idSheet IdSheet, userID uID) error {
	sheet, err := SelectSheet(ct.db, idSheet)
	if err != nil {
		return utils.SQLError(err)
	}

	cl, err := teacher.SelectClassroom(ct.db, sheet.IdClassroom)
	if err != nil {
		return utils.SQLError(err)
	}

	if cl.IdTeacher != userID {
		return accessForbidden
	}

	_, err = DeleteSheetById(ct.db, idSheet)
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}
