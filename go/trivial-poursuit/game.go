package trivialpoursuit

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
	"github.com/gorilla/websocket"
)

var (
	WarningLogger  = log.New(os.Stdout, "trivial-poursuit:ERROR:", log.LstdFlags)
	ProgressLogger = log.New(io.Discard, "trivial-poursuit:INFO:", log.LstdFlags)
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// RegisterAndStart starts a new game at the given `apiPath`,
// registering the route.
func RegisterAndStart(apiPath string, options GameOptions) {
	ct := newGameController(options)
	go ct.startLoop()

	http.HandleFunc(apiPath, ct.setupWebSocket)
}

func (ct *gameController) setupWebSocket(w http.ResponseWriter, r *http.Request) {
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

	client.startLoop()
}

type client struct {
	conn *websocket.Conn
	game *gameController // to accept user events
}

func (cl *client) sendEvent(er game.EventList) error { return cl.conn.WriteJSON(er) }

// startLoop listens for new messages being sent to our WebSocket
// endpoint, only returning on error
// the connection is not closed yet
func (cl *client) startLoop() {
	defer func() {
		cl.game.leave <- cl
	}()

	for {
		// read in a message
		_, r, err := cl.conn.NextReader()
		if err != nil {
			// fmt.Errorf("client connection closed: %s", err)
			return
		}

		var event game.ClientEvent
		err = json.NewDecoder(r).Decode(&event)
		if err != nil {
			WarningLogger.Printf("invalid event format: %s", err)

			// return an error to the client and close
			message := websocket.FormatCloseMessage(websocket.CloseUnsupportedData, err.Error())
			cl.conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))

			return
		}

		// the player is deduced from the client pointer
		event.Player = cl.game.clients[cl]

		// process the event
		cl.game.incomingEvents <- event
	}
}

type GameOptions struct {
	PlayersNumber   int
	QuestionTimeout time.Duration
}

// gameController handle one game session
type gameController struct {
	join, leave chan *client

	incomingEvents  chan game.ClientEvent
	broadcastEvents chan []game.GameEvents
	clients         map[*client]game.PlayerID // current clients in the game

	game game.Game // game logic

	options GameOptions
}

func newGameController(options GameOptions) *gameController {
	return &gameController{
		join:            make(chan *client),
		leave:           make(chan *client),
		incomingEvents:  make(chan game.ClientEvent),
		broadcastEvents: make(chan []game.GameEvents, 1), // the main loop write in this channel
		clients:         map[*client]game.PlayerID{},
		game:            *game.NewGame(options.QuestionTimeout),
		options:         options,
	}
}

func (gc *gameController) startLoop() {
	for {
		select {
		case event := <-gc.broadcastEvents:
			ProgressLogger.Println("Broadcasting...")
			for client, clientID := range gc.clients {
				err := client.sendEvent(event)
				if err != nil {
					WarningLogger.Printf("Broadcasting to client %d failed: %s", clientID, err)
				}
			}

		case client := <-gc.leave:
			ProgressLogger.Println("Removing player...")

			event := gc.game.RemovePlayer(gc.clients[client])
			delete(gc.clients, client)

			if gc.game.NumberPlayers() == 0 { // reset the game
				// TODO: higher level gestion of multiples games
				gc.game = *game.NewGame(gc.options.QuestionTimeout)
			} else if !gc.game.HasStarted() { // update the lobby
				gc.broadcastEvents <- []game.GameEvents{{
					Events: game.Events{event},
					State:  gc.game.GameState,
				}}
			}

		case <-gc.game.QuestionTimeout.C:
			ProgressLogger.Println("QuestionTimeoutAction...")

			events := gc.game.QuestionTimeoutAction()
			gc.broadcastEvents <- events

		case client := <-gc.join:
			ProgressLogger.Println("Adding player...")

			// we do not allow connection into an already started game
			if gc.game.HasStarted() { // TODO: notify the client
				continue
			}

			event := gc.game.AddPlayer()
			gc.clients[client] = event.Player

			// only notifie the player who joined ...
			client.sendEvent(game.EventList{{
				Events: game.Events{game.PlayerJoin{Player: event.Player}},
				State:  gc.game.GameState,
			}})

			// ... check if the new player triggers a game start
			if gc.game.NumberPlayers() >= gc.options.PlayersNumber {
				events := gc.game.StartGame()
				gc.broadcastEvents <- []game.GameEvents{events}
			} else { // update the lobby
				gc.broadcastEvents <- []game.GameEvents{{
					Events: game.Events{event},
					State:  gc.game.GameState,
				}}
			}

		case message := <-gc.incomingEvents:
			ProgressLogger.Println("HandleClientEvent...")

			out, err := gc.game.HandleClientEvent(message)
			if err != nil { // malicious client: ignore the query
				WarningLogger.Println(err)
			}

			if len(out) != 0 {
				gc.broadcastEvents <- out
			}

		}
	}
}
