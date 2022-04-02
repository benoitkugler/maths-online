package editor

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/benoitkugler/maths-online/maths/exercice/client"
	"github.com/benoitkugler/maths-online/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type loopbackController struct {
	broadcast  chan client.Question // the question to display in the preview
	clientLeft chan bool
	client     *previewClient // initialy empty
	sessionID  string
}

func newLoopbackController(sessionID string) *loopbackController {
	return &loopbackController{
		broadcast:  make(chan client.Question),
		clientLeft: make(chan bool),
		sessionID:  sessionID,
	}
}

func (ct *loopbackController) startLoop(ctx context.Context) {
	for {
		select {
		case <-ct.clientLeft: // client is done
			log.Println("Client is done")
			return

		case <-ctx.Done(): // terminate the session on timeout
			log.Println("Session timed out")
			if ct.client != nil {
				utils.WebsocketError(ct.client.conn, errors.New("session timeout reached"))
			}
			return

		case question := <-ct.broadcast:
			log.Println("Sending question...")
			if ct.client != nil {
				err := ct.client.sendQuestion(question)
				if err != nil {
					log.Printf("Broadcasting to client (session %s) failed: %s", ct.sessionID, err)
				}
			}
		}
	}
}

type previewClient struct {
	conn       *websocket.Conn
	controller *loopbackController
}

func (cl *previewClient) sendQuestion(question client.Question) error {
	return cl.conn.WriteJSON(question)
}

// startLoop listens for new messages being sent to our WebSocket
// endpoint, only returning on error.
// the connection is not closed yet
func (cl *previewClient) startLoop() {
	defer func() {
		cl.controller.clientLeft <- true
	}()

	for {
		// read in a message, expecting a stay alive ping
		_, _, err := cl.conn.ReadMessage()
		if err != nil {
			log.Printf("client connection: %s", err)
			return
		}
	}
}

func (ct *loopbackController) setupWebSocket(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to init websocket: ", err)
		return
	}
	defer ws.Close()

	client := &previewClient{conn: ws, controller: ct}
	ct.client = client

	client.startLoop()
}
