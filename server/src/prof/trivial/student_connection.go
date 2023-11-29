package trivial

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/benoitkugler/maths-online/server/src/pass"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	tv "github.com/benoitkugler/maths-online/server/src/trivial"
	"github.com/benoitkugler/maths-online/server/src/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// handle connexion logic
// the client make a first http request
// the server return a (crypted) url containing all
// the parameters needed to acces the game room
// (or an error)
// the client then connect on an websocket with this payload.

// sessionID is a 4 digit identifier internaly used
// to map from teacher permanent ID to active game sessions
type sessionID = string

// gameConnection stores the information needed to access
// the proper game room and to reconnect in an already started game.
type gameConnection struct {
	GameID tv.RoomID

	// PlayerID is a generated UID used to
	// link incomming connections to internal game players,
	// used for reconnection.
	PlayerID tv.PlayerID

	// StudentID is the (crypted) ID of the student client.
	// It default to -1 for anonymous connections.
	StudentID pass.EncryptedID
}

type SetupStudentClientOut struct {
	GameMeta string
}

func firstBytes(s string, n int) string {
	if len(s) < n {
		return s
	}
	return s[:n]
}

// SetupStudentClient handles the connection of one student to the activity
// It is responsible for checking the credentials, creating games if needed,
// and returning the resolved game URL param used by `ConnectStudentSession`.
func (ct *Controller) SetupStudentClient(c echo.Context) error {
	completeID := c.QueryParam("session-id")
	clientID := c.QueryParam("client-id")
	gameMetaString := c.QueryParam("game-meta") // optional, used to reconnect

	ProgressLogger.Printf("Setup client <%s> for code %s (incomming meta: %s)", clientID, completeID, firstBytes(gameMetaString, 20))

	out, err := ct.setupStudentClient(completeID, clientID, gameMetaString)
	if err != nil {
		WarningLogger.Printf("setting up student: %s", err)
		return err
	}

	ProgressLogger.Printf("Setup done : %v", out)

	gameMeta, err := ct.studentKey.EncryptJSON(out)
	if err != nil {
		return fmt.Errorf("interal error: %s", err)
	}

	return c.JSON(200, SetupStudentClientOut{GameMeta: gameMeta})
}

type gameID interface {
	isGameID()
	String() string
}

func (demoCode) isGameID()       {}
func (teacherCode) isGameID()    {}
func (selfaccessCode) isGameID() {}

// <demoPin>.<room>.<number>
type demoCode struct {
	demoPin   string
	room      int // printed with 0 padding
	nbPlayers int
}

func (code demoCode) String() string {
	return fmt.Sprintf("%s.%02d.%d", code.demoPin, code.room, code.nbPlayers)
}

func (ct *Controller) createDemoGame(code demoCode) error {
	// check if the game is running and waiting for players
	ct.store.lock.Lock()
	_, ok := ct.store.games[code]
	ct.store.lock.Unlock()

	if !ok {
		// create a game on the fly
		questionPool, err := selectQuestions(ct.db, demoQuestions, ct.admin.Id, true)
		if err != nil {
			return err
		}

		nbSuccess := code.room % (tv.NbCategories + 1)

		options := tv.Options{
			Launch:          tv.LaunchStrategy{Manual: false, Max: code.nbPlayers},
			QuestionTimeout: time.Second * 120,
			ShowDecrassage:  true,
			Questions:       questionPool,
			StartNbSuccess:  nbSuccess,
		}

		ct.store.createGame(createGame{
			ID:      code,
			Options: options,
		})
	}

	ProgressLogger.Printf("Setting up student at (demo) %s", code.String())

	return nil
}

// <sessionID>.<gameID>, where sessionID is 4 digits
type teacherCode struct {
	sessionID sessionID
	gameID    string
}

func (code teacherCode) String() string {
	return fmt.Sprintf("%s.%s", code.sessionID, code.gameID)
}

// <gameID> (5 digits)
type selfaccessCode string

func (code selfaccessCode) String() string { return string(code) }

