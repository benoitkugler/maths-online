package trivialpoursuit

// Game represents an on-going game,
// defined by an initial state and list of events.
type Game struct {
	events []event

	initialState GameState
}

func (g Game) currentState() GameState {
	out := g.initialState
	for _, event := range g.events {
		event.apply(&out)
	}
	return out
}

// newGameState returns an initial game state
func newGameState(nbPlayers int) GameState {
	return GameState{Successes: make([]success, nbPlayers)}
}

// winners returns the player who win
func (gs *GameState) winners() (out []int) {
	for player, success := range gs.Successes {
		if success.isDone() {
			out = append(out, player)
		}
	}
	return out
}

// the number of categories a player should complete
const nbCategories = 5

// the categories completed by a player
type success [nbCategories]bool

func (sc success) isDone() bool {
	for _, b := range sc {
		if !b {
			return false
		}
	}
	return true
}

// // the players who won the game
// type endGame []int
