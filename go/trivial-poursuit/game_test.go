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
	// ProgressLogger.SetOutput(os.Stdout)

	ct := newGameController(GameOptions{4, 0}) // do not start a game
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
			err := client.WriteJSON(game.ClientEvent{Event: game.Ping{Info: fmt.Sprintf("message %d", i)}})
			if err != nil {
				panic(err)
			}
		}
	}

	go clientLoop(client1)
	go clientLoop(client2)
	go clientLoop(client3)

	time.Sleep(time.Second / 10)
}

func TestEvents(t *testing.T) {
	game.DebugLogger.SetOutput(io.Discard)

	ct := newGameController(GameOptions{4, time.Millisecond * 50})
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

	ct := newGameController(GameOptions{2, 0})
	go ct.startLoop()

	server := httptest.NewServer(http.HandlerFunc(ct.setupWebSocket))
	defer server.Close()

	client, _, err := websocket.DefaultDialer.Dial(websocketURL(t, server.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	err = client.ReadJSON(&[]game.StateUpdate{}) // player join
	if err != nil {
		t.Fatal(err)
	}
	err = client.ReadJSON(&[]game.StateUpdate{}) // game lobby
	if err != nil {
		t.Fatal(err)
	}

	err = client.WriteMessage(websocket.TextMessage, []byte("BAD"))
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = client.ReadMessage()
	if err == nil {
		t.Fatal("expected error on invalid input")
	}
}

func TestStartGame(t *testing.T) {
	WarningLogger.SetOutput(io.Discard)

	ct := newGameController(GameOptions{2, 0})

	go ct.startLoop()

	server := httptest.NewServer(http.HandlerFunc(ct.setupWebSocket))
	defer server.Close()

	client1, _, err := websocket.DefaultDialer.Dial(websocketURL(t, server.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	var events game.StateUpdates
	if err = client1.ReadJSON(&events); err != nil {
		t.Fatal(err)
	}

	ct.gameLock.Lock()
	if ct.game.IsPlaying() {
		t.Fatal("game should not have started")
	}
	ct.gameLock.Unlock()

	client2, _, err := websocket.DefaultDialer.Dial(websocketURL(t, server.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Millisecond * 20)

	if err = client2.ReadJSON(&events); err != nil { // PlayerJoin event
		t.Fatal(err)
	}
	if err = client2.ReadJSON(&events); err != nil { // GameStart event
		t.Fatal(err)
	}

	ct.gameLock.Lock()
	if !ct.game.IsPlaying() {
		t.Fatal("game should have started")
	}
	ct.gameLock.Unlock()
}

func TestInvalidJoin(t *testing.T) {
	WarningLogger.SetOutput(io.Discard)

	ct := newGameController(GameOptions{1, 0})

	go ct.startLoop()

	server := httptest.NewServer(http.HandlerFunc(ct.setupWebSocket))
	defer server.Close()

	// first client join and launch the game
	_, _, err := websocket.DefaultDialer.Dial(websocketURL(t, server.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	// second client try to join
	client2, _, err := websocket.DefaultDialer.Dial(websocketURL(t, server.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = client2.ReadMessage()
	if err == nil {
		t.Fatal("expected error when joining started game")
	}
}
