package game

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/benoitkugler/maths-online/maths/questions/client"
)

func TestEventsJSON(t *testing.T) {
	dice := newDiceThrow()
	moves := Board.choices(0, int(dice.Face)).list()
	payload := StateUpdate{
		Events: []GameEvent{
			GameStart{},
			PlayerTurn{"Haha", 2},
			dice,
			PossibleMoves{"", moves, 2},
			Move{Tile: moves[0]},
			PlayerLeft{1},
			ShowQuestion{ID: 1, Categorie: 0},
			PlayerAnswerResults{
				Results: map[int]playerAnswerResult{
					0: {Success: true},
					1: {Success: false},
					2: {Success: true},
				},
			},
			PlayerTurn{"", 0},
			DiceThrow{3},
			Move{Tile: 4},
			ShowQuestion{ID: 2, Categorie: 1},
			PlayerAnswerResults{
				Results: map[int]playerAnswerResult{
					0: {Success: true},
					1: {Success: false},
					2: {Success: true},
				},
			},
			PlayerTurn{"", 1},
		},
	}

	b, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))

	err = json.Unmarshal(b, &StateUpdate{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGameStateJSON(t *testing.T) {
	payload := GameState{
		Players: map[int]*PlayerStatus{
			0: {Success: Success{true, false}},
			1: {Success: Success{false, true, true}},
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
	for _, event := range []clientEventData{
		ClientMove{Tile: 4, Path: []int{}},
		Answer{client.QuestionAnswersIn{Data: make(map[int]client.Answer)}},
		Ping{"Test"},
	} {
		payload := ClientEvent{Event: event, Player: 0}
		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(string(b))

		var payload2 ClientEvent
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
	GameStart{}.isGameEvent()
	PlayerLeft{}.isGameEvent()
	PlayerTurn{}.isGameEvent()
	DiceThrow{}.isGameEvent()
	Move{}.isGameEvent()
	PossibleMoves{}.isGameEvent()
	ShowQuestion{}.isGameEvent()
	PlayerAnswerResults{}.isGameEvent()
	GameEnd{}.isGameEvent()

	ClientMove{}.isClientEvent()
	Answer{}.isClientEvent()
}
