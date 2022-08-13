// Package tasks exposes the data structure
// required to assign exercices during activities,
// and tracking the progression of the students.
package tasks

import (
	"database/sql"
	"sort"

	ed "github.com/benoitkugler/maths-online/prof/editor"
	"github.com/benoitkugler/maths-online/prof/teacher"
	"github.com/benoitkugler/maths-online/utils"
)

type QuestionHistory ed.QuestionHistory

func (l ProgressionQuestions) ensureOrder() {
	sort.Slice(l, func(i, j int) bool { return l[i].Index < l[j].Index })
}

// Student API

type TaskProgressionHeader struct {
	Id            IdTask
	IdExercice    ed.IdExercice
	TitleExercice string

	HasProgression bool
	// empty if HasProgression is false
	Progression  ed.ProgressionExt `gomacro-extern:"editor:dart:../shared_gen.dart"`
	Mark, Bareme int               // student mark / exercice total
}

// LoadTasksProgression fetches the progression of one student against
// the given tasks.
func LoadTasksProgression(db DB, idStudent teacher.IdStudent, idTasks []IdTask) (map[IdTask]TaskProgressionHeader, error) {
	tasks, err := SelectTasks(db, idTasks...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	exercices, err := ed.SelectExercices(db, tasks.IdExercices()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	links1, err := SelectProgressionsByIdStudents(db, idStudent)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	// collect the student progressions (we load all the student progression)
	extendedProgressions, err := loadProgressions(db, links1)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	// load the questions of each exercice
	links2, err := ed.SelectExerciceQuestionsByIdExercices(db, exercices.IDs()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	questionsByExercice := links2.ByIdExercice()

	progressionsByTask := links1.ByIdTask()

	out := make(map[IdTask]TaskProgressionHeader, len(tasks))
	for _, task := range tasks {
		// select the right progression, which may be empty
		// before the student starts the exercice,
		// that is progs has either length one or zero
		progs, hasProg := progressionsByTask[task.Id]
		var idProg IdProgression
		if hasProg {
			idProg = progs.IDs()[0]
		}

		exercice := exercices[task.IdExercice]
		questions := questionsByExercice[task.IdExercice]
		progression := extendedProgressions[idProg]

		out[task.Id] = TaskProgressionHeader{
			Id:            task.Id,
			IdExercice:    exercice.Id,
			TitleExercice: exercice.Title,

			HasProgression: hasProg,
			Progression:    progression,
			Bareme:         questions.Bareme(),
			Mark:           computeMark(questions, progression.Questions),
		}
	}

	return out, nil
}

// loadProgressions loads the progression contents
func loadProgressions(db DB, prs Progressions) (map[IdProgression]ed.ProgressionExt, error) {
	// select the associated exercices
	exercices := make(ed.IdExerciceSet)
	for _, pr := range prs {
		exercices.Add(pr.IdExercice)
	}

	links1, err := ed.SelectExerciceQuestionsByIdExercices(db, exercices.Keys()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	questionsExesMap := links1.ByIdExercice() // reference from the exercice

	links2, err := SelectProgressionQuestionsByIdProgressions(db, prs.IDs()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	questionsProgMap := links2.ByIdProgression() // (incomplete) progression of the student

	out := make(map[IdProgression]ed.ProgressionExt, len(prs))
	for _, pr := range prs {
		questions := questionsExesMap[pr.IdExercice]
		questions.EnsureOrder()

		// beware that some questions may not have a link item yet
		progExt := ed.ProgressionExt{
			Questions: make([]ed.QuestionHistory, len(questions)),
		}
		prog := questionsProgMap[pr.Id]
		for _, link := range prog {
			progExt.Questions[link.Index] = ed.QuestionHistory(link.History)
		}
		progExt.InferNextQuestion()

		out[pr.Id] = progExt
	}

	return out, nil
}

// updateProgression write the question results for the given progression.
func updateProgression(db *sql.DB, prog Progression, questions []ed.QuestionHistory) error {
	tx, err := db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}

	_, err = DeleteProgressionQuestionsByIdProgressions(tx, prog.Id)
	if err != nil {
		return utils.SQLError(err)
	}

	links := make(ProgressionQuestions, len(questions))
	for i, qu := range questions {
		links[i] = ProgressionQuestion{
			IdProgression: prog.Id,
			IdExercice:    prog.IdExercice,
			Index:         i,
			History:       QuestionHistory(qu),
		}
	}
	err = InsertManyProgressionQuestions(tx, links...)
	if err != nil {
		return utils.SQLError(err)
	}

	err = tx.Commit()
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}

// loadOrCreateProgressionFor ensures a progression item exists for
// the given task and student, and returns it.
func loadOrCreateProgressionFor(db DB, idTask IdTask, idStudent teacher.IdStudent) (Progression, error) {
	task, err := SelectTask(db, idTask)
	if err != nil {
		return Progression{}, utils.SQLError(err)
	}

	prog, has, err := SelectProgressionByIdStudentAndIdTask(db, idStudent, idTask)
	if err != nil {
		return Progression{}, utils.SQLError(err)
	}
	if has {
		return prog, nil
	}

	// else, create an entry
	prog, err = Progression{IdStudent: idStudent, IdTask: idTask, IdExercice: task.IdExercice}.Insert(db)
	if err != nil {
		return Progression{}, utils.SQLError(err)
	}

	return prog, nil
}

// EvaluateTaskExercice calls `editor.EvaluateTaskExercice` and registers
// the student progression, returning the updated mark.
// If needed, a new progression item is created.
func EvaluateTaskExercice(db *sql.DB, idTask IdTask, idStudent teacher.IdStudent, ex ed.EvaluateExerciceIn) (out ed.EvaluateExerciceOut, mark int, err error) {
	out, err = ed.EvaluateExercice(db, ex)
	if err != nil {
		return
	}

	// persists the progression on DB ...
	prog, err := loadOrCreateProgressionFor(db, idTask, idStudent)
	if err != nil {
		return
	}

	// ... update the progression questions
	err = updateProgression(db, prog, out.Progression.Questions)
	if err != nil {
		return
	}

	// and compute the new mark
	questions, err := ed.SelectExerciceQuestionsByIdExercices(db, ex.IdExercice)
	if err != nil {
		return out, mark, utils.SQLError(err)
	}
	mark = computeMark(questions, out.Progression.Questions)

	return out, mark, nil
}

// compute the student mark
// an empty progression is supported and returns 0
func computeMark(questions ed.ExerciceQuestions, progression []ed.QuestionHistory) int {
	if len(progression) == 0 {
		return 0
	}

	questions.EnsureOrder()

	var out int
	for index, qu := range questions {
		results := progression[index]
		if results.Success() {
			out += qu.Bareme
		}
	}
	return out
}
