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
	WarningLogger  = log.New(os.Stdout, "tv-game:ERROR:", 0)
	ProgressLogger = log.New(io.Discard, "tv-game:INFO:", 0)
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// GameID is the full identifier of a game room,
// usually of the form <sessionID><gameID> (excepted for demonstration games).
type GameID string

// AddClient uses the given connection to start a web socket, registred
// with given `player`.
// Errors are send to the websocket, and the function blocks until the game ends
func (ct *GameController) AddClient(w http.ResponseWriter, r *http.Request, player Player, playerID game.PlayerSerial) *Client {
	// upgrade this connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		WarningLogger.Println("Failed to init websocket: ", err)
		return nil
	}

	client := &Client{
		WS:         ws,
		game:       ct,
		isAccepted: make(chan bool),
		player:     player,
		PlayerID:   playerID,
	}

	ct.join <- client

	// wait for the controller to check the access
	isAccepted := <-client.isAccepted
	if !isAccepted {
		ProgressLogger.Printf("Rejecting connection at %s", ct.ID)
		// the game at this end point is not usable: close the connection with an error
		utils.WebsocketError(ws, errors.New("game is closed"))
		ws.Close()
		return nil
	}

	return client
}

type Client struct {
	// WS should be close when StartLoop ends
	WS   *websocket.Conn
	game *GameController // to accept user events

	isAccepted chan bool // valid the access to the game

	player   Player
	PlayerID game.PlayerSerial // used to handle reconnection
}

func (cl *Client) sendEvent(er game.StateUpdate) error { return cl.WS.WriteJSON(er) }

