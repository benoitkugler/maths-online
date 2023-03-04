package repere

import (
	"fmt"
	"math"
)

func (c Coord) Round() IntCoord {
	return IntCoord{
		X: int(math.Round(c.X)),
		Y: int(math.Round(c.Y)),
	}
}

type IntCoord struct {
	X, Y int
}

// ColorHex is an hex or ahex color string, with
// #FFFFFF or #AAFFFFFF format
type ColorHex string

func (c ColorHex) ToARGB() (a, r, g, b uint8) {
	if c == "" {
		return
	}
	c = c[1:]
	if len(c) == 6 {
		c = "FF" + c
	} else if len(c) != 8 {
		return
	}
	_, _ = fmt.Sscanf(string(c), "%2x%2x%2x%2x", &a, &r, &g, &b)
	return
}
