package trivial

import (
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
	tv "github.com/benoitkugler/maths-online/trivial-poursuit"
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
		demoPin:  demoPin,
		admin:    admin,
		sessions: make(map[string]*gameSession),
	}
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

type RunningSessionMetaOut struct {
	NbGames int
}

// getSession locks and may return nil if no games has started yet
func (ct *Controller) getSession(userID int64) *gameSession {
	ct.lock.Lock()
	defer ct.lock.Unlock()

	for _, session := range ct.sessions {
		if session.idTeacher == userID {
			return session
		}
	}

	return nil
}

func (ct *Controller) GetTrivialRunningSessions(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	var out RunningSessionMetaOut
	if session := ct.getSession(user.Id); session != nil {
		out = RunningSessionMetaOut{NbGames: len(session.games)}
	}

	return c.JSON(200, out)
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
		out = append(out, config.withDetails(tagsDict, origin))
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

	out := config.withDetails(tags.ByIdQuestion(), teacher.Origin{
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

	out := config.withDetails(tags.ByIdQuestion(), teacher.Origin{
		Visibility: teacher.Personnal,
		Owner:      user.Mail,
		IsPublic:   config.Public,
	})

	return c.JSON(200, out)
}

type UpdateTrivialVisiblityIn struct {
	ConfigID int64
	Public   bool
}

func (ct *Controller) UpdateTrivialVisiblity(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	var args UpdateTrivialVisiblityIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	qu, err := SelectTrivialConfig(ct.db, args.ConfigID)
	if err != nil {
		return utils.SQLError(err)
	}
	if qu.IdTeacher != user.Id {
		return accessForbidden
	}

	qu.Public = args.Public
	qu, err = qu.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	return c.NoContent(200)
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

// LaunchSessionTrivialPoursuit starts a new TrivialPoursuit session
// where the games are setup with the given config.
func (ct *Controller) LaunchSessionTrivialPoursuit(c echo.Context) error {
	var in LaunchSessionIn
	if err := c.Bind(&in); err != nil {
		return fmt.Errorf("invalid parameters format: %s", err)
	}

	user := teacher.JWTTeacher(c)

	out, err := ct.launchConfig(in, user.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) launchConfig(params LaunchSessionIn, userID int64) (LaunchSessionOut, error) {
	config, err := SelectTrivialConfig(ct.db, params.IdConfig)
	if err != nil {
		return LaunchSessionOut{}, utils.SQLError(err)
	}

	if config.IdTeacher != userID {
		return LaunchSessionOut{}, accessForbidden
	}

	session := ct.getOrCreateSession(userID)

	// populate the session with the required games :

	// select the questions
	questionPool, err := config.Questions.selectQuestions(ct.db, userID)
	if err != nil {
		return LaunchSessionOut{}, err
	}

	var out LaunchSessionOut
	for _, nbPlayers := range params.Groups {
		options := tv.GameOptions{
			PlayersNumber:   nbPlayers,
			QuestionTimeout: time.Second * time.Duration(config.QuestionTimeout),
			ShowDecrassage:  config.ShowDecrassage,
		}

		gameID := session.newGameID()
		session.createGame(createGame{
			ID:        gameID,
			Questions: questionPool,
			Options:   options,
		})
		out.GameIDs = append(out.GameIDs, string(gameID))
	}

	return out, nil
}

// StopTrivialGame stops a game, optionnaly restarting it,
// interrupting the connection with clients.
func (ct *Controller) StopTrivialGame(c echo.Context) error {
	user := teacher.JWTTeacher(c)

	var in stopGame
	if err := c.Bind(&in); err != nil {
		return fmt.Errorf("invalid parameters format: %s", err)
	}

	session := ct.getSession(user.Id)
	if session == nil {
		// gracefully exit
	} else {
		in.terminateChanel = true
		session.stopGameEvents <- in
	}

	return c.NoContent(200)
}
