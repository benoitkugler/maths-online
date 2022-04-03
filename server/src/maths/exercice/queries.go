package exercice

import "encoding/json"

func (ct Enonce) MarshalJSON() ([]byte, error) {
	tmp := make([]BlockWrapper, len(ct))
	for i, v := range ct {
		tmp[i].Data = v
	}
	return json.Marshal(tmp)
}

func (ct *Enonce) UnmarshalJSON(data []byte) error {
	var tmp []BlockWrapper
	err := json.Unmarshal(data, &tmp)
	*ct = make(Enonce, len(tmp))
	for i, v := range tmp {
		(*ct)[i] = v.Data
	}
	return err
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
