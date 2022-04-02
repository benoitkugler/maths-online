package utils

import (
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
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
