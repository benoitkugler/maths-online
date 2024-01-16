package ceintures

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	"github.com/benoitkugler/maths-online/server/src/pass"
	ce "github.com/benoitkugler/maths-online/server/src/sql/ceintures"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/benoitkugler/maths-online/server/src/tasks"
	"github.com/benoitkugler/maths-online/server/src/utils"
	"github.com/labstack/echo/v4"
)

type StudentTokens struct {
	AnonymousID string           // may be empty
	ClientID    pass.EncryptedID // may be empty
}

type anonymousEvolutions struct {
	m    map[string]ce.Beltevolution
	lock sync.Mutex
}

func newAnonymousEvolutions() anonymousEvolutions {
	return anonymousEvolutions{m: make(map[string]ce.Beltevolution)}
}

// locks and adds a new anonymous student and returns the generated ID
func (ae *anonymousEvolutions) add(level ce.Level) string {
	ae.lock.Lock()
	defer ae.lock.Unlock()

	id := utils.RandomID(false, 32, func(s string) bool { _, has := ae.m[s]; return has })
	ae.m[id] = ce.Beltevolution{Level: level}

	return id
}

// locks
func (ae *anonymousEvolutions) get(id string) (ce.Beltevolution, bool) {
	ae.lock.Lock()
	defer ae.lock.Unlock()

	v, ok := ae.m[id]
	return v, ok
}

// locks
func (ae *anonymousEvolutions) set(id string, ev ce.Beltevolution) {
	ae.lock.Lock()
	defer ae.lock.Unlock()

	ae.m[id] = ev
}

func (ct *Controller) CeinturesCreateEvolution(c echo.Context) error {
	var args CreateEvolutionIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.createEvolution(args)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) createEvolution(args CreateEvolutionIn) (out CreateEvolutionOut, err error) {
	var tokens StudentTokens
	if ci := args.ClientID; ci != "" {
		id_, err := ct.studentKey.DecryptID(ci)
		if err != nil {
			return out, fmt.Errorf("Erreur interne: %s", err)
		}
		idStudent := teacher.IdStudent(id_)
		err = ce.InsertBeltevolution(ct.db, ce.Beltevolution{IdStudent: idStudent, Level: args.Level})
		if err != nil {
			return out, utils.SQLError(err)
		}
		tokens.ClientID = ci
	} else {
		id := ct.anons.add(args.Level)
		out.AnonymousID = id
		tokens.AnonymousID = id
	}

	out.Evolution, _, err = ct.getEvolution(tokens)
	if err != nil {
		return out, err
	}

	return out, nil
}

func (ct *Controller) CeinturesGetEvolution(c echo.Context) error {
	var args StudentTokens
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, has, err := ct.getEvolution(args)
	if err != nil {
		return err
	}

	return c.JSON(200, GetEvolutionOut{has, out})
}

func (ct *Controller) getEvolution(args StudentTokens) (out StudentEvolution, has bool, err error) {
	var ev ce.Beltevolution
	if ci := args.ClientID; ci != "" {
		id_, err := ct.studentKey.DecryptID(ci)
		if err != nil {
			return out, has, fmt.Errorf("Erreur interne: %s", err)
		}
		idStudent := teacher.IdStudent(id_)

		ev, has, err = ce.SelectBeltevolutionByIdStudent(ct.db, idStudent)
		if err != nil {
			return out, has, utils.SQLError(err)
		}

	} else {
		ev, has = ct.anons.get(args.AnonymousID)
	}

	if !has {
		return out, false, nil
	}

	return newStudentEvolution(ev), true, nil
}

