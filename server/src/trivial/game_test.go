package trivial

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	"github.com/benoitkugler/maths-online/server/src/sql/editor"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func init() {
	GameStartDelay = time.Millisecond
}

var (
	exQu = WeigthedQuestions{
		Questions: []editor.Question{{Id: 1}, {Id: 2}, {Id: 3}, {Id: 4}},
		Weights:   []float64{1. / 4, 1. / 4, 1. / 4, 1. / 4},
	}
	exPool = QuestionPool{exQu, exQu, exQu, exQu, exQu}
)

func playersFromSuccess(scs ...Success) map[serial]*playerConn {
	out := make(map[serial]*playerConn)
	for i, s := range scs {
		id := serial(fmt.Sprintf("%d", i))
		out[id] = &playerConn{pl: Player{ID: id}, advance: playerAdvance{success: s}}
	}
	return out
}

func Test_gameState_winners(t *testing.T) {
	tests := []struct {
		sc      map[serial]*playerConn
		wantOut []serial
	}{
		{playersFromSuccess(Success{true}, Success{true}, Success{true, true, true, true, true}), []serial{"2"}},
		{playersFromSuccess(Success{true}, Success{true}, Success{true, true, true, true}), nil},
		{playersFromSuccess(Success{true, true, true, true, true}, Success{true}, Success{true, true, true, true, true}), []serial{"0", "2"}},
	}
	for _, tt := range tests {
		r := Room{
			players: tt.sc,
		}
		if gotOut := r.winners(); !reflect.DeepEqual(gotOut, tt.wantOut) {
			t.Errorf("gameState.winners() = %v, want %v", gotOut, tt.wantOut)
		}
	}
}

func playersFromIds(scs ...string) map[serial]*playerConn {
	out := make(map[serial]*playerConn)
	for _, id := range scs {
		out[serial(id)] = &playerConn{pl: Player{ID: serial(id)}, conn: &clientOut{}}
	}
	return out
}

func TestGameState_nextPlayer(t *testing.T) {
	tests := []struct {
		players map[serial]*playerConn
		current serial
		want    serial
	}{
		{
			players: playersFromIds("0", "1", "4"),
			current: "0",
			want:    "1",
		},
		{
			players: playersFromIds("0", "1", "4"),
			current: "1",
			want:    "4",
		},
		{
			players: playersFromIds("0", "1", "4"),
			current: "2",
			want:    "4",
		},
		{
			players: playersFromIds("0", "1", "4"),
			current: "4",
			want:    "0",
		},
	}
	for _, tt := range tests {
		r := Room{
			game:    game{playerTurn: tt.current},
			players: tt.players,
		}
		if got := r.nextPlayer(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("GameState.nextPlayer() = %v, want %v", got, tt.want)
		}
	}
}

// use a new connection
func (r *Room) mustJoin(t *testing.T, id PlayerID) {
	t.Helper()

	r.mustJoinConn(t, id, &clientOut{})
}

func (r *Room) mustJoinConn(t *testing.T, id PlayerID, client Connection) {
	t.Helper()

	if err := r.Join(Player{ID: id}, client); err != nil {
		t.Fatal(err)
	}
}

func (r *Room) throwAndMove(player PlayerID) {
	r.Event <- ClientEvent{Event: DiceClicked{}, Player: player}
	r.Event <- ClientEvent{Event: ClientMove{Tile: int(r.lg().dice.Face)}, Player: player}
}

func (r *Room) lg() game {
	time.Sleep(50 * time.Microsecond)
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.game
}

func (r *Room) lp() map[serial]*playerConn {
	time.Sleep(50 * time.Microsecond)
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.players
}

func TestStartAuto(t *testing.T) {
	r := NewRoom("0", Options{Launch: LaunchStrategy{Max: 3}}, noOpSuccesHandler{})
	go r.Listen(context.Background())

	err := r.StartGame()
	tu.Assert(t, err != nil) // auto mode

	r.mustJoin(t, "0")
	r.mustJoin(t, "1")

	if g := r.lg(); g.phase != pGameLobby {
		t.Fatalf("unexpected game phase %v", r.game.phase)
	}

	r.Leave <- "1"

	if players := r.lp(); len(players) != 1 {
		t.Fatalf("unexpected number of players %d", len(players))
	}

	// check that invalid ID is just a no op and not a crash
	r.Leave <- "bad ID"
	if players := r.lp(); len(players) != 1 {
		t.Fatalf("unexpected number of players %d", len(players))
	}

	if err := r.Join(Player{ID: "1"}, &clientOut{}); err != nil {
		t.Fatal(err)
	}
	if err := r.Join(Player{ID: "2"}, &clientOut{}); err != nil {
		t.Fatal(err)
	}

	if g := r.lg(); g.phase != pTurnStarted {
		t.Fatalf("unexpected game phase %v", r.game.phase)
	}

	if err := r.Join(Player{ID: "4"}, &clientOut{}); err == nil {
		t.Fatal("expected error on joining already started game")
	}
}

