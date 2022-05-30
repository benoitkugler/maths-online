package expression

import (
	"bytes"
	"strings"
	"unicode"
)

type token struct {
	data tokenData // a nil data field means EOF
	pos  int
}

type tokenData interface {
	isToken()
}

func (symbol) isToken()          {}
func (numberText) isToken()      {}
func (constant) isToken()        {}
func (Variable) isToken()        {}
func (function) isToken()        {}
func (roundFn) isToken()         {}
func (specialFunction) isToken() {}
func (randVariable) isToken()    {}
func (operator) isToken()        {}

// differs on regular function by the parsing
// of its arguments
type specialFunction uint8

const (
	randInt specialFunction = iota
	randPrime
	randChoice
	randDenominator

	invalidSpecialFunction
)

type symbol uint8

const (
	openPar   symbol = iota // (
	closePar                // )
	semicolon               // ;

	invalidSymbol
)

type numberText string

type tokenizer struct {
	// caches used for peek and to handle implicit multiplication
	// with implicit multiplication, and additional token may be added
	lastToken, currentToken, nextToken token

	src []rune
	pos int
}

func newTokenizer(text []byte) *tokenizer { return &tokenizer{src: bytes.Runes(text)} }

// Peek returns the next token, without advancing the position.
func (tk *tokenizer) Peek() token {
	// check if we have a cached value in currentToken or nextToken
	if tk.currentToken.data != nil { // return it without changing the state
		return tk.currentToken
	}

	// trigger a new read
	current, next := tk.readTokenImplicitMult()
	// save both tokens
	tk.currentToken = current
	tk.nextToken = next

	return tk.currentToken
}

func (tk *tokenizer) Next() token {
	// check if we have already read the next token
	if out := tk.currentToken; out.data != nil { // returns it without triggering a read
		tk.lastToken = out
		tk.currentToken = tk.nextToken
		tk.nextToken = token{}
		return out
	}

	// trigger a new read
	current, next := tk.readTokenImplicitMult()
	tk.currentToken = next // only save the potential next token, not current
	return current
}

// readTokenImplicitMult wraps readToken and handles implicit
// multiplication.
// When an implicit multiplication is detected, a corresponding
// token is returned and `currentToken` and `currentToken` are updated accordingly
func (tk *tokenizer) readTokenImplicitMult() (current, next token) {
	nextToken := tk.readToken()
	if isImplicitMult(tk.lastToken, nextToken) { // insert a mult
		current = token{data: mult, pos: nextToken.pos}
		next = nextToken
	} else {
		current = nextToken
		next = token{}
	}
	tk.lastToken = current // update lastToken for next calls

	return
}

func isIdentifier(r rune) bool {
	return 'A' <= r && r <= 'z' || '0' <= r && r <= '9'
}

func isWhiteSpace(r rune) bool {
	switch r {
	case ' ', '\n', '\t', '\r', '\f':
		return true
	default:
		return false
	}
}

func isRemainder(src []rune) bool {
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

	isRem := isRemainder(tk.src[tk.pos:])
	op, isOp := isOperator(c)
	switch {
	case c == '(':
		out.data = openPar
		tk.pos++
	case c == ')':
		out.data = closePar
		tk.pos++
	case c == ';':
		out.data = semicolon
		tk.pos++
	case isRem:
		out.data = rem
		tk.pos += 2
	case isOp:
		out.data = op
		tk.pos++
	case unicode.IsLetter(c): // either a function, a variable, Inf/inf a constant
		if tk.tryReadRandVariable() {
			out.data = randVariable{}
		} else if isInf := tk.tryReadInf(); isInf {
			out.data = numberText("inf")
		} else if isRound := tk.tryReadRoundFunction(); isRound {
			out.data = roundFn{}
		} else if fn, isSpecial := tk.tryReadSpecialFunction(); isSpecial {
			out.data = fn
		} else if fn, isFunction := tk.tryReadFunction(); isFunction {
			out.data = fn
		} else if ct, isConst := tk.tryReadConstant(); isConst {
			out.data = ct
		} else {
			// read the symbol as variable
			v := tk.readVariable()
			out.data = v
		}
	case unicode.Is(unicode.Co, c): // custom variables
		v := tk.readVariable()
		out.data = v
	default: // number start
		out.data = tk.readNumber()
	}

	return out
}

