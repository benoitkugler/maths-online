package expression

import (
	"bytes"
	"unicode"
)

type token struct {
	data tokenData // a nil data field means EOF
	pos  int
}

type tokenData interface {
	isToken()
}

func (symbol) isToken()       {}
func (numberText) isToken()   {}
func (constant) isToken()     {}
func (Variable) isToken()     {}
func (function) isToken()     {}
func (randFunction) isToken() {}
func (operator) isToken()     {}

type randFunction struct {
	isPrime bool
}

type symbol uint8

const (
	openPar  symbol = iota // (
	closePar               // )
	comma                  // ,

	invalidSymbol
)

type numberText string

type tokenizer struct {
	nextToken token // cache used for peek

	src []rune
	pos int
}

func newTokenizer(text []byte) *tokenizer { return &tokenizer{src: bytes.Runes(text)} }

func (tk *tokenizer) peek() token {
	if tk.nextToken.data != nil {
		return tk.nextToken
	}

	tk.nextToken = tk.readToken() // save it so that next do not advance again
	return tk.nextToken
}

func (tk *tokenizer) next() token {
	if out := tk.nextToken; out.data != nil {
		tk.nextToken = token{}
		return out
	}

	out := tk.readToken()
	return out
}

func isWhiteSpace(r rune) bool {
	switch r {
	case ' ', '\n', '\t', '\r', '\f':
		return true
	default:
		return false
	}
}

func isRemaider(src []rune) bool {
	return len(src) >= 2 && (string(src[0:2]) == "//")
}

func isOperator(r rune) (operator, bool) {
	var op operator
	switch r {
	case '+', '\ufe62':
		op = plus
	case '-', '\u2212':
		op = minus
	case '/', '\u00F7':
		op = div
	case '*', '\u00D7':
		op = mult
	case '%':
		op = mod
	case '^':
		op = pow
	default:
		return 0, false
	}
	return op, true
}

// advance pos to the next non whitespace char
func (tk *tokenizer) skipWhitespace() {
	for tk.pos < len(tk.src) && isWhiteSpace(tk.src[tk.pos]) {
		tk.pos++
	}
}

// advance the position
func (tk *tokenizer) readToken() (tok token) {
	tk.skipWhitespace()

	if tk.pos >= len(tk.src) {
		return token{pos: tk.pos}
	}

	out := token{pos: tk.pos}
	c := tk.src[tk.pos]

	isRem := isRemaider(tk.src[tk.pos:])
	op, isOp := isOperator(c)
	switch {
	case c == '(':
		out.data = openPar
		tk.pos++
	case c == ')':
		out.data = closePar
		tk.pos++
	case c == ',':
		out.data = comma
		tk.pos++
	case isRem:
		out.data = rem
		tk.pos += 2
	case isOp:
		out.data = op
		tk.pos++
	case unicode.IsLetter(c): // either a function, a variable or a constant
		if tk.tryReadRandint() {
			out.data = randFunction{isPrime: false}
		} else if tk.tryReadRandPrime() {
			out.data = randFunction{isPrime: true}
		} else if fn, isFunction := tk.tryReadFunction(); isFunction {
			out.data = fn
		} else {
			// read the symbol as variable or predefined constants
			v := tk.readConstantOrVariable()
			out.data = v
		}
	case unicode.Is(unicode.Co, c): // custom variables
		v := tk.readConstantOrVariable()
		out.data = v
	default: // number start
		out.data = tk.readNumber()
	}

	return out
}

func (tk *tokenizer) tryReadRandint() bool {
	const runeLen = len("randint")

	if len(tk.src) < tk.pos+runeLen {
		return false
	}

	word := string(tk.src[tk.pos : tk.pos+runeLen])
	if word == "randInt" || word == "randint" {
		tk.pos += runeLen
		return true
	}

	return false
}

func (tk *tokenizer) tryReadRandPrime() bool {
	const runeLen = len("randPrime")

	if len(tk.src) < tk.pos+runeLen {
		return false
	}

	word := string(tk.src[tk.pos : tk.pos+runeLen])
	if word == "randPrime" || word == "randprime" {
		tk.pos += runeLen
		return true
	}

	return false
}

func (tk *tokenizer) tryReadFunction() (function, bool) {
	L := len(tk.src)

	// read subsequent letters
	var letters []rune
	for i := tk.pos; i < L && unicode.IsLetter(tk.src[i]); i++ {
		letters = append(letters, tk.src[i])
	}

	var fn function
	switch string(letters) {
	case "exp":
		fn = expFn
	case "ln", "log":
		fn = logFn
	case "sin":
		fn = sinFn
	case "cos":
		fn = cosFn
	case "abs":
		fn = absFn
	case "sqrt":
		fn = sqrtFn
	case "sgn":
		fn = sgnFn
	case "isPrime":
		fn = isPrimeFn
	default: // no  matching function name
		_ = exhaustiveFunctionSwitch
		return 0, false
	}

	// found a function, advance the position
	tk.pos += len(letters)
	return fn, true
}

const piRune = '\u03C0'

func (tk *tokenizer) readConstantOrVariable() (out tokenData) {
	switch c := tk.src[tk.pos]; c {
	case 'e':
		out = eConstant
	case piRune:
		out = piConstant
	default:
		out = Variable(c)
	}
	tk.pos++

	return out
}

func (tk *tokenizer) readNumber() numberText {
	var buffer []rune
	L := len(tk.src)
	for ; tk.pos < L; tk.pos++ {
		r := tk.src[tk.pos]
		if unicode.IsDigit(r) || r == '.' {
			buffer = append(buffer, r)
		} else {
			break
		}
	}

	return numberText(buffer)
}
