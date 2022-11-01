// Package homework implements an activity for student
// consisting in personal, at home training on exercices given by the
// teacher.
package homework

import (
	"database/sql"
	"fmt"
	"sort"

	ho "github.com/benoitkugler/maths-online/sql/homework"
	"github.com/benoitkugler/maths-online/sql/tasks"
	"github.com/benoitkugler/maths-online/sql/teacher"
	taAPI "github.com/benoitkugler/maths-online/tasks"
	"github.com/benoitkugler/maths-online/utils"
)

type ClassroomSheets struct {
	Classroom teacher.Classroom
	Sheets    []SheetExt
}

func newClassroomSheets(cl teacher.Classroom, sheepMap map[ho.IdSheet]SheetExt) ClassroomSheets {
	out := ClassroomSheets{Classroom: cl}
	for _, sheet := range sheepMap {
		if sheet.Sheet.IdClassroom == cl.Id {
			out.Sheets = append(out.Sheets, sheet)
		}
	}
	sort.Slice(out.Sheets, func(i, j int) bool { return out.Sheets[i].Sheet.Id < out.Sheets[j].Sheet.Id })
	return out
}

type TaskExt struct {
	Id             tasks.IdTask
	IdWork         taAPI.WorkID
	Title          string // title of the underlying exercice or question
	Subtitle       string // subtitle of the underlying exercice or question
	NbProgressions int    // the number of student having started this task
	Baremes        taAPI.TaskBareme
}

func newTaskExt(task tasks.Task, work taAPI.Work, progressions tasks.Progressions) TaskExt {
	_, baremes := work.QuestionsList()
	return TaskExt{
		Id:             task.Id,
		IdWork:         taAPI.NewWorkID(task),
		Title:          work.Title(),
		Subtitle:       work.Subtitle(),
		NbProgressions: len(progressions),
		Baremes:        baremes,
	}
}

func loadTaskExt(db ho.DB, idTask tasks.IdTask) (TaskExt, error) {
	loader, err := taAPI.NewTasksContents(db, []tasks.IdTask{idTask})
	if err != nil {
		return TaskExt{}, err
	}

	task := loader.Tasks[idTask]

	progressions, err := tasks.SelectProgressionsByIdTasks(db, idTask)
	if err != nil {
		return TaskExt{}, utils.SQLError(err)
	}

	return newTaskExt(task, loader.GetWork(task), progressions), nil
}

type SheetExt struct {
	Sheet ho.Sheet
	Tasks []TaskExt
}

// sheetLoader is an helper type to
// unify sheet tasks loading
type sheetLoader struct {
	links map[ho.IdSheet]ho.SheetTasks

	tasks taAPI.TasksContents

	progressions map[tasks.IdTask]tasks.Progressions
}

func newSheetsLoader(db ho.DB, idSheets []ho.IdSheet) (out sheetLoader, err error) {
	// load all the tasks and exercices required...
	links1, err := ho.SelectSheetTasksByIdSheets(db, idSheets...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	out.tasks, err = taAPI.NewTasksContents(db, links1.IdTasks())
	if err != nil {
		return out, err
	}

	// also load the current progressions
	links2, err := tasks.SelectProgressionsByIdTasks(db, out.tasks.Tasks.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	out.links = links1.ByIdSheet()
	out.progressions = links2.ByIdTask()

	return out, nil
}

func (loader sheetLoader) allProgressions() (out tasks.Progressions) {
	out = make(tasks.Progressions)
	for _, progressionForTask := range loader.progressions {
		for _, prog := range progressionForTask {
			out[prog.Id] = prog
		}
	}
	return out
}

func (loader sheetLoader) taskForSheet(id ho.IdSheet) ho.SheetTasks {
	links := loader.links[id]
	links.EnsureOrder()
	return links
}

func (loader sheetLoader) newSheetExt(sheet ho.Sheet) SheetExt {
	out := SheetExt{Sheet: sheet}
	links := loader.taskForSheet(sheet.Id)
	for _, link := range links {
		task := loader.tasks.Tasks[link.IdTask]
		work := loader.tasks.GetWork(task)
		out.Tasks = append(out.Tasks, newTaskExt(task, work, loader.progressions[task.Id]))
	}
	return out
}

func (loader sheetLoader) buildSheetExts(sheets ho.Sheets) map[ho.IdSheet]SheetExt {
	out := make(map[ho.IdSheet]SheetExt, len(sheets))
	for idSheet, v := range sheets {
		out[idSheet] = loader.newSheetExt(v)
	}
	return out
}

func updateSheetTasksOrder(tx *sql.Tx, idSheet ho.IdSheet, l []tasks.IdTask) error {
	links := make(ho.SheetTasks, len(l))
	for i, idTask := range l { // enforce correct index
		links[i] = ho.SheetTask{IdTask: idTask, IdSheet: idSheet, Index: i}
	}

	_, err := ho.DeleteSheetTasksByIdSheets(tx, idSheet)
	if err != nil {
		return utils.SQLError(err)
	}

	err = ho.InsertManySheetTasks(tx, links...)
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}

func sheetFromTask(db tasks.DB, idTask tasks.IdTask) (ho.Sheet, error) {
	link, found, err := ho.SelectSheetTaskByIdTask(db, idTask)
	if err != nil {
		return ho.Sheet{}, utils.SQLError(err)
	}

	if !found {
		return ho.Sheet{}, fmt.Errorf("internal error: task %d without sheet", idTask)
	}

	return ho.SelectSheet(db, link.IdSheet)
}
