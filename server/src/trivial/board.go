// Package trivialpoursuit implements a backend for
// a multi player trivial poursuit game, where questions
// are (short) maths questions.
package trivial

import (
	"sort"
)

// To support arbitrary trivial poursuit board,
// we model it by a graph

const nbSquares = 17

var (
	// Board is the Trivial-Poursuit board game shape.
	Board = board{ // the first tile is in the center
		0:  {1: true, nbSquares - 1: true},
		1:  {0: true, 2: true},
		2:  {1: true, 3: true},
		3:  {2: true, 4: true, nbSquares - 3: true}, // cross
		4:  {3: true, 5: true},
		5:  {4: true, 6: true},
		6:  {5: true, 7: true},
		7:  {6: true, 8: true},
		8:  {7: true, 9: true},
		9:  {8: true, 10: true, nbSquares - 2: true}, // cross
		10: {9: true, 11: true},
		11: {10: true, 12: true},
		12: {11: true, 13: true},
		13: {12: true, 14: true},
		14: {13: true, 3: true},
		15: {16: true, 9: true},
		16: {15: true, 0: true},
	}

	categories = [nbSquares]Categorie{
		Purple,
		Blue,
		Green,
		Blue, // cross
		Green,
		Orange,
		Purple,
		Yellow,
		Orange,
		Yellow, // cross
		Orange,
		Yellow,
		Purple,
		Orange,
		Green,
		Blue,
		Yellow,
	}
)

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

// tileSet is a set of reachable tiles, associated to a path towards them
type tileSet map[int]tilePath

// list returns a sorted slice, for reproducibility
func (ts tileSet) list() []int {
	var out []int
	for pos := range ts {
		out = append(out, pos)
	}
	sort.Ints(out)
	return out
}

// a path of contiguous tiles towards the last element of the slice
type tilePath []int

func (tp tilePath) pos() int { return tp[len(tp)-1] }

func (tp tilePath) origin() int {
	if len(tp) == 1 {
		return -1
	}
	return tp[len(tp)-2]
}

// choices returns the square indices where the player located at `currentPos`
// may advances with `nbMoves`.
// In this game, you can't go back when you have made one step.
func (b board) choices(currentPos, nbMoves int) tileSet {
	// start with the current pos, with no constraint
	var (
		currentPaths = []tilePath{{currentPos}}
		buffer       []tilePath
	)
	for i := 0; i < nbMoves; i++ { // one step at a time
		for _, path := range currentPaths {
			// advance everywhere expect to the origin
			candidates := b.adjacents(path.pos())

			// remove the previous square
			for _, candidate := range candidates {
				if candidate == path.origin() {
					continue
				}
				// extend the path : we need to copy to avoid sharing the same underlying array
				newPath := append(append(tilePath(nil), path...), candidate)
				buffer = append(buffer, newPath)
			}
		}

		currentPaths = buffer
		buffer = nil
	}

	// in case a square is reachable from two paths, remove duplicates
	out := make(tileSet)
	for _, t := range currentPaths {
		out[t.pos()] = t
	}

	return out
}
