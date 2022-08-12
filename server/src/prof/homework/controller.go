package homework

import (
	"database/sql"
	"errors"
	"sort"
	"time"

	"github.com/benoitkugler/maths-online/pass"
	ed "github.com/benoitkugler/maths-online/prof/editor"
	"github.com/benoitkugler/maths-online/prof/teacher"
	"github.com/benoitkugler/maths-online/tasks"
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

	loader, err := newSheetLoader(ct.db, sheetsDict.IDs(), userID, ct.admin.Id)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	// finally agregate the results
	sheets := loader.buildSheetExts(sheetsDict)
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

func (ct *Controller) checkSheetOwner(idSheet IdSheet, userID uID) error {
	sheet, err := SelectSheet(ct.db, idSheet)
	if err != nil {
		return utils.SQLError(err)
	}

	// check the classroom is owned by the user
	class, err := teacher.SelectClassroom(ct.db, sheet.IdClassroom)
	if err != nil {
		return utils.SQLError(err)
	}

	if class.IdTeacher != userID {
		return accessForbidden
	}

	return nil
}

func (ct *Controller) HomeworkUpdateSheet(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	var args Sheet
	if err := c.Bind(&args); err != nil {
		return err
	}

	err := ct.updateSheet(args, user.Id)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) updateSheet(sheet Sheet, userID uID) error {
	if err := ct.checkSheetOwner(sheet.Id, userID); err != nil {
		return err
	}

	_, err := sheet.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}

type AddTaskIn struct {
	IdSheet    IdSheet
	IdExercice ed.IdExercice
}

func (ct *Controller) HomeworkAddTask(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	var args AddTaskIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	task, err := ct.addTaskTo(args, user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, task)
}

func (ct *Controller) addTaskTo(args AddTaskIn, userID uID) (tasks.Task, error) {
	if err := ct.checkSheetOwner(args.IdSheet, userID); err != nil {
		return tasks.Task{}, err
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return tasks.Task{}, utils.SQLError(err)
	}

	task, err := tasks.Task{IdExercice: args.IdExercice}.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return tasks.Task{}, utils.SQLError(err)
	}

	// link the task to the sheet, appending
	links, err := SelectSheetTasksByIdSheets(tx, args.IdSheet)
	if err != nil {
		_ = tx.Rollback()
		return tasks.Task{}, utils.SQLError(err)
	}

	err = InsertManySheetTasks(tx, SheetTask{IdSheet: args.IdSheet, IdTask: task.Id, Index: len(links)})
	if err != nil {
		_ = tx.Rollback()
		return tasks.Task{}, utils.SQLError(err)
	}

	err = tx.Commit()
	if err != nil {
		return tasks.Task{}, utils.SQLError(err)
	}

	return task, nil
}

// HomeworkRemoveTask deletes the given tasks, also removing
// all potential student progressions
func (ct *Controller) HomeworkRemoveTask(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id-task")
	if err != nil {
		return err
	}

	err = ct.removeTask(tasks.IdTask(id), user.Id)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) removeTask(idTask tasks.IdTask, userID uID) error {
	tx, err := ct.db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}

	link, found, err := SelectSheetTaskByIdTask(tx, idTask)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	if !found {
		_ = tx.Rollback()
		return errors.New("internal error: task without sheet")
	}

	if err := ct.checkSheetOwner(link.IdSheet, userID); err != nil {
		_ = tx.Rollback()
		return err
	}

	currentLinks, err := SelectSheetTasksByIdSheets(tx, link.IdSheet)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}
	currentLinks.ensureOrder()

	var filtered []tasks.IdTask
	for _, link := range currentLinks {
		if link.IdTask != idTask {
			filtered = append(filtered, link.IdTask)
		}
	}

	// start by updating the links
	err = updateSheetTasksOrder(tx, link.IdSheet, filtered)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	// and then cleanup the unused task
	_, err = tasks.DeleteTaskById(tx, idTask)
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

type ReorderSheetTasksIn struct {
	IdSheet IdSheet
	Tasks   []tasks.IdTask
}

func (ct *Controller) HomeworkReorderSheetTasks(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	var args ReorderSheetTasksIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	err := ct.reorderSheetTasks(args, user.Id)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) reorderSheetTasks(args ReorderSheetTasksIn, userID uID) error {
	if err := ct.checkSheetOwner(args.IdSheet, userID); err != nil {
		return err
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}
	err = updateSheetTasksOrder(tx, args.IdSheet, args.Tasks)
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

	links, err := SelectSheetTasksByIdSheets(ct.db, sheet.Id)
	if err != nil {
		return SheetExt{}, utils.SQLError(err)
	}

	taskMap, err := tasks.SelectTasks(ct.db, links.IdTasks()...)
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

	// create new tasks : a task can't be be shared
	newLinks := make(SheetTasks, len(links))
	for i, link := range links {
		ta, err := tasks.Task{IdExercice: taskMap[link.IdTask].IdExercice}.Insert(tx)
		if err != nil {
			_ = tx.Rollback()
			return SheetExt{}, utils.SQLError(err)
		}
		newLinks[i] = SheetTask{IdSheet: newSheet.Id, IdTask: ta.Id, Index: i}
	}

	err = InsertManySheetTasks(tx, newLinks...)
	if err != nil {
		_ = tx.Rollback()
		return SheetExt{}, utils.SQLError(err)
	}

	loader, err := newSheetLoader(tx, []IdSheet{newSheet.Id}, userID, ct.admin.Id)
	if err != nil {
		_ = tx.Rollback()
		return SheetExt{}, utils.SQLError(err)
	}

	err = tx.Commit()
	if err != nil {
		return SheetExt{}, utils.SQLError(err)
	}

	out := loader.newSheetExt(newSheet)

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

	links1, err := SelectSheetTasksByIdSheets(ct.db, sheets.IDs()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	sheetToTasks := links1.ByIdSheet()

	// collect the student progressions
	progMap, err := tasks.LoadTasksProgression(ct.db, idStudent, links1.IdTasks())
	if err != nil {
		return nil, utils.SQLError(err)
	}

	for _, sheet := range sheets {
		if !sheet.Activated { // ignore hidden sheets
			continue
		}

		tasksForSheet := sheetToTasks[sheet.Id] // defined exercices
		taskList := make([]tasks.TaskProgressionHeader, len(tasksForSheet))
		for i, exLink := range tasksForSheet {
			taskList[i] = progMap[exLink.IdTask]
		}
		out = append(out, SheetProgression{
			Sheet: sheet,
			Tasks: taskList,
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

	ex, mark, err := tasks.EvaluateTaskExercice(ct.db, args.IdTask, teacher.IdStudent(idStudent), args.Ex)
	if err != nil {
		return err
	}
	out := StudentEvaluateExerciceOut{Ex: ex, Mark: mark}

	return c.JSON(200, out)
}
