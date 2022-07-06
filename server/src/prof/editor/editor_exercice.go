package editor

import (
	"fmt"
	"sort"

	"github.com/benoitkugler/maths-online/maths/questions"
	"github.com/benoitkugler/maths-online/prof/teacher"
	"github.com/benoitkugler/maths-online/utils"
	"github.com/labstack/echo/v4"
)

func (l ExerciceQuestions) ensureIndex() {
	sort.Slice(l, func(i, j int) bool { return l[i].Index < l[j].Index })
}

func (l ProgressionQuestions) ensureIndex() {
	sort.Slice(l, func(i, j int) bool { return l[i].Index < l[j].Index })
}

type ExerciceExt struct {
	Exercice  Exercice
	Origin    teacher.Origin
	Questions []ExerciceQuestionExt
}

type ExerciceQuestionExt struct {
	Link     ExerciceQuestion
	Question Question
	Origin   teacher.Origin
}

func fillQuestions(l ExerciceQuestions, dict Questions, userID, adminID int64) []ExerciceQuestionExt {
	l.ensureIndex()
	out := make([]ExerciceQuestionExt, len(l))
	for i, qu := range l {
		question := dict[qu.IdQuestion]
		out[i].Link = qu
		out[i].Question = question
		origin, _ := question.origin(userID, adminID)
		out[i].Origin = origin
	}
	return out
}

type ExerciceHeader struct {
	Exercice  Exercice
	Origin    teacher.Origin
	Questions ExerciceQuestions
}

