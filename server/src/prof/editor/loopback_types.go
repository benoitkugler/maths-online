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
	Question client.Question `gomacro-extern:"client:dart:questions/types.gen.dart"`
}

type loopbackQuestionValidOut struct {
	Answers client.QuestionAnswersOut `gomacro-extern:"client:dart:questions/types.gen.dart"`
}

type loopbackQuestionCorrectAnswersOut struct {
	Answers client.QuestionAnswersIn `gomacro-extern:"client:dart:questions/types.gen.dart"`
}

type loopbackShowExercice struct {
	Exercice    tasks.InstantiatedWork `gomacro-extern:"editor:dart:shared_gen.dart"`
	Progression tasks.ProgressionExt   `gomacro-extern:"editor:dart:shared_gen.dart"`
}

func (loopbackPaused) isLoopbackServerEvent()                    {}
func (loopbackQuestion) isLoopbackServerEvent()                  {}
func (loopbackQuestionValidOut) isLoopbackServerEvent()          {}
func (loopbackQuestionCorrectAnswersOut) isLoopbackServerEvent() {}
func (loopbackShowExercice) isLoopbackServerEvent()              {}

type LoopbackClientEvent interface {
	isLoopbackClientEvent()
}

type loopbackPing struct{}

type loopbackQuestionValidIn struct {
	Answers client.QuestionAnswersIn `gomacro-extern:"client:dart:questions/types.gen.dart"`
}

type loopbackQuestionCorrectAnswersIn struct{}

type loopbackExerciceValidIn struct {
	Answer tasks.EvaluateWorkIn `gomacro-extern:"editor:dart:shared_gen.dart"`
}

func (loopbackPing) isLoopbackClientEvent()                     {}
func (loopbackQuestionValidIn) isLoopbackClientEvent()          {}
func (loopbackQuestionCorrectAnswersIn) isLoopbackClientEvent() {}
func (loopbackExerciceValidIn) isLoopbackClientEvent()          {}