// isImplicitMultRight returns true if `t` may be on
// the right of an implicit multiplication
func isImplicitMultRight(t token) bool {
	switch t := t.data.(type) {
	case symbol:
		return t == openPar // (...)(...)
	case Variable, constant: // (...)y or ...(pi)
		return true
	case function: // (...)log()
		return true
	case specialFunction: // (...)randPrime()
		return true
	default:
		return false
	}
}

// isImplicitMultLeft returns true if `t` may be on
// the left of an implicit multiplication
func isImplicitMultLeft(t token) bool {
	switch t := t.data.(type) {
	case numberText, constant: // 4x
		return true
	case symbol:
		return t == closePar // (...)(...)
	case Variable: // y(...)
		return true
	default:
		return false
	}
}

// isImplicitMult returns true if there is an implicit multiplication
// between t1 and t2
func isImplicitMult(t1, t2 token) bool {
	return isImplicitMultLeft(t1) && isImplicitMultRight(t2)
}

func (tk *tokenizer) tryReadRoundFunction() bool {
	if s := string(tk.peekLetters()); s == "round" {
		tk.pos += len("round")
		return true
	}
	return false
}

func (tk *tokenizer) tryReadInf() bool {
	if s := strings.ToLower(string(tk.peekLetters())); s == "inf" {
		tk.pos += len("inf")
		return true
	}
	return false
}

func (tk *tokenizer) tryReadRandVariable() bool {
	if s := string(tk.peekLetters()); s == "randLetter" || s == "randletter" {
		tk.pos += len("randLetter")
		return true
	}
	return false
}

func (tk *tokenizer) tryReadSpecialFunction() (specialFunction, bool) {
	letters := tk.peekLetters()
	var fn specialFunction
	switch string(letters) {
	case "randint", "randInt":
		fn = randInt
	case "randPrime", "randprime":
		fn = randPrime
	case "randChoice", "randchoice":
		fn = randChoice
	case "randDecDen", "randdecden":
		fn = randDenominator
	default:
		_ = exhaustiveSpecialFunctionSwitch
		return 0, false
	}

	// found a function, advance the position
	tk.pos += len(letters)
	return fn, true
}

// return the next letters; without advancing
func (tk *tokenizer) peekLetters() []rune {
	L := len(tk.src)

	// read subsequent letters
	var letters []rune
	for i := tk.pos; i < L && unicode.IsLetter(tk.src[i]); i++ {
		letters = append(letters, tk.src[i])
	}
	return letters
}

func (tk *tokenizer) tryReadFunction() (function, bool) {
	letters := tk.peekLetters()

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
	case "tan":
		fn = tanFn
	case "asin", "arcsin":
		fn = asinFn
	case "acos", "arccos":
		fn = acosFn
	case "atan", "arctan":
		fn = atanFn
	case "abs":
		fn = absFn
	case "sqrt":
		fn = sqrtFn
	case "sgn":
		fn = sgnFn
	case "isZero":
		fn = isZeroFn
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

func (tk *tokenizer) tryReadConstant() (constant, bool) {
	switch c := tk.src[tk.pos]; c {
	case 'e':
		tk.pos++
		return eConstant, true
	case piRune:
		tk.pos++
		return piConstant, true
	default:
		_ = exhaustiveConstantSwitch
		return 0, false
	}
}

func (tk *tokenizer) readVariable() Variable {
	c := tk.src[tk.pos]
	out := Variable{Name: c}
	tk.pos++
	if tk.pos < len(tk.src) && tk.src[tk.pos] == '_' { // indice start
		tk.pos++
		start := tk.pos
		for ; tk.pos < len(tk.src); tk.pos++ {
			if !isIdentifier(tk.src[tk.pos]) {
				break
			}
		}
		out.Indice = string(tk.src[start:tk.pos])
	}
	return out
}

func (tk *tokenizer) readNumber() numberText {
	var buffer []rune
	L := len(tk.src)
	for ; tk.pos < L; tk.pos++ {
		r := tk.src[tk.pos]
		if unicode.IsDigit(r) || r == '.' {
			buffer = append(buffer, r)
		} else if r == ',' { // also accept comma a decimal separator
			buffer = append(buffer, '.')
		} else {
			break
		}
	}

	return numberText(buffer)
}
