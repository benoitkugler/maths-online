package homework

import (
	"database/sql"
	"errors"
	"sort"
	"time"

	"github.com/benoitkugler/maths-online/pass"
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

	studentKey pass.Encrypter
}

func NewController(db *sql.DB, admin teacher.Teacher, studentKey pass.Encrypter) *Controller {
	return &Controller{
		db:         db,
		admin:      admin,
		studentKey: studentKey,
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

	// load all the exercices required...
	links, err := SelectSheetExercicesByIdSheets(ct.db, sheetsDict.IDs()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	exes, err := editor.SelectExercices(ct.db, links.IdExercices()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	// ... and their questions
	questions, err := editor.SelectExerciceQuestionsByIdExercices(ct.db, exes.IDs()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	exerciceMap := editor.BuildExerciceHeaders(userID, ct.admin.Id, exes, questions)

	// finally agregate the results
	sheets := buildSheetExts(sheetsDict, links, exerciceMap)
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
		Deadline:    Time(time.Now().Add(time.Hour * 24 * 14).Round(time.Hour)), // two weeks
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

type CopySheetIn struct {
	IdSheet     IdSheet
	IdClassroom teacher.IdClassroom
}

func (ct *Controller) HomeworkCopySheet(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	var args CopySheetIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.copySheetTo(args, user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) copySheetTo(args CopySheetIn, userID uID) (SheetExt, error) {
	cl, err := teacher.SelectClassroom(ct.db, args.IdClassroom)
	if err != nil {
		return SheetExt{}, utils.SQLError(err)
	}

	if cl.IdTeacher != userID {
		return SheetExt{}, accessForbidden
	}

	sheet, err := SelectSheet(ct.db, args.IdSheet)
	if err != nil {
		return SheetExt{}, utils.SQLError(err)
	}

	links, err := SelectSheetExercicesByIdSheets(ct.db, sheet.Id)
	if err != nil {
		return SheetExt{}, utils.SQLError(err)
	}

	exes, err := editor.SelectExercices(ct.db, links.IdExercices()...)
	if err != nil {
		return SheetExt{}, utils.SQLError(err)
	}

	questions, err := editor.SelectExerciceQuestionsByIdExercices(ct.db, exes.IDs()...)
	if err != nil {
		return SheetExt{}, utils.SQLError(err)
	}

	// shallow copy of the item ...
	tx, err := ct.db.Begin()
	if err != nil {
		return SheetExt{}, utils.SQLError(err)
	}

	sheet.IdClassroom = args.IdClassroom
	newSheet, err := sheet.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return SheetExt{}, utils.SQLError(err)
	}

	// and copy of the links
	for i := range links {
		links[i].IdSheet = newSheet.Id
	}
	err = InsertManySheetExercices(tx, links...)
	if err != nil {
		_ = tx.Rollback()
		return SheetExt{}, utils.SQLError(err)
	}

	err = tx.Commit()
	if err != nil {
		return SheetExt{}, utils.SQLError(err)
	}

	exerciceMap := editor.BuildExerciceHeaders(userID, ct.admin.Id, exes, questions)
	out := newSheetExt(newSheet, links, exerciceMap)

	return out, nil
}

// Student API

// StudentGetSheets returns the sheet for the given student
func (ct *Controller) StudentGetSheets(c echo.Context) error {
	idCrypted := pass.EncryptedID(c.QueryParam("client-id"))

	idStudent, err := ct.studentKey.DecryptID(idCrypted)
	if err != nil {
		return err
	}

	out, err := ct.getStudentSheets(teacher.IdStudent(idStudent))
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) getStudentSheets(idStudent teacher.IdStudent) (out StudentSheets, err error) {
	student, err := teacher.SelectStudent(ct.db, idStudent)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	sheets, err := SelectSheetsByIdClassrooms(ct.db, student.IdClassroom)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	links1, err := SelectSheetExercicesByIdSheets(ct.db, sheets.IDs()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	sheetToExercices := links1.ByIdSheet()

	studentProgs, err := SelectStudentProgressionsByIdStudents(ct.db, teacher.IdStudent(idStudent))
	if err != nil {
		return nil, utils.SQLError(err)
	}
	// since we only load data for one student, we map by IdSheet and Index
	sheetExToProg := bySheetAndIndex(studentProgs)

	// collect the student progressions
	idProgressions := editor.NewIdProgressionSetFrom(studentProgs.IdProgressions())
	progMap, err := editor.LoadProgressions(ct.db, idProgressions)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	exercices, err := editor.SelectExercices(ct.db, studentProgs.IdExercices()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	for _, sheet := range sheets {
		if !sheet.Activated { // ignore hidden sheets
			continue
		}

		exerciceForSheet := sheetToExercices[sheet.Id] // defined exercices
		exWithProg := make([]ExerciceProgressionHeader, len(exerciceForSheet))
		for i, ex := range exerciceForSheet {
			// select the right progression, which may be empty
			// before the student starts the exercice
			prog, hasProg := sheetExToProg[sheetAndIndex{IdSheet: sheet.Id, Index: ex.Index}]
			exWithProg[i] = ExerciceProgressionHeader{
				Exercice:       exercices[ex.IdExercice],
				HasProgression: hasProg,
				Progression:    progMap[prog],
			}
		}
		out = append(out, SheetProgression{
			Sheet:     sheet,
			Exercices: exWithProg,
		})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Sheet.Id < out[j].Sheet.Id })

	return out, nil
}
