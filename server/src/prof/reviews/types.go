package reviews

import (
	edAPI "github.com/benoitkugler/maths-online/server/src/prof/editor"
	"github.com/benoitkugler/maths-online/server/src/prof/homework"
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
func (TargetSheet) isTargetContent()    {}

type TargetTrivial struct {
	Config trivial.Trivial

	NbQuestionsByCategories [trivGame.NbCategories]int
}

type TargetQuestion struct {
	Group    edAPI.QuestiongroupExt
	Variants []editor.Question
	AllTags  edAPI.TagsDB
}

type TargetExercice struct {
	Group   edAPI.ExercicegroupExt
	AllTags edAPI.TagsDB
}

type TargetSheet struct {
	Sheet homework.SheetExt
}

type LoadTargetOut struct {
	Content TargetContent
}
