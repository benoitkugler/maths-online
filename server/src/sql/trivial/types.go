package trivial

import (
	"github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/trivial"
)

// QuestionCriterion is an union of intersection of tags.
type QuestionCriterion [][]editor.TagSection

// normalize removes empty intersections and normalizes tags
func (qc QuestionCriterion) normalize() (out QuestionCriterion) {
	for _, q := range qc {
		for i, t := range q {
			q[i].Tag = editor.NormalizeTag(t.Tag)
		}

		if len(q) != 0 {
			out = append(out, q)
		}
	}
	return out
}

// CategoriesQuestions defines a union of intersection of tags,
// for every category.
type CategoriesQuestions struct {
	Tags [trivial.NbCategories]QuestionCriterion
	// Union. An empty slice means no selection : all variants are accepted.
	Difficulties editor.DifficultyQuery
}

// Normalize removes empty intersections and normalizes tags, for each
// categories
func (query *CategoriesQuestions) Normalize() {
	for i := range query.Tags {
		query.Tags[i] = query.Tags[i].normalize()
	}
}
