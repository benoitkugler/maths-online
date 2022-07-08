package trivial

import (
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/trivial-poursuit/game"
)

type client struct {
	events []interface{}
}

func (c *client) WriteJSON(v interface{}) error {
	err := json.NewEncoder(io.Discard).Encode(v)
	c.events = append(c.events, v)
	return err
}

func TestTerminate(t *testing.T) {
	r := NewRoom("<test>", Options{PlayersNumber: 3})

	isNaturalEnd := true
	go func() {
		_, isNaturalEnd = r.Listen()
	}()

	c := client{}
	if err := r.Join(Player{ID: "<p1>"}, &c); err != nil {
		t.Fatal(err)
	}
	if err := r.Join(Player{ID: "<p2>"}, &client{}); err != nil {
		t.Fatal(err)
	}

	r.Terminate <- true

	time.Sleep(time.Millisecond)

	if isNaturalEnd {
		t.Fatal("expected forced end")
	}

	L := len(c.events)
	last := c.events[L-1].(game.StateUpdate)
	if len(last.Events) != 1 {
		t.Fatal(last.Events)
	}
	if _, ok := last.Events[0].(game.GameTerminated); !ok {
		t.Fatalf("expected GameTerminated, got %v (%T)", last.Events[0], last.Events[0])
	}
}
