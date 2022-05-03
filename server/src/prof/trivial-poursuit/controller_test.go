package trivialpoursuit

import (
	"sort"
	"testing"

	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
)

func TestGameTermination(t *testing.T) {
	ct := newGameSession("test", nil, TrivialConfig{}, RandomGroupStrategy{2, 2}, game.QuestionPool{})

	id := ct.createGame(2)

	if len(ct.games) != 1 {
		t.Fatal("expected one game")
	}

	ct.games[id].Terminate <- true

	if len(ct.games) != 0 {
		t.Fatal("game should have been removed")
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
