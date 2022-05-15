package trivialpoursuit

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/benoitkugler/maths-online/maths/exercice"
	"github.com/benoitkugler/maths-online/pass"
	"github.com/benoitkugler/maths-online/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	WarningLogger  = log.New(os.Stdout, "trivial-poursuit-session:ERROR:", log.LstdFlags)
	ProgressLogger = log.New(os.Stdout, "trivial-poursuit-session:INFO:", log.LstdFlags)
)

const (
	GameEndPoint = "/trivial/game/:session-id"
)

var sessionTimeout = 12 * time.Hour

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

// lookupSession locks and revese the sessionID -> session map
func (ct *Controller) sessionMap() map[int64]LaunchSessionOut {
	ct.lock.Lock()
	defer ct.lock.Unlock()

	out := make(map[int64]LaunchSessionOut, len(ct.sessions))

	for sessionID, session := range ct.sessions {
		out[session.config.Id] = LaunchSessionOut{
			SessionID:         sessionID,
			GroupStrategyKind: session.group.kind(),
			GroupsID:          session.groupIDs(),
		}
	}

	return out
}

func (ct *Controller) GetTrivialPoursuit(c echo.Context) error {
	configs, err := SelectAllTrivialConfigs(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	tags, err := exercice.SelectAllQuestionTags(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	dict := ct.sessionMap()

	tagsDict := tags.ByIdQuestion()
	var out []TrivialConfigExt
	for _, config := range configs {
		out = append(out, config.withDetails(tagsDict, dict))
	}

	return c.JSON(200, out)
}

func (ct *Controller) CreateTrivialPoursuit(c echo.Context) error {
	tc, err := TrivialConfig{
		QuestionTimeout: 120,
		ShowDecrassage:  true,
	}.Insert(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	out := TrivialConfigExt{Config: tc} // 0 questions by categories, not running

	return c.JSON(200, out)
}

func (ct *Controller) UpdateTrivialPoursuit(c echo.Context) error {
	var params TrivialConfig
	if err := c.Bind(&params); err != nil {
		return err
	}

	config, err := params.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	tags, err := exercice.SelectAllQuestionTags(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	out := config.withDetails(tags.ByIdQuestion(), ct.sessionMap())

	return c.JSON(200, out)
}

type CheckMissingQuestionsOut struct {
	Pattern []string   // empty if no pattern is found
	Missing [][]string // missing tags, may be empty if OK
}

// CheckMissingQuestions is an hint to avoid forgetting a tag
// when setting up the questions.
// It first finds the current common tags, and then checks that no
// questions sharing the same tags are left behind.
func (ct *Controller) CheckMissingQuestions(c echo.Context) error {
	var criteria CategoriesQuestions
	if err := c.Bind(&criteria); err != nil {
		return err
	}

	out, err := ct.checkMissingQuestions(criteria)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) checkMissingQuestions(criteria CategoriesQuestions) (CheckMissingQuestionsOut, error) {
	criteria.normalize()

	pattern := criteria.commonTags()
	if len(pattern) == 0 { // no pattern found, return early
		return CheckMissingQuestionsOut{}, nil
	}

	existingQuestions, err := exercice.SelectQuestionByTags(ct.db, pattern...)
	if err != nil {
		return CheckMissingQuestionsOut{}, utils.SQLError(err)
	}

	usedQuestions, err := criteria.selectQuestionIds(ct.db)
	if err != nil {
		return CheckMissingQuestionsOut{}, err
	}

	// check is existingQuestions is included in usedQuestions
	// if not, add the tags as hint
	hint := exercice.NewTagListSet()
	for idExisting, tags := range existingQuestions {
		if !usedQuestions.Has(idExisting) {
			hint.Add(tags.List())
		}
	}

	return CheckMissingQuestionsOut{
		Pattern: pattern,
		Missing: hint.List(),
	}, nil
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

// LaunchSessionTrivialPoursuit starts a new TrivialPoursuit session with
// the given config, and returns the updated version of the config.
func (ct *Controller) LaunchSessionTrivialPoursuit(c echo.Context) error {
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
	if dict := ct.sessionMap(); dict[params.IdConfig].SessionID != "" {
		return LaunchSessionOut{}, errors.New("session already running")
	}

	config, err := SelectTrivialConfig(ct.db, params.IdConfig)
	if err != nil {
		return LaunchSessionOut{}, utils.SQLError(err)
	}

	// select the questions
	questionPool, err := config.Questions.selectQuestions(ct.db)
	if err != nil {
		return LaunchSessionOut{}, err
	}

	var out LaunchSessionOut
	out.GroupStrategyKind = params.GroupStrategy.kind()

	ct.lock.Lock()
	out.SessionID = utils.RandomID(true, 4, func(s string) bool {
		_, taken := ct.sessions[s]
		return taken
	})
	ct.lock.Unlock()

	session := newGameSession(out.SessionID, ct.db, config, params.GroupStrategy, questionPool)

	// the rooms may be either created initially, or during client connection
	// (see connectStudent)
	session.group.initGames(session) // initial setup of rooms
	out.GroupsID = session.groupIDs()

	ct.lock.Lock()
	// register the controller...
	ct.sessions[out.SessionID] = session
	ct.lock.Unlock()

	ProgressLogger.Printf("Launching session %s", out.SessionID)
	// ...and start it
	go func() {
		ctx, cancelFunc := context.WithTimeout(context.Background(), sessionTimeout)
		session.startLoop(ctx)

		cancelFunc()

		// remove the game controller when the game is over
		ct.lock.Lock()
		defer ct.lock.Unlock()
		delete(ct.sessions, out.SessionID)

		ProgressLogger.Printf("Removed session %s", out.SessionID)
	}()

	return out, nil
}

// StopSessionTrivialPoursuit stops a running session,
// cleaning the controllers and interrupting the connection with clients.
func (ct *Controller) StopSessionTrivialPoursuit(c echo.Context) error {
	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	sessionID := ct.sessionMap()[id].SessionID

	ct.lock.Lock()
	session, ok := ct.sessions[sessionID]
	if !ok {
		// gracefully exit
	}
	ct.lock.Unlock()

	session.quit <- true

	return c.NoContent(200)
}

func (ct *Controller) ConnectTeacherMonitor(c echo.Context) error {
	sessionID := c.QueryParam("session-id")

	ct.lock.Lock()
	session, ok := ct.sessions[sessionID]
	if !ok {
		return fmt.Errorf("invalid session %s", sessionID)
	}
	ct.lock.Unlock()

	// upgrade this connection to a WebSocket connection
	ws, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	client := session.connectTeacher(ws)

	client.startLoop() // block

	session.lock.Lock() // remove the client
	delete(session.teacherClients, client)
	session.lock.Unlock()

	return nil
}

// ConnectStudentSession handles the connection of one student to the activity
func (ct *Controller) ConnectStudentSession(c echo.Context) error {
	fmt.Println("ConnectStudentSession")

	completeID := c.Param("session-id")

	if len(completeID) < 4 {
		return fmt.Errorf("invalid ID %s", completeID)
	}
	sessionID := completeID[:4]
	var student studentMeta
	student.gameID = completeID[4:]

	ct.lock.Lock()
	session, ok := ct.sessions[sessionID]
	if !ok {

		WarningLogger.Printf("invalid session ID %s", sessionID)

		return fmt.Errorf("L'activité n'existe pas ou est déjà terminée.")
	}
	ct.lock.Unlock()

	student.id = pass.EncryptedID(c.QueryParam("client-id"))
	student.pseudo = c.QueryParam("client-pseudo")

	ProgressLogger.Printf("Connecting student %v at %s", student, sessionID)

	err := session.connectStudent(c, student, ct.key)

	return err
}
