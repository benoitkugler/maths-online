package trivialpoursuit

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/benoitkugler/maths-online/pass"
	"github.com/labstack/echo/v4"
)

// handle connexion logic
// the client make a first http request
// the server return a (crypted) url containing all
// the parameters needed to acces the game room
// (or an error)
// the client then connect on an websocket with this payload.

// gameConnection stores the information needed to access
// the proper game room and to reconnect in an already started game.
type gameConnection struct {
	SessionID SessionID
	GameID    GameID

	// PlayerID is a generated UID used to
	// link incomming connections to internal game players,
	// used for reconnection.
	PlayerID string

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

	if _, has := session.playerIDs[meta.PlayerID]; !has {
		return false
	}

	return true
}

func (ct *Controller) setupStudentClient(clientSessionID, clientID, gameMetaString string) (gameConnection, error) {
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
	if room, nbPlayers := ct.isDemoSessionID(clientSessionID); nbPlayers != 0 {
		return ct.setupStudentDemo(room, nbPlayers)
	}

	if len(clientSessionID) < 4 {
		return gameConnection{}, fmt.Errorf("Code %s invalide", clientSessionID)
	}
	sessionID := clientSessionID[:4]
	gameID := clientSessionID[4:]

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
		},
			ct.admin.Id,
		)
		if err != nil {
			return gameConnection{}, fmt.Errorf("Erreur interne : %s", err)
		}
	}

	ProgressLogger.Printf("Setting up student at (demo) %s", sessionID)

	return session.setupStudent("", "", ct.key)
}

// setupStudent create a game if needed and return the game
// room meta data.
func (gs *gameSession) setupStudent(studentID pass.EncryptedID, requestedGameID GameID, key pass.Encrypter) (gameConnection, error) {
	playerID := gs.registerPlayer()

	out := gameConnection{
		SessionID: gs.id,
		PlayerID:  playerID,
		StudentID: studentID,
	}
	var studentIDInt int64 = -1
	if studentID != "" { // decode ID
		var err error
		studentIDInt, err = key.DecryptID(studentID)
		if err != nil {
			return out, fmt.Errorf("ID personnel %s invalide (%s).", studentID, err)
		}
	}

	// select or create the game
	gameID, err := gs.group.selectOrCreateGame(requestedGameID, studentIDInt, gs)
	if err != nil {
		return out, fmt.Errorf("Code de salle %s invalide (%s).", requestedGameID, err)
	}
	out.GameID = gameID

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

	err = session.connectStudent(c, meta, clientPseudo, ct.key)
	if err != nil {
		WarningLogger.Printf("connecting student: %s", err)
	}

	return err
}
