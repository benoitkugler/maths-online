package trivialpoursuit

import (
	"fmt"
	"math/rand"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

// exposes some routes used to control and launch games

const GameEndPoint = "/trivial/game/:game_id"

type Controller struct {
	host  *url.URL                   // current URL start, such as http://localhost:1323, or https://www.deployed.fr
	games map[string]*gameController // active games
}

func NewController(host string) (*Controller, error) {
	u, err := url.Parse(host)
	return &Controller{
		host:  u,
		games: make(map[string]*gameController),
	}, err
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
	const choices = "ABCDEFGHIJKLMNOP"
	var out [6]byte
	for i := range out {
		out[i] = choices[rand.Intn(len(choices))]
	}
	return string(out[:])
}

func (ct *Controller) launchGame(options LaunchGameIn) LaunchGameOut {
	newID := randGameID()
	for _, taken := ct.games[newID]; taken; newID = randGameID() { // avoid (unlikely) collisions
	}

	game := newGameController(GameOptions{PlayersNumber: options.NbPlayers})
	// register the controller...
	ct.games[newID] = game
	// ...and start it
	go func() {
		game.startLoop()

		delete(ct.games, newID)
	}()

	var out LaunchGameOut
	out.URL = ct.buildURL(strings.ReplaceAll(GameEndPoint, ":game_id", newID), true)
	return out
}

func (ct *Controller) AccessGame(c echo.Context) error {
	gameID := c.Param("game_id")

	game, ok := ct.games[gameID]
	if !ok {
		return fmt.Errorf("La partie n'existe pas ou est déjà terminée")
	}

	game.setupWebSocket(c.Response().Writer, c.Request())
	return nil
}
