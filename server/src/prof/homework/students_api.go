package homework

import (
	"database/sql"
	"errors"
	"sort"
	"time"

	"github.com/benoitkugler/maths-online/server/src/pass"
	"github.com/benoitkugler/maths-online/server/src/sql/events"
	ho "github.com/benoitkugler/maths-online/server/src/sql/homework"
	"github.com/benoitkugler/maths-online/server/src/sql/tasks"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	taAPI "github.com/benoitkugler/maths-online/server/src/tasks"
	"github.com/benoitkugler/maths-online/server/src/utils"

	"github.com/labstack/echo/v4"
)

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

// do not perform any filtering nor sorting
func loadSheetProgressions(db ho.DB, idStudent teacher.IdStudent, travaux ho.Travails) ([]SheetProgression, error) {
	sheets, err := ho.SelectSheets(db, travaux.IdSheets()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	links1, err := ho.SelectSheetTasksByIdSheets(db, sheets.IDs()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	sheetToTasks := links1.ByIdSheet()

	// collect the student progressions
	progMap, err := taAPI.LoadTasksProgression(db, idStudent, links1.IdTasks())
	if err != nil {
		return nil, utils.SQLError(err)
	}

	// load the potential exceptions
	links, err := ho.SelectTravailExceptionsByIdStudents(db, idStudent)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	excepts := links.ByIdTravail()

	out := make([]SheetProgression, 0, len(travaux))
	for _, travail := range travaux {
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

		matiere := sheet.Matiere
		if sheet.Anonymous.Valid {
			// use the content
			if len(taskList) != 0 {
				matiere = taAPI.MatiereFromTasks(taskList)
			}
		}
		out = append(out, SheetProgression{
			IdTravail: travail.Id,
			Sheet: Sheet{
				sheet.Id,
				sheet.Title,
				travail.Noted,
				deadline,
				exp.IgnoreForMark,
				matiere,
				travail.QuestionRepeat,
				travail.QuestionTimeLimit,
			},
			Tasks: taskList,
		})
	}
	return out, nil
}

func (ct *Controller) getStudentSheets(idStudent teacher.IdStudent, noted bool) (StudentSheets, error) {
	student, err := teacher.SelectStudent(ct.db, idStudent)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	travaux, err := ho.SelectTravailsByIdClassrooms(ct.db, student.IdClassroom)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	for id, travail := range travaux {
		if travail.Noted != noted { // select noted / free travaux
			delete(travaux, id)
		}

		// check the start field
		if time.Now().Before(time.Time(travail.ShowAfter)) { // hide the work for now
			delete(travaux, id)
		}
	}

	out, err := loadSheetProgressions(ct.db, idStudent, travaux)
	if err != nil {
		return nil, err
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

func (ct *Controller) StudentLoadTravail(c echo.Context) error {
	idCrypted := pass.EncryptedID(c.QueryParam("client-id"))
	idStudent, err := ct.studentKey.DecryptID(idCrypted)
	if err != nil {
		return err
	}
	idTravail, err := utils.QueryParamInt[ho.IdTravail](c, "idTravail")
	if err != nil {
		return err
	}
	out, err := ct.getStudentTravail(teacher.IdStudent(idStudent), idTravail)
	if err != nil {
		return err
	}
	return c.JSON(200, out)
}

func (ct *Controller) getStudentTravail(idStudent teacher.IdStudent, idTravail ho.IdTravail) (SheetProgression, error) {
	travail, err := ho.SelectTravail(ct.db, idTravail)
	if err != nil {
		return SheetProgression{}, utils.SQLError(err)
	}
	out, err := loadSheetProgressions(ct.db, idStudent, ho.Travails{travail.Id: travail})
	if err != nil {
		return SheetProgression{}, err
	}
	return out[0], nil
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

	// register success
	ev := events.E_All_QuestionWrong
	if ex.Result.IsCorrect() {
		ev = events.E_All_QuestionRight
	}
	evL := []events.EventK{ev}
	if ex.Progression.IsComplete() {
		evL = append(evL, events.E_Homework_TaskDone)
	}
	// do we have finished the whole travail ?
	ld, err := newSheetsLoader(ct.db, []ho.IdSheet{travail.IdSheet})
	if err != nil {
		return StudentEvaluateTaskOut{}, err
	}
	if len(ld.tasksForSheet(travail.IdSheet)) >= 2 {
		travailComplete, err := ld.isSheetComplete(ct.db, idStudent, travail.IdSheet)
		if err != nil {
			return StudentEvaluateTaskOut{}, err
		}
		if travailComplete {
			evL = append(evL, events.E_Homework_TravailDone)
		}
	}

	notif, err := events.RegisterEvents(ct.db, idStudent, evL...)
	if err != nil {
		return StudentEvaluateTaskOut{}, err
	}

	return StudentEvaluateTaskOut{Ex: ex, Mark: mark, WasProgressionRegistred: registerProgression, Advance: notif}, nil
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

	return utils.InTx(ct.db, func(tx *sql.Tx) error {
		// remove any progression
		_, err = tasks.DeleteProgressionsByIdStudentAndIdTask(tx, idStudent, task.Id)
		if err != nil {
			return err
		}

		// for random monoquestion, remove the selected variants
		if id := task.IdRandomMonoquestion; id.Valid {
			_, err = tasks.DeleteRandomMonoquestionVariantsByIdStudentAndIdRandomMonoquestion(tx, idStudent, id.ID)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
