package exercice

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
