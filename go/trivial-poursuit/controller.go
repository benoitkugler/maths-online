package trivialpoursuit

import (
	"fmt"
	"log"
	"net/http"

	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// SetupRoutes start the controllers and registers the end points
func SetupRoutes(apiPath string) {
	ct := newGameController()
	go ct.startLoop()

	http.HandleFunc(apiPath, ct.setupWebSocket)
}

func (ct *controller) setupWebSocket(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to init websocket: ", err)
		return
	}

	defer ws.Close()

	// err = ws.WriteMessage(websocket.TextMessage, []byte("Hi Client!"))
	// if err != nil {
	// 	log.Println("failed to greet", err)
	// 	return
	// }

	// create a client object ...
	client := &client{conn: ws, game: ct}
	// ... and adds it to the current game
	ct.join <- client

	err = client.startLoop()
	if err != nil {
		log.Println("error reading client:", err)
		ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
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
			return fmt.Errorf("invalid client format: %s", err)
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

		if timeout := gc.game.QuestionTimeout; timeout != nil {
			if _, has := <-timeout.C; has { // force the question ending
				fmt.Println("QuestionTimeoutAction")

				events := gc.game.QuestionTimeoutAction()
				gc.broadcastEvents <- events
				continue
			}
		}

		select {
		case client := <-gc.join:
			fmt.Println("AddPlayer")

			playerID := gc.game.AddPlayer()
			gc.clients[client] = playerID
			fmt.Println(gc.clients)
		case client := <-gc.leave:
			fmt.Println("RemovePlayer")

			gc.game.RemovePlayer(gc.clients[client])
			delete(gc.clients, client)
		case message := <-gc.incomingEvents:
			fmt.Println("HandleClientEvent")

			out, err := gc.game.HandleClientEvent(message)
			if err != nil { // malicious client: ignore the query
				log.Println(err)
			}
			gc.broadcastEvents <- out

		case event := <-gc.broadcastEvents:
			fmt.Println("Sending events to all clients")

			for client, clientID := range gc.clients {
				err := client.sendEvent(event)
				if err != nil {
					log.Printf("Broadcasting to client %d: %s", clientID, err)
				}
			}

		}
	}
}
