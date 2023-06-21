// Package tasks implements the server instatiation/validation of question and exercices,
// and the logic needed to handle a session.
// It is build upon the data structures defined in sql/tasks and exposes its API
// via package level functions.
package tasks

import (
	"errors"
	"fmt"

	"github.com/benoitkugler/maths-online/server/src/maths/expression"
	"github.com/benoitkugler/maths-online/server/src/maths/questions"
	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	ed "github.com/benoitkugler/maths-online/server/src/sql/editor"
	ta "github.com/benoitkugler/maths-online/server/src/sql/tasks"
	tc "github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/benoitkugler/maths-online/server/src/utils"
)

type InstantiatedQuestion struct {
	Id       ed.IdQuestion
	Question client.Question
	Params   Params
}

type AnswerP struct {
	Params Params
	Answer client.QuestionAnswersIn
}

type InstantiateQuestionsOut []InstantiatedQuestion

type VarEntry struct {
	Variable expression.Variable
	Resolved string
}

// Params is a serialized version of [expression.Vars],
// used by clients.
type Params []VarEntry

// NewParams serialize the given map.
func NewParams(vars expression.Vars) []VarEntry {
	varList := make([]VarEntry, 0, len(vars))
	for k, v := range vars {
		varList = append(varList, VarEntry{Variable: k, Resolved: v.Serialize()})
	}
	return varList
}

// ToMap parse the [Params]
func (params Params) ToMap() (expression.Vars, error) {
	paramsDict := make(expression.Vars)
	for _, entry := range params {
		var err error
		paramsDict[entry.Variable], err = expression.Parse(entry.Resolved)
		if err != nil {
			return nil, err
		}
	}
	return paramsDict, nil
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
		instance, vars, err := qu.Page().InstantiateErr()
		if err != nil {
			return nil, err
		}
		out[index] = InstantiatedQuestion{
			Id:       id,
			Question: instance.ToClient(),
			Params:   NewParams(vars),
		}
	}

	return out, nil
}

type EvaluateQuestionIn struct {
	Answer     AnswerP
	IdQuestion ed.IdQuestion
}

// Evaluate instantiate the given question with the given parameters,
// and evaluate the given answer.
func (params EvaluateQuestionIn) Evaluate(db ed.DB) (client.QuestionAnswersOut, error) {
	qu, err := ed.SelectQuestion(db, params.IdQuestion)
	if err != nil {
		return client.QuestionAnswersOut{}, utils.SQLError(err)
	}
	return EvaluateQuestion(qu.Enonce, params.Answer)
}

// EvaluateQuestion instantiate [qu] against the given [answer.Params]
// and evaluate the given [answer.Answer]
func EvaluateQuestion(qu questions.Enonce, answer AnswerP) (client.QuestionAnswersOut, error) {
	paramsDict, err := answer.Params.ToMap()
	if err != nil {
		return client.QuestionAnswersOut{}, err
	}

	instance, err := qu.InstantiateWith(paramsDict)
	if err != nil {
		return client.QuestionAnswersOut{}, err
	}

	return instance.EvaluateAnswer(answer.Answer), nil
}

const (
	WorkExercice uint8 = iota
	WorkMonoquestion
	WorkRandomMonoquestion
)

// WorkID identifies either an exercice or a (possibly random) monoquestion
type WorkID struct {
	ID   int64
	Kind uint8
}

func NewWorkID(task ta.Task) WorkID {
	switch {
	case task.IdExercice.Valid:
		return newWorkIDFromEx(task.IdExercice.ID)
	case task.IdMonoquestion.Valid:
		return newWorkIDFromMono(task.IdMonoquestion.ID)
	case task.IdRandomMonoquestion.Valid:
		return newWorkIDFromRandomMono(task.IdRandomMonoquestion.ID)
	default:
		panic("unexpected task target")
	}
}

func newWorkIDFromEx(id ed.IdExercice) WorkID { return WorkID{ID: int64(id), Kind: WorkExercice} }

func newWorkIDFromMono(id ta.IdMonoquestion) WorkID {
	return WorkID{ID: int64(id), Kind: WorkMonoquestion}
}

func newWorkIDFromRandomMono(id ta.IdRandomMonoquestion) WorkID {
	return WorkID{ID: int64(id), Kind: WorkRandomMonoquestion}
}

// Work is the common interface for exercices and mono-questions.
type Work interface {
	Title() string // as presented to the student
	flow() ed.Flow
	QuestionsList() ([]ed.Question, TaskBareme)
	Instantiate() (InstantiatedWork, error)
}

// student is only required for RandomMonoquestion
func newWorkLoader(db ed.DB, work WorkID, student tc.IdStudent) (Work, error) {
	switch work.Kind {
	case WorkExercice:
		return NewExerciceData(db, ed.IdExercice(work.ID))
	case WorkMonoquestion:
		return newMonoquestionData(db, ta.IdMonoquestion(work.ID))
	case WorkRandomMonoquestion:
		return newRandomMonoquestionData(db, ta.IdRandomMonoquestion(work.ID), student)
	default:
		return nil, errors.New("internal error: unexpected Work kind")
	}
}

// InstantiatedWork is an instance of an exercice, more precisely
// of a full Exercice or Monoquestion (one question duplicated)
type InstantiatedWork struct {
	ID WorkID

	Title     string
	Flow      ed.Flow
	Questions []InstantiatedQuestion
	Baremes   []int
}

