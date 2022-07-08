package trivial

import (
	"fmt"
	"sort"

	"github.com/benoitkugler/maths-online/prof/editor"
	tv "github.com/benoitkugler/maths-online/trivial-poursuit/game"
	"github.com/benoitkugler/maths-online/utils"
)

var demoQuestions = CategoriesQuestions{
	{
		{"Pourcentages", "Valeur initiale"},
		{"Pourcentages", "Valeur finale"},
	},
	{
		{"Pourcentages", "Taux global"},
		{"Pourcentages", "Taux réciproque"},
	},
	{
		{"Pourcentages", "Proportion"},
		{"Pourcentages", "Proportion de proportion"},
		{"Pourcentages", "Pourcentage d'un nombre"},
	},
	{
		{"Pourcentages", "Evolutions identiques"},
		{"Pourcentages", "Evolution unique"},
		{"Pourcentages", "Evolutions successives"},
	},
	{
		{"Pourcentages", "Coefficient multiplicateur"},
		{"Pourcentages", "Taux d'évolution"},
	},
}

// remove empty intersection and normalizes tags
func (qc QuestionCriterion) normalize() (out QuestionCriterion) {
	for _, q := range qc {
		for i, t := range q {
			q[i] = editor.NormalizeTag(t)
		}

		if len(q) != 0 {
			out = append(out, q)
		}
	}
	return out
}

func (qc QuestionCriterion) filter(tags editor.QuestionTags) (out IDs) {
	for idQuestion, questions := range tags.ByIdQuestion() {
		questionTags := questions.Crible()
		for _, union := range qc { // at least one union must match
			if questionTags.HasAll(union) {
				out = append(out, idQuestion)
				break // no need to check the other unions
			}
		}
	}
	return out
}

// selectQuestions selects the questions matching the criterion, available
// to the user, and not needing an exercice context.
func (qc QuestionCriterion) selectQuestions(db DB, userID int64) (editor.Questions, error) {
	qc = qc.normalize()

	// an empty criterion is interpreted as an invalid criterion,
	// since it is never something you want in practice (at least the class level should be specified)
	if len(qc) == 0 {
		return nil, nil
	}

	tags, err := editor.SelectAllQuestionTags(db)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	tmp := qc.filter(tags)

	// restrict to user visible and standalone
	questionsDict, err := editor.SelectQuestions(db, tmp...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	questionsDict.RestrictVisible(userID)

	questionsDict.RestrictNeedExercice()

	return questionsDict, nil
}

func (cats *CategoriesQuestions) normalize() {
	for i := range cats {
		cats[i] = cats[i].normalize()
	}
}

func (cats CategoriesQuestions) selectQuestions(db DB, userID int64) (out tv.QuestionPool, err error) {
	// select the questions...
	for i, cat := range cats {
		questionsDict, err := cat.selectQuestions(db, userID)
		if err != nil {
			return out, err
		}

		// this should be avoided by the client side validation
		if len(questionsDict) == 0 {
			return out, fmt.Errorf("La catégorie %d n'a aucune question.", i+1)
		}

		// select the tags, required for difficulty groups
		tags, err := editor.SelectQuestionTagsByIdQuestions(db, questionsDict.IDs()...)
		if err != nil {
			return out, utils.SQLError(err)
		}
		tagsDict := tags.ByIdQuestion()

		tmp := make([]questionDiff, 0, len(questionsDict))
		questions := make([]editor.Question, 0, len(questionsDict))
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
	diff     editor.DifficultyTag
	question editor.Question
}

// weightQuestions compute the probabilities of each question in
// the given slice to account for implicit groups defined by title and difficulties
func weightQuestions(questions []questionDiff) []float64 {
	// form title groups
	groups := make(map[string][]questionDiff)
	for _, qu := range questions {
		groups[qu.question.Page.Title] = append(groups[qu.question.Page.Title], qu)
	}
	// now differentiate against the difficulty;
	// to simplify, we consider that question without difficulty form a sub-group of their own
	difficulties := make(map[string]map[editor.DifficultyTag][]questionDiff)
	for ti, group := range groups {
		perDifficulty := make(map[editor.DifficultyTag][]questionDiff)
		for _, question := range group {
			perDifficulty[question.diff] = append(perDifficulty[question.diff], question)
		}
		difficulties[ti] = perDifficulty
	}

	NG := len(groups)
	out := make([]float64, len(questions))
	for i, qu := range questions {
		perDifficulty := difficulties[qu.question.Page.Title]
		ND := len(perDifficulty) // number of difficulties in the group
		// each group must have a resulting proba of 1/NG,
		// now, each subgroup must have a resulting proba of 1/ND,
		// meaning a question if a subgroup with length nd has proba (1/NG) * (1/ND) * (1/nd)
		nd := len(perDifficulty[qu.diff])
		out[i] = 1 / float64(NG*ND*nd)
	}

	return out
}

// commonTags returns the tags shared by all categories
func (cats CategoriesQuestions) commonTags() []string {
	var allUnions [][]string
	for _, cat := range cats {
		allUnions = append(allUnions, cat...)
	}
	return editor.CommonTags(allUnions)
}

// returns the questions available to `userID` and matching one of the
// categorie criteria
func (cats CategoriesQuestions) selectQuestionIds(db DB, userID int64) (Set, error) {
	crible := NewSet()

	for _, cat := range cats {
		questions, err := cat.selectQuestions(db, userID)
		if err != nil {
			return nil, err
		}
		for id := range questions {
			crible.Add(id)
		}
	}

	return crible, nil
}
