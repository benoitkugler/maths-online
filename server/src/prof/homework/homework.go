// Package homework implements an activity for student
// consisting in personal, at home training on exercices given by the
// teacher.
package homework

import (
	"database/sql"
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
}

type SheetExt struct {
	Sheet ho.Sheet
	Tasks []TaskExt
}

type sheetLoader struct {
	links map[ho.IdSheet]ho.SheetTasks
	// tasks        tasks.Tasks
	// exes         map[ed.IdExercice]edAPI.ExerciceHeader

	tasks taAPI.TasksContents

	progressions map[tasks.IdTask]tasks.Progressions
}

func newSheetLoader(db ho.DB, idSheets []ho.IdSheet, userID, adminID uID) (out sheetLoader, err error) {
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

func (loader sheetLoader) newSheetExt(sheet ho.Sheet) SheetExt {
	out := SheetExt{Sheet: sheet}
	links := loader.links[sheet.Id]
	links.EnsureOrder()
	for _, link := range links {
		task := loader.tasks.Tasks[link.IdTask]
		work := loader.tasks.GetWork(task)
		out.Tasks = append(out.Tasks, TaskExt{
			Id:             task.Id,
			IdWork:         taAPI.NewWorkID(task),
			Title:          work.Title(),
			Subtitle:       work.Subtitle(),
			NbProgressions: len(loader.progressions[task.Id]),
		})
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

// Student API

// SheetProgression is the summary of the progression
// of one student for one sheet
type SheetProgression struct {
	Sheet ho.Sheet
	Tasks []taAPI.TaskProgressionHeader
}
