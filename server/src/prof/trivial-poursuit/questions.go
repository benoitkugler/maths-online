package trivialpoursuit

import (
	"fmt"
	"sort"

	"github.com/benoitkugler/maths-online/maths/exercice"
	tv "github.com/benoitkugler/maths-online/trivial-poursuit/game"
	"github.com/benoitkugler/maths-online/utils"
)

func (cats CategoriesQuestions) selectQuestions(db DB) (out tv.QuestionPool, err error) {
	// select the questions...
	tags, err := exercice.SelectAllQuestionTags(db)
	if err != nil {
		return out, utils.SQLError(err)
	}

	tagsDict := tags.ByIdQuestion()
	for i, cat := range cats {
		idQuestions := cat.filter(tagsDict)
		// this should be avoided by the client side validation
		if len(idQuestions) == 0 {
			return out, fmt.Errorf("La catégorie %d n'a aucune question.", i+1)
		}

		questionsDict, err := exercice.SelectQuestions(db, idQuestions...)
		if err != nil {
			return out, utils.SQLError(err)
		}

		// select the tags, required for difficulty groups
		tags, err := exercice.SelectQuestionTagsByIdQuestions(db, idQuestions...)
		if err != nil {
			return out, utils.SQLError(err)
		}
		tagsDict := tags.ByIdQuestion()

		tmp := make([]questionDiff, 0, len(questionsDict))
		questions := make([]exercice.Question, 0, len(questionsDict))
		for _, question := range questionsDict {
			diff := tagsDict[question.Id].Crible().Difficulty()
			tmp = append(tmp, questionDiff{question: question, diff: diff})
			questions = append(questions, question)
		}

		// sorting is not required, but make tests easier to write
		sort.Slice(tmp, func(i, j int) bool { return tmp[i].question.Id < tmp[j].question.Id })
		sort.Slice(questions, func(i, j int) bool { return questions[i].Id < questions[j].Id })

		weights := weightQuestions(tmp)
		out[i] = tv.WeigthedQuestions{
			Questions: questions,
			Weights:   weights,
		}
	}

	return out, nil
}

type questionDiff struct {
	diff     exercice.DifficultyTag
	question exercice.Question
}

// weightQuestions compute the probabilities of each question in
// the given slice to account for implicit groups defined by title and difficulties
func weightQuestions(questions []questionDiff) []float64 {
	// form title groups
	groups := make(map[string][]questionDiff)
	for _, qu := range questions {
		groups[qu.question.Title] = append(groups[qu.question.Title], qu)
	}
	// now differentiate against the difficulty;
	// we consider that question without difficulty form a sub-group of their own
	difficulties := make(map[string]map[exercice.DifficultyTag][]questionDiff)
	for ti, group := range groups {
		perDifficulty := make(map[exercice.DifficultyTag][]questionDiff)
		for _, question := range group {
			perDifficulty[question.diff] = append(perDifficulty[question.diff], question)
		}
		difficulties[ti] = perDifficulty
	}

	NG := len(groups)
	out := make([]float64, len(questions))
	for i, qu := range questions {
		perDifficulty := difficulties[qu.question.Title]
		ND := len(perDifficulty) // number of difficulties in the group
		// each group must have a resulting proba of 1/NG,
		// now, each subgroup must have a resulting proba of 1/ND,
		// meaning a question if a subgroup with length nd has proba (1/NG) * (1/ND) * (1/nd)
		nd := len(perDifficulty[qu.diff])
		out[i] = 1 / float64(NG*ND*nd)
	}

	return out
}