package main

import (
	"github.com/benoitkugler/maths-online/maths/expression"
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
