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

// SetupStudentClient handles the connection of one student to the activity
// It is responsible for checking the credentials, creating games if needed,
// and returning the resolved game URL param used by `ConnectStudentSession`.
func (ct *Controller) SetupStudentClient(c echo.Context) error {
	completeID := c.QueryParam("session-id")
	clientID := c.QueryParam("client-id")
	gameMetaString := c.QueryParam("game-meta") // optional, used to reconnect

	out, err := ct.setupStudentClient(completeID, clientID, gameMetaString)
	if err != nil {
		return err
	}

	gameMeta, err := ct.studentKey.EncryptJSON(out)
	if err != nil {
		return fmt.Errorf("Erreur interne (%s)", err)
	}

	return c.JSON(200, SetupStudentClientOut{GameMeta: gameMeta})
}

type gameID interface {
	roomID() tv.RoomID
	setupStudentClient(ct *Controller, studentID pass.EncryptedID) (gameConnection, error)
}

// <demoPin>.<room>.<number>
type demoCode struct {
	demoPin   string
	room      string
	nbPlayers int
}

func (code demoCode) roomID() tv.RoomID {
	return tv.RoomID(fmt.Sprintf("%s.%s.%d", code.demoPin, code.room, code.nbPlayers))
}

func (code demoCode) setupStudentClient(ct *Controller, _ pass.EncryptedID) (gameConnection, error) {
	// check if the game is running and waiting for players
	ct.store.lock.Lock()
	_, ok := ct.store.games[code]
	ct.store.lock.Unlock()

	if !ok {
		// create a game on the fly
		questionPool, err := selectQuestions(ct.db, demoQuestions, ct.admin.Id)
		if err != nil {
			return gameConnection{}, err
		}

		options := tv.Options{
			Launch:          tv.LaunchStrategy{Manual: false, Max: code.nbPlayers},
			QuestionTimeout: time.Second * 120,
			ShowDecrassage:  true,
			Questions:       questionPool,
		}

		ct.store.createGame(createGame{
			ID:      code,
			Options: options,
		})
	}

	ProgressLogger.Printf("Setting up student at (demo) %s", code.roomID())

	return ct.store.setupStudent("", code, ct.studentKey)
}

// <sessionID>.<gameID>, where sessionID is 4 digits
type teacherCode struct {
	sessionID sessionID
	gameID    string
}

func (code teacherCode) roomID() tv.RoomID {
	return tv.RoomID(fmt.Sprintf("%s.%s", code.sessionID, code.gameID))
}

func (tc teacherCode) setupStudentClient(ct *Controller, studentID pass.EncryptedID) (gameConnection, error) {
	return ct.store.setupStudent(studentID, tc, ct.studentKey)
}

// <gameID> (5 digits)
type selfaccessCode string

func (code selfaccessCode) roomID() tv.RoomID { return tv.RoomID(code) }

func (code selfaccessCode) setupStudentClient(ct *Controller, clientID pass.EncryptedID) (gameConnection, error) {
	// TODO:
	return gameConnection{GameID: code.roomID()}, nil
}

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
		out.room = cuts[1]
		if len(out.room) < 2 {
			return nil, fmt.Errorf("Code de partie de démonstration %s invalide", cuts[1])
		}
		var err error
		out.nbPlayers, err = strconv.Atoi(cuts[2])
		if err != nil {
			return nil, fmt.Errorf("Numéro de partie de démonstration %s invalide: %s", cuts[2], err)
		}
		return out, nil
	case 2: // teacher
		sessionID, gameID := cuts[0], cuts[1]
		if len(sessionID) < 4 || len(gameID) < 2 {
			return nil, fmt.Errorf("Code (type classe) %s invalide", clientGameCode)
		}
		return teacherCode{sessionID: sessionID, gameID: gameID}, nil
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
	if gID != gameID2 {
		return false
	}

	return true
}

func (ct *Controller) setupStudentClient(clientGameCode, clientID, gameMetaString string) (gameConnection, error) {
	if gameMetaString != "" {
		var incomingGameMeta gameConnection
		err := ct.studentKey.DecryptJSON(gameMetaString, &incomingGameMeta)
		if err != nil {
			return gameConnection{}, err
		}

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

	return codeKind.setupStudentClient(ct, pass.EncryptedID(clientID))
}

// setupStudent returns the game room meta data.
func (gs *gameStore) setupStudent(studentID pass.EncryptedID, requestedGameID gameID, key pass.Encrypter) (gameConnection, error) {
	gs.lock.Lock()
	game := gs.games[requestedGameID]
	gs.lock.Unlock()

	if game == nil {
		return gameConnection{}, fmt.Errorf("Code de salle %s invalide.", requestedGameID.roomID())
	}

	playerID := gs.registerPlayer(requestedGameID)
	out := gameConnection{
		PlayerID:  playerID,
		StudentID: studentID,
		GameID:    game.ID,
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
		WarningLogger.Println("Failed to init websocket: ", err)
		return nil
	}

	err = game.Join(player, ws) // check the access
	if err != nil {
		ProgressLogger.Printf("Rejecting connection to game %s", game.ID)
		// the game at this end point is not usable: close the connection with an error
		utils.WebsocketError(ws, errors.New("game is closed"))
		ws.Close()

		return err
	}

	client := &studentClient{
		game:     game,
		playerID: player.ID,
		WS:       ws,
	}

	client.listen() // block until client leaves

	client.WS.Close()

	return nil
}
