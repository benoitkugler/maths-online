package expression

import (
	"reflect"
	"testing"
)

func Test_parseExpression(t *testing.T) {
	tests := []struct {
		args    string
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
			" x ", newVariable('x'), false,
		},
		{
			" 3 ", newNumber(3), false,
		},
		{
			" 3.14 ", newNumber(3.14), false,
		},
		{
			" e ", &Expression{atom: numberE}, false,
		},
		{
			" \u03C0 ", &Expression{atom: numberPi}, false,
		},
		{
			" (e) ", &Expression{atom: numberE}, false,
		},
		{
			" ((e)) ", &Expression{atom: numberE}, false,
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
			"ln(x)", &Expression{atom: log, left: nil, right: newVariable('x')}, false,
		},
		{
			"exp(x)", &Expression{atom: exp, left: nil, right: newVariable('x')}, false,
		},
		{
			"sin(x)", &Expression{atom: sin, left: nil, right: newVariable('x')}, false,
		},
		{
			"cos(x)", &Expression{atom: cos, left: nil, right: newVariable('x')}, false,
		},
		{
			"abs(x)", &Expression{atom: abs, left: nil, right: newVariable('x')}, false,
		},
		{
			"ln(x + y)", &Expression{atom: log, left: nil, right: &Expression{atom: plus, left: newVariable('x'), right: newVariable('y')}}, false,
		},
	}
	for _, tt := range tests {
		got, err := Parse(tt.args)
		if err != nil {
			_ = err.Error()
		}

		if (err != nil) != tt.wantErr {
			t.Fatalf("parseExpression(%s) error = %v, wantErr %v", tt.args, err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("parseExpression(%s) = %v, want %v", tt.args, got, tt.want)
		}
	}
}
