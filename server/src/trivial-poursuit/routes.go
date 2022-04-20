package trivialpoursuit

// exposes some routes used to control and launch games

// type Controller struct {
// 	host *url.URL // current URL start, such as http://localhost:1323, or https://www.deployed.fr

// 	lock sync.Mutex

// 	games       map[GameID]*GameController // active games
// 	gameTimeout time.Duration              // max duration of a game (useful is nobody joins)

// 	monitor Monitor
// }

// // NewController initialize the controller. It will panic
// // if the given host is an invalid url
// func NewController(host string) *Controller {
// 	u, err := url.Parse("http://" + host)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return &Controller{
// 		host:        u,
// 		games:       make(map[string]*GameController),
// 		gameTimeout: gameTimeout,
// 	}
// }

// // return the active game and the number of players in it
// func (ct *Controller) stats() map[GameID]int {
// 	ct.lock.Lock()
// 	defer ct.lock.Unlock()

// 	out := make(map[GameID]int)
// 	for k, v := range ct.games {
// 		out[k] = len(v.clients)
// 	}
// 	return out
// }

// func (ct *Controller) buildURL(path string, isWebSocket bool) string {
// 	u := *ct.host
// 	u.RawPath = path
// 	u.Path = path

// 	if isWebSocket {
// 		u.Scheme = "ws"
// 	}

// 	return u.String()
// }

// type LaunchGameIn struct {
// 	NbPlayers      int
// 	TimeoutSeconds int
// }

// type LaunchGameOut struct {
// 	URL string
// }

// // LaunchGame starts a new game and return the WebSocket URL to
// // access it.
// func (ct *Controller) LaunchGame(c echo.Context) error {
// 	var in LaunchGameIn
// 	if err := c.Bind(&in); err != nil {
// 		return fmt.Errorf("invalid parameters format: %s", err)
// 	}

// 	out := ct.launchGame(in)

// 	return c.JSON(200, out)
// }

// func (ct *Controller) launchGame(options LaunchGameIn) LaunchGameOut {
// 	ct.lock.Lock()
// 	defer ct.lock.Unlock()

// 	newID := utils.RandomID(true, 3, func(s string) bool {
// 		_, taken := ct.games[s]
// 		return taken
// 	})

// 	game := NewGameController(
// 		newID,
// 		GameOptions{
// 			PlayersNumber:   options.NbPlayers,
// 			QuestionTimeout: time.Second * time.Duration(options.TimeoutSeconds),
// 		},
// 		ct.monitor)
// 	// register the controller...
// 	ct.games[newID] = game
// 	// ...and start it
// 	go func() {
// 		ctx, cancelFunc := context.WithTimeout(context.Background(), ct.gameTimeout)
// 		game.StartLoop(ctx)

// 		cancelFunc()

// 		// remove the game controller when the game is over
// 		ct.lock.Lock()
// 		defer ct.lock.Unlock()
// 		delete(ct.games, newID)
// 	}()

// 	var out LaunchGameOut
// 	out.URL = ct.buildURL(strings.ReplaceAll(GameEndPoint, ":game_id", newID), true)
// 	return out
// }

// // AccessGame establish a connection to a game session, using
// // WebSockets
// func (ct *Controller) AccessGame(c echo.Context) error {
// 	gameID := c.Param("game_id")

// 	game, ok := ct.games[gameID]
// 	if !ok {
// 		WarningLogger.Printf("invalid game ID %s", gameID)
// 		return fmt.Errorf("La partie n'existe pas ou est déjà terminée.")
// 	}

// 	// connect to the websocket handler, which handle errors
// 	clientID := c.QueryParam("client_id")
// 	game.AddClient(c.Response().Writer, c.Request(), Player{ID: pass.EncryptedID(clientID), Name: ""}) // TODO: name
// 	return nil
// }

// // ShowStats returns a JSON dict of current games, usable as
// // debug tool.
// // TODO: protect the access
// func (ct *Controller) ShowStats(c echo.Context) error {
// 	return c.JSON(200, ct.stats())
// }