// parse a client game code, returning an error on invalid/malicious inputs
func (gs *gameStore) parseCode(clientGameCode string) (gameID, error) {
	cuts := strings.Split(clientGameCode, ".")
	switch len(cuts) {
	case 3: // demo
		if gs.demoPin != cuts[0] {
			return nil, fmt.Errorf("Code pin de démonstration %s invalide", cuts[0])
		}
		var out demoCode
		out.demoPin = cuts[0]
		if len(cuts[1]) < 2 {
			return nil, fmt.Errorf("Code de partie de démonstration %s invalide", cuts[1])
		}
		var err error
		out.room, err = strconv.Atoi(cuts[1])
		if err != nil {
			return nil, fmt.Errorf("Numéro de partie de démonstration %s invalide: %s", cuts[1], err)
		}
		out.nbPlayers, err = strconv.Atoi(cuts[2])
		if err != nil {
			return nil, fmt.Errorf("Nombre de joueurs %s invalide: %s", cuts[2], err)
		}
		return out, nil
	case 2: // teacher
		sID, gameID := cuts[0], cuts[1]
		if len(sID) < 4 || len(gameID) < 2 {
			return nil, fmt.Errorf("Code (type classe) %s invalide", clientGameCode)
		}
		return teacherCode{sessionID: sID, gameID: gameID}, nil
	case 1: // selfaccess
		if len(clientGameCode) != 5 {
			return nil, fmt.Errorf("Code (type perso.) %s invalide", clientGameCode)
		}
		return selfaccessCode(clientGameCode), nil
	default:
		return nil, fmt.Errorf("Code <%s> invalide", clientGameCode)
	}
}

// checkGameConnection checks if the given `meta` corresponds to a currently valid game.
// Clients may cache the returned connections so that expired data may be send back.
func (ct *gameStore) checkGameConnection(meta gameConnection) bool {
	gID, err := ct.parseCode(string(meta.GameID))
	if err != nil {
		return false
	}

	ct.lock.Lock()
	defer ct.lock.Unlock()

	if _, has := ct.games[gID]; !has {
		return false
	}

	gameID2, has := ct.playerIDs[meta.PlayerID]
	if !has {
		return false
	}
	// check that it is the correct game
	if gID != gameID2.game {
		return false
	}

	return true
}

func (ct *Controller) setupStudentClient(clientGameCode, clientID, gameMetaString string) (gameConnection, error) {
	if gameMetaString != "" {
		var incomingGameMeta gameConnection
		err := ct.studentKey.DecryptJSON(gameMetaString, &incomingGameMeta)
		if err != nil {
			return gameConnection{}, fmt.Errorf("internal error: %s", err)
		}

		ProgressLogger.Printf("checking client provided game meta: %v", incomingGameMeta)

		if ct.store.checkGameConnection(incomingGameMeta) {
			// simply the return the valid information
			return incomingGameMeta, nil
		}
	}

	// connect student according to the connection mode
	codeKind, err := ct.store.parseCode(clientGameCode)
	if err != nil {
		return gameConnection{}, err
	}

	clientID_ := pass.EncryptedID(clientID)
	// special case for demo
	if demo, isDemo := codeKind.(demoCode); isDemo {
		err = ct.createDemoGame(demo)
		if err != nil {
			return gameConnection{}, err
		}
		clientID_ = ""
	}

	return ct.store.setupStudent(clientID_, codeKind, ct.studentKey)
}

// setupStudent returns the game room meta data.
func (gs *gameStore) setupStudent(studentID pass.EncryptedID, requestedGameID gameID, key pass.Encrypter) (gameConnection, error) {
	gs.lock.Lock()
	game := gs.games[requestedGameID]
	gs.lock.Unlock()

	if game == nil {
		return gameConnection{}, fmt.Errorf("Code de salle %s invalide.", requestedGameID.String())
	}

	// if the game has already started, return an error early
	if game.HasStarted() {
		return gameConnection{}, fmt.Errorf("La partie %s a déjà commencée.", requestedGameID.String())
	}

	playerID := gs.registerPlayer(requestedGameID, studentID)
	out := gameConnection{
		GameID:    game.ID,
		PlayerID:  playerID,
		StudentID: studentID,
	}

	return out, nil
}

