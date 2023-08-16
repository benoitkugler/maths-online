package editor

import (
	"fmt"
	"sort"

	"github.com/benoitkugler/maths-online/server/src/maths/questions"
	tcAPI "github.com/benoitkugler/maths-online/server/src/prof/teacher"
	ed "github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/utils"
	"github.com/labstack/echo/v4"
)

// EditorGetTags return all tags currently used by questions.
// It also add the special level tags.
func (ct *Controller) EditorGetTags(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	filtred, err := LoadTags(ct.db, userID)
	if err != nil {
		return err
	}

	return c.JSON(200, filtred)
}

type TagsDB struct {
	Levels           []string
	ChaptersByLevel  map[string][]string            // level -> chapters
	TrivByChapters   map[string]map[string][]string // level -> chapter -> triv maths
	SubLevelsByLevel map[string][]string            // level -> sublevels
}

// LoadTags returns all the tags visible by [userID], merging
// questions and exercices.
func LoadTags(db ed.DB, userID uID) (out TagsDB, _ error) {
	// only return tags used by visible groups
	questionGroups, err := ed.SelectAllQuestiongroups(db)
	if err != nil {
		return out, utils.SQLError(err)
	}
	exerciceGroups, err := ed.SelectAllExercicegroups(db)
	if err != nil {
		return out, utils.SQLError(err)
	}

	questionGroups.RestrictVisible(userID)
	exerciceGroups.RestrictVisible(userID)

	questionTags, err := ed.SelectQuestiongroupTagsByIdQuestiongroups(db, questionGroups.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}
	exerciceTags, err := ed.SelectExercicegroupTagsByIdExercicegroups(db, exerciceGroups.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	// build the map between level to chapters
	indexes := append(
		questionsToTagGroups(questionGroups, questionTags),
		exercicesToTagGroups(exerciceGroups, exerciceTags)...,
	)
	// add common suggestions
	indexes = append(indexes,
		ed.TagGroup{TagIndex: ed.TagIndex{Level: ed.Seconde}},
		ed.TagGroup{TagIndex: ed.TagIndex{Level: ed.Premiere}},
		ed.TagGroup{TagIndex: ed.TagIndex{Level: ed.Terminale}},
		ed.TagGroup{TagIndex: ed.TagIndex{Level: ed.CPGE}},
	)

	return buildTagsDB(indexes), nil
}

type GenerateSyntaxHintIn struct {
	Block              questions.ExpressionFieldBlock
	SharedParameters   questions.Parameters // empty for standalone questions
	QuestionParameters questions.Parameters
}

// EditorGenerateSyntaxHint generate a TextBlock with hints about the command
// used in the given expression field
func (ct *Controller) EditorGenerateSyntaxHint(c echo.Context) error {
	var args GenerateSyntaxHintIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	args.SharedParameters = append(args.SharedParameters, args.QuestionParameters...)
	out, err := args.Block.SyntaxHint(args.SharedParameters)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
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

type Query struct {
	TitleQuery   string   // empty means all
	LevelTags    []string // union, empty means all; an empty tag means "with no level"
	ChapterTags  []string // union, empty means all; an empty tag means "with no chapter"
	SubLevelTags []string // union, empty means all
	Origin       OriginKind
}

func (query Query) normalize() {
	// normalize query
	for i, t := range query.LevelTags {
		query.LevelTags[i] = ed.NormalizeTag(t)
	}
	for i, t := range query.ChapterTags {
		query.ChapterTags[i] = ed.NormalizeTag(t)
	}
	for i, t := range query.SubLevelTags {
		query.SubLevelTags[i] = ed.NormalizeTag(t)
	}
}

func (query Query) matchOrigin(vis tcAPI.Visibility) bool {
	switch query.Origin {
	case OnlyPersonnal:
		return vis == tcAPI.Personnal
	case OnlyAdmin:
		return vis == tcAPI.Admin
	default:
		return true
	}
}

func (query Query) matchLevel(level ed.LevelTag) bool {
	if len(query.LevelTags) == 0 {
		return true
	}
	for _, tag := range query.LevelTags {
		if tag == string(level) {
			return true
		}
	}
	return false
}

func (query Query) matchChapter(chapter string) bool {
	if len(query.ChapterTags) == 0 {
		return true
	}
	for _, tag := range query.ChapterTags {
		if tag == string(chapter) {
			return true
		}
	}
	return false
}

func (query Query) matchSubLevel(subLevels []string) bool {
	if len(query.SubLevelTags) == 0 {
		return true
	}
	for _, tag := range query.SubLevelTags {
		for _, resourceTag := range subLevels {
			if tag == resourceTag {
				return true
			}
		}
	}
	return false
}

type ExerciceUpdateVisiblityIn struct {
	ID     ed.IdExercicegroup
	Public bool
}

func (ct *Controller) EditorUpdateExercicegroupVis(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	// we only accept public question from admin account
	if userID != ct.admin.Id {
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
	if ex.IdTeacher != userID {
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
	userID := tcAPI.JWTTeacher(c)

	var args SaveExerciceAndPreviewIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	out, err := ct.saveExerciceAndPreview(args, userID)
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
		out = append(out, ls.Tags().BySection().TagIndex)
	}
	return out
}

func exercicesToIndex(exercices ed.Exercicegroups, tags ed.ExercicegroupTags) []ed.TagIndex {
	m := tags.ByIdExercicegroup()
	out := make([]ed.TagIndex, 0, len(m))
	for id := range exercices {
		ls := m[id]
		out = append(out, ls.Tags().BySection().TagIndex)
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
		sort.Slice(items.Chapters, func(i, j int) bool { return items.Chapters[i].Chapter > items.Chapters[j].Chapter })
		out = append(out, items)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Level > out[j].Level })
	return out
}

func questionsToTagGroups(questions ed.Questiongroups, tags ed.QuestiongroupTags) []ed.TagGroup {
	m := tags.ByIdQuestiongroup()
	out := make([]ed.TagGroup, 0, len(m))
	for id := range questions {
		ls := m[id]
		out = append(out, ls.Tags().BySection())
	}
	return out
}

func exercicesToTagGroups(exercices ed.Exercicegroups, tags ed.ExercicegroupTags) []ed.TagGroup {
	m := tags.ByIdExercicegroup()
	out := make([]ed.TagGroup, 0, len(m))
	for id := range exercices {
		ls := m[id]
		out = append(out, ls.Tags().BySection())
	}
	return out
}

func buildTagsDB(tags []ed.TagGroup) (out TagsDB) {
	tmp := make(map[ed.LevelTag]map[string]map[string]bool) // level -> chapter -> trivs (maybe empty)
	tmp2 := make(map[ed.LevelTag]map[string]bool)           // level -> sublevels
	for _, tag := range tags {
		level, chapter := tag.Level, tag.Chapter
		m1 := tmp[level]
		if m1 == nil {
			m1 = make(map[string]map[string]bool)
		}
		m2 := m1[chapter]
		if m2 == nil {
			m2 = make(map[string]bool)
		}
		for _, triv := range tag.TrivMaths {
			m2[triv] = true
		}
		m1[chapter] = m2
		tmp[level] = m1

		mSub := tmp2[level]
		if mSub == nil {
			mSub = make(map[string]bool)
		}
		for _, sub := range tag.SubLevels {
			mSub[sub] = true
		}
		tmp2[level] = mSub
	}
	out.ChaptersByLevel = make(map[string][]string)
	out.TrivByChapters = make(map[string]map[string][]string)
	for level, m := range tmp {
		level := string(level)
		// fill the level
		if level != "" {
			out.Levels = append(out.Levels, level)
		}
		// fill the chapters
		trivOut := out.TrivByChapters[level]
		if trivOut == nil {
			trivOut = make(map[string][]string)
		}
		for chapter, trivs := range m {
			if chapter != "" {
				out.ChaptersByLevel[level] = append(out.ChaptersByLevel[level], chapter)
			}

			// fill the trivMaths
			for triv := range trivs {
				trivOut[chapter] = append(trivOut[chapter], triv)
			}

			sort.Strings(trivOut[chapter])
		}
		out.TrivByChapters[level] = trivOut

		sort.Strings(out.ChaptersByLevel[level])
	}
	sort.Strings(out.Levels)

	out.SubLevelsByLevel = make(map[string][]string)
	for level, subs := range tmp2 {
		l := make([]string, 0, len(subs))
		for sub := range subs {
			l = append(l, sub)
		}
		sort.Strings(l)
		out.SubLevelsByLevel[string(level)] = l
	}

	return out
}
