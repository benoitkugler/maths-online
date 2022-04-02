package editor

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

const LoopbackEndpoint = "/prof-loopback/:session_id"

// EditStartSession setup a new editing session.
// In particular, it launches in the background a
// `loopbackController` instance to handle preview requests.
func (ct *Controller) EditStartSession(c echo.Context) error {
	out := ct.startSession()

	return c.JSON(200, out)
}

type SaveAndPreviewIn struct {
	SessionID string
	// TODO: read the actual question
}

func (ct *Controller) EditSaveAndPreview(c echo.Context) error {
	var args SaveAndPreviewIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	err := ct.saveAndPreview(args.SessionID)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

// AccessLoopback establish a connection with the embedded preview app
func (ct *Controller) AccessLoopback(c echo.Context) error {
	sessionID := c.Param("session_id")

	loopback, ok := ct.sessions[sessionID]
	if !ok {
		return fmt.Errorf("invalid session ID %s", sessionID)
	}

	// connect to the websocket handler, which handle errors
	loopback.setupWebSocket(c.Response().Writer, c.Request())

	return nil
}
