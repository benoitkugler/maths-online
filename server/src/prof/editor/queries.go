package editor

import (
	"github.com/benoitkugler/maths-online/utils"
)

// SelectAllTags returns all the tags already used.
// Note that it does not include the special DifficultyTags
func SelectAllTags(db DB) ([]string, error) {
	rs, err := db.Query("SELECT DISTINCT tag FROM question_tags ORDER BY tag")
	if err != nil {
		return nil, err
	}
	var out []string
	defer func() {
		errClose := rs.Close()
		if err == nil {
			err = errClose
		}
	}()
	for rs.Next() {
		var s string
		err = rs.Scan(&s)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	if err = rs.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

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