// InstantiateWork load an exercice (or a monoquestion) and its questions.
func InstantiateWork(db ed.DB, work WorkID, student tc.IdStudent) (InstantiatedWork, error) {
	loader, err := newWorkLoader(db, work, student)
	if err != nil {
		return InstantiatedWork{}, err
	}
	return loader.Instantiate()
}

func instantiateQuestions(questions []ed.Question, sharedVars expression.Vars) ([]InstantiatedQuestion, error) {
	out := make([]InstantiatedQuestion, len(questions))

	for index, question := range questions {
		ownVars, err := question.Parameters.ToMap().Instantiate()
		if err != nil {
			return nil, err
		}

		if question.NeedExercice.Valid {
			// merge the parameters, given higher precedence to question
			ownVars.CompleteFrom(sharedVars)
		}

		instance, err := question.Page().InstantiateWith(ownVars)
		if err != nil {
			return nil, err
		}

		out[index] = InstantiatedQuestion{
			Id:       question.Id,
			Question: instance.ToClient(),
			Params:   NewParams(ownVars),
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

func (ex ExerciceData) Title() string { return ex.Group.Title }
func (ExerciceData) flow() ed.Flow    { return ed.Sequencial }

// QuestionsList resolve the links list using `source`,
// returning lists of length `len(Links)`
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

type monoquestionData struct {
	params        ta.Monoquestion
	questiongroup ed.Questiongroup
	question      ed.Question
}

// newMonoquestionData loads the given Monoquestion and the associated question
func newMonoquestionData(db ed.DB, id ta.IdMonoquestion) (out monoquestionData, err error) {
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

func (data monoquestionData) Title() string { return data.questiongroup.Title }
func (monoquestionData) flow() ed.Flow      { return ed.Parallel }

// QuestionsList returns the generated list of questions
func (data monoquestionData) QuestionsList() ([]ed.Question, TaskBareme) {
	questions := make([]ed.Question, data.params.NbRepeat)
	baremes := make([]int, data.params.NbRepeat)
	// repeat the question
	for i := range questions {
		questions[i] = data.question
		baremes[i] = data.params.Bareme
	}
	return questions, baremes
}

func (data monoquestionData) Instantiate() (InstantiatedWork, error) {
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

type randomMonoquestionData struct {
	params            ta.RandomMonoquestion
	questiongroup     ed.Questiongroup
	selectedQuestions []ed.Question // with length params.NbRepeat
}

// assume a progression is registred for this student
func newRandomMonoquestionData(db ed.DB, id ta.IdRandomMonoquestion, student tc.IdStudent) (out randomMonoquestionData, err error) {
	out.params, err = ta.SelectRandomMonoquestion(db, id)
	if err != nil {
		return out, utils.SQLError(err)
	}

	out.questiongroup, err = ed.SelectQuestiongroup(db, out.params.IdQuestiongroup)
	if err != nil {
		return out, utils.SQLError(err)
	}

	variants, err := ta.SelectRandomMonoquestionVariantsByIdRandomMonoquestions(db, id)
	if err != nil {
		return out, utils.SQLError(err)
	}
	variants = variants.ByIdStudent()[student]
	variants.EnsureOrder()

	if len(variants) != out.params.NbRepeat {
		return out, errors.New("internal error: inconsitent length of variant for RandomMonoquestion")
	}

	dict, err := ed.SelectQuestions(db, variants.IdQuestions()...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	out.selectedQuestions = make([]ed.Question, len(variants))
	for i, v := range variants {
		out.selectedQuestions[i] = dict[v.IdQuestion]
	}

	return out, nil
}

func (data randomMonoquestionData) Title() string { return data.questiongroup.Title }
func (randomMonoquestionData) flow() ed.Flow      { return ed.Parallel }

// QuestionsList returns the generated list of questions
func (data randomMonoquestionData) QuestionsList() ([]ed.Question, TaskBareme) {
	baremes := make([]int, len(data.selectedQuestions))
	// repeat the bareme
	for i := range baremes {
		baremes[i] = data.params.Bareme
	}
	return data.selectedQuestions, baremes
}

func (data randomMonoquestionData) Instantiate() (InstantiatedWork, error) {
	questions, baremes := data.QuestionsList()
	out := InstantiatedWork{
		ID:      newWorkIDFromRandomMono(data.params.Id),
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

// ------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------
// ------------------------------------------------------------------------------------

type EvaluateWorkIn struct {
	ID WorkID

	IdStudent tc.IdStudent

	Answers map[int]AnswerP // by question index (not ID)

	// the current progression, as send by the server,
	// to update with the given answers
	Progression ProgressionExt
}

type EvaluateWorkOut struct {
	Results      map[int]client.QuestionAnswersOut
	Progression  ProgressionExt         // the updated progression
	NewQuestions []InstantiatedQuestion // only non empty if the answer is not correct
}

// Evaluate checks the answer provided for the given exercice and
// update the in-memory progression.
// The given progression must either be empty or have same length
// as the exercice.
func (args EvaluateWorkIn) Evaluate(db ed.DB) (EvaluateWorkOut, error) {
	data, err := newWorkLoader(db, args.ID, args.IdStudent)
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
				resp, err := EvaluateQuestion(question.Enonce, answer)
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

		resp, err := EvaluateQuestion(qus[questionIndex].Enonce, answer)
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
