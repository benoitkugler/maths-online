package trivial

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"
)

type clientOut struct {
	updates []StateUpdate
}

func (c *clientOut) WriteJSON(v interface{}) error {
	err := json.NewEncoder(io.Discard).Encode(v)
	c.updates = append(c.updates, v.(StateUpdate))
	return err
}

func (c *clientOut) lastU(lock *sync.Mutex) StateUpdate {
	lock.Lock()
	defer lock.Unlock()
	return c.updates[len(c.updates)-1]
}

// tests that the room may concurrently receive events,
// without race conditions
func TestConcurrentEvents(t *testing.T) {
	ProgressLogger.SetOutput(io.Discard) // hide verbose log

	r := NewRoom("<test>", Options{Launch: LaunchStrategy{Max: 4}}) // do not start a game to simplify
	go r.Listen(context.Background())

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
	r := NewRoom("<test>", Options{Launch: LaunchStrategy{Max: 3}})

	isNaturalEnd := make(chan bool)
	go func() {
		_, ok := r.Listen(context.Background())
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

	last := c.lastU(&r.lock)
	if len(last.Events) != 1 {
		t.Fatal(last.Events)
	}
	if _, ok := last.Events[0].(GameTerminated); !ok {
		t.Fatalf("expected GameTerminated, got %v (%T)", last.Events[0], last.Events[0])
	}
}

func (r *Room) decoReco(player PlayerID, errC chan<- error) {
	time.Sleep(20 * time.Millisecond)
	r.Leave <- player

	time.Sleep(time.Millisecond)

	err := r.Join(Player{ID: player}, &clientOut{})
	errC <- err
}

func TestReconnection(t *testing.T) {
	r := NewRoom("<test>", Options{Launch: LaunchStrategy{Max: 3}})
	go r.Listen(context.Background())

	r.mustJoin(t, "p1")
	r.mustJoin(t, "p2")
	r.mustJoin(t, "p3")

	// simulate a deconection
	errC1 := make(chan error)
	errC2 := make(chan error)

	go r.decoReco("p1", errC1)
	go r.decoReco("p2", errC2)

	if err := <-errC1; err != nil {
		t.Fatal(err)
	}
	if err := <-errC2; err != nil {
		t.Fatal(err)
	}

	if pl := r.lp(); len(pl) != 3 {
		t.Fatal()
	}
}
