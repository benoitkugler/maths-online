package homework

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"sort"
	"time"

	"github.com/benoitkugler/maths-online/server/src/pass"
	tcAPI "github.com/benoitkugler/maths-online/server/src/prof/teacher"
	"github.com/benoitkugler/maths-online/server/src/sql/editor"
	ed "github.com/benoitkugler/maths-online/server/src/sql/editor"
	ho "github.com/benoitkugler/maths-online/server/src/sql/homework"
	"github.com/benoitkugler/maths-online/server/src/sql/tasks"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	taAPI "github.com/benoitkugler/maths-online/server/src/tasks"
	"github.com/benoitkugler/maths-online/server/src/utils"
	"github.com/labstack/echo/v4"
)

var Logger = log.New(os.Stdout, "homework:INFO:", 0)

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

	mat_ := c.QueryParam("matiere")
	out, err := ct.getSheets(userID, teacher.MatiereTag(mat_))
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

func (ct *Controller) getSheets(userID uID, matiere teacher.MatiereTag) (out Homeworks, err error) {
	// load the classrooms
	classrooms, err := teacher.SelectClassroomsByIdTeachers(ct.db, userID)
	if err != nil {
		return out, utils.SQLError(err)
	}

	// load all the available [Sheets] (including admin) ...
	sheetsDict, err := ho.SelectSheetsByIdTeachers(ct.db, userID, ct.admin.Id)
	if err != nil {
		return out, utils.SQLError(err)
	}
	sheetsDict.RestrictVisible(userID)

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
	tmp := loader.buildSheetExts(sheetsDict, userID, ct.admin.Id)

	// only keep the sheets needed by travaux and the ones with the given topic
	out.Sheets = make(map[ho.IdSheet]SheetExt)
	for _, travail := range travauxDict {
		out.Sheets[travail.IdSheet] = tmp[travail.IdSheet]
	}
	for _, sheet := range tmp {
		if sheet.Sheet.Matiere == matiere {
			out.Sheets[sheet.Sheet.Id] = sheet
		}
	}

	for _, class := range classrooms {
		out.Travaux = append(out.Travaux, newClassroomTravaux(class, travauxDict))
	}

	// sort by classrooms
	sort.Slice(out.Travaux, func(i, j int) bool { return out.Travaux[i].Classroom.Id < out.Travaux[j].Classroom.Id })

	return out, nil
}

func (ct *Controller) HomeworkCreateSheet(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	sheet, err := ct.createSheet(userID)
	if err != nil {
		return err
	}

	out := SheetExt{Sheet: sheet, Origin: tcAPI.Origin{Visibility: tcAPI.Personnal}}
	return c.JSON(200, out)
}

func (ct *Controller) createSheet(userID uID) (ho.Sheet, error) {
	user, err := teacher.SelectTeacher(ct.db, userID)
	if err != nil {
		return ho.Sheet{}, utils.SQLError(err)
	}

	sheet, err := ho.Sheet{
		IdTeacher: userID,
		Title:     "Feuille d'exercices",
		Matiere:   user.FavoriteMatiere,
	}.Insert(ct.db)
	if err != nil {
		return ho.Sheet{}, utils.SQLError(err)
	}

	return sheet, nil
}

type CreateTravailWithIn struct {
	IdSheet     ho.IdSheet
	IdClassroom teacher.IdClassroom
}

// [HomeworkCreateTravailWith] creates a new [Travail] entry for the
// given classroom, with the given [Sheet]
func (ct *Controller) HomeworkCreateTravailWith(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args CreateTravailWithIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.assignSheetTo(args, userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) assignSheetTo(args CreateTravailWithIn, userID uID) (ho.Travail, error) {
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
		ShowAfter:   ho.Time(time.Now().Round(10 * time.Minute)),
		Deadline:    ho.Time(time.Now().Add(time.Hour * 7 * 14).Round(time.Hour)), // one week
	}
	tr, err = tr.Insert(ct.db)
	if err != nil {
		return ho.Travail{}, utils.SQLError(err)
	}
	return tr, nil
}

