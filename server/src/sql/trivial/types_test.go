package trivial

import (
	"testing"

	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func TestMatchMatiere(t *testing.T) {
	tu.Assert(t, !CategoriesQuestions{}.MatchMatiere(teacher.Mathematiques))
	tu.Assert(t, CategoriesQuestions{}.MatchMatiere(teacher.Autre))
}
