package game

import (
	"reflect"
	"testing"
)

func Test_board_adjacents(t *testing.T) {
	tests := []struct {
		b       board
		args    int
		wantOut []int
	}{
		{
			board{{true, false, true}, {true, true, true}},
			0,
			[]int{2},
		},
		{
			board{{true, false, true, true, true}, {true, true, true}},
			0,
			[]int{2, 3, 4},
		},
		{
			board{{true, false, true, true, true}, {true, true, true, true}},
			1,
			[]int{0, 2, 3},
		},
		{
			Board,
			18,
			[]int{0, 17},
		},
	}
	for _, tt := range tests {
		if gotOut := tt.b.adjacents(tt.args); !reflect.DeepEqual(gotOut, tt.wantOut) {
			t.Errorf("board.adjacents() = %v, want %v", gotOut, tt.wantOut)
		}
	}
}

func Test_board_choices(t *testing.T) {
	type args struct {
		currentPos int
		nbMoves    int
	}
	tests := []struct {
		b    board
		args args
		want []int
	}{
		{
			board{{false, true, true}, {true, true, true}, {false, true, false}},
			args{currentPos: 0, nbMoves: 1},
			[]int{1, 2},
		},
		{
			board{{false, true, true}, {true, true, true}, {false, true, false}},
			args{currentPos: 0, nbMoves: 2},
			[]int{1, 2}, // 0 -> 1 -> 2 ; 0 -> 2 -> 1
		},
		{
			board{{false, true, true}, {true, true, true}, {false, true, false}},
			args{currentPos: 0, nbMoves: 3},
			[]int{0}, // 0 -> 1 -> 2 -> nothing; 0 -> 2 -> 1 -> 0
		},
	}
	for _, tt := range tests {
		if got := tt.b.choices(tt.args.currentPos, tt.args.nbMoves).list(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("board.choices() = %v, want %v", got, tt.want)
		}
	}
}

func TestBoardMoves(t *testing.T) {
	// test the real board
	tests := []struct {
		pos     int
		nbMoves int
		want    []int
	}{
		{18, 2, []int{1, 16}},
		{0, 1, []int{1, nbSquares - 1}},
		{4, 1, []int{3, 5, 15}},
		{2, 3, []int{5, 15, nbSquares - 1}},
	}
	for _, tt := range tests {
		if got := Board.choices(tt.pos, tt.nbMoves).list(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Board.choices() = %v, want %v", got, tt.want)
		}
	}
}

func TestBoardPaths(t *testing.T) {
	// test the real board
	tests := []struct {
		pos     int
		nbMoves int
		want    tileSet
	}{
		// {18, 2, tileSet{1: []int{18, 0, 1}, 16: []int{18, 17, 16}}},
		// {0, 1, tileSet{1: []int{0, 1}, nbSquares - 1: []int{0, nbSquares - 1}}},
		{2, 3, tileSet{5: []int{2, 3, 4, 5}, 15: []int{2, 3, 4, 15}, nbSquares - 1: []int{2, 1, 0, nbSquares - 1}}},
	}
	for _, tt := range tests {
		if got := Board.choices(tt.pos, tt.nbMoves); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Board.choices() = %v, want %v", got, tt.want)
		}
	}
}
