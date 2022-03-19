package expression

import (
	"fmt"
	"sort"
	"strconv"
)

// InvalidExpr is returned by `Parse`.
type InvalidExpr struct {
	Reason string // in french
	Pos    int    // index in the input rune slice
}

func (inv InvalidExpr) Error() string {
	return fmt.Sprintf("expression invalide : position %d : %s", inv.Pos, inv.Reason)
}

func (inv InvalidExpr) PortionOf(expr string) string {
	return string([]rune(expr)[:inv.Pos])
}

// VarMap register the position of each variable occurence in
// an expression rune slice.
type VarMap map[int]Variable

// Positions returns the indices in the original expression rune slice
// of every occurence of the given variables
func (vm VarMap) Positions(subset RandomParameters) []int {
	var out []int
	for index, v := range vm {
		if _, has := subset[v]; has {
			out = append(out, index)
		}
	}

	sort.Ints(out)

	return out
}

// Parse parses a mathematical expression. If invalid, an `InvalidExpr` is returned.
func Parse(s string) (*Expression, VarMap, error) {
	return parseBytes([]byte(s))
}

// parseBytes parses a mathematical expression. If invalid, an `InvalidExpr` is returned.
func parseBytes(text []byte) (*Expression, VarMap, error) {
	pr := newParser(text)
	e, err := pr.parseExpression()

	return e, pr.variablePos, err
}

// parser transforms an utf8 byte input into an ast
type parser struct {
	tk *tokenizer

	variablePos VarMap // index of variable in input rune slice

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

func (pr *parser) parseExpression() (*Expression, error) {
	for pr.tk.peek().data != nil {
		node, err := pr.parseOneNode()
		if err != nil {
			return nil, err
		}

		pr.stack = append(pr.stack, node)
	}

	if len(pr.stack) != 1 {
		return nil, InvalidExpr{
			Reason: fmt.Sprintf("expression mal formée (%d termes trouvés)", len(pr.stack)),
			Pos:    0,
		}
	}

	return pr.stack[0], nil
}

// the next token has already been checked for emptyness
func (pr *parser) parseOneNode() (*Expression, error) {
	tok := pr.tk.next()
	c := tok.data
	switch data := c.(type) {
	case symbol:
		switch data {
		case openPar:
			return pr.parseParenthesisBlock(tok.pos)
		case closePar:
			return nil, InvalidExpr{
				Reason: "parenthèse fermante invalide",
				Pos:    tok.pos,
			}
		case semicolon:
			return nil, InvalidExpr{
				Reason: "virgule inattendue",
				Pos:    tok.pos,
			}
		default:
			panic(exhaustiveSymbolSwitch)
		}
	case operator:
		return pr.parseOperator(data, tok.pos)
	case randFunction:
		rd, err := pr.parseRandFunction(tok.pos, data.isPrime)
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
func (pr *parser) parseOperator(op operator, pos int) (*Expression, error) {
	var leftIsOptional bool
	switch op {
	case plus, minus: // an expression before the sign is optional
		leftIsOptional = true
	}

	// pop the left member
	left := pr.pop()
	if left == nil && !leftIsOptional {
		return nil, InvalidExpr{
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
	right, err := pr.parseUntil(stopOp)
	if err != nil {
		return nil, err
	}

	if right == nil {
		return nil, InvalidExpr{
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
func (pr *parser) parseUntil(op operator) (*Expression, error) {
	for {
		tok := pr.tk.peek()
		// if we reach EOF, return
		// same if we encouter a closing )
		// such as log(  2 + x  )
		if tok.data == closePar || tok.data == nil {
			break
		}

		// if the next token is an operator with same or lower precedence
		// stop the parsing here
		if tokOperator, isOp := tok.data.(operator); isOp && tokOperator <= op {
			break
		}

		node, err := pr.parseOneNode()
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
		tok := pr.tk.peek()
		if tok.data == closePar {
			// consume the closing ')'
			pr.tk.next()
			out := pr.pop()
			return out, nil
		}

		if tok.data == nil { // unexpected EOF
			return nil, InvalidExpr{
				Reason: "parenthèse fermante manquante",
				Pos:    pos,
			}
		}

		arg, err := pr.parseOneNode()
		if err != nil {
			return nil, err
		}

		pr.stack = append(pr.stack, arg)
	}
}

func (pr *parser) parseFunction(fn function, pos int) (*Expression, error) {
	// after a function name, their must be a (
	// with optional whitespaces
	par := pr.tk.next()

	if par.data != openPar {
		return nil, InvalidExpr{
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
		return nil, InvalidExpr{
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
		return 0, InvalidExpr{
			Reason: fmt.Sprintf(`syntaxe "%s" non reconnue pour un nombre`, v),
			Pos:    pos,
		}
	}

	return Number(out), nil
}

// accept a possibly negative integer
func (pr *parser) parseInt() (int, error) {
	arg := pr.tk.next()

	var isNegative bool
	if arg.data == minus { // read another token
		isNegative = true
		arg = pr.tk.next()
	}

	v, ok := arg.data.(numberText)
	if !ok {
		return 0, InvalidExpr{
			Reason: "nombre entier attendu",
			Pos:    arg.pos,
		}
	}

	n, err := parseNumber(v, arg.pos)
	if err != nil {
		return 0, err
	}

	if _, isInt := isInt(float64(n)); !isInt {
		return 0, InvalidExpr{
			Reason: "nombre entier attendu",
			Pos:    arg.pos,
		}
	}

	out := int(n)
	if isNegative {
		out = -out
	}
	return out, nil
}

func isInt(v float64) (int, bool) {
	out := int(v)
	if float64(out) == v {
		return out, true
	}
	return 0, false
}

// special case for randInt or randPrime
// to keep the parser simple, we only accept randXXX(number, number)
func (pr *parser) parseRandFunction(pos int, isPrime bool) (rd random, err error) {
	// after a function name, their must be a (
	// with optional whitespaces
	par := pr.tk.next()

	if par.data != openPar {
		return random{}, InvalidExpr{
			Reason: "parenthèse ouvrante manquante après randInt",
			Pos:    pos,
		}
	}

	rd.isPrime = isPrime

	rd.start, err = pr.parseInt()
	if err != nil {
		return random{}, err
	}

	if tok := pr.tk.next(); tok.data != semicolon {
		return random{}, InvalidExpr{
			Reason: "virgule manquant entre les arguments de randInt",
			Pos:    tok.pos,
		}
	}

	rd.end, err = pr.parseInt()
	if err != nil {
		return random{}, err
	}

	if tok := pr.tk.next(); tok.data != closePar {
		return random{}, InvalidExpr{
			Reason: "parenthèse fermante manquante après randInt",
			Pos:    tok.pos,
		}
	}

	err = rd.validate(pos)

	return rd, err
}
