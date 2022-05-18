package trivialpoursuit

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/benoitkugler/maths-online/pass"
	"github.com/benoitkugler/maths-online/prof/students"
	tv "github.com/benoitkugler/maths-online/trivial-poursuit"
	ga "github.com/benoitkugler/maths-online/trivial-poursuit/game"
	"github.com/benoitkugler/maths-online/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// GameID is an in-memory identifier for a game room,
// with length 2 (meaning only 100 games per session are allowed)
// It is meant to be associated with a session ID.
type GameID = string

// gameSession monitor the games of one session (think one classroom)
// and broadcast the main advances from all the games to the teacher client
type gameSession struct {
	id                   SessionID
	quit                 chan bool
	notifyMonitorEndGame chan string // partial game ID

	lock sync.Mutex

	db *sql.DB

	teacherClients map[*teacherClient]bool
	monitor        chan tv.GameSummary
	questions      ga.QuestionPool
	games          map[GameID]*tv.GameController // active games, either in lobby or playing

	config TrivialConfig // database entry, cached for simplicity
	group  GroupStrategy // specified when starting the session
}

func newGameSession(id SessionID, db *sql.DB, config TrivialConfig, group GroupStrategy, questions ga.QuestionPool) *gameSession {
	return &gameSession{
		id:                   id,
		db:                   db,
		config:               config,
		group:                group,
		quit:                 make(chan bool),
		notifyMonitorEndGame: make(chan string),
		monitor:              make(chan tv.GameSummary),
		games:                make(map[string]*tv.GameController),
		teacherClients:       make(map[*teacherClient]bool),
		questions:            questions,
	}
}

// make sure game id are properly sorted when
// created from groups
func gameIDFromSerial(serial int) string {
	return fmt.Sprintf("%02d", serial+1)
}

// createGame locks and starts a new game
func (gs *gameSession) createGame(nbPlayers int) GameID {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	gameID := gameIDFromSerial(len(gs.games))

	game := tv.NewGameController(gameID,
		gs.questions,
		tv.GameOptions{
			PlayersNumber:   nbPlayers,
			QuestionTimeout: time.Second * time.Duration(gs.config.QuestionTimeout),
			ShowDecrassage:  gs.config.ShowDecrassage,
		},
		gs.monitor)
	// register the controller...
	gs.games[gameID] = game
	// ...and start it
	go func() {
		review, ok := game.StartLoop()
		if ok { // exploit the review
			gs.exploitReview(review)
		}

		gs.lock.Lock()
		defer gs.lock.Unlock()
		// remove the game controller when the game is over
		delete(gs.games, gameID)

		gs.notifyMonitorEndGame <- gameID
	}()

	return gameID
}

// TODO:
func (gs *gameSession) exploitReview(review tv.Review) {
	ProgressLogger.Printf("GAME REVIEW: %v", review)
}

type studentMeta struct {
	id     pass.EncryptedID
	gameID string // as requested by the client
	pseudo string // used for annonymous connection
}

func (gs *gameSession) connectStudent(c echo.Context, student studentMeta, key pass.Encrypter) error {
	player := tv.Player{ID: student.id}
	var studentID int64 = -1
	if student.id == "" { // anonymous connection
		player.Name = student.pseudo
		if player.Name == "" {
			player.Name = gs.generateName() // finally generate a random pseudo
		}
	} else { // fetch name from DB
		var err error
		studentID, err = student.id.Decrypt(key)
		if err != nil {
			return fmt.Errorf("invalid student ID: %s", err)
		}

		student, err := students.SelectStudent(gs.db, studentID)
		if err != nil {
			return utils.SQLError(err)
		}

		player.Name = student.Surname
	}

	// select or create the game
	gameID, err := gs.group.selectOrCreateGame(student.gameID, studentID, gs)
	if err != nil {
		return err
	}

	// then add the player
	gs.lock.Lock()
	game := gs.games[gameID]
	gs.lock.Unlock()

	// connect to the websocket handler, which handle errors
	game.AddClient(c.Response().Writer, c.Request(), player) // block on the client WS

	return nil
}

func (gs *gameSession) connectTeacher(ws *websocket.Conn) *teacherClient {
	client := &teacherClient{conn: ws, currentSummaries: make(map[string]tv.GameSummary)}

	// start with the current summary for all running sessions
	gs.lock.Lock()
	defer gs.lock.Unlock()

	for _, ga := range gs.games {
		su := ga.Summary()
		su.ID = gs.id + su.ID
		client.currentSummaries[su.ID] = su
	}
	gs.teacherClients[client] = true

	client.conn.WriteJSON(client.socketData())

	return client
}

func (gs *gameSession) isDemo() bool {
	return gs.config.Id == -1
}

func (gs *gameSession) startLoop(ctx context.Context) {
	for {
		select {
		case <-gs.quit:
			for _, game := range gs.games {
				game.Terminate <- true
			}
			return
		case <-ctx.Done():
			for _, game := range gs.games {
				game.Terminate <- true
			}
			return
		case gameID := <-gs.notifyMonitorEndGame:
			// notify the monitor
			for client := range gs.teacherClients {
				client.removeGame(gs.id + gameID)
			}

			// for demonstration session, we simply terminate the session
			// if there is no more games
			if gs.isDemo() && len(gs.games) == 0 {
				return
			}
		case summary := <-gs.monitor:
			summary.ID = gs.id + summary.ID
			for client := range gs.teacherClients {
				client.sendSummary(summary)
			}
		}
	}
}

// nbPlayers returns the number of players currently connected
func (gs *gameSession) nbPlayers() int {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	var allPlayers int
	for _, game := range gs.games {
		allPlayers += len(game.Summary().Successes)
	}
	return allPlayers
}

// groupIDs locks and returns the COMPLETE group ids
func (gs *gameSession) groupIDs() []string {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	var out []string
	for _, game := range gs.games {
		out = append(out, gs.id+game.ID)
	}

	sort.Strings(out)
	return out
}

func (gs *gameSession) generateName() string {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	allPlayers := make(map[tv.Player]bool)
	for _, game := range gs.games {
		for p := range game.Summary().Successes {
			allPlayers[p] = true
		}
	}

	nameFromID := func(s string) string {
		return fmt.Sprintf("Annonyme %s", s)
	}

	id := utils.RandomID(true, 5, func(s string) bool {
		return allPlayers[tv.Player{Name: nameFromID(s)}]
	})

	return nameFromID(id)
}

// teacherClient represents the teacher browser
type teacherClient struct {
	conn             *websocket.Conn
	currentSummaries map[GameID]tv.GameSummary
}

func (tc *teacherClient) socketData() (out teacherSocketData) {
	for _, su := range tc.currentSummaries {
		out.Games = append(out.Games, newGameSummary(su))
	}

	sort.Slice(out.Games, func(i, j int) bool {
		return out.Games[i].GameID < out.Games[j].GameID
	})

	return out
}

func (tc *teacherClient) removeGame(fullID string) {
	delete(tc.currentSummaries, fullID)
	tc.conn.WriteJSON(tc.socketData())
}

func (tc *teacherClient) sendSummary(gs tv.GameSummary) {
	// update the list
	tc.currentSummaries[gs.ID] = gs

	tc.conn.WriteJSON(tc.socketData())
}

// start loop listen for ping messages
func (tc *teacherClient) startLoop() {
	for {
		// read in a message
		_, _, err := tc.conn.ReadMessage()
		if err != nil {
			WarningLogger.Printf("teacher connection: %s", err)
			return
		}
	}
}
