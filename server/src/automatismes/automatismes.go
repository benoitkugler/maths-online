package automatismes

import (
	"cmp"
	"database/sql"
	"fmt"
	"maps"
	"math/rand"
	"slices"
	"strings"

	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	"github.com/benoitkugler/maths-online/server/src/pass"
	"github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/benoitkugler/maths-online/server/src/sql/trivial"
	"github.com/benoitkugler/maths-online/server/src/tasks"
	"github.com/benoitkugler/maths-online/server/src/utils"
	"github.com/labstack/echo/v4"
)

// the trivial related logic is in prof/trivial

type Controller struct {
	db         *sql.DB
	admin      teacher.Teacher
	studentKey pass.Encrypter
}

func NewController(db *sql.DB, admin teacher.Teacher, studentKey pass.Encrypter) *Controller {
	return &Controller{db, admin, studentKey}
}

// return the admin trivials and questions matching the given level
func (ct *Controller) StudentGetAutomatismes(c echo.Context) error {
	var args GetAutomatismesIn
	if err := c.Bind(&args); err != nil {
		return err
	}
	ts, err := ct.getTrivials(args)
	if err != nil {
		return err
	}
	qs, err := ct.getQuestions(args)
	if err != nil {
		return err
	}
	return c.JSON(200, GetAutomatismesOut{ts, qs})
}

func (args GetAutomatismesIn) match(tags editor.TagGroup) bool {
	return tags.Chapter == "AUTOMATISMES" && tags.Level == args.Level &&
		(args.Sublevel == "" || slices.Contains(tags.SubLevels, args.Sublevel))
}

func (ct *Controller) getTrivials(args GetAutomatismesIn) ([]trivial.Trivial, error) {
	adminTrivials, err := trivial.SelectTrivialsByIdTeachers(ct.db, ct.admin.Id)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	var out []trivial.Trivial
	for _, trivial := range adminTrivials {
		if tags := trivial.Questions.Common().BySection(); args.match(tags) {
			out = append(out, trivial)
		}
	}

	slices.SortFunc(out, func(a, b trivial.Trivial) int { return strings.Compare(a.Name, b.Name) })

	return out, nil
}

func (ct *Controller) getQuestions(args GetAutomatismesIn) ([]editor.Questiongroup, error) {
	groups, err := editor.SelectQuestiongroupsByIdTeachers(ct.db, ct.admin.Id)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	tmp, err := editor.SelectQuestiongroupTagsByIdQuestiongroups(ct.db, groups.IDs()...)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	byQuestion := tmp.ByIdQuestiongroup()

	var out []editor.Questiongroup
	for id, tags := range byQuestion {
		if kind := tags.Tags().BySection(); args.match(kind) {
			out = append(out, groups[id])
		}
	}

	slices.SortFunc(out, func(a, b editor.Questiongroup) int { return cmp.Compare(a.Id, b.Id) })

	return out, nil
}

// StudentInstantiateQuestion selects a question and instantiate random variables.
func (ct *Controller) StudentInstantiateQuestion(c echo.Context) error {
	id, err := utils.QueryParamInt[editor.IdQuestiongroup](c, "idQuestion")
	if err != nil {
		return err
	}
	out, err := ct.instantiateQuestion(id)
	if err != nil {
		return err
	}
	return c.JSON(200, out)
}

func (ct *Controller) instantiateQuestion(id editor.IdQuestiongroup) (InstantiatedAutomatismeQuestion, error) {
	tmp, err := editor.SelectQuestionsByIdGroups(ct.db, id)
	if err != nil {
		return InstantiatedAutomatismeQuestion{}, err
	}
	variants := slices.Collect(maps.Values(tmp))
	if len(variants) == 0 {
		return InstantiatedAutomatismeQuestion{}, fmt.Errorf("internal error: question %d with no variants", id)
	}
	question := variants[rand.Intn(len(variants))]
	content, params := question.Page().Instantiate()
	return InstantiatedAutomatismeQuestion{
		Id:       question.Id,
		Question: content.ToClient(),
		Params:   tasks.NewParams(params),
	}, nil
}

func (ct *Controller) StudentEvaluateQuestion(c echo.Context) error {
	var args EvaluateAutomatismeIn
	if err := c.Bind(&args); err != nil {
		return err
	}
	out, err := ct.evaluateQuestion(args)
	if err != nil {
		return err
	}
	return c.JSON(200, out)
}

func (ct *Controller) evaluateQuestion(args EvaluateAutomatismeIn) (client.QuestionAnswersOut, error) {
	question, err := editor.SelectQuestion(ct.db, args.Id)
	if err != nil {
		return client.QuestionAnswersOut{}, utils.SQLError(err)
	}
	out, err := tasks.EvaluateQuestion(question.Enonce, args.Answer)
	if err != nil {
		return client.QuestionAnswersOut{}, err
	}
	return out, nil
}
