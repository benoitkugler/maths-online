package editor

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"

	"github.com/benoitkugler/maths-online/server/src/maths/questions"
	tcAPI "github.com/benoitkugler/maths-online/server/src/prof/teacher"
	ed "github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/sql/reviews"
	ta "github.com/benoitkugler/maths-online/server/src/sql/tasks"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	taAPI "github.com/benoitkugler/maths-online/server/src/tasks"
	"github.com/benoitkugler/maths-online/server/src/utils"
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
	Tags     ed.Tags
	Variants []ExerciceHeader
}

func NewExercicegroupExt(group ed.Exercicegroup, variants []ed.Exercice, tags ed.Tags, inReview tcAPI.OptionalIdReview, userID, adminID uID,
) ExercicegroupExt {
	origin, _ := exerciceOrigin(group, inReview, userID, adminID)
	groupExt := ExercicegroupExt{
		Group:  group,
		Origin: origin,
		Tags:   tags,
	}

	for _, exercice := range variants {
		groupExt.Variants = append(groupExt.Variants, newExerciceHeader(exercice))
	}

	// sort to make sure the display is consistent between two queries
	sort.Slice(groupExt.Variants, func(i, j int) bool { return groupExt.Variants[i].Id < groupExt.Variants[j].Id })
	sort.SliceStable(groupExt.Variants, func(i, j int) bool { return groupExt.Variants[i].Difficulty < groupExt.Variants[j].Difficulty })

	return groupExt
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
		Difficulty: exercice.Difficulty,
	}
}

