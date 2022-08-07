package trivial

import "github.com/benoitkugler/maths-online/prof/teacher"

type IdTrivial int64

// Trivial is a trivial game configuration
// stored in the DB, one per activity.
type Trivial struct {
	Id              IdTrivial
	Questions       CategoriesQuestions
	QuestionTimeout int // in seconds
	ShowDecrassage  bool
	Public          bool
	IdTeacher       teacher.IdTeacher `json:"id_teacher"`
	Name            string
}
