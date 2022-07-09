package trivial

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"
	"time"
)

type clientOut struct {
	events []interface{}
}

func (c *clientOut) WriteJSON(v interface{}) error {
	err := json.NewEncoder(io.Discard).Encode(v)
	c.events = append(c.events, v)
	return err
}

func (c *clientOut) reset() {
	c.events = c.events[:0]
}

// tests that the room may concurrently receive events,
// without race conditions
func TestConcurrentEvents(t *testing.T) {
	ProgressLogger.SetOutput(io.Discard) // hide verbose log

	r := NewRoom("<test>", Options{PlayersNumber: 4}) // do not start a game to simplify
	go r.Listen()

	var client1, client2, client3 clientOut
	if err := r.Join(Player{ID: "<p1>"}, &client1); err != nil {
		t.Fatal(err)
	}
	if err := r.Join(Player{ID: "<p2>"}, &client2); err != nil {
		t.Fatal(err)
	}
	if err := r.Join(Player{ID: "<p3>"}, &client3); err != nil {
		t.Fatal(err)
	}

	const nbSend = 100
	clientLoop := func() {
		for i := range [nbSend]int{} {
			r.Event <- ClientEvent{Event: Ping{Info: fmt.Sprintf("message %d", i)}}
		}
	}

	go clientLoop()
	go clientLoop()
	go clientLoop()

	time.Sleep(time.Second / 10)
}

func TestTerminate(t *testing.T) {
	r := NewRoom("<test>", Options{PlayersNumber: 3})

	isNaturalEnd := make(chan bool)
	go func() {
		_, ok := r.Listen()
		isNaturalEnd <- ok
	}()

	c := clientOut{}
	if err := r.Join(Player{ID: "<p1>"}, &c); err != nil {
		t.Fatal(err)
	}
	if err := r.Join(Player{ID: "<p2>"}, &clientOut{}); err != nil {
		t.Fatal(err)
	}

	r.Terminate <- true

	if <-isNaturalEnd {
		t.Fatal("expected forced end")
	}

	L := len(c.events)
	last := c.events[L-1].(StateUpdate)
	if len(last.Events) != 1 {
		t.Fatal(last.Events)
	}
	if _, ok := last.Events[0].(GameTerminated); !ok {
		t.Fatalf("expected GameTerminated, got %v (%T)", last.Events[0], last.Events[0])
	}
}
