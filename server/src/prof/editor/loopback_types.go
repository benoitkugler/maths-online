package editor

import (
	"github.com/benoitkugler/maths-online/maths/questions/client"
	"github.com/benoitkugler/maths-online/tasks"
)

type LoopbackServerEvent interface {
	isLoopbackServerEvent()
}

type loopbackPaused struct{}

type loopbackQuestion struct {
	Question client.Question `gomacro-extern:"client#dart#questions/types.gen.dart"`
}

type loopbackQuestionValidOut struct {
	Answers client.QuestionAnswersOut `gomacro-extern:"client#dart#questions/types.gen.dart"`
}

type loopbackQuestionCorrectAnswersOut struct {
	Answers client.QuestionAnswersIn `gomacro-extern:"client#dart#questions/types.gen.dart"`
}

type loopbackShowExercice struct {
	Exercice    tasks.InstantiatedWork `gomacro-extern:"tasks#dart#shared_gen.dart"`
	Progression tasks.ProgressionExt   `gomacro-extern:"tasks#dart#shared_gen.dart"`
}

type loopbackExerciceCorrectAnswersOut struct {
	Answers       client.QuestionAnswersIn `gomacro-extern:"client#dart#questions/types.gen.dart"`
	QuestionIndex int
}

func (loopbackPaused) isLoopbackServerEvent()                    {}
func (loopbackQuestion) isLoopbackServerEvent()                  {}
func (loopbackQuestionValidOut) isLoopbackServerEvent()          {}
func (loopbackQuestionCorrectAnswersOut) isLoopbackServerEvent() {}
func (loopbackShowExercice) isLoopbackServerEvent()              {}
func (loopbackExerciceCorrectAnswersOut) isLoopbackServerEvent() {}

type LoopbackClientEvent interface {
	isLoopbackClientEvent()
}

type loopbackPing struct{}

type loopbackQuestionValidIn struct {
	Answers client.QuestionAnswersIn `gomacro-extern:"client#dart#questions/types.gen.dart"`
}

type loopbackQuestionCorrectAnswersIn struct{}

type loopbackExerciceCorrectAnswsersIn struct {
	QuestionIndex int
}

func (loopbackPing) isLoopbackClientEvent()                      {}
func (loopbackQuestionValidIn) isLoopbackClientEvent()           {}
func (loopbackQuestionCorrectAnswersIn) isLoopbackClientEvent()  {}
func (loopbackExerciceCorrectAnswsersIn) isLoopbackClientEvent() {}
