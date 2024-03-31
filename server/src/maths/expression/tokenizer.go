package expression

import (
	"bytes"
	"strconv"
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

func (symbol) isToken()              {}
func (nT) isToken()                  {}
func (constant) isToken()            {}
func (Variable) isToken()            {}
func (function) isToken()            {}
func (roundFunc) isToken()           {}
func (specialFunctionKind) isToken() {}
func (operator) isToken()            {}

// differs on regular function by the parsing
// of its arguments
type specialFunctionKind uint8

const (
	randInt specialFunctionKind = iota
	randPrime
	randChoice
	randMatrixInt // random matrix with integer coeff
	choiceFrom    // the last argument is used as selector
	randDenominator
	minFn
	maxFn
	sumFn  // sum of an expression over a variable
	prodFn // product of an expression over a variable
	matCoeff
	matSet   // update a matrix
	binomial // coefficient binomial (n, k)

	invalidSpecialFunction
)

func (sf specialFunctionKind) String() string {
	switch sf {
	case randInt:
		return "randInt"
	case randPrime:
		return "randPrime"
	case randChoice:
		return "randChoice"
	case randMatrixInt:
		return "randMatrix"
	case choiceFrom:
		return "choiceFrom"
	case randDenominator:
		return "randDecDen"
	case minFn:
		return "min"
	case maxFn:
		return "max"
	case sumFn:
		return "sum"
	case prodFn:
		return "prod"
	case matCoeff:
		return "coeff"
	case matSet:
		return "set"
	case binomial:
		return "binom"
	default:
		panic(exhaustiveSpecialFunctionSwitch)
	}
}

type symbol uint8

const (
	openPar      symbol = iota // (
	closePar                   // )
	semicolon                  // ;
	openBracket                // [
	closeBracket               // ]
	openCurly                  // {
	closeCurly                 // }
	underscore                 // _
	openMatrix                 // [[

	invalidSymbol
)

func (sy symbol) String() string {
	switch sy {
	case openPar:
		return "("
	case closePar:
		return ")"
	case semicolon:
		return ";"
	case openBracket:
		return "["
	case closeBracket:
		return "]"
	case openCurly:
		return "{"
	case closeCurly:
		return "}"
	case underscore:
		return "_"
	case openMatrix:
		return "[["
	default:
		panic(exhaustiveSymbolSwitch)
	}
}

// number token
type nT string

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
	current, next := tk.readToken()
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
	current, next := tk.readToken()
	tk.currentToken = next // only save the potential next token, not current
	return current
}

// readToken wraps readRawToken and handles implicit
// multiplication, and special exponents
//
// When an implicit multiplication is detected, a corresponding
// token is returned and `lastToken` and `currentToken` are updated accordingly
//
// Similarly, when an exponent is detected, a ^ token is returned
// and `lastToken` and `currentToken` are updated accordingly
func (tk *tokenizer) readToken() (current, next token) {
	nextToken := tk.readRawToken()

	if va, isVar := nextToken.data.(Variable); isVar && va.Name == runeIsExponent { // insert a ^
		current = token{data: pow, pos: nextToken.pos}
		next = token{data: nT(va.Indice), pos: nextToken.pos}
	} else if isImplicitMult(tk.lastToken, nextToken) { // insert a mult
		current = token{data: mult, pos: nextToken.pos}
		next = nextToken
	} else {
		current = nextToken
		next = token{}
	}
	tk.lastToken = current // update lastToken for next calls

	return
}

