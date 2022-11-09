package reviews

import (
	"github.com/benoitkugler/maths-online/sql/editor"
	"github.com/benoitkugler/maths-online/sql/teacher"
	"github.com/benoitkugler/maths-online/sql/trivial"
)

type IdReview int64

// Review stores the messages and evaluation about publishing a user
// created content in the admin account.
// An implicit invariant is that each [Review] is mapped to exactly one item
// from the tables [ReviewTrivial], [ReviewQuestion], [ReviewExercice].
type Review struct {
	Id   IdReview
	Kind Kind
}

// gomacro:SQL ADD FOREIGN KEY (IdReview, Kind) REFERENCES Review (ID, Kind)
// gomacro:SQL ADD CHECK (Kind = #[Kind.KQuestion])
// gomacro:SQL ADD UNIQUE (IdQuestion)
// gomacro:SQL ADD UNIQUE (IdReview)
type ReviewQuestion struct {
	IdReview   IdReview `gomacro-sql-on-delete:"CASCADE"`
	IdQuestion editor.IdQuestiongroup
	Kind       Kind // used for integrity
}

// gomacro:SQL ADD FOREIGN KEY (IdReview, Kind) REFERENCES Review (ID, Kind)
// gomacro:SQL ADD CHECK (Kind = #[Kind.KExercice])
// gomacro:SQL ADD UNIQUE (IdExercice)
// gomacro:SQL ADD UNIQUE (IdReview)
type ReviewExercice struct {
	IdReview   IdReview `gomacro-sql-on-delete:"CASCADE"`
	IdExercice editor.IdExercicegroup
	Kind       Kind // used for integrity
}

// gomacro:SQL ADD FOREIGN KEY (IdReview, Kind) REFERENCES Review (ID, Kind)
// gomacro:SQL ADD CHECK (Kind = #[Kind.KTrivial])
// gomacro:SQL ADD UNIQUE (IdTrivial)
// gomacro:SQL ADD UNIQUE (IdReview)
type ReviewTrivial struct {
	IdReview  IdReview `gomacro-sql-on-delete:"CASCADE"`
	IdTrivial trivial.IdTrivial
	Kind      Kind // used for integrity
}

type Comments []Comment

// gomacro:SQL ADD UNIQUE (IdReview, IdTeacher)
type ReviewParticipation struct {
	IdReview  IdReview          `gomacro-sql-on-delete:"CASCADE"`
	IdTeacher teacher.IdTeacher `gomacro-sql-on-delete:"CASCADE"`
	Approval  Approval
	Comments  Comments
}
