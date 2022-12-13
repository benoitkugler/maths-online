package editor

import (
	"github.com/benoitkugler/maths-online/maths/questions"
	"github.com/benoitkugler/maths-online/maths/questions/client"
	"github.com/benoitkugler/maths-online/tasks"
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
	Question client.Question `gomacro-extern:"client#dart#../questions/types.gen.dart"`
	Params   tasks.Params    `gomacro-extern:"tasks#dart#../shared_gen.dart"`

	Origin questions.QuestionPage `gomacro-extern:"questions#dart#editor_question.dart"`
}

type LoopbackShowExercice struct {
	Exercice    tasks.InstantiatedWork `gomacro-extern:"tasks#dart#../shared_gen.dart"`
	Progression tasks.ProgressionExt   `gomacro-extern:"tasks#dart#../shared_gen.dart"`

	Origin []questions.QuestionPage `gomacro-extern:"questions#dart#editor_question.dart"`
}

func (LoopbackPaused) isLoopbackServerEvent()       {}
func (LoopbackShowQuestion) isLoopbackServerEvent() {}
func (LoopbackShowExercice) isLoopbackServerEvent() {}

type LoopackEvaluateQuestionIn struct {
	Question questions.QuestionPage `gomacro-extern:"questions#dart#editor_question.dart"`
	Answer   tasks.Answer
}

type LoopbackEvaluateQuestionOut struct {
	Answers client.QuestionAnswersOut `gomacro-extern:"client#dart#../questions/types.gen.dart"`
}

type LoopbackShowQuestionAnswerIn struct {
	Question questions.QuestionPage `gomacro-extern:"questions#dart#editor_question.dart"`
	Params   tasks.Params
}

type LoopbackShowQuestionAnswerOut struct {
	Answers client.QuestionAnswersIn `gomacro-extern:"client#dart#../questions/types.gen.dart"`
}
type loopbackExerciceCorrectAnswersOut struct {
	Answers       client.QuestionAnswersIn `gomacro-extern:"client#dart#../questions/types.gen.dart"`
	QuestionIndex int
}
