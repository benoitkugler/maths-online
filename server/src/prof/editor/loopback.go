package editor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/benoitkugler/maths-online/maths/questions"
	"github.com/benoitkugler/maths-online/utils"
	"github.com/gorilla/websocket"
)

var LoopackLogger = log.New(os.Stdout, "editor-loopback:", log.LstdFlags)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type loopbackController struct {
	incomingClient chan *previewClient
	broadcast      chan LoopbackServerEvent

	clientLeft chan bool

	client *previewClient // initialy empty

	sessionID       string
	currentQuestion questions.QuestionInstance

	currentExercice      InstantiatedExercice
	currentQuestionIndex int // in the current exercice
}

func newLoopbackController(sessionID string) *loopbackController {
	return &loopbackController{
		incomingClient: make(chan *previewClient),
		broadcast:      make(chan LoopbackServerEvent),
		clientLeft:     make(chan bool),
		sessionID:      sessionID,
	}
}

func (ct *loopbackController) setQuestion(question questions.QuestionInstance) {
	ct.currentQuestion = question
	ct.broadcast <- loopbackQuestion{Question: question.ToClient()}
}

func (ct *loopbackController) setExercice(exercice InstantiatedExercice) {
	ct.currentExercice = exercice
	ct.broadcast <- loopbackShowExercice{Exercice: exercice, Progression: ProgressionExt{
		Progression:  Progression{}, // ignored by the client
		NextQuestion: ct.currentQuestionIndex,
		Questions:    make([]QuestionHistory, len(exercice.Questions)),
	}}
}

func (ct *loopbackController) pause() {
	ct.currentQuestion = questions.QuestionInstance{}
	ct.currentExercice = InstantiatedExercice{}
	ct.currentQuestionIndex = 0

	ct.broadcast <- loopbackPaused{}
}

func (ct *loopbackController) startLoop(ctx context.Context) {
	for {
		select {
		case client := <-ct.incomingClient:
			ct.client = client
		case <-ct.clientLeft: // client is done
			LoopackLogger.Println("Client is done, terminating session.")
			return
		case <-ctx.Done(): // terminate the session on timeout
			LoopackLogger.Println("Session timed out")
			if ct.client != nil {
				utils.WebsocketError(ct.client.conn, errors.New("session timeout reached"))
			}
			return
		case data := <-ct.broadcast:
			if ct.client != nil {
				LoopackLogger.Println("Broadcasting...")
				ct.client.send(data)
			}
		}
	}
}

type previewClient struct {
	conn       *websocket.Conn
	controller *loopbackController
}

func (cl *previewClient) send(data LoopbackServerEvent) {
	err := cl.conn.WriteJSON(LoopbackServerEventWrapper{Data: data})
	if err != nil {
		LoopackLogger.Printf("Broadcasting to client (session %s) failed: %s", cl.controller.sessionID, err)
	}
}

// startLoop listens for new messages being sent to our WebSocket
// endpoint, only returning on error.
// the connection is not closed yet
func (cl *previewClient) startLoop() {
	for {
		// read in a message
		var data LoopbackClientEventWrapper
		err := cl.conn.ReadJSON(&data)
		if err != nil {
			LoopackLogger.Printf("invalid client messsage: %s", err)
			return
		}

		switch data := data.Data.(type) {
		case loopbackPing:
			LoopackLogger.Println("Ping (ignoring)")
		case loopbackQuestionValidIn:
			out := cl.controller.currentQuestion.EvaluateAnswer(data.Answers)
			cl.controller.broadcast <- loopbackQuestionValidOut{out}
		case loopbackQuestionCorrectAnswersIn:
			out := cl.controller.currentQuestion.CorrectAnswer()
			cl.controller.broadcast <- loopbackQuestionCorrectAnswersOut{out}
		case loopbackExerciceValidIn:
			// TODO:
			fmt.Println(data)
		}
	}

	// do not close the connection when the client is leaving
	// so that iframe reloads may use the same controller
}

func (ct *loopbackController) setupWebSocket(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		LoopackLogger.Println("Failed to init websocket: ", err)
		return
	}
	defer ws.Close()

	client := &previewClient{conn: ws, controller: ct}
	ct.incomingClient <- client

	client.startLoop()
}
