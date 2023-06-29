package trivial

import "github.com/benoitkugler/maths-online/server/src/sql/teacher"

type IdTrivial int64

// Trivial is a trivial game configuration
// stored in the DB, one per activity.
type Trivial struct {
	Id              IdTrivial
	Questions       CategoriesQuestions
	QuestionTimeout int // in seconds
	ShowDecrassage  bool
	Public          bool
	IdTeacher       teacher.IdTeacher
	Name            string
}

// SelfaccessTrivial is a link table enabling a teacher
// to publish (or hide) a [Trivial] for the students of a
// classroom.
// gomacro:SQL ADD FOREIGN KEY (IdClassroom, IdTeacher) REFERENCES Classrooms (Id, IdTeacher) ON DELETE CASCADE
// gomacro:SQL _SELECT KEY (IdTrivial, IdTeacher)
type SelfaccessTrivial struct {
	IdClassroom teacher.IdClassroom `gomacro-sql-on-delete:"CASCADE"`
	IdTrivial   IdTrivial           `gomacro-sql-on-delete:"CASCADE"`
	IdTeacher   teacher.IdTeacher
}
