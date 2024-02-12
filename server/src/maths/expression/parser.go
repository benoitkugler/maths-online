package expression

import (
	"fmt"
	"sort"
	"strconv"
)

// ErrInvalidExpr is returned by `Parse` and `ParseCompound`.
type ErrInvalidExpr struct {
	Input  string
	Reason string // in french
	Pos    int    // index in the input rune slice
}

func (inv ErrInvalidExpr) Error() string {
	return fmt.Sprintf(`L'expression %s est invalide : %s`, inv.Input, inv.Reason)
}

// Portion returns the start of `expr` until the error.
func (inv ErrInvalidExpr) Portion() string {
	return string([]rune(inv.Input)[:inv.Pos])
}

// varMap register the position of each variable occurence in
// an expression rune slice.
type varMap map[int]Variable

// Positions returns the indices in the original expression rune slice
// of every occurence of the given variables
func (vm varMap) Positions(subset map[Variable]bool) []int {
	var out []int
	for index, v := range vm {
		if _, has := subset[v]; has {
			out = append(out, index)
		}
	}

	sort.Ints(out)

	return out
}

// Parse parses a mathematical expression. If invalid, an `ErrInvalidExpr` is returned.
func Parse(s string) (*Expr, error) {
	expr, _, err := parseBytes([]byte(s))
	if err != nil {
		errV := err.(ErrInvalidExpr)
		errV.Input = s
		return nil, errV
	}
	return expr, nil
}

// MustParse is the same as Parse but panics on invalid expressions.
func MustParse(s string) *Expr {
	expr, err := Parse(s)
	if err != nil {
		panic(fmt.Sprintf("%s: %s", s, err))
	}
	return expr
}

// parseBytes parses a mathematical expression. If invalid, an `InvalidExpr` is returned.
func parseBytes(text []byte) (*Expr, varMap, error) {
	pr := newParser(text)
	e, err := pr.parseExpression(false)
	if err != nil {
		return nil, nil, err
	}

	// check if we properly reached EOF
	if tk := pr.tk.Peek().data; tk != nil {
		return nil, nil, ErrInvalidExpr{
			Reason: fmt.Sprintf("expression mal formée (symbol %s restant)", tk),
			Pos:    0,
		}
	}

	return e, pr.variablePos, err
}

// parser transforms an utf8 byte input into an ast
type parser struct {
	tk *tokenizer

	variablePos varMap // index of variable in input rune slice

	stack []*Expr // waiting to be consumed by operators
}

func newParser(text []byte) *parser {
	return &parser{tk: newTokenizer(text), variablePos: make(map[int]Variable)}
}

// return nil if the stack is empty
func (pr *parser) pop() *Expr {
	if len(pr.stack) == 0 {
		return nil
	}
	out := pr.stack[len(pr.stack)-1]
	pr.stack = pr.stack[:len(pr.stack)-1]
	return out
}

// if `acceptSemiColon` is true, a semi colon at the end
// of the expression is interpreted as EOF (but not consumed)
// a single [ is accepted as EOF
func (pr *parser) parseExpression(acceptSemiColon bool) (*Expr, error) {
	for {
		peeked := pr.tk.Peek().data
		if peeked == nil || peeked == closePar ||
			peeked == closeCurly || peeked == closeBracket || peeked == openBracket ||
			(acceptSemiColon && peeked == semicolon) {
			break
		}

		node, err := pr.parseOneNode(acceptSemiColon)
		if err != nil {
			return nil, err
		}

		pr.stack = append(pr.stack, node)
	}

	if len(pr.stack) != 1 {
		return nil, ErrInvalidExpr{
			Reason: fmt.Sprintf("expression mal formée (%d termes trouvés)", len(pr.stack)),
			Pos:    0,
		}
	}

	return pr.pop(), nil
}

// the next token has already been checked for emptyness,
// and is assumed not to be a closing delimiter )
func (pr *parser) parseOneNode(acceptSemiColon bool) (*Expr, error) {
	tok := pr.tk.Next()
	c := tok.data
	switch data := c.(type) {
	case symbol:
		switch data {
		case openPar:
			return pr.parseParenthesisBlock(tok.pos)
		case closePar:
			panic("internal error")
		case openMatrix: // matrix/vector start
			return pr.parseMatrix(tok.pos)
		case semicolon, openBracket, closeBracket, openCurly, closeCurly:
			return nil, ErrInvalidExpr{
				Reason: fmt.Sprintf("symbole %s inattendu", data),
				Pos:    tok.pos,
			}
		case underscore:
			// we have an indice
			return pr.parseIndice(tok.pos)
		default:
			panic(exhaustiveSymbolSwitch)
		}
	case operator:
		return pr.parseOperator(data, tok.pos, acceptSemiColon)
	case roundFunc:
		return pr.parseRoundFunction(tok.pos)
	case specialFunctionKind:
		rd, err := pr.parseSpecialFunction(tok.pos, data)
		if err != nil {
			return nil, err
		}
		return &Expr{atom: rd}, nil
	case function:
		return pr.parseFunction(data, tok.pos)
	case constant:
		return &Expr{atom: data}, nil
	case Variable:
		// register the variable position
		pr.variablePos[tok.pos] = data

		return &Expr{atom: data}, nil
	case nT:
		nb, err := parseNumber(data, tok.pos)
		if err != nil {
			return nil, err
		}
		return &Expr{atom: nb}, nil
	default:
		panic(exhaustiveTokenSwitch)
	}
}