func (ct *Controller) ExercicesGetList(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	out, err := ct.getExercices(user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ex Exercice) origin(userID, adminID int64) (teacher.Origin, bool) {
	vis, ok := teacher.NewVisibility(ex.IdTeacher, userID, adminID, ex.Public)
	if !ok {
		return teacher.Origin{}, false
	}
	return teacher.Origin{
		AllowPublish: userID == adminID,
		IsPublic:     ex.Public,
		Visibility:   vis,
	}, true
}

func (ct *Controller) getExercices(userID int64) ([]ExerciceHeader, error) {
	exs, err := SelectAllExercices(ct.db)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	var (
		out []ExerciceHeader
		tmp IDs
	)
	for _, ex := range exs {
		origin, ok := ex.origin(userID, ct.admin.Id)
		if !ok {
			continue
		}
		out = append(out, ExerciceHeader{
			Exercice: ex,
			Origin:   origin,
		})
		tmp = append(tmp, ex.Id)
	}

	links, err := SelectExerciceQuestionsByIdExercices(ct.db, tmp...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	dict := links.ByIdExercice()

	for i, ex := range out {
		s := dict[ex.Exercice.Id]
		s.ensureIndex()
		out[i].Questions = s
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Exercice.Id < out[j].Exercice.Id })
	sort.SliceStable(out, func(i, j int) bool { return out[i].Exercice.Title < out[j].Exercice.Title })

	return out, nil
}

// ExerciceGetContent loads the questions associated with the given exercice
func (ct *Controller) ExerciceGetContent(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	idExercice, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	data, err := ct.loadExercice(idExercice)
	if err != nil {
		return err
	}
	ex := data.exercice

	origin, ok := ex.origin(user.Id, ct.admin.Id)
	if !ok {
		return accessForbidden
	}

	out := ExerciceExt{
		Exercice:  ex,
		Origin:    origin,
		Questions: fillQuestions(data.links, data.dict, user.Id, ct.admin.Id),
	}

	return c.JSON(200, out)
}

func (ct *Controller) ExerciceCreate(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	out, err := ct.createExercice(user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) createExercice(userID int64) (ExerciceHeader, error) {
	ex, err := Exercice{IdTeacher: userID, Flow: Parallel}.Insert(ct.db)
	if err != nil {
		return ExerciceHeader{}, utils.SQLError(err)
	}

	out := ExerciceHeader{
		Exercice: ex,
		Origin: teacher.Origin{
			AllowPublish: userID == ct.admin.Id,
			IsPublic:     false,
			Visibility:   teacher.Personnal,
		},
	}

	return out, nil
}

func (ct *Controller) ExerciceDelete(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	idExercice, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	deleteQuestions := utils.QueryParamBool(c, "delete_questions")

	err = ct.deleteExercice(idExercice, deleteQuestions, user.Id)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) checkExerciceOwner(idExercice, userID int64) error {
	ex, err := SelectExercice(ct.db, idExercice)
	if err != nil {
		return utils.SQLError(err)
	}

	if ex.IdTeacher != userID {
		return accessForbidden
	}

	return nil
}

func (ct *Controller) deleteExercice(idExercice int64, deleteQuestions bool, userID int64) error {
	if err := ct.checkExerciceOwner(idExercice, userID); err != nil {
		return err
	}

	links, err := SelectExerciceQuestionsByIdExercices(ct.db, idExercice)
	if err != nil {
		return utils.SQLError(err)
	}
	qus, err := SelectQuestions(ct.db, links.IdQuestions()...)
	if err != nil {
		return utils.SQLError(err)
	}

	// we always delete not standalone questions linked to the exercice
	var toDelete IDs
	for _, question := range qus {
		if question.NeedExercice || deleteQuestions {
			toDelete = append(toDelete, question.Id)
		}
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}

	// remove the links
	_, err = DeleteExerciceQuestionsByIdExercices(tx, idExercice)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	// remove the actual questions
	_, err = DeleteQuestionsByIDs(tx, toDelete...)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	// finaly remove the exercice
	_, err = DeleteExerciceById(tx, idExercice)
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

type ExerciceCreateQuestionIn struct {
	IdExercice int64
}

// ExerciceCreateQuestion creates a question and appends it
// to the given exercice.
func (ct *Controller) ExerciceCreateQuestion(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	var args ExerciceCreateQuestionIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.createQuestionEx(args, user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) createQuestionEx(args ExerciceCreateQuestionIn, userID int64) ([]ExerciceQuestionExt, error) {
	if err := ct.checkExerciceOwner(args.IdExercice, userID); err != nil {
		return nil, err
	}

	// creates a question
	question, err := Question{IdTeacher: userID, Public: false, NeedExercice: true}.Insert(ct.db)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	// append if to the current questions
	existing, err := SelectExerciceQuestionsByIdExercices(ct.db, args.IdExercice)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	existing = append(existing, ExerciceQuestion{IdExercice: args.IdExercice, IdQuestion: question.Id, Bareme: 1})

	out, err := updateExerciceQuestionList(ct.db, args.IdExercice, existing, userID, ct.admin.Id)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	return out, nil
}

type ExerciceUpdateQuestionsIn struct {
	Questions  ExerciceQuestions
	IdExercice int64
}

func (ct *Controller) ExerciceUpdateQuestions(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	var args ExerciceUpdateQuestionsIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.updateQuestionsEx(args, user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) updateQuestionsEx(args ExerciceUpdateQuestionsIn, userID int64) ([]ExerciceQuestionExt, error) {
	if err := ct.checkExerciceOwner(args.IdExercice, userID); err != nil {
		return nil, err
	}

	// garbage collect the question only used by this exercice
	links, err := SelectExerciceQuestionsByIdExercices(ct.db, args.IdExercice)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	questions, err := SelectQuestions(ct.db, args.Questions.IdQuestions()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	var (
		toDelete       IDs
		newQuestionIDs = args.Questions.ByIdQuestion()
	)
	for _, link := range links {
		_, willBeUsed := newQuestionIDs[link.IdQuestion]
		if shouldDelete := questions[link.IdQuestion].NeedExercice && !willBeUsed; shouldDelete {
			toDelete = append(toDelete, link.IdQuestion)
		}
	}

	out, err := updateExerciceQuestionList(ct.db, args.IdExercice, args.Questions, userID, ct.admin.Id)
	if err != nil {
		return nil, err
	}

	_, err = DeleteQuestionsByIDs(ct.db, toDelete...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	return out, nil
}

func (ct *Controller) ExerciceUpdate(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	var args Exercice
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.updateExercice(args, user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) updateExercice(in Exercice, userID int64) (Exercice, error) {
	if err := ct.checkExerciceOwner(in.Id, userID); err != nil {
		return Exercice{}, err
	}

	ex, err := SelectExercice(ct.db, in.Id)
	if err != nil {
		return Exercice{}, err
	}

	// only update meta data
	ex.Description = in.Description
	ex.Title = in.Title
	ex.Flow = in.Flow
	return ex.Update(ct.db)
}

type CheckExerciceParametersIn struct {
	IdExercice         int64
	SharedParameters   questions.Parameters
	QuestionParameters []questions.Parameters
}

type CheckExerciceParametersOut struct {
	ErrDefinition questions.ErrParameters
	QuestionIndex int // ignored if ErrDefinition is empty
}

// checks that the merging of SharedParameters and QuestionParameters is valid
func (ct *Controller) checkExerciceParameters(params CheckExerciceParametersIn) (CheckExerciceParametersOut, error) {
	// fetch the mode of each question
	data, err := ct.loadExercice(params.IdExercice)
	if err != nil {
		return CheckExerciceParametersOut{}, err
	}
	qus := data.questions()

	if L1, L2 := len(params.QuestionParameters), len(qus); L1 != L2 {
		return CheckExerciceParametersOut{}, fmt.Errorf("internal error: mismatched question length (%d != %d)", L1, L2)
	}

	for index, question := range qus {
		toCheck := params.QuestionParameters[index]
		if question.NeedExercice { // add the shared parameters
			toCheck = toCheck.Append(params.SharedParameters)
		}

		err := toCheck.Validate()
		if err != nil {
			return CheckExerciceParametersOut{
				ErrDefinition: err.(questions.ErrParameters),
				QuestionIndex: index,
			}, nil
		}
	}

	return CheckExerciceParametersOut{}, nil
}

type SaveExerciceAndPreviewIn struct {
	SessionID  string
	IdExercice int64
	Parameters questions.Parameters // shared parameters
	Questions  []Question           // questions content
}

type SaveExerciceAndPreviewOut struct {
	Error         questions.ErrQuestionInvalid
	QuestionIndex int
	IsValid       bool
}

func (ct *Controller) saveExerciceAndPreview(params SaveExerciceAndPreviewIn, userID int64) (SaveExerciceAndPreviewOut, error) {
	data, err := ct.loadExercice(params.IdExercice)
	if err != nil {
		return SaveExerciceAndPreviewOut{}, err
	}
	ex, qus := data.exercice, data.questions()

	if !ex.IsVisibleBy(userID) {
		return SaveExerciceAndPreviewOut{}, accessForbidden
	}

	if L1, L2 := len(params.Questions), len(qus); L1 != L2 {
		return SaveExerciceAndPreviewOut{}, fmt.Errorf("internal error: mismatched question length (%d != %d)", L1, L2)
	}
	for index := range qus {
		if id1, id2 := params.Questions[index].Id, qus[index].Id; id1 != id2 {
			return SaveExerciceAndPreviewOut{}, fmt.Errorf("internal error: inconsistent question ID (%d != %d)", id1, id2)
		}
	}

	// validate all the questions, using shared parameters if needed
	for index, question := range qus {
		toCheck := params.Questions[index].Page
		if question.NeedExercice { // add the shared parameters
			toCheck.Parameters = toCheck.Parameters.Append(params.Parameters)
		}

		err = toCheck.Validate()
		if err != nil {
			return SaveExerciceAndPreviewOut{
				Error:         err.(questions.ErrQuestionInvalid),
				QuestionIndex: index,
			}, nil
		}
	}

	// if the exerice is owned : save it, else only preview
	if ex.IdTeacher == userID {
		tx, err := ct.db.Begin()
		if err != nil {
			return SaveExerciceAndPreviewOut{}, utils.SQLError(err)
		}

		// save the shared parameters
		ex.Parameters = params.Parameters
		_, err = ex.Update(tx)
		if err != nil {
			_ = tx.Rollback()
			return SaveExerciceAndPreviewOut{}, utils.SQLError(err)
		}

		// if owner, save the linked questions
		for index, qu := range qus {
			if qu.IdTeacher != userID {
				continue
			}
			// update the content
			qu.Page = params.Questions[index].Page
			qu.Description = params.Questions[index].Description
			_, err = qu.Update(tx)
			if err != nil {
				_ = tx.Rollback()
				return SaveExerciceAndPreviewOut{}, utils.SQLError(err)
			}
		}

		if err := tx.Commit(); err != nil {
			return SaveExerciceAndPreviewOut{}, utils.SQLError(err)
		}
	}

	// eventually, instantiate the exercice and send preview data
	instance, err := ct.instantiateExercice(ex.Id)
	if err != nil {
		return SaveExerciceAndPreviewOut{}, err
	}

	ct.lock.Lock()
	defer ct.lock.Unlock()

	loopback, ok := ct.sessions[params.SessionID]
	if !ok {
		return SaveExerciceAndPreviewOut{}, fmt.Errorf("invalid session ID %s", params.SessionID)
	}

	loopback.setExercice(instance)

	return SaveExerciceAndPreviewOut{IsValid: true}, nil
}
