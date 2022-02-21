package expression

import (
	"reflect"
	"testing"
)

var expressions = []struct {
	expr    string
	want    *Expression
	wantErr bool
}{
	{
		"1 + (x + z) * (x + 1)", &Expression{
			atom: plus,
			left: newNumber(1),
			right: &Expression{
				atom:  mult,
				left:  &Expression{atom: plus, left: newVariable('x'), right: newVariable('z')},
				right: &Expression{atom: plus, left: newVariable('x'), right: newNumber(1)},
			},
		}, false,
	},
	{
		"1 + 2 * 3 * 4", &Expression{
			atom: plus,
			left: newNumber(1),
			right: &Expression{
				atom:  mult,
				left:  &Expression{atom: mult, left: newNumber(2), right: newNumber(3)},
				right: newNumber(4),
			},
		}, false,
	},
	{
		"1 + (x + z) * (x + 1) * z ", &Expression{
			atom: plus,
			left: newNumber(1),
			right: &Expression{
				atom: mult,
				left: &Expression{
					atom:  mult,
					left:  &Expression{atom: plus, left: newVariable('x'), right: newVariable('z')},
					right: &Expression{atom: plus, left: newVariable('x'), right: newNumber(1)},
				},
				right: newVariable('z'),
			},
		}, false,
	},
	{
		"x", newVariable('x'), false,
	},
	{
		"t * q", &Expression{atom: mult, left: newVariable('t'), right: newVariable('q')}, false,
	},
	{
		" x ", newVariable('x'), false,
	},
	{
		" 3 ", newNumber(3), false,
	},
	{
		" 3.14 ", newNumber(3.14), false,
	},
	{
		" e ", &Expression{atom: eConstant}, false,
	},
	// custom variables
	{
		" \uE000 + 2", &Expression{atom: plus, left: newVariable('\uE000'), right: newNumber(2)}, false,
	},
	{
		" \u03C0 ", &Expression{atom: piConstant}, false,
	},
	{
		" (e) ", &Expression{atom: eConstant}, false,
	},
	{
		" ((e)) ", &Expression{atom: eConstant}, false,
	},
	{
		"", nil, true,
	},
	{
		"2  +", nil, true,
	},
	{
		"3..", nil, true,
	},
	{
		" / 4", nil, true,
	},
	{
		"2 + 1)", nil, true,
	},
	{
		"2 + 1 * (1 + ", nil, true,
	},
	{
		" ( ", nil, true,
	},
	{
		" () ( ", nil, true,
	},
	{
		" (( ) ", nil, true,
	},
	{
		" ( y + (1 / 2 ) ", nil, true,
	},
	{
		" (1+ 2 ) (x + 3) ", nil, true,
	},
	{
		" log 3 ", nil, true,
	},
	{
		" log( /2)", nil, true,
	},
	{
		" x + 3 ", &Expression{atom: plus, left: newVariable('x'), right: newNumber(3)}, false,
	},
	{
		" x - 3 ", &Expression{atom: minus, left: newVariable('x'), right: newNumber(3)}, false,
	},
	{
		" x ^ 3 ", &Expression{atom: pow, left: newVariable('x'), right: newNumber(3)}, false,
	},
	// ^ is not associative !
	{
		" a^b^c ", &Expression{
			atom: pow,
			left: newVariable('a'),
			right: &Expression{
				atom:  pow,
				left:  newVariable('b'),
				right: newVariable('c'),
			},
		}, false,
	},
	{
		" (a^b)^c ", &Expression{
			atom: pow,
			left: &Expression{
				atom:  pow,
				left:  newVariable('a'),
				right: newVariable('b'),
			},
			right: newVariable('c'),
		}, false,
	},
	{
		" a^b  * c ", &Expression{
			atom: mult,
			left: &Expression{
				atom:  pow,
				left:  newVariable('a'),
				right: newVariable('b'),
			},
			right: newVariable('c'),
		}, false,
	},
	{
		" x / 3 ", &Expression{atom: div, left: newVariable('x'), right: newNumber(3)}, false,
	},
	{
		" x * 3 ", &Expression{atom: mult, left: newVariable('x'), right: newNumber(3)}, false,
	},
	{
		"(x + 3 )", &Expression{atom: plus, left: newVariable('x'), right: newNumber(3)}, false,
	},
	{
		"1 + 2 * 3", &Expression{
			atom:  plus,
			left:  newNumber(1),
			right: &Expression{atom: mult, left: newNumber(2), right: newNumber(3)},
		}, false,
	},
	{
		"1 + 2 * 3 ^ 4", &Expression{
			atom: plus,
			left: newNumber(1),
			right: &Expression{
				atom: mult,
				left: newNumber(2),
				right: &Expression{
					atom:  pow,
					left:  newNumber(3),
					right: newNumber(4),
				},
			},
		}, false,
	},
	{
		"1 + (x + z) * (x + 1) * z ", &Expression{
			atom: plus,
			left: newNumber(1),
			right: &Expression{
				atom: mult,
				left: &Expression{
					atom:  mult,
					left:  &Expression{atom: plus, left: newVariable('x'), right: newVariable('z')},
					right: &Expression{atom: plus, left: newVariable('x'), right: newNumber(1)},
				},
				right: newVariable('z'),
			},
		}, false,
	},
	{
		"1 + 2 * (x + 1)", &Expression{
			atom:  plus,
			left:  newNumber(1),
			right: &Expression{atom: mult, left: newNumber(2), right: &Expression{atom: plus, left: newVariable('x'), right: newNumber(1)}},
		}, false,
	},
	{
		"3 - (x + 3 )", &Expression{atom: minus, left: newNumber(3), right: &Expression{atom: plus, left: newVariable('x'), right: newNumber(3)}}, false,
	},
	{
		"(1 + y) / (3 - (x + 3 ))",

		&Expression{
			atom:  div,
			left:  &Expression{atom: plus, left: newNumber(1), right: newVariable('y')},
			right: &Expression{atom: minus, left: newNumber(3), right: &Expression{atom: plus, left: newVariable('x'), right: newNumber(3)}},
		},
		false,
	},
	{
		"sqrt(x)", &Expression{atom: sqrtFn, left: nil, right: newVariable('x')}, false,
	},
	{
		"ln(x)", &Expression{atom: logFn, left: nil, right: newVariable('x')}, false,
	},
	{
		"exp(x)", &Expression{atom: expFn, left: nil, right: newVariable('x')}, false,
	},
	{
		"sin(x)", &Expression{atom: sinFn, left: nil, right: newVariable('x')}, false,
	},
	{
		"cos(x)", &Expression{atom: cosFn, left: nil, right: newVariable('x')}, false,
	},
	{
		"abs(x)", &Expression{atom: absFn, left: nil, right: newVariable('x')}, false,
	},
	{
		"ln(x + y)", &Expression{atom: logFn, left: nil, right: &Expression{atom: plus, left: newVariable('x'), right: newVariable('y')}}, false,
	},
}

func Test_parseExpression(t *testing.T) {
	for _, tt := range expressions {
		got, err := Parse(tt.expr)
		if err != nil {
			_ = err.Error()
		}

		if (err != nil) != tt.wantErr {
			t.Fatalf("parseExpression(%s) error = %v, wantErr %v", tt.expr, err, tt.wantErr)
			return
		}

		if !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("parseExpression(%s) = %v, want %v", tt.expr, got, tt.want)
		}
	}
}
