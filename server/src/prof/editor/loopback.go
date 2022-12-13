package editor

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/benoitkugler/maths-online/server/src/maths/questions"
	"github.com/benoitkugler/maths-online/server/src/sql/tasks"
	taAPI "github.com/benoitkugler/maths-online/server/src/tasks"
	"github.com/benoitkugler/maths-online/server/src/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
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

	currentExercice taAPI.InstantiatedWork
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
	ct.broadcast <- LoopbackShowQuestion{Question: question.ToClient()}
}

func (ct *loopbackController) setExercice(exercice taAPI.InstantiatedWork) {
	ct.currentExercice = exercice
	ct.broadcast <- LoopbackShowExercice{Exercice: exercice, Progression: taAPI.ProgressionExt{
		NextQuestion: 0,
		Questions:    make([]tasks.QuestionHistory, len(exercice.Questions)),
	}}
}

func (ct *loopbackController) pause() {
	ct.currentQuestion = questions.QuestionInstance{}
	ct.currentExercice = taAPI.InstantiatedWork{}

	ct.broadcast <- LoopbackPaused{}
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
	// for {
	// 	// read in a message
	// 	var data LoopbackClientEventWrapper
	// 	err := cl.conn.ReadJSON(&data)
	// 	if err != nil {
	// 		LoopackLogger.Printf("invalid client messsage: %s", err)
	// 		return
	// 	}

	// 	switch data := data.Data.(type) {
	// 	case loopbackPing:
	// 		LoopackLogger.Println("Ping (ignoring)")
	// 	case loopbackQuestionValidIn:
	// 		out := cl.controller.currentQuestion.EvaluateAnswer(data.Answers)
	// 		cl.controller.broadcast <- loopbackQuestionValidOut{out}
	// 	case loopbackQuestionCorrectAnswersIn:
	// 		out := cl.controller.currentQuestion.CorrectAnswer()
	// 		cl.controller.broadcast <- loopbackQuestionCorrectAnswersOut{out}
	// 	case loopbackExerciceCorrectAnswsersIn:
	// 		questions := cl.controller.currentExercice.Questions
	// 		if data.QuestionIndex < 0 || data.QuestionIndex >= len(questions) {
	// 			LoopackLogger.Printf("invalid question index: %d", data.QuestionIndex)
	// 			continue
	// 		}
	// 		answers := questions[data.QuestionIndex].Instance().CorrectAnswer()
	// 		cl.controller.broadcast <- loopbackExerciceCorrectAnswersOut{
	// 			QuestionIndex: data.QuestionIndex,
	// 			Answers:       answers,
	// 		}
	// 	}
	// }

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

// LoopackEvaluateQuestion expects a question definition, a set of
// random variables, and an answer, and performs the evaluation.
func (ct *Controller) LoopackEvaluateQuestion(c echo.Context) error {
	var args LoopackEvaluateQuestionIn

	if err := c.Bind(&args); err != nil {
		return err
	}

	ans, err := taAPI.EvaluateQuestion(args.Question, args.Answer)
	if err != nil {
		return err
	}

	out := LoopbackEvaluateQuestionOut{ans}

	return c.JSON(200, out)
}

// LoopbackShowQuestionAnswer expects a question, random parameters,
// and returns the correct answer for these parameters
func (ct *Controller) LoopbackShowQuestionAnswer(c echo.Context) error {
	var args LoopbackShowQuestionAnswerIn

	if err := c.Bind(&args); err != nil {
		return err
	}

	p, err := args.Params.ToMap()
	if err != nil {
		return err
	}
	instance, err := args.Question.InstantiateWith(p)
	if err != nil {
		return err
	}
	ans := instance.CorrectAnswer()

	out := LoopbackShowQuestionAnswerOut{ans}
	return c.JSON(200, out)
}
