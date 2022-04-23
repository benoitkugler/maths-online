package main

import (
	"github.com/benoitkugler/maths-online/prof/editor"
	trivialpoursuit "github.com/benoitkugler/maths-online/prof/trivial-poursuit"
	"github.com/labstack/echo/v4"
)

//go:generate ../../../structgen/apigen -source=prof_api.go -out=../../prof/src/controller/api_gen.ts

func setupProfAPI(e *echo.Echo, trivial *trivialpoursuit.Controller, edit *editor.Controller) {
	e.GET("/prof/trivial/config", trivial.GetTrivialPoursuit)
	e.PUT("/prof/trivial/config", trivial.CreateTrivialPoursuit)
	e.POST("/prof/trivial/config", trivial.UpdateTrivialPoursuit)
	e.DELETE("/prof/trivial/config", trivial.DeleteTrivialPoursuit)

	// trivialpoursuit game server
	e.POST("/trivial/launch_session", trivial.LaunchSessionTrivialPoursuit)

	e.PUT("/prof/editor/api/new", edit.EditorStartSession)
	e.GET("/prof/editor/api/tags", edit.EditorGetTags)
	e.POST("/prof/editor/api/questions", edit.EditorSearchQuestions)
	e.GET("/prof/editor/api/question-duplicate", edit.EditorDuplicateQuestion)
	e.PUT("/prof/editor/api/question", edit.EditorCreateQuestion)
	e.GET("/prof/editor/api/question", edit.EditorGetQuestion)
	e.DELETE("/prof/editor/api/question", edit.EditorDeleteQuestion)
	e.POST("/prof/editor/api/question", edit.EditorSaveAndPreview)
	e.GET("/prof/editor/api/pause-preview", edit.EditorPausePreview)
	e.POST("/prof/editor/api/question/tags", edit.EditorUpdateTags)
	e.POST("/prof/editor/api/check-params", edit.EditorCheckParameters)
}
