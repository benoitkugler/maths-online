package trivialpoursuit

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestEventsJSON(t *testing.T) {
	payload := EventRange{
		Events: []event{
			playerTurn{2},
			newDiceThrow(),
			move{3},
			showQuestion{Question: "Super", Categorie: 0},
			playerAnswerSuccess{0, true},
			playerAnswerSuccess{1, false},
			playerAnswerSuccess{2, true},
			playerTurn{0},
			diceThrow{3},
			move{4},
			showQuestion{Question: "Super", Categorie: 1},
			playerAnswerSuccess{0, false},
			playerAnswerSuccess{1, true},
			playerAnswerSuccess{2, true},
			playerTurn{1},
		},
		Start: 0,
	}

	b, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
}

func TestGameStateJSON(t *testing.T) {
	payload := GameState{
		Successes: []success{
			{true, false},
			{false, true, true},
		},
		PawnTile: 2,
	}
	b, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
}
