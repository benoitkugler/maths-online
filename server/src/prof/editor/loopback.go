package editor

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/benoitkugler/maths-online/maths/exercice"
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
	incomingClient chan *previewClient
	broadcast      chan serverData

	clientLeft chan bool

	client *previewClient // initialy empty

	sessionID       string
	currentQuestion exercice.QuestionInstance
}

func newLoopbackController(sessionID string) *loopbackController {
	return &loopbackController{
		incomingClient: make(chan *previewClient),
		broadcast:      make(chan serverData),
		clientLeft:     make(chan bool),
		sessionID:      sessionID,
	}
}

func (ct *loopbackController) setQuestion(question exercice.QuestionInstance) {
	ct.currentQuestion = question
	ct.broadcast <- serverData{Kind: State, Data: LoopbackState{Question: question.ToClient()}}
}

func (ct *loopbackController) unsetQuestion() {
	ct.broadcast <- serverData{Kind: State, Data: LoopbackState{IsPaused: true}}
}

func (ct *loopbackController) startLoop(ctx context.Context) {
	for {
		select {
		case client := <-ct.incomingClient:
			ct.client = client
		case <-ct.clientLeft: // client is done
			log.Println("Client is done")
			return
		case <-ctx.Done(): // terminate the session on timeout
			log.Println("Session timed out")
			if ct.client != nil {
				utils.WebsocketError(ct.client.conn, errors.New("session timeout reached"))
			}
			return

		case data := <-ct.broadcast:
			if ct.client != nil {
				log.Println("Sending data...")
				ct.client.send(data)
			}
		}
	}
}

type serverData struct {
	Data interface{}
	Kind loopbackServerDataKind
}

type clientData struct {
	Data interface{}
	Kind loopbackClientDataKind
}

func (cld *clientData) UnmarshalJSON(data []byte) error {
	var wr struct {
		Kind loopbackClientDataKind
		Data json.RawMessage
	}
	if err := json.Unmarshal(data, &wr); err != nil {
		return err
	}
	cld.Kind = wr.Kind
	switch wr.Kind {
	case Ping: // nothing to do
	case CheckSyntaxIn:
		var content client.QuestionSyntaxCheckIn
		err := json.Unmarshal(wr.Data, &content)
		if err != nil {
			return err
		}
		cld.Data = content
	case ValidAnswerIn:
		var content client.QuestionAnswersIn
		err := json.Unmarshal(wr.Data, &content)
		if err != nil {
			return err
		}
		cld.Data = content
	}
	return nil
}

type previewClient struct {
	conn       *websocket.Conn
	controller *loopbackController
}

func (cl *previewClient) send(data serverData) {
	err := cl.conn.WriteJSON(data)
	if err != nil {
		log.Printf("Broadcasting to client (session %s) failed: %s", cl.controller.sessionID, err)
	}
}

// startLoop listens for new messages being sent to our WebSocket
// endpoint, only returning on error.
// the connection is not closed yet
func (cl *previewClient) startLoop() {
	for {
		// read in a message
		var data clientData
		err := cl.conn.ReadJSON(&data)
		if err != nil {
			log.Printf("invalid client messsage: %s", err)
			return
		}

		switch data.Kind {
		case CheckSyntaxIn:
			out := cl.controller.currentQuestion.CheckSyntaxe(data.Data.(client.QuestionSyntaxCheckIn))
			cl.controller.broadcast <- serverData{Kind: CheckSyntaxeOut, Data: out}
		case ValidAnswerIn:
			out := cl.controller.currentQuestion.EvaluateAnswer(data.Data.(client.QuestionAnswersIn))
			cl.controller.broadcast <- serverData{Kind: ValidAnswerOut, Data: out}
		}
	}

	// do not close the connection when the client is leaving
	// so that iframe reloads may use the same controller
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
	ct.incomingClient <- client

	client.startLoop()
}
