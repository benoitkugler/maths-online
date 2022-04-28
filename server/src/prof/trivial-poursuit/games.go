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

// GameID is an in-memory identifier for a game room.
type GameID = string

// gameSession monitor the games of one session (think one classroom)
// and broadcast the main advances from all the games to the teacher client
type gameSession struct {
	quit chan bool

	lock sync.Mutex

	db *sql.DB

	teacherClients map[*teacherClient]bool
	monitor        chan tv.GameSummary
	questions      ga.QuestionPool
	games          map[GameID]*tv.GameController // active games, either in lobby or playing

	config TrivialConfig // database entry, cached for simplicity
	group  GroupStrategy // specified when starting the session
}

func newGameSession(db *sql.DB, config TrivialConfig, group GroupStrategy, questions ga.QuestionPool) *gameSession {
	return &gameSession{
		db:             db,
		config:         config,
		group:          group,
		quit:           make(chan bool),
		monitor:        make(chan tv.GameSummary),
		games:          make(map[string]*tv.GameController),
		teacherClients: make(map[*teacherClient]bool),
		questions:      questions,
	}
}

func (gs *gameSession) createGame(nbPlayers int) GameID {
	gs.lock.Lock()
	defer gs.lock.Unlock()

	gameID := utils.RandomID(false, 10, func(s string) bool {
		_, taken := gs.games[s]
		return taken
	})

	game := tv.NewGameController(gameID,
		gs.questions,
		tv.GameOptions{
			PlayersNumber:   nbPlayers,
			QuestionTimeout: time.Second * time.Duration(gs.config.QuestionTimeout),
		},
		gs.monitor)
	// register the controller...
	gs.games[gameID] = game
	// ...and start it
	go func() {
		ctx, cancelFunc := context.WithTimeout(context.Background(), gameTimeout)
		review, ok := game.StartLoop(ctx)
		if ok { // exploit the review
			gs.exploitReview(review)
		}

		cancelFunc()

		// remove the game controller when the game is over
		gs.lock.Lock()
		defer gs.lock.Unlock()
		delete(gs.games, gameID)
	}()

	return gameID
}

// TODO:
func (gs *gameSession) exploitReview(review tv.Review) {
	ProgressLogger.Printf("GAME REVIEW: %v", review)
}

func (gs *gameSession) connectStudent(c echo.Context, clientID pass.EncryptedID, key pass.Encrypter) error {
	player := tv.Player{ID: clientID}
	var studentID int64 = -1
	if clientID == "" { // anonymous connection
		player.Name = gs.generateName()
	} else { // fetch name from DB
		var err error
		studentID, err = clientID.Decrypt(key)
		if err != nil {
			return fmt.Errorf("invalid student ID: %s", err)
		}

		student, err := students.SelectStudent(gs.db, studentID)
		if err != nil {
			return utils.SQLError(err)
		}

		player.Name = student.Surname
	}

	// select the game
	gameID, newGamePlayers := gs.group.selectGame(studentID, gs)

	if gameID == "" { // first create a new game

		ProgressLogger.Printf("Creating game room (for %d players)", newGamePlayers)

		gameID = gs.createGame(newGamePlayers)
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

	for k, ga := range gs.games {
		client.currentSummaries[k] = ga.Summary()
	}
	gs.teacherClients[client] = true

	client.conn.WriteJSON(client.socketData())

	return client
}

func (gs *gameSession) startLoop(ctx context.Context) {
	for {
		select {
		case <-gs.quit:
			return
		case <-ctx.Done():
			return
		case summary := <-gs.monitor:
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
