package editor

import (
	"fmt"
	"sort"

	"github.com/benoitkugler/maths-online/server/src/prof/teacher"
	ed "github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/utils"
	"github.com/labstack/echo/v4"
)

// EditorGetTags return all tags currently used by questions.
// It also add the special level tags.
func (ct *Controller) EditorGetTags(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	filtred, err := LoadTags(ct.db, user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, filtred)
}

// LoadTags returns all the tags visible by [userID], merging
// questions and exercices.
func LoadTags(db ed.DB, userID uID) ([]string, error) {
	questionTags, err := ed.SelectAllQuestiongroupTags(db)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	exerciceTags, err := ed.SelectAllExercicegroupTags(db)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	// only return tags used by visible groups
	questionGroups, err := ed.SelectAllQuestiongroups(db)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	exerciceGroups, err := ed.SelectAllExercicegroups(db)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	// add the special difficulty and level tags among the proposition,
	// in first choices
	var (
		filtred []string
		allTags = map[string]bool{}
	)
	for _, tag := range questionTags {
		if !questionGroups[tag.IdQuestiongroup].IsVisibleBy(userID) {
			continue
		}
		allTags[tag.Tag] = true
	}
	for _, tag := range exerciceTags {
		if !exerciceGroups[tag.IdExercicegroup].IsVisibleBy(userID) {
			continue
		}
		allTags[tag.Tag] = true
	}

	for tag := range allTags {
		switch tag {
		// case string(Diff1), string(Diff2), string(Diff3): // added after
		case string(ed.Seconde), string(ed.Premiere), string(ed.Terminale): // added after
		default:
			filtred = append(filtred, tag)
		}
	}

	// sort by name but make sure special tags come first
	sort.Strings(filtred)

	filtred = append([]string{
		// string(Diff1), string(Diff2), string(Diff3),
		string(ed.Seconde), string(ed.Premiere), string(ed.Terminale),
	}, filtred...)

	return filtred, nil
}

func (ct *Controller) EditorCheckExerciceParameters(c echo.Context) error {
	var args CheckExerciceParametersIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	out, err := ct.checkExerciceParameters(args)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

type ExerciceUpdateVisiblityIn struct {
	ID     ed.IdExercicegroup
	Public bool
}

func (ct *Controller) EditorUpdateExercicegroupVis(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	// we only accept public question from admin account
	if user.Id != ct.admin.Id {
		return accessForbidden
	}

	var args ExerciceUpdateVisiblityIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	ex, err := ed.SelectExercicegroup(ct.db, args.ID)
	if err != nil {
		return utils.SQLError(err)
	}
	if ex.IdTeacher != user.Id {
		return accessForbidden
	}

	if !args.Public {
		// TODO: check that it is not harmful to hide the exercice group again
	}
	ex.Public = args.Public
	ex, err = ex.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	return c.NoContent(200)
}

// For non personnal questions, only preview.
func (ct *Controller) EditorSaveExerciceAndPreview(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	var args SaveExerciceAndPreviewIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	out, err := ct.saveExerciceAndPreview(args, user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}
