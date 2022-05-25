package trivialpoursuit

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/benoitkugler/maths-online/pass"
	"github.com/benoitkugler/maths-online/prof/editor"
	"github.com/benoitkugler/maths-online/prof/teacher"
	"github.com/benoitkugler/maths-online/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	WarningLogger  = log.New(os.Stdout, "trivial-poursuit-session:ERROR:", log.LstdFlags)
	ProgressLogger = log.New(os.Stdout, "trivial-poursuit-session:INFO:", log.LstdFlags)
)

var sessionTimeout = 12 * time.Hour

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var accessForbidden = errors.New("trivial config access forbidden")

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

	admin teacher.Teacher
}

func NewController(db *sql.DB, key pass.Encrypter, demoPin string, admin teacher.Teacher) *Controller {
	return &Controller{
		db:       db,
		key:      key,
		sessions: make(map[string]*gameSession),
		demoPin:  demoPin,
		admin:    admin,
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

func (ct *Controller) checkOwner(configID, userID int64) error {
	in, err := SelectTrivialConfig(ct.db, configID)
	if err != nil {
		return utils.SQLError(err)
	}

	if in.IdTeacher != userID {
		return accessForbidden
	}

	return nil
}

func (ct *Controller) getTrivialPoursuits(userID int64) ([]TrivialConfigExt, error) {
	configs, err := SelectAllTrivialConfigs(ct.db)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	teachers, err := teacher.SelectAllTeachers(ct.db)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	tags, err := editor.SelectAllQuestionTags(ct.db)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	dict := ct.sessionMap()

	tagsDict := tags.ByIdQuestion()
	var out []TrivialConfigExt
	for _, config := range configs {
		vis, hasAcces := teacher.NewVisibility(config.IdTeacher, userID, ct.admin.Id, config.Public)
		if !hasAcces {
			continue // do not expose
		}
		origin := teacher.Origin{
			Visibility: vis,
			Owner:      teachers[config.IdTeacher].Mail,
			IsPublic:   config.Public,
		}
		out = append(out, config.withDetails(tagsDict, dict, origin))
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Config.Id < out[j].Config.Id })

	return out, nil
}

func (ct *Controller) GetTrivialPoursuit(c echo.Context) error {
	teacher := teacher.JWTTeacher(c)
	out, err := ct.getTrivialPoursuits(teacher.Id)
	if err != nil {
		return err
	}
	return c.JSON(200, out)
}

func (ct *Controller) CreateTrivialPoursuit(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	tc, err := TrivialConfig{
		QuestionTimeout: 120,
		ShowDecrassage:  true,
		IdTeacher:       user.Id,
	}.Insert(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	out := TrivialConfigExt{Config: tc} // 0 questions by categories, not running

	return c.JSON(200, out)
}

func (ct *Controller) DeleteTrivialPoursuit(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	if err = ct.checkOwner(id, user.Id); err != nil {
		return err
	}

	_, err = DeleteTrivialConfigById(ct.db, id)
	if err != nil {
		return utils.SQLError(err)
	}

	return c.NoContent(200)
}

func (ct *Controller) DuplicateTrivialPoursuit(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	in, err := SelectTrivialConfig(ct.db, id)
	if err != nil {
		return utils.SQLError(err)
	}

	_, hasAcces := teacher.NewVisibility(in.IdTeacher, user.Id, ct.admin.Id, in.Public)
	if !hasAcces {
		return accessForbidden
	}

	// attribute the new copy to the current owner, and make it private
	config := in
	config.IdTeacher = user.Id
	config.Public = false
	config, err = config.Insert(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	tags, err := editor.SelectAllQuestionTags(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	out := config.withDetails(tags.ByIdQuestion(), ct.sessionMap(), teacher.Origin{
		Visibility: teacher.Personnal,
		Owner:      user.Mail,
		IsPublic:   config.Public,
	},
	)

	return c.JSON(200, out)
}

func (ct *Controller) UpdateTrivialPoursuit(c echo.Context) error {
	var params TrivialConfig
	if err := c.Bind(&params); err != nil {
		return err
	}

	user := teacher.JWTTeacher(c)

	if err := ct.checkOwner(params.Id, user.Id); err != nil {
		return err
	}

	// ensure correct owner
	params.IdTeacher = user.Id

	config, err := params.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	tags, err := editor.SelectAllQuestionTags(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	out := config.withDetails(tags.ByIdQuestion(), ct.sessionMap(), teacher.Origin{
		Visibility: teacher.Personnal,
		Owner:      user.Mail,
		IsPublic:   config.Public,
	})

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

	user := teacher.JWTTeacher(c)

	out, err := ct.checkMissingQuestions(criteria, user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) checkMissingQuestions(criteria CategoriesQuestions, userID int64) (CheckMissingQuestionsOut, error) {
	criteria.normalize()

	pattern := criteria.commonTags()
	if len(pattern) == 0 { // no pattern found, return early
		return CheckMissingQuestionsOut{}, nil
	}

	existingQuestions, err := editor.SelectQuestionByTags(ct.db, userID, pattern...)
	if err != nil {
		return CheckMissingQuestionsOut{}, utils.SQLError(err)
	}

	usedQuestions, err := criteria.selectQuestionIds(ct.db, userID)
	if err != nil {
		return CheckMissingQuestionsOut{}, err
	}

	// check is existingQuestions is included in usedQuestions
	// if not, add the tags as hint
	hint := editor.NewTagListSet()
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

// LaunchSessionTrivialPoursuit starts a new TrivialPoursuit session with
// the given config, and returns the updated version of the config.
func (ct *Controller) LaunchSessionTrivialPoursuit(c echo.Context) error {
	var in LaunchSessionIn
	if err := c.Bind(&in); err != nil {
		return fmt.Errorf("invalid parameters format: %s", err)
	}

	user := teacher.JWTTeacher(c)

	out, err := ct.launchSession(in, user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) createGameSession(sessionID string, config TrivialConfig, group GroupStrategy, userID int64) (*gameSession, error) {
	// select the questions
	questionPool, err := config.Questions.selectQuestions(ct.db, userID)
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

func (ct *Controller) launchSession(params LaunchSessionIn, userID int64) (LaunchSessionOut, error) {
	if dict := ct.sessionMap(); dict[params.IdConfig].SessionID != "" {
		return LaunchSessionOut{}, errors.New("session already running")
	}

	config, err := SelectTrivialConfig(ct.db, params.IdConfig)
	if err != nil {
		return LaunchSessionOut{}, utils.SQLError(err)
	}

	if config.IdTeacher != userID {
		return LaunchSessionOut{}, accessForbidden
	}

	ct.lock.Lock()
	sessionID := utils.RandomID(true, 4, func(s string) bool {
		_, taken := ct.sessions[s]
		return taken
	})
	ct.lock.Unlock()

	session, err := ct.createGameSession(sessionID, config, params.GroupStrategy, userID)
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

	user := teacher.JWTTeacher(c)

	if err := ct.checkOwner(id, user.Id); err != nil {
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

	ProgressLogger.Println("Connecting teacher on session", sessionID)

	client := session.connectTeacher(ws)

	client.startLoop() // block

	session.lock.Lock() // remove the client
	delete(session.teacherClients, client)
	session.lock.Unlock()

	return nil
}
