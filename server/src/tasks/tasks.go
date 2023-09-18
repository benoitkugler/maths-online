package tasks

import (
	"database/sql"
	"fmt"
	"sort"

	ed "github.com/benoitkugler/maths-online/server/src/sql/editor"
	ta "github.com/benoitkugler/maths-online/server/src/sql/tasks"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/benoitkugler/maths-online/server/src/utils"
)

type ProgressionExt struct {
	// Questions stores the progression for each question of the task.
	Questions    []ta.QuestionHistory
	NextQuestion int
}

func NewProgressionExt(progressions ta.Progressions, nbQuestions int) (out ProgressionExt) {
	out.Questions = make([]ta.QuestionHistory, nbQuestions)
	for _, link := range progressions {
		out.Questions[link.Index] = link.History
	}
	out.inferNextQuestion()
	return out
}

func (qh ProgressionExt) Copy() ProgressionExt {
	return ProgressionExt{
		NextQuestion: qh.NextQuestion,
		Questions:    append([]ta.QuestionHistory(nil), qh.Questions...),
	}
}

// NbTries returns the total number of tries for this progression.
func (qh ProgressionExt) NbTries() int {
	s := 0
	for _, qu := range qh.Questions {
		s += len(qu)
	}
	return s
}

func (qh ProgressionExt) IsComplete() bool { return qh.NextQuestion == -1 }

// inferNextQuestion stores into `NextQuestion` the first question not passed by the student,
// according to `QuestionHistory.Success`.
// If all the questions are successul, it sets it to -1
func (qh *ProgressionExt) inferNextQuestion() {
	for i, question := range qh.Questions {
		if !question.Success() {
			qh.NextQuestion = i
			return
		}
	}
	qh.NextQuestion = -1
}

// TasksContents is an helper struct to unify tasks loading.
type TasksContents struct {
	Tasks          ta.Tasks
	exercicegroups ed.Exercicegroups
	exercices      ed.Exercices
	exToQuestions  map[ed.IdExercice]ed.ExerciceQuestions

	monoquestions       ta.Monoquestions
	randomMonoquestions ta.RandomMonoquestions

	questiongroups ed.Questiongroups // for questions in [monoquestions] and [randomMonoquestions]

	questions ed.Questions // provide exercices and monoquestions contents

	exerciceTags map[ed.IdExercicegroup]ed.ExercicegroupTags
	questionTags map[ed.IdQuestiongroup]ed.QuestiongroupTags // for [monoquestions] and [randomMonoquestions]

	selectedVariants map[teacher.IdStudent]ta.RandomMonoquestionVariants // for [randomMonoquestions], only used in [ResolveQuestions]
}

