package editor

import (
	"database/sql"

	"github.com/benoitkugler/maths-online/utils"
)

func newID(ID int64) sql.NullInt64 { return sql.NullInt64{Valid: true, Int64: ID} }

// selectQuestionByTag returns the question matching the given tag
func selectQuestionByTag(db DB, tag string) (Questions, error) {
	rs, err := db.Query(`SELECT questions.*
	FROM questions
	JOIN question_tags ON questions.id = question_tags.id_question
   	WHERE question_tags.tag = $1`, tag)
	if err != nil {
		return nil, err
	}
	return ScanQuestions(rs)
}

// IsVisibleBy returns `true` if the question is public or
// owned by `userID`
func (qu Question) IsVisibleBy(userID int64) bool {
	return qu.Public || qu.IdTeacher == userID
}

// RestrictVisible remove the questions not visible by `userID`
func (qus Questions) RestrictVisible(userID int64) {
	for id, qu := range qus {
		if !qu.IsVisibleBy(userID) {
			delete(qus, id)
		}
	}
}

// RestrictNeedExercice remove the questions marked as requiring an
// exercice
func (qus Questions) RestrictNeedExercice() {
	for id, question := range qus {
		if question.NeedExercice.Valid {
			delete(qus, id)
		}
	}
}

// SelectQuestionByTags returns the question matching ALL the tags given,
// and available to `userID`, returning a map IdQuestion -> Tags
// It panics if tags is empty.
func SelectQuestionByTags(db DB, userID int64, tags ...string) (map[int64]QuestionTags, error) {
	firstSelection, err := selectQuestionByTag(db, tags[0])
	if err != nil {
		return nil, err
	}

	quTags, err := SelectQuestionTagsByIdQuestions(db, firstSelection.IDs()...)
	if err != nil {
		return nil, err
	}

	dict := quTags.ByIdQuestion()
	var selectedIDs IDs
	// remove questions not matching all the tags
	for idQuestion, cr := range dict {
		hasAll := cr.Crible().HasAll(tags)
		if !hasAll {
			delete(dict, idQuestion)
		} else {
			selectedIDs = append(selectedIDs, idQuestion)
		}
	}

	// restrict to available questions
	questions, err := SelectQuestions(db, selectedIDs...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	for _, qu := range questions {
		if !qu.IsVisibleBy(userID) {
			delete(dict, qu.Id)
		}
	}

	return dict, nil
}

// updateExerciceQuestionList set the questions for the given exercice,
// overiding `IdExercice` and `index` fields of the list items.
func updateExerciceQuestionList(db *sql.DB, idExercice int64, l ExerciceQuestions) error {
	// enforce fields
	for i := range l {
		l[i].Index = i
		l[i].IdExercice = idExercice
	}
	tx, err := db.Begin()
	if err != nil {
		return utils.SQLError(err)
	}
	_, err = DeleteExerciceQuestionsByIdExercices(db, idExercice)
	if err != nil {
		_ = tx.Rollback()
		return utils.SQLError(err)
	}

	err = InsertManyExerciceQuestions(tx, l...)
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

// IsVisibleBy returns `true` if the question is public or
// owned by `userID`
func (qu Exercice) IsVisibleBy(userID int64) bool {
	return qu.Public || qu.IdTeacher == userID
}
