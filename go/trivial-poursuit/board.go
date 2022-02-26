// Package trivialpoursuit implements a backend for
// a multi player trivial poursuit game, where questions
// are (short) maths questions.
package trivialpoursuit

import "sort"

// To support arbitrary trivial poursuit board,
// we model it by a graph

const nbSquares = 19

// Board is the Trivial-Poursuit board game shape.
var Board = board{ // the first tile is in the center
	{1: true, nbSquares - 1: true},
	{0: true, 2: true},
	{1: true, 3: true},
	{2: true, 4: true},
	4: {3: true, 5: true, nbSquares - 4: true},
	{4: true, 6: true},
	{5: true, 7: true},
	{6: true, 8: true},
	{7: true, 9: true},
	{8: true, 10: true},
	10: {11: true, 9: true, nbSquares - 3: true},
	{10: true, 12: true},
	{11: true, 13: true},
	{12: true, 14: true},
	{13: true, 15: true},
	15: {14: true, 4: true},
	16: {17: true, 10: true},
	{16: true, 18: true},
	{17: true, 0: true},
}

// board is a trivial-poursuit board, stored as an adjency matrix
// even if board[i][i] is true, it is ignored.
type board [nbSquares][nbSquares]bool

// adjacents returns the squares accessible from `pos`
// in one move.
func (b board) adjacents(pos int) (out []int) {
	for i, b := range b[pos] {
		if b && i != pos { // do not add
			out = append(out, i)
		}
	}
	return out
}

// choices returns the square indices where the player located at `currentPos`
// may advances with `nbMoves`.
// In this game, you can go back when you have made one step.
func (b board) choices(currentPos, nbMoves int) []int {
	type target struct {
		pos    int
		origin int // -1 for no constraint
	}
	// start with the current pos, with no constraint
	var (
		targets = []target{{pos: currentPos}}
		buffer  []target
	)
	for i := 0; i < nbMoves; i++ { // one step at a time
		for _, t := range targets {
			// advance everywhere expect to the origin
			candidates := b.adjacents(t.pos)
			// remove the previous square
			for _, candidate := range candidates {
				if candidate == t.origin {
					continue
				}
				buffer = append(buffer, target{pos: candidate, origin: t.pos})
			}
		}

		targets = buffer
		buffer = nil
	}

	// in case a square is reachable from two paths, remove duplicates
	uniqueSquares := make(map[int]bool)
	for _, t := range targets {
		uniqueSquares[t.pos] = true
	}

	// convert back to a sorted slice (for reproducibility)
	var out []int
	for pos := range uniqueSquares {
		out = append(out, pos)
	}
	sort.Ints(out)

	return out
}
