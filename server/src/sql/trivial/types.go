package trivial

import (
	"github.com/benoitkugler/maths-online/sql/editor"
	"github.com/benoitkugler/maths-online/trivial"
)

// QuestionCriterion is an union of intersection of tags.
type QuestionCriterion [][]string

// normalize removes empty intersections and normalizes tags
func (qc QuestionCriterion) normalize() (out QuestionCriterion) {
	for _, q := range qc {
		for i, t := range q {
			q[i] = editor.NormalizeTag(t)
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
