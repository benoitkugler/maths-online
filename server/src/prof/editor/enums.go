package editor

const (
	Ping          loopbackClientDataKind = iota // nil
	CheckSyntaxIn                               // exercice/client.QuestionSyntaxCheckIn
	ValidAnswerIn                               // exercice/client.QuestionAnswersIn
)

const (
	State           loopbackServerDataKind = iota // LoopbackState
	CheckSyntaxeOut                               // exercice/client.QuestionSyntaxCheckOut
	ValidAnswerOut                                // exercice/client.QuestionAnswersOut
)
