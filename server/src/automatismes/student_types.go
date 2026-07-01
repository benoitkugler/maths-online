package automatismes

import (
	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	"github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/sql/trivial"
	"github.com/benoitkugler/maths-online/server/src/tasks"
)

type GetAutomatismesIn struct {
	Level    editor.LevelTag
	Sublevel string
}

type GetAutomatismesOut struct {
	Trivials  []trivial.Trivial
	Questions []editor.Questiongroup
}

// Automatisme variant

type InstantiatedAutomatismeQuestion struct {
	Id       editor.IdQuestion
	Question client.Question
	Params   tasks.Params // for the evaluation
}

type EvaluateAutomatismeIn struct {
	Id     editor.IdQuestion
	Answer tasks.AnswerP
}
