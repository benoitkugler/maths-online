package expression

import (
	"fmt"
	"sort"
	"strconv"
)

// ErrInvalidExpr is returned by `Parse`.
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
func (vm varMap) Positions(subset RandomParameters) []int {
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
func Parse(s string) (*Expression, error) {
	expr, _, err := parseBytes([]byte(s))
	if err != nil {
		errV := err.(ErrInvalidExpr)
		errV.Input = s
		return nil, errV
	}
	return expr, nil
}

// MustParse is the same as Parse but panics on invalid expressions.
func MustParse(s string) *Expression {
	expr, err := Parse(s)
	if err != nil {
		panic(fmt.Sprintf("%s: %s", s, err))
	}
	return expr
}

// parseBytes parses a mathematical expression. If invalid, an `InvalidExpr` is returned.
func parseBytes(text []byte) (*Expression, varMap, error) {
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

	stack []*Expression // waiting to be consumed by operators
}

func newParser(text []byte) *parser {
	return &parser{tk: newTokenizer(text), variablePos: make(map[int]Variable)}
}

// return nil if the stack is empty
func (pr *parser) pop() *Expression {
	if len(pr.stack) == 0 {
		return nil
	}
	out := pr.stack[len(pr.stack)-1]
	pr.stack = pr.stack[:len(pr.stack)-1]
	return out
}

// if `acceptSemiColon` is true, a semi colon at the end
// of the expression is interpreted as EOF (but not consumed)
func (pr *parser) parseExpression(acceptSemiColon bool) (*Expression, error) {
	for {
		peeked := pr.tk.Peek().data
		if peeked == nil || peeked == closePar || (acceptSemiColon && peeked == semicolon) {
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
// and is assumed not to be a closing )
func (pr *parser) parseOneNode(acceptSemiColon bool) (*Expression, error) {
	tok := pr.tk.Next()
	c := tok.data
	switch data := c.(type) {
	case symbol:
		switch data {
		case openPar:
			return pr.parseParenthesisBlock(tok.pos)
		case closePar:
			panic("internal error")
		case semicolon:
			return nil, ErrInvalidExpr{
				Reason: "point-virgule inattendu",
				Pos:    tok.pos,
			}
		default:
			panic(exhaustiveSymbolSwitch)
		}
	case operator:
		return pr.parseOperator(data, tok.pos, acceptSemiColon)
	case randVariable:
		rd, err := pr.parseRandVariable(tok.pos)
		if err != nil {
			return nil, err
		}
		return &Expression{atom: rd}, nil
	case roundFn:
		return pr.parseRoundFunction(tok.pos)
	case specialFunction:
		rd, err := pr.parseSpecialFunction(tok.pos, data)
		if err != nil {
			return nil, err
		}
		return &Expression{atom: rd}, nil
	case function:
		return pr.parseFunction(data, tok.pos)
	case constant:
		return &Expression{atom: data}, nil
	case Variable:
		// register the variable position
		pr.variablePos[tok.pos] = data

		return &Expression{atom: data}, nil
	case numberText:
		nb, err := parseNumber(data, tok.pos)
		if err != nil {
			return nil, err
		}
		return &Expression{atom: nb}, nil
	default:
		panic(exhaustiveTokenSwitch)
	}
}

// pr.src[pr.pos] must be an operator
func (pr *parser) parseOperator(op operator, pos int, acceptSemiColon bool) (*Expression, error) {
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

	// parse the remaining expression

	// since a^b^c = a^(b^c) and not (a^b)^c
	// adjust operator precedence
	stopOp := op
	if op == pow {
		stopOp = mult
	}
	right, err := pr.parseUntil(stopOp, acceptSemiColon)
	if err != nil {
		return nil, err
	}

	if right == nil {
		return nil, ErrInvalidExpr{
			Reason: fmt.Sprintf("expression manquante après l'opération %s", op),
			Pos:    pos,
		}
	}

	return &Expression{
		atom:  op,
		left:  left,
		right: right,
	}, nil
}

// parse while the operator have strictly higher precedence than `op`
func (pr *parser) parseUntil(op operator, acceptSemiColon bool) (*Expression, error) {
	for {
		tok := pr.tk.Peek()
		// if we reach EOF, return
		// same if we encouter a closing )
		// such as log(  2 + x  )
		if tok.data == closePar || tok.data == nil || (acceptSemiColon && tok.data == semicolon) {
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
func (pr *parser) parseParenthesisBlock(pos int) (*Expression, error) {
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

func (pr *parser) parseFunction(fn function, pos int) (*Expression, error) {
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

	return &Expression{
		left:  nil,
		right: arg,
		atom:  fn,
	}, nil
}

func parseNumber(v numberText, pos int) (Number, error) {
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

	v, ok := arg.data.(numberText)
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

	nI, ok := isInt(float64(n))
	if !ok {
		return 0, ErrInvalidExpr{
			Reason: "nombre entier attendu",
			Pos:    arg.pos,
		}
	}

	return nI, nil
}

func isInt(v float64) (int, bool) {
	out := int(v)
	if float64(out) == v {
		return out, true
	}
	return 0, false
}

func (pr *parser) parseRandVariable(pos int) (rv randVariable, err error) {
	// after a function name, their must be a (
	// with optional whitespaces
	par := pr.tk.Next()

	if par.data != openPar {
		return nil, ErrInvalidExpr{
			Reason: "parenthèse ouvrante manquante après randLetter",
			Pos:    pos,
		}
	}

	for {
		if t := pr.tk.Peek().data; t == closePar || t == nil {
			break
		}

		arg := pr.tk.Next()
		if v, isVariable := arg.data.(Variable); isVariable {
			rv = append(rv, v)
		} else {
			return nil, ErrInvalidExpr{
				Reason: "randLetter n'accepte que des variables comme arguments",
				Pos:    pos,
			}
		}

		// four cases here :
		//	- valid closing parenthesis, closed after the loop
		//	- valid separator ;
		//  - invalid character, probably an invalid separator like ,
		//	- EOF
		switch next := pr.tk.Peek().data; next {
		case closePar:
			continue // without consuming
		case semicolon: // consume and continue
			_ = pr.tk.Next()
		case nil: // EOF, error reported after the loop
			break
		default: // invalid separator
			return nil, ErrInvalidExpr{
				Reason: "randLetter utilise ';' comme séparateur",
				Pos:    pos,
			}
		}
	}

	if tok := pr.tk.Next(); tok.data != closePar {
		return nil, ErrInvalidExpr{
			Reason: "parenthèse fermante manquante après randLetter",
			Pos:    tok.pos,
		}
	}

	if len(rv) == 0 {
		return nil, ErrInvalidExpr{
			Reason: "randLetter attend au moins un argument",
			Pos:    pos,
		}
	}

	return rv, nil
}

// special case for the round function,
// which accept one expression and number of digits
func (pr *parser) parseRoundFunction(pos int) (expr *Expression, err error) {
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

	return &Expression{atom: roundFn{nbDigits: nbDigitsI}, right: arg}, nil
}

// special case for special functions, which are of the form
// <function>(<expr>; <expr> ...)
func (pr *parser) parseSpecialFunction(pos int, fn specialFunction) (rd specialFunctionA, err error) {
	// after a function name, their must be a (
	// with optional whitespaces
	par := pr.tk.Next()
	if par.data != openPar {
		return rd, ErrInvalidExpr{
			Reason: "parenthèse ouvrante manquante après randXXX",
			Pos:    pos,
		}
	}

	rd.kind = fn

	for {
		if t := pr.tk.Peek().data; t == closePar || t == nil {
			break
		}

		// we then accept a whole expression
		arg, err := pr.parseExpression(true)
		if err != nil {
			return rd, err
		}
		rd.args = append(rd.args, arg)

		if pr.tk.Peek().data == closePar {
			continue
		}

		// consume the separator
		if pr.tk.Next().data != semicolon {
			return rd, ErrInvalidExpr{
				Reason: "randXXX utilise ';' comme séparateur",
				Pos:    pos,
			}
		}
	}

	if tok := pr.tk.Next(); tok.data != closePar {
		return rd, ErrInvalidExpr{
			Reason: "parenthèse fermante manquante après randXXX",
			Pos:    tok.pos,
		}
	}

	err = rd.validate(pos)

	return rd, err
}
