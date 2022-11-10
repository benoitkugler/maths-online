package reviews

import (
	"database/sql"
	"errors"
	"time"

	"github.com/benoitkugler/maths-online/sql/editor"
	"github.com/benoitkugler/maths-online/sql/teacher"
	"github.com/benoitkugler/maths-online/sql/trivial"
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

type TargetHeader struct {
	Title string
	Owner teacher.IdTeacher
}

type Target interface {
	Review() IdReview
	WithIdReview(IdReview) Target

	// Insert inserts the target in the proper table
	Insert(tx *sql.Tx) error

	// Errors should not be wrapped
	Load(db DB) (TargetHeader, error)
}

func (tr ReviewQuestion) Review() IdReview { return tr.IdReview }
func (tr ReviewExercice) Review() IdReview { return tr.IdReview }
func (tr ReviewTrivial) Review() IdReview  { return tr.IdReview }

func (tr ReviewQuestion) WithIdReview(r IdReview) Target { tr.IdReview = r; return tr }
func (tr ReviewExercice) WithIdReview(r IdReview) Target { tr.IdReview = r; return tr }
func (tr ReviewTrivial) WithIdReview(r IdReview) Target  { tr.IdReview = r; return tr }

func (tr ReviewQuestion) Insert(tx *sql.Tx) error { return InsertManyReviewQuestions(tx, tr) }
func (tr ReviewExercice) Insert(tx *sql.Tx) error { return InsertManyReviewExercices(tx, tr) }
func (tr ReviewTrivial) Insert(tx *sql.Tx) error  { return InsertManyReviewTrivials(tx, tr) }

func (tr ReviewQuestion) Load(db DB) (TargetHeader, error) {
	item, err := editor.SelectQuestiongroup(db, tr.IdQuestion)
	if err != nil {
		return TargetHeader{}, err
	}
	return TargetHeader{Title: item.Title, Owner: item.IdTeacher}, nil
}

func (tr ReviewExercice) Load(db DB) (TargetHeader, error) {
	item, err := editor.SelectExercicegroup(db, tr.IdExercice)
	if err != nil {
		return TargetHeader{}, err
	}
	return TargetHeader{Title: item.Title, Owner: item.IdTeacher}, nil
}

func (tr ReviewTrivial) Load(db DB) (TargetHeader, error) {
	item, err := trivial.SelectTrivial(db, tr.IdTrivial)
	if err != nil {
		return TargetHeader{}, err
	}
	return TargetHeader{Title: item.Name, Owner: item.IdTeacher}, nil
}

// LoadTarget load all the target associated to the given review
func LoadTarget(db DB, id IdReview) (Target, error) {
	question, isQuestion, err := SelectReviewQuestionByIdReview(db, id)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	exercice, isExercice, err := SelectReviewExerciceByIdReview(db, id)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	trivial, isTrivial, err := SelectReviewTrivialByIdReview(db, id)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	switch {
	case isQuestion:
		return question, nil
	case isExercice:
		return exercice, nil
	case isTrivial:
		return trivial, nil
	default:
		return nil, errors.New("internal error: review without target")
	}
}

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

// UpdateParticipation update the fields given by [part]
func UpdateParticipation(db *sql.DB, part ReviewParticipation) error {
	// insert back into the DB : delete and insert
	tx, err := db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}
	err = part.Delete(tx)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}
	err = InsertManyReviewParticipations(tx, part)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}
	err = tx.Commit()
	if err != nil {
		return utils.SQLError(err)
	}
	return nil
}
