package trivial

import (
	"reflect"
	"testing"

	"github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func TestMatchMatiere(t *testing.T) {
	tu.Assert(t, !CategoriesQuestions{}.MatchMatiere(teacher.Mathematiques))
	tu.Assert(t, CategoriesQuestions{}.MatchMatiere(teacher.Autre))
}

func TestCommon(t *testing.T) {
	cq := CategoriesQuestions{Tags: [5]QuestionCriterion{
		{
			{{Section: editor.Chapter, Tag: "AUTO"}, {Section: editor.Level, Tag: "SNDE"}},
			{{Section: editor.Chapter, Tag: "AUTO"}, {Section: editor.Level, Tag: "1ERE"}},
		},
		{{{Section: editor.Chapter, Tag: "AUTO"}, {Section: editor.Level, Tag: "SNDE"}}},
		{{{Section: editor.Chapter, Tag: "AUTO"}, {Section: editor.Level, Tag: "SNDE"}}},
		{{{Section: editor.Chapter, Tag: "AUTO"}, {Section: editor.Level, Tag: "SNDE"}}},
		{{{Section: editor.Chapter, Tag: "AUTO"}, {Section: editor.Level, Tag: "SNDE"}}},
	}}
	tu.Assert(t, reflect.DeepEqual(cq.Common(), editor.Tags{{Section: editor.Chapter, Tag: "AUTO"}}))
}
