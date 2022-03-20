package repere

import "math"

func (c Coord) Round() IntCoord {
	return IntCoord{
		X: int(math.Round(c.X)),
		Y: int(math.Round(c.Y)),
	}
}

type IntCoord struct {
	X, Y int
}
