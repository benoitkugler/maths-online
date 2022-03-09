package exercice

import "encoding/json"

func (ct Content) MarshalJSON() ([]byte, error) {
	tmp := make([]blockWrapper, len(ct))
	for i, v := range ct {
		tmp[i].Data = v
	}
	return json.Marshal(tmp)
}

func (ct *Content) UnmarshalJSON(data []byte) error {
	var tmp []blockWrapper
	err := json.Unmarshal(data, &tmp)
	*ct = make(Content, len(tmp))
	for i, v := range tmp {
		(*ct)[i] = v.Data
	}
	return err
}

type ExerciceQuestions struct {
	Exercice
	Questions Questions
}

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