// ConnectStudentSession handles the connection of one student to the activity,
// using the meta data returned by a previous call to `SetupStudentClient`.
func (ct *Controller) ConnectStudentSession(c echo.Context) error {
	cryptedMeta := c.QueryParam("game-meta")
	clientPseudo := c.QueryParam("client-pseudo")

	var meta gameConnection
	err := ct.studentKey.DecryptJSON(cryptedMeta, &meta)
	if err != nil {
		return err
	}

	ProgressLogger.Printf("Connecting student %v", meta)

	err = ct.connectStudentTo(c, meta, clientPseudo)
	if err != nil {
		WarningLogger.Printf("connecting student: %s", err)
	}

	return err
}

type studentClient struct {
	// WS should be close when StartLoop ends
	WS   *websocket.Conn
	game *tv.Room // to accept user events

	playerID tv.PlayerID // used to handle reconnection and identifie client events
}

// listen listens for new messages being sent to our WebSocket
// endpoint, only returning on error.
// the connection is not closed yet
func (cl *studentClient) listen() {
	defer func() {
		cl.game.Leave <- cl.playerID
	}()

	for {
		// read in a message
		_, r, err := cl.WS.NextReader()
		if err, isClose := err.(*websocket.CloseError); isClose {
			ProgressLogger.Printf("Client left (%v)", err)
			return
		}
		if err != nil {
			WarningLogger.Printf("unexpected client error: %s", err)
			return
		}

		var event tv.ClientEventITFWrapper
		err = json.NewDecoder(r).Decode(&event)
		if err != nil {
			WarningLogger.Printf("invalid event format: %s", err)

			// return an error to the client and close
			utils.WebsocketError(cl.WS, err)
			return
		}

		// process the event
		cl.game.Event <- tv.ClientEvent{Event: event.Data, Player: cl.playerID}
	}
}

func (ct *Controller) connectStudentTo(c echo.Context, student gameConnection, pseudo string) error {
	gameID, err := ct.store.parseCode(string(student.GameID))
	if err != nil {
		return err
	}

	player := tv.Player{ID: student.PlayerID}
	var studentID int64 = -1
	if student.StudentID == "" { // anonymous connection
		player.Pseudo = pseudo
		if player.Pseudo == "" {
			player.Pseudo = ct.store.generateName() // finally generate a random pseudo
		}
	} else { // fetch name from DB
		studentID, err = ct.studentKey.DecryptID(student.StudentID)
		if err != nil {
			return fmt.Errorf("invalid student ID: %s", err)
		}

		student, err := teacher.SelectStudent(ct.db, teacher.IdStudent(studentID))
		if err != nil {
			return utils.SQLError(err)
		}

		player.Pseudo = student.Surname
		player.PseudoSuffix = student.Name
	}

	// then add the player
	ct.store.lock.Lock()
	game := ct.store.games[gameID]
	ct.store.lock.Unlock()

	if game == nil {
		return fmt.Errorf("internal error: invalid game ID %s", student.GameID)
	}

	// upgrade this connection to a WebSocket connection
	ws, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		WarningLogger.Println("internal error: failed to upgrade websocket: ", err)
		return nil
	}

	err = game.Join(player, ws) // check the access
	if err != nil {
		ProgressLogger.Printf("Rejecting connection for playerID %s to game %s: %s", student.PlayerID, game.ID, err)
		// the game at this end point is not usable: close the connection with an error
		utils.WebsocketError(ws, errors.New("game is closed"))
		ws.Close()

		return nil
	}

	client := &studentClient{
		game:     game,
		playerID: player.ID,
		WS:       ws,
	}

	client.listen() // block until client leaves

	ProgressLogger.Println("closing client connection", client.WS.RemoteAddr())
	client.WS.Close()

	return nil
}
