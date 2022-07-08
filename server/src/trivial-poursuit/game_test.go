package trivialpoursuit

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/maths/questions"
	exClient "github.com/benoitkugler/maths-online/maths/questions/client"
	"github.com/benoitkugler/maths-online/pass"
	"github.com/benoitkugler/maths-online/prof/editor"
	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
	"github.com/benoitkugler/maths-online/utils/testutils"
	"github.com/gorilla/websocket"
)

var page = questions.QuestionPage{
	Enonce: questions.Enonce{
		questions.NumberFieldBlock{Expression: "1"},
	},
}

var exQu = game.WeigthedQuestions{
	Questions: []editor.Question{{Id: 1, Page: page}, {Id: 2, Page: page}},
	Weights:   []float64{1. / 3, 2. / 3},
}

var questionsList = game.QuestionPool{exQu, exQu, exQu, exQu, exQu}

func websocketURLWithClientID(t *testing.T, urlS, clientID string) string {
	return testutils.WebsocketURL(urlS) + "?client_id=" + clientID
}

func (ct *GameController) setupWebSocket(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	client := ct.AddClient(w, r, Player{Pseudo: "testName", ID: pass.EncryptedID(clientID)}, -1)
	if client != nil {
		client.StartLoop()
		client.WS.Close()
	}
}

