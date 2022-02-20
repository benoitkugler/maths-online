package expression

import (
	"fmt"
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

// Parse parses a mathematical expression. If invalid, an `InvalidExpr` is returned.
func Parse(s string) (*Expression, error) {
	pr := parser{tk: newTokenizer(s)}
	return pr.parseExpression()
}

// parser transforms an utf8 byte input
// to an ast
type parser struct {
	tk *tokenizer

	stack []*Expression // waiting to be consumed by operators
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
		default:
			panic("unknown symbol")
		}
	case operator:
		return pr.parseOperator(data, tok.pos)
	case function:
		return pr.parseFunction(data, tok.pos)
	case constant:
		return &Expression{atom: data}, nil
	case Variable:
		return &Expression{atom: data}, nil
	case numberText:
		nb, err := pr.parseNumber(data, tok.pos)
		if err != nil {
			return nil, err
		}
		return &Expression{atom: nb}, nil
	default:
		panic("unknown token type")
	}
}

// pr.src[pr.pos] must be an operator
func (pr *parser) parseOperator(op operator, pos int) (*Expression, error) {
	var leftIsOptional bool
	switch op {
	case plus, minus:
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
	right, err := pr.parseUntil(op)
	if err != nil {
		return nil, err
	}

	if right == nil {
		return nil, InvalidExpr{
			Reason: fmt.Sprintf("expression manquante après l'opération %s", op),
			Pos:    pos,
		}
	}

	// an expression before the sign is optional
	return &Expression{
		atom:  op,
		right: right,
		left:  left,
	}, nil
}

// parse while the operator have higher precedence than `op`
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

	return &Expression{
		left:  nil,
		right: arg,
		atom:  fn,
	}, nil
}

func (pr *parser) parseNumber(v numberText, pos int) (number, error) {
	out, err := strconv.ParseFloat(string(v), 64)
	if err != nil {
		return 0, InvalidExpr{
			Reason: fmt.Sprintf("syntaxe %s non reconnue pour un nombre", v),
			Pos:    pos,
		}
	}

	return number(out), nil
}
