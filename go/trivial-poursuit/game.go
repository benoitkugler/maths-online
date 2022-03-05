package trivialpoursuit

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
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

// RegisterTestGame starts a new game at the given `apiPath`,
// and registers the route.
func RegisterTestGame(apiPath string, options GameOptions) {
	ct := newGameController(options)
	go func() {
		for {
			ct.startLoop()

			// when game ends, just reset it and start again
			ct.game = *game.NewGame(options.QuestionTimeout)
		}
	}()

	http.HandleFunc(apiPath, ct.setupWebSocket)
}

func websocketError(ws *websocket.Conn, err error) {
	message := websocket.FormatCloseMessage(websocket.CloseUnsupportedData, err.Error())
	ws.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
}

func (ct *gameController) setupWebSocket(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		WarningLogger.Println("Failed to init websocket: ", err)
		return
	}
	defer ws.Close()

	client, hasJoined := ct.tryJoin(ws)
	if !hasJoined { // the game at this end point is not usable: close the connection with an error
		websocketError(ws, errors.New("game is closed"))
		return
	}

	// all good, start listening for client messages
	client.startLoop()
}

type client struct {
	conn *websocket.Conn
	game *gameController // to accept user events
}

func (cl *client) sendEvent(er game.StateUpdates) error { return cl.conn.WriteJSON(er) }

// startLoop listens for new messages being sent to our WebSocket
// endpoint, only returning on error.
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
			websocketError(cl.conn, err)

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
	broadcastEvents chan []game.StateUpdate
	clients         map[*client]game.PlayerID // current clients in the game

	game     game.Game // game logic
	gameLock sync.Mutex

	options GameOptions
}

func newGameController(options GameOptions) *gameController {
	return &gameController{
		join:            make(chan *client, 1),
		leave:           make(chan *client),
		incomingEvents:  make(chan game.ClientEvent),
		broadcastEvents: make(chan []game.StateUpdate, 1), // the main loop write in this channel
		clients:         map[*client]game.PlayerID{},
		game:            *game.NewGame(options.QuestionTimeout),
		options:         options,
	}
}

// check if the game can accept a new player
// if so, also create a new client instance and sends it
// to the join channel
func (gc *gameController) tryJoin(ws *websocket.Conn) (*client, bool) {
	gc.gameLock.Lock()
	defer gc.gameLock.Unlock()

	// we do not allow connection into an already started game
	if gc.game.IsPlaying() {
		return nil, false
	}

	// create a client object ...
	cl := &client{conn: ws, game: gc}
	// ... and adds it to the current game
	gc.join <- cl

	return cl, true
}

func (gc *gameController) startLoop() {
	var isGameOver bool // if true, broadcast the last events and quit
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

			if isGameOver {
				ProgressLogger.Println("Game is over: exitting game loop.")
				return
			}

		case client := <-gc.leave:
			ProgressLogger.Println("Removing player...")

			gc.gameLock.Lock()
			event := gc.game.RemovePlayer(gc.clients[client])
			hasStarted := gc.game.IsPlaying()
			nbPlayers := gc.game.NumberPlayers()
			gc.gameLock.Unlock()

			delete(gc.clients, client)

			// end the game only if the game has already started and all
			// players have left
			if hasStarted && nbPlayers == 0 {
				return
			} else if !hasStarted { // update the lobby if the game has to started
				gc.broadcastEvents <- []game.StateUpdate{{
					Events: game.Events{event},
					State:  gc.game.GameState,
				}}
			}

		case <-gc.game.QuestionTimeout.C:
			ProgressLogger.Println("QuestionTimeoutAction...")

			var events game.StateUpdates
			events, isGameOver = gc.game.QuestionTimeoutAction()
			gc.broadcastEvents <- events

		case client := <-gc.join:
			ProgressLogger.Println("Adding player...")

			event := gc.game.AddPlayer()
			gc.clients[client] = event.Player

			// only notifie the player who joined ...
			client.sendEvent(game.StateUpdates{{
				Events: game.Events{game.PlayerJoin{Player: event.Player}},
				State:  gc.game.GameState,
			}})

			// ... check if the new player triggers a game start
			if gc.game.NumberPlayers() >= gc.options.PlayersNumber {
				events := gc.game.StartGame()
				gc.broadcastEvents <- []game.StateUpdate{events}
			} else { // update the lobby
				gc.broadcastEvents <- []game.StateUpdate{{
					Events: game.Events{event},
					State:  gc.game.GameState,
				}}
			}

		case message := <-gc.incomingEvents:
			ProgressLogger.Println("HandleClientEvent...")

			var (
				events game.StateUpdates
				err    error
			)
			events, isGameOver, err = gc.game.HandleClientEvent(message)
			if err != nil { // malicious client: ignore the query
				WarningLogger.Println(err)
				continue
			}

			if len(events) != 0 {
				gc.broadcastEvents <- events
			}

		}
	}
}
