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

	"github.com/benoitkugler/maths-online/pass"
	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
	"github.com/benoitkugler/maths-online/utils"
	"github.com/gorilla/websocket"
)

var (
	WarningLogger  = log.New(os.Stdout, "trivial-poursuit-game:ERROR:", log.LstdFlags)
	ProgressLogger = log.New(io.Discard, "trivial-poursuit-game:INFO:", log.LstdFlags)
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// GameID is an in-memory identifier for a game room.
type GameID = string

// AddClient uses the given connection to start a web socket, registred
// with given `player`.
// Errors are send to the websocket, and the function blocks until the game ends
func (ct *GameController) AddClient(w http.ResponseWriter, r *http.Request, player Player) {
	// upgrade this connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		WarningLogger.Println("Failed to init websocket: ", err)
		return
	}
	defer ws.Close()

	client := &client{conn: ws, game: ct, isAccepted: make(chan bool), player: player}
	ct.join <- client

	isAccepted := <-client.isAccepted // wait for the controller to check the access
	if !isAccepted {
		// the game at this end point is not usable: close the connection with an error
		utils.WebsocketError(ws, errors.New("game is closed"))
		return
	}

	client.startLoop()
}

type client struct {
	conn *websocket.Conn
	game *GameController // to accept user events

	isAccepted chan bool // valid the access to the game

	player Player
}

func (cl *client) sendEvent(er game.StateUpdate) error { return cl.conn.WriteJSON(er) }

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

		if err, isClose := err.(*websocket.CloseError); isClose {
			ProgressLogger.Printf("Client left (%v)", err)
			return
		}

		if err != nil {
			WarningLogger.Printf("unexpected client error: %s", err)
			return
		}

		var event game.ClientEvent
		err = json.NewDecoder(r).Decode(&event)
		if err != nil {
			WarningLogger.Printf("invalid event format: %s", err)

			// return an error to the client and close
			utils.WebsocketError(cl.conn, err)

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
	ShowDecrassage  bool
}

// GameController handles one game room
type GameController struct {
	ID GameID

	// Terminate may be used to cleanly exit the game,
	// noticing clients and exiting the main goroutine.
	// It is however not considered as a normal exit,
	// so that the review is not emitted.
	Terminate chan bool

	monitor     chan GameSummary
	join, leave chan *client

	incomingEvents  chan game.ClientEvent
	broadcastEvents chan game.StateUpdate
	clients         map[*client]game.PlayerID // current clients in the game

	game     game.Game // game logic
	gameLock sync.Mutex

	Options GameOptions
}

// NewGameController creates a new game, with given `id` and `options`.
// `monitor` is an optionnal channel to write back the main progress of the game.
func NewGameController(id GameID, questions game.QuestionPool, options GameOptions, monitor chan GameSummary) *GameController {
	return &GameController{
		ID:              id,
		monitor:         monitor,
		Terminate:       make(chan bool),
		join:            make(chan *client, 1),
		leave:           make(chan *client),
		incomingEvents:  make(chan game.ClientEvent),
		broadcastEvents: make(chan game.StateUpdate, 1), // the main loop write in this channel
		clients:         map[*client]game.PlayerID{},
		game:            *game.NewGame(options.QuestionTimeout, options.ShowDecrassage, questions),
		Options:         options,
	}
}

func (gc *GameController) playerIDsToClients() map[game.PlayerID]*client {
	players := make(map[game.PlayerID]*client)
	for k, v := range gc.clients {
		players[v] = k
	}
	return players
}