func TestStartManual(t *testing.T) {
	r := NewRoom("0", Options{Launch: LaunchStrategy{Manual: true}}, noOpSuccesHandler{})
	go r.Listen(context.Background())

	err := r.StartGame()
	tu.Assert(t, err != nil) // no players

	r.mustJoin(t, "0")

	r.Leave <- "0"

	time.Sleep(time.Millisecond)

	err = r.StartGame()
	tu.Assert(t, err != nil) // no players

	r.mustJoin(t, "0")
	r.mustJoin(t, "1")
	r.mustJoin(t, "2")

	g := r.lg()
	tu.Assert(t, g.phase == pGameLobby)

	err = r.StartGame()
	tu.AssertNoErr(t, err)

	g = r.lg()
	tu.Assert(t, g.phase == pTurnStarted)

	err = r.StartGame()
	tu.Assert(t, err != nil) // game already started
}

func TestEmitQuestion(t *testing.T) {
	r := NewRoom("", Options{Launch: LaunchStrategy{Max: 1}, Questions: exPool, QuestionTimeout: time.Minute}, noOpSuccesHandler{})
	go r.Listen(context.Background())

	r.mustJoin(t, "p1")

	if g := r.lg(); g.questionTimer.Stop() {
		t.Fatal("timer should not being running")
	}

	r.throwAndMove("p1")

	if g := r.lg(); g.phase != pDoingQuestion {
		t.Fatalf("unexpected phase %v", g.phase)
	}

	r.Event <- ClientEvent{Event: Answer{}, Player: "p1"}

	if g := r.lg(); g.phase != pQuestionResult {
		t.Fatalf("unexpected phase %v", g.phase)
	}
	if g := r.lg(); g.questionTimer.Stop() {
		t.Fatal("timer should not being running")
	}
}

func TestQuestionTimeout(t *testing.T) {
	r := NewRoom("", Options{Launch: LaunchStrategy{Max: 1}, Questions: exPool, QuestionTimeout: 5 * time.Millisecond}, noOpSuccesHandler{})

	r.mustJoin(t, "p1")
	r.players["p1"].advance.success = Success{true, true, true, true, true}

	go r.Listen(context.Background())

	r.throwAndMove("p1")

	qu := r.lg().question
	if g := r.lg(); g.phase != pDoingQuestion {
		t.Fatalf("unexpected phase %v", g.phase)
	}

	time.Sleep(10 * time.Millisecond) // trigger question timeout

	if g := r.lg(); g.phase != pQuestionResult {
		t.Fatalf("unexpected phase %v", g.phase)
	}

	if pl := r.lp(); pl["p1"].advance.success[qu.Categorie] {
		t.Fatal("success must have been lost")
	}
	if pl := r.lp(); len(pl["p1"].advance.review.QuestionHistory) != 1 {
		t.Fatal("missing question in history")
	}
}

func TestDisconnectStartTurn(t *testing.T) {
	r := NewRoom("", Options{Launch: LaunchStrategy{Max: 2}}, noOpSuccesHandler{})
	go r.Listen(context.Background())

	p2 := &clientOut{}

	r.mustJoin(t, "p1")
	r.mustJoinConn(t, "p2", p2)

	r.Leave <- "p1"
	time.Sleep(time.Millisecond)

	events := p2.lastU(&r.lock).Events
	if len(events) != 2 {
		t.Fatal()
	}
	resetTurn := events[1].(PlayerTurn)
	if resetTurn.Player != "p2" {
		t.Fatal(resetTurn)
	}
}

