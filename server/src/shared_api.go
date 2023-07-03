package main

import (
	"github.com/benoitkugler/maths-online/server/src/maths/expression"
	ed "github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/tasks"
	"github.com/labstack/echo/v4"
)

type CheckExpressionOut struct {
	Reason  string
	IsValid bool
}

// utility end point used by clients to perform
// on the fly validation of expression fields
func checkExpressionSyntax(c echo.Context) error {
	expr := c.QueryParam("expression")
	_, err := expression.ParseCompound(expr)
	out := CheckExpressionOut{IsValid: err == nil}
	if err != nil {
		out.Reason = err.Error()
	}
	return c.JSON(200, out)
}

type InstantiateQuestionsOut = tasks.InstantiateQuestionsOut

// standalone endpoint to check if an answer is correct
func instantiateQuestions(db ed.DB, c echo.Context) error {
	var args []ed.IdQuestion
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := tasks.InstantiateQuestions(db, args)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

type EvaluateQuestionIn = tasks.EvaluateQuestionIn

// standalone endpoint to check if an answer is correct
func evaluateQuestion(db ed.DB, c echo.Context) error {
	var args EvaluateQuestionIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := args.Evaluate(db)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

// import for dart code generation
type (
	EvaluateWorkIn  = tasks.EvaluateWorkIn
	EvaluateWorkOut = tasks.EvaluateWorkOut
)

type StudentWork struct {
	Exercice    tasks.InstantiatedWork
	Progression tasks.ProgressionExt
}

// standalone endpoint to check if an exercice answer is correct
// note that this API does not handle progression persistence,
// and do not support RandomMonoquestion either
func evaluateExercice(db ed.DB, c echo.Context) error {
	var args EvaluateWorkIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := args.Evaluate(db, -1)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}
