package editor

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/benoitkugler/maths-online/maths/questions"
	"github.com/gorilla/websocket"
)

func websocketURL(t *testing.T, s string) string {
	t.Helper()

	u, err := url.Parse(s)
	if err != nil {
		t.Fatal(err)
	}

	u.Scheme = "ws"
	return u.String()
}

func TestLoopback(t *testing.T) {
	loopback := newLoopbackController("newID")

	// start the websocket for the loopback
	go func() {
		ctx, cancelFunc := context.WithTimeout(context.Background(), sessionTimeout)
		loopback.startLoop(ctx) // block
		cancelFunc()
	}()

	server := httptest.NewServer(http.HandlerFunc(loopback.setupWebSocket))
	defer server.Close()

	cl, _, err := websocket.DefaultDialer.Dial(websocketURL(t, server.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	err = cl.WriteJSON(LoopbackClientEventWrapper{loopbackPing{}})
	if err != nil {
		t.Fatal(err)
	}

	loopback.setQuestion(questions.QuestionInstance{Title: "Test", Enonce: questions.EnonceInstance{
		questions.NumberFieldInstance{ID: 0},
	}})

	loopback.unsetQuestion()

	_, json, err := cl.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(json))

	_, json, err = cl.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(json))

	err = cl.WriteJSON(LoopbackClientEventWrapper{loopbackQuestionValidIn{}})
	if err != nil {
		t.Fatal(err)
	}

	_, json, err = cl.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(json))
}