func TestDisconnectInQuestion(t *testing.T) {
	r := NewRoom("", Options{Launch: LaunchStrategy{Max: 2}, Questions: exPool, QuestionTimeout: time.Minute}, noOpSuccesHandler{})

	go r.Listen(context.Background())

	r.mustJoin(t, "p1")
	r.mustJoin(t, "p2")

	r.throwAndMove("p1")

	r.Event <- ClientEvent{Event: Answer{}, Player: "p1"}

	// player1 has answered, waiting for player2
	r.Leave <- "p2"

	time.Sleep(time.Millisecond)

	if g := r.lg(); g.questionTimer.Stop() {
		t.Fatal("question should have been closed")
	}
	if g := r.lg(); g.phase != pQuestionResult {
		t.Fatalf("unexpected phase %v", g.phase)
	}
}

func TestDisconnectInQuestionResult(t *testing.T) {
	r := NewRoom("", Options{Launch: LaunchStrategy{Max: 3}, Questions: exPool, QuestionTimeout: time.Minute}, noOpSuccesHandler{})

	go r.Listen(context.Background())

	r.mustJoin(t, "p1")
	r.mustJoin(t, "p2")
	r.mustJoin(t, "p3")

	r.throwAndMove("p1")

	r.Event <- ClientEvent{Event: Answer{}, Player: "p1"}
	r.Event <- ClientEvent{Event: Answer{}, Player: "p2"}
	r.Event <- ClientEvent{Event: Answer{}, Player: "p3"}

	time.Sleep(time.Millisecond)

	r.Event <- ClientEvent{Event: WantNextTurn{}, Player: "p1"}
	r.Event <- ClientEvent{Event: WantNextTurn{}, Player: "p3"}

	// player1 wants next turn, waiting for player2
	if g := r.lg(); g.playerTurn != "p1" {
		t.Fatal(g.playerTurn)
	}

	r.Leave <- "p2"

	if g := r.lg(); g.playerTurn != "p3" || g.phase != pTurnStarted {
		t.Fatal(g)
	}
}

func TestAllDisconnectInQuestionResult(t *testing.T) {
	r := NewRoom("<test>", Options{Launch: LaunchStrategy{Max: 2}, Questions: exPool, QuestionTimeout: time.Minute}, noOpSuccesHandler{})

	go r.Listen(context.Background())

	r.mustJoin(t, "p1")
	r.mustJoin(t, "p2")

	r.throwAndMove("p1")

	r.Event <- ClientEvent{Event: Answer{}, Player: "p1"}
	r.Event <- ClientEvent{Event: Answer{}, Player: "p2"}

	time.Sleep(time.Millisecond)

	r.Leave <- "p1"

	time.Sleep(time.Millisecond)

	r.Leave <- "p2"

	// no more players, the game is "staled"

	if g := r.lg(); g.phase != pQuestionResult {
		t.Fatal(g)
	}

	// upon reconnection, starts a new turn
	r.mustJoin(t, "p2")

	if g := r.lg(); g.phase != pTurnStarted || g.playerTurn != "p2" {
		t.Fatal(g)
	}
}

func TestReconnectInQuestion(t *testing.T) {
	r := NewRoom("", Options{Launch: LaunchStrategy{Max: 2}, Questions: exPool, QuestionTimeout: 3 * time.Second}, noOpSuccesHandler{})

	go r.Listen(context.Background())

	r.mustJoin(t, "p1")
	r.mustJoin(t, "p2")

	r.throwAndMove("p1")

	// disconnect p2 ...
	r.Leave <- "p2"

	time.Sleep(200 * time.Millisecond)

	// ...and reconnect it before p1 has answered
	p2 := &clientOut{}
	r.mustJoinConn(t, "p2", p2)
	question2 := p2.lastU(&r.lock).Events[0].(ShowQuestion)
	if question2.TimeoutSeconds != 2 {
		t.Fatal(question2)
	}

	r.Event <- ClientEvent{Event: Answer{}, Player: "p1"}
	r.Event <- ClientEvent{Event: Answer{}, Player: "p2"}

	time.Sleep(time.Millisecond)

	if g := r.lg(); g.questionTimer.Stop() {
		t.Fatal("question should have been closed")
	}
	if g := r.lg(); g.phase != pQuestionResult {
		t.Fatalf("unexpected phase %v", g.phase)
	}
}

