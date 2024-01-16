package ceintures

import (
	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	"github.com/benoitkugler/maths-online/server/src/pass"
	ce "github.com/benoitkugler/maths-online/server/src/sql/ceintures"
	"github.com/benoitkugler/maths-online/server/src/tasks"
)

// used to autogenerate Dart types

type CreateEvolutionIn struct {
	ClientID pass.EncryptedID // empty for anonymous users
	Level    ce.Level
}

type CreateEvolutionOut struct {
	AnonymousID string // optional
	Evolution   StudentEvolution
}

type GetEvolutionOut struct {
	Has       bool
	Evolution StudentEvolution
}

type StudentEvolution struct {
	Scheme          Scheme
	Level           ce.Level
	Advance         ce.Advance
	Stats           ce.Stats
	Pending         []Stage
	SuggestionIndex int // index in Pending, or -1 if Pending is empty
}

func newStudentEvolution(ev ce.Beltevolution) StudentEvolution {
	pending := mathScheme.Pending(ev.Advance, ev.Level)
	return StudentEvolution{
		Scheme:          mathScheme,
		Level:           ev.Level,
		Advance:         ev.Advance,
		Stats:           ev.Stats,
		Pending:         pending,
		SuggestionIndex: mathScheme.suggestionIndex(pending),
	}
}

type SelectQuestionsIn struct {
	Tokens StudentTokens
	Stage  Stage
}

type InstantiatedBeltQuestion struct {
	Id       ce.IdBeltquestion
	Question client.Question
	Params   tasks.Params // for the evaluation
}

type SelectQuestionsOut struct {
	Questions []InstantiatedBeltQuestion // Id and Difficulty should be ignored
}

type EvaluateAnswersIn struct {
	Tokens    StudentTokens
	Stage     Stage
	Questions []ce.IdBeltquestion // as send by the server; the Question field is ignored
	Answers   []tasks.AnswerP
}

type EvaluateAnswersOut struct {
	Answers   []client.QuestionAnswersOut
	Evolution StudentEvolution // updated evolution
}
