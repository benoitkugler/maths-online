package editor

import (
	"sort"

	"github.com/benoitkugler/maths-online/prof/teacher"
	"github.com/benoitkugler/maths-online/utils"
	"github.com/labstack/echo/v4"
)

// CRUD

type ExerciceExt struct {
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

	for i, ex := range out {
		s := dict[ex.Exercice.Id]
		out[i].Questions = s
		sort.Slice(s, func(i, j int) bool { return s[i].index < s[j].index })
	}

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
	ex, err := Exercice{IdTeacher: userID}.Insert(ct.db)
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

type ExerciceAddQuestionIn struct {
	IdExercice int64
	IdQuestion int64 // < 0 means create a new question
}

// ExerciceAddQuestion creates or importd a question and appends it
// to the given exercice.
func (ct *Controller) ExerciceAddQuestion(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	var args ExerciceAddQuestionIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.addQuestionToExercice(args, user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) addQuestionToExercice(args ExerciceAddQuestionIn, userID int64) (ExerciceQuestions, error) {
	if err := ct.checkAcces(args.IdExercice, userID); err != nil {
		return nil, err
	}

	// creates a question if needed
	if args.IdQuestion < 0 {
		question, err := Question{IdTeacher: userID, Public: false, NeedExercice: true}.Insert(ct.db)
		if err != nil {
			return nil, utils.SQLError(err)
		}
		args.IdQuestion = question.Id
	}

	// make sure the question is visible by the user
	question, err := SelectQuestion(ct.db, args.IdQuestion)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	if !question.IsVisibleBy(userID) {
		return nil, accessForbidden
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return nil, utils.SQLError(err)
	}

	existing, err := DeleteExerciceQuestionsByIdExercices(tx, args.IdExercice)
	if err != nil {
		_ = tx.Rollback()
		return nil, utils.SQLError(err)
	}

	existing = append(existing, ExerciceQuestion{IdExercice: args.IdExercice, IdQuestion: args.IdQuestion, Bareme: 1})
	err = insertExerciceQuestionList(tx, existing)
	if err != nil {
		_ = tx.Rollback()
		return nil, utils.SQLError(err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, utils.SQLError(err)
	}

	return existing, nil
}

type ExerciceRemoveQuestionIn struct {
	IdExercice int64
	IdQuestion int64
}

func (ct *Controller) ExerciceRemoveQuestion(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	var args ExerciceRemoveQuestionIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.removeQuestionFromExercice(args, user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) removeQuestionFromExercice(args ExerciceRemoveQuestionIn, userID int64) (ExerciceQuestions, error) {
	if err := ct.checkAcces(args.IdExercice, userID); err != nil {
		return nil, err
	}

	tx, err := ct.db.Begin()
	if err != nil {
		return nil, utils.SQLError(err)
	}

	current, err := DeleteExerciceQuestionsByIdExercices(tx, args.IdExercice)
	if err != nil {
		_ = tx.Rollback()
		return nil, utils.SQLError(err)
	}
	newList := make(ExerciceQuestions, 0, len(current)-1)
	for _, v := range current {
		if v.IdQuestion != args.IdQuestion {
			newList = append(newList, v)
		}
	}
	err = insertExerciceQuestionList(tx, newList)
	if err != nil {
		_ = tx.Rollback()
		return nil, utils.SQLError(err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, utils.SQLError(err)
	}

	return newList, nil
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
