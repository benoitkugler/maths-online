package reviews

import (
	"time"

	"github.com/benoitkugler/maths-online/utils"
)

// Comment stores the content of a comment in a
// review
type Comment struct {
	Time    time.Time
	Message string
}

// Approval is the evaluation of one teacher
// about the review
type Approval uint8

const (
	// The teacher does not pronounce itself
	Neutral Approval = iota
	InFavor
	Opposed
)

// Kind is an enum describing the kind of item
// which may be in a review
type Kind uint8

const (
	KQuestion Kind = iota
	KExercice
	KTrivial
)

type Target interface {
	Review() IdReview
}

func (i ReviewTrivial) Review() IdReview  { return i.IdReview }
func (i ReviewQuestion) Review() IdReview { return i.IdReview }
func (i ReviewExercice) Review() IdReview { return i.IdReview }

// LoadTargets load all the targets associated to the reviews
func LoadTargets(db DB) (map[IdReview]Target, error) {
	questions, err := SelectAllReviewQuestions(db)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	exercices, err := SelectAllReviewExercices(db)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	trivials, err := SelectAllReviewTrivials(db)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	out := make(map[IdReview]Target)
	for _, target := range questions {
		out[target.IdReview] = target
	}
	for _, target := range exercices {
		out[target.IdReview] = target
	}
	for _, target := range trivials {
		out[target.IdReview] = target
	}
	return out, nil
}
