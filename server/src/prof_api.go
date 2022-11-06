package main

import (
	"github.com/benoitkugler/maths-online/prof/editor"
	"github.com/benoitkugler/maths-online/prof/homework"
	"github.com/benoitkugler/maths-online/prof/teacher"
	"github.com/benoitkugler/maths-online/prof/trivial"
	"github.com/labstack/echo/v4"
)

func setupProfAPI(e *echo.Echo, tvc *trivial.Controller,
	edit *editor.Controller, tc *teacher.Controller,
	home *homework.Controller,
) {
	e.POST("/prof/inscription", tc.AskInscription)
	e.GET(teacher.ValidateInscriptionEndPoint, tc.ValidateInscription)
	e.POST("/prof/loggin", tc.Loggin)

	gr := e.Group("", tc.JWTMiddleware())

	// classrooms
	gr.GET("/api/prof/classrooms", tc.TeacherGetClassrooms)
	gr.PUT("/api/prof/classrooms", tc.TeacherCreateClassroom)
	gr.POST("/api/prof/classrooms", tc.TeacherUpdateClassroom)
	gr.DELETE("/api/prof/classrooms", tc.TeacherDeleteClassroom)

	gr.GET("/api/prof/classrooms/students", tc.TeacherGetClassroomStudents)
	gr.PUT("/api/prof/classrooms/students", tc.TeacherAddStudent)
	gr.POST("/api/prof/classrooms/students", tc.TeacherUpdateStudent)
	gr.DELETE("/api/prof/classrooms/students", tc.TeacherDeleteStudent)
	gr.POST("/api/prof/classrooms/students/import", tc.TeacherImportStudents)
	gr.GET("/api/prof/classrooms/students/connect", tc.TeacherGenerateClassroomCode)

	// trivial activity
	gr.GET("/api/prof/trivial/config", tvc.GetTrivialPoursuit)
	gr.PUT("/api/prof/trivial/config", tvc.CreateTrivialPoursuit)
	gr.POST("/api/prof/trivial/config", tvc.UpdateTrivialPoursuit)
	gr.DELETE("/api/prof/trivial/config", tvc.DeleteTrivialPoursuit)
	gr.POST("/api/prof/trivial/config/visibility", tvc.UpdateTrivialVisiblity)
	gr.GET("/api/prof/trivial/config/duplicate", tvc.DuplicateTrivialPoursuit)
	gr.POST("/api/prof/trivial/config/check-missing-questions", tvc.CheckMissingQuestions)
	gr.GET("/api/prof/trivial/monitor", tvc.TrivialTeacherMonitor)
	// e.GET("/prof/trivial/monitor", tvc.TrivialTeacherMonitor, tc.JWTMiddlewareForQuery())

	// trivialpoursuit game server
	gr.GET("/api/trivial/sessions", tvc.GetTrivialRunningSessions)
	gr.PUT("/api/trivial/sessions", tvc.LaunchSessionTrivialPoursuit)
	gr.POST("/api/trivial/sessions/stop", tvc.StopTrivialGame)

	// question editor
	gr.PUT("/api/prof/editor/new", edit.EditorStartSession)
	gr.GET("/api/prof/editor/pause-preview", edit.EditorPausePreview)
	gr.GET("/api/prof/editor/tags", edit.EditorGetTags)

	gr.POST("/api/prof/editor/questiongroups", edit.EditorSearchQuestions)
	gr.GET("/api/prof/editor/question/duplicate", edit.EditorDuplicateQuestion)
	gr.GET("/api/prof/editor/questiongroup/duplicate", edit.EditorDuplicateQuestiongroup)
	gr.PUT("/api/prof/editor/questiongroup", edit.EditorCreateQuestiongroup)
	gr.POST("/api/prof/editor/questiongroup", edit.EditorUpdateQuestiongroup)
	gr.POST("/api/prof/editor/questiongroup/tags", edit.EditorUpdateQuestionTags)
	gr.POST("/api/prof/editor/questiongroup/visibility", edit.EditorUpdateQuestiongroupVis)
	gr.GET("/api/prof/editor/question", edit.EditorGetQuestions)
	gr.DELETE("/api/prof/editor/question", edit.EditorDeleteQuestion)
	gr.POST("/api/prof/editor/question/variant", edit.EditorSaveQuestionMeta)
	gr.POST("/api/prof/editor/question/check-params", edit.EditorCheckQuestionParameters)
	gr.POST("/api/prof/editor/question/preview", edit.EditorSaveQuestionAndPreview)

	// exercice editor
	gr.POST("/api/prof/editor/exercicegroups", edit.EditorSearchExercices)
	gr.POST("/api/prof/editor/exercicegroup", edit.EditorUpdateExercicegroup)
	gr.POST("/api/prof/editor/exercicegroup/tags", edit.EditorUpdateExerciceTags)
	gr.GET("/api/prof/editor/exercicegroup/duplicate", edit.EditorDuplicateExercicegroup)
	gr.GET("/api/prof/editor/exercice", edit.EditorGetExerciceContent)
	gr.PUT("/api/prof/editor/exercice", edit.EditorCreateExercice)
	gr.DELETE("/api/prof/editor/exercice", edit.EditorDeleteExercice)
	gr.POST("/api/prof/editor/exercice", edit.EditorSaveExerciceMeta)
	gr.GET("/api/prof/editor/exercice/duplicate", edit.EditorDuplicateExercice)
	gr.PUT("/api/prof/editor/exercice/questions", edit.EditorExerciceCreateQuestion)
	gr.POST("/api/prof/editor/exercice/questions/import", edit.EditorExerciceImportQuestion)
	gr.POST("/api/prof/editor/exercice/questions/duplicate", edit.EditorExerciceDuplicateQuestion)
	gr.POST("/api/prof/editor/exercice/questions", edit.EditorExerciceUpdateQuestions)
	gr.POST("/api/prof/editor/exercicegroup/visibility", edit.EditorUpdateExercicegroupVis)
	gr.POST("/api/prof/editor/exercice/check-params", edit.EditorCheckExerciceParameters)
	gr.POST("/api/prof/editor/exercice/preview", edit.EditorSaveExerciceAndPreview)

	// homework activity
	gr.GET("/api/prof/homework", home.HomeworkGetSheets)
	gr.PUT("/api/prof/homework", home.HomeworkCreateSheet)
	gr.POST("/api/prof/homework", home.HomeworkUpdateSheet)
	gr.DELETE("/api/prof/homework", home.HomeworkDeleteSheet)
	gr.POST("/api/prof/homework/copy-sheet", home.HomeworkCopySheet)
	gr.DELETE("/api/prof/homework/sheet", home.HomeworkRemoveTask)
	gr.PUT("/api/prof/homework/sheet/exercice", home.HomeworkAddExercice)
	gr.PUT("/api/prof/homework/sheet/monoquestion", home.HomeworkAddMonoquestion)
	gr.POST("/api/prof/homework/sheet/monoquestion", home.HomeworkUpdateMonoquestion)
	gr.POST("/api/prof/homework/sheet", home.HomeworkReorderSheetTasks)
	gr.POST("/api/prof/homework/marks", home.HomeworkGetMarks)
}
