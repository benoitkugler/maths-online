package trivialpoursuit

import (
	"sort"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
)

func TestGameTimeout(t *testing.T) {
	// ProgressLogger.SetOutput(os.Stdout)
	const timeout = time.Second / 10

	ct := newGameSession("test", nil, TrivialConfig{}, RandomGroupStrategy{2, 2}, game.QuestionPool{})
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

func TestGameID(t *testing.T) {
	s := make([]string, 20)
	for i := range s {
		s[i] = gameIDFromSerial(i)
	}

	if !sort.StringsAreSorted(s) {
		t.Fatal("game ids are not sorted")
	}
}
