// Package runs a server handling predefined questions
package main

import (
	"net/http"
	"strconv"

	"github.com/benoitkugler/maths-online/maths/exercice"
	"github.com/benoitkugler/maths-online/maths/exercice/client"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods:  append(middleware.DefaultCORSConfig.AllowMethods, http.MethodOptions),
		AllowHeaders:  []string{"Authorization", "Content-Type", "Access-Control-Allow-Origin"},
		ExposeHeaders: []string{"Content-Disposition"},
	}))

	e.GET("/", func(c echo.Context) error {
		var out []client.Question
		for _, qu := range exercice.PredefinedQuestions {
			out = append(out, qu.ToClient())
		}
		return c.JSON(200, out)
	})

	e.POST("/syntaxe/:index", func(c echo.Context) error {
		index, _ := strconv.Atoi(c.Param("index"))
		var data client.QuestionSyntaxCheckIn
		err := c.Bind(&data)
		if err != nil {
			return err
		}

		var out client.QuestionSyntaxCheckOut
		err = exercice.PredefinedQuestions[index].CheckSyntaxe(data)
		if err != nil {
			out.Reason = err.Error()
		} else {
			out.IsValid = true
		}

		c.JSON(200, out)

		return nil
	})

	e.POST("/answer/:index", func(c echo.Context) error {
		index, _ := strconv.Atoi(c.Param("index"))

		var data client.QuestionAnswersIn
		err := c.Bind(&data)
		if err != nil {
			return err
		}

		out := exercice.PredefinedQuestions[index].EvaluateAnswer(data)
		c.JSON(200, out)

		return nil
	})

	e.Start("localhost:3030")
}
