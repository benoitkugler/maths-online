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

const (
	Diff1 DifficultyTag = "\u2605"             // 1 étoile
	Diff2 DifficultyTag = "\u2605\u2605"       // 2 étoiles
	Diff3 DifficultyTag = "\u2605\u2605\u2605" // 3 étoiles
)

type Flow uint8

const (
	Parallel   Flow = iota // Questions indépendantes
	Sequencial             // Questions liées
)
