package editor

import (
	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/questions/client"
	"github.com/benoitkugler/maths-online/utils"
)

// implementation of the server instatiation/validation of question/exercices

type InstantiatedQuestion struct {
	Id       int64
	Question client.Question `dart-extern:"exercices/types.gen.dart"`
	Params   []VarEntry
}

type InstantiateQuestionsOut []InstantiatedQuestion

type InstantiatedExercice struct {
	Exercice  Exercice
	Questions []InstantiatedQuestion
}

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
func (ct *Controller) InstantiateQuestions(ids []int64) (InstantiateQuestionsOut, error) {
	questions, err := SelectQuestions(ct.db, ids...)
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
		}
	}

	return out, nil
}

// EvaluateQuestion instantiate the given question with the given parameters,
// and evaluate the given answer.
func (ct *Controller) EvaluateQuestion(id int64, params []VarEntry, answer client.QuestionAnswersIn) (client.QuestionAnswersOut, error) {
	qu, err := SelectQuestion(ct.db, id)
	if err != nil {
		return client.QuestionAnswersOut{}, utils.SQLError(err)
	}

	paramsDict := make(expression.Vars)
	for _, entry := range params {
		paramsDict[entry.Variable], err = expression.Parse(entry.Resolved)
		if err != nil {
			return client.QuestionAnswersOut{}, err
		}
	}

	instance, err := qu.Page.InstantiateWith(paramsDict)
	if err != nil {
		return client.QuestionAnswersOut{}, err
	}

	return instance.EvaluateAnswer(answer), nil
}

// QuestionHistory stores the successes for one question,
// in chronological order.
// For instance, [true, false, true] means : first try: correct, second: wrong answer,third: correct
type QuestionHistory []bool

// success return true if the last try is sucesseful
func (qh QuestionHistory) success() bool {
	return len(qh) > 0 && qh[len(qh)-1]
}

type ProgressionExt struct {
	Progression  Progression
	Questions    []QuestionHistory
	NextQuestion int
}

// inferNextQuestion stores the first question not passed by the student
// note that only the last try is taken in account
// if all the questions are successul, it set it to -1
func (qh *ProgressionExt) inferNextQuestion() {
	for i, question := range qh.Questions {
		if !question.success() {
			qh.NextQuestion = i
			return
		}
	}
	qh.NextQuestion = -1
}

// load the whole progression
func (ct *Controller) fetchProgression(id int64) (ProgressionExt, error) {
	pr, err := SelectProgression(ct.db, id)
	if err != nil {
		return ProgressionExt{}, utils.SQLError(err)
	}

	questions, err := SelectExerciceQuestionsByIdExercices(ct.db, pr.IdExercice)
	if err != nil {
		return ProgressionExt{}, utils.SQLError(err)
	}
	questions.ensureIndex()

	links, err := SelectProgressionQuestionsByIdProgressions(ct.db, id)
	if err != nil {
		return ProgressionExt{}, utils.SQLError(err)
	}
	// beware that some questions may not have a link item yet
	out := ProgressionExt{
		Progression: pr,
		Questions:   make([]QuestionHistory, len(questions)),
	}
	for _, link := range links {
		out.Questions[link.Index] = link.History
	}
	out.inferNextQuestion()
	return out, nil
}

// instantiateExercice loads the given exercice, the associated questions,
// and instantiates them with the same random parameters
func (ct *Controller) instantiateExercice(id int64) (InstantiatedExercice, error) {
	ex, err := SelectExercice(ct.db, id)
	if err != nil {
		return InstantiatedExercice{}, utils.SQLError(err)
	}
	links, err := SelectExerciceQuestionsByIdExercices(ct.db, id)
	if err != nil {
		return InstantiatedExercice{}, utils.SQLError(err)
	}
	links.ensureIndex()

	// load the question contents
	qus, err := SelectQuestions(ct.db, links.IdQuestions()...)
	if err != nil {
		return InstantiatedExercice{}, utils.SQLError(err)
	}

	out := InstantiatedExercice{
		Exercice:  ex,
		Questions: make([]InstantiatedQuestion, len(links)),
	}

	// instantiate the questions
	commonParams := ex.Parameters.ToMap()
	for index, link := range links {
		question := qus[link.IdQuestion]
		ownParams := question.Page.Parameters.ToMap()

		// merge the parameters, given higher precedence to question
		for c, v := range commonParams {
			if _, has := ownParams[c]; !has {
				ownParams[c] = v
			}
		}

		vars, err := ownParams.Instantiate()
		if err != nil {
			return InstantiatedExercice{}, err
		}

		instance, err := question.Page.InstantiateWith(vars)
		if err != nil {
			return InstantiatedExercice{}, err
		}

		out.Questions[index] = InstantiatedQuestion{
			Id:       link.IdQuestion,
			Question: instance.ToClient(),
			Params:   newVarList(vars),
		}
	}

	return out, nil
}
