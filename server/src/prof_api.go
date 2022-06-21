package main

import (
	"github.com/benoitkugler/maths-online/prof/editor"
	"github.com/benoitkugler/maths-online/prof/teacher"
	"github.com/benoitkugler/maths-online/prof/trivial"
	"github.com/labstack/echo/v4"
)

//go:generate ../../../structgen/apigen -source=prof_api.go -out=../../prof/src/controller/api_gen.ts

func setupProfAPI(e *echo.Echo, tvc *trivial.Controller,
	edit *editor.Controller, tc *teacher.Controller,
) {
	e.POST("/prof/inscription", tc.AskInscription)
	e.GET(teacher.ValidateInscriptionEndPoint, tc.ValidateInscription)
	e.POST("/prof/loggin", tc.Loggin)

	gr := e.Group("", tc.JWTMiddleware())

	// classrooms
	gr.GET("/prof/classrooms/api", tc.TeacherGetClassrooms)
	gr.PUT("/prof/classrooms/api", tc.TeacherCreateClassroom)
	gr.POST("/prof/classrooms/api", tc.TeacherUpdateClassroom)
	gr.DELETE("/prof/classrooms/api", tc.TeacherDeleteClassroom)

	gr.GET("/prof/classrooms/api/students", tc.TeacherGetClassroomStudents)
	gr.POST("/prof/classrooms/api/students/import", tc.TeacherImportStudents)
	gr.GET("/prof/classrooms/api/students/connect", tc.TeacherGenerateClassroomCode)

	// trivial configurations
	gr.GET("/prof/trivial/config", tvc.GetTrivialPoursuit)
	gr.PUT("/prof/trivial/config", tvc.CreateTrivialPoursuit)
	gr.POST("/prof/trivial/config", tvc.UpdateTrivialPoursuit)
	gr.DELETE("/prof/trivial/config", tvc.DeleteTrivialPoursuit)
	gr.POST("/prof/trivial/config/visibility", tvc.UpdateTrivialVisiblity)
	gr.GET("/prof/trivial/config/duplicate", tvc.DuplicateTrivialPoursuit)
	gr.POST("/prof/trivial/config/check-missing-questions", tvc.CheckMissingQuestions)

	// trivialpoursuit game server
	gr.GET("/trivial/sessions", tvc.GetTrivialRunningSessions)
	gr.PUT("/trivial/sessions", tvc.LaunchSessionTrivialPoursuit)
	gr.POST("/trivial/sessions/stop", tvc.StopTrivialGame)

	// question editor
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
	gr.POST("/prof/editor/api/question/group-tags", edit.EditorUpdateGroupTags)
	gr.POST("/prof/editor/api/question/visibility", edit.QuestionUpdateVisiblity)
	gr.POST("/prof/editor/api/check-params", edit.EditorCheckParameters)

	// exercice editor
	gr.GET("/prof/editor/api/exercices", edit.ExercicesGetList)
	gr.PUT("/prof/editor/api/exercice", edit.ExerciceCreate)
	gr.DELETE("/prof/editor/api/exercice", edit.ExerciceDelete)
	gr.POST("/prof/editor/api/exercice", edit.ExerciceUpdate)
	gr.PUT("/prof/editor/api/exercice/questions", edit.ExerciceCreateQuestion)
	gr.POST("/prof/editor/api/exercice/questions", edit.ExerciceUpdateQuestions)
}