// StartLoop starts the main game loop.
// The function blocks until the game is over,
// and then returns the game review.
// It returns false if the game ended abnormally, due to forced termination or all players leaving
func (gc *GameController) StartLoop() (Review, bool) {
	var isGameOver bool // if true, broadcast the last events and quit
	for {
		select {
		case <-gc.Terminate:
			ProgressLogger.Println("Terminating game")

			for client, clientID := range gc.clients {
				err := client.sendEvent(game.StateUpdate{
					Events: game.Events{
						game.GameTerminated{},
					},
					State: gc.game.GameState,
				})
				if err != nil {
					WarningLogger.Printf("Broadcasting to client %d failed: %s", clientID, err)
				}
			}

			return Review{}, false

		case event := <-gc.broadcastEvents:
			ProgressLogger.Println("Broadcasting...")
			for client, clientID := range gc.clients {
				err := client.sendEvent(event)
				if err != nil {
					WarningLogger.Printf("Broadcasting to client %d failed: %s", clientID, err)
				}
			}

			if gc.monitor != nil { // notify the monitor
				gc.monitor <- gc.Summary()
			}

			if isGameOver {
				ProgressLogger.Println("Game is over: exitting game loop.")
				return gc.review(), true
			}

		case client := <-gc.leave:
			if _, in := gc.clients[client]; !in { // client who never joined may still end up here
				continue
			}

			ProgressLogger.Printf("Removing player %d...", gc.clients[client])

			gc.gameLock.Lock()
			event := gc.game.RemovePlayer(gc.clients[client])
			hasStarted := gc.game.HasStarted()
			nbPlayers := gc.game.NumberPlayers()
			delete(gc.clients, client)
			gc.gameLock.Unlock()

			if gc.monitor != nil { // notify the monitor
				gc.monitor <- gc.Summary()
			}

			// end the game only if the game has already started and all
			// players have left
			if hasStarted && nbPlayers == 0 {
				// we consider all player leaving early means the game
				// did not end properly
				// also, game.Players would be empty here
				return Review{}, false
			} else if !hasStarted { // update the lobby if the game has to started
				gc.broadcastEvents <- game.StateUpdate{
					Events: game.Events{event},
					State:  gc.game.GameState,
				}
			}

		case <-gc.game.QuestionTimeout.C:
			ProgressLogger.Println("QuestionTimeoutAction...")

			events := gc.game.QuestionTimeoutAction()
			if events != nil {
				gc.broadcastEvents <- *events
			}

		case client := <-gc.join:
			ProgressLogger.Println("Adding player...")

			// we do not allow connection into an already started game
			if gc.game.HasStarted() {
				// the game at this end point is not usable: close the connection with an error
				client.isAccepted <- false
				continue
			} else {
				client.isAccepted <- true
			}

			event := gc.game.AddPlayer(client.player.Name)
			gc.clients[client] = event.Player

			// only notifie the player who joined ...
			client.sendEvent(game.StateUpdate{
				Events: game.Events{game.PlayerJoin{Player: event.Player}},
				State:  gc.game.GameState,
			})

			// ... check if the new player triggers a game start
			if gc.game.NumberPlayers() >= gc.Options.PlayersNumber {
				events := gc.game.StartGame()
				gc.broadcastEvents <- events
			} else { // update the lobby
				gc.broadcastEvents <- game.StateUpdate{
					Events: game.Events{event},
					State:  gc.game.GameState,
				}
			}

			if gc.monitor != nil { // notify the monitor
				gc.monitor <- gc.Summary()
			}

		case message := <-gc.incomingEvents:
			ProgressLogger.Println("HandleClientEvent...")

			var (
				events game.MaybeUpdate
				err    error
			)
			events, isGameOver, err = gc.game.HandleClientEvent(message)
			if err != nil { // malicious client: ignore the query
				WarningLogger.Println(err)
				continue
			}

			if events != nil {
				gc.broadcastEvents <- *events
			}

		}
	}
}

type Player struct {
	ID   pass.EncryptedID
	Name string // used for anonymous players
}

// GameSummary is emitted back to the teacher monitor,
// and provides an high level overview of one game.
type GameSummary struct {
	PlayerTurn *Player // nil before game start
	Successes  map[Player]game.Success
	ID         GameID
	RoomSize   int // number of player expected
}

// Summary locks and returns the current game summary.
func (gc *GameController) Summary() GameSummary {
	gc.gameLock.Lock()
	defer gc.gameLock.Unlock()

	state := gc.game.GameState
	players := gc.playerIDsToClients()

	successes := make(map[Player]game.Success)
	for k, v := range state.Players {
		successes[players[k].player] = v.Success
	}
	out := GameSummary{
		ID:        gc.ID,
		Successes: successes,
		RoomSize:  gc.Options.PlayersNumber,
	}
	if id := state.Player; id != -1 {
		out.PlayerTurn = &players[id].player
	}

	return out
}

// Review contains the information at the end of a game room,
// and should be used to persit information over sessions.
type Review struct {
	QuestionHistory map[Player]game.QuestionReview
	ID              GameID
}

// return the current game review
func (gc *GameController) review() Review {
	gc.gameLock.Lock()
	defer gc.gameLock.Unlock()

	out := Review{
		ID:              gc.ID,
		QuestionHistory: make(map[Player]game.QuestionReview),
	}

	players := gc.playerIDsToClients()
	for k, v := range gc.game.Players {
		out.QuestionHistory[players[k].player] = v.Review
	}
	return out
}
