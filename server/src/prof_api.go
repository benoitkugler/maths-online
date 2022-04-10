package main

import (
	"github.com/benoitkugler/maths-online/prof/editor"
	trivialpoursuit "github.com/benoitkugler/maths-online/trivial-poursuit"
	"github.com/labstack/echo/v4"
)

//go:generate ../../../structgen/apigen -source=prof_api.go -out=../../prof/src/controller/api_gen.ts

func setupProfAPI(e *echo.Echo, trivial *trivialpoursuit.Controller, edit *editor.Controller) {
	// trivialpoursuit game server
	e.POST("/trivial/launch_game", trivial.LaunchGame)

	e.PUT("/prof/editor/api/new", edit.EditStartSession)
	e.POST("/prof/editor/api/check-params", edit.EditCheckParameters)
	e.POST("/prof/editor/api/save", edit.EditSaveAndPreview)
}
