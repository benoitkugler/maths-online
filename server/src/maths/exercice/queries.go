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

// SelectQuestionByTags returns the question matching ALL the tags given
// It panis if tags is empty.
func SelectQuestionByTags(db DB, queryTags ...string) (Questions, error) {
	firstSelection, err := selectQuestionByTag(db, queryTags[0])
	if err != nil {
		return nil, err
	}

	quTags, err := SelectQuestionTagsByIdQuestions(db, firstSelection.IDs()...)
	if err != nil {
		return nil, err
	}

	var ids IDs
	for idQuestion, thisTags := range quTags.ByIdQuestion() {
		crible := thisTags.Crible()

		if crible.HasAll(queryTags) {
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

// Crible is a set of tags.
type Crible map[string]bool

// HasAll returns `true` is all the `tags` are present in the crible.
func (cr Crible) HasAll(tags []string) bool {
	for _, tag := range tags {
		if !cr[tag] {
			return false
		}
	}
	return true
}

// Crible build a set from the tags
func (qus QuestionTags) Crible() Crible {
	out := make(Crible, len(qus))
	for _, qt := range qus {
		out[qt.Tag] = true
	}
	return out
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