// accepts <left>_{<right>}, where <left> might only be a one letter name
// the tokenizer is at _
func (pr *parser) parseIndice(pos int) (*Expr, error) {
	// pop the left member
	left := pr.pop()
	if left == nil {
		return nil, ErrInvalidExpr{
			Reason: "variable manquante avant l'indice _",
			Pos:    pos,
		}
	}

	if v, isVar := left.atom.(Variable); !(isVar && v.Indice == "") {
		return nil, ErrInvalidExpr{
			Reason: "seule une variable simple est acceptée avant l'indice _",
			Pos:    pos,
		}
	}

	// check and consume the starting {
	// the tokenizer made sure it is present
	pr.tk.Next()

	// read the full indice expression
	right, err := pr.parseExpression(false)
	if err != nil {
		return nil, err
	}

	// check and consume the closing }
	tk := pr.tk.Next()
	if tk.data != closeCurly {
		return nil, ErrInvalidExpr{
			Reason: "accolade fermante manquante après _",
			Pos:    tk.pos,
		}
	}

	return &Expr{left: left, atom: indice{}, right: right}, nil
}

// pr.src[pr.pos] must be an operator
func (pr *parser) parseOperator(op operator, pos int, acceptSemiColon bool) (*Expr, error) {
	_ = exhaustiveOperatorSwitch

	var leftIsOptional bool
	switch op {
	case plus, minus: // an expression before the sign is optional
		leftIsOptional = true
	}

	// pop the left member
	left := pr.pop()
	if left == nil && !leftIsOptional {
		return nil, ErrInvalidExpr{
			Reason: fmt.Sprintf("expression manquante avant l'opération %s", op),
			Pos:    pos,
		}
	}

	// parse the remaining expression, except for factorial
	if op == factorial {
		return &Expr{
			atom: op,
			left: left,
		}, nil
	}

	// since a^b^c = a^(b^c) and not (a^b)^c
	// adjust operator precedence threshold
	stopOp := op
	if op == pow {
		stopOp = pow - 1
	}

	var (
		right *Expr
		err   error
	)
	// special case ..^-<exponent>, quite natural in practice
	if op == pow && pr.tk.Peek().data == minus {
		_ = pr.tk.Next() // consume minus
		if next := pr.tk.Peek(); next.data == nil || next.data == closePar {
			return nil, ErrInvalidExpr{
				Reason: fmt.Sprintf("expression manquante après une puissance négative"),
				Pos:    pos,
			}
		}
		expr, err := pr.parseOneNode(acceptSemiColon)
		if err != nil {
			return nil, err
		}
		right = &Expr{atom: minus, right: expr}
	} else {
		right, err = pr.parseUntil(stopOp, acceptSemiColon)
		if err != nil {
			return nil, err
		}
	}

	if right == nil {
		return nil, ErrInvalidExpr{
			Reason: fmt.Sprintf("expression manquante après l'opération %s", op),
			Pos:    pos,
		}
	}

	return &Expr{
		atom:  op,
		left:  left,
		right: right,
	}, nil
}

// parse while the found operator has strictly higher precedence than `op`,
// that is returns if `found <= op`
func (pr *parser) parseUntil(op operator, acceptSemiColon bool) (*Expr, error) {
	for {
		tok := pr.tk.Peek()
		// if we reach EOF, return
		// same if we encouter a closing )
		// such as log(  2 + x  )
		if tok.data == closePar ||
			tok.data == closeBracket || tok.data == openBracket || tok.data == closeCurly ||
			tok.data == nil || (acceptSemiColon && tok.data == semicolon) {
			break
		}

		// if the next token is an operator with same or lower precedence
		// stop the parsing here
		if tokOperator, isOp := tok.data.(operator); isOp && tokOperator <= op {
			break
		}

		node, err := pr.parseOneNode(acceptSemiColon)
		if err != nil {
			return nil, err
		}

		pr.stack = append(pr.stack, node)
	}

	return pr.pop(), nil
}

