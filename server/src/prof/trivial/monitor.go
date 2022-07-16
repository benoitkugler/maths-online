package trivial

import (
	"fmt"
	"sort"

	"github.com/benoitkugler/maths-online/prof/teacher"
	tv "github.com/benoitkugler/maths-online/trivial"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// teacherClient represents the teacher browser
type teacherClient struct {
	conn *websocket.Conn
}

func newGameSummary(s tv.Summary) (out gameSummary) {
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

func socketData(summaries map[tv.RoomID]tv.Summary) (out teacherSocketData) {
	for _, su := range summaries {
		out.Games = append(out.Games, newGameSummary(su))
	}

	sort.Slice(out.Games, func(i, j int) bool {
		return out.Games[i].GameID < out.Games[j].GameID
	})

	return out
}

func (tc *teacherClient) sendSummary(summaries map[tv.RoomID]tv.Summary) {
	tc.conn.WriteJSON(socketData(summaries))
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

// lock and fetch summaries
func (gs *gameSession) collectSummaries() map[tv.RoomID]tv.Summary {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	out := make(map[tv.RoomID]tv.Summary)
	for _, ga := range gs.games {
		su := ga.Summary()
		out[su.ID] = su
	}

	return out
}

func (gs *gameSession) connectTeacher(ws *websocket.Conn) *teacherClient {
	client := &teacherClient{conn: ws}

	gs.lock.Lock()
	gs.teacherClients[client] = true // register the client
	gs.lock.Unlock()

	// start with the current summary for all running sessions
	client.sendSummary(gs.collectSummaries())

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
