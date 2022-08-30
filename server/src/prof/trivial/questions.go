package trivial

import (
	"sort"

	ed "github.com/benoitkugler/maths-online/sql/editor"
	tc "github.com/benoitkugler/maths-online/sql/trivial"
	tv "github.com/benoitkugler/maths-online/trivial"
	"github.com/benoitkugler/maths-online/utils"
)

var demoQuestions = tc.CategoriesQuestions{
	Tags: [...]tc.QuestionCriterion{
		{
			{"POURCENTAGES", "EVOLUTION UNIQUE"},
		},
		{
			{"POURCENTAGES", "TAUX GLOBAL"},
			{"POURCENTAGES", "TAUX RÉCIPROQUE"},
		},
		{
			{"POURCENTAGES", "PROPORTION"},
			{"POURCENTAGES", "PROPORTION DE PROPORTION"},
			{"POURCENTAGES", "POURCENTAGE D'UN NOMBRE"},
		},
		{
			{"POURCENTAGES", "EVOLUTIONS IDENTIQUES"},
			{"POURCENTAGES", "EVOLUTIONS SUCCESSIVES"},
		},
		{
			{"POURCENTAGES", "COEFFICIENT MULTIPLICATEUR"},
			{"POURCENTAGES", "TAUX D'EVOLUTION"},
		},
	},
}

func filterTags(qc tc.QuestionCriterion, tags ed.QuestiongroupTags) (out []ed.IdQuestiongroup) {
	for idGroup, groupTags := range tags.ByIdQuestiongroup() {
		questionTags := groupTags.Crible()
		for _, intersection := range qc { // at least one intersection must match
			if questionTags.HasAll(intersection) {
				out = append(out, idGroup)
				break // no need to check the other unions
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

type questionSelector struct {
	db ed.DB // used to lazily load question content

	tags           ed.QuestiongroupTags // all the DB
	questionGroups ed.Questiongroups    // all the DB
	questions      ed.Questions         // all the ones coming from groups
}

// load once for all all the tags and questionGroups
func newQuestionSelector(db ed.DB) (out questionSelector, err error) {
	out.db = db

	out.tags, err = ed.SelectAllQuestiongroupTags(db)
	if err != nil {
		return out, utils.SQLError(err)
	}

	out.questionGroups, err = ed.SelectAllQuestiongroups(db)
	if err != nil {
		return out, utils.SQLError(err)
	}

	out.questions, err = ed.SelectQuestionsByIdGroups(db, out.questionGroups.IDs()...)
	if err != nil {
		return out, utils.SQLError(err)
	}

	return out, nil
}

// do not error on empty catagories
func (sel questionSelector) search(query tc.CategoriesQuestions, userID uID) (out tv.QuestionPool, err error) {
	query.Normalize()

	// select the questions...
	for i, cat := range query.Tags {
		// an empty criterion is interpreted as an never matched criterion,
		// since it is never something you want in practice (at least the class level should be specified)
		if len(cat) == 0 {
			return out, nil
		}

		idGroups := ed.NewIdQuestiongroupSetFrom(filterTags(cat, sel.tags))

		groups := make(ed.Questiongroups, len(idGroups)) // select the groups
		for id := range idGroups {
			groups[id] = sel.questionGroups[id]
		}

		// restrict to user visible
		groups.RestrictVisible(userID)

		questionsDict := make(ed.Questions) // select the questions
		for _, qu := range sel.questions {
			if idGroups.Has(qu.IdGroup.ID) {
				questionsDict[qu.Id] = qu
			}
		}

		// filter by difficulty
		for _, question := range questionsDict {
			if !query.Difficulties.Match(question.Difficulty) {
				delete(questionsDict, question.Id)
			}
		}

		wQuestions := weightQuestions(questionsDict)

		out[i] = wQuestions
	}

	return out, nil
}

func selectQuestions(db ed.DB, query tc.CategoriesQuestions, userID uID) (out tv.QuestionPool, err error) {
	sel, err := newQuestionSelector(db)
	if err != nil {
		return out, err
	}
	return sel.search(query, userID)
}

type sorter tv.WeigthedQuestions

func (wq sorter) Len() int { return len(wq.Questions) }
func (wq sorter) Swap(i, j int) {
	wq.Questions[i], wq.Questions[j] = wq.Questions[j], wq.Questions[i]
	wq.Weights[i], wq.Weights[j] = wq.Weights[j], wq.Weights[i]
}
func (wq sorter) Less(i, j int) bool { return wq.Questions[i].Id < wq.Questions[j].Id }

// weightQuestions compute the probabilities of each question in
// the given set to account for groups and difficulties
func weightQuestions(questions ed.Questions) tv.WeigthedQuestions {
	// form groups
	groups := questions.ByGroup()

	// now differentiate against the difficulty;
	// to simplify, we consider that question without difficulty form a common sub-group of their own
	difficulties := make(map[ed.IdQuestiongroup]map[ed.DifficultyTag][]ed.Question)
	for idGroup, group := range groups {
		perDifficulty := make(map[ed.DifficultyTag][]ed.Question)
		for _, question := range group {
			perDifficulty[question.Difficulty] = append(perDifficulty[question.Difficulty], question)
		}
		difficulties[idGroup] = perDifficulty
	}

	NG := len(groups)
	out := tv.WeigthedQuestions{
		Questions: make([]ed.Question, 0, len(questions)),
		Weights:   make([]float64, 0, len(questions)),
	}
	for _, qu := range questions {
		perDifficulty := difficulties[qu.IdGroup.ID]
		ND := len(perDifficulty) // number of difficulties in the group
		// each group must have a resulting proba of 1/NG,
		// now, each subgroup must have a resulting proba of 1/ND,
		// meaning a question in a sub-group with length nd has proba (1/NG) * (1/ND) * (1/nd)
		nd := len(perDifficulty[qu.Difficulty])

		out.Questions = append(out.Questions, qu)
		out.Weights = append(out.Weights, 1/float64(NG*ND*nd))
	}

	// sorting is not required, but make tests easier to write
	sort.Sort(sorter(out))

	return out
}

// commonTags returns the tags shared by all categories
func commonTags(cats tc.CategoriesQuestions) []string {
	var allUnions [][]string
	for _, cat := range cats.Tags {
		allUnions = append(allUnions, cat...)
	}
	return ed.CommonTags(allUnions)
}

// returns the union of all the question groups in the pool,
// that is, question matching at least one criteria
func allQuestions(pool tv.QuestionPool) ed.IdQuestiongroupSet {
	crible := make(ed.IdQuestiongroupSet)

	for _, cat := range pool {
		for _, question := range cat.Questions {
			crible.Add(question.IdGroup.ID)
		}
	}

	return crible
}
