package reviews

import (
	edAPI "github.com/benoitkugler/maths-online/server/src/prof/editor"
	"github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/sql/trivial"
)

type TargetContent interface {
	isTargetContent()
}

func (TargetTrivial) isTargetContent()  {}
func (TargetQuestion) isTargetContent() {}
func (TargetExercice) isTargetContent() {}

type TargetTrivial struct {
	Config trivial.Trivial

	// TODO: add question numbers
}

type TargetQuestion struct {
	Group    edAPI.QuestiongroupExt
	Variants []editor.Question
	AllTags  []string
}

type TargetExercice struct {
	Group   edAPI.ExercicegroupExt
	AllTags []string
}

type LoadTargetOut struct {
	Content TargetContent
}