// StartLoop listens for new messages being sent to our WebSocket
// endpoint, only returning on error.
// the connection is not closed yet
func (cl *Client) StartLoop() {
	defer func() {
		cl.game.leave <- cl
	}()

	for {

		// read in a message
		_, r, err := cl.WS.NextReader()

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
			utils.WebsocketError(cl.WS, err)

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
	join, leave chan *Client

	incomingEvents chan game.ClientEvent
	clients        map[*Client]game.PlayerSerial // current clients in the game

	Game     game.Game // game logic
	gameLock sync.Mutex

	Options GameOptions
}

// NewGameController creates a new game, with given `id` and `options`.
// `monitor` is an optionnal channel to write back the main progress of the game.
func NewGameController(id GameID, questions game.QuestionPool, options GameOptions, monitor chan GameSummary) *GameController {
	return &GameController{
		ID:             id,
		monitor:        monitor,
		Terminate:      make(chan bool),
		join:           make(chan *Client, 1),
		leave:          make(chan *Client),
		incomingEvents: make(chan game.ClientEvent),
		clients:        map[*Client]game.PlayerSerial{},
		Game:           game.NewGame(options.QuestionTimeout, options.ShowDecrassage, questions),
		Options:        options,
	}
}

func (gc *GameController) playerIDsToClients() map[game.PlayerSerial]*Client {
	players := make(map[game.PlayerSerial]*Client)
	for k, v := range gc.clients {
		players[v] = k
	}
	return players
}

func (gc *GameController) broadcastEvents(events game.StateUpdate) {
	ProgressLogger.Printf("Game %s : broadcasting...", gc.ID)

	for client, clientID := range gc.clients {
		err := client.sendEvent(events)
		if err != nil {
			WarningLogger.Printf("Broadcasting to client %d failed: %s", clientID, err)
		}
	}

	if gc.monitor != nil { // notify the monitor
		gc.monitor <- gc.Summary()
	}
}

// StartLoop starts the main game loop.
// The function blocks until the game is over,
// and then returns the game review.
// It returns false if the game ended abnormally, due to forced termination.
func (gc *GameController) StartLoop() (Review, bool) {
	for {
		select {
		case <-gc.Terminate:
			ProgressLogger.Println("Terminating game", gc.ID)

			for client, clientID := range gc.clients {
				err := client.sendEvent(game.StateUpdate{
					Events: game.Events{
						game.GameTerminated{},
					},
					State: gc.Game.GameState,
				})
				if err != nil {
					WarningLogger.Printf("Broadcasting to client %d failed: %s", clientID, err)
				}
			}

			return Review{}, false
		case client := <-gc.leave:
			if _, in := gc.clients[client]; !in { // client who never joined may still end up here
				continue
			}

			ProgressLogger.Printf("Game %s : removing player %d...", gc.ID, gc.clients[client])

			playerID := gc.clients[client]

			// check if the player is not already removed
			if gc.Game.Players[playerID] == nil {
				continue
			}

			gc.gameLock.Lock()
			update := gc.Game.RemovePlayer(playerID)
			gc.gameLock.Unlock()

			delete(gc.clients, client)

			if gc.monitor != nil { // notify the monitor
				gc.monitor <- gc.Summary()
			}

			gc.broadcastEvents(update)

		case <-gc.Game.QuestionTimeout.C:
			ProgressLogger.Printf("Game %s : questionTimeoutAction...", gc.ID)

			events := gc.Game.QuestionTimeoutAction()
			if events != nil {
				gc.broadcastEvents(*events)
			}

		case client := <-gc.join:
			// Here are the following cases :
			//	- totally fresh connection
			//	- reconnection in an already started game
			//	- reconnection (as detected at the request level) in a waiting game :
			//		for simplicity we considered this a new connection

			gc.gameLock.Lock()
			reconnection := client.PlayerID != -1 && gc.Game.Players[client.PlayerID] != nil && gc.Game.HasStarted()
			gc.gameLock.Unlock()

			if !reconnection { // fresh connection
				if gc.Game.HasStarted() {
					// we do not allow fresh connection into an already started game
					client.isAccepted <- false
					continue
				}
				ProgressLogger.Printf("Game %s : adding player...", gc.ID)

				playerID, event := gc.Game.AddPlayer(client.player.Pseudo)
				// register the playerID so that it can be send back
				client.PlayerID = playerID
				gc.clients[client] = playerID

				client.isAccepted <- true

				// only notifie the player who joined ...
				client.sendEvent(game.StateUpdate{
					Events: game.Events{game.PlayerJoin{Player: playerID}},
					State:  gc.Game.GameState,
				})

				// ... check if the new player triggers a game start
				if gc.Game.NumberPlayers(true) >= gc.Options.PlayersNumber {
					ProgressLogger.Printf("Game %s : starting", gc.ID)

					events := gc.Game.StartGame()
					gc.broadcastEvents(events)
				} else { // update the lobby
					gc.broadcastEvents(game.StateUpdate{
						Events: game.Events{event},
						State:  gc.Game.GameState,
					})
				}
			} else { // reconnection
				ProgressLogger.Printf("Game %s : reconnecting player %d...", gc.ID, client.PlayerID)

				gc.clients[client] = client.PlayerID // register the new client connection
				client.isAccepted <- true
				gc.gameLock.Lock()
				event := gc.Game.ReconnectPlayer(client.PlayerID, client.player.Pseudo)
				gc.gameLock.Unlock()

				gc.broadcastEvents(game.StateUpdate{
					Events: game.Events{event},
					State:  gc.Game.GameState,
				})
			}

		case message := <-gc.incomingEvents:
			ProgressLogger.Printf("Game %s : handleClientEvent...", gc.ID)

			events, isGameOver, err := gc.Game.HandleClientEvent(message)
			if err != nil { // malicious client: ignore the query
				WarningLogger.Println(err)
				continue
			}

			if events != nil {
				gc.broadcastEvents(*events)
			}

			if isGameOver {
				ProgressLogger.Printf("Game %s is over: exitting game loop.", gc.ID)
				return gc.review(), true
			}
		}
	}
}

type Player struct {
	ID     pass.EncryptedID
	Pseudo string // used for anonymous players
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

	state := gc.Game.GameState
	players := gc.playerIDsToClients()

	successes := make(map[Player]game.Success)
	for k, v := range state.Players {
		client := players[k]
		if client == nil { // inactive player
			continue
		}
		successes[players[k].player] = v.Success
	}
	out := GameSummary{
		ID:        gc.ID,
		Successes: successes,
		RoomSize:  gc.Options.PlayersNumber,
	}
	if id := state.Player; id != -1 {
		if pl := players[id]; pl != nil {
			out.PlayerTurn = &pl.player
		}
	}

	return out
}

// Review contains the information at the end of a game room,
// and should be used to persist information over sessions.
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
	for k, v := range gc.Game.Players {
		if players[k] == nil { // player not connected anymore
			continue
		}
		out.QuestionHistory[players[k].player] = v.Review
	}
	return out
}
