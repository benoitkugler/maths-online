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

// Sample choose an index between 0 and len(weights)-1 at random, with the given weights,
// which must sum up to 1.
func SampleIndex(weights []float64) int {
	cumWeights := make([]float64, len(weights)) // last entry is 1
	sum := 0.
	for i, w := range weights {
		sum += w
		cumWeights[i] = sum
	}
	alea := rand.Float64()
	for i, cumWeight := range cumWeights {
		if alea < cumWeight {
			return i
		}
	}
	return len(weights) - 1
}
