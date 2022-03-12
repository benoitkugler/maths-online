package game

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestEventsJSON(t *testing.T) {
	dice := newDiceThrow()
	moves := Board.choices(0, int(dice.Face)).list()
	payload := StateUpdates{{
		Events: []GameEvent{
			gameStart{},
			playerTurn{"Haha", 2},
			dice,
			possibleMoves{"", moves, 2},
			move{Tile: moves[0]},
			playerLeft{1},
			showQuestion{Question: "Super", Categorie: 0},
			playerAnswerResult{0, true},
			playerAnswerResult{1, false},
			playerAnswerResult{2, true},
			playerTurn{"", 0},
			diceThrow{3},
			move{Tile: 4},
			showQuestion{Question: "Super", Categorie: 1},
			playerAnswerResult{0, false},
			playerAnswerResult{1, true},
			playerAnswerResult{2, true},
			playerTurn{"", 1},
		},
	}}

	b, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))

	err = json.Unmarshal(b, &[]StateUpdate{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGameStateJSON(t *testing.T) {
	payload := GameState{
		Players: map[int]*PlayerStatus{
			0: {Success: success{true, false}},
			1: {Success: success{false, true, true}},
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
		move{Tile: 4},
		answer{"ma réponse"},
		Ping{"Test"},
	} {
		paylod := ClientEvent{Event: event, Player: 0}
		b, err := json.Marshal(paylod)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(string(b))

		var paylod2 ClientEvent
		err = json.Unmarshal(b, &paylod2)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(paylod, paylod2) {
			t.Fatal()
		}
	}
}

func TestMethodTag(t *testing.T) {
	gameStart{}.isGameEvent()
	playerLeft{}.isGameEvent()
	playerTurn{}.isGameEvent()
	diceThrow{}.isGameEvent()
	move{}.isGameEvent()
	possibleMoves{}.isGameEvent()
	showQuestion{}.isGameEvent()
	playerAnswerResult{}.isGameEvent()
	gameEnd{}.isGameEvent()
	move{}.isClientEvent()
	answer{}.isClientEvent()
}
