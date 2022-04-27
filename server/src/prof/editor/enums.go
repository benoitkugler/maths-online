package editor

const (
	Ping          loopbackClientDataKind = iota // nil
	ValidAnswerIn                               // exercice/client.QuestionAnswersIn
)

const (
	State          loopbackServerDataKind = iota // LoopbackState
	ValidAnswerOut                               // exercice/client.QuestionAnswersOut
)