// [HomeworkCreateTravail] creates a [Travail] for the given classroom,
// linked to an anonymous [Sheet]
func (ct *Controller) HomeworkCreateTravail(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	id_, err := utils.QueryParamInt64(c, "id-classroom")
	if err != nil {
		return err
	}

	out, err := ct.createTravail(teacher.IdClassroom(id_), userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

type CreateTravailOut struct {
	Sheet   SheetExt
	Travail ho.Travail
}

func (ct *Controller) createTravail(idClassroom teacher.IdClassroom, userID uID) (CreateTravailOut, error) {
	// check classroom owner
	classroom, err := teacher.SelectClassroom(ct.db, idClassroom)
	if err != nil {
		return CreateTravailOut{}, utils.SQLError(err)
	}
	if classroom.IdTeacher != userID {
		return CreateTravailOut{}, errAccessForbidden
	}

	user, err := teacher.SelectTeacher(ct.db, userID)
	if err != nil {
		return CreateTravailOut{}, utils.SQLError(err)
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return CreateTravailOut{}, utils.SQLError(err)
	}
	// create anonymous Sheet
	sheet, err := ho.Sheet{IdTeacher: userID, Title: "Feuille d'exercices", Matiere: user.FavoriteMatiere}.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return CreateTravailOut{}, utils.SQLError(err)
	}
	// create Travail with this sheet
	tr := ho.Travail{
		IdSheet:     sheet.Id,
		IdClassroom: idClassroom,
		Noted:       true,
		ShowAfter:   ho.Time(time.Now().Round(10 * time.Minute)),
		Deadline:    ho.Time(time.Now().Add(time.Hour * 24 * 7).Round(time.Hour)), // one week
	}
	tr, err = tr.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return CreateTravailOut{}, utils.SQLError(err)
	}
	// mark the sheet as anonymous
	sheet.Anonymous = tr.Id.AsOptional()
	sheet, err = sheet.Update(tx)
	if err != nil {
		_ = tx.Rollback()
		return CreateTravailOut{}, utils.SQLError(err)
	}

	err = tx.Commit()
	if err != nil {
		return CreateTravailOut{}, utils.SQLError(err)
	}

	return CreateTravailOut{Sheet: SheetExt{Sheet: sheet, Origin: tcAPI.Origin{Visibility: tcAPI.Personnal}}, Travail: tr}, nil
}

func checkSheetOwner(db ho.DB, idSheet ho.IdSheet, userID uID) error {
	sheet, err := ho.SelectSheet(db, idSheet)
	if err != nil {
		return utils.SQLError(err)
	}

	// check if the sheet is owned by the user
	if sheet.IdTeacher != userID {
		return errAccessForbidden
	}

	return nil
}

func (ct *Controller) checkSheetOwner(idSheet ho.IdSheet, userID uID) error {
	return checkSheetOwner(ct.db, idSheet, userID)
}

