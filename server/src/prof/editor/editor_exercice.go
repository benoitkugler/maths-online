package editor

import (
	"sort"

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

type ExerciceQuestionExt struct {
	Title    string
	Question ExerciceQuestion
}

func fillQuestions(l ExerciceQuestions, dict Questions) []ExerciceQuestionExt {
	l.ensureIndex()
	out := make([]ExerciceQuestionExt, len(l))
	for i, qu := range l {
		out[i].Question = qu
		out[i].Title = dict[qu.IdQuestion].Page.Title
	}
	return out
}

type ExerciceExt struct {
	Exercice  Exercice
	Origin    teacher.Origin
	Questions []ExerciceQuestionExt
}

func (ct *Controller) ExercicesGetList(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	out, err := ct.getExercices(user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) getExercices(userID int64) ([]ExerciceExt, error) {
	exs, err := SelectAllExercices(ct.db)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	var (
		out []ExerciceExt
		tmp IDs
	)
	for _, ex := range exs {
		vis, ok := teacher.NewVisibility(ex.IdTeacher, userID, ct.admin.Id, ex.Public)
		if !ok {
			continue
		}
		out = append(out, ExerciceExt{
			Exercice: ex,
			Origin: teacher.Origin{
				AllowPublish: userID == ct.admin.Id,
				IsPublic:     ex.Public,
				Visibility:   vis,
			},
		})
		tmp = append(tmp, ex.Id)
	}

	links, err := SelectExerciceQuestionsByIdExercices(ct.db, tmp...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	dict := links.ByIdExercice()

	qus, err := SelectQuestions(ct.db, links.IdQuestions()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	for i, ex := range out {
		s := dict[ex.Exercice.Id]
		out[i].Questions = fillQuestions(s, qus)
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Exercice.Id < out[j].Exercice.Id })
	sort.SliceStable(out, func(i, j int) bool { return out[i].Exercice.Title < out[j].Exercice.Title })

	return out, nil
}

func (ct *Controller) ExerciceCreate(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	out, err := ct.createExercice(user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) createExercice(userID int64) (ExerciceExt, error) {
	ex, err := Exercice{IdTeacher: userID, Flow: Parallel}.Insert(ct.db)
	if err != nil {
		return ExerciceExt{}, utils.SQLError(err)
	}

	out := ExerciceExt{
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

func (ct *Controller) checkAcces(idExercice, userID int64) error {
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
	if err := ct.checkAcces(idExercice, userID); err != nil {
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
	if err := ct.checkAcces(args.IdExercice, userID); err != nil {
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

	out, err := updateExerciceQuestionList(ct.db, args.IdExercice, existing)
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
	if err := ct.checkAcces(args.IdExercice, userID); err != nil {
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

	out, err := updateExerciceQuestionList(ct.db, args.IdExercice, args.Questions)
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
	if err := ct.checkAcces(in.Id, userID); err != nil {
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
