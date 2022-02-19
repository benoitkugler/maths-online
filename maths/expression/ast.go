package expression

type node struct {
	left, right *node
	atom        atom
}

// atom is either an operator, a function,
// a variable, a predefined constant or a numerical value
type atom interface{}

type operator uint8

const (
	plus operator = iota
	minus
	mult
	div
	pow // x^2
)

type function uint8

const (
	log function = iota
	exp
	sin
	cos
	abs
	// round
)

// such as a, b in  (a + b)^2
// of x in 2x + 3
type variable rune

func newVariable(r rune) *node {
	return &node{atom: variable(r)}
}

type constant uint8

const (
	numberPi constant = iota
	numberE
	// i
)

type number float64

func newNumber(v float64) *node {
	return &node{atom: number(v)}
}
