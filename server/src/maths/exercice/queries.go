package exercice

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

// SelectQuestionByTags returns the question matching ALL the tags given
// It panis if tags is empty.
func SelectQuestionByTags(db DB, tags ...string) (Questions, error) {
	firstSelection, err := selectQuestionByTag(db, tags[0])
	if err != nil {
		return nil, err
	}

	quTags, err := SelectQuestionTagsByIdQuestions(db, firstSelection.IDs()...)
	if err != nil {
		return nil, err
	}

	crible := make(map[int64]map[string]bool)
	for _, tag := range quTags {
		m := crible[tag.IdQuestion]
		if m == nil {
			m = make(map[string]bool)
		}
		m[tag.Tag] = true
		crible[tag.IdQuestion] = m
	}

	var ids IDs
	for idQuestion, m := range crible {
		hasAll := true
		for _, tag := range tags {
			if !m[tag] {
				hasAll = false
				break
			}
		}
		if hasAll {
			ids = append(ids, idQuestion)
		}
	}

	return SelectQuestions(db, ids...)
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

// func SelectExerciceQuestions(db DB, id int64) (ExerciceQuestions, error) {
// 	ex, err := SelectExercice(db, id)
// 	if err != nil {
// 		return ExerciceQuestions{}, err
// 	}

// 	qu, err := SelectQuestionsByIdExercices(db, id)
// 	if err != nil {
// 		return ExerciceQuestions{}, err
// 	}

// 	return ExerciceQuestions{Exercice: ex, Questions: qu}, nil
// }
