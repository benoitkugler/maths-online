package homework

import (
	"database/sql"
	"errors"
	"sort"
	"time"

	"github.com/benoitkugler/maths-online/server/src/pass"
	tcAPI "github.com/benoitkugler/maths-online/server/src/prof/teacher"
	ed "github.com/benoitkugler/maths-online/server/src/sql/editor"
	ho "github.com/benoitkugler/maths-online/server/src/sql/homework"
	"github.com/benoitkugler/maths-online/server/src/sql/tasks"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	taAPI "github.com/benoitkugler/maths-online/server/src/tasks"
	"github.com/benoitkugler/maths-online/server/src/utils"
	"github.com/labstack/echo/v4"
)

type uID = teacher.IdTeacher

var errAccessForbidden = errors.New("resource access forbidden")

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
	userID := tcAPI.JWTTeacher(c)

	out, err := ct.getSheets(userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

// Homeworks stores the [Travail]s and [Sheet]s available to
// one teacher
type Homeworks struct {
	Sheets  map[ho.IdSheet]SheetExt
	Travaux []ClassroomTravaux // one per classroom
}

func (ct *Controller) getSheets(userID uID) (out Homeworks, err error) {
	// load the classrooms
	classrooms, err := teacher.SelectClassroomsByIdTeachers(ct.db, userID)
	if err != nil {
		return out, utils.SQLError(err)
	}

	// load all the available [Sheets] ...
	sheetsDict, err := ho.SelectSheetsByIdTeachers(ct.db, userID)
	if err != nil {
		return out, utils.SQLError(err)
	}

	// .. and all the [Travail]s
	travauxDict, err := ho.SelectTravailsByIdClassrooms(ct.db, classrooms.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	loader, err := newSheetsLoader(ct.db, sheetsDict.IDs())
	if err != nil {
		return out, utils.SQLError(err)
	}

	// finally agregate the results
	out.Sheets = loader.buildSheetExts(sheetsDict)
	for _, class := range classrooms {
		out.Travaux = append(out.Travaux, newClassroomTravaux(class, travauxDict))
	}
	sort.Slice(out.Travaux, func(i, j int) bool { return out.Travaux[i].Classroom.Id < out.Travaux[j].Classroom.Id })

	return out, nil
}

func (ct *Controller) HomeworkCreateSheet(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	sheet, err := ct.createSheet(userID)
	if err != nil {
		return err
	}

	out := SheetExt{Sheet: sheet}
	return c.JSON(200, out)
}

func (ct *Controller) createSheet(userID uID) (ho.Sheet, error) {
	sheet, err := ho.Sheet{
		IdTeacher: userID,
		Title:     "Feuille d'exercices",
	}.Insert(ct.db)
	if err != nil {
		return ho.Sheet{}, utils.SQLError(err)
	}

	return sheet, nil
}

type CreateTravailIn struct {
	IdSheet     ho.IdSheet
	IdClassroom teacher.IdClassroom
}

// [HomeworkCreateTravail] creates a new [Travail] entry for the
// given classroom, with the given [Sheet]
func (ct *Controller) HomeworkCreateTravail(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args CreateTravailIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.assignSheetTo(args, userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) assignSheetTo(args CreateTravailIn, userID uID) (ho.Travail, error) {
	// check classroom owner
	classroom, err := teacher.SelectClassroom(ct.db, args.IdClassroom)
	if err != nil {
		return ho.Travail{}, utils.SQLError(err)
	}
	if classroom.IdTeacher != userID {
		return ho.Travail{}, errAccessForbidden
	}

	tr := ho.Travail{
		IdSheet:     args.IdSheet,
		IdClassroom: args.IdClassroom,
		Noted:       true,
		Deadline:    ho.Time(time.Now().Add(time.Hour * 24 * 14).Round(time.Hour)), // two weeks
	}
	tr, err = tr.Insert(ct.db)
	if err != nil {
		return ho.Travail{}, utils.SQLError(err)
	}
	return tr, nil
}

func (ct *Controller) checkSheetOwner(idSheet ho.IdSheet, userID uID) error {
	sheet, err := ho.SelectSheet(ct.db, idSheet)
	if err != nil {
		return utils.SQLError(err)
	}

	// check if the sheet is owned by the user
	if sheet.IdTeacher != userID {
		return errAccessForbidden
	}

	return nil
}

func (ct *Controller) checkTravailOwner(idTravail ho.IdTravail, userID uID) error {
	travail, err := ho.SelectTravail(ct.db, idTravail)
	if err != nil {
		return utils.SQLError(err)
	}

	classroom, err := teacher.SelectClassroom(ct.db, travail.IdClassroom)
	if err != nil {
		return utils.SQLError(err)
	}

	// check if the travail is owned by the user
	if classroom.IdTeacher != userID {
		return errAccessForbidden
	}

	return nil
}

func (ct *Controller) HomeworkUpdateSheet(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args ho.Sheet
	if err := c.Bind(&args); err != nil {
		return err
	}

	err := ct.updateSheet(args, userID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) updateSheet(sheet ho.Sheet, userID uID) error {
	if err := ct.checkSheetOwner(sheet.Id, userID); err != nil {
		return err
	}

	_, err := sheet.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}

func (ct *Controller) HomeworkUpdateTravail(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args ho.Travail
	if err := c.Bind(&args); err != nil {
		return err
	}

	err := ct.updateTravail(args, userID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) updateTravail(travail ho.Travail, userID uID) error {
	if err := ct.checkTravailOwner(travail.Id, userID); err != nil {
		return err
	}

	_, err := travail.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}

type AddExerciceToTaskIn struct {
	IdSheet    ho.IdSheet
	IdExercice ed.IdExercice
}

type AddMonoquestionToTaskIn struct {
	IdSheet    ho.IdSheet
	IdQuestion ed.IdQuestion
}

type AddRandomMonoquestionToTaskIn struct {
	IdSheet         ho.IdSheet
	IdQuestiongroup ed.IdQuestiongroup
}

func (ct *Controller) HomeworkAddExercice(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args AddExerciceToTaskIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	task, err := ct.addExerciceTo(args, userID)
	if err != nil {
		return err
	}

	return c.JSON(200, task)
}

func (ct *Controller) HomeworkAddMonoquestion(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args AddMonoquestionToTaskIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	task, err := ct.addMonoquestionTo(args, userID)
	if err != nil {
		return err
	}

	return c.JSON(200, task)
}

func (ct *Controller) HomeworkAddRandomMonoquestion(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args AddRandomMonoquestionToTaskIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	task, err := ct.addRandomMonoquestionTo(args, userID)
	if err != nil {
		return err
	}

	return c.JSON(200, task)
}

func (ct *Controller) addExerciceTo(args AddExerciceToTaskIn, userID uID) (TaskExt, error) {
	task := tasks.Task{IdExercice: args.IdExercice.AsOptional()}
	return ct.addTaskTo(args.IdSheet, task, userID)
}

// used defaut value of  Bareme: 1, NbRepeat: 3
func (ct *Controller) addMonoquestionTo(args AddMonoquestionToTaskIn, userID uID) (TaskExt, error) {
	// create the monoquestion, checking that the question has a group
	question, err := ed.SelectQuestion(ct.db, args.IdQuestion)
	if err != nil {
		return TaskExt{}, utils.SQLError(err)
	}

	if !question.IdGroup.Valid {
		return TaskExt{}, errors.New("internal error: (mono)question not included in a group")
	}

	mono, err := tasks.Monoquestion{IdQuestion: args.IdQuestion, Bareme: 1, NbRepeat: 3}.Insert(ct.db)
	if err != nil {
		return TaskExt{}, utils.SQLError(err)
	}
	task := tasks.Task{IdMonoquestion: mono.Id.AsOptional()}
	out, err := ct.addTaskTo(args.IdSheet, task, userID)
	if err != nil {
		// cleanup the monoquestion
		_, _ = tasks.DeleteMonoquestionById(ct.db, mono.Id)
		return TaskExt{}, err
	}

	return out, nil
}

// used defaut value of Bareme: 1, NbRepeat: 3
func (ct *Controller) addRandomMonoquestionTo(args AddRandomMonoquestionToTaskIn, userID uID) (TaskExt, error) {
	mono, err := tasks.RandomMonoquestion{IdQuestiongroup: args.IdQuestiongroup, Bareme: 1, NbRepeat: 3}.Insert(ct.db)
	if err != nil {
		return TaskExt{}, utils.SQLError(err)
	}
	task := tasks.Task{IdRandomMonoquestion: mono.Id.AsOptional()}
	out, err := ct.addTaskTo(args.IdSheet, task, userID)
	if err != nil {
		// cleanup the monoquestion
		_, _ = tasks.DeleteRandomMonoquestionById(ct.db, mono.Id)
		return TaskExt{}, err
	}

	return out, nil
}

// addTaskTo insert the given new task in the DB, and adds it to sheet
func (ct *Controller) addTaskTo(sheet ho.IdSheet, task tasks.Task, userID uID) (TaskExt, error) {
	if err := ct.checkSheetOwner(sheet, userID); err != nil {
		return TaskExt{}, err
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return TaskExt{}, utils.SQLError(err)
	}

	task, err = task.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return TaskExt{}, utils.SQLError(err)
	}

	// link the task to the sheet, appending
	links, err := ho.SelectSheetTasksByIdSheets(tx, sheet)
	if err != nil {
		_ = tx.Rollback()
		return TaskExt{}, utils.SQLError(err)
	}

	err = ho.InsertManySheetTasks(tx, ho.SheetTask{IdSheet: sheet, IdTask: task.Id, Index: len(links)})
	if err != nil {
		_ = tx.Rollback()
		return TaskExt{}, utils.SQLError(err)
	}

	err = tx.Commit()
	if err != nil {
		return TaskExt{}, utils.SQLError(err)
	}

	out, err := loadTaskExt(ct.db, task.Id)
	if err != nil {
		return TaskExt{}, err
	}

	return out, nil
}

// HomeworkRemoveTask deletes the given tasks, also removing
// all potential student progressions
func (ct *Controller) HomeworkRemoveTask(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id-task")
	if err != nil {
		return err
	}

	err = ct.removeTask(tasks.IdTask(id), userID)
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

	sheet, err := sheetFromTask(tx, idTask)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := ct.checkSheetOwner(sheet.Id, userID); err != nil {
		_ = tx.Rollback()
		return err
	}

	currentLinks, err := ho.SelectSheetTasksByIdSheets(tx, sheet.Id)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}
	currentLinks.EnsureOrder()

	var filtered []tasks.IdTask
	for _, link := range currentLinks {
		if link.IdTask != idTask {
			filtered = append(filtered, link.IdTask)
		}
	}

	// start by updating the links
	err = updateSheetTasksOrder(tx, sheet.Id, filtered)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	// and then cleanup the unused task
	removedTask, err := tasks.DeleteTaskById(tx, idTask)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	// delete the potential associated monoquestion
	if id := removedTask.IdMonoquestion; id.Valid {
		_, err = tasks.DeleteMonoquestionById(tx, id.ID)
		if err != nil {
			_ = tx.Rollback()
			return utils.SQLError(err)
		}
	} else if id := removedTask.IdRandomMonoquestion; id.Valid {
		_, err = tasks.DeleteRandomMonoquestionById(tx, id.ID)
		if err != nil {
			_ = tx.Rollback()
			return utils.SQLError(err)
		}

	}

	err = tx.Commit()
	if err != nil {
		return utils.SQLError(err)
	}
	return nil
}

func (ct *Controller) HomeworkUpdateMonoquestion(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args tasks.Monoquestion
	if err := c.Bind(&args); err != nil {
		return err
	}

	// check that the monoquestion is in a sheet owner by user
	idTask, idSheet, err := ho.LoadMonoquestionSheet(ct.db, args.Id)
	if err != nil {
		return err
	}
	err = ct.checkSheetOwner(idSheet, userID)
	if err != nil {
		return err
	}

	// only update bareme and repetitions
	mono, err := tasks.SelectMonoquestion(ct.db, args.Id)
	if err != nil {
		return utils.SQLError(err)
	}

	mono.Bareme = args.Bareme
	mono.NbRepeat = args.NbRepeat
	_, err = mono.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	// reload the task to properly update the UI
	out, err := loadTaskExt(ct.db, idTask)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) HomeworkUpdateRandomMonoquestion(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args tasks.RandomMonoquestion
	if err := c.Bind(&args); err != nil {
		return err
	}

	// check that the monoquestion is in a sheet owner by user
	idTask, idSheet, err := ho.LoadRandomMonoquestionSheet(ct.db, args.Id)
	if err != nil {
		return err
	}
	err = ct.checkSheetOwner(idSheet, userID)
	if err != nil {
		return err
	}

	// only update bareme and repetitions
	mono, err := tasks.SelectRandomMonoquestion(ct.db, args.Id)
	if err != nil {
		return utils.SQLError(err)
	}

	mono.Bareme = args.Bareme
	mono.NbRepeat = args.NbRepeat
	mono.Difficulty = args.Difficulty
	_, err = mono.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	// reload the task to properly update the UI
	out, err := loadTaskExt(ct.db, idTask)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

type ReorderSheetTasksIn struct {
	IdSheet ho.IdSheet
	Tasks   []tasks.IdTask
}

func (ct *Controller) HomeworkReorderSheetTasks(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args ReorderSheetTasksIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	err := ct.reorderSheetTasks(args, userID)
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
	userID := tcAPI.JWTTeacher(c)

	idSheet, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	err = ct.deleteSheet(ho.IdSheet(idSheet), userID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) deleteSheet(idSheet ho.IdSheet, userID uID) error {
	sheet, err := ho.SelectSheet(ct.db, idSheet)
	if err != nil {
		return utils.SQLError(err)
	}

	if sheet.IdTeacher != userID {
		return errAccessForbidden
	}

	// garbage collect the associated tasks :
	// the link table "sheet_tasks" is automatically cleaned up, but not the "tasks" table
	ts, err := ho.SelectSheetTasksByIdSheets(ct.db, idSheet)
	if err != nil {
		return utils.SQLError(err)
	}

	// we also need to remove the monoquestions and randommonoquestion associated
	tasksMap, err := tasks.SelectTasks(ct.db, ts.IdTasks()...)
	if err != nil {
		return utils.SQLError(err)
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}

	_, err = ho.DeleteSheetById(tx, idSheet)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	_, err = tasks.DeleteTasksByIDs(tx, ts.IdTasks()...)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	// delete the potential associated monoquestion
	for _, removedTask := range tasksMap {
		if id := removedTask.IdMonoquestion; id.Valid {
			_, err = tasks.DeleteMonoquestionById(tx, id.ID)
			if err != nil {
				_ = tx.Rollback()
				return utils.SQLError(err)
			}
		} else if id := removedTask.IdRandomMonoquestion; id.Valid {
			_, err = tasks.DeleteRandomMonoquestionById(tx, id.ID)
			if err != nil {
				_ = tx.Rollback()
				return utils.SQLError(err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}

func (ct *Controller) HomeworkDeleteTravail(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	idSheet, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	err = ct.deleteTravail(ho.IdTravail(idSheet), userID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

// remove the travail entry, but not the sheet neither the progressions
func (ct *Controller) deleteTravail(id ho.IdTravail, userID uID) error {
	_, err := ho.DeleteTravailById(ct.db, id)
	if err != nil {
		return utils.SQLError(err)
	}
	return nil
}

type CopySheetIn struct {
	IdSheet ho.IdSheet
}

func (ct *Controller) HomeworkCopySheet(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args CopySheetIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.duplicateSheet(args, userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) duplicateSheet(args CopySheetIn, userID uID) (SheetExt, error) {
	if err := ct.checkSheetOwner(args.IdSheet, userID); err != nil {
		return SheetExt{}, err
	}

	sheet, err := ho.SelectSheet(ct.db, args.IdSheet)
	if err != nil {
		return SheetExt{}, utils.SQLError(err)
	}

	links, err := ho.SelectSheetTasksByIdSheets(ct.db, sheet.Id)
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

	newSheet, err := sheet.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return SheetExt{}, utils.SQLError(err)
	}

	// create new tasks : a task can't be be shared
	newLinks := make(ho.SheetTasks, len(links))
	for i, link := range links {
		task := taskMap[link.IdTask]
		newTask := task
		// for monoquestion, also copy the monoquestion
		if id := task.IdMonoquestion; id.Valid {
			monoquestion, err := tasks.SelectMonoquestion(tx, id.ID)
			if err != nil {
				_ = tx.Rollback()
				return SheetExt{}, utils.SQLError(err)
			}
			monoquestion, err = monoquestion.Insert(tx)
			if err != nil {
				_ = tx.Rollback()
				return SheetExt{}, utils.SQLError(err)
			}
			newTask.IdMonoquestion = monoquestion.Id.AsOptional()
		} else if id := task.IdRandomMonoquestion; id.Valid {
			monoquestion, err := tasks.SelectRandomMonoquestion(tx, id.ID)
			if err != nil {
				_ = tx.Rollback()
				return SheetExt{}, utils.SQLError(err)
			}
			monoquestion, err = monoquestion.Insert(tx)
			if err != nil {
				_ = tx.Rollback()
				return SheetExt{}, utils.SQLError(err)
			}
			newTask.IdRandomMonoquestion = monoquestion.Id.AsOptional()
		}

		newTask, err = newTask.Insert(tx)
		if err != nil {
			_ = tx.Rollback()
			return SheetExt{}, utils.SQLError(err)
		}
		newLinks[i] = ho.SheetTask{IdSheet: newSheet.Id, IdTask: newTask.Id, Index: i}
	}

	err = ho.InsertManySheetTasks(tx, newLinks...)
	if err != nil {
		_ = tx.Rollback()
		return SheetExt{}, utils.SQLError(err)
	}

	loader, err := newSheetsLoader(tx, []ho.IdSheet{newSheet.Id})
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

type CopyTravailIn struct {
	IdTravail   ho.IdTravail
	IdClassroom teacher.IdClassroom
}

func (ct *Controller) HomeworkCopyTravail(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args CopyTravailIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.copyTravailTo(args, userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) copyTravailTo(args CopyTravailIn, userID uID) (ho.Travail, error) {
	cl, err := teacher.SelectClassroom(ct.db, args.IdClassroom)
	if err != nil {
		return ho.Travail{}, utils.SQLError(err)
	}

	if cl.IdTeacher != userID {
		return ho.Travail{}, errAccessForbidden
	}

	travail, err := ho.SelectTravail(ct.db, args.IdTravail)
	if err != nil {
		return ho.Travail{}, utils.SQLError(err)
	}

	// shallow copy is enough
	travail.IdClassroom = args.IdClassroom
	travail, err = travail.Insert(ct.db)
	if err != nil {
		return ho.Travail{}, utils.SQLError(err)
	}

	return travail, nil
}

type HowemorkMarksIn struct {
	IdClassroom teacher.IdClassroom
	IdTravaux   []ho.IdTravail
}

type HomeworkMarksOut struct {
	Students []tcAPI.StudentHeader                          // the students of the classroom
	Marks    map[ho.IdTravail]map[teacher.IdStudent]float64 // the notes for each travail and student, /20
}

func (ct *Controller) HomeworkGetMarks(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args HowemorkMarksIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.getMarks(args, userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) getMarks(args HowemorkMarksIn, userID uID) (HomeworkMarksOut, error) {
	classroom, err := teacher.SelectClassroom(ct.db, args.IdClassroom)
	if err != nil {
		return HomeworkMarksOut{}, utils.SQLError(err)
	}
	if classroom.IdTeacher != userID {
		return HomeworkMarksOut{}, errAccessForbidden
	}

	students, err := tcAPI.LoadClassroomStudents(ct.db, classroom.Id)
	if err != nil {
		return HomeworkMarksOut{}, err
	}
	travaux, err := ho.SelectTravails(ct.db, args.IdTravaux...)
	if err != nil {
		return HomeworkMarksOut{}, utils.SQLError(err)
	}

	out := HomeworkMarksOut{
		Students: make([]tcAPI.StudentHeader, len(students)),
		Marks:    make(map[ho.IdTravail]map[teacher.IdStudent]float64),
	}
	// student list
	for i, s := range students {
		out.Students[i] = tcAPI.NewStudentHeader(s)
	}
	// compute the sheets marks :
	loader, err := newSheetsLoader(ct.db, travaux.IdSheets())
	if err != nil {
		return HomeworkMarksOut{}, err
	}
	// load all the progressions : for each task and student
	progressions, err := loader.tasks.LoadProgressions(ct.db)
	if err != nil {
		return HomeworkMarksOut{}, err
	}

	for id, travail := range travaux {
		if travail.IdClassroom != classroom.Id {
			return HomeworkMarksOut{}, errors.New("internal error: inconsitent classroom ID")
		}

		markByStudent := make(map[teacher.IdStudent]float64)
		var sheetTotal int
		// for each student, get its progression for each task
		tasks := loader.tasksForSheet(travail.IdSheet)
		for _, link := range tasks {
			work := loader.tasks.GetWork(loader.tasks.Tasks[link.IdTask])
			bareme := work.Bareme()
			taskTotal := bareme.Total()
			sheetTotal += taskTotal
			byStudent := progressions[link.IdTask]
			// add each progression to the student note
			for idStudent, prog := range byStudent {
				studentMark := bareme.ComputeMark(prog.Questions)
				markByStudent[idStudent] = markByStudent[idStudent] + float64(studentMark)
			}
		}
		// normalize the mark / 20
		for id, mark := range markByStudent {
			markByStudent[id] = 20 * mark / float64(sheetTotal)
		}
		out.Marks[id] = markByStudent
	}

	return out, nil
}

// Student API

// StudentGetTravaux returns the sheets for the given student
// Only the mandatory, noted one are returned.
func (ct *Controller) StudentGetTravaux(c echo.Context) error {
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

	travaux, err := ho.SelectTravailsByIdClassrooms(ct.db, student.IdClassroom)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	sheets, err := ho.SelectSheets(ct.db, travaux.IdSheets()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	links1, err := ho.SelectSheetTasksByIdSheets(ct.db, sheets.IDs()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	sheetToTasks := links1.ByIdSheet()

	// collect the student progressions
	progMap, err := taAPI.LoadTasksProgression(ct.db, idStudent, links1.IdTasks())
	if err != nil {
		return nil, utils.SQLError(err)
	}

	for _, travail := range travaux {
		if !travail.Noted { // ignore free sheets
			continue
		}

		sheet := sheets[travail.IdSheet]
		tasksForSheet := sheetToTasks[sheet.Id] // defined exercices
		taskList := make([]taAPI.TaskProgressionHeader, len(tasksForSheet))
		for i, exLink := range tasksForSheet {
			taskList[i] = progMap[exLink.IdTask]
		}
		out = append(out, SheetProgression{
			IdTravail: travail.Id,
			Sheet: Sheet{
				Id:       sheet.Id,
				Title:    sheet.Title,
				Deadline: travail.Deadline,
			},
			Tasks: taskList,
		})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Sheet.Id < out[j].Sheet.Id })

	return out, nil
}

func (ct *Controller) StudentInstantiateTask(c echo.Context) error {
	idCrypted := pass.EncryptedID(c.QueryParam("client-id"))
	idStudent, err := ct.studentKey.DecryptID(idCrypted)
	if err != nil {
		return err
	}

	idTask, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	task, err := tasks.SelectTask(ct.db, tasks.IdTask(idTask))
	if err != nil {
		return utils.SQLError(err)
	}

	out, err := taAPI.InstantiateWork(ct.db, taAPI.NewWorkID(task), teacher.IdStudent(idStudent))
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

// StudentEvaluateTask calls ed.EvaluteExercice and registers
// the student progression, returning the update mark.
// However, if the sheet is expired, it does not register the progression.
func (ct *Controller) StudentEvaluateTask(c echo.Context) error {
	var args StudentEvaluateTaskIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.studentEvaluateTask(args)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) studentEvaluateTask(args StudentEvaluateTaskIn) (StudentEvaluateTaskOut, error) {
	idStudent, err := ct.studentKey.DecryptID(args.StudentID)
	if err != nil {
		return StudentEvaluateTaskOut{}, err
	}

	registerProgression := true

	// to preserve compatibily, we accept empty [travail] field
	if args.IdTravail != 0 {
		travail, err := ho.SelectTravail(ct.db, args.IdTravail)
		if err != nil {
			return StudentEvaluateTaskOut{}, utils.SQLError(err)
		}
		// Always register progression for free travail
		registerProgression = !travail.Noted || !travail.IsExpired()
	}

	ex, mark, err := taAPI.EvaluateTaskExercice(ct.db, args.IdTask, teacher.IdStudent(idStudent), args.Ex, registerProgression)
	if err != nil {
		return StudentEvaluateTaskOut{}, err
	}
	return StudentEvaluateTaskOut{Ex: ex, Mark: mark}, nil
}