// assume that ( token has already been read
func (pr *parser) parseParenthesisBlock(pos int) (*Expr, error) {
	// parse the content, until the closing )

	for {
		tok := pr.tk.Peek()
		if tok.data == closePar {
			// consume the closing ')'
			pr.tk.Next()
			out := pr.pop()
			return out, nil
		}

		if tok.data == nil { // unexpected EOF
			return nil, ErrInvalidExpr{
				Reason: "parenthèse fermante manquante",
				Pos:    pos,
			}
		}

		arg, err := pr.parseOneNode(false)
		if err != nil {
			return nil, err
		}

		pr.stack = append(pr.stack, arg)
	}
}

func (pr *parser) parseFunction(fn function, pos int) (*Expr, error) {
	// after a function name, their must be a (
	// with optional whitespaces
	par := pr.tk.Next()

	if par.data != openPar {
		return nil, ErrInvalidExpr{
			Reason: "parenthèse ouvrante manquante après une fonction",
			Pos:    pos,
		}
	}

	// parse the argument
	arg, err := pr.parseParenthesisBlock(pos)
	if err != nil {
		return nil, err
	}

	if arg == nil {
		return nil, ErrInvalidExpr{
			Reason: fmt.Sprintf("argument manquant pour la fonction %s", fn),
			Pos:    pos,
		}
	}

	return &Expr{
		left:  nil,
		right: arg,
		atom:  fn,
	}, nil
}

func parseNumber(v nT, pos int) (Number, error) {
	out, err := strconv.ParseFloat(string(v), 64)
	if err != nil {
		return 0, ErrInvalidExpr{
			Reason: fmt.Sprintf(`syntaxe "%s" non reconnue pour un nombre`, v),
			Pos:    pos,
		}
	}

	return Number(out), nil
}

// also accept a negative value
func (pr *parser) parsePositiveInt() (int, error) {
	arg := pr.tk.Next()

	v, ok := arg.data.(nT)
	if !ok {
		return 0, ErrInvalidExpr{
			Reason: "nombre attendu",
			Pos:    arg.pos,
		}
	}

	n, err := parseNumber(v, arg.pos)
	if err != nil {
		return 0, err
	}

	nI, ok := IsInt(float64(n))
	if !ok {
		return 0, ErrInvalidExpr{
			Reason: "nombre entier attendu",
			Pos:    arg.pos,
		}
	}

	return nI, nil
}

// IsInt returns `true` if `v` is a finite integer number.
func IsInt(v float64) (int, bool) {
	out := int(v)
	if float64(out) == v {
		return out, true
	}
	return 0, false
}

// special case for the round function,
// which accept one expression and number of digits
func (pr *parser) parseRoundFunction(pos int) (expr *Expr, err error) {
	// after a function name, their must be a (
	// with optional whitespaces
	par := pr.tk.Next()
	if par.data != openPar {
		return nil, ErrInvalidExpr{
			Reason: "parenthèse ouvrante manquante après round",
			Pos:    pos,
		}
	}

	// we then accept a whole expression
	arg, err := pr.parseExpression(true)
	if err != nil {
		return nil, err
	}

	// we then look for a ; token ...
	if tok := pr.tk.Next(); tok.data != semicolon {
		return nil, ErrInvalidExpr{
			Reason: "point-virgule manquant entre les arguments de round",
			Pos:    tok.pos,
		}
	}

	// ... and finally for a positive integer
	nbDigitsI, err := pr.parsePositiveInt()
	if err != nil {
		return nil, err
	}

	if tok := pr.tk.Next(); tok.data != closePar {
		return nil, ErrInvalidExpr{
			Reason: "parenthèse fermante manquante après round",
			Pos:    tok.pos,
		}
	}

	return &Expr{atom: roundFunc{nbDigits: nbDigitsI}, right: arg}, nil
}

// parses a semi colon separated list of expressions, followed by a closing delimiter,
// such as :
//
//		<expr> ; <expr> ; <expr> }
//		<expr>; <expr> [
//	 <expr>; <expr>)
//
// [isCloser] must return [true] for symbol which are valid closing delimiters
//
// Note that the oppening symbol must have been read.
// The closing delimiter is consumed and returned.
func (pr *parser) parseExpressionList(isCloser func(symbol) bool) ([]*Expr, symbol, error) {
	var args []*Expr
	for {
		t := pr.tk.Peek().data
		if ts, isSymbol := t.(symbol); t == nil || (isSymbol && isCloser(ts)) {
			break
		}

		// we then accept a whole expression
		arg, err := pr.parseExpression(true)
		if err != nil {
			return nil, 0, err
		}
		args = append(args, arg)

		// four cases here :
		//	- valid closing delimiter, closed after the loop
		//	- valid separator ;
		//  - invalid character, probably an invalid separator like ,
		//	- EOF
		next := pr.tk.Peek()
		ns, isSymbol := next.data.(symbol)
		switch {
		case isSymbol && isCloser(ns):
			continue // without consuming
		case next.data == semicolon: // consume the separator and continue
			_ = pr.tk.Next()
		case next.data == nil: // EOF, error reported after the loop

		default: // invalid separator
			return nil, 0, ErrInvalidExpr{
				Reason: "le séparateur attendu est ';'",
				Pos:    next.pos,
			}
		}
	}

	tok := pr.tk.Next()
	ts, isSymbol := tok.data.(symbol)
	if !(isSymbol && isCloser(ts)) {
		return nil, 0, ErrInvalidExpr{
			Reason: "délimiteur manquant",
			Pos:    tok.pos,
		}
	}

	return args, ts, nil
}

