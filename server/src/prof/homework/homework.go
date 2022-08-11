// Package homework implements an activity for student
// consisting in personal, at home training on exercices given by the
// teacher.
package homework

import (
	"database/sql"
	"sort"

	ed "github.com/benoitkugler/maths-online/prof/editor"
	"github.com/benoitkugler/maths-online/prof/teacher"
	"github.com/benoitkugler/maths-online/tasks"
	"github.com/benoitkugler/maths-online/utils"
)

type ClassroomSheets struct {
	Classroom teacher.Classroom
	Sheets    []SheetExt
}

func newClassroomSheets(cl teacher.Classroom, sheepMap map[IdSheet]SheetExt) ClassroomSheets {
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
	Exercice       ed.ExerciceHeader
	NbProgressions int // the number of student having started this task
}

type SheetExt struct {
	Sheet Sheet
	Tasks []TaskExt
}

type sheetLoader struct {
	links        map[IdSheet]SheetTasks
	tasks        tasks.Tasks
	exes         map[ed.IdExercice]ed.ExerciceHeader
	progressions map[tasks.IdTask]tasks.Progressions
}

func newSheetLoader(db DB, idSheets []IdSheet, userID, adminID uID) (out sheetLoader, err error) {
	// load all the tasks and exercices required...
	links1, err := SelectSheetTasksByIdSheets(db, idSheets...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	tasksMap, err := tasks.SelectTasks(db, links1.IdTasks()...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	exes, err := ed.SelectExercices(db, tasksMap.IdExercices()...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	// ... and their questions
	questions, err := ed.SelectExerciceQuestionsByIdExercices(db, exes.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	// also load the current progressions
	links2, err := tasks.SelectProgressionsByIdTasks(db, tasksMap.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	out.links = links1.ByIdSheet()
	out.tasks = tasksMap
	out.exes = ed.BuildExerciceHeaders(userID, adminID, exes, questions)
	out.progressions = links2.ByIdTask()

	return out, nil
}

func (loader sheetLoader) newSheetExt(sheet Sheet) SheetExt {
	out := SheetExt{Sheet: sheet}
	links := loader.links[sheet.Id]
	links.ensureOrder()
	for _, link := range links {
		idExercice := loader.tasks[link.IdTask].IdExercice
		out.Tasks = append(out.Tasks, TaskExt{
			Id:             link.IdTask,
			Exercice:       loader.exes[idExercice],
			NbProgressions: len(loader.progressions[link.IdTask]),
		})
	}
	return out
}

func (loader sheetLoader) buildSheetExts(sheets Sheets) map[IdSheet]SheetExt {
	out := make(map[IdSheet]SheetExt, len(sheets))
	for idSheet, v := range sheets {
		out[idSheet] = loader.newSheetExt(v)
	}
	return out
}

func (l SheetTasks) ensureOrder() {
	sort.Slice(l, func(i, j int) bool { return l[i].Index < l[j].Index })
}

func updateSheetTasksOrder(tx *sql.Tx, idSheet IdSheet, l []tasks.IdTask) error {
	links := make(SheetTasks, len(l))
	for i, idTask := range l { // enforce correct index
		links[i] = SheetTask{IdTask: idTask, IdSheet: idSheet, Index: i}
	}

	_, err := DeleteSheetTasksByIdSheets(tx, idSheet)
	if err != nil {
		return utils.SQLError(err)
	}

	err = InsertManySheetTasks(tx, links...)
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}

// Student API

// SheetProgression is the summary of the progression
// of one student for one sheet
type SheetProgression struct {
	Sheet Sheet
	Tasks []tasks.TaskProgressionHeader
}
