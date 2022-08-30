package editor

import (
	"fmt"
	"sort"

	"github.com/benoitkugler/maths-online/maths/questions"
	tcAPI "github.com/benoitkugler/maths-online/prof/teacher"
	ed "github.com/benoitkugler/maths-online/sql/editor"
	"github.com/benoitkugler/maths-online/sql/teacher"
	"github.com/benoitkugler/maths-online/tasks"
	"github.com/benoitkugler/maths-online/utils"
	"github.com/labstack/echo/v4"
)

type uID = teacher.IdTeacher

type QuestionOrigin struct {
	Question ed.Question
	Origin   tcAPI.Origin
}

type ExerciceExt struct {
	Exercice        ed.Exercice
	Origin          tcAPI.Origin
	Questions       ed.ExerciceQuestions
	QuestionsSource map[ed.IdQuestion]ed.Question
}

type ExerciceHeader struct {
	Exercice  ed.Exercice
	Origin    tcAPI.Origin
	Questions ed.ExerciceQuestions
}

func (ct *Controller) ExercicesGetList(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	out, err := ct.getExercices(user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func exerciceOrigin(ex ed.Exercicegroup, userID, adminID uID) (tcAPI.Origin, bool) {
	vis := tcAPI.NewVisibility(ex.IdTeacher, userID, adminID, ex.Public)
	if vis.Restricted() {
		return tcAPI.Origin{}, false
	}
	return tcAPI.Origin{
		AllowPublish: userID == adminID,
		IsPublic:     ex.Public,
		Visibility:   vis,
	}, true
}

// buildExerciceHeaders aggregates the content of the tables Exercice and ExerciceQuestion,
// selecting only the exercices visible by `userID`
func buildExerciceHeaders(userID, adminID teacher.IdTeacher, groups ed.Exercicegroups, exes ed.Exercices, links ed.ExerciceQuestions) map[ed.IdExercice]ExerciceHeader {
	out := make(map[ed.IdExercice]ExerciceHeader, len(exes))
	questionDict := links.ByIdExercice()

	for _, ex := range exes {
		group := groups[ex.IdGroup]
		origin, ok := exerciceOrigin(group, userID, adminID)
		if !ok {
			continue
		}

		questions := questionDict[ex.Id]
		questions.EnsureOrder()
		out[ex.Id] = ExerciceHeader{
			Exercice:  ex,
			Origin:    origin,
			Questions: questions,
		}
	}

	return out
}

func (ct *Controller) getExercices(userID uID) ([]ExerciceHeader, error) {
	groups, err := ed.SelectAllExercicegroups(ct.db)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	exs, err := ed.SelectAllExercices(ct.db)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	links, err := ed.SelectExerciceQuestionsByIdExercices(ct.db, exs.IDs()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	tmp := buildExerciceHeaders(userID, ct.admin.Id, groups, exs, links)
	out := make([]ExerciceHeader, 0, len(tmp))
	for _, ex := range tmp {
		out = append(out, ex)
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Exercice.Id < out[j].Exercice.Id })
	sort.SliceStable(out, func(i, j int) bool { return out[i].Exercice.Subtitle < out[j].Exercice.Subtitle })

	return out, nil
}

// ExerciceGetContent loads the questions associated with the given exercice
func (ct *Controller) ExerciceGetContent(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	idExercice, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	out, err := ct.getExercice(ed.IdExercice(idExercice), user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) getExercice(exerciceID ed.IdExercice, userID uID) (ExerciceExt, error) {
	data, err := tasks.NewExerciceData(ct.db, exerciceID)
	if err != nil {
		return ExerciceExt{}, err
	}

	origin, ok := exerciceOrigin(data.Group, userID, ct.admin.Id)
	if !ok {
		return ExerciceExt{}, accessForbidden
	}

	out := ExerciceExt{
		Exercice:        data.Exercice,
		Origin:          origin,
		Questions:       data.Links,
		QuestionsSource: data.QuestionsSource,
	}

	return out, nil
}

// ExerciceCreate creates a new exercice group with one exercice
func (ct *Controller) ExerciceCreate(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	out, err := ct.createExercice(user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) createExercice(userID uID) (ExerciceHeader, error) {
	tx, err := ct.db.Begin()
	if err != nil {
		return ExerciceHeader{}, utils.SQLError(err)
	}

	group, err := ed.Exercicegroup{IdTeacher: userID}.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return ExerciceHeader{}, utils.SQLError(err)
	}

	ex, err := ed.Exercice{IdGroup: group.Id, Flow: ed.Parallel}.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return ExerciceHeader{}, utils.SQLError(err)
	}

	err = tx.Commit()
	if err != nil {
		return ExerciceHeader{}, utils.SQLError(err)
	}

	origin, _ := exerciceOrigin(group, userID, ct.admin.Id)
	out := ExerciceHeader{
		Exercice: ex,
		Origin:   origin,
	}

	return out, nil
}

// ExerciceDelete remove the given exercice, also cleaning
// up the exercice group if needed.
func (ct *Controller) ExerciceDelete(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	idExercice, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	err = ct.deleteExercice(ed.IdExercice(idExercice), user.Id)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) checkExerciceOwner(idExercice ed.IdExercice, userID uID) error {
	ex, err := ed.SelectExercice(ct.db, idExercice)
	if err != nil {
		return utils.SQLError(err)
	}

	group, err := ed.SelectExercicegroup(ct.db, ex.IdGroup)
	if err != nil {
		return utils.SQLError(err)
	}

	if group.IdTeacher != userID {
		return accessForbidden
	}

	return nil
}

func (ct *Controller) deleteExercice(idExercice ed.IdExercice, userID uID) error {
	if err := ct.checkExerciceOwner(idExercice, userID); err != nil {
		return err
	}

	links, err := ed.SelectExerciceQuestionsByIdExercices(ct.db, idExercice)
	if err != nil {
		return utils.SQLError(err)
	}
	qus, err := ed.SelectQuestions(ct.db, links.IdQuestions()...)
	if err != nil {
		return utils.SQLError(err)
	}

	// delete not standalone questions linked to the exercice
	var toDelete []ed.IdQuestion
	for _, question := range qus {
		if question.NeedExercice.Valid {
			toDelete = append(toDelete, question.Id)
		}
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}

	// remove the links
	_, err = ed.DeleteExerciceQuestionsByIdExercices(tx, idExercice)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	// remove the actual questions
	_, err = ed.DeleteQuestionsByIDs(tx, toDelete...)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	// finaly remove the exercice
	_, err = ed.DeleteExerciceById(tx, idExercice)
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
	SessionID  string
	IdExercice ed.IdExercice
}

// ExerciceCreateQuestion creates a question and appends it
// to the given exercice.
func (ct *Controller) ExerciceCreateQuestion(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	var args ExerciceCreateQuestionIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.createQuestionEx(args, user.Id)
	if err != nil {
		return err
	}

	data, err := tasks.NewExerciceData(ct.db, args.IdExercice)
	if err != nil {
		return err
	}

	err = ct.updateExercicePreview(data, args.SessionID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) createQuestionEx(args ExerciceCreateQuestionIn, userID uID) (ExerciceExt, error) {
	if err := ct.checkExerciceOwner(args.IdExercice, userID); err != nil {
		return ExerciceExt{}, err
	}

	existingLinks, err := ed.SelectExerciceQuestionsByIdExercices(ct.db, args.IdExercice)
	if err != nil {
		return ExerciceExt{}, utils.SQLError(err)
	}

	// creates a question linked to the given exercice
	question, err := ed.Question{NeedExercice: args.IdExercice.AsOptional()}.Insert(ct.db)
	if err != nil {
		return ExerciceExt{}, utils.SQLError(err)
	}

	// append it to the current questions
	existingLinks = append(existingLinks, ed.ExerciceQuestion{IdExercice: args.IdExercice, IdQuestion: question.Id, Bareme: 1})

	err = ed.UpdateExerciceQuestionList(ct.db, args.IdExercice, existingLinks)
	if err != nil {
		return ExerciceExt{}, err
	}

	return ct.getExercice(args.IdExercice, userID)
}

type ExerciceUpdateQuestionsIn struct {
	Questions  ed.ExerciceQuestions
	IdExercice ed.IdExercice
	SessionID  string
}

// ExerciceUpdateQuestions updates the question links and
// the preview
func (ct *Controller) ExerciceUpdateQuestions(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	var args ExerciceUpdateQuestionsIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.updateQuestionsEx(args, user.Id)
	if err != nil {
		return err
	}

	data, err := tasks.NewExerciceData(ct.db, args.IdExercice)
	if err != nil {
		return err
	}

	err = ct.updateExercicePreview(data, args.SessionID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) updateQuestionsEx(args ExerciceUpdateQuestionsIn, userID uID) (ExerciceExt, error) {
	if err := ct.checkExerciceOwner(args.IdExercice, userID); err != nil {
		return ExerciceExt{}, err
	}

	// garbage collect the question only used by this exercice
	links, err := ed.SelectExerciceQuestionsByIdExercices(ct.db, args.IdExercice)
	if err != nil {
		return ExerciceExt{}, utils.SQLError(err)
	}
	questions, err := ed.SelectQuestions(ct.db, args.Questions.IdQuestions()...)
	if err != nil {
		return ExerciceExt{}, utils.SQLError(err)
	}

	var (
		toDelete       []ed.IdQuestion
		newQuestionIDs = args.Questions.ByIdQuestion()
	)
	for _, link := range links {
		_, willBeUsed := newQuestionIDs[link.IdQuestion]
		if shouldDelete := questions[link.IdQuestion].NeedExercice.Valid && !willBeUsed; shouldDelete {
			toDelete = append(toDelete, link.IdQuestion)
		}
	}

	err = ed.UpdateExerciceQuestionList(ct.db, args.IdExercice, args.Questions)
	if err != nil {
		return ExerciceExt{}, err
	}

	_, err = ed.DeleteQuestionsByIDs(ct.db, toDelete...)
	if err != nil {
		return ExerciceExt{}, utils.SQLError(err)
	}

	return ct.getExercice(args.IdExercice, userID)
}

type ExerciceUpdateIn struct {
	Exercice  ed.Exercice
	SessionID string
}

// ExerciceUpdate update the exercice metadata and
// update the preview
func (ct *Controller) ExerciceUpdate(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	var args ExerciceUpdateIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.updateExercice(args.Exercice, user.Id)
	if err != nil {
		return err
	}

	data, err := tasks.NewExerciceData(ct.db, args.Exercice.Id)
	if err != nil {
		return err
	}

	err = ct.updateExercicePreview(data, args.SessionID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) updateExercice(in ed.Exercice, userID uID) (ed.Exercice, error) {
	if err := ct.checkExerciceOwner(in.Id, userID); err != nil {
		return ed.Exercice{}, err
	}

	ex, err := ed.SelectExercice(ct.db, in.Id)
	if err != nil {
		return ed.Exercice{}, err
	}

	// only update meta data
	ex.Description = in.Description
	ex.Subtitle = in.Subtitle
	ex.Flow = in.Flow
	ex, err = ex.Update(ct.db)
	if err != nil {
		return ed.Exercice{}, utils.SQLError(err)
	}

	return ex, nil
}

type CheckExerciceParametersIn struct {
	IdExercice         ed.IdExercice
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
	data, err := tasks.NewExerciceData(ct.db, params.IdExercice)
	if err != nil {
		return CheckExerciceParametersOut{}, err
	}
	qus, _ := data.QuestionsList()

	if L1, L2 := len(params.QuestionParameters), len(qus); L1 != L2 {
		return CheckExerciceParametersOut{}, fmt.Errorf("internal error: mismatched question length (%d != %d)", L1, L2)
	}

	for index, question := range qus {
		toCheck := params.QuestionParameters[index]
		if question.NeedExercice.Valid { // add the shared parameters // TODO: probably cleanup the check
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
	IdExercice ed.IdExercice
	Parameters questions.Parameters // shared parameters
	Questions  ed.Questions         // questions content
}

type SaveExerciceAndPreviewOut struct {
	Error         questions.ErrQuestionInvalid
	QuestionIndex int
	IsValid       bool
}

func (ct *Controller) saveExerciceAndPreview(params SaveExerciceAndPreviewIn, userID uID) (SaveExerciceAndPreviewOut, error) {
	data, err := tasks.NewExerciceData(ct.db, params.IdExercice)
	if err != nil {
		return SaveExerciceAndPreviewOut{}, err
	}
	ex := &data.Exercice

	if !data.Group.IsVisibleBy(userID) {
		return SaveExerciceAndPreviewOut{}, accessForbidden
	}

	// validate all the questions, using shared parameters if needed
	qus, _ := data.QuestionsList()
	for index, question := range qus {
		toCheck := params.Questions[question.Id].Page
		if question.NeedExercice.Valid { // add the shared parameters
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

	// always apply change in memory, so that preview is correctly updated
	ex.Parameters = params.Parameters // save the shared parameters
	for _, incomming := range params.Questions {
		qu := data.QuestionsSource[incomming.Id]
		// update the content
		qu.Page = incomming.Page
		qu.Description = incomming.Description
		data.QuestionsSource[incomming.Id] = qu
	}

	// if the exercice is owned : save it, else only preview
	if data.Group.IdTeacher == userID {
		tx, err := ct.db.Begin()
		if err != nil {
			return SaveExerciceAndPreviewOut{}, utils.SQLError(err)
		}

		_, err = ex.Update(tx)
		if err != nil {
			_ = tx.Rollback()
			return SaveExerciceAndPreviewOut{}, utils.SQLError(err)
		}

		// update the linked questions
		// TODO: only do it for sequencial exercices
		for _, qu := range data.QuestionsSource {
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

	err = ct.updateExercicePreview(data, params.SessionID)
	if err != nil {
		return SaveExerciceAndPreviewOut{}, err
	}

	return SaveExerciceAndPreviewOut{IsValid: true}, nil
}

// updateExercicePreview instantiates the exercice and send preview data
func (ct *Controller) updateExercicePreview(content tasks.ExerciceData, sessionID string) error {
	instance, err := content.Instantiate()
	if err != nil {
		return err
	}

	ct.lock.Lock()
	defer ct.lock.Unlock()

	loopback, ok := ct.sessions[sessionID]
	if !ok {
		return fmt.Errorf("invalid session ID %s", sessionID)
	}

	loopback.setExercice(instance)
	return nil
}
