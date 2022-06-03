package trivial

import (
	"fmt"
	"sort"

	"github.com/benoitkugler/maths-online/prof/teacher"
	tv "github.com/benoitkugler/maths-online/trivial-poursuit"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// teacherClient represents the teacher browser
type teacherClient struct {
	conn             *websocket.Conn
	currentSummaries map[tv.GameID]tv.GameSummary
}

func newGameSummary(s tv.GameSummary) (out gameSummary) {
	out.GameID = s.ID
	out.RoomSize = s.RoomSize
	for p, su := range s.Successes {
		out.Players = append(out.Players, gamePlayers{
			Player:    p.Pseudo,
			Successes: su,
		})
	}

	sort.Slice(out.Players, func(i, j int) bool { return out.Players[i].Player < out.Players[j].Player })

	if s.PlayerTurn != nil {
		out.CurrentPlayer = s.PlayerTurn.Pseudo
	}
	return out
}

func (tc *teacherClient) socketData() (out teacherSocketData) {
	for _, su := range tc.currentSummaries {
		out.Games = append(out.Games, newGameSummary(su))
	}

	sort.Slice(out.Games, func(i, j int) bool {
		return out.Games[i].GameID < out.Games[j].GameID
	})

	return out
}

func (tc *teacherClient) removeGame(id tv.GameID) {
	delete(tc.currentSummaries, id)
	tc.conn.WriteJSON(tc.socketData())
}

func (tc *teacherClient) sendSummary(gs tv.GameSummary) {
	// update the list
	tc.currentSummaries[gs.ID] = gs

	tc.conn.WriteJSON(tc.socketData())
}

// start loop listenning for ping messages
func (tc *teacherClient) startLoop() {
	for {
		// read in a message
		_, _, err := tc.conn.ReadMessage()
		if err != nil {
			WarningLogger.Printf("teacher connection: %s", err)
			return
		}
	}
}

func (gs *gameSession) connectTeacher(ws *websocket.Conn) *teacherClient {
	client := &teacherClient{conn: ws, currentSummaries: make(map[tv.GameID]tv.GameSummary)}

	// start with the current summary for all running sessions
	gs.lock.Lock()
	for _, ga := range gs.games {
		su := ga.Summary()
		client.currentSummaries[su.ID] = su
	}

	gs.teacherClients[client] = true // register the client
	gs.lock.Unlock()

	client.conn.WriteJSON(client.socketData())

	return client
}

func (ct *Controller) ConnectTeacherMonitor(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	session := ct.getSession(user.Id)
	if session == nil {
		return fmt.Errorf("no running session for %d", user.Id)
	}

	// upgrade this connection to a WebSocket connection
	ws, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	ProgressLogger.Printf("Connecting teacher %d", user.Id)

	client := session.connectTeacher(ws)

	client.startLoop() // block

	session.lock.Lock() // remove the client
	defer session.lock.Unlock()
	delete(session.teacherClients, client)

	return nil
}
