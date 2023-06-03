package trivial

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/benoitkugler/maths-online/server/src/pass"
	tcAPI "github.com/benoitkugler/maths-online/server/src/prof/teacher"
	"github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/sql/reviews"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	tc "github.com/benoitkugler/maths-online/server/src/sql/trivial"
	tv "github.com/benoitkugler/maths-online/server/src/trivial"
	"github.com/benoitkugler/maths-online/server/src/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	WarningLogger  = log.New(os.Stdout, "tv-store:ERROR:", 0)
	ProgressLogger = log.New(os.Stdout, "tv-store:INFO:", 0)
)

var sessionTimeout = 12 * time.Hour

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var errAccessForbidden = errors.New("trivial config access forbidden")

type uID = teacher.IdTeacher

// Controller is the top level (singleton) objects
// handling requests related to trivial pousuit setups
// It delegates to trivial-poursuit.GameController for the
// actual game logic handling.
type Controller struct {
	db         *sql.DB
	studentKey pass.Encrypter

	store gameStore

	admin teacher.Teacher
}

func NewController(db *sql.DB, key pass.Encrypter, demoPin string, admin teacher.Teacher) *Controller {
	out := &Controller{
		db:         db,
		studentKey: key,
		store:      newGameStore(demoPin),
		admin:      admin,
	}

	return out
}

func (ct *Controller) checkOwner(configID tc.IdTrivial, userID uID) error {
	in, err := tc.SelectTrivial(ct.db, configID)
	if err != nil {
		return utils.SQLError(err)
	}

	if in.IdTeacher != userID {
		return errAccessForbidden
	}

	return nil
}

type RunningSessionMetaOut struct {
	NbGames int
}

func (ct *Controller) GetTrivialRunningSessions(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	session := ct.store.getSessionID(userID)
	ct.store.lock.Lock()
	out := RunningSessionMetaOut{NbGames: len(ct.store.getSession(session))}
	ct.store.lock.Unlock()

	return c.JSON(200, out)
}

func (ct *Controller) getTrivialPoursuits(userID uID) ([]TrivialExt, error) {
	configs, err := tc.SelectAllTrivials(ct.db)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	revs, err := reviews.SelectAllReviewTrivials(ct.db)
	if err != nil {
		return nil, utils.SQLError(err)
	}
	revsMap := revs.ByIdTrivial()

	sel, err := newQuestionSelector(ct.db)
	if err != nil {
		return nil, err
	}

	var out []TrivialExt
	for _, config := range configs {
		var inReview tcAPI.OptionalIdReview
		link, isInReview := revsMap[config.Id]
		if isInReview {
			inReview = tcAPI.OptionalIdReview{InReview: true, Id: link.IdReview}
		}

		item, err := newTrivialExt(sel, config, inReview, userID, ct.admin.Id)
		if err != nil {
			return nil, err
		}
		if item.Origin.Visibility.Restricted() {
			continue // do not expose
		}
		out = append(out, item)
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Config.Id < out[j].Config.Id })

	return out, nil
}

func (ct *Controller) GetTrivialPoursuit(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)
	out, err := ct.getTrivialPoursuits(userID)
	if err != nil {
		return err
	}
	return c.JSON(200, out)
}