// returns false if `r` should stop variable indices read
func isVariableChar(r rune) bool {
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

// check for 1 or 2 runes operators, returning the number
// of them, with zero meaning it is not an operator
func isOperator(src []rune) (op operator, n int) {
	// starts with two runes ops
	if len(src) >= 2 {
		s := string(src[0:2])
		switch s {
		case "==":
			return equals, 2
		case ">=":
			return greater, 2
		case "<=":
			return lesser, 2
		case "//":
			return rem, 2
		}
	}
	switch src[0] {
	case '>':
		return strictlyGreater, 1
	case '<':
		return strictlyLesser, 1
	case '+', '\ufe62':
		return plus, 1
	case '-', '\u2212':
		return minus, 1
	case '/', '\u00F7':
		return div, 1
	case '*', '\u00D7':
		return mult, 1
	case '%':
		return mod, 1
	case '^':
		return pow, 1
	case '!':
		return factorial, 1
	case '\u222A':
		return union, 1
	case '\u2229':
		return intersection, 1
	case '\u00AC':
		return complement, 1
	}
	_ = exhaustiveOperatorSwitch
	return op, 0
}

// match the special runes x¹, x², etc..
func isExponent(c rune) int {
	switch c {
	case 0x00b9:
		return 1
	case 0x00b2:
		return 2
	case 0x00b3:
		return 3
	case 0x2074:
		return 4
	case 0x2075:
		return 5
	case 0x2076:
		return 6
	case 0x2077:
		return 7
	case 0x2078:
		return 8
	case 0x2079:
		return 9
	default:
		return -1
	}
}

// advance pos to the next non whitespace char
func (tk *tokenizer) skipWhitespace() {
	for tk.pos < len(tk.src) && isWhiteSpace(tk.src[tk.pos]) {
		tk.pos++
	}
}

const runeIsExponent rune = -2

// advance the position
func (tk *tokenizer) readRawToken() (out token) {
	tk.skipWhitespace()

	if tk.pos >= len(tk.src) {
		return token{pos: tk.pos} // nil data for EOF
	}

	out = token{pos: tk.pos}
	c := tk.src[tk.pos]

	if exp := isExponent(c); exp >= 0 {
		tk.pos++
		out.data = Variable{Name: runeIsExponent, Indice: strconv.Itoa(exp)}
		return
	}

	op, isOpRunes := isOperator(tk.src[tk.pos:])
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
	case c == '[':
		// look ahead for [[
		tk.pos++
		tk.skipWhitespace()
		if tk.pos < len(tk.src) && tk.src[tk.pos] == '[' {
			out.data = openMatrix
			tk.pos++
		} else {
			out.data = openBracket
		}
	case c == ']':
		out.data = closeBracket
		tk.pos++
	case c == '{':
		out.data = openCurly
		tk.pos++
	case c == '}':
		out.data = closeCurly
		tk.pos++
	case c == '_':
		out.data = underscore
		tk.pos++
	case isOpRunes != 0:
		out.data = op
		tk.pos += isOpRunes

	case c == '"': // custom symbol
		out.data = tk.readCustomSymbol()
	case unicode.IsLetter(c): // either a function, a variable, Inf/inf or a constant
		if isInf := tk.tryReadInf(); isInf {
			out.data = nT("inf")
		} else if isRound := tk.tryReadRoundFunction(); isRound {
			out.data = roundFunc{}
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
		currentPos := tk.pos
		out.data = tk.readNumber()
		if tk.pos == currentPos { // invalid input, avoid infinite loop
			tk.pos++
		}
	}

	return out
}

// isImplicitMultRight returns true if `t` may be on
// the right of an implicit multiplication
func isImplicitMultRight(t token) bool {
	switch t := t.data.(type) {
	case symbol:
		return t == openPar || t == openMatrix // (...)(...) or ... [[ 2 ]]
	case Variable, constant: // (...)y or ...(pi)
		return true
	case function: // (...)log()
		return true
	case specialFunctionKind: // (...)randPrime()
		return true
	default:
		return false
	}
}

// isImplicitMultLeft returns true if `t` may be on
// the left of an implicit multiplication
func isImplicitMultLeft(t token) bool {
	switch t := t.data.(type) {
	case nT, constant: // 4x
		return true
	case operator:
		return t == factorial // n! n
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

func (tk *tokenizer) tryReadSpecialFunction() (specialFunctionKind, bool) {
	letters := tk.peekLetters()
	var fn specialFunctionKind
	switch strings.ToLower(string(letters)) {
	case "randint":
		fn = randInt
	case "randprime":
		fn = randPrime
	case "randmatrix":
		fn = randMatrixInt
	case "randchoice":
		fn = randChoice
	case "choicefrom":
		fn = choiceFrom
	case "randdecden":
		fn = randDenominator
	case "min":
		fn = minFn
	case "max":
		fn = maxFn
	case "sum":
		fn = sumFn
	case "prod":
		fn = prodFn
	case "coeff":
		fn = matCoeff
	case "set":
		fn = matSet
	case "binom":
		fn = binomial
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
	case "floor":
		fn = floorFn
	case "sqrt":
		fn = sqrtFn
	case "sgn":
		fn = sgnFn
	case "isPrime", "isprime":
		fn = isPrimeFn
	case "forceDecimal", "forcedecimal":
		fn = forceDecimalFn
	case "trace":
		fn = traceFn
	case "trans", "transpose":
		fn = transposeFn
	case "det":
		fn = detFn
	case "inv", "inverse":
		fn = invertFn
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
	case 'p': // also accepts pi for π
		if tk.pos+1 < len(tk.src) && tk.src[tk.pos+1] == 'i' {
			tk.pos += 2
			return piConstant, true
		}
		return 0, false
	default:
		_ = exhaustiveConstantSwitch
		return 0, false
	}
}

// assume we are at the starting "
func (tk *tokenizer) readCustomSymbol() Variable {
	tk.pos++ // consume the starting "
	start := tk.pos
	for ; tk.pos < len(tk.src); tk.pos++ {
		if tk.src[tk.pos] == '"' {
			break
		}
	}
	content := string(tk.src[start:tk.pos])
	tk.pos++ // consume the closing "

	return Variable{Indice: content}
}

func (tk *tokenizer) readVariable() Variable {
	c := tk.src[tk.pos]
	out := Variable{Name: c}
	tk.pos++

	// do not read _{...} as variable indice
	if tk.pos+1 < len(tk.src) && tk.src[tk.pos] == '_' && tk.src[tk.pos+1] == '{' {
		return out
	}

	if tk.pos < len(tk.src) && tk.src[tk.pos] == '_' { // indice start
		tk.pos++
		start := tk.pos
		for ; tk.pos < len(tk.src); tk.pos++ {
			if !isVariableChar(tk.src[tk.pos]) {
				break
			}
		}
		out.Indice = string(tk.src[start:tk.pos])
	}
	return out
}

func (tk *tokenizer) readNumber() nT {
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

	return nT(buffer)
}
