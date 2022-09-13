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

type ExercicegroupExt struct {
	Group    ed.Exercicegroup
	Origin   tcAPI.Origin
	Tags     []string
	Variants []ExerciceHeader
}

type ExerciceHeader struct {
	Id         ed.IdExercice
	Subtitle   string
	Difficulty ed.DifficultyTag
}

func newExerciceHeader(exercice ed.Exercice) ExerciceHeader {
	return ExerciceHeader{
		Id:         exercice.Id,
		Subtitle:   exercice.Subtitle,
		Difficulty: ed.DiffEmpty, // for now, we don't support difficulty on exercices
	}
}

type ExerciceQuestionExt struct {
	Question ed.Question
	Bareme   int
}

// type ExerciceHeader struct {
// 	Exercice  ed.Exercice
// 	Questions ed.ExerciceQuestions
// }

type ExerciceExt struct {
	Exercice  ed.Exercice
	Questions []ExerciceQuestionExt
}

type SearchExercicesIn = Query

func (ct *Controller) EditorSearchExercices(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	var args SearchExercicesIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	out, err := ct.searchExercices(args, user.Id)
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

type ListExercicesOut struct {
	Groups      []ExercicegroupExt // limited by `pagination`
	NbGroups    int                // total number of groups (passing the given filter)
	NbExercices int                // total number of exercices contained in the groups
}

func (ct *Controller) searchExercices(query Query, userID uID) (out ListExercicesOut, err error) {
	const pagination = 30 // number of groups

	groups, err := ed.SelectAllExercicegroups(ct.db)
	if err != nil {
		return out, utils.SQLError(err)
	}
	groups.RestrictVisible(userID)

	query.normalize()

	// restrict the groups to matching title
	matcher, err := newQuery(query.TitleQuery)
	if err != nil {
		return out, err
	}
	for _, group := range groups {
		if !matcher.match(int64(group.Id), group.Title) {
			delete(groups, group.Id)
		}
	}

	// load the tags ...
	tags, err := ed.SelectExercicegroupTagsByIdExercicegroups(ct.db, groups.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}
	tagsMap := tags.ByIdExercicegroup()

	// ... and the exercices
	tmp, err := ed.SelectExercicesByIdGroups(ct.db, groups.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}
	exercicesByGroup := tmp.ByGroup()

	// .. and build the groups, restricting to the ones matching the given tags
	for _, group := range groups {
		crible := tagsMap[group.Id].Crible()
		if !crible.HasAll(query.Tags) {
			continue
		}
		exercices := exercicesByGroup[group.Id]
		if len(exercices) == 0 { // ignore empty groupExts
			continue
		}

		origin, _ := exerciceOrigin(group, userID, ct.admin.Id)
		groupExt := ExercicegroupExt{
			Group:  group,
			Origin: origin,
			Tags:   tagsMap[group.Id].List(),
		}

		for _, exercice := range exercices {
			groupExt.Variants = append(groupExt.Variants, newExerciceHeader(exercice))
		}

		// sort to make sure the display is consistent between two queries
		sort.Slice(groupExt.Variants, func(i, j int) bool { return groupExt.Variants[i].Id < groupExt.Variants[j].Id })
		sort.SliceStable(groupExt.Variants, func(i, j int) bool { return groupExt.Variants[i].Difficulty < groupExt.Variants[j].Difficulty })

		out.NbExercices += len(groupExt.Variants)

		out.Groups = append(out.Groups, groupExt)
	}

	// sort before pagination
	sort.Slice(out.Groups, func(i, j int) bool { return out.Groups[i].Group.Title < out.Groups[j].Group.Title })

	out.NbGroups = len(out.Groups)
	if len(out.Groups) > pagination {
		out.Groups = out.Groups[:pagination]
	}

	return out, nil
}

func (ct *Controller) EditorUpdateExercicegroup(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	var args ed.Exercicegroup
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	group, err := ed.SelectExercicegroup(ct.db, args.Id)
	if err != nil {
		return utils.SQLError(err)
	}

	if group.IdTeacher != user.Id {
		return accessForbidden
	}

	group.Title = args.Title
	_, err = group.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	return c.NoContent(200)
}

// EditorDuplicateExercice duplicate one variant inside a group.
func (ct *Controller) EditorDuplicateExercice(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	out, err := ct.duplicateExercice(ed.IdExercice(id), user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) getExerciceGroup(qu ed.Exercice) (ed.Exercicegroup, error) {
	group, err := ed.SelectExercicegroup(ct.db, qu.IdGroup)
	if err != nil {
		return ed.Exercicegroup{}, utils.SQLError(err)
	}
	return group, nil
}

// duplicateExercice duplicate the given exercice, returning
// the newly created one
func (ct *Controller) duplicateExercice(idExercice ed.IdExercice, userID uID) (ExerciceHeader, error) {
	ex, err := ed.SelectExercice(ct.db, idExercice)
	if err != nil {
		return ExerciceHeader{}, utils.SQLError(err)
	}

	group, err := ct.getExerciceGroup(ex)
	if err != nil {
		return ExerciceHeader{}, err
	}

	if !group.IsVisibleBy(userID) {
		return ExerciceHeader{}, accessForbidden
	}

	// shallow copy is enough
	newExercice := ex
	newExercice, err = newExercice.Insert(ct.db)
	if err != nil {
		return ExerciceHeader{}, utils.SQLError(err)
	}

	// also copy the questions
	links, err := ed.SelectExerciceQuestionsByIdExercices(ct.db, ex.Id)
	if err != nil {
		return ExerciceHeader{}, utils.SQLError(err)
	}
	questions, err := ed.SelectQuestions(ct.db, links.IdQuestions()...)
	if err != nil {
		return ExerciceHeader{}, utils.SQLError(err)
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return ExerciceHeader{}, utils.SQLError(err)
	}

	for index, link := range links {
		question := questions[link.IdQuestion]
		question.NeedExercice = newExercice.Id.AsOptional()
		question, err = question.Insert(tx)
		if err != nil {
			_ = tx.Rollback()
			return ExerciceHeader{}, utils.SQLError(err)
		}
		links[index].IdExercice = newExercice.Id
		links[index].IdQuestion = question.Id
	}

	err = ed.UpdateExerciceQuestionList(tx, newExercice.Id, links)
	if err != nil {
		_ = tx.Rollback()
		return ExerciceHeader{}, err
	}

	err = tx.Commit()
	if err != nil {
		return ExerciceHeader{}, utils.SQLError(err)
	}

	return newExerciceHeader(newExercice), nil
}

// EditorGetExerciceContent loads the questions associated with the given exercice
func (ct *Controller) EditorGetExerciceContent(c echo.Context) error {
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

	questions, baremes := data.QuestionsList()
	l := make([]ExerciceQuestionExt, len(questions))
	for i := range questions {
		l[i] = ExerciceQuestionExt{Question: questions[i], Bareme: baremes[i]}
	}
	out := ExerciceExt{
		Exercice:  data.Exercice,
		Questions: l,
	}

	return out, nil
}

// EditorCreateExercice creates a new exercice group with one exercice
func (ct *Controller) EditorCreateExercice(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	out, err := ct.createExercice(user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) createExercice(userID uID) (ExercicegroupExt, error) {
	tx, err := ct.db.Begin()
	if err != nil {
		return ExercicegroupExt{}, utils.SQLError(err)
	}

	group, err := ed.Exercicegroup{IdTeacher: userID}.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return ExercicegroupExt{}, utils.SQLError(err)
	}

	ex, err := ed.Exercice{IdGroup: group.Id}.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return ExercicegroupExt{}, utils.SQLError(err)
	}

	err = tx.Commit()
	if err != nil {
		return ExercicegroupExt{}, utils.SQLError(err)
	}

	origin, _ := exerciceOrigin(group, userID, ct.admin.Id)
	out := ExercicegroupExt{
		Group:    group,
		Origin:   origin,
		Tags:     nil,
		Variants: []ExerciceHeader{newExerciceHeader(ex)},
	}
	return out, nil
}

type UpdateExercicegroupTagsIn struct {
	Id   ed.IdExercicegroup
	Tags []string
}

func (ct *Controller) EditorUpdateExerciceTags(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	var args UpdateExercicegroupTagsIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	err := ct.updateExerciceTags(args, user.Id)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

func (ct *Controller) updateExerciceTags(params UpdateExercicegroupTagsIn, userID uID) error {
	group, err := ed.SelectExercicegroup(ct.db, params.Id)
	if err != nil {
		return utils.SQLError(err)
	}
	if group.IdTeacher != userID {
		return accessForbidden
	}

	var tags ed.ExercicegroupTags
	for _, tag := range params.Tags {
		// enforce proper tags
		tag = ed.NormalizeTag(tag)
		if tag == "" {
			continue
		}

		tags = append(tags, ed.ExercicegroupTag{IdExercicegroup: params.Id, Tag: tag})
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return err
	}

	err = ed.UpdateExercicegroupTags(tx, tags, params.Id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

// EditorDeleteExercice remove the given exercice, also cleaning
// up the exercice group if needed.
func (ct *Controller) EditorDeleteExercice(c echo.Context) error {
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

// EditorExerciceCreateQuestion creates a question and appends it
// to the given exercice.
func (ct *Controller) EditorExerciceCreateQuestion(c echo.Context) error {
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
	tx, err := ct.db.Begin()
	if err != nil {
		return ExerciceExt{}, utils.SQLError(err)
	}

	question, err := ed.Question{NeedExercice: args.IdExercice.AsOptional()}.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return ExerciceExt{}, utils.SQLError(err)
	}

	// append it to the current questions
	existingLinks = append(existingLinks, ed.ExerciceQuestion{IdExercice: args.IdExercice, IdQuestion: question.Id, Bareme: 1})

	err = ed.UpdateExerciceQuestionList(tx, args.IdExercice, existingLinks)
	if err != nil {
		_ = tx.Rollback()
		return ExerciceExt{}, err
	}

	err = tx.Commit()
	if err != nil {
		return ExerciceExt{}, utils.SQLError(err)
	}

	return ct.getExercice(args.IdExercice, userID)
}

type ExerciceUpdateQuestionsIn struct {
	Questions  ed.ExerciceQuestions
	IdExercice ed.IdExercice
	SessionID  string
}

// EditorExerciceUpdateQuestions updates the question links and
// the preview
func (ct *Controller) EditorExerciceUpdateQuestions(c echo.Context) error {
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

	tx, err := ct.db.Begin()
	if err != nil {
		return ExerciceExt{}, utils.SQLError(err)
	}

	err = ed.UpdateExerciceQuestionList(tx, args.IdExercice, args.Questions)
	if err != nil {
		_ = tx.Rollback()
		return ExerciceExt{}, err
	}

	_, err = ed.DeleteQuestionsByIDs(tx, toDelete...)
	if err != nil {
		_ = tx.Rollback()
		return ExerciceExt{}, utils.SQLError(err)
	}

	err = tx.Commit()
	if err != nil {
		return ExerciceExt{}, utils.SQLError(err)
	}

	return ct.getExercice(args.IdExercice, userID)
}

type ExerciceUpdateIn struct {
	Exercice  ed.Exercice
	SessionID string
}

// EditorUpdateExercice update the exercice metadata and
// update the preview
func (ct *Controller) EditorUpdateExercice(c echo.Context) error {
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

	for index := range qus {
		toCheck := params.QuestionParameters[index]
		toCheck = toCheck.Append(params.SharedParameters)

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
	Questions  []ed.Question        // questions content
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
		qu := data.QuestionsMap[incomming.Id]
		// update the content
		qu.Page = incomming.Page
		qu.Description = incomming.Description
		data.QuestionsMap[incomming.Id] = qu
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
		for _, qu := range data.QuestionsMap {
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
