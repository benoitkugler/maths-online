package trivial

import (
	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	"github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/tasks"

	tv "github.com/benoitkugler/maths-online/server/src/trivial"
)

type QuestionContent struct {
	Id        editor.IdQuestion `gomacro-opaque:"dart"` // 0 if empty
	Categorie tv.Categorie

	// instance, to be used by the embeded preview

	Question client.Question `gomacro-opaque:"typescript"`
	Params   tasks.Params    `gomacro-opaque:"typescript"`
}
