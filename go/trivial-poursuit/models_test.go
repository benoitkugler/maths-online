package trivialpoursuit

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestEventsJSON(t *testing.T) {
	dice := newDiceThrow()
	moves := Board.choices(0, int(dice.Face))
	payload := EventRange{
		Events: []event{
			playerTurn{2},
			dice,
			possibleMoves{moves},
			move{moves[0]},
			showQuestion{Question: "Super", Categorie: 0},
			playerAnswerResult{0, true},
			playerAnswerResult{1, false},
			playerAnswerResult{2, true},
			playerTurn{0},
			diceThrow{3},
			move{4},
			showQuestion{Question: "Super", Categorie: 1},
			playerAnswerResult{0, false},
			playerAnswerResult{1, true},
			playerAnswerResult{2, true},
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
