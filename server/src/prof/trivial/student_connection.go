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

// SessionID is a 4 digit identifier internaly used
// to access one activity
// It is usually contained in a larger tv.RoomID
type SessionID = string

// gameConnection stores the information needed to access
// the proper game room and to reconnect in an already started game.
type gameConnection struct {
	SessionID SessionID
	GameID    tv.RoomID

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

	gameMeta, err := ct.key.EncryptJSON(out)
	if err != nil {
		return fmt.Errorf("Erreur interne (%s)", err)
	}

	return c.JSON(200, SetupStudentClientOut{GameMeta: gameMeta})
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

// checkGameConnection checks if the given `meta` corresponds to a currently valid game.
// Clients may cache the returned connections so that expired data may be send back.
func (ct *Controller) checkGameConnection(meta gameConnection) bool {
	ct.lock.Lock()
	defer ct.lock.Unlock()
	session, ok := ct.sessions[meta.SessionID]
	if !ok {
		return false
	}

	session.lock.Lock()
	defer session.lock.Unlock()

	if _, has := session.games[meta.GameID]; !has {
		return false
	}

	gameID, has := session.playerIDs[meta.PlayerID]
	if !has {
		return false
	}
	// check that it is the correct game
	if gameID != meta.GameID {
		return false
	}

	return true
}

func (ct *Controller) setupStudentClient(clientGameCode, clientID, gameMetaString string) (gameConnection, error) {
	if gameMetaString != "" {
		var incomingGameMeta gameConnection
		err := ct.key.DecryptJSON(gameMetaString, &incomingGameMeta)
		if err != nil {
			return gameConnection{}, err
		}

		if ct.checkGameConnection(incomingGameMeta) {
			// simply the return the valid information
			return incomingGameMeta, nil
		}
	}

	// special case for demonstration sessions
	if room, nbPlayers := ct.isDemoSessionID(clientGameCode); nbPlayers != 0 {
		return ct.setupStudentDemo(room, nbPlayers)
	}

	if len(clientGameCode) < 4 {
		return gameConnection{}, fmt.Errorf("Code %s invalide", clientGameCode)
	}
	sessionID := clientGameCode[:4]
	gameID := tv.RoomID(clientGameCode)

	ct.lock.Lock()
	session, ok := ct.sessions[sessionID]
	ct.lock.Unlock()
	if !ok {
		WarningLogger.Printf("invalid session ID %s", sessionID)
		return gameConnection{}, fmt.Errorf("L'activité n'existe pas ou est déjà terminée.")
	}

	studentID := pass.EncryptedID(clientID)

	return session.setupStudent(studentID, gameID, ct.key)
}

func (ct *Controller) setupStudentDemo(room string, nbPlayers int) (gameConnection, error) {
	sessionID := fmt.Sprintf("%s.%s.%d", ct.demoPin, room, nbPlayers)

	// check if the session is running and waiting for players
	ct.lock.Lock()
	session, ok := ct.sessions[sessionID]
	ct.lock.Unlock()

	if !ok {
		// create and launch the session ...
		session = ct.createSession(sessionID, -1)

		// ... and add one game on the fly
		questionPool, err := selectQuestions(ct.db, demoQuestions, ct.admin.Id)
		if err != nil {
			return gameConnection{}, err
		}

		options := tv.Options{
			Launch:          tv.LaunchStrategy{Manual: false, Max: nbPlayers},
			QuestionTimeout: time.Second * 120,
			ShowDecrassage:  true,
			Questions:       questionPool,
		}

		// we only build one game per session, so use the sessionID
		// as gameID for simplicity
		session.createGame(createGame{
			ID:      tv.RoomID(sessionID),
			Options: options,
		})
	}

	ProgressLogger.Printf("Setting up student at (demo) %s", sessionID)

	return session.setupStudent("", tv.RoomID(sessionID), ct.key)
}

// setupStudent returns the game room meta data.
func (gs *gameSession) setupStudent(studentID pass.EncryptedID, requestedGameID tv.RoomID, key pass.Encrypter) (gameConnection, error) {
	gs.lock.Lock()
	game := gs.games[requestedGameID]
	gs.lock.Unlock()

	if game == nil {
		return gameConnection{}, fmt.Errorf("Code de salle %s invalide.", requestedGameID)
	}

	playerID := gs.registerPlayer(game.ID)
	out := gameConnection{
		SessionID: gs.id,
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
	err := ct.key.DecryptJSON(cryptedMeta, &meta)
	if err != nil {
		return err
	}

	ct.lock.Lock()
	session, ok := ct.sessions[meta.SessionID]
	ct.lock.Unlock()
	if !ok {
		WarningLogger.Printf("unused session ID: %s", meta.SessionID)
		return fmt.Errorf("L'activité n'existe pas ou est déjà terminée.")
	}

	ProgressLogger.Printf("Connecting student %v", meta)

	err = ct.connectStudentTo(session, c, meta, clientPseudo)
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

func (ct *Controller) connectStudentTo(session *gameSession, c echo.Context, student gameConnection, pseudo string) error {
	player := tv.Player{ID: student.PlayerID}
	var studentID int64 = -1
	if student.StudentID == "" { // anonymous connection
		player.Pseudo = pseudo
		if player.Pseudo == "" {
			player.Pseudo = session.generateName() // finally generate a random pseudo
		}
	} else { // fetch name from DB
		var err error
		studentID, err = ct.key.DecryptID(student.StudentID)
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
	session.lock.Lock()
	game := session.games[student.GameID]
	session.lock.Unlock()

	if game == nil {
		return fmt.Errorf("invalid game ID %s", student.StudentID)
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
