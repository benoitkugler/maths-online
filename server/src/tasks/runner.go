// Package tasks implements the server instatiation/validation of question and exercices,
// and the logic needed to handle a session.
// It is build upon the data structures defined in sql/tasks and exposes its API
// via package level functions.
package tasks

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"sort"

	"github.com/benoitkugler/maths-online/server/src/maths/expression"
	"github.com/benoitkugler/maths-online/server/src/maths/questions"
	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	ce "github.com/benoitkugler/maths-online/server/src/sql/ceintures"
	ed "github.com/benoitkugler/maths-online/server/src/sql/editor"
	ta "github.com/benoitkugler/maths-online/server/src/sql/tasks"
	tc "github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/benoitkugler/maths-online/server/src/utils"
)

type InstantiatedQuestion struct {
	Id         ed.IdQuestion
	Question   client.Question
	Difficulty ed.DifficultyTag
	Params     Params
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
func NewParams(vars expression.Vars) Params {
	varList := make(Params, 0, len(vars))
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
			Id:         id,
			Question:   instance.ToClient(),
			Difficulty: qu.Difficulty,
			Params:     NewParams(vars),
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

type WorkKind uint8

const (
	WorkExercice WorkKind = iota + 1
	WorkMonoquestion
	WorkRandomMonoquestion
)

// WorkID identifies either an exercice or a (possibly random) monoquestion
type WorkID struct {
	ID   int64
	Kind WorkKind

	IsExercice bool // deprecated in version 1.5
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

func newWorkIDFromEx(id ed.IdExercice) WorkID {
	return WorkID{ID: int64(id), IsExercice: true, Kind: WorkExercice}
}

func newWorkIDFromMono(id ta.IdMonoquestion) WorkID {
	return WorkID{ID: int64(id), Kind: WorkMonoquestion}
}

func newWorkIDFromRandomMono(id ta.IdRandomMonoquestion) WorkID {
	return WorkID{ID: int64(id), Kind: WorkRandomMonoquestion}
}

// Work is the common interface for exercices and mono-questions.
type WorkMeta interface {
	Title() string // as presented to the student
	flow() ed.Flow
	Bareme() TaskBareme
	Chapter() string
}

type Work interface {
	WorkMeta
	Questions() []ed.Question
	Instantiate() (InstantiatedWork, error)
}

// student is only required for RandomMonoquestion
func newWorkLoader(db ed.DB, work WorkID, student tc.IdStudent) (Work, error) {
	if work.Kind == 0 { // backward compatiblity
		if work.IsExercice {
			work.Kind = WorkExercice
		} else {
			work.Kind = WorkMonoquestion
		}
	}

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
// For new RandomMonoquestions, the actual list of questions is also generated and saved.
func InstantiateWork(db *sql.DB, work WorkID, student tc.IdStudent) (InstantiatedWork, error) {
	loader, err := newWorkLoader(db, work, student)
	if err != nil {
		return InstantiatedWork{}, err
	}

	if random, ok := loader.(RandomMonoquestionData); ok && len(random.selectedQuestions) == 0 {
		loader, err = random.selectQuestions(db, student)
		if err != nil {
			return InstantiatedWork{}, err
		}
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
			Id:         question.Id,
			Question:   instance.ToClient(),
			Difficulty: question.Difficulty,
			Params:     NewParams(ownVars),
		}
	}

	return out, nil
}

// ExerciceData is an helper struct to unify question loading
// for an exercice.
type ExerciceData struct {
	Group   ed.Exercicegroup // the exercice group
	chapter string

	Exercice     ed.Exercice
	links        ed.ExerciceQuestions
	QuestionsMap ed.Questions
}

// NewExerciceData loads the given exercice and the associated questions
func NewExerciceData(db ed.DB, id ed.IdExercice) (out ExerciceData, _ error) {
	ex, err := ed.SelectExercice(db, id)
	if err != nil {
		return out, utils.SQLError(err)
	}

	group, err := ed.SelectExercicegroup(db, ex.IdGroup)
	if err != nil {
		return out, utils.SQLError(err)
	}

	tags, err := ed.SelectExercicegroupTagsByIdExercicegroups(db, group.Id)
	if err != nil {
		return out, utils.SQLError(err)
	}
	chapter := tags.Tags().BySection().Chapter

	links, err := ed.SelectExerciceQuestionsByIdExercices(db, id)
	if err != nil {
		return out, utils.SQLError(err)
	}
	links.EnsureOrder()

	// load the question contents
	dict, err := ed.SelectQuestions(db, links.IdQuestions()...)
	if err != nil {
		return out, utils.SQLError(err)
	}
	return ExerciceData{
		Group:        group,
		chapter:      chapter,
		Exercice:     ex,
		links:        links,
		QuestionsMap: dict,
	}, nil
}

func (ex ExerciceData) Title() string   { return ex.Group.Title }
func (ExerciceData) flow() ed.Flow      { return ed.Sequencial }
func (ex ExerciceData) Chapter() string { return ex.chapter }

func (ex ExerciceData) Questions() []ed.Question {
	questions := make([]ed.Question, len(ex.links))
	for i, link := range ex.links {
		questions[i] = ex.QuestionsMap[link.IdQuestion]
		// copy the exercice difficulty
		questions[i].Difficulty = ex.Exercice.Difficulty
	}
	return questions
}

func (ex ExerciceData) Bareme() TaskBareme {
	baremes := make([]int, len(ex.links))
	for i, link := range ex.links {
		baremes[i] = link.Bareme
	}
	return baremes
}

// Instantiate instantiates the questions, using a fixed shared instance of the exercice parameters
// for each question
func (data ExerciceData) Instantiate() (InstantiatedWork, error) {
	ex := data.Exercice
	questions := data.Questions()

	out := InstantiatedWork{
		ID:      newWorkIDFromEx(ex.Id),
		Title:   data.Title(),
		Flow:    data.flow(),
		Baremes: data.Bareme(),
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
	params   ta.Monoquestion
	Group    ed.Questiongroup
	chapter  string
	question ed.Question
}

// newMonoquestionData loads the given Monoquestion and the associated question
func newMonoquestionData(db ed.DB, id ta.IdMonoquestion) (out MonoquestionData, err error) {
	out.params, err = ta.SelectMonoquestion(db, id)
	if err != nil {
		return out, utils.SQLError(err)
	}

	out.question, err = ed.SelectQuestion(db, out.params.IdQuestion)
	if err != nil {
		return out, utils.SQLError(err)
	}

	out.Group, err = ed.SelectQuestiongroup(db, out.question.IdGroup.ID)
	if err != nil {
		return out, utils.SQLError(err)
	}

	tags, err := ed.SelectQuestiongroupTagsByIdQuestiongroups(db, out.Group.Id)
	if err != nil {
		return out, utils.SQLError(err)
	}
	out.chapter = tags.Tags().BySection().Chapter

	return out, nil
}

func (data MonoquestionData) Title() string   { return data.Group.Title }
func (MonoquestionData) flow() ed.Flow        { return ed.Parallel }
func (data MonoquestionData) Chapter() string { return data.chapter }

// Questions returns the generated list of questions
func (data MonoquestionData) Questions() []ed.Question {
	questions := make([]ed.Question, data.params.NbRepeat)
	// repeat the question
	for i := range questions {
		questions[i] = data.question
	}
	return questions
}

func (data MonoquestionData) Bareme() TaskBareme {
	baremes := make([]int, data.params.NbRepeat)
	// repeat the question
	for i := range baremes {
		baremes[i] = data.params.Bareme
	}
	return baremes
}

func (data MonoquestionData) Instantiate() (InstantiatedWork, error) {
	questions := data.Questions()
	out := InstantiatedWork{
		ID:      newWorkIDFromMono(data.params.Id),
		Title:   data.Title(),
		Flow:    data.flow(),
		Baremes: data.Bareme(),
	}

	var err error
	out.Questions, err = instantiateQuestions(questions, nil)
	if err != nil {
		return InstantiatedWork{}, err
	}
	return out, nil
}

type RandomMonoquestionData struct {
	params  ta.RandomMonoquestion
	Group   ed.Questiongroup
	chapter string

	// with length params.NbRepeat, or empty before
	// instanciation
	selectedQuestions []ed.Question
}

// `selectedQuestions` may be empty if the student has not instanciated this task yet
// See also [selectQuestions]
func newRandomMonoquestionData(db ed.DB, id ta.IdRandomMonoquestion, student tc.IdStudent) (out RandomMonoquestionData, err error) {
	out.params, err = ta.SelectRandomMonoquestion(db, id)
	if err != nil {
		return out, utils.SQLError(err)
	}

	out.Group, err = ed.SelectQuestiongroup(db, out.params.IdQuestiongroup)
	if err != nil {
		return out, utils.SQLError(err)
	}

	tags, err := ed.SelectQuestiongroupTagsByIdQuestiongroups(db, out.Group.Id)
	if err != nil {
		return out, utils.SQLError(err)
	}
	out.chapter = tags.Tags().BySection().Chapter

	variants, err := ta.SelectRandomMonoquestionVariantsByIdRandomMonoquestions(db, id)
	if err != nil {
		return out, utils.SQLError(err)
	}
	variants = variants.ByIdStudent()[student]
	variants.EnsureOrder()

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

// selectQuestions chooses the variants (with respect to the teacher settings), update the DB,
// and returns the updated struct
func (data RandomMonoquestionData) selectQuestions(db *sql.DB, idStudent tc.IdStudent) (RandomMonoquestionData, error) {
	// load all the variants
	questions, err := ed.SelectQuestionsByIdGroups(db, data.params.IdQuestiongroup)
	if err != nil {
		return data, utils.SQLError(err)
	}

	var filtered []ed.Question
	for _, qu := range questions {
		if data.params.Difficulty.Match(qu.Difficulty) {
			filtered = append(filtered, qu)
		}
	}

	if len(filtered) == 0 {
		// this should not happen, since the prof. API has a check for it
		return data, errors.New("Aucune question n'est disponible pour ce travail !")
	}

	data.selectedQuestions = selectVariants(data.params.NbRepeat, filtered)

	links := make(ta.RandomMonoquestionVariants, len(data.selectedQuestions))
	for i, qu := range data.selectedQuestions {
		links[i] = ta.RandomMonoquestionVariant{
			IdStudent:            idStudent,
			IdRandomMonoquestion: data.params.Id,
			Index:                int16(i),
			IdQuestion:           qu.Id,
		}
	}
	tx, err := db.Begin()
	if err != nil {
		return data, utils.SQLError(err)
	}
	err = ta.InsertManyRandomMonoquestionVariants(tx, links...)
	if err != nil {
		_ = tx.Rollback()
		return data, utils.SQLError(err)
	}
	err = tx.Commit()
	if err != nil {
		return data, utils.SQLError(err)
	}

	return data, nil
}

// assume len(among) > 0
func selectVariants(nbToSelect int, among []ed.Question) []ed.Question {
	selected := make([]ed.Question, 0, nbToSelect)

	nbAvail := len(among)
	// we prioritize diversity
	quotient, remainder := nbToSelect/nbAvail, nbToSelect%nbAvail
	for i := 0; i < quotient; i++ {
		selected = append(selected, among...) // repeat all questions
	}

	perm := rand.Perm(nbAvail)[:remainder]
	for _, index := range perm {
		selected = append(selected, among[index])
	}

	// sort by difficulty
	sort.Slice(selected, func(i, j int) bool { return selected[i].Difficulty < selected[j].Difficulty })

	return selected
}

func (data RandomMonoquestionData) Title() string   { return data.Group.Title }
func (RandomMonoquestionData) flow() ed.Flow        { return ed.Parallel }
func (data RandomMonoquestionData) Chapter() string { return data.chapter }

// Questions is only valid when the student specific variants have been loaded
func (data RandomMonoquestionData) Questions() []ed.Question { return data.selectedQuestions }

func (data RandomMonoquestionData) Bareme() TaskBareme {
	baremes := make([]int, data.params.NbRepeat)
	// repeat the bareme
	for i := range baremes {
		baremes[i] = data.params.Bareme
	}
	return baremes
}

func (data RandomMonoquestionData) Instantiate() (InstantiatedWork, error) {
	questions := data.Questions()
	out := InstantiatedWork{
		ID:      newWorkIDFromRandomMono(data.params.Id),
		Title:   data.Title(),
		Flow:    data.flow(),
		Baremes: data.Bareme(),
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

	// the current progression, as send by the server,
	// to update with the given answers
	Progression ProgressionExt

	AnswerIndex int     // new in v1.7
	Answer      AnswerP // new in v1.7

	// Deprecated
	Answers map[int]AnswerP `gomacro:"ignore"` // by question index (not ID)
}

func (ew *EvaluateWorkIn) fillFromMap() error {
	// client is using new API
	if len(ew.Answers) == 0 {
		return nil
	}

	if len(ew.Answers) != 1 {
		return errors.New("internal error: expected only one question")
	}
	for k, v := range ew.Answers {
		ew.AnswerIndex = k
		ew.Answer = v
	}
	return nil
}

type EvaluateWorkOut struct {
	Progression  ProgressionExt         // the updated progression
	NewQuestions []InstantiatedQuestion // only non empty if the answer is not correct

	AnswerIndex int
	Result      client.QuestionAnswersOut

	// Deprecated
	Results map[int]client.QuestionAnswersOut `gomacro:"ignore"`
}

func (ew *EvaluateWorkOut) fillMap() {
	ew.Results = map[int]client.QuestionAnswersOut{ew.AnswerIndex: ew.Result}
}

// Evaluate checks the answer provided for the given exercice and
// update the in-memory progression.
// The given progression must either be empty or have same length
// as the exercice.
// [idStudent] is only used to handle RandomMonoquestions
func (args EvaluateWorkIn) Evaluate(db ed.DB, idStudent tc.IdStudent) (EvaluateWorkOut, error) {
	data, err := newWorkLoader(db, args.ID, idStudent)
	if err != nil {
		return EvaluateWorkOut{}, utils.SQLError(err)
	}

	qus := data.Questions()

	// handle initial empty progressions
	if len(args.Progression.Questions) == 0 {
		args.Progression.Questions = make([]ta.QuestionHistory, len(qus))
	}

	// enforce invariant for existing progressions
	if L1, L2 := len(qus), len(args.Progression.Questions); L1 != L2 {
		return EvaluateWorkOut{}, fmt.Errorf("internal error: inconsistent length %d != %d", L1, L2)
	}

	// compat mode
	if err := args.fillFromMap(); err != nil {
		return EvaluateWorkOut{}, err
	}

	if args.AnswerIndex < 0 || args.AnswerIndex >= len(qus) {
		return EvaluateWorkOut{}, fmt.Errorf("internal error: invalid answer index %d", args.AnswerIndex)
	}

	// depending on the flow, check question index
	switch data.flow() {
	case ed.Parallel: // all questions are accessible
	case ed.Sequencial: // only the current question is accessible
		if exp := args.Progression.NextQuestion; args.AnswerIndex != exp {
			return EvaluateWorkOut{}, fmt.Errorf("internal error: expected answer for %d, got %d", exp, args.AnswerIndex)
		}
	}

	resp, err := EvaluateQuestion(qus[args.AnswerIndex].Enonce, args.Answer)
	if err != nil {
		return EvaluateWorkOut{}, err
	}

	updatedProgression := args.Progression.Copy()
	l := &updatedProgression.Questions[args.AnswerIndex]
	*l = append(*l, resp.IsCorrect())
	updatedProgression.inferNextQuestion() // update in case of success

	newVersion, err := data.Instantiate()
	if err != nil {
		return EvaluateWorkOut{}, err
	}

	out := EvaluateWorkOut{
		AnswerIndex:  args.AnswerIndex,
		Result:       resp,
		Progression:  updatedProgression,
		NewQuestions: newVersion.Questions,
	}

	// compat
	out.fillMap()

	return out, nil
}

// Ceintures variant

type InstantiatedBeltQuestion struct {
	Id       ce.IdBeltquestion
	Question client.Question
	Params   Params // for the evaluation
}

type BeltResult []client.QuestionAnswersOut

func EvaluateBelt(db ce.DB, questions []ce.IdBeltquestion, answers []AnswerP) (BeltResult, error) {
	if len(questions) != len(answers) {
		return nil, fmt.Errorf("internal error: length mistmatch")
	}

	questionsSource, err := ce.SelectBeltquestions(db, questions...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	out := make([]client.QuestionAnswersOut, len(answers))
	for index, idQuestion := range questions {
		answer := answers[index]
		qu := questionsSource[idQuestion]
		out[index], err = EvaluateQuestion(qu.Enonce, answer)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

func (br BeltResult) Stats() (hasPassed bool, stat ce.Stat) {
	hasPassed = true
	for _, res := range br {
		correct := res.IsCorrect()
		hasPassed = hasPassed && correct
		if correct {
			stat.Success += 1
		} else {
			stat.Failure += 1
		}
	}
	return
}
