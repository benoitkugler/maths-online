package trivialpoursuit

import (
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
)

func TestGameTimeout(t *testing.T) {
	// ProgressLogger.SetOutput(os.Stdout)
	const timeout = time.Second / 10

	ct := newGameSession(nil, TrivialConfig{}, RandomGroupStrategy{2, 2}, game.QuestionPool{})
	gameTimeout = timeout

	ct.createGame(2)

	if len(ct.games) != 1 {
		t.Fatal("expected one game")
	}

	time.Sleep(2 * timeout) // wait for the timeout

	if len(ct.games) != 0 {
		t.Fatal("game should have timed out")
	}
}
