package editor

import (
	"fmt"

	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/questions/client"
	"github.com/benoitkugler/maths-online/utils"
)

// implementation of the server instatiation/validation of question/exercices

type InstantiatedQuestion struct {
	Id       int64
	Question client.Question `dart-extern:"client:questions/types.gen.dart"`
	Params   []VarEntry
}

type Answer struct {
	Params []VarEntry
	Answer client.QuestionAnswersIn `dart-extern:"client:questions/types.gen.dart"`
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
func (ct *Controller) EvaluateQuestion(id int64, answer Answer) (client.QuestionAnswersOut, error) {
	qu, err := SelectQuestion(ct.db, id)
	if err != nil {
		return client.QuestionAnswersOut{}, utils.SQLError(err)
	}
	return qu.evaluate(answer)
}

func (qu Question) evaluate(answer Answer) (client.QuestionAnswersOut, error) {
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

type InstantiatedExercice struct {
	Id        int64
	Title     string
	Flow      Flow
	Questions []InstantiatedQuestion
	Baremes   []int
}

func (ct *Controller) loadExercice(id int64) (Exercice, ExerciceQuestions, []Question, error) {
	ex, err := SelectExercice(ct.db, id)
	if err != nil {
		return Exercice{}, nil, nil, utils.SQLError(err)
	}
	links, err := SelectExerciceQuestionsByIdExercices(ct.db, id)
	if err != nil {
		return Exercice{}, nil, nil, utils.SQLError(err)
	}
	links.ensureIndex()

	// load the question contents
	qus, err := SelectQuestions(ct.db, links.IdQuestions()...)
	if err != nil {
		return Exercice{}, nil, nil, utils.SQLError(err)
	}

	out := make([]Question, len(links))
	for i, link := range links {
		out[i] = qus[link.IdQuestion]
	}

	return ex, links, out, nil
}

// instantiateExercice loads the given exercice, the associated questions,
// and instantiates them with the same random parameters
func (ct *Controller) instantiateExercice(id int64) (InstantiatedExercice, error) {
	ex, links, qus, err := ct.loadExercice(id)
	if err != nil {
		return InstantiatedExercice{}, err
	}

	out := InstantiatedExercice{
		Id:        ex.Id,
		Title:     ex.Title,
		Flow:      ex.Flow,
		Questions: make([]InstantiatedQuestion, len(qus)),
		Baremes:   make([]int, len(qus)),
	}

	// instantiate the questions
	commonParams := ex.Parameters.ToMap()
	for index, question := range qus {
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
			Id:       question.Id,
			Question: instance.ToClient(),
			Params:   newVarList(vars),
		}
		out.Baremes[index] = links[index].Bareme
	}

	return out, nil
}

type EvaluateExerciceIn struct {
	IdExercice int64
	Answers    map[int]Answer // by question index (not ID)
	// the current progression, as send by the server,
	// to update with the given answers
	Progression ProgressionExt
}

type EvaluateExerciceOut struct {
	Results      map[int]client.QuestionAnswersOut `dart-extern:"client:questions/types.gen.dart"`
	Progression  ProgressionExt                    // the updated progression
	NewQuestions []InstantiatedQuestion            // only non empty if the answer is not correct
}

// EvaluateExercice checks the answer provided for the given exercice and
// update the progression.
func (ct *Controller) EvaluateExercice(args EvaluateExerciceIn) (EvaluateExerciceOut, error) {
	ex, _, qus, err := ct.loadExercice(args.IdExercice)
	if err != nil {
		return EvaluateExerciceOut{}, utils.SQLError(err)
	}

	if L1, L2 := len(qus), len(args.Progression.Questions); L1 != L2 {
		return EvaluateExerciceOut{}, fmt.Errorf("internal error: inconsistent length %d != %d", L1, L2)
	}

	updatedProgression := args.Progression // shallow copy is enough
	results := make(map[int]client.QuestionAnswersOut)

	// depending on the flow, we either evaluate only one question,
	// or all the ones given
	switch ex.Flow {
	case Parallel: // all questions
		for questionIndex, question := range qus {
			if answer, hasAnswer := args.Answers[questionIndex]; hasAnswer {
				resp, err := question.evaluate(answer)
				if err != nil {
					return EvaluateExerciceOut{}, err
				}

				results[questionIndex] = resp
				l := &updatedProgression.Questions[questionIndex]
				*l = append(*l, resp.IsCorrect())
			}
		}
	case Sequencial: // only the current question
		questionIndex := args.Progression.NextQuestion
		if questionIndex < 0 || questionIndex >= len(qus) {
			return EvaluateExerciceOut{}, fmt.Errorf("internal error: invalid question index %d", questionIndex)
		}

		answer, has := args.Answers[questionIndex]
		if !has {
			return EvaluateExerciceOut{}, fmt.Errorf("internal error: missing answer for %d", questionIndex)
		}

		resp, err := qus[questionIndex].evaluate(answer)
		if err != nil {
			return EvaluateExerciceOut{}, err
		}

		results[questionIndex] = resp
		l := &updatedProgression.Questions[questionIndex]
		*l = append(*l, resp.IsCorrect())
	}

	updatedProgression.inferNextQuestion() // update in case of success

	newVersion, err := ct.instantiateExercice(args.IdExercice)
	if err != nil {
		return EvaluateExerciceOut{}, err
	}

	return EvaluateExerciceOut{Results: results, Progression: updatedProgression, NewQuestions: newVersion.Questions}, nil
}
