package editor

import (
	taAPI "github.com/benoitkugler/maths-online/server/src/tasks"
	"github.com/labstack/echo/v4"
)

// LoopackEvaluateQuestion expects a question definition, a set of
// random variables, an answer, and performs the evaluation.
func (ct *Controller) LoopackEvaluateQuestion(c echo.Context) error {
	var args LoopackEvaluateQuestionIn

	if err := c.Bind(&args); err != nil {
		return err
	}

	ans, err := taAPI.EvaluateQuestion(args.Question.Enonce, args.Answer)
	if err != nil {
		return err
	}

	out := LoopbackEvaluateQuestionOut{ans}

	return c.JSON(200, out)
}

// LoopbackShowQuestionAnswer expects a question, random parameters,
// and returns the correct answer for these parameters.
// It is shared between question and exercice loopback preview.
func (ct *Controller) LoopbackShowQuestionAnswer(c echo.Context) error {
	var args LoopbackShowQuestionAnswerIn

	if err := c.Bind(&args); err != nil {
		return err
	}

	p, err := args.Params.ToMap()
	if err != nil {
		return err
	}
	instance, err := args.Question.Enonce.InstantiateWith(p)
	if err != nil {
		return err
	}
	ans := instance.CorrectAnswer()

	out := LoopbackShowQuestionAnswerOut{ans}
	return c.JSON(200, out)
}
