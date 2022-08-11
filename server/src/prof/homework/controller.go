package homework

import (
	"database/sql"
	"errors"
	"sort"
	"time"

	"github.com/benoitkugler/maths-online/pass"
	ed "github.com/benoitkugler/maths-online/prof/editor"
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

	exes, err := ed.SelectExercices(ct.db, links.IdExercices()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	// ... and their questions
	questions, err := ed.SelectExerciceQuestionsByIdExercices(ct.db, exes.IDs()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	exerciceMap := ed.BuildExerciceHeaders(userID, ct.admin.Id, exes, questions)

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
	Exercices []ed.IdExercice
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

	exes, err := ed.SelectExercices(ct.db, links.IdExercices()...)
	if err != nil {
		return SheetExt{}, utils.SQLError(err)
	}

	questions, err := ed.SelectExerciceQuestionsByIdExercices(ct.db, exes.IDs()...)
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

	exerciceMap := ed.BuildExerciceHeaders(userID, ct.admin.Id, exes, questions)
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
	sheetExToProg := studentProgs.bySheetAndIndex()

	// collect the student progressions
	idProgressions := ed.NewIdProgressionSetFrom(studentProgs.IdProgressions())
	progMap, err := ed.LoadProgressions(ct.db, idProgressions)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	exercices, err := ed.SelectExercices(ct.db, studentProgs.IdExercices()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	links, err := ed.SelectExerciceQuestionsByIdExercices(ct.db, exercices.IDs()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	questionsMap := links.ByIdExercice()

	for _, sheet := range sheets {
		if !sheet.Activated { // ignore hidden sheets
			continue
		}

		exerciceForSheet := sheetToExercices[sheet.Id] // defined exercices
		exWithProg := make([]ExerciceProgressionHeader, len(exerciceForSheet))
		for i, exLink := range exerciceForSheet {
			// select the right progression, which may be empty
			// before the student starts the exercice
			idProg, hasProg := sheetExToProg[sheetAndIndex{IdSheet: sheet.Id, Index: exLink.Index}]

			exercice := exercices[exLink.IdExercice]
			questions := questionsMap[exLink.IdExercice]
			progression := progMap[idProg]

			exWithProg[i] = ExerciceProgressionHeader{
				Exercice:       exercice,
				HasProgression: hasProg,
				Progression:    progression,
				Bareme:         questions.Bareme(),
				Mark:           mark(questions, progression),
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

func (ct *Controller) StudentInstantiateExercice(c echo.Context) error {
	idE, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	out, err := ed.InstantiateExercice(ct.db, ed.IdExercice(idE))
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

// StudentEvaluateExercice calls ed.EvaluteExercice and registers
// the student progression, returning the update mark.
func (ct *Controller) StudentEvaluateExercice(c echo.Context) error {
	var args StudentEvaluateExerciceIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	idStudent, err := ct.studentKey.DecryptID(args.StudentID)
	if err != nil {
		return err
	}

	out, err := ct.evaluateExercice(args.IdSheet, args.Index, teacher.IdStudent(idStudent), args.Ex)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) evaluateExercice(idSheet IdSheet, index int, idStudent teacher.IdStudent, ex ed.EvaluateExerciceIn) (StudentEvaluateExerciceOut, error) {
	out, err := ed.EvaluateExercice(ct.db, ex)
	if err != nil {
		return StudentEvaluateExerciceOut{}, err
	}

	// persists the progression on DB ...
	prog, err := ct.loadOrCreateProgressionFor(idSheet, index, idStudent)
	if err != nil {
		return StudentEvaluateExerciceOut{}, err
	}

	// ... update the progression questions
	err = ed.UpdateProgression(ct.db, prog, out.Progression.Questions)
	if err != nil {
		return StudentEvaluateExerciceOut{}, err
	}

	// compute the new mark
	questions, err := ed.SelectExerciceQuestionsByIdExercices(ct.db, ex.IdExercice)
	if err != nil {
		return StudentEvaluateExerciceOut{}, utils.SQLError(err)
	}
	m := mark(questions, out.Progression)

	return StudentEvaluateExerciceOut{Ex: out, Mark: m}, nil
}

// in the first try, the progression does not exists : create it
func (ct *Controller) loadOrCreateProgressionFor(idSheet IdSheet, index int, idStudent teacher.IdStudent) (ed.Progression, error) {
	links, err := SelectStudentProgressionsByIdStudents(ct.db, idStudent)
	if err != nil {
		return ed.Progression{}, utils.SQLError(err)
	}
	forThisSheet := links.ByIdSheet()[idSheet].bySheetAndIndex()
	idProg, hasProg := forThisSheet[sheetAndIndex{idSheet, index}]
	if !hasProg { // create an entry
		return ct.createProgressionFor(idSheet, index, idStudent)
	}

	// load the existing one
	prog, err := ed.SelectProgression(ct.db, idProg)
	if err != nil {
		return ed.Progression{}, utils.SQLError(err)
	}

	return prog, nil
}

func (ct *Controller) createProgressionFor(idSheet IdSheet, index int, idStudent teacher.IdStudent) (ed.Progression, error) {
	tx, err := ct.db.Begin()
	if err != nil {
		return ed.Progression{}, utils.SQLError(err)
	}
	links, err := SelectSheetExercicesByIdSheets(tx, idSheet)
	if err != nil {
		_ = tx.Rollback()
		return ed.Progression{}, utils.SQLError(err)
	}
	idExercice := links.bySheetAndIndex()[sheetAndIndex{idSheet, index}]

	prog, err := ed.Progression{IdExercice: idExercice}.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return ed.Progression{}, utils.SQLError(err)
	}

	err = InsertManyStudentProgressions(tx, StudentProgression{
		IdStudent:     idStudent,
		IdSheet:       idSheet,
		Index:         index,
		IdExercice:    idExercice,
		IdProgression: prog.Id,
	})
	if err != nil {
		_ = tx.Rollback()
		return ed.Progression{}, utils.SQLError(err)
	}

	err = tx.Commit()
	if err != nil {
		return ed.Progression{}, utils.SQLError(err)
	}

	return prog, nil
}