func NewTasksContents(db ta.DB, ids []ta.IdTask) (out TasksContents, err error) {
	out.Tasks, err = ta.SelectTasks(db, ids...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	// fetch the associated exerciceIDs or monoquestionIDs
	exerciceIDs, monoquestionIDs, randomMonoquestionIDs := make(ed.IdExerciceSet), make(ta.IdMonoquestionSet), make(ta.IdRandomMonoquestionSet)
	for _, task := range out.Tasks {
		if task.IdExercice.Valid {
			exerciceIDs.Add(task.IdExercice.ID)
		} else if task.IdMonoquestion.Valid {
			monoquestionIDs.Add(task.IdMonoquestion.ID)
		} else if task.IdRandomMonoquestion.Valid {
			randomMonoquestionIDs.Add(task.IdRandomMonoquestion.ID)
		}
	}

	out.exercices, err = ed.SelectExercices(db, exerciceIDs.Keys()...)
	if err != nil {
		return out, utils.SQLError(err)
	}
	out.exercicegroups, err = ed.SelectExercicegroups(db, out.exercices.IdGroups()...)
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

	out.randomMonoquestions, err = ta.SelectRandomMonoquestions(db, randomMonoquestionIDs.Keys()...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	groupsFromRandom, err := ed.SelectQuestiongroups(db, out.randomMonoquestions.IdQuestiongroups()...)
	if err != nil {
		return out, utils.SQLError(err)
	}
	questionsFromRandom, err := ed.SelectQuestionsByIdGroups(db, groupsFromRandom.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	questionsIds := append(out.monoquestions.IdQuestions(), links.IdQuestions()...)
	out.questions, err = ed.SelectQuestions(db, questionsIds...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	questiongroupsId := make(ed.IdQuestiongroupSet) // select the groups required
	for _, qu := range out.questions {
		if qu.IdGroup.Valid {
			questiongroupsId.Add(qu.IdGroup.ID)
		}
	}

	out.questiongroups, err = ed.SelectQuestiongroups(db, questiongroupsId.Keys()...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	// merge from random
	for k, v := range groupsFromRandom {
		out.questiongroups[k] = v
	}
	for k, v := range questionsFromRandom {
		out.questions[k] = v
	}

	// load tags
	eTags, err := ed.SelectExercicegroupTagsByIdExercicegroups(db, out.exercicegroups.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}
	out.exerciceTags = eTags.ByIdExercicegroup()
	qTags, err := ed.SelectQuestiongroupTagsByIdQuestiongroups(db, out.questiongroups.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}
	out.questionTags = qTags.ByIdQuestiongroup()

	tmp, err := ta.SelectRandomMonoquestionVariantsByIdRandomMonoquestions(db, randomMonoquestionIDs.Keys()...)
	if err != nil {
		return out, utils.SQLError(err)
	}
	out.selectedVariants = tmp.ByIdStudent()

	return out, nil
}

// GetWork returns the task content for `task`.
func (contents TasksContents) GetWork(task ta.Task) WorkMeta {
	switch {
	case task.IdExercice.Valid:
		ex := contents.exercices[task.IdExercice.ID]
		questions := contents.exToQuestions[task.IdExercice.ID]
		tags := contents.exerciceTags[ex.IdGroup]
		return ExerciceData{
			Group:        contents.exercicegroups[ex.IdGroup],
			Exercice:     ex,
			links:        questions,
			QuestionsMap: contents.questions,
			chapter:      tags.Tags().BySection().Chapter,
		}
	case task.IdMonoquestion.Valid:
		monoquestion := contents.monoquestions[task.IdMonoquestion.ID]
		question := contents.questions[monoquestion.IdQuestion]
		tags := contents.questionTags[question.IdGroup.ID]
		return monoquestionData{
			params:        monoquestion,
			question:      question,
			questiongroup: contents.questiongroups[question.IdGroup.ID],
			chapter:       tags.Tags().BySection().Chapter,
		}
	case task.IdRandomMonoquestion.Valid:
		mono := contents.randomMonoquestions[task.IdRandomMonoquestion.ID]
		tags := contents.questionTags[mono.IdQuestiongroup]
		// for this use case, leaving [selectedQuestions] is OK
		return randomMonoquestionData{
			params:        mono,
			questiongroup: contents.questiongroups[mono.IdQuestiongroup],
			chapter:       tags.Tags().BySection().Chapter,
		}
	default: // should not happen (enforced by SQL constraint)
		return nil
	}
}

// LoadProgressions load the question progression related to the tasks
// in [contents].
func (contents TasksContents) LoadProgressions(db ta.DB) (map[ta.IdTask]map[teacher.IdStudent]ProgressionExt, error) {
	tmp, err := ta.SelectProgressionsByIdTasks(db, contents.Tasks.IDs()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	byTask := tmp.ByIdTask() // (incomplete) progression of the students

	out := make(map[ta.IdTask]map[teacher.IdStudent]ProgressionExt)
	for _, task := range contents.Tasks {
		taskMap := make(map[teacher.IdStudent]ProgressionExt)
		work := contents.GetWork(task)
		// get the questions length
		L := len(work.Bareme())

		byStudent := byTask[task.Id].ByIdStudent() // for one task
		for idStudent, progressions := range byStudent {
			// beware that some questions may not have a link item for the student yet
			// so that we take L as reference
			progExt := NewProgressionExt(progressions, L)
			taskMap[idStudent] = progExt
		}

		out[task.Id] = taskMap
	}

	return out, nil
}

// ResolveQuestions returns the question variants actually done by the student.
// For [RandomMonoquestion]s, the student is expected to have actually started it.
func (contents TasksContents) ResolveQuestions(idStudent teacher.IdStudent, work WorkMeta) []ed.Question {
	switch work := work.(type) {
	case ExerciceData:
		return work.Questions()
	case monoquestionData:
		return work.Questions()
	case randomMonoquestionData:
		l := contents.selectedVariants[idStudent].ByIdRandomMonoquestion()[work.params.Id]
		l.EnsureOrder()
		out := make([]ed.Question, len(l))
		for i, item := range l {
			out[i] = contents.questions[item.IdQuestion]
		}
		return out
	default:
		panic("exhaustive switch")
	}
}

func (contents TasksContents) OrderQuestions(work WorkMeta) []ed.Question {
	switch work := work.(type) {
	case ExerciceData:
		return work.Questions()
	case monoquestionData:
		return []ed.Question{work.question}
	case randomMonoquestionData:
		l := contents.questions.ByGroup()[work.questiongroup.Id]
		sort.Slice(l, func(i, j int) bool { return l[i].Difficulty < l[j].Difficulty })
		return l
	default:
		panic("exhaustive switch")
	}
}

// updateProgression write the question results for the given progression.
func updateProgression(db *sql.DB, idStudent teacher.IdStudent, idTask ta.IdTask, questions []ta.QuestionHistory) error {
	// sanity checks
	task, err := ta.SelectTask(db, idTask)
	if err != nil {
		return utils.SQLError(err)
	}

	work, err := newWorkLoader(db, NewWorkID(task), idStudent)
	if err != nil {
		return err
	}

	expectedLength := len(work.Bareme())

	if len(questions) != expectedLength {
		return fmt.Errorf("internal error: inconsistent questions length %d != %d", len(questions), expectedLength)
	}

	tx, err := db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}

	_, err = ta.DeleteProgressionsByIdStudentAndIdTask(tx, idStudent, idTask)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	links := make(ta.Progressions, len(questions))
	for i, qu := range questions {
		links[i] = ta.Progression{
			IdStudent: idStudent,
			IdTask:    idTask,
			Index:     i,
			History:   qu,
		}
	}
	err = ta.InsertManyProgressions(tx, links...)
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

// Student API

type TaskProgressionHeader struct {
	Id    ta.IdTask
	Title string

	// The chapter of the task content,
	// maybe empty
	Chapter string

	// HasProgression is false if [Progression] is invalid
	HasProgression bool
	// empty if HasProgression is false
	Progression  ProgressionExt
	Mark, Bareme int // student mark / exercice total
}

// LoadTaskProgression is a convenience wrapper around [LoadTasksProgression]
func LoadTaskProgression(db ta.DB, idStudent teacher.IdStudent, idTask ta.IdTask) (TaskProgressionHeader, error) {
	pr, err := LoadTasksProgression(db, idStudent, []ta.IdTask{idTask})
	if err != nil {
		return TaskProgressionHeader{}, err
	}
	return pr[idTask], nil
}

// LoadTasksProgression fetches the progression of one student against
// the given tasks.
func LoadTasksProgression(db ta.DB, idStudent teacher.IdStudent, idTasks []ta.IdTask) (map[ta.IdTask]TaskProgressionHeader, error) {
	contents, err := NewTasksContents(db, idTasks)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	links1, err := ta.SelectProgressionsByIdStudents(db, idStudent)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	progressionsByTask := links1.ByIdTask() // for one student

	out := make(map[ta.IdTask]TaskProgressionHeader, len(contents.Tasks))
	for _, task := range contents.Tasks {
		work := contents.GetWork(task)
		baremes := work.Bareme()
		progs := progressionsByTask[task.Id]
		// the progression may be empty if the student has not started it
		hasProg := len(progs) != 0
		progression := NewProgressionExt(progs, len(baremes))

		out[task.Id] = TaskProgressionHeader{
			Id:      task.Id,
			Title:   work.Title(),
			Chapter: work.Chapter(),

			HasProgression: hasProg,
			Progression:    progression,
			Bareme:         baremes.Total(),
			Mark:           baremes.ComputeMark(progression.Questions),
		}
	}

	return out, nil
}

// EvaluateTaskExercice calls `EvaluateExercice` and, if [registerProgression] is true, registers
// the student progression, returning the updated mark.
// If needed, a new progression item is created.
// If [registerProgression] is false, no progression is created.
func EvaluateTaskExercice(db *sql.DB, idTask ta.IdTask, idStudent teacher.IdStudent, ex EvaluateWorkIn, registerProgression bool) (out EvaluateWorkOut, mark int, err error) {
	out, err = ex.Evaluate(db, idStudent)
	if err != nil {
		return
	}

	if registerProgression {
		// persists the progression on DB
		err = updateProgression(db, idStudent, idTask, out.Progression.Questions)
		if err != nil {
			return out, 0, err
		}
	}

	// in any case compute the (new) mark
	task, err := ta.SelectTask(db, idTask)
	if err != nil {
		return out, mark, utils.SQLError(err)
	}

	loader, err := newWorkLoader(db, NewWorkID(task), idStudent)
	if err != nil {
		return
	}
	baremes := loader.Bareme()
	mark = baremes.ComputeMark(out.Progression.Questions)

	return out, mark, nil
}

// TaskBareme stores the baremes of a task, for each question.
type TaskBareme []int

// Total returns the bareme for one task, that is the sum
// of each question's bareme
func (bareme TaskBareme) Total() int {
	var out int
	for _, questionBareme := range bareme {
		out += questionBareme
	}
	return out
}

// ComputeMark computes the student mark.
// An empty [progression] is supported and returns 0.
// Otherwise, the length of [progression] must match the length of [bareme]
func (bareme TaskBareme) ComputeMark(progression []ta.QuestionHistory) int {
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