func (ct *Controller) EditorGetExercicesIndex(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	out, err := ct.loadExercicesIndex(user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) loadExercicesIndex(userID uID) (Index, error) {
	groups, err := ed.SelectAllExercicegroups(ct.db)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	groups.RestrictVisible(userID)

	// load the tags ...
	tags, err := ed.SelectExercicegroupTagsByIdExercicegroups(ct.db, groups.IDs()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	return buildIndex(exercicesToIndex(groups, tags)), nil
}

type ExerciceQuestionExt struct {
	Question ed.Question
	Bareme   int
}

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

func exerciceOrigin(ex ed.Exercicegroup, inReview tcAPI.OptionalIdReview, userID, adminID uID) (tcAPI.Origin, bool) {
	vis := tcAPI.NewVisibility(ex.IdTeacher, userID, adminID, ex.Public)
	if vis.Restricted() {
		return tcAPI.Origin{}, false
	}
	return tcAPI.Origin{
		AllowPublish: userID == adminID,
		IsPublic:     ex.Public,
		Visibility:   vis,
		IsInReview:   inReview,
	}, true
}

type ListExercicesOut struct {
	Groups      []ExercicegroupExt // limited by `pagination`
	NbGroups    int                // total number of groups (passing the given filter)
	NbExercices int                // total number of exercices contained in the groups
}

func (ct *Controller) searchExercices(query Query, userID uID) (out ListExercicesOut, err error) {
	const pagination = 10 // number of groups

	groups, err := ed.SelectAllExercicegroups(ct.db)
	if err != nil {
		return out, utils.SQLError(err)
	}
	groups.RestrictVisible(userID)

	query.normalize()

	// restrict the groups to matching title and origin
	matcher, err := newQuery(query.TitleQuery)
	if err != nil {
		return out, err
	}
	for _, group := range groups {
		vis := tcAPI.NewVisibility(group.IdTeacher, userID, ct.admin.Id, group.Public)

		keep := query.matchOrigin(vis) && matcher.match(int64(group.Id), group.Title)
		if !keep {
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

	revs, err := reviews.SelectAllReviewExercices(ct.db)
	if err != nil {
		return out, utils.SQLError(err)
	}
	revsMap := revs.ByIdExercice()

	// .. and build the groups, restricting to the ones matching the given tags
	for _, group := range groups {
		tagIndex := tagsMap[group.Id].Tags().BySection()
		if !(query.matchLevel(tagIndex.Level) && query.matchChapter(tagIndex.Chapter)) {
			continue
		}

		variants := exercicesByGroup[group.Id]
		if len(variants) == 0 { // ignore empty groupExts
			continue
		}

		var inReview tcAPI.OptionalIdReview
		link, isInReview := revsMap[group.Id]
		if isInReview {
			inReview = tcAPI.OptionalIdReview{InReview: true, Id: link.IdReview}
		}

		groupExt := NewExercicegroupExt(group, variants, tagsMap[group.Id].Tags(), inReview, userID, ct.admin.Id)

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
		return errAccessForbidden
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

// duplicateExercice duplicate the given exercice (variant), returning
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
		return ExerciceHeader{}, errAccessForbidden
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return ExerciceHeader{}, utils.SQLError(err)
	}

	newExercice, err := duplicateExerciceTo(tx, idExercice, group)
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

// duplicateExerciceTo duplicate the given exercice (variant) and assign it
// to the given group.
// It does NOT rollback or commit
func duplicateExerciceTo(tx *sql.Tx, idExercice ed.IdExercice, group ed.Exercicegroup) (ed.Exercice, error) {
	ex, err := ed.SelectExercice(tx, idExercice)
	if err != nil {
		return ed.Exercice{}, utils.SQLError(err)
	}

	// shallow copy is enough
	newExercice := ex
	newExercice.Subtitle += " (2)" // to distinguish
	newExercice.IdGroup = group.Id // re-direct to the given group
	newExercice, err = newExercice.Insert(tx)
	if err != nil {
		return ed.Exercice{}, utils.SQLError(err)
	}

	// also copy the questions
	links, err := ed.SelectExerciceQuestionsByIdExercices(tx, ex.Id)
	if err != nil {
		return ed.Exercice{}, utils.SQLError(err)
	}
	questions, err := ed.SelectQuestions(tx, links.IdQuestions()...)
	if err != nil {
		return ed.Exercice{}, utils.SQLError(err)
	}

	for index, link := range links {
		question := questions[link.IdQuestion]
		question.NeedExercice = newExercice.Id.AsOptional() // redirect the questions to the new exerice variant
		question, err = question.Insert(tx)
		if err != nil {
			return ed.Exercice{}, utils.SQLError(err)
		}
		links[index].IdExercice = newExercice.Id
		links[index].IdQuestion = question.Id
	}

	err = ed.UpdateExerciceQuestionList(tx, newExercice.Id, links)
	if err != nil {
		return ed.Exercice{}, err
	}

	return newExercice, nil
}

// EditorDuplicateExercicegroup duplicates the whole group, deep copying
// every exercice (variant), and assigns it to the current user.
func (ct *Controller) EditorDuplicateExercicegroup(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	err = ct.duplicateExercicegroup(ed.IdExercicegroup(id), user.Id)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

// duplicateExercicegroup creates a new group with the same title, (copied) exercices and tags
func (ct *Controller) duplicateExercicegroup(idGroup ed.IdExercicegroup, userID uID) error {
	group, err := ed.SelectExercicegroup(ct.db, idGroup)
	if err != nil {
		return utils.SQLError(err)
	}

	if !group.IsVisibleBy(userID) {
		return errAccessForbidden
	}

	tags, err := ed.SelectExercicegroupTagsByIdExercicegroups(ct.db, group.Id)
	if err != nil {
		return utils.SQLError(err)
	}

	variants, err := ed.SelectExercicesByIdGroups(ct.db, group.Id)
	if err != nil {
		return utils.SQLError(err)
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}

	// start by inserting a new group ...
	newGroup := group
	newGroup.IdTeacher = userID // redirect to the current user
	newGroup.Public = false
	newGroup, err = newGroup.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}
	// .. then add its tags ..
	err = ed.UpdateExercicegroupTags(tx, tags.Tags(), newGroup.Id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	// finaly, copy the exercice variantes...
	for _, variant := range variants {
		_, err = duplicateExerciceTo(tx, variant.Id, newGroup)
		if err != nil {
			_ = tx.Rollback()
			return utils.SQLError(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}

func (ct *Controller) getExerciceGroup(qu ed.Exercice) (ed.Exercicegroup, error) {
	group, err := ed.SelectExercicegroup(ct.db, qu.IdGroup)
	if err != nil {
		return ed.Exercicegroup{}, utils.SQLError(err)
	}
	return group, nil
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
	data, err := taAPI.NewExerciceData(ct.db, exerciceID)
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

	origin, _ := exerciceOrigin(group, tcAPI.OptionalIdReview{}, userID, ct.admin.Id)
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
	Tags ed.Tags
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
		return errAccessForbidden
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return err
	}

	err = ed.UpdateExercicegroupTags(tx, params.Tags, params.Id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

// EditorDeleteExercice remove the given exercice, also cleaning
// up the exercice group if needed.
// It returns information if the exercice is used in tasks
func (ct *Controller) EditorDeleteExercice(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	idExercice, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	out, err := ct.deleteExercice(ed.IdExercice(idExercice), user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
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
		return errAccessForbidden
	}

	return nil
}

// getExerciceUses returns the item using the given exercice
func getExerciceUses(db ed.DB, id ed.IdExercice) (out QuestionExerciceUses, err error) {
	tas, err := ta.SelectTasksByIdExercices(db, id)
	if err != nil {
		return out, utils.SQLError(err)
	}

	return newQuestionExericeUses(db, tas.IDs())
}

type DeleteExerciceOut struct {
	Deleted   bool
	BlockedBy QuestionExerciceUses // non empty iff Deleted == false
}

func (ct *Controller) deleteExercice(idExercice ed.IdExercice, userID uID) (DeleteExerciceOut, error) {
	if err := ct.checkExerciceOwner(idExercice, userID); err != nil {
		return DeleteExerciceOut{}, err
	}

	uses, err := getExerciceUses(ct.db, idExercice)
	if err != nil {
		return DeleteExerciceOut{}, err
	}
	if len(uses) != 0 {
		return DeleteExerciceOut{
			Deleted:   false,
			BlockedBy: uses,
		}, nil
	}

	links, err := ed.SelectExerciceQuestionsByIdExercices(ct.db, idExercice)
	if err != nil {
		return DeleteExerciceOut{}, utils.SQLError(err)
	}
	qus, err := ed.SelectQuestions(ct.db, links.IdQuestions()...)
	if err != nil {
		return DeleteExerciceOut{}, utils.SQLError(err)
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
		return DeleteExerciceOut{}, utils.SQLError(err)
	}

	// remove the links
	_, err = ed.DeleteExerciceQuestionsByIdExercices(tx, idExercice)
	if err != nil {
		_ = tx.Rollback()
		return DeleteExerciceOut{}, utils.SQLError(err)
	}

	// remove the actual questions
	_, err = ed.DeleteQuestionsByIDs(tx, toDelete...)
	if err != nil {
		_ = tx.Rollback()
		return DeleteExerciceOut{}, utils.SQLError(err)
	}

	// finaly remove the exercice
	_, err = ed.DeleteExerciceById(tx, idExercice)
	if err != nil {
		_ = tx.Rollback()
		return DeleteExerciceOut{}, utils.SQLError(err)
	}

	err = tx.Commit()
	if err != nil {
		return DeleteExerciceOut{}, utils.SQLError(err)
	}

	return DeleteExerciceOut{Deleted: true}, nil
}

type ExerciceWithPreview struct {
	Ex      ExerciceExt
	Preview LoopbackShowExercice
}

type ExerciceCreateQuestionIn struct {
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

	ex, err := ct.createQuestionEx(args, user.Id)
	if err != nil {
		return err
	}

	data, err := taAPI.NewExerciceData(ct.db, args.IdExercice)
	if err != nil {
		return err
	}

	preview, err := newExercicePreview(data, -1)
	if err != nil {
		return err
	}

	out := ExerciceWithPreview{ex, preview}
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

type ExerciceImportQuestionIn struct {
	IdQuestion ed.IdQuestion
	IdExercice ed.IdExercice
}

// EditorExerciceImportQuestion imports an already existing question,
// making a copy and associating it with the given exercice.
// It also updates the preview
func (ct *Controller) EditorExerciceImportQuestion(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	var args ExerciceImportQuestionIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	ex, err := ct.importQuestionEx(args, user.Id)
	if err != nil {
		return err
	}

	data, err := taAPI.NewExerciceData(ct.db, args.IdExercice)
	if err != nil {
		return err
	}

	preview, err := newExercicePreview(data, -1)
	if err != nil {
		return err
	}

	out := ExerciceWithPreview{ex, preview}

	return c.JSON(200, out)
}

func (ct *Controller) importQuestionEx(args ExerciceImportQuestionIn, userID uID) (ExerciceExt, error) {
	if err := ct.checkExerciceOwner(args.IdExercice, userID); err != nil {
		return ExerciceExt{}, err
	}

	existingLinks, err := ed.SelectExerciceQuestionsByIdExercices(ct.db, args.IdExercice)
	if err != nil {
		return ExerciceExt{}, utils.SQLError(err)
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return ExerciceExt{}, utils.SQLError(err)
	}

	// copy the question to import : shallow copy is enough
	question, err := ed.SelectQuestion(ct.db, args.IdQuestion)
	if err != nil {
		return ExerciceExt{}, utils.SQLError(err)
	}
	question.IdGroup = ed.OptionalIdQuestiongroup{}
	question.NeedExercice = args.IdExercice.AsOptional()
	question, err = question.Insert(tx)
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

type ExerciceDuplicateQuestionIn struct {
	QuestionIndex int
	IdExercice    ed.IdExercice
}

// EditorExerciceDuplicateQuestion duplicate the question at the given
// index in the given exercice (variant), also updating the preview
func (ct *Controller) EditorExerciceDuplicateQuestion(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	var args ExerciceDuplicateQuestionIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	ex, err := ct.duplicateQuestionEx(args, user.Id)
	if err != nil {
		return err
	}

	data, err := taAPI.NewExerciceData(ct.db, args.IdExercice)
	if err != nil {
		return err
	}

	preview, err := newExercicePreview(data, -1)
	if err != nil {
		return err
	}

	out := ExerciceWithPreview{ex, preview}
	return c.JSON(200, out)
}

func (ct *Controller) duplicateQuestionEx(args ExerciceDuplicateQuestionIn, userID uID) (ExerciceExt, error) {
	if err := ct.checkExerciceOwner(args.IdExercice, userID); err != nil {
		return ExerciceExt{}, err
	}

	existingLinks, err := ed.SelectExerciceQuestionsByIdExercices(ct.db, args.IdExercice)
	if err != nil {
		return ExerciceExt{}, utils.SQLError(err)
	}
	existingLinks.EnsureOrder()

	if args.QuestionIndex >= len(existingLinks) {
		return ExerciceExt{}, errors.New("internal error: invalid question index")
	}
	link := existingLinks[args.QuestionIndex]

	question, err := ed.SelectQuestion(ct.db, link.IdQuestion)
	if err != nil {
		return ExerciceExt{}, utils.SQLError(err)
	}

	// duplicate the question : shallow copy is enough
	tx, err := ct.db.Begin()
	if err != nil {
		return ExerciceExt{}, utils.SQLError(err)
	}
	question, err = question.Insert(tx)
	if err != nil {
		_ = tx.Rollback()
		return ExerciceExt{}, utils.SQLError(err)
	}

	// and insert it to the current questions
	newLink := ed.ExerciceQuestion{IdExercice: args.IdExercice, IdQuestion: question.Id, Bareme: link.Bareme}
	index := args.QuestionIndex + 1
	existingLinks = append(existingLinks[:index], append([]ed.ExerciceQuestion{newLink}, existingLinks[index:]...)...)

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
}

// EditorExerciceUpdateQuestions updates the question links and
// the preview
func (ct *Controller) EditorExerciceUpdateQuestions(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	var args ExerciceUpdateQuestionsIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	ex, err := ct.updateQuestionsEx(args, user.Id)
	if err != nil {
		return err
	}

	data, err := taAPI.NewExerciceData(ct.db, args.IdExercice)
	if err != nil {
		return err
	}

	preview, err := newExercicePreview(data, -1)
	if err != nil {
		return err
	}

	out := ExerciceWithPreview{ex, preview}
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

	var (
		toDelete       []ed.IdQuestion
		newQuestionIDs = args.Questions.ByIdQuestion()
	)
	for _, link := range links {
		_, willBeUsed := newQuestionIDs[link.IdQuestion]
		if !willBeUsed {
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

type ExerciceUpdateIn = ExerciceHeader

// EditorSaveExerciceMeta update the exercice metadata.
func (ct *Controller) EditorSaveExerciceMeta(c echo.Context) error {
	user := tcAPI.JWTTeacher(c)

	var args ExerciceUpdateIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.updateExercice(args, user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) updateExercice(in ExerciceHeader, userID uID) (ed.Exercice, error) {
	if err := ct.checkExerciceOwner(in.Id, userID); err != nil {
		return ed.Exercice{}, err
	}

	ex, err := ed.SelectExercice(ct.db, in.Id)
	if err != nil {
		return ed.Exercice{}, err
	}

	// only update meta data
	// ex.Description = in.Description
	ex.Subtitle = in.Subtitle
	ex.Difficulty = in.Difficulty
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
	data, err := taAPI.NewExerciceData(ct.db, params.IdExercice)
	if err != nil {
		return CheckExerciceParametersOut{}, err
	}
	qus, _ := data.QuestionsList()

	if L1, L2 := len(params.QuestionParameters), len(qus); L1 != L2 {
		return CheckExerciceParametersOut{}, fmt.Errorf("internal error: mismatched question length (%d != %d)", L1, L2)
	}

	for index := range qus {
		toCheck := params.QuestionParameters[index]
		toCheck = append(toCheck, params.SharedParameters...)

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
	OnlyPreview     bool // if true, skip the save part (Parameters and Questions are thus ignored)
	IdExercice      ed.IdExercice
	Parameters      questions.Parameters // shared parameters
	Questions       []ed.Question        // questions content
	CurrentQuestion int                  // to update the preview accordingly
}

type SaveExerciceAndPreviewOut struct {
	Error questions.ErrQuestionInvalid
	// Index of the error, meaningful only if Error is non empty
	QuestionIndex int

	IsValid bool
	Preview LoopbackShowExercice
}

func (ct *Controller) saveExerciceAndPreview(params SaveExerciceAndPreviewIn, userID uID) (SaveExerciceAndPreviewOut, error) {
	data, err := taAPI.NewExerciceData(ct.db, params.IdExercice)
	if err != nil {
		return SaveExerciceAndPreviewOut{}, err
	}
	ex := &data.Exercice

	// if the exercice is in review, allow external user to preview it
	_, inReview, err := reviews.SelectReviewExerciceByIdExercice(ct.db, data.Group.Id)
	if err != nil {
		return SaveExerciceAndPreviewOut{}, utils.SQLError(err)
	}

	if !inReview && !data.Group.IsVisibleBy(userID) {
		return SaveExerciceAndPreviewOut{}, errAccessForbidden
	}

	qus, _ := data.QuestionsList()
	// validate all the questions, using shared parameters if needed
	if !params.OnlyPreview {
		if len(params.Questions) != len(qus) {
			return SaveExerciceAndPreviewOut{}, errors.New("internal error: inconsistent questions length")
		}
		for index, question := range qus {
			toCheck := params.Questions[index].Parameters
			if question.NeedExercice.Valid { // add the shared parameters
				toCheck = append(toCheck, params.Parameters...)
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
			qu.Enonce = incomming.Enonce
			qu.Parameters = incomming.Parameters
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
	}

	preview, err := newExercicePreview(data, params.CurrentQuestion)
	if err != nil {
		return SaveExerciceAndPreviewOut{}, err
	}

	return SaveExerciceAndPreviewOut{IsValid: true, Preview: preview}, nil
}

// newExercicePreview instantiates the exercice and return preview data
// [nextQuestion] is the index of the question to show in the preview,
// or -1 for the summary
func newExercicePreview(content taAPI.ExerciceData, nextQuestion int) (LoopbackShowExercice, error) {
	instance, err := content.Instantiate()
	if err != nil {
		return LoopbackShowExercice{}, err
	}

	qus, _ := content.QuestionsList()
	questionOrigins := make([]questions.QuestionPage, len(qus))
	for i, qu := range qus {
		questionOrigins[i] = questions.QuestionPage{Enonce: qu.Enonce, Parameters: qu.Parameters}
	}

	if nextQuestion >= len(instance.Questions) {
		return LoopbackShowExercice{}, fmt.Errorf("internal error: invalid question index %d", nextQuestion)
	}
	progression := make([]ta.QuestionHistory, len(instance.Questions))
	// mark previous question as validated for better consistency on the preview
	for i := 0; i < nextQuestion; i++ {
		progression[i] = ta.QuestionHistory{true}
	}
	return LoopbackShowExercice{
		Exercice: instance,
		Progression: taAPI.ProgressionExt{
			NextQuestion: nextQuestion,
			Questions:    progression,
		},
		Origin: questionOrigins,
	}, nil
}

type ExportExerciceLatexIn struct {
	Parameters questions.Parameters     // shared parameters
	Questions  []questions.QuestionPage // questions content
}

type ExportExerciceLatexOut struct {
	Error         questions.ErrQuestionInvalid
	QuestionIndex int // of the error
	IsValid       bool

	Latex string
}

// EditorExerciceExportLateX instantiate the given exercice and generates a LaTeX version,
// returning the code as a string
func (ct *Controller) EditorExerciceExportLateX(c echo.Context) error {
	var args ExportExerciceLatexIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	out, err := exportExerciceLatex(args)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func exportExerciceLatex(exercice ExportExerciceLatexIn) (ExportExerciceLatexOut, error) {
	ques := make([]questions.QuestionInstance, len(exercice.Questions))
	for index, question := range exercice.Questions {
		toCheck := question.Parameters
		toCheck = append(toCheck, exercice.Parameters...) // add the shared parameters

		err := toCheck.Validate()
		if err != nil {
			return ExportExerciceLatexOut{
				Error:         err.(questions.ErrQuestionInvalid),
				QuestionIndex: index,
			}, nil
		}

		instanceParams, err := toCheck.ToMap().Instantiate()
		if err != nil {
			return ExportExerciceLatexOut{}, err
		}
		ques[index], err = question.Enonce.InstantiateWith(instanceParams)
		if err != nil {
			return ExportExerciceLatexOut{}, err
		}
	}

	return ExportExerciceLatexOut{
		IsValid: true,
		Latex:   questions.InstancesToLatex(ques),
	}, nil
}