func (ct *Controller) CeinturesSelectQuestions(c echo.Context) error {
	var args SelectQuestionsIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.selectQuestions(args)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) selectQuestions(args SelectQuestionsIn) (SelectQuestionsOut, error) {
	// We could check that the stage is actually reachable by the student,
	// but we "trust" the client for now
	questions, err := ce.SelectAllBeltquestions(ct.db)
	if err != nil {
		return SelectQuestionsOut{}, utils.SQLError(err)
	}
	byStage := map[Stage][]ce.Beltquestion{}
	for _, qu := range questions {
		key := Stage{qu.Domain, qu.Rank}
		byStage[key] = append(byStage[key], qu)
	}

	// we include every question, with repetitions if there is less
	// than 3
	selected := byStage[args.Stage]
	switch len(selected) {
	case 0:
		return SelectQuestionsOut{}, fmt.Errorf("Erreur interne: question manquante !")
	case 1:
		qu := selected[0]
		selected = append(selected, qu, qu)
	case 2:
		qu := selected[rand.Intn(2)]
		selected = append(selected, qu)
	}

	// we include one question from the previous rank
	if rank := args.Stage.Rank; rank >= ce.Jaune {
		previous := byStage[Stage{args.Stage.Domain, rank - 1}]
		if len(previous) == 0 {
			return SelectQuestionsOut{}, fmt.Errorf("Erreur interne: question manquante !")
		}
		qu := previous[rand.Intn(len(previous))]
		selected = append(selected, qu)
	}

	out := SelectQuestionsOut{Questions: make([]InstantiatedBeltQuestion, len(selected))}
	for i, qu := range selected {
		inst, params, err := qu.Page().InstantiateErr()
		if err != nil {
			return SelectQuestionsOut{}, err
		}
		out.Questions[i] = InstantiatedBeltQuestion{
			Id:       qu.Id,
			Question: inst.ToClient(),
			Params:   tasks.NewParams(params),
		}
	}

	return out, nil
}

func (ct *Controller) CeinturesEvaluateAnswers(c echo.Context) error {
	var args EvaluateAnswersIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.evaluateAnswers(args)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) evaluateAnswers(args EvaluateAnswersIn) (EvaluateAnswersOut, error) {
	current, has, err := ct.getEvolution(args.Tokens)
	if err != nil {
		return EvaluateAnswersOut{}, err
	}
	if !has {
		return EvaluateAnswersOut{}, fmt.Errorf("Erreur interne (évolution manquante)")
	}

	// We should check that the stage and the questions are correct,
	// but we "trust" the client for now
	if len(args.Answers) != len(args.Questions) {
		return EvaluateAnswersOut{}, fmt.Errorf("Erreur interne (questions/réponses incohérentes)")
	}

	questionsSource, err := ce.SelectBeltquestionsByDomainAndRank(ct.db, args.Stage.Domain, args.Stage.Rank)
	if err != nil {
		return EvaluateAnswersOut{}, utils.SQLError(err)
	}

	out := EvaluateAnswersOut{Answers: make([]client.QuestionAnswersOut, len(args.Answers))}
	hasPassed := true
	stats := current.Stats[args.Stage.Domain][args.Stage.Rank]
	for i, idQuestion := range args.Questions {
		answer := args.Answers[i]
		qu := questionsSource[idQuestion]
		out.Answers[i], err = tasks.EvaluateQuestion(qu.Enonce, answer)
		if err != nil {
			return EvaluateAnswersOut{}, err
		}
		correct := out.Answers[i].IsCorrect()
		hasPassed = hasPassed && correct
		if correct {
			stats.Success += 1
		} else {
			stats.Failure += 1
		}
	}

	// update the evolution
	newStats := current.Stats
	newStats[args.Stage.Domain][args.Stage.Rank] = stats
	newAdvance := current.Advance
	if hasPassed {
		newAdvance[args.Stage.Domain] += 1
	}

	err = ct.setEvolution(args.Tokens, newAdvance, newStats)
	if err != nil {
		return EvaluateAnswersOut{}, err
	}

	out.Evolution, _, err = ct.getEvolution(args.Tokens)
	if err != nil {
		return EvaluateAnswersOut{}, err
	}

	return out, nil
}

func (ct *Controller) setEvolution(tokens StudentTokens, adv ce.Advance, stats ce.Stats) error {
	if ci := tokens.ClientID; ci != "" {
		id_, err := ct.studentKey.DecryptID(ci)
		if err != nil {
			return fmt.Errorf("Erreur interne: %s", err)
		}
		idStudent := teacher.IdStudent(id_)

		tx, err := ct.db.Begin()
		if err != nil {
			return utils.SQLError(err)
		}

		ev, _, err := ce.SelectBeltevolutionByIdStudent(tx, idStudent)
		if err != nil {
			_ = tx.Rollback()
			return utils.SQLError(err)
		}
		ev.Advance = adv
		ev.Stats = stats

		err = ev.Delete(tx)
		if err != nil {
			_ = tx.Rollback()
			return utils.SQLError(err)
		}

		err = ce.InsertBeltevolution(tx, ev)
		if err != nil {
			_ = tx.Rollback()
			return utils.SQLError(err)
		}

		err = tx.Commit()
		if err != nil {
			return utils.SQLError(err)
		}
	} else {
		ev, _ := ct.anons.get(tokens.AnonymousID)
		ev.Advance = adv
		ev.Stats = stats
		ct.anons.set(tokens.AnonymousID, ev)
	}

	return nil
}