func TestConcurrentEvents(t *testing.T) {
	// ProgressLogger.SetOutput(os.Stdout)

	ct := NewGameController("testGame", questionsList, GameOptions{4, 0, true}, nil) // do not start a game
	go ct.StartLoop()

	server := httptest.NewServer(http.HandlerFunc(ct.setupWebSocket))
	defer server.Close()

	client1, _, err := websocket.DefaultDialer.Dial(testutils.WebsocketURL(server.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	client2, _, err := websocket.DefaultDialer.Dial(testutils.WebsocketURL(server.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	client3, _, err := websocket.DefaultDialer.Dial(testutils.WebsocketURL(server.URL), nil)
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
	ct := NewGameController("testGame", questionsList, GameOptions{4, time.Millisecond * 50, true}, nil)
	go ct.StartLoop()

	server := httptest.NewServer(http.HandlerFunc(ct.setupWebSocket))
	defer server.Close()

	client, _, err := websocket.DefaultDialer.Dial(testutils.WebsocketURL(server.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		_, _, err = client.ReadMessage()
		if err != nil {
			panic(err)
		}
	}()

	ct.Game.EmitQuestion() // launch the timer

	time.Sleep(time.Millisecond * 200)
}

func TestClientInvalidMessage(t *testing.T) {
	WarningLogger.SetOutput(os.Stdout)
	ProgressLogger.SetOutput(os.Stdout)

	ct := NewGameController("testGame", questionsList, GameOptions{2, 0, true}, nil)
	go ct.StartLoop()

	server := httptest.NewServer(http.HandlerFunc(ct.setupWebSocket))
	defer server.Close()

	client, _, err := websocket.DefaultDialer.Dial(testutils.WebsocketURL(server.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	var ev game.StateUpdate
	err = client.ReadJSON(&ev) // player join
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ev)

	err = client.ReadJSON(&ev) // game lobby
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ev)

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

	ct := NewGameController("testGame", questionsList, GameOptions{2, 0, true}, nil)

	go ct.StartLoop()

	server := httptest.NewServer(http.HandlerFunc(ct.setupWebSocket))
	defer server.Close()

	client1, _, err := websocket.DefaultDialer.Dial(testutils.WebsocketURL(server.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	var events game.MaybeUpdate
	if err = client1.ReadJSON(&events); err != nil {
		t.Fatal(err)
	}

	ct.gameLock.Lock()
	if ct.Game.IsPlaying() {
		t.Fatal("game should not have started")
	}
	ct.gameLock.Unlock()

	client2, _, err := websocket.DefaultDialer.Dial(testutils.WebsocketURL(server.URL), nil)
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
	if !ct.Game.IsPlaying() {
		t.Fatal("game should have started")
	}
	ct.gameLock.Unlock()
}

func TestInvalidJoin(t *testing.T) {
	WarningLogger.SetOutput(io.Discard)

	ct := NewGameController("testGame", questionsList, GameOptions{1, 0, true}, nil)

	go ct.StartLoop()

	server := httptest.NewServer(http.HandlerFunc(ct.setupWebSocket))
	defer server.Close()

	// first client join and launch the game
	_, _, err := websocket.DefaultDialer.Dial(testutils.WebsocketURL(server.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	// second client try to join
	client2, _, err := websocket.DefaultDialer.Dial(testutils.WebsocketURL(server.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = client2.ReadMessage()
	if err == nil {
		t.Fatal("expected error when joining started game")
	}
}

func TestSummary(t *testing.T) {
	WarningLogger.SetOutput(io.Discard)

	notif := make(chan GameSummary)
	go func() {
		for {
			if sum, has := <-notif; has {
				fmt.Println(sum)
			}
		}
	}()

	ct := NewGameController("testGame", questionsList, GameOptions{2, 0, true}, notif)

	go ct.StartLoop()

	server := httptest.NewServer(http.HandlerFunc(ct.setupWebSocket))
	defer server.Close()

	client1, _, err := websocket.DefaultDialer.Dial(websocketURLWithClientID(t, server.URL, "client1"), nil)
	if err != nil {
		t.Fatal(err)
	}

	var events game.MaybeUpdate
	if err = client1.ReadJSON(&events); err != nil {
		t.Fatal(err)
	}

	if sum := ct.Summary(); len(sum.Successes) != 1 || sum.PlayerTurn != nil {
		t.Fatal(sum)
	}

	client2, _, err := websocket.DefaultDialer.Dial(testutils.WebsocketURL(server.URL), nil)
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

	if sum := ct.Summary(); len(sum.Successes) != 2 || sum.PlayerTurn == nil {
		t.Fatal(sum)
	}
}

func TestReview(t *testing.T) {
	WarningLogger.SetOutput(io.Discard)

	ct := NewGameController("testGame", questionsList, GameOptions{2, 0, true}, nil)

	go ct.StartLoop()

	server := httptest.NewServer(http.HandlerFunc(ct.setupWebSocket))
	defer server.Close()

	client1, _, err := websocket.DefaultDialer.Dial(websocketURLWithClientID(t, server.URL, "client1"), nil)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Millisecond * 20)

	client2, _, err := websocket.DefaultDialer.Dial(testutils.WebsocketURL(server.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	_, _, _ = client1.ReadMessage() // player joined
	_, _, _ = client2.ReadMessage() // player joined
	_, _, _ = client1.ReadMessage() // player joined
	_, _, _ = client2.ReadMessage() // player joined

	time.Sleep(time.Millisecond * 20)

	_, _, _ = client1.ReadMessage() // game start & playerTurn

	err = client1.WriteJSON(game.ClientEvent{Event: game.DiceClicked{}, Player: 0})
	if err != nil {
		t.Fatal(err)
	}

	var update game.StateUpdate
	if err = client1.ReadJSON(&update); err != nil { // dice throw and possibleMoves
		t.Fatal(err)
	}
	choosenTile := update.Events[1].(game.PossibleMoves).Tiles[0]

	err = client1.WriteJSON(game.ClientEvent{Event: game.ClientMove{Tile: choosenTile}, Player: 0})
	if err != nil {
		t.Fatal(err)
	}

	if err = client1.ReadJSON(&update); err != nil { // move and showQuestion
		t.Fatal(err)
	}

	err = client1.WriteJSON(game.ClientEvent{Event: game.Answer{Answer: exClient.QuestionAnswersIn{}}, Player: 0})
	if err != nil {
		t.Fatal(err)
	}
	err = client2.WriteJSON(game.ClientEvent{Event: game.Answer{Answer: exClient.QuestionAnswersIn{}}, Player: 0})
	if err != nil {
		t.Fatal(err)
	}

	if err = client1.ReadJSON(&update); err != nil { // playerAnswerResult 1 & 2
		t.Fatal(err)
	}

	err = client1.WriteJSON(game.ClientEvent{Event: game.WantNextTurn{MarkQuestion: true}, Player: 0})
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Millisecond * 20)

	ct.Terminate <- true

	time.Sleep(time.Millisecond * 20)

	if err = client1.ReadJSON(&update); err != nil { // playerAnswerResult 1 & 2
		t.Fatal(err)
	}
	if _, isTerm := update.Events[0].(game.GameTerminated); !isTerm {
		t.Fatal("unexepected events", update.Events)
	}

	history := ct.review().QuestionHistory[Player{ID: "client1", Pseudo: "testName"}]
	if len(history.MarkedQuestions) != 1 || len(history.QuestionHistory) != 1 {
		t.Fatal(history)
	}
}
