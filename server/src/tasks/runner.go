// Package tasks implements the server instatiation/validation of question and exercices,
// and the logic needed to handle a session.
// It is build upon the data structures defined in sql/tasks and exposes its API
// via package level functions.
package tasks

import (
	"fmt"

	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/questions"
	"github.com/benoitkugler/maths-online/maths/questions/client"
	ed "github.com/benoitkugler/maths-online/sql/editor"
	ta "github.com/benoitkugler/maths-online/sql/tasks"
	"github.com/benoitkugler/maths-online/utils"
)

type InstantiatedQuestion struct {
	Id       ed.IdQuestion
	Question client.Question `gomacro-extern:"client#dart#questions/types.gen.dart"`
	Params   []VarEntry

	// this field is private, not exported as JSON,
	// and used to simplify the loopback logic
	instance questions.QuestionInstance
}

func (iq InstantiatedQuestion) Instance() questions.QuestionInstance { return iq.instance }

type Answer struct {
	Params []VarEntry
	Answer client.QuestionAnswersIn `gomacro-extern:"client#dart#questions/types.gen.dart"`
}

type InstantiateQuestionsOut []InstantiatedQuestion

type VarEntry struct {
	Variable expression.Variable
	Resolved string
}

func newVarList(vars expression.Vars) []VarEntry {
	varList := make([]VarEntry, 0, len(vars))
	for k, v := range vars {
		varList = append(varList, VarEntry{Variable: k, Resolved: v.Serialize()})
	}
	return varList
}