func TestHandleClientEvent(t *testing.T) {
	r := NewRoom("", Options{Launch: LaunchStrategy{Max: 2}, Questions: exPool, QuestionTimeout: time.Minute}, noOpSuccesHandler{})

	fmt.Println("Phase :", r.game.phase)

	r.mustJoin(t, "p1")
	r.mustJoin(t, "p2")

	up, isOver, err := r.handleClientEvent(Ping{}, Player{ID: "p1"})
	if err != nil {
		t.Fatal(err)
	}
	if up != nil || isOver {
		t.Fatal("Ping should be ignored")
	}

	fmt.Println("Phase :", r.game.phase)

	_, _, err = r.handleClientEvent(nil, Player{ID: "p1"})
	if err == nil {
		t.Fatal("expected error for invalid client event type")
	}
	_, _, err = r.handleClientEvent(DiceClicked{}, Player{ID: "p2"})
	if err == nil {
		t.Fatal("expected error for incorrect player")
	}
	_, _, err = r.handleClientEvent(ClientMove{}, Player{ID: "p1"})
	if err == nil {
		t.Fatal("expected error for inconsitent phase")
	}

	_, _, _ = r.handleClientEvent(DiceClicked{}, Player{ID: "p1"}) // trigger a move

	fmt.Println("Phase :", r.game.phase)

	_, _, err = r.handleClientEvent(ClientMove{}, Player{ID: "p2"})
	if err == nil {
		t.Fatal("expected error for invalid move")
	}
	_, _, err = r.handleClientEvent(ClientMove{Tile: 89}, Player{ID: "p1"})
	if err == nil {
		t.Fatal("expected error for invalid tile")
	}
	_, _, err = r.handleClientEvent(Answer{}, Player{ID: "p1"})
	if err == nil {
		t.Fatal("expected error for inconsitent phase")
	}

	expected := r.game.dice.Face // since we start at tile 0
	_, _, err = r.handleClientEvent(ClientMove{Tile: int(expected)}, Player{ID: "p1"})
	if err != nil {
		t.Fatal(err)
	}
	if r.game.pawnTile != int(expected) {
		t.Fatalf("expected %d, got %d", expected, r.game.pawnTile)
	}
	if r.game.dice != (DiceThrow{}) { // dice is reset on move
		t.Fatal(r.game.dice)
	}

	fmt.Println("Phase :", r.game.phase)

	_, _, err = r.handleClientEvent(WantNextTurn{}, Player{ID: "p1"})
	if err == nil {
		t.Fatal("expected error for inconsitent phase")
	}
	_, _, err = r.handleClientEvent(Answer{client.QuestionAnswersIn{}}, Player{ID: "p1"})
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = r.handleClientEvent(Answer{client.QuestionAnswersIn{}}, Player{ID: "p2"})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Phase :", r.game.phase)

	_, _, err = r.handleClientEvent(DiceClicked{}, Player{ID: "p1"})
	if err == nil {
		t.Fatal("expected error for inconsitent phase")
	}

	// check if nextTurn is properly reset
	r.game.currentWantNextTurn = map[serial]bool{"p3": true}
	_, _, err = r.handleClientEvent(WantNextTurn{}, Player{ID: "p1"})
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = r.handleClientEvent(WantNextTurn{}, Player{ID: "p2"})
	if err != nil {
		t.Fatal(err)
	}
	if len(r.game.currentWantNextTurn) != 0 {
		t.Fatal("currentWantNextTurn should be reset")
	}

	fmt.Println("Phase :", r.game.phase)

	r.game.pawnTile = 0
	_, _, err = r.handleClientEvent(DiceClicked{}, Player{ID: "p2"})
	if err != nil {
		t.Fatal("expected error for invalid click")
	}
	expected = r.game.dice.Face // since we reset at tile 0
	_, _, err = r.handleClientEvent(ClientMove{Tile: int(expected)}, Player{ID: "p2"})
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = r.handleClientEvent(Answer{client.QuestionAnswersIn{}}, Player{ID: "p1"})
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = r.handleClientEvent(Answer{client.QuestionAnswersIn{}}, Player{ID: "p2"})
	if err != nil {
		t.Fatal(err)
	}

	// check the end game
	r.players["p1"].advance.success = Success{true, true, true, true, true}

	if !reflect.DeepEqual(r.winners(), []serial{"p1"}) {
		t.Fatal()
	}

	if r.game.phase == pGameOver {
		t.Fatal("game should still be active")
	}

	_, _, err = r.handleClientEvent(WantNextTurn{}, Player{ID: "p1"})
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = r.handleClientEvent(WantNextTurn{}, Player{ID: "p2"})
	if err != nil {
		t.Fatal(err)
	}

	if r.game.phase != pGameOver {
		t.Fatal("game should be over")
	}

	fmt.Println("Phase :", r.game.phase)
}

