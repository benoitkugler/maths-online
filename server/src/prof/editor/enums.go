package editor

const (
	State                loopbackServerDataKind = iota // LoopbackState
	ValidAnswerOut                                     // exercice/client.QuestionAnswersOut
	ShowCorrectAnswerOut                               // exercice/client.QuestionAnswersIn
)

const (
	Ping                loopbackClientDataKind = iota // nil
	ValidAnswerIn                                     // exercice/client.QuestionAnswersIn
	ShowCorrectAnswerIn                               // nil
)
