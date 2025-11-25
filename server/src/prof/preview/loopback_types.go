package preview

import (
	"github.com/benoitkugler/maths-online/server/src/maths/questions"
	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	"github.com/benoitkugler/maths-online/server/src/sql/ceintures"
	"github.com/benoitkugler/maths-online/server/src/tasks"
)

// LoopbackServerEvent describes an event triggered
// by the editor web app (usually in response to a server call)
// It is emitted in the HTML/Javascript side and received in the Dart
// code.
type LoopbackServerEvent interface {
	isLoopbackServerEvent()
}

type LoopbackPaused struct{}

type LoopbackShowQuestion struct {
	Question client.Question `gomacro-opaque:"typescript"`
	Params   tasks.Params    `gomacro-opaque:"typescript"`

	// Set the initial view to display the correction,
	// instead of the enonce.
	ShowCorrection bool

	Origin questions.QuestionPage
}

type LoopbackShowExercice struct {
	Exercice     tasks.InstantiatedWork `gomacro-opaque:"typescript"`
	Progression  tasks.Progression      `gomacro-opaque:"typescript"`
	NextQuestion int

	// Set the initial view to display the correction,
	// instead of the enonce.
	ShowCorrection bool

	Origin []questions.QuestionPage
}

type LoopbackShowCeinture struct {
	Questions []tasks.InstantiatedBeltQuestion `gomacro-opaque:"typescript"`

	QuestionIndex int

	Origin []questions.QuestionPage

	// Set the initial view to display the correction,
	// instead of the enonce.
	ShowCorrection bool
}

func (LoopbackPaused) isLoopbackServerEvent()       {}
func (LoopbackShowQuestion) isLoopbackServerEvent() {}
func (LoopbackShowExercice) isLoopbackServerEvent() {}
func (LoopbackShowCeinture) isLoopbackServerEvent() {}

type LoopackEvaluateQuestionIn struct {
	Question questions.QuestionPage
	Answer   tasks.AnswerP `gomacro-opaque:"typescript"`
}

type LoopbackEvaluateQuestionOut struct {
	Answers client.QuestionAnswersOut `gomacro-opaque:"typescript"`
}

type LoopbackShowQuestionAnswerIn struct {
	Question questions.QuestionPage
	Params   tasks.Params `gomacro-opaque:"typescript"`
}

type LoopbackShowQuestionAnswerOut struct {
	Answers client.QuestionAnswersIn `gomacro-opaque:"typescript"`
}

type LoopbackEvaluateCeintureIn struct {
	Questions []ceintures.IdBeltquestion // as send by the server
	Answers   []tasks.AnswerP            `gomacro-opaque:"typescript"`
}

type LoopbackEvaluateCeintureOut struct {
	Answers []client.QuestionAnswersOut `gomacro-opaque:"typescript"`
}
