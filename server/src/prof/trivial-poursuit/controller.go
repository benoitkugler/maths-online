package trivialpoursuit

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	// demoPin is used to create testing games on the fly
	demoPin string
}

func NewController(db *sql.DB, key pass.Encrypter, demoPin string) *Controller {
	return &Controller{
		db:       db,
		key:      key,
		sessions: make(map[string]*gameSession),
		demoPin:  demoPin,
	}
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

func (ct *Controller) getTrivialPoursuits() ([]TrivialConfigExt, error) {
	configs, err := SelectAllTrivialConfigs(ct.db)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	tags, err := exercice.SelectAllQuestionTags(ct.db)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	dict := ct.sessionMap()

	tagsDict := tags.ByIdQuestion()
	var out []TrivialConfigExt
	for _, config := range configs {
		out = append(out, config.withDetails(tagsDict, dict))
	}
	return out, nil
}

func (ct *Controller) GetTrivialPoursuit(c echo.Context) error {
	out, err := ct.getTrivialPoursuits()
	if err != nil {
		return err
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

func (ct *Controller) createGameSession(sessionID string, config TrivialConfig, group GroupStrategy) (*gameSession, error) {
	// select the questions
	questionPool, err := config.Questions.selectQuestions(ct.db)
	if err != nil {
		return nil, err
	}

	session := newGameSession(sessionID, ct.db, config, group, questionPool)

	// the rooms may be either created initially, or during client connection
	// (see connectStudent)
	session.group.initGames(session) // initial setup of rooms

	ct.lock.Lock()
	// register the controller...
	ct.sessions[sessionID] = session
	ct.lock.Unlock()

	ProgressLogger.Printf("Launching session %s", sessionID)
	// ...and start it
	go func() {
		ctx, cancelFunc := context.WithTimeout(context.Background(), sessionTimeout)
		session.startLoop(ctx)

		cancelFunc()

		// remove the game controller when the game is over
		ct.lock.Lock()
		defer ct.lock.Unlock()
		delete(ct.sessions, sessionID)

		ProgressLogger.Printf("Removed session %s", sessionID)
	}()

	return session, nil
}

func (ct *Controller) launchSession(params LaunchSessionIn) (LaunchSessionOut, error) {
	if dict := ct.sessionMap(); dict[params.IdConfig].SessionID != "" {
		return LaunchSessionOut{}, errors.New("session already running")
	}

	config, err := SelectTrivialConfig(ct.db, params.IdConfig)
	if err != nil {
		return LaunchSessionOut{}, utils.SQLError(err)
	}

	ct.lock.Lock()
	sessionID := utils.RandomID(true, 4, func(s string) bool {
		_, taken := ct.sessions[s]
		return taken
	})
	ct.lock.Unlock()

	session, err := ct.createGameSession(sessionID, config, params.GroupStrategy)
	if err != nil {
		return LaunchSessionOut{}, err
	}

	out := LaunchSessionOut{
		GroupStrategyKind: params.GroupStrategy.kind(),
		SessionID:         sessionID,
		GroupsID:          session.groupIDs(),
	}
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

// expects <demoPin>.<number>
// or return 0
func (ct *Controller) isDemoSessionID(completeID string) (room string, nbPlayers int) {
	cuts := strings.Split(completeID, ".")
	if len(cuts) != 3 {
		return "", 0
	}
	if ct.demoPin != cuts[0] {
		return "", 0
	}
	room = cuts[1]
	if len(room) < 2 {
		return "", 0
	}
	nbPlayers, _ = strconv.Atoi(cuts[2])
	return room, nbPlayers
}

func (ct *Controller) connectDemo(c echo.Context, room string, nbPlayers int) error {
	sessionID := fmt.Sprintf("%s.%s.%d", ct.demoPin, room, nbPlayers)

	// check if the session is running and waiting for players
	ct.lock.Lock()
	session, ok := ct.sessions[sessionID]
	ct.lock.Unlock()

	if !ok {
		// create the session
		var err error
		session, err = ct.createGameSession(sessionID, TrivialConfig{
			Id:              -1,
			Questions:       demoQuestions,
			QuestionTimeout: 120,
			ShowDecrassage:  true,
		}, RandomGroupStrategy{
			MaxPlayersPerGroup: nbPlayers,
			TotalPlayersNumber: nbPlayers,
		})
		if err != nil {
			return err
		}
	}

	student := studentMeta{
		pseudo: c.QueryParam("client-pseudo"),
	}

	ProgressLogger.Printf("Connecting student %v at (demo) %s", student, sessionID)

	err := session.connectStudent(c, student, ct.key)

	return err
}

// ConnectStudentSession handles the connection of one student to the activity
func (ct *Controller) ConnectStudentSession(c echo.Context) error {
	fmt.Println("ConnectStudentSession")

	completeID := c.Param("session-id")

	if room, nbPlayers := ct.isDemoSessionID(completeID); nbPlayers != 0 {
		return ct.connectDemo(c, room, nbPlayers)
	}

	if len(completeID) < 4 {
		return fmt.Errorf("invalid ID %s", completeID)
	}
	sessionID := completeID[:4]
	var student studentMeta
	student.gameID = completeID[4:]

	ct.lock.Lock()
	session, ok := ct.sessions[sessionID]
	ct.lock.Unlock()
	if !ok {
		WarningLogger.Printf("invalid session ID %s", sessionID)
		return fmt.Errorf("L'activité n'existe pas ou est déjà terminée.")
	}

	student.id = pass.EncryptedID(c.QueryParam("client-id"))
	student.pseudo = c.QueryParam("client-pseudo")

	ProgressLogger.Printf("Connecting student %v at %s", student, sessionID)

	err := session.connectStudent(c, student, ct.key)

	return err
}
