package tasks

import (
	"database/sql"
	"fmt"

	ed "github.com/benoitkugler/maths-online/sql/editor"
	ta "github.com/benoitkugler/maths-online/sql/tasks"
	"github.com/benoitkugler/maths-online/sql/teacher"
	"github.com/benoitkugler/maths-online/utils"
)

type ProgressionExt struct {
	Questions    []ta.QuestionHistory
	NextQuestion int
}

// InferNextQuestion stores into `NextQuestion` the first question not passed by the student,
// according to `QuestionHistory.Success`.
// If all the questions are successul, it sets it to -1
func (qh *ProgressionExt) InferNextQuestion() {
	for i, question := range qh.Questions {
		if !question.Success() {
			qh.NextQuestion = i
			return
		}
	}
	qh.NextQuestion = -1
}

type tasksContents struct {
	tasks         ta.Tasks
	groups        ed.Exercicegroups
	exercices     ed.Exercices
	exToQuestions map[ed.IdExercice]ed.ExerciceQuestions

	monoquestions ta.Monoquestions

	questions ed.Questions // provide both exercices and monoquestions contents
}

func loadTasksContents(db ta.DB, ids []ta.IdTask) (out tasksContents, err error) {
	out.tasks, err = ta.SelectTasks(db, ids...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	// fetch the associated exerciceIDs or monoquestionIDs
	exerciceIDs, monoquestionIDs := make(ed.IdExerciceSet), make(ta.IdMonoquestionSet)
	for _, task := range out.tasks {
		if task.IdExercice.Valid {
			exerciceIDs.Add(task.IdExercice.ID)
		} else {
			monoquestionIDs.Add(task.IdMonoquestion.Id)
		}
	}

	out.exercices, err = ed.SelectExercices(db, exerciceIDs.Keys()...)
	if err != nil {
		return out, utils.SQLError(err)
	}
	out.groups, err = ed.SelectExercicegroups(db, out.exercices.IdGroups()...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	links, err := ed.SelectExerciceQuestionsByIdExercices(db, exerciceIDs.Keys()...)
	if err != nil {
		return out, utils.SQLError(err)
	}
	out.exToQuestions = links.ByIdExercice()

	out.monoquestions, err = ta.SelectMonoquestions(db, monoquestionIDs.Keys()...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	questionsIds := append(out.monoquestions.IdQuestions(), links.IdQuestions()...)
	out.questions, err = ed.SelectQuestions(db, questionsIds...)
	if err != nil {
		return out, utils.SQLError(err)
	}
	return out, nil
}

func (contents tasksContents) getWork(task ta.Task) workLoader {
	if task.IdExercice.Valid {
		ex := contents.exercices[task.IdExercice.ID]
		questions := contents.exToQuestions[task.IdExercice.ID]
		return ExerciceData{Group: contents.groups[ex.IdGroup], Exercice: ex, Links: questions, QuestionsSource: contents.questions}
	}

	monoquestion := contents.monoquestions[task.IdMonoquestion.Id]
	return MonoquestionData{params: monoquestion, question: contents.questions[monoquestion.IdQuestion]}
}

// loadProgressions loads the progression contents
func loadProgressions(db ta.DB, prs ta.Progressions) (map[ta.IdProgression]ProgressionExt, error) {
	// fetch the associated tasks
	tasks, err := loadTasksContents(db, prs.IdTasks())
	if err != nil {
		return nil, err
	}

	links2, err := ta.SelectProgressionQuestionsByIdProgressions(db, prs.IDs()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	questionsProgMap := links2.ByIdProgression() // (incomplete) progression of the student

	out := make(map[ta.IdProgression]ProgressionExt, len(prs))
	for _, pr := range prs {
		task := tasks.tasks[pr.IdTask]
		work := tasks.getWork(task)
		// get the questions length
		tmp, _ := work.QuestionsList()
		L := len(tmp)

		// beware that some questions may not have a link item for the student yet
		// so that we take L as reference
		progExt := ProgressionExt{
			Questions: make([]ta.QuestionHistory, L),
		}
		prog := questionsProgMap[pr.Id]
		for _, link := range prog {
			progExt.Questions[link.Index] = link.History
		}
		progExt.InferNextQuestion()

		out[pr.Id] = progExt
	}

	return out, nil
}

// updateProgression write the question results for the given progression.
func updateProgression(db *sql.DB, prog ta.Progression, questions []ta.QuestionHistory) error {
	// sanity checks
	task, err := ta.SelectTask(db, prog.IdTask)
	if err != nil {
		return utils.SQLError(err)
	}

	var expectedLength int
	if task.IdExercice.Valid {
		reference, err := ed.SelectExerciceQuestionsByIdExercices(db, task.IdExercice.ID)
		if err != nil {
			return utils.SQLError(err)
		}
		expectedLength = len(reference)
	} else {
		monoquestion, err := ta.SelectMonoquestion(db, task.IdMonoquestion.Id)
		if err != nil {
			return utils.SQLError(err)
		}
		expectedLength = monoquestion.NbRepeat
	}

	if len(questions) != expectedLength {
		return fmt.Errorf("internal error: inconsistent questions length %d != %d", len(questions), expectedLength)
	}

	tx, err := db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}

	_, err = ta.DeleteProgressionQuestionsByIdProgressions(tx, prog.Id)
	if err != nil {
		return utils.SQLError(err)
	}

	links := make(ta.ProgressionQuestions, len(questions))
	for i, qu := range questions {
		links[i] = ta.ProgressionQuestion{
			IdProgression: prog.Id,
			Index:         i,
			History:       qu,
		}
	}
	err = ta.InsertManyProgressionQuestions(tx, links...)
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
func loadOrCreateProgressionFor(db ta.DB, idTask ta.IdTask, idStudent teacher.IdStudent) (ta.Progression, error) {
	prog, has, err := ta.SelectProgressionByIdStudentAndIdTask(db, idStudent, idTask)
	if err != nil {
		return ta.Progression{}, utils.SQLError(err)
	}
	if has {
		return prog, nil
	}

	// else, create an entry
	prog, err = ta.Progression{IdStudent: idStudent, IdTask: idTask}.Insert(db)
	if err != nil {
		return ta.Progression{}, utils.SQLError(err)
	}

	return prog, nil
}

// Student API

type TaskProgressionHeader struct {
	Id    ta.IdTask
	Title string

	HasProgression bool
	// empty if HasProgression is false
	Progression  ProgressionExt `gomacro-extern:"editor:dart:../shared_gen.dart"`
	Mark, Bareme int            // student mark / exercice total
}

// LoadTasksProgression fetches the progression of one student against
// the given tasks.
func LoadTasksProgression(db ta.DB, idStudent teacher.IdStudent, idTasks []ta.IdTask) (map[ta.IdTask]TaskProgressionHeader, error) {
	contents, err := loadTasksContents(db, idTasks)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	links1, err := ta.SelectProgressionsByIdStudents(db, idStudent)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	// collect the student progressions (we load all the student progression)
	extendedProgressions, err := loadProgressions(db, links1)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	progressionsByTask := links1.ByIdTask()

	out := make(map[ta.IdTask]TaskProgressionHeader, len(contents.tasks))
	for _, task := range contents.tasks {
		// select the right progression, which may be empty
		// before the student starts the exercice,
		// that is progs has either length one or zero
		progs, hasProg := progressionsByTask[task.Id]
		var idProg ta.IdProgression
		if hasProg {
			idProg = progs.IDs()[0]
		}
		progression := extendedProgressions[idProg] // may be empty

		work := contents.getWork(task)
		_, baremes := work.QuestionsList()
		title := work.title()

		out[task.Id] = TaskProgressionHeader{
			Id:    task.Id,
			Title: title,

			HasProgression: hasProg,
			Progression:    progression,
			Bareme:         baremes.total(),
			Mark:           baremes.computeMark(progression.Questions),
		}
	}

	return out, nil
}

// EvaluateTaskExercice calls `EvaluateExercice` and registers
// the student progression, returning the updated mark.
// If needed, a new progression item is created.
func EvaluateTaskExercice(db *sql.DB, idTask ta.IdTask, idStudent teacher.IdStudent, ex EvaluateWorkIn) (out EvaluateWorkOut, mark int, err error) {
	out, err = EvaluateWork(db, ex)
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
	task, err := ta.SelectTask(db, idTask)
	if err != nil {
		return out, mark, utils.SQLError(err)
	}

	loader, err := newWorkLoader(db, newWorkID(task))
	if err != nil {
		return
	}
	_, baremes := loader.QuestionsList()
	mark = baremes.computeMark(out.Progression.Questions)

	return out, mark, nil
}

type taskBareme []int // for each question

// total returns the bareme for one task, that is the sum
// of each question's bareme
func (bareme taskBareme) total() int {
	var out int
	for _, questionBareme := range bareme {
		out += questionBareme
	}
	return out
}

// compute the student mark
// an empty progression is supported and returns 0
func (bareme taskBareme) computeMark(progression []ta.QuestionHistory) int {
	if len(progression) == 0 {
		return 0
	}

	var out int
	for index, baremeQuestion := range bareme {
		results := progression[index]
		if results.Success() {
			out += baremeQuestion
		}
	}
	return out
}
