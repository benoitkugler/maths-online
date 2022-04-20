package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func randomID(numberOnly bool, length int) string {
	choices := "abcdefghijklmnnopqrst0123456789"
	if numberOnly {
		choices = "0123456789"
	}
	out := make([]byte, length)
	for i := range out {
		out[i] = choices[rand.Intn(len(choices))]
	}
	return string(out[:])
}

// RandomID generates a random ID, for which `isTaken` is false.
func RandomID(numberOnly bool, length int, isTaken func(string) bool) string {
	newID := randomID(numberOnly, length)
	// avoid (unlikely) collisions
	for taken := isTaken(newID); taken; newID = randomID(numberOnly, length) {
	}
	return newID
}

// WebsocketError format `err` and send a Control message to `ws`
func WebsocketError(ws *websocket.Conn, err error) {
	message := websocket.FormatCloseMessage(websocket.CloseUnsupportedData, err.Error())
	ws.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
}

func SQLError(err error) error {
	return fmt.Errorf("SQL request failed : %s", err)
}

// QueryParamInt64 parse the query param `name` to an int64
func QueryParamInt64(c echo.Context, name string) (int64, error) {
	idS := c.QueryParam(name)
	id, err := strconv.ParseInt(idS, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid ID parameter %s : %s", idS, err)
	}
	return id, nil
}
