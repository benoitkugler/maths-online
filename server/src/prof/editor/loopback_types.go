package editor

import "github.com/benoitkugler/maths-online/maths/questions/client"

type LoopbackServerEvent interface {
	isLoopbackServerEvent()
}

type loopbackPaused struct{}

type loopbackQuestion struct {
	Question client.Question `dart-extern:"client:questions/types.gen.dart"`
}

type loopbackQuestionValidOut struct {
	Answers client.QuestionAnswersOut `dart-extern:"client:questions/types.gen.dart"`
}

type loopbackQuestionCorrectAnswersOut struct {
	Answers client.QuestionAnswersIn `dart-extern:"client:questions/types.gen.dart"`
}

type loopbackShowExercice struct {
	Exercice    InstantiatedExercice `dart-extern:"editor:shared_gen.dart"`
	Progression ProgressionExt       `dart-extern:"editor:shared_gen.dart"`
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
	Answers client.QuestionAnswersIn `dart-extern:"client:questions/types.gen.dart"`
}

type loopbackQuestionCorrectAnswersIn struct{}

type loopbackExerciceValidIn struct {
	Answer EvaluateExerciceIn `dart-extern:"editor:shared_gen.dart"`
}

func (loopbackPing) isLoopbackClientEvent()                     {}
func (loopbackQuestionValidIn) isLoopbackClientEvent()          {}
func (loopbackQuestionCorrectAnswersIn) isLoopbackClientEvent() {}
func (loopbackExerciceValidIn) isLoopbackClientEvent()          {}