func TestGameEnd(t *testing.T) {
	r := NewRoom("<test>", Options{Launch: LaunchStrategy{Max: 3}, Questions: exPool, QuestionTimeout: 10 * time.Millisecond, ShowDecrassage: true}, noOpSuccesHandler{})

	var c1 clientOut

	r.mustJoinConn(t, "p1", &c1)
	r.mustJoin(t, "p2")
	r.mustJoin(t, "p3")

	r.players["p1"].advance = playerAdvance{
		success: Success{true, true, true, true, true},
		review: QuestionReview{
			MarkedQuestions: []editor.IdQuestion{1, 1, 1, 2, 3, 4},
		},
	}

	r.players["p2"].advance.success = Success{true, true, true, true, true}

	rep := make(chan Replay)
	go func() {
		replay, _ := r.Listen(context.Background())
		rep <- replay
	}()

	r.Leave <- "p3"

	r.throwAndMove("p1")

	questionID := r.lg().question.ID

	r.Event <- ClientEvent{Event: Answer{client.QuestionAnswersIn{}}, Player: "p1"} // correct
	// p2 is incorrect due to timeout

	time.Sleep(50 * time.Millisecond)

	r.Event <- ClientEvent{Event: WantNextTurn{}, Player: "p1"}
	r.Event <- ClientEvent{Event: WantNextTurn{true}, Player: "p2"}

	events := c1.lastU(&r.lock).Events
	if len(events) != 1 {
		t.Fatal(events)
	}
	gameEnd, ok := events[0].(GameEnd)
	if !ok {
		t.Fatal(events[0])
	}

	if !reflect.DeepEqual(gameEnd.Winners, []serial{"p1"}) {
		t.Fatal(gameEnd.Winners)
	}

	decrP1 := editor.NewIdQuestionSetFrom(gameEnd.QuestionDecrassageIds["p1"])
	decrP2 := editor.NewIdQuestionSetFrom(gameEnd.QuestionDecrassageIds["p2"])

	if !decrP2.Has(questionID) {
		t.Fatal(decrP2)
	}
	if len(decrP1) != 3 {
		t.Fatal(decrP1)
	}

	replay := <-rep
	if len(replay.QuestionHistory) != 3 {
		t.Fatal()
	}
}

func TestGameTimeout(t *testing.T) {
	r := NewRoom("<test>", Options{Launch: LaunchStrategy{Max: 2}, Questions: exPool, QuestionTimeout: 10 * time.Millisecond, ShowDecrassage: true}, noOpSuccesHandler{})

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancelFunc()

	done := make(chan bool)
	go func() {
		_, _ = r.Listen(ctx)
		done <- true
	}()

	r.mustJoin(t, "p1")
	r.mustJoin(t, "p2")

	// wait for timeout
	isDone := <-done
	tu.Assert(t, isDone)
}

func TestSampleQuestion(t *testing.T) {
	wq := WeigthedQuestions{
		Questions: []editor.Question{
			{Id: 0},
			{Id: 1},
			{Id: 2},
			{Id: 3},
			{Id: 4},
		},
		Weights: []float64{
			0.2,
			0.2,
			0.2,
			0.2,
			0.2,
		},
	}

	for range [1000]int{} {
		qu := wq.sample(questionHistory{0: 1})
		tu.Assert(t, qu.Id != 0)
	}
	for range [1000]int{} {
		qu := wq.sample(questionHistory{0: 1, 1: 1})
		tu.Assert(t, qu.Id != 0 && qu.Id != 1)
	}
	for range [1000]int{} {
		_ = wq.sample(questionHistory{0: 1, 1: 1, 2: 1, 3: 1, 4: 1})
	}
	for range [1000]int{} {
		qu := wq.sample(questionHistory{0: 1, 1: 1, 2: 1, 3: 1, 4: 2})
		tu.Assert(t, qu.Id != 4)
	}
}

func TestPseudos(t *testing.T) {
	r := Room{players: map[PlayerID]*playerConn{
		"1": {pl: Player{"1", "Ben", "Kugler", 0}},
		"2": {pl: Player{"2", "Ben", "Symp", 0}},
		"3": {pl: Player{"3", "Ben", "", 0}},
		"4": {pl: Player{"4", "George", "Kugler", 0}},
	}}
	tu.Assert(t, r.serialToPseudo("1") == "Ben K.")
	tu.Assert(t, r.serialToPseudo("2") == "Ben S.")
	tu.Assert(t, r.serialToPseudo("3") == "Ben 3.")
	tu.Assert(t, r.serialToPseudo("4") == "George")
}