// special case for special functions, which are of the form
// <function>(<expr>; <expr> ...)
func (pr *parser) parseSpecialFunction(pos int, fn specialFunctionKind) (rd specialFunction, err error) {
	// after a function name, their must be a (
	// with optional whitespaces
	if par := pr.tk.Next(); par.data != openPar {
		return rd, ErrInvalidExpr{
			Reason: fmt.Sprintf("parenthèse ouvrante manquante après %s", fn),
			Pos:    pos,
		}
	}

	rd.kind = fn

	rd.args, _, err = pr.parseExpressionList(func(s symbol) bool { return s == closePar })
	if err != nil {
		return rd, err
	}

	err = rd.validate(pos)

	return rd, err
}

func (rd specialFunction) validate(pos int) error {
	switch rd.kind {
	case randInt, randPrime:
		if len(rd.args) < 2 {
			return ErrInvalidExpr{
				Reason: fmt.Sprintf("%s attend deux paramètres", rd.kind),
				Pos:    pos,
			}
		}

		// eagerly try to eval start and end in case their are constant,
		// so that the error is detected during parameter setup
		start, end, err := startEnd(rd.args[0], rd.args[1], nil)
		if err == nil {
			return rd.validateStartEnd(start, end, pos)
		}
	case randChoice:
		if len(rd.args) == 0 {
			return ErrInvalidExpr{
				Reason: "randChoice doit préciser au moins un argument",
				Pos:    pos,
			}
		}
	case choiceFrom:
		if len(rd.args) < 2 {
			return ErrInvalidExpr{
				Reason: "choiceFrom doit préciser au moins deux arguments",
				Pos:    pos,
			}
		}
	case randDenominator: // nothing to validate

	case minFn, maxFn:
		if len(rd.args) == 0 {
			return ErrInvalidExpr{
				Reason: "min et max requierent au moins un argument",
				Pos:    pos,
			}
		}
	case sumFn:
		if len(rd.args) != 4 {
			return ErrInvalidExpr{
				Reason: "sum requiert exactement 4 arguments",
				Pos:    pos,
			}
		}

		// eagerly try to eval start and end in case their are constant,
		// so that a potential error is detected during parameter setup
		start, end, err := startEnd(rd.args[1], rd.args[2], nil)
		if err == nil {
			return rd.validateStartEnd(start, end, pos)
		}
	case matCoeff:
		if len(rd.args) != 3 {
			return ErrInvalidExpr{
				Reason: "coeff requiert exactement 3 arguments",
				Pos:    pos,
			}
		}
	case binomial:
		if len(rd.args) != 2 {
			return ErrInvalidExpr{
				Reason: "binom requiert exactement 2 arguments",
				Pos:    pos,
			}
		}
	default:
		panic(exhaustiveSpecialFunctionSwitch)
	}

	return nil
}

// assume the first [[ has been read
func (pr *parser) parseMatrix(pos int) (*Expr, error) {
	isFirst := true
	var mt matrix
	for { // parse each rows
		if !isFirst {
			if open := pr.tk.Next(); open.data != openBracket {
				return nil, ErrInvalidExpr{
					Reason: fmt.Sprintf("ligne d'une matrice attendue (crochet ouvrant manquant)"),
					Pos:    open.pos,
				}
			}
		}

		row, _, err := pr.parseExpressionList(func(s symbol) bool { return s == closeBracket })
		if err != nil {
			return nil, err
		}

		mt = append(mt, row)
		isFirst = false

		// two valid cases here : "]" to end the matrix, ";" to add a new row
		next := pr.tk.Next()
		if next.data == semicolon {
			continue
		} else if next.data == closeBracket {
			break
		} else {
			return nil, ErrInvalidExpr{
				Reason: fmt.Sprintf("] ou ; attendu dans une matrice"),
				Pos:    next.pos,
			}
		}
	}

	// validate rows length : by construction, there is always a row
	L := len(mt[0])
	for _, row := range mt {
		if len(row) != L {
			return nil, ErrInvalidExpr{
				Reason: fmt.Sprintf("les lignes d'une matrice doivent être de même taille (%d != %d)", len(row), L),
				Pos:    pos,
			}
		}
	}

	return &Expr{atom: mt}, nil
}
