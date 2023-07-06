// Package homework implements an activity for student
// consisting in personal, at home training on exercices given by the
// teacher.
package homework

import (
	"database/sql"
	"fmt"
	"sort"

	ho "github.com/benoitkugler/maths-online/server/src/sql/homework"
	"github.com/benoitkugler/maths-online/server/src/sql/tasks"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	taAPI "github.com/benoitkugler/maths-online/server/src/tasks"
	"github.com/benoitkugler/maths-online/server/src/utils"
)

type ClassroomTravaux struct {
	Classroom teacher.Classroom
	Travaux   []ho.Travail
}

func newClassroomTravaux(cl teacher.Classroom, travailMap ho.Travails) ClassroomTravaux {
	out := ClassroomTravaux{Classroom: cl}
	for _, tr := range travailMap {
		if tr.IdClassroom == cl.Id {
			out.Travaux = append(out.Travaux, tr)
		}
	}

	// show Noted first
	sort.Slice(out.Travaux, func(i, j int) bool { return out.Travaux[i].Id < out.Travaux[j].Id })
	sort.SliceStable(out.Travaux, func(i, j int) bool { return out.Travaux[i].Noted && !out.Travaux[j].Noted })
	return out
}

type TaskExt struct {
	Id             tasks.IdTask
	IdWork         taAPI.WorkID
	Title          string // title of the underlying exercice or question
	Bareme         taAPI.TaskBareme
	NbProgressions int // the number of student having started this task
}

// [progressions] is the list of all links item related to [task]
func newTaskExt(task tasks.Task, work taAPI.WorkMeta, progressions tasks.Progressions) TaskExt {
	baremes := work.Bareme()
	return TaskExt{
		Id:             task.Id,
		IdWork:         taAPI.NewWorkID(task),
		Title:          work.Title(),
		NbProgressions: len(progressions.ByIdStudent()),
		Bareme:         baremes,
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
	Sheet     ho.Sheet
	Tasks     []TaskExt
	NbTravaux int
}

// sheetLoader is an helper type to
// unify sheet tasks loading
type sheetLoader struct {
	links map[ho.IdSheet]ho.SheetTasks

	tasks taAPI.TasksContents

	progressions map[tasks.IdTask]tasks.Progressions

	travaux map[ho.IdSheet]ho.Travails
}

func newSheetsLoader(db ho.DB, idSheets []ho.IdSheet) (out sheetLoader, err error) {
	travaux, err := ho.SelectTravailsByIdSheets(db, idSheets...)
	if err != nil {
		return out, utils.SQLError(err)
	}
	out.travaux = travaux.ByIdSheet()

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

func (loader sheetLoader) tasksForSheet(id ho.IdSheet) ho.SheetTasks {
	links := loader.links[id]
	links.EnsureOrder()
	return links
}

func (loader sheetLoader) newSheetExt(sheet ho.Sheet) SheetExt {
	out := SheetExt{
		Sheet:     sheet,
		NbTravaux: len(loader.travaux[sheet.Id]),
	}
	links := loader.tasksForSheet(sheet.Id)
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
