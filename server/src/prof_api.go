package main

import (
	"github.com/benoitkugler/maths-online/server/src/prof/ceintures"
	"github.com/benoitkugler/maths-online/server/src/prof/editor"
	"github.com/benoitkugler/maths-online/server/src/prof/homework"
	"github.com/benoitkugler/maths-online/server/src/prof/reviews"
	"github.com/benoitkugler/maths-online/server/src/prof/teacher"
	"github.com/benoitkugler/maths-online/server/src/prof/trivial"
	"github.com/labstack/echo/v4"
)

func setupProfAPI(e *echo.Echo, tvc *trivial.Controller,
	edit *editor.Controller, tc *teacher.Controller,
	home *homework.Controller, review *reviews.Controller,
	ce *ceintures.Controller,
) {
	e.POST("/prof/inscription", tc.AskInscription)
	e.GET(teacher.ValidateInscriptionEndPoint, tc.ValidateInscription)
	e.POST("/prof/loggin", tc.Loggin)
	e.GET("/api/prof/reset", tc.TeacherResetPassword)

	gr := e.Group("", tc.JWTMiddleware())

	// settings
	gr.GET("/api/prof/settings", tc.TeacherGetSettings)
	gr.POST("/api/prof/settings", tc.TeacherUpdateSettings)

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
	// trivial self-access
	gr.GET("/api/prof/trivial/selfaccess", tvc.TrivialGetSelfaccess)
	gr.POST("/api/prof/trivial/selfaccess", tvc.TrivialUpdateSelfaccess)

	// trivialpoursuit game server
	gr.GET("/api/trivial/sessions", tvc.GetTrivialRunningSessions)
	gr.PUT("/api/trivial/sessions", tvc.LaunchSessionTrivialPoursuit)
	gr.POST("/api/trivial/sessions/start", tvc.StartTrivialGame)
	gr.POST("/api/trivial/sessions/stop", tvc.StopTrivialGame)

	// question editor
	gr.GET("/api/prof/editor/tags", edit.EditorGetTags)
	gr.POST("/api/prof/editor/syntax-hint", edit.EditorGenerateSyntaxHint)

	gr.GET("/api/prof/editor/questiongroups", edit.EditorGetQuestionsIndex)
	gr.POST("/api/prof/editor/questiongroups", edit.EditorSearchQuestions)
	gr.GET("/api/prof/editor/question/duplicate", edit.EditorDuplicateQuestion)
	gr.GET("/api/prof/editor/questiongroup/duplicate", edit.EditorDuplicateQuestiongroup)
	gr.PUT("/api/prof/editor/questiongroup", edit.EditorCreateQuestiongroup)
	gr.POST("/api/prof/editor/questiongroup", edit.EditorUpdateQuestiongroup)
	gr.DELETE("/api/prof/editor/questiongroup", edit.EditorDeleteQuestiongroup)
	gr.POST("/api/prof/editor/questiongroup/tags", edit.EditorUpdateQuestionTags)
	gr.POST("/api/prof/editor/questiongroup/visibility", edit.EditorUpdateQuestiongroupVis)
	gr.GET("/api/prof/editor/question", edit.EditorGetQuestions)
	gr.DELETE("/api/prof/editor/question", edit.EditorDeleteQuestion)
	gr.POST("/api/prof/editor/question/variant", edit.EditorSaveQuestionMeta)
	gr.POST("/api/prof/editor/question/check-params", edit.EditorCheckQuestionParameters)
	gr.POST("/api/prof/editor/question/preview", edit.EditorSaveQuestionAndPreview)
	gr.POST("/api/prof/editor/question/export/latex", edit.EditorQuestionExportLateX)

	// exercice editor
	gr.GET("/api/prof/editor/exercicegroups", edit.EditorGetExercicesIndex)
	gr.POST("/api/prof/editor/exercicegroups", edit.EditorSearchExercices)
	gr.POST("/api/prof/editor/exercicegroup", edit.EditorUpdateExercicegroup)
	gr.DELETE("/api/prof/editor/exercicegroup", edit.EditorDeleteExercicegroup)
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
	gr.POST("/api/prof/editor/exercice/export/latex", edit.EditorExerciceExportLateX)

	// homework activity
	gr.GET("/api/prof/homework", home.HomeworkGetSheets)
	gr.PUT("/api/prof/homework", home.HomeworkCreateSheet)
	gr.POST("/api/prof/homework", home.HomeworkUpdateSheet)
	gr.DELETE("/api/prof/homework", home.HomeworkDeleteSheet)
	gr.POST("/api/prof/homework/sheet/copy", home.HomeworkCopySheet)
	gr.PUT("/api/prof/homework/travail", home.HomeworkCreateTravailWith)
	gr.PUT("/api/prof/homework/travail/anonymous", home.HomeworkCreateTravail)
	gr.POST("/api/prof/homework/travail", home.HomeworkUpdateTravail)
	gr.DELETE("/api/prof/homework/travail", home.HomeworkDeleteTravail)
	gr.POST("/api/prof/homework/travail/copy", home.HomeworkCopyTravail)
	gr.DELETE("/api/prof/homework/sheet", home.HomeworkRemoveTask)
	gr.PUT("/api/prof/homework/sheet/exercice", home.HomeworkAddExercice)
	gr.GET("/api/prof/homework/sheet/monoquestion", home.HomeworkGetMonoquestion)
	gr.PUT("/api/prof/homework/sheet/monoquestion", home.HomeworkAddMonoquestion)
	gr.GET("/api/prof/homework/sheet/randommonoquestion", home.HomeworkGetRandomMonoquestion)
	gr.PUT("/api/prof/homework/sheet/randommonoquestion", home.HomeworkAddRandomMonoquestion)
	gr.POST("/api/prof/homework/sheet/monoquestion", home.HomeworkUpdateMonoquestion)
	gr.POST("/api/prof/homework/sheet/randommonoquestion", home.HomeworkUpdateRandomMonoquestion)
	gr.POST("/api/prof/homework/sheet", home.HomeworkReorderSheetTasks)
	gr.GET("/api/prof/homework/sheet/missing-hint", home.HomeworkMissingTasksHint)
	gr.POST("/api/prof/homework/marks", home.HomeworkGetMarks)
	gr.GET("/api/prof/homework/dispences", home.HomeworkGetDispenses)
	gr.POST("/api/prof/homework/dispences", home.HomeworkSetDispense)

	// ceintures
	gr.GET("/api/prof/ceintures/scheme", ce.CeinturesGetScheme)
	gr.POST("/api/prof/ceintures/pending", ce.CeinturesGetPending)
	gr.POST("/api/prof/ceintures/questions", ce.CeinturesGetQuestions)
	gr.POST("/api/prof/ceintures/question-meta", ce.CeinturesUpdateQuestion)
	gr.PUT("/api/prof/ceintures/question", ce.CeinturesCreateQuestion)
	gr.POST("/api/prof/ceintures/question", ce.CeinturesSaveQuestion)
	gr.DELETE("/api/prof/ceintures/question", ce.CeinturesDeleteQuestion)
	gr.POST("/api/prof/ceintures/question-duplicate", ce.CeinturesDuplicateQuestion)
	gr.GET("/api/prof/ceintures/advance", ce.CeinturesGetStudentsAdvance)

	// review page
	gr.PUT("/api/prof/review", review.ReviewCreate)
	gr.GET("/api/prof/reviews", review.ReviewsList)
	gr.GET("/api/prof/review", review.ReviewLoad)
	gr.GET("/api/prof/review/target", review.ReviewLoadTarget)
	gr.DELETE("/api/prof/review", review.ReviewDelete)
	gr.POST("/api/prof/review/approval", review.ReviewUpdateApproval)
	gr.POST("/api/prof/review/comments", review.ReviewUpdateCommnents)
	gr.POST("/api/prof/review/accept", review.ReviewAccept)
}