// InstantiateQuestions loads and instantiates the given questions,
// also returning the paramerters used to do so.
func InstantiateQuestions(db ed.DB, ids []ed.IdQuestion) (InstantiateQuestionsOut, error) {
	questions, err := ed.SelectQuestions(db, ids...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	out := make(InstantiateQuestionsOut, len(ids))
	for index, id := range ids {
		qu := questions[id]
		vars, err := qu.Page.Parameters.ToMap().Instantiate()
		if err != nil {
			return nil, err
		}
		instance, err := qu.Page.InstantiateWith(vars)
		if err != nil {
			return nil, err
		}
		out[index] = InstantiatedQuestion{
			Id:       id,
			Question: instance.ToClient(),
			Params:   newVarList(vars),
			instance: instance,
		}
	}

	return out, nil
}

type EvaluateQuestionIn struct {
	Answer     Answer
	IdQuestion ed.IdQuestion
}

// Evaluate instantiate the given question with the given parameters,
// and evaluate the given answer.
func (params EvaluateQuestionIn) Evaluate(db ed.DB) (client.QuestionAnswersOut, error) {
	qu, err := ed.SelectQuestion(db, params.IdQuestion)
	if err != nil {
		return client.QuestionAnswersOut{}, utils.SQLError(err)
	}
	return evaluateQuestion(qu, params.Answer)
}

func evaluateQuestion(qu ed.Question, answer Answer) (client.QuestionAnswersOut, error) {
	var err error
	paramsDict := make(expression.Vars)
	for _, entry := range answer.Params {
		paramsDict[entry.Variable], err = expression.Parse(entry.Resolved)
		if err != nil {
			return client.QuestionAnswersOut{}, err
		}
	}

	instance, err := qu.Page.InstantiateWith(paramsDict)
	if err != nil {
		return client.QuestionAnswersOut{}, err
	}

	return instance.EvaluateAnswer(answer.Answer), nil
}

// WorkID identifies either an exercice or a monoquestion
type WorkID struct {
	ID         int64
	IsExercice bool
}

func NewWorkID(task ta.Task) WorkID {
	if task.IdExercice.Valid {
		return newWorkIDFromEx(task.IdExercice.ID)
	}
	return newWorkIDFromMono(task.IdMonoquestion.ID)
}

func newWorkIDFromEx(id ed.IdExercice) WorkID { return WorkID{ID: int64(id), IsExercice: true} }

func newWorkIDFromMono(id ta.IdMonoquestion) WorkID { return WorkID{ID: int64(id), IsExercice: false} }

// Work is the common interface for exerices and mono-questions.
type Work interface {
	Title() string    // as presented to the student
	Subtitle() string // used by the teacher
	flow() ed.Flow
	QuestionsList() ([]ed.Question, TaskBareme)
	Instantiate() (InstantiatedWork, error)
}

func newWorkLoader(db ed.DB, work WorkID) (Work, error) {
	if work.IsExercice {
		return NewExerciceData(db, ed.IdExercice(work.ID))
	}
	return NewMonoquestionData(db, ta.IdMonoquestion(work.ID))
}

type InstantiatedWork struct {
	ID WorkID

	Title     string
	Flow      ed.Flow
	Questions []InstantiatedQuestion
	Baremes   []int
}

// InstantiateWork load an exercice (or a monoquestion) and its questions.
func InstantiateWork(db ed.DB, work WorkID) (InstantiatedWork, error) {
	loader, err := newWorkLoader(db, work)
	if err != nil {
		return InstantiatedWork{}, err
	}
	return loader.Instantiate()
}

func instantiateQuestions(questions []ed.Question, sharedVars expression.Vars) ([]InstantiatedQuestion, error) {
	out := make([]InstantiatedQuestion, len(questions))

	for index, question := range questions {
		ownVars, err := question.Page.Parameters.ToMap().Instantiate()
		if err != nil {
			return nil, err
		}

		if question.NeedExercice.Valid {
			// merge the parameters, given higher precedence to question
			for c, v := range sharedVars {
				if _, has := ownVars[c]; !has {
					ownVars[c] = v
				}
			}
		}

		instance, err := question.Page.InstantiateWith(ownVars)
		if err != nil {
			return nil, err
		}

		out[index] = InstantiatedQuestion{
			Id:       question.Id,
			Question: instance.ToClient(),
			Params:   newVarList(ownVars),
			instance: instance,
		}
	}

	return out, nil
}

// ExerciceData is an helper struct to unify question loading
// for an exercice.
type ExerciceData struct {
	Group        ed.Exercicegroup // the exercice group
	Exercice     ed.Exercice
	Links        ed.ExerciceQuestions
	QuestionsMap ed.Questions
}

// NewExerciceData loads the given exercice and the associated questions
func NewExerciceData(db ed.DB, id ed.IdExercice) (ExerciceData, error) {
	ex, err := ed.SelectExercice(db, id)
	if err != nil {
		return ExerciceData{}, utils.SQLError(err)
	}

	group, err := ed.SelectExercicegroup(db, ex.IdGroup)
	if err != nil {
		return ExerciceData{}, utils.SQLError(err)
	}

	links, err := ed.SelectExerciceQuestionsByIdExercices(db, id)
	if err != nil {
		return ExerciceData{}, utils.SQLError(err)
	}
	links.EnsureOrder()

	// load the question contents
	dict, err := ed.SelectQuestions(db, links.IdQuestions()...)
	if err != nil {
		return ExerciceData{}, utils.SQLError(err)
	}
	return ExerciceData{Group: group, Exercice: ex, Links: links, QuestionsMap: dict}, nil
}

func (ex ExerciceData) Title() string    { return ex.Group.Title }
func (ex ExerciceData) Subtitle() string { return ex.Exercice.Subtitle }
func (ExerciceData) flow() ed.Flow       { return ed.Sequencial }

// QuestionsList resolve the links list using `source`
func (ex ExerciceData) QuestionsList() ([]ed.Question, TaskBareme) {
	questions := make([]ed.Question, len(ex.Links))
	baremes := make([]int, len(ex.Links))
	for i, link := range ex.Links {
		questions[i] = ex.QuestionsMap[link.IdQuestion]
		baremes[i] = ex.Links[i].Bareme
	}
	return questions, baremes
}

// Instantiate instantiates the questions, using a fixed shared instance of the exercice parameters
// for each question
func (data ExerciceData) Instantiate() (InstantiatedWork, error) {
	ex := data.Exercice
	questions, baremes := data.QuestionsList()

	out := InstantiatedWork{
		ID:      newWorkIDFromEx(ex.Id),
		Title:   data.Title(),
		Flow:    data.flow(),
		Baremes: baremes,
	}

	// instantiate the questions :
	// start with the shared paremeters, which must be instantiated only once
	sharedVars, err := ex.Parameters.ToMap().Instantiate()
	if err != nil {
		return InstantiatedWork{}, err
	}

	out.Questions, err = instantiateQuestions(questions, sharedVars)
	if err != nil {
		return InstantiatedWork{}, err
	}

	return out, nil
}

type MonoquestionData struct {
	params        ta.Monoquestion
	questiongroup ed.Questiongroup
	question      ed.Question
}

// NewMonoquestionData loads the given Monoquestion and the associated question
func NewMonoquestionData(db ed.DB, id ta.IdMonoquestion) (out MonoquestionData, err error) {
	out.params, err = ta.SelectMonoquestion(db, id)
	if err != nil {
		return out, utils.SQLError(err)
	}

	out.question, err = ed.SelectQuestion(db, out.params.IdQuestion)
	if err != nil {
		return out, utils.SQLError(err)
	}

	out.questiongroup, err = ed.SelectQuestiongroup(db, out.question.IdGroup.ID)
	if err != nil {
		return out, utils.SQLError(err)
	}

	return out, nil
}

func (data MonoquestionData) Title() string    { return data.questiongroup.Title }
func (data MonoquestionData) Subtitle() string { return data.question.Subtitle }
func (MonoquestionData) flow() ed.Flow         { return ed.Parallel }

// QuestionsList returns the generated list of questions
func (data MonoquestionData) QuestionsList() ([]ed.Question, TaskBareme) {
	questions := make([]ed.Question, data.params.NbRepeat)
	baremes := make([]int, data.params.NbRepeat)
	// repeat the question
	for i := range questions {
		questions[i] = data.question
		baremes[i] = data.params.Bareme
	}
	return questions, baremes
}

func (data MonoquestionData) Instantiate() (InstantiatedWork, error) {
	questions, baremes := data.QuestionsList()
	out := InstantiatedWork{
		ID:      newWorkIDFromMono(data.params.Id),
		Title:   data.Title(),
		Flow:    ed.Parallel,
		Baremes: baremes,
	}

	var err error
	out.Questions, err = instantiateQuestions(questions, nil)
	if err != nil {
		return InstantiatedWork{}, err
	}
	return out, nil
}

type EvaluateWorkIn struct {
	ID WorkID

	Answers map[int]Answer // by question index (not ID)
	// the current progression, as send by the server,
	// to update with the given answers
	Progression ProgressionExt
}

type EvaluateWorkOut struct {
	Results      map[int]client.QuestionAnswersOut `gomacro-extern:"client#dart#questions/types.gen.dart"`
	Progression  ProgressionExt                    // the updated progression
	NewQuestions []InstantiatedQuestion            // only non empty if the answer is not correct
}

// Evaluate checks the answer provided for the given exercice and
// update the in-memory progression.
// The given progression must either be empty or have same length
// as the exercice.
func (args EvaluateWorkIn) Evaluate(db ed.DB) (EvaluateWorkOut, error) {
	data, err := newWorkLoader(db, args.ID)
	if err != nil {
		return EvaluateWorkOut{}, utils.SQLError(err)
	}

	qus, _ := data.QuestionsList()

	// handle initial empty progressions
	if len(args.Progression.Questions) == 0 {
		args.Progression.Questions = make([]ta.QuestionHistory, len(qus))
	}

	// enforce invariant for existing progressions
	if L1, L2 := len(qus), len(args.Progression.Questions); L1 != L2 {
		return EvaluateWorkOut{}, fmt.Errorf("internal error: inconsistent length %d != %d", L1, L2)
	}

	updatedProgression := args.Progression // shallow copy is enough
	results := make(map[int]client.QuestionAnswersOut)

	// depending on the flow, we either evaluate only one question,
	// or all the ones given
	switch data.flow() {
	case ed.Parallel: // all questions
		for questionIndex, question := range qus {
			if answer, hasAnswer := args.Answers[questionIndex]; hasAnswer {
				resp, err := evaluateQuestion(question, answer)
				if err != nil {
					return EvaluateWorkOut{}, err
				}

				results[questionIndex] = resp
				l := &updatedProgression.Questions[questionIndex]
				*l = append(*l, resp.IsCorrect())
			}
		}
	case ed.Sequencial: // only the current question
		questionIndex := args.Progression.NextQuestion
		if questionIndex < 0 || questionIndex >= len(qus) {
			return EvaluateWorkOut{}, fmt.Errorf("internal error: invalid question index %d", questionIndex)
		}

		answer, has := args.Answers[questionIndex]
		if !has {
			return EvaluateWorkOut{}, fmt.Errorf("internal error: missing answer for %d", questionIndex)
		}

		resp, err := evaluateQuestion(qus[questionIndex], answer)
		if err != nil {
			return EvaluateWorkOut{}, err
		}

		results[questionIndex] = resp
		l := &updatedProgression.Questions[questionIndex]
		*l = append(*l, resp.IsCorrect())
	}

	updatedProgression.InferNextQuestion() // update in case of success

	newVersion, err := data.Instantiate()
	if err != nil {
		return EvaluateWorkOut{}, err
	}

	return EvaluateWorkOut{Results: results, Progression: updatedProgression, NewQuestions: newVersion.Questions}, nil
}