func (ct *Controller) checkTravailOwner(idTravail ho.IdTravail, userID uID) (ho.Travail, error) {
	travail, err := ho.SelectTravail(ct.db, idTravail)
	if err != nil {
		return travail, utils.SQLError(err)
	}

	classroom, err := teacher.SelectClassroom(ct.db, travail.IdClassroom)
	if err != nil {
		return travail, utils.SQLError(err)
	}

	// check if the travail is owned by the user
	if classroom.IdTeacher != userID {
		return travail, errAccessForbidden
	}

	return travail, nil
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
	if _, err := ct.checkTravailOwner(travail.Id, userID); err != nil {
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

	err = ho.InsertSheetTask(tx, ho.SheetTask{IdSheet: sheet, IdTask: task.Id, Index: len(links)})
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

	// check that the question group has questions matching the difficulty
	// to avoid future errors
	variants, err := editor.SelectQuestionsByIdGroups(ct.db, mono.IdQuestiongroup)
	if err != nil {
		return utils.SQLError(err)
	}

	hasOne := false
	for _, qu := range variants {
		if args.Difficulty.Match(qu.Difficulty) {
			hasOne = true
			break
		}
	}
	if !hasOne {
		return errors.New("Aucune variante n'est disponible pour cette difficulté.")
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

func (ct *Controller) HomeworkGetMonoquestion(c echo.Context) error {
	id, err := utils.QueryParamInt64(c, "id-monoquestion")
	if err != nil {
		return err
	}

	out, err := tasks.SelectMonoquestion(ct.db, tasks.IdMonoquestion(id))
	if err != nil {
		return utils.SQLError(err)
	}

	return c.JSON(200, out)
}

func (ct *Controller) HomeworkGetRandomMonoquestion(c echo.Context) error {
	id, err := utils.QueryParamInt64(c, "id-randommonoquestion")
	if err != nil {
		return err
	}

	out, err := tasks.SelectRandomMonoquestion(ct.db, tasks.IdRandomMonoquestion(id))
	if err != nil {
		return utils.SQLError(err)
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

	out, err := ct.duplicateSheet(args.IdSheet, userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) duplicateSheet(idSheet ho.IdSheet, userID uID) (SheetExt, error) {
	tx, err := ct.db.Begin()
	if err != nil {
		return SheetExt{}, utils.SQLError(err)
	}

	newSheet, err := duplicateSheetTx(tx, idSheet, userID, ct.admin.Id)
	if err != nil {
		_ = tx.Rollback()
		return SheetExt{}, err
	}

	err = tx.Commit()
	if err != nil {
		return SheetExt{}, utils.SQLError(err)
	}

	return LoadSheet(ct.db, newSheet.Id, userID, ct.admin.Id)
}

// DO NOT COMMIT, DO NOT ROLLBACK
func duplicateSheetTx(tx *sql.Tx, idSheet ho.IdSheet, userID, adminID uID) (ho.Sheet, error) {
	sheet, err := ho.SelectSheet(tx, idSheet)
	if err != nil {
		return ho.Sheet{}, utils.SQLError(err)
	}

	// duplicate is allowed for public sheet or personnal ones
	if !sheet.IsVisibleBy(userID) {
		return ho.Sheet{}, errAccessForbidden
	}

	links, err := ho.SelectSheetTasksByIdSheets(tx, sheet.Id)
	if err != nil {
		return ho.Sheet{}, utils.SQLError(err)
	}

	taskMap, err := tasks.SelectTasks(tx, links.IdTasks()...)
	if err != nil {
		return ho.Sheet{}, utils.SQLError(err)
	}

	// shallow copy of the item ...
	sheet.Anonymous = ho.OptionalIdTravail{} // new sheet can't have the right travail id, since it does not exists
	// attribute the new copy to the current user, and make it private
	sheet.IdTeacher = userID
	sheet.Public = false
	newSheet, err := sheet.Insert(tx)
	if err != nil {
		return ho.Sheet{}, utils.SQLError(err)
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
				return ho.Sheet{}, utils.SQLError(err)
			}
			monoquestion, err = monoquestion.Insert(tx)
			if err != nil {
				return ho.Sheet{}, utils.SQLError(err)
			}
			newTask.IdMonoquestion = monoquestion.Id.AsOptional()
		} else if id := task.IdRandomMonoquestion; id.Valid {
			monoquestion, err := tasks.SelectRandomMonoquestion(tx, id.ID)
			if err != nil {
				return ho.Sheet{}, utils.SQLError(err)
			}
			monoquestion, err = monoquestion.Insert(tx)
			if err != nil {
				return ho.Sheet{}, utils.SQLError(err)
			}
			newTask.IdRandomMonoquestion = monoquestion.Id.AsOptional()
		}

		newTask, err = newTask.Insert(tx)
		if err != nil {
			return ho.Sheet{}, utils.SQLError(err)
		}
		newLinks[i] = ho.SheetTask{IdSheet: newSheet.Id, IdTask: newTask.Id, Index: i}
	}

	err = ho.InsertManySheetTasks(tx, newLinks...)
	if err != nil {
		return ho.Sheet{}, utils.SQLError(err)
	}

	return newSheet, nil
}

type CopyTravailIn struct {
	IdTravail   ho.IdTravail
	IdClassroom teacher.IdClassroom
}

// HomeworkCopyTravail duplicate the given [Travail] entry,
// updating its Classroom.
// Anonymous [Sheet] are also duplicated.
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

type CopyTravailToOut struct {
	Travail     ho.Travail
	HasNewSheet bool
	NewSheet    SheetExt // valid if and only if HasNewSheet is true
}

func (ct *Controller) copyTravailTo(args CopyTravailIn, userID uID) (CopyTravailToOut, error) {
	cl, err := teacher.SelectClassroom(ct.db, args.IdClassroom)
	if err != nil {
		return CopyTravailToOut{}, utils.SQLError(err)
	}

	if cl.IdTeacher != userID {
		return CopyTravailToOut{}, errAccessForbidden
	}

	travail, err := ho.SelectTravail(ct.db, args.IdTravail)
	if err != nil {
		return CopyTravailToOut{}, utils.SQLError(err)
	}

	sheet, err := ho.SelectSheet(ct.db, travail.IdSheet)
	if err != nil {
		return CopyTravailToOut{}, utils.SQLError(err)
	}

	isAnonymous := sheet.Anonymous.Valid

	tx, err := ct.db.Begin()
	if err != nil {
		return CopyTravailToOut{}, utils.SQLError(err)
	}

	var newSheet SheetExt
	if isAnonymous {
		// also duplicate the underlying sheet
		newSheet.Sheet, err = duplicateSheetTx(tx, sheet.Id, userID, ct.admin.Id)
		if err != nil {
			_ = tx.Rollback()
			return CopyTravailToOut{}, utils.SQLError(err)
		}
		travail.IdSheet = newSheet.Sheet.Id
	}

	// shallow copy is enough
	travail.IdClassroom = args.IdClassroom
	travail, err = travail.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return CopyTravailToOut{}, utils.SQLError(err)
	}

	if isAnonymous {
		// map the new sheet to its new travail
		newSheet.Sheet.Anonymous = travail.Id.AsOptional()
		_, err = newSheet.Sheet.Update(tx)
		if err != nil {
			_ = tx.Rollback()
			return CopyTravailToOut{}, utils.SQLError(err)
		}

		newSheet, err = LoadSheet(tx, newSheet.Sheet.Id, userID, ct.admin.Id)
		if err != nil {
			_ = tx.Rollback()
			return CopyTravailToOut{}, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return CopyTravailToOut{}, utils.SQLError(err)
	}

	out := CopyTravailToOut{Travail: travail}
	if isAnonymous {
		out.HasNewSheet = true
		out.NewSheet = newSheet
	}
	return out, nil
}

func (ct *Controller) HomeworkGetDispenses(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)
	id_, err := utils.QueryParamInt64(c, "id-travail")
	if err != nil {
		return err
	}

	out, err := ct.getDispenses(ho.IdTravail(id_), userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

type Exceptions struct {
	Exceptions ho.TravailExceptions
	Students   teacher.Students // for the classroom
}

func (ct *Controller) getDispenses(idTravail ho.IdTravail, userID uID) (Exceptions, error) {
	travail, err := ct.checkTravailOwner(idTravail, userID)
	if err != nil {
		return Exceptions{}, err
	}
	out, err := ho.SelectTravailExceptionsByIdTravails(ct.db, idTravail)
	if err != nil {
		return Exceptions{}, utils.SQLError(err)
	}

	students, err := teacher.SelectStudentsByIdClassrooms(ct.db, travail.IdClassroom)
	if err != nil {
		return Exceptions{}, utils.SQLError(err)
	}

	return Exceptions{Exceptions: out, Students: students}, nil
}

func (ct *Controller) HomeworkSetDispense(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args ho.TravailException
	if err := c.Bind(&args); err != nil {
		return err
	}

	err := ct.setDispense(args, userID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) setDispense(args ho.TravailException, userID uID) error {
	if _, err := ct.checkTravailOwner(args.IdTravail, userID); err != nil {
		return err
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}

	// remove existing item
	_, err = ho.DeleteTravailExceptionsByIdStudentAndIdTravail(tx, args.IdStudent, args.IdTravail)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	// delete empty items
	if !args.Deadline.Valid && !args.IgnoreForMark {
		//
	} else {
		err = ho.InsertTravailException(tx, args)
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

//
// Student API
//

// StudentGetTravaux returns the sheets for the given student
// Only the mandatory, noted one are returned.
func (ct *Controller) StudentGetTravaux(c echo.Context) error {
	idCrypted := pass.EncryptedID(c.QueryParam("client-id"))

	idStudent, err := ct.studentKey.DecryptID(idCrypted)
	if err != nil {
		return err
	}

	out, err := ct.getStudentSheets(teacher.IdStudent(idStudent), true)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

// StudentGetFreeTravaux returns the sheets for the given student
// Only the optional, non noted are returned.
func (ct *Controller) StudentGetFreeTravaux(c echo.Context) error {
	idCrypted := pass.EncryptedID(c.QueryParam("client-id"))

	idStudent, err := ct.studentKey.DecryptID(idCrypted)
	if err != nil {
		return err
	}

	out, err := ct.getStudentSheets(teacher.IdStudent(idStudent), false)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) getStudentSheets(idStudent teacher.IdStudent, noted bool) (out StudentSheets, err error) {
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

	// load the potential exceptions
	links, err := ho.SelectTravailExceptionsByIdStudents(ct.db, idStudent)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	excepts := links.ByIdTravail()

	for _, travail := range travaux {
		if travail.Noted != noted { // select noted / free travaux
			continue
		}

		// check the start field
		if time.Now().Before(time.Time(travail.ShowAfter)) { // hide the work for now
			continue
		}

		var exp ho.TravailException
		if l := excepts[travail.Id]; len(l) != 0 {
			exp = l[0] // by design there is at most 1 entry for a student and travail
		}
		deadline := travail.Deadline
		if exp.Deadline.Valid {
			deadline = ho.Time(exp.Deadline.Time)
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
				Id:            sheet.Id,
				Title:         sheet.Title,
				Deadline:      deadline,
				Noted:         travail.Noted,
				IgnoreForMark: exp.IgnoreForMark,

				// TODO: cleanup these unused fields
				Notation:    0,
				Activated:   true,
				IdClassroom: 0,
			},
			Tasks: taskList,
		})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Sheet.Id < out[j].Sheet.Id })
	// for noted work, show the most recent first
	if noted {
		sort.SliceStable(out, func(i, j int) bool {
			di, dj := out[i].Sheet.Deadline, out[j].Sheet.Deadline
			return time.Time(di).After(time.Time(dj))
		})
	}
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
	idStudent_, err := ct.studentKey.DecryptID(args.StudentID)
	if err != nil {
		return StudentEvaluateTaskOut{}, err
	}
	idStudent := teacher.IdStudent(idStudent_)

	Logger.Printf("evaluating Task %d (Travail %d) for Student %d", args.IdTask, args.IdTravail, idStudent)

	pr, err := taAPI.LoadTaskProgression(ct.db, idStudent, args.IdTask)
	if err != nil {
		return StudentEvaluateTaskOut{}, err
	}
	isTaskComplete := pr.HasProgression && pr.Progression.IsComplete()

	travail, err := ho.SelectTravail(ct.db, args.IdTravail)
	if err != nil {
		return StudentEvaluateTaskOut{}, utils.SQLError(err)
	}
	var registerProgression bool
	if travail.Noted {
		exp, has, err := ho.SelectTravailExceptionByIdStudentAndIdTravail(ct.db, idStudent, travail.Id)
		if err != nil {
			return StudentEvaluateTaskOut{}, utils.SQLError(err)
		}
		deadline := time.Time(travail.Deadline)
		if has && exp.Deadline.Valid {
			deadline = exp.Deadline.Time
		}
		isExpired := deadline.Before(time.Now())

		// only register progression for non expired, non completed
		registerProgression = !isTaskComplete && !isExpired
	} else {
		// Always register progression for free travail
		registerProgression = true
	}

	ex, mark, err := taAPI.EvaluateTaskExercice(ct.db, args.IdTask, idStudent, args.Ex, registerProgression)
	if err != nil {
		return StudentEvaluateTaskOut{}, err
	}
	return StudentEvaluateTaskOut{Ex: ex, Mark: mark, WasProgressionRegistred: registerProgression}, nil
}

// StudentResetTask remove the progression for the given student
// and task. It is only allowed for free travaux.
func (ct *Controller) StudentResetTask(c echo.Context) error {
	var args StudentResetTaskIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	err := ct.studentResetTask(args)
	if err != nil {
		return err
	}

	return c.JSON(200, true)
}

func (ct *Controller) studentResetTask(args StudentResetTaskIn) error {
	idStudent_, err := ct.studentKey.DecryptID(args.StudentID)
	if err != nil {
		return err
	}
	idStudent := teacher.IdStudent(idStudent_)

	travail, err := ho.SelectTravail(ct.db, args.IdTravail)
	if err != nil {
		return utils.SQLError(err)
	}

	if travail.Noted {
		return errors.New("internal error: travail noted may not be reset")
	}

	task, err := tasks.SelectTask(ct.db, args.IdTask)
	if err != nil {
		return utils.SQLError(err)
	}

	// remove any progression
	tx, err := ct.db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}

	_, err = tasks.DeleteProgressionsByIdStudentAndIdTask(ct.db, idStudent, task.Id)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	// for random monoquestion, remove the selected variants
	if id := task.IdRandomMonoquestion; id.Valid {
		_, err = tasks.DeleteRandomMonoquestionVariantsByIdStudentAndIdRandomMonoquestion(ct.db, idStudent, id.ID)
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
