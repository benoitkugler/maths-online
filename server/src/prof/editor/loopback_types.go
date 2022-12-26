package editor

import (
	"github.com/benoitkugler/maths-online/server/src/maths/questions"
	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
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

	Origin questions.QuestionPage
}

type LoopbackShowExercice struct {
	Exercice    tasks.InstantiatedWork `gomacro-opaque:"typescript"`
	Progression tasks.ProgressionExt   `gomacro-opaque:"typescript"`

	Origin []questions.QuestionPage
}

func (LoopbackPaused) isLoopbackServerEvent()       {}
func (LoopbackShowQuestion) isLoopbackServerEvent() {}
func (LoopbackShowExercice) isLoopbackServerEvent() {}

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
