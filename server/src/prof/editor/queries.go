package editor

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

// SelectQuestionByTags returns the question matching ALL the tags given
// It panics if tags is empty.
func SelectQuestionByTags(db DB, tags ...string) (map[int64]QuestionTags, error) {
	firstSelection, err := selectQuestionByTag(db, tags[0])
	if err != nil {
		return nil, err
	}

	quTags, err := SelectQuestionTagsByIdQuestions(db, firstSelection.IDs()...)
	if err != nil {
		return nil, err
	}

	dict := quTags.ByIdQuestion()

	// remove questions not matching all the tags
	for idQuestion, cr := range dict {
		hasAll := cr.Crible().HasAll(tags)
		if !hasAll {
			delete(dict, idQuestion)
		}
	}

	return dict, nil
}
