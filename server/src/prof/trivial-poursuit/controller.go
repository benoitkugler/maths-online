package trivialpoursuit

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/benoitkugler/maths-online/pass"
	"github.com/benoitkugler/maths-online/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

const (
	GameEndPoint = "/trivial/game/:session-id"
)

var (
	sessionTimeout = 12 * time.Hour
	gameTimeout    = 6 * time.Hour
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// SessionID is a 4 digit identifier used
// by students to access one activity
type SessionID = string

// Controller is the top level (singleton) objects
// handling requests related to trivial pousuit setups
// It delegates to trivial-poursuit.GameController for the
// actual game logic handling.
type Controller struct {
	db       *sql.DB
	key      pass.Encrypter
	lock     sync.Mutex
	sessions map[SessionID]*gameSession
}

func NewController(db *sql.DB, key pass.Encrypter) *Controller {
	return &Controller{db: db, key: key, sessions: make(map[string]*gameSession)}
}

func (ct *Controller) GetTrivialPoursuit(c echo.Context) error {
	out, err := SelectAllTrivialConfigs(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	return c.JSON(200, out)
}

func (ct *Controller) CreateTrivialPoursuit(c echo.Context) error {
	out, err := TrivialConfig{QuestionTimeout: 60}.Insert(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	return c.JSON(200, out)
}

func (ct *Controller) UpdateTrivialPoursuit(c echo.Context) error {
	var params TrivialConfig
	if err := c.Bind(&params); err != nil {
		return err
	}

	out, err := params.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	return c.JSON(200, out)
}

func (ct *Controller) DeleteTrivialPoursuit(c echo.Context) error {
	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}
	_, err = DeleteTrivialConfigById(ct.db, id)
	if err != nil {
		return utils.SQLError(err)
	}

	return c.NoContent(200)
}

type LaunchSessionOut struct {
	SessionID SessionID
}

// LaunchSession starts a new TrivialPoursuit session with
// the given config.
func (ct *Controller) LaunchSession(c echo.Context) error {
	var in LaunchSessionIn
	if err := c.Bind(&in); err != nil {
		return fmt.Errorf("invalid parameters format: %s", err)
	}

	out, err := ct.launchSession(in)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) launchSession(params LaunchSessionIn) (LaunchSessionOut, error) {
	config, err := SelectTrivialConfig(ct.db, params.IdConfig)
	if err != nil {
		return LaunchSessionOut{}, utils.SQLError(err)
	}

	if config.IsLaunched {
		return LaunchSessionOut{}, errors.New("session already running")
	}

	ct.lock.Lock()
	newID := utils.RandomID(true, 4, func(s string) bool {
		_, taken := ct.sessions[s]
		return taken
	})
	ct.lock.Unlock()

	session := newGameSession(ct.db, config, params.GroupStrategy)

	session.group.initGames(session) // initial setup of rooms

	// register the controller...
	ct.sessions[newID] = session
	// ...and start it
	go func() {
		ctx, cancelFunc := context.WithTimeout(context.Background(), sessionTimeout)
		session.startLoop(ctx)

		cancelFunc()

		// remove the game controller when the game is over
		ct.lock.Lock()
		defer ct.lock.Unlock()
		delete(ct.sessions, newID)
	}()

	// mark the session as started
	config.IsLaunched = true
	_, err = config.Update(ct.db)
	if err != nil {
		return LaunchSessionOut{}, utils.SQLError(err)
	}

	return LaunchSessionOut{SessionID: newID}, nil
}

func (ct *Controller) ConnectTeacherMonitor(c echo.Context) error {
	sessionID := c.QueryParam("session-id")
	ct.lock.Lock()
	defer ct.lock.Unlock()

	session, ok := ct.sessions[sessionID]
	if !ok {
		return fmt.Errorf("invalid session %s", sessionID)
	}

	// upgrade this connection to a WebSocket connection
	ws, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	client := &teacherClient{conn: ws}

	session.lock.Lock()
	session.teacherClient = client
	session.lock.Unlock()

	client.startLoop() // block
	return nil
}

// ConnectStudentSession handles the connection of one student to the activity
func (ct *Controller) ConnectStudentSession(c echo.Context) error {
	ct.lock.Lock()
	defer ct.lock.Unlock()

	sessionID := c.Param("session-id")

	session, ok := ct.sessions[sessionID]
	if !ok {
		log.Printf("invalid session ID %s", sessionID)
		return fmt.Errorf("L'activité n'existe pas ou est déjà terminée.")
	}

	clientID := pass.EncryptedID(c.QueryParam("client-id"))

	err := session.connectStudent(c, clientID, ct.key)

	return err
}