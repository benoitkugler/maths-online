package exercice

func SelectExerciceQuestions(db DB, id int64) (ExerciceQuestions, error) {
	ex, err := SelectExercice(db, id)
	if err != nil {
		return ExerciceQuestions{}, err
	}

	qu, err := SelectQuestionsByIdExercices(db, id)
	if err != nil {
		return ExerciceQuestions{}, err
	}

	return ExerciceQuestions{Exercice: ex, Questions: qu}, nil
}
