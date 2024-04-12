package trivial

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	"github.com/benoitkugler/maths-online/server/src/sql/editor"
)

func TestEventsJSON(t *testing.T) {
	dice := newDiceThrow()
	moves := Board.choices(0, int(dice.Face)).list()
	question := client.Question{Enonce: client.Enonce{client.NumberFieldBlock{}}, Correction: client.Enonce{}}
	payload := StateUpdate{
		Events: []ServerEvent{
			PlayerJoin{},
			PlayerReconnected{},
			LobbyUpdate{PlayerPseudos: map[serial]string{"0": "Paul"}},
			GameStart{},
			PlayerLeft{"1"},
			PlayerTurn{"Haha", "2"},
			DiceThrow{3},
			Move{Tile: moves[0], Path: []int{0}},
			PossibleMoves{"", moves, "2"},
			ShowQuestion{ID: 1, Categorie: 0, Question: question},
			PlayerAnswerResults{
				Results: map[serial]playerAnswerResult{
					"0": {Success: true},
					"1": {Success: false},
					"2": {Success: true},
				},
			},
			GameEnd{
				QuestionDecrassageIds: map[serial][]editor.IdQuestion{"0": {1}},
				Winners:               []serial{"2"},
				WinnerNames:           []string{"Paul"},
			},
			GameTerminated{},
		},
	}

	_ = payload.String()

	b, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))

	var payload2 StateUpdate
	err = json.Unmarshal(b, &payload2)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(payload, payload2) {
		t.Fatalf("expected %#v, got %#v", payload, payload2)
	}
}

func TestGameStateJSON(t *testing.T) {
	payload := GameState{
		Players: map[serial]PlayerStatus{
			"0": {Success: Success{true, false}},
			"1": {Success: Success{false, true, true}},
		},
		PawnTile: 2,
	}
	b, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
}

func TestClientEventJSON(t *testing.T) {
	for _, event := range []ClientEventITF{
		ClientMove{Tile: 4, Path: []int{}},
		Answer{client.QuestionAnswersIn{Data: make(map[int]client.Answer)}},
		DiceClicked{},
		WantNextTurn{true},
		Ping{"Test"},
	} {
		payload := ClientEventITFWrapper{event}
		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(string(b))

		var payload2 ClientEventITFWrapper
		err = json.Unmarshal(b, &payload2)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(payload, payload2) {
			t.Fatal()
		}
	}
}

func TestMethodTag(t *testing.T) {
	PlayerJoin{}.isServerEvent()
	PlayerReconnected{}.isServerEvent()
	LobbyUpdate{}.isServerEvent()
	GameStart{}.isServerEvent()
	PlayerLeft{}.isServerEvent()
	PlayerTurn{}.isServerEvent()
	DiceThrow{}.isServerEvent()
	Move{}.isServerEvent()
	PossibleMoves{}.isServerEvent()
	ShowQuestion{}.isServerEvent()
	PlayerAnswerResults{}.isServerEvent()
	GameEnd{}.isServerEvent()
	GameTerminated{}.isServerEvent()

	ClientMove{}.isClientEvent()
	Answer{}.isClientEvent()
	DiceClicked{}.isClientEvent()
	WantNextTurn{}.isClientEvent()
	Ping{}.isClientEvent()
}

func TestSuccess_isDone(t *testing.T) {
	tests := []struct {
		sc   Success
		want bool
	}{
		{Success{false, true, false, true, true}, false},
		{Success{true, true, false, true, true}, false},
		{Success{true, true, true, true, true}, true},
	}
	for _, tt := range tests {
		if got := tt.sc.isDone(); got != tt.want {
			t.Errorf("Success.isDone() = %v, want %v", got, tt.want)
		}
	}
}

func TestQuestionReview_streak(t *testing.T) {
	tests := []struct {
		QuestionHistory []QR
		want            int
		want2           bool
	}{
		{
			[]QR{},
			0, false,
		},
		{
			[]QR{{Success: false}},
			0, false,
		},
		{
			[]QR{{Success: false}, {Success: true}},
			1, false,
		},
		{
			[]QR{{Success: false}, {Success: true}, {Success: false}, {Success: true}, {Success: true}},
			2, false,
		},
		{
			[]QR{{Success: true}, {Success: true}, {Success: true}},
			3, true,
		},
		{
			[]QR{{Success: true}, {Success: true}, {Success: true}, {Success: true}},
			4, false,
		},
		{
			[]QR{{Success: true}, {Success: true}, {Success: true}, {Success: true}, {Success: true}, {Success: true}},
			6, true,
		},
	}
	for _, tt := range tests {
		qr := QuestionReview{
			QuestionHistory: tt.QuestionHistory,
		}
		if got := qr.streak(); got != tt.want {
			t.Errorf("QuestionReview.streak() = %v, want %v", got, tt.want)
		}
		if got := qr.hasStreak3(); got != tt.want2 {
			t.Errorf("QuestionReview.hasStreak3() = %v, want %v", got, tt.want)
		}
	}
}
