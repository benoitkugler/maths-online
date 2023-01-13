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
func LoadTags(db ed.DB, userID uID) (map[ed.Section][]string, error) {
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

	allTags := map[ed.TagSection]bool{}

	for _, tag := range questionTags {
		if !questionGroups[tag.IdQuestiongroup].IsVisibleBy(userID) {
			continue
		}
		allTags[ed.TagSection{Tag: tag.Tag, Section: tag.Section}] = true
	}
	for _, tag := range exerciceTags {
		if !exerciceGroups[tag.IdExercicegroup].IsVisibleBy(userID) {
			continue
		}
		allTags[ed.TagSection{Tag: tag.Tag, Section: tag.Section}] = true
	}

	// add common suggestions
	allTags[ed.TagSection{Tag: string(ed.Seconde), Section: ed.Level}] = true
	allTags[ed.TagSection{Tag: string(ed.Premiere), Section: ed.Level}] = true
	allTags[ed.TagSection{Tag: string(ed.Terminale), Section: ed.Level}] = true
	allTags[ed.TagSection{Tag: string(ed.CPGE), Section: ed.Level}] = true

	out := make(map[ed.Section][]string)
	for tag := range allTags {
		out[tag.Section] = append(out[tag.Section], tag.Tag)
	}

	// sort by name
	for _, l := range out {
		sort.Strings(l)
	}

	return out, nil
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
		return errAccessForbidden
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
		return errAccessForbidden
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

// the question/exercice for one level
type LevelItems struct {
	Level    ed.LevelTag
	Chapters []ChapterItems // sorted by Chapter
}

// the question/exercice for one chapter
type ChapterItems struct {
	Chapter    string
	GroupCount int
}

// Index exposes the structure of the resources
type Index []LevelItems

func questionsToIndex(questions ed.Questiongroups, tags ed.QuestiongroupTags) []ed.TagIndex {
	m := tags.ByIdQuestiongroup()
	out := make([]ed.TagIndex, 0, len(m))
	for id := range questions {
		ls := m[id]
		out = append(out, ls.Tags().Index())
	}
	return out
}

func exercicesToIndex(exercices ed.Exercicegroups, tags ed.ExercicegroupTags) []ed.TagIndex {
	m := tags.ByIdExercicegroup()
	out := make([]ed.TagIndex, 0, len(m))
	for id := range exercices {
		ls := m[id]
		out = append(out, ls.Tags().Index())
	}
	return out
}

func buildIndex(tags []ed.TagIndex) Index {
	tmp := map[ed.LevelTag]map[string]int{}
	for _, item := range tags {
		byLevel := tmp[item.Level]
		if byLevel == nil {
			byLevel = make(map[string]int)
		}
		byLevel[item.Chapter] = byLevel[item.Chapter] + 1
		tmp[item.Level] = byLevel
	}
	out := make(Index, 0, len(tmp))
	for level, byLevel := range tmp {
		items := LevelItems{
			Level:    level,
			Chapters: make([]ChapterItems, 0, len(byLevel)),
		}
		for chapter, count := range byLevel {
			items.Chapters = append(items.Chapters, ChapterItems{Chapter: chapter, GroupCount: count})
		}
		sort.Slice(items.Chapters, func(i, j int) bool { return items.Chapters[i].Chapter < items.Chapters[j].Chapter })
		out = append(out, items)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Level < out[j].Level })
	return out
}
