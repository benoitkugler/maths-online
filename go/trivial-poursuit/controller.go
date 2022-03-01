package trivialpoursuit

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
	"github.com/gorilla/websocket"
)

var WarningLogger = log.New(os.Stdout, "trivial-poursuit:", log.LstdFlags)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// RegisterAndStart start the controllers and registers the end points
func RegisterAndStart(apiPath string) {
	ct := newGameController()
	go ct.startLoop()

	http.HandleFunc(apiPath, ct.setupWebSocket)
}

func (ct *controller) setupWebSocket(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		WarningLogger.Println("Failed to init websocket: ", err)
		return
	}

	defer ws.Close()

	// create a client object ...
	client := &client{conn: ws, game: ct}
	// ... and adds it to the current game
	ct.join <- client

	err = client.startLoop()
	if err != nil {
		WarningLogger.Println(err)
		ws.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseUnsupportedData, err.Error()),
			time.Now().Add(time.Second))
	}
}

type client struct {
	conn *websocket.Conn
	game *controller // to accept user events
}

func (cl *client) sendEvent(er game.GameEvents) error { return cl.conn.WriteJSON(er) }

// startLoop listens for new messages being sent to our WebSocket
// endpoint, only returning on error
// the connection is not closed yet
func (cl *client) startLoop() error {
	defer func() {
		cl.game.leave <- cl
	}()

	for {
		// read in a message
		var event game.ClientEvent
		err := cl.conn.ReadJSON(&event)
		if err != nil {
			return fmt.Errorf("invalid event: %s", err)
		}

		// the player is deduced from the client pointer
		event.Player = cl.game.clients[cl]

		// process the event
		cl.game.incomingEvents <- event
	}
}

type controller struct {
	join, leave chan *client

	incomingEvents  chan game.ClientEvent
	broadcastEvents chan game.GameEvents
	clients         map[*client]game.PlayerID // current clients in the game

	game game.Game // game logic
}

func newGameController() *controller {
	return &controller{
		join:            make(chan *client),
		leave:           make(chan *client),
		incomingEvents:  make(chan game.ClientEvent),
		broadcastEvents: make(chan game.GameEvents, 1), // the main loop write in this channel
		clients:         map[*client]game.PlayerID{},
		game:            *game.NewGame(),
	}
}

func (gc *controller) startLoop() {
	for {
		select {
		case <-gc.game.QuestionTimeout.C:
			events := gc.game.QuestionTimeoutAction()
			gc.broadcastEvents <- events
		case client := <-gc.join:
			playerID := gc.game.AddPlayer()
			gc.clients[client] = playerID
		case client := <-gc.leave:
			gc.game.RemovePlayer(gc.clients[client])
			delete(gc.clients, client)
		case message := <-gc.incomingEvents:
			out, err := gc.game.HandleClientEvent(message)
			if err != nil { // malicious client: ignore the query
				WarningLogger.Println(err)
			}

			if !out.IsEmpty() {
				gc.broadcastEvents <- out
			}

		case event := <-gc.broadcastEvents:
			for client, clientID := range gc.clients {
				err := client.sendEvent(event)
				if err != nil {
					WarningLogger.Printf("Broadcasting to client %d failed: %s", clientID, err)
				}
			}

		}
	}
}
