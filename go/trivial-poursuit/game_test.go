package trivialpoursuit

import (
	"reflect"
	"testing"
)

func TestGame_currentState(t *testing.T) {
	tests := []struct {
		events []event
		want   GameState
	}{
		{
			[]event{
				diceThrow{2},
			},
			GameState{
				Successes: make([]success, 3),
				Dice:      2,
				// dice is reset on move
			},
		},
		{
			[]event{
				playerTurn{2},
				newDiceThrow(),
				move{3},
			},
			GameState{
				Successes: make([]success, 3),
				PawnTile:  3,
				Player:    2,
				// dice is reset on move
			},
		},
		{
			[]event{
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
			GameState{
				Successes: []success{
					{true},
					{false, true},
					{true, true},
				},
				PawnTile: 4,
				Player:   1,
			},
		},
	}
	for _, tt := range tests {
		g := Game{
			events:       tt.events,
			initialState: newGameState(3),
		}
		if got := g.currentState(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Game.currentState() = %v, want %v", got, tt.want)
		}
	}
}

func Test_gameState_winners(t *testing.T) {
	tests := []struct {
		sc      []success
		wantOut []int
	}{
		{[]success{{true}, {true}, {true, true, true, true, true}}, []int{2}},
		{[]success{{true}, {true}, {true, true, true, true}}, nil},
		{[]success{{true, true, true, true, true}, {true}, {true, true, true, true, true}}, []int{0, 2}},
	}
	for _, tt := range tests {
		gs := &GameState{
			Successes: tt.sc,
		}
		if gotOut := gs.winners(); !reflect.DeepEqual(gotOut, tt.wantOut) {
			t.Errorf("gameState.winners() = %v, want %v", gotOut, tt.wantOut)
		}
	}
}
