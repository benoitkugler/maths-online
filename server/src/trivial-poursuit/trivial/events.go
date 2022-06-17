package trivial

import (
	"log"
	"os"

	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
)

var (
	WarningLogger  = log.New(os.Stdout, "tv-game:ERROR:", 0)
	ProgressLogger = log.New(os.Stdout, "tv-game:INFO :", 0)
)

func (r *Room) broadcastEvents(events game.StateUpdate) {
	ProgressLogger.Printf("Game %s : broadcasting...", r.ID)

	for client, conn := range r.currentPlayers {
		if conn == nil { // ignore disconnected players
			continue
		}
		err := conn.WriteJSON(events)
		if err != nil {
			WarningLogger.Printf("Broadcasting to client %d failed: %s", client.ID, err)
		}
	}
}

// Listen starts the main game loop, listening
// on the game channels.
// It returns `true` when the game is over, or `false` when it
// is manually terminated.
// Care should be taken to make sure no more events are send on the
// channels when this method has returned.
func (r *Room) Listen() (replay Replay, naturalEnding bool) {
	for {
		select {
		case <-r.Terminate:
			r.onTerminate()
			return Replay{}, false
		case client := <-r.Leave:
			r.onLeave(client)
		case <-r.game.QuestionTimeout.C:
			r.onQuestionTimeout()
		case message := <-r.Event:
			isGameOver := r.onEvent(message)
			if isGameOver {
				return r.review(), true
			}
		}
	}
}

func (r *Room) onTerminate() {
	r.lock.Lock()
	defer r.lock.Unlock()

	ProgressLogger.Println("Terminating game", r.ID)

	// for client, clientID := range r.currentPlayers {
	// 	err := client.sendEvent(game.StateUpdate{
	// 		Events: game.Events{
	// 			game.GameTerminated{},
	// 		},
	// 		State: r.Game.GameState,
	// 	})
	// 	if err != nil {
	// 		WarningLogger.Printf("Broadcasting to client %d failed: %s", clientID, err)
	// 	}
	// }
}

// Join should used on a new connection, on one the of the following cases:
//	- totally fresh connection
//	- reconnection in an already started game
//	- reconnection in a waiting game
// It is safe for concurrent uses.
// It returns an error for instance if the game has already started.
func (r *Room) Join(player Player, connection Connection) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	// Here are the following cases :
	//	- totally fresh connection
	//	- reconnection in an already started game
	//	- reconnection (as detected at the request level) in a waiting game :
	//		for simplicity we considered this a new connection

	// r.gameLock.Lock()
	// reconnection := client.PlayerID != -1 && r.Game.Players[client.PlayerID] != nil && r.Game.HasStarted()
	// r.gameLock.Unlock()

	// if !reconnection { // fresh connection
	// 	if r.Game.HasStarted() {
	// 		// we do not allow fresh connection into an already started game
	// 		client.isAccepted <- false
	// 		continue
	// 	}
	// 	ProgressLogger.Printf("Game %s : adding player...", r.ID)

	// 	playerID, event := r.Game.AddPlayer(client.player.Pseudo)
	// 	// register the playerID so that it can be send back
	// 	client.PlayerID = playerID
	// 	r.currentPlayers[client] = playerID

	// 	client.isAccepted <- true

	// 	// only notifie the player who joined ...
	// 	client.sendEvent(game.StateUpdate{
	// 		Events: game.Events{game.PlayerJoin{Player: playerID}},
	// 		State:  r.Game.GameState,
	// 	})

	// 	// ... check if the new player triggers a game start
	// 	if r.Game.NumberPlayers(true) >= r.Options.PlayersNumber {
	// 		ProgressLogger.Printf("Game %s : starting", r.ID)

	// 		events := r.Game.StartGame()
	// 		r.broadcastEvents(events)
	// 	} else { // update the lobby
	// 		r.broadcastEvents(game.StateUpdate{
	// 			Events: game.Events{event},
	// 			State:  r.Game.GameState,
	// 		})
	// 	}
	// } else { // reconnection
	// 	ProgressLogger.Printf("Game %s : reconnecting player %d...", r.ID, client.PlayerID)

	// 	r.currentPlayers[client] = client.PlayerID // register the new client connection
	// 	client.isAccepted <- true
	// 	r.gameLock.Lock()
	// 	event := r.Game.ReconnectPlayer(client.PlayerID, client.player.Pseudo)
	// 	r.gameLock.Unlock()

	// 	r.broadcastEvents(game.StateUpdate{
	// 		Events: game.Events{event},
	// 		State:  r.Game.GameState,
	// 	})
	// }
}

func (r *Room) onLeave(player Player) {
	r.lock.Lock()
	defer r.lock.Unlock()

	// if _, in := r.currentPlayers[client]; !in { // client who never joined may still end up here
	// 	continue
	// }
	// ProgressLogger.Printf("Game %s : removing player %d...", r.ID, r.currentPlayers[client])

	// playerID := r.currentPlayers[client]

	// // check if the player is not already removed
	// if r.Game.Players[playerID] == nil {
	// 	continue
	// }

	// r.gameLock.Lock()
	// update := r.Game.RemovePlayer(playerID)
	// r.gameLock.Unlock()

	// delete(r.currentPlayers, client)

	// if r.monitor != nil { // notify the monitor
	// 	r.monitor <- r.Summary()
	// }

	// r.broadcastEvents(update)
}

func (r *Room) onEvent(event game.ClientEvent) (isGameOver bool) {
	r.lock.Lock()
	defer r.lock.Unlock()

	ProgressLogger.Printf("Game %s : handleClientEvent...", r.ID)

	// events, isGameOver, err := r.Game.HandleClientEvent(message)
	// if err != nil { // malicious client: ignore the query
	// 	WarningLogger.Println(err)
	// 	continue
	// }

	// if events != nil {
	// 	r.broadcastEvents(*events)
	// }

	// if isGameOver {
	// 	ProgressLogger.Printf("Game %s is over: exitting game loop.", r.ID)
	// 	return r.review(), true
	// }
}

func (r *Room) onQuestionTimeout() {
	r.lock.Lock()
	defer r.lock.Unlock()
	// ProgressLogger.Printf("Game %s : questionTimeoutAction...", r.ID)

	// 		events := r.Game.QuestionTimeoutAction()
	// 		if events != nil {
	// 			r.broadcastEvents(*events)
	// }
}
