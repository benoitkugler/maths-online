package expression

import (
	"fmt"
	"strconv"
	"unicode"
)

type InvalidExpr struct {
	Reason string
	Pos    int
}

func (inv InvalidExpr) Error() string {
	return fmt.Sprintf("invalid expression: at %d: %s", inv.Pos, inv.Reason)
}

func parseExpression(s string) (*node, error) {
	pr := parser{src: []rune(s)}
	return pr.parseExpression()
}

// parser transforms an utf8 byte input
// to an ast
type parser struct {
	stack []*node // wainting to be consumed by operators

	src []rune
	pos int
}

// return nil if the stack is empty
func (pr *parser) pop() *node {
	if len(pr.stack) == 0 {
		return nil
	}
	out := pr.stack[len(pr.stack)-1]
	pr.stack = pr.stack[:len(pr.stack)-1]
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

func isOperator(r rune) bool {
	switch r {
	case '+', '-', '*', '/', '^', '\u00F7', '\u00D7':
		return true
	default:
		return false
	}
}

// advance pos to the next non whitespace char
func (pr *parser) skipWhitespace() {
	for pr.pos < len(pr.src) && isWhiteSpace(pr.src[pr.pos]) {
		pr.pos++
	}
}

func (pr *parser) parseExpression() (*node, error) {
	L := len(pr.src)

	pr.skipWhitespace()

	for pr.pos < L {
		node, err := pr.parseOneNode()
		if err != nil {
			return nil, err
		}

		pr.stack = append(pr.stack, node)

		pr.skipWhitespace()
	}

	if len(pr.stack) != 1 {
		return nil, InvalidExpr{
			Reason: fmt.Sprintf("expression mal formée (%d termes trouvés)", len(pr.stack)),
			Pos:    0,
		}
	}

	return pr.stack[0], nil
}

func (pr *parser) parseOneNode() (*node, error) {
	if pr.pos >= len(pr.src) {
		return nil, InvalidExpr{
			Reason: "expression vide ou incomplète",
			Pos:    pr.pos,
		}
	}

	c := pr.src[pr.pos]
	switch {
	case c == '(':
		return pr.parseParenthesisBlock()
	case c == ')':
		return nil, InvalidExpr{
			Reason: "parenthèse fermante invalide",
			Pos:    pr.pos,
		}
	case isOperator(c): // we have an operator
		return pr.parseOperator()
	case unicode.IsLetter(c): // either a function, a variable or a constant
		fn, isFunction, err := pr.tryParseFunction()
		if err != nil {
			return nil, err
		}
		if isFunction {
			return fn, nil
		}

		// read the symbol as variable or predefined constants
		v := pr.parseConstantOrVariable()
		return &node{atom: v}, nil
	default:
		nb, err := pr.parseNumber()
		if err != nil {
			return nil, err
		}
		return &node{atom: nb}, nil
	}
}

// pr.src[pr.pos] must be an operator
func (pr *parser) parseOperator() (*node, error) {
	var (
		op             operator
		leftIsOptional bool
		c              = pr.src[pr.pos]
	)
	switch c {
	case '+':
		op = plus
		leftIsOptional = true
	case '-':
		op = minus
		leftIsOptional = true
	case '/', '\u00F7':
		op = div
	case '*', '\u00D7':
		op = mult
	case '^':
		op = pow
	}

	// parse the remaining expression
	pr.pos++
	pr.skipWhitespace()
	right, err := pr.parseOneNode()
	if err != nil {
		return nil, err
	}

	// pop the left member
	left := pr.pop()
	if left == nil && !leftIsOptional {
		return nil, InvalidExpr{
			Reason: fmt.Sprintf("expression manquante avant l'opération %s", string(c)),
			Pos:    pr.pos,
		}
	}

	// an expression before the sign is optional
	return &node{
		atom:  op,
		right: right,
		left:  left,
	}, nil
}

// assume pr.src[pr.pos] == '(' and consume the block
func (pr *parser) parseParenthesisBlock() (*node, error) {
	startPos := pr.pos
	// consume the '('
	pr.pos++

	L := len(pr.src)

	// parse the content, until the closing )
	pr.skipWhitespace()
	for pr.pos < L {
		if pr.src[pr.pos] == ')' {
			// consume the closing ')'
			pr.pos++
			out := pr.pop()
			return out, nil
		}

		arg, err := pr.parseOneNode()
		if err != nil {
			return nil, err
		}

		pr.stack = append(pr.stack, arg)

		pr.skipWhitespace()
	}

	return nil, InvalidExpr{
		Reason: "parenthèse fermante manquante",
		Pos:    startPos,
	}
}

func (pr *parser) tryParseFunction() (*node, bool, error) {
	L := len(pr.src)

	// read subsequent letters
	var letters []rune
	for i := pr.pos; i < L && unicode.IsLetter(pr.src[i]); i++ {
		letters = append(letters, pr.src[i])
	}

	var fn function
	switch string(letters) {
	case "exp":
		fn = exp
	case "ln", "log":
		fn = log
	case "sin":
		fn = sin
	case "cos":
		fn = cos
	case "abs":
		fn = abs
	default: // no  matching function name
		return nil, false, nil
	}

	// after a function name, their must be a (
	// with optional whitespaces
	pr.pos += len(letters)
	endFunctionPos := pr.pos
	pr.skipWhitespace()

	if pr.pos >= L || pr.src[pr.pos] != '(' {
		return nil, false, InvalidExpr{
			Reason: "parenthèse ouvrante manquante après une fonction",
			Pos:    endFunctionPos,
		}
	}

	// parse the argument
	arg, err := pr.parseParenthesisBlock()
	if err != nil {
		return nil, false, err
	}

	return &node{
		left:  nil,
		right: arg,
		atom:  fn,
	}, true, nil
}

func (pr *parser) parseConstantOrVariable() (out atom) {
	switch pr.src[pr.pos] {
	case 'e':
		out = numberE
	case '\u03C0':
		out = numberPi
	default:
		out = variable(pr.src[pr.pos])
	}
	pr.pos++
	return out
}

func (pr *parser) parseNumber() (number, error) {
	var buffer []rune
	originalPos, L := pr.pos, len(pr.src)
	for ; pr.pos < L; pr.pos++ {
		r := pr.src[pr.pos]
		if unicode.IsDigit(r) || r == '.' {
			buffer = append(buffer, r)
		} else {
			break
		}
	}

	s := string(buffer)
	out, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, InvalidExpr{
			Reason: fmt.Sprintf("syntaxe %s non reconnue pour un nombre", s),
			Pos:    originalPos,
		}
	}

	return number(out), nil
}
