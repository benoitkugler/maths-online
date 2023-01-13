package reviews

import (
	edAPI "github.com/benoitkugler/maths-online/server/src/prof/editor"
	"github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/sql/trivial"
	trivGame "github.com/benoitkugler/maths-online/server/src/trivial"
)

type TargetContent interface {
	isTargetContent()
}

func (TargetTrivial) isTargetContent()  {}
func (TargetQuestion) isTargetContent() {}
func (TargetExercice) isTargetContent() {}

type TargetTrivial struct {
	Config trivial.Trivial

	NbQuestionsByCategories [trivGame.NbCategories]int
}

type TargetQuestion struct {
	Group    edAPI.QuestiongroupExt
	Variants []editor.Question
	AllTags  map[editor.Section][]string
}

type TargetExercice struct {
	Group   edAPI.ExercicegroupExt
	AllTags map[editor.Section][]string
}

type LoadTargetOut struct {
	Content TargetContent
}
