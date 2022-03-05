package trivialpoursuit

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// exposes some routes used to control and launch games

const (
	GameEndPoint = "/trivial/game/:game_id"
	gameTimeout  = time.Hour * 24
)

type gameID = string

type Controller struct {
	host *url.URL // current URL start, such as http://localhost:1323, or https://www.deployed.fr

	lock sync.Mutex

	games       map[gameID]*gameController // active games
	gameTimeout time.Duration              // max duration of a game (useful is nobody joins)
}

// NewController initialize the controller. It will panic
// if the given host is an invalid url
func NewController(host string) *Controller {
	u, err := url.Parse("http://" + host)
	if err != nil {
		panic(err)
	}
	return &Controller{
		host:        u,
		games:       make(map[string]*gameController),
		gameTimeout: gameTimeout,
	}
}

// return the active game and the number of players in it
func (ct *Controller) stats() map[gameID]int {
	ct.lock.Lock()
	defer ct.lock.Unlock()

	out := make(map[gameID]int)
	for k, v := range ct.games {
		out[k] = len(v.clients)
	}
	return out
}

func (ct *Controller) buildURL(path string, isWebSocket bool) string {
	u := *ct.host
	u.RawPath = path
	u.Path = path

	if isWebSocket {
		u.Scheme = "ws"
	}

	return u.String()
}

type LaunchGameIn struct {
	NbPlayers int
}

type LaunchGameOut struct {
	URL string
}

// LaunchGame starts a new game and return the WebSocket URL to
// access it.
func (ct *Controller) LaunchGame(c echo.Context) error {
	var in LaunchGameIn
	if err := c.Bind(&in); err != nil {
		return fmt.Errorf("invalid parameters format: %s", err)
	}

	out := ct.launchGame(in)

	return c.JSON(200, out)
}

func randGameID() string {
	const choices = "abcdefghijklmnopqrstuvwxyz0123456789"
	var out [6]byte
	for i := range out {
		out[i] = choices[rand.Intn(len(choices))]
	}
	return string(out[:])
}

func (ct *Controller) launchGame(options LaunchGameIn) LaunchGameOut {
	ct.lock.Lock()
	defer ct.lock.Unlock()

	newID := randGameID()
	for _, taken := ct.games[newID]; taken; newID = randGameID() { // avoid (unlikely) collisions
	}

	game := newGameController(GameOptions{PlayersNumber: options.NbPlayers})
	// register the controller...
	ct.games[newID] = game
	// ...and start it
	go func() {
		ctx, cancelFunc := context.WithTimeout(context.Background(), ct.gameTimeout)
		game.startLoop(ctx)

		cancelFunc()

		// remove the game controller when the game is over
		ct.lock.Lock()
		defer ct.lock.Unlock()
		delete(ct.games, newID)
	}()

	var out LaunchGameOut
	out.URL = ct.buildURL(strings.ReplaceAll(GameEndPoint, ":game_id", newID), true)
	return out
}

// AccessGame establish a connection to a game session, using
// WebSockets
func (ct *Controller) AccessGame(c echo.Context) error {
	gameID := c.Param("game_id")

	game, ok := ct.games[gameID]
	if !ok {
		WarningLogger.Printf("invalid game ID %s", gameID)
		return fmt.Errorf("La partie n'existe pas ou est déjà terminée.")
	}

	// connect to the websocket handler, which handle errors
	game.setupWebSocket(c.Response().Writer, c.Request())
	return nil
}

// ShowStats returns a JSON dict of current games, usable as
// debug tool.
// TODO: protect the access
func (ct *Controller) ShowStats(c echo.Context) error {
	return c.JSON(200, ct.stats())
}
