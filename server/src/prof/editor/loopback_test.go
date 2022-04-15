package editor

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/benoitkugler/maths-online/maths/exercice"
	"github.com/benoitkugler/maths-online/maths/exercice/client"
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

	err = cl.WriteJSON(clientData{Kind: Ping})
	if err != nil {
		t.Fatal(err)
	}

	loopback.setQuestion(exercice.QuestionInstance{Title: "Test", Enonce: exercice.EnonceInstance{
		exercice.NumberFieldInstance{ID: 0},
	}})

	err = cl.WriteJSON(clientData{Kind: CheckSyntaxIn, Data: client.QuestionSyntaxCheckIn{Answer: client.NumberAnswer{}}})
	if err != nil {
		t.Fatal(err)
	}

	_, json, err := cl.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(json))

	loopback.unsetQuestion()

	_, json, err = cl.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(json))

	_, json, err = cl.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(json))

	err = cl.WriteJSON(clientData{Kind: ValidAnswerIn, Data: client.QuestionAnswersIn{}})
	if err != nil {
		t.Fatal(err)
	}

	_, json, err = cl.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(json))
}
