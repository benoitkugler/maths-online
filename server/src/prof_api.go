package main

import (
	"github.com/benoitkugler/maths-online/prof/editor"
	"github.com/benoitkugler/maths-online/prof/teacher"
	trivialpoursuit "github.com/benoitkugler/maths-online/prof/trivial-poursuit"
	"github.com/labstack/echo/v4"
)

//go:generate ../../../structgen/apigen -source=prof_api.go -out=../../prof/src/controller/api_gen.ts

func setupProfAPI(e *echo.Echo, trivial *trivialpoursuit.Controller,
	edit *editor.Controller, tc *teacher.Controller,
) {
	e.POST("/prof/inscription", tc.AskInscription)
	e.GET(teacher.ValidateInscriptionEndPoint, tc.ValidateInscription)
	e.POST("/prof/loggin", tc.Loggin)

	gr := e.Group("", tc.JWTMiddleware())

	gr.GET("/prof/trivial/config", trivial.GetTrivialPoursuit)
	gr.PUT("/prof/trivial/config", trivial.CreateTrivialPoursuit)
	gr.POST("/prof/trivial/config", trivial.UpdateTrivialPoursuit)
	gr.DELETE("/prof/trivial/config", trivial.DeleteTrivialPoursuit)
	gr.GET("/prof/trivial/config/duplicate", trivial.DuplicateTrivialPoursuit)
	gr.POST("/prof/trivial/config/check-missing-questions", trivial.CheckMissingQuestions)

	// trivialpoursuit game server
	gr.POST("/trivial/launch_session", trivial.LaunchSessionTrivialPoursuit)
	gr.DELETE("/trivial/launch_session", trivial.StopSessionTrivialPoursuit)

	gr.PUT("/prof/editor/api/new", edit.EditorStartSession)
	gr.GET("/prof/editor/api/tags", edit.EditorGetTags)
	gr.POST("/prof/editor/api/questions", edit.EditorSearchQuestions)
	gr.GET("/prof/editor/api/question-duplicate-one", edit.EditorDuplicateQuestion)
	gr.GET("/prof/editor/api/question-duplicate", edit.EditorDuplicateQuestionWithDifficulty)
	gr.PUT("/prof/editor/api/question", edit.EditorCreateQuestion)
	gr.GET("/prof/editor/api/question", edit.EditorGetQuestion)
	gr.DELETE("/prof/editor/api/question", edit.EditorDeleteQuestion)
	gr.POST("/prof/editor/api/question", edit.EditorSaveAndPreview)
	gr.GET("/prof/editor/api/pause-preview", edit.EditorPausePreview)
	gr.POST("/prof/editor/api/question/tags", edit.EditorUpdateTags)
	gr.POST("/prof/editor/api/check-params", edit.EditorCheckParameters)
}
