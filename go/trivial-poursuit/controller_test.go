package trivialpoursuit

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
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

func TestConcurrentEvents(t *testing.T) {
	game.DebugLogger.SetOutput(io.Discard)

	ct := newGameController()
	go ct.startLoop()

	server := httptest.NewServer(http.HandlerFunc(ct.setupWebSocket))
	defer server.Close()

	client1, _, err := websocket.DefaultDialer.Dial(websocketURL(t, server.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	client2, _, err := websocket.DefaultDialer.Dial(websocketURL(t, server.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	client3, _, err := websocket.DefaultDialer.Dial(websocketURL(t, server.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	clientLoop := func(client *websocket.Conn) {
		for i := range [200]int{} {
			err := client.WriteJSON(game.ClientEvent{Event: game.Ping(fmt.Sprintf("message %d", i))})
			if err != nil {
				panic(err)
			}
		}
	}

	go clientLoop(client1)
	go clientLoop(client2)
	go clientLoop(client3)

	for range [200]int{} {
		ct.broadcastEvents <- game.GameEvents{}
	}
}

func TestEvents(t *testing.T) {
	game.QuestionDurationLimit = time.Millisecond * 50

	ct := newGameController()
	go ct.startLoop()

	server := httptest.NewServer(http.HandlerFunc(ct.setupWebSocket))
	defer server.Close()

	client, _, err := websocket.DefaultDialer.Dial(websocketURL(t, server.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		_, _, err = client.ReadMessage()
		if err != nil {
			panic(err)
		}
	}()

	ct.game.EmitQuestion() // launch the timer

	time.Sleep(time.Millisecond * 200)
}

func TestClientInvalidMessage(t *testing.T) {
	WarningLogger.SetOutput(io.Discard)

	ct := newGameController()
	go ct.startLoop()

	server := httptest.NewServer(http.HandlerFunc(ct.setupWebSocket))
	defer server.Close()

	client, _, err := websocket.DefaultDialer.Dial(websocketURL(t, server.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	err = client.WriteMessage(websocket.TextMessage, []byte("BAD"))
	if err != nil {
		t.Fatal(err)
	}

	if _, _, err := client.ReadMessage(); err == nil {
		t.Fatal("expected error for invalid event")
	}
}
