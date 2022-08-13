package editor

import (
	"fmt"

	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/questions/client"
	"github.com/benoitkugler/maths-online/utils"
)

// implementation of the server instatiation/validation of question/exercices

type InstantiatedQuestion struct {
	Id       IdQuestion
	Question client.Question `gomacro-extern:"client:dart:questions/types.gen.dart"`
	Params   []VarEntry
}

type Answer struct {
	Params []VarEntry
	Answer client.QuestionAnswersIn `gomacro-extern:"client:dart:questions/types.gen.dart"`
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
func (ct *Controller) InstantiateQuestions(ids []IdQuestion) (InstantiateQuestionsOut, error) {
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
func (ct *Controller) EvaluateQuestion(id IdQuestion, answer Answer) (client.QuestionAnswersOut, error) {
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

// Success return true if at least one try is sucessful
func (qh QuestionHistory) Success() bool {
	for _, try := range qh {
		if try {
			return true
		}
	}
	return false
}

type ProgressionExt struct {
	Questions    []QuestionHistory
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

type InstantiatedExercice struct {
	Exercice  Exercice
	Questions []InstantiatedQuestion
	Baremes   []int
}

// InstantiateExercice load an exercice and its questions.
func InstantiateExercice(db DB, id IdExercice) (InstantiatedExercice, error) {
	content, err := loadExercice(db, id)
	if err != nil {
		return InstantiatedExercice{}, err
	}
	return content.instantiate()
}

// helper to unify question loading
type exerciceContent struct {
	exercice Exercice
	links    ExerciceQuestions
	dict     Questions
}

func (ex exerciceContent) questions() []Question {
	out := make([]Question, len(ex.links))
	for i, link := range ex.links {
		out[i] = ex.dict[link.IdQuestion]
	}
	return out
}

func (ex exerciceContent) questionsSource(userID, adminID uID) map[IdQuestion]QuestionOrigin {
	out := make(map[IdQuestion]QuestionOrigin, len(ex.dict))
	for i, qu := range ex.dict {
		origin, _ := qu.origin(userID, adminID)
		out[i] = QuestionOrigin{Question: qu, Origin: origin}
	}
	return out
}

// instantiates the questions, using a fixed shared instance of the exercice parameters
// for each question
func (data exerciceContent) instantiate() (InstantiatedExercice, error) {
	ex, links, questions := data.exercice, data.links, data.questions()

	out := InstantiatedExercice{
		Exercice:  ex,
		Questions: make([]InstantiatedQuestion, len(questions)),
		Baremes:   make([]int, len(questions)),
	}

	// instantiate the questions :
	// start with the shared paremeters, which must be instantiated only once
	sharedVars, err := ex.Parameters.ToMap().Instantiate()
	if err != nil {
		return InstantiatedExercice{}, err
	}

	for index, question := range questions {
		ownVars, err := question.Page.Parameters.ToMap().Instantiate()
		if err != nil {
			return InstantiatedExercice{}, err
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
			return InstantiatedExercice{}, err
		}

		out.Questions[index] = InstantiatedQuestion{
			Id:       question.Id,
			Question: instance.ToClient(),
			Params:   newVarList(ownVars),
		}
		out.Baremes[index] = links[index].Bareme
	}

	return out, nil
}

// loadExercice loads the given exercice and the associated questions
func loadExercice(db DB, exerciceID IdExercice) (exerciceContent, error) {
	ex, err := SelectExercice(db, exerciceID)
	if err != nil {
		return exerciceContent{}, utils.SQLError(err)
	}
	links, err := SelectExerciceQuestionsByIdExercices(db, exerciceID)
	if err != nil {
		return exerciceContent{}, utils.SQLError(err)
	}
	links.EnsureOrder()

	// load the question contents
	dict, err := SelectQuestions(db, links.IdQuestions()...)
	if err != nil {
		return exerciceContent{}, utils.SQLError(err)
	}
	return exerciceContent{exercice: ex, links: links, dict: dict}, nil
}

type EvaluateExerciceIn struct {
	IdExercice IdExercice
	Answers    map[int]Answer // by question index (not ID)
	// the current progression, as send by the server,
	// to update with the given answers
	Progression ProgressionExt
}

type EvaluateExerciceOut struct {
	Results      map[int]client.QuestionAnswersOut `gomacro-extern:"client:dart:questions/types.gen.dart"`
	Progression  ProgressionExt                    // the updated progression
	NewQuestions []InstantiatedQuestion            // only non empty if the answer is not correct
}

func (ct *Controller) EvaluateExercice(args EvaluateExerciceIn) (EvaluateExerciceOut, error) {
	return EvaluateExercice(ct.db, args)
}

// EvaluateExercice checks the answer provided for the given exercice and
// update the in-memory progression.
// The given progression must either be empty or having same length
// as the exercice.
func EvaluateExercice(db DB, args EvaluateExerciceIn) (EvaluateExerciceOut, error) {
	data, err := loadExercice(db, args.IdExercice)
	if err != nil {
		return EvaluateExerciceOut{}, utils.SQLError(err)
	}
	ex, qus := data.exercice, data.questions()

	// handle initial empty progressions
	if len(args.Progression.Questions) == 0 {
		args.Progression.Questions = make([]QuestionHistory, len(qus))
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

	updatedProgression.InferNextQuestion() // update in case of success

	newVersion, err := data.instantiate()
	if err != nil {
		return EvaluateExerciceOut{}, err
	}

	return EvaluateExerciceOut{Results: results, Progression: updatedProgression, NewQuestions: newVersion.Questions}, nil
}
