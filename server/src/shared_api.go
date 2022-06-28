package main

import (
	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/questions/client"
	"github.com/benoitkugler/maths-online/prof/editor"
	"github.com/labstack/echo/v4"
)

//go:generate ../../../structgen/structgen -source=shared_api.go -mode=dart:../../eleve/lib/shared_gen.dart

type CheckExpressionOut struct {
	Reason  string
	IsValid bool
}

// utility end point used by clients to perform
// on the fly validation of expression fields
func checkExpressionSyntax(c echo.Context) error {
	expr := c.QueryParam("expression")
	_, err := expression.Parse(expr)
	out := CheckExpressionOut{IsValid: err == nil}
	if err != nil {
		out.Reason = err.Error()
	}
	return c.JSON(200, out)
}

type InstantiateQuestionsOut = editor.InstantiateQuestionsOut

// standalone endpoint to check if an answer is correct
func instantiateQuestions(ct *editor.Controller, c echo.Context) error {
	var args []int64
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.InstantiateQuestions(args)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

type EvaluateQuestionIn struct {
	Answer     client.QuestionAnswersIn `dart-extern:"client:questions/types.gen.dart"`
	Params     []editor.VarEntry
	IdQuestion int64
}

// standalone endpoint to check if an answer is correct
func evaluateQuestion(ct *editor.Controller, c echo.Context) error {
	var args EvaluateQuestionIn
	if err := c.Bind(&args); err != nil {
		return err
	}

	out, err := ct.EvaluateQuestion(args.IdQuestion, editor.Answer{Params: args.Params, Answer: args.Answer})
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

// import for dart code generation
type (
	EvaluateExerciceIn  = editor.EvaluateExerciceIn
	EvaluateExerciceOut = editor.EvaluateExerciceOut
)

type Exercice struct {
	Exercice    editor.InstantiatedExercice
	Progression editor.ProgressionExt
}