func (ct *Controller) CreateTrivialPoursuit(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	item, err := tc.Trivial{
		QuestionTimeout: 120,
		ShowDecrassage:  true,
		IdTeacher:       userID,
	}.Insert(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	sel, err := newQuestionSelector(ct.db)
	if err != nil {
		return err
	}
	out, err := newTrivialExt(sel, item, tcAPI.OptionalIdReview{}, userID, ct.admin.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) DeleteTrivialPoursuit(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	if err = ct.checkOwner(tc.IdTrivial(id), userID); err != nil {
		return err
	}

	_, err = tc.DeleteTrivialById(ct.db, tc.IdTrivial(id))
	if err != nil {
		return utils.SQLError(err)
	}

	return c.NoContent(200)
}

func (ct *Controller) DuplicateTrivialPoursuit(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	id, err := utils.QueryParamInt64(c, "id")
	if err != nil {
		return err
	}

	in, err := tc.SelectTrivial(ct.db, tc.IdTrivial(id))
	if err != nil {
		return utils.SQLError(err)
	}

	vis := tcAPI.NewVisibility(in.IdTeacher, userID, ct.admin.Id, in.Public)
	if vis.Restricted() {
		return errAccessForbidden
	}

	// attribute the new copy to the current owner, and make it private
	config := in
	config.IdTeacher = userID
	config.Public = false
	config, err = config.Insert(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	sel, err := newQuestionSelector(ct.db)
	if err != nil {
		return err
	}
	out, err := newTrivialExt(sel, config, tcAPI.OptionalIdReview{}, userID, ct.admin.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) UpdateTrivialPoursuit(c echo.Context) error {
	var params tc.Trivial
	if err := c.Bind(&params); err != nil {
		return err
	}

	userID := tcAPI.JWTTeacher(c)

	if err := ct.checkOwner(params.Id, userID); err != nil {
		return err
	}

	// ensure correct owner
	params.IdTeacher = userID

	config, err := params.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	sel, err := newQuestionSelector(ct.db)
	if err != nil {
		return err
	}

	var inReview tcAPI.OptionalIdReview
	item, isInReview, err := reviews.SelectReviewTrivialByIdTrivial(ct.db, config.Id)
	if err != nil {
		return utils.SQLError(err)
	}
	if isInReview {
		inReview = tcAPI.OptionalIdReview{InReview: true, Id: item.IdReview}
	}

	out, err := newTrivialExt(sel, config, inReview, userID, ct.admin.Id)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

type UpdateTrivialVisiblityIn struct {
	ConfigID tc.IdTrivial
	Public   bool
}

func (ct *Controller) UpdateTrivialVisiblity(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var args UpdateTrivialVisiblityIn
	if err := c.Bind(&args); err != nil {
		return fmt.Errorf("invalid parameters: %s", err)
	}

	qu, err := tc.SelectTrivial(ct.db, args.ConfigID)
	if err != nil {
		return utils.SQLError(err)
	}
	if qu.IdTeacher != userID {
		return errAccessForbidden
	}

	qu.Public = args.Public
	qu, err = qu.Update(ct.db)
	if err != nil {
		return utils.SQLError(err)
	}

	return c.NoContent(200)
}

type CheckMissingQuestionsOut struct {
	Pattern editor.Tags   // empty if no pattern is found
	Missing []editor.Tags // missing tags, may be empty if OK
}

// CheckMissingQuestions is an hint to avoid forgetting a tag
// when setting up the questions.
// It first finds the current common tags, and then checks that no
// questions sharing the same tags are left behind.
func (ct *Controller) CheckMissingQuestions(c echo.Context) error {
	var criteria tc.CategoriesQuestions
	if err := c.Bind(&criteria); err != nil {
		return err
	}

	userID := tcAPI.JWTTeacher(c)

	out, err := ct.checkMissingQuestions(criteria, userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) checkMissingQuestions(criteria tc.CategoriesQuestions, userID uID) (CheckMissingQuestionsOut, error) {
	criteria.Normalize()

	pattern := commonTags(criteria)
	if len(pattern) == 0 { // no pattern found, return early
		return CheckMissingQuestionsOut{}, nil
	}

	existingQuestiongroups, err := editor.SelectQuestiongroupByTags(ct.db, userID, pattern)
	if err != nil {
		return CheckMissingQuestionsOut{}, utils.SQLError(err)
	}
	var ids []editor.IdQuestiongroup
	for i := range existingQuestiongroups {
		ids = append(ids, i)
	}

	// restrict the search to the selected difficulties, if any,
	// to avoid false alerts
	questions, err := editor.SelectQuestionsByIdGroups(ct.db, ids...)
	if err != nil {
		return CheckMissingQuestionsOut{}, err
	}
	for id, question := range questions {
		if !criteria.Difficulties.Match(question.Difficulty) {
			delete(questions, id)
		}
	}

	pool, err := selectQuestions(ct.db, criteria, userID)
	if err != nil {
		return CheckMissingQuestionsOut{}, err
	}
	usedQuestions := allQuestions(pool)

	// check is existingQuestions is included in usedQuestions
	// if not, add the tags as hint
	hint := editor.NewTagListSet()
	for idExisting, question := range questions {
		if !usedQuestions.Has(idExisting) {
			tags := existingQuestiongroups[question.IdGroup.ID]
			hint.Add(tags.Tags())
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

	userID := tcAPI.JWTTeacher(c)

	out, err := ct.launchConfig(in, userID)
	if err != nil {
		return err
	}

	return c.JSON(200, out)
}

func (ct *Controller) launchConfig(params LaunchSessionIn, userID uID) (LaunchSessionOut, error) {
	config, err := tc.SelectTrivial(ct.db, params.IdConfig)
	if err != nil {
		return LaunchSessionOut{}, utils.SQLError(err)
	}

	// admin config may be launched, since it is a readonly operation
	if config.IdTeacher != userID && config.IdTeacher != ct.admin.Id {
		return LaunchSessionOut{}, errAccessForbidden
	}

	session := ct.store.getOrCreateSession(userID)

	// populate the session with the required games :

	// select the questions
	questionPool, err := selectQuestions(ct.db, config.Questions, userID)
	if err != nil {
		return LaunchSessionOut{}, err
	}

	groups, err := params.Groups.groups()
	if err != nil {
		return LaunchSessionOut{}, err
	}

	var out LaunchSessionOut
	for _, groupStrategy := range groups {
		options := tv.Options{
			Launch:          groupStrategy,
			QuestionTimeout: time.Second * time.Duration(config.QuestionTimeout),
			ShowDecrassage:  config.ShowDecrassage,
			Questions:       questionPool,
		}

		gameID := ct.store.newTeacherGameID(session)
		ct.store.createGame(createGame{
			ID:      gameID,
			Options: options,
		})
		out.GameIDs = append(out.GameIDs, gameID.roomID())
	}

	return out, nil
}

// check the code refers to a session actually owned by the teacher [userID]
func (ct *Controller) parseTeacherCode(gameID string, userID uID) (teacherCode, error) {
	parsed, err := ct.store.parseCode(gameID)
	if err != nil {
		return teacherCode{}, err
	}
	tc, ok := parsed.(teacherCode)
	if !ok {
		return teacherCode{}, errors.New("internal error: expected teacher code")
	}

	ct.store.lock.Lock()
	idTeacher := ct.store.teacherSessions[tc.sessionID]
	ct.store.lock.Unlock()
	if userID != idTeacher {
		return teacherCode{}, errAccessForbidden
	}

	return tc, nil
}

// StartTtrivialGame starts a game created in manual mode
func (ct *Controller) StartTrivialGame(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	gameID := c.QueryParam("game-id")

	parsed, err := ct.parseTeacherCode(gameID, userID)
	if err != nil {
		return err
	}
	err = ct.store.startGame(parsed)
	if err != nil {
		return err
	}

	return c.NoContent(200)
}

type stopGame struct {
	ID      string
	Restart bool // if false, definitively close the game
}

// StopTrivialGame stops a game, optionnaly restarting it,
// interrupting the connection with clients.
func (ct *Controller) StopTrivialGame(c echo.Context) error {
	userID := tcAPI.JWTTeacher(c)

	var in stopGame
	if err := c.Bind(&in); err != nil {
		return fmt.Errorf("invalid parameters format: %s", err)
	}

	parsed, err := ct.parseTeacherCode(in.ID, userID)
	if err != nil {
		return err
	}

	ct.store.stopGame(parsed, in.Restart)

	return c.NoContent(200)
}
