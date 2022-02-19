package expression

import (
	"reflect"
	"testing"
)

func Test_parseExpression(t *testing.T) {
	tests := []struct {
		args    string
		want    *node
		wantErr bool
	}{
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
			" e ", &node{atom: numberE}, false,
		},
		{
			" \u03C0 ", &node{atom: numberPi}, false,
		},
		{
			" (e) ", &node{atom: numberE}, false,
		},
		{
			" ((e)) ", &node{atom: numberE}, false,
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
			" x + 3 ", &node{atom: plus, left: newVariable('x'), right: newNumber(3)}, false,
		},
		{
			" x - 3 ", &node{atom: minus, left: newVariable('x'), right: newNumber(3)}, false,
		},
		{
			" x ^ 3 ", &node{atom: pow, left: newVariable('x'), right: newNumber(3)}, false,
		},
		{
			" x / 3 ", &node{atom: div, left: newVariable('x'), right: newNumber(3)}, false,
		},
		{
			" x * 3 ", &node{atom: mult, left: newVariable('x'), right: newNumber(3)}, false,
		},
		{
			"(x + 3 )", &node{atom: plus, left: newVariable('x'), right: newNumber(3)}, false,
		},
		{
			"3 - (x + 3 )", &node{atom: minus, left: newNumber(3), right: &node{atom: plus, left: newVariable('x'), right: newNumber(3)}}, false,
		},
		{
			"(1 + y) / (3 - (x + 3 ))",

			&node{
				atom:  div,
				left:  &node{atom: plus, left: newNumber(1), right: newVariable('y')},
				right: &node{atom: minus, left: newNumber(3), right: &node{atom: plus, left: newVariable('x'), right: newNumber(3)}},
			},
			false,
		},
		{
			"ln(x)", &node{atom: log, left: nil, right: newVariable('x')}, false,
		},
		{
			"exp(x)", &node{atom: exp, left: nil, right: newVariable('x')}, false,
		},
		{
			"sin(x)", &node{atom: sin, left: nil, right: newVariable('x')}, false,
		},
		{
			"cos(x)", &node{atom: cos, left: nil, right: newVariable('x')}, false,
		},
		{
			"abs(x)", &node{atom: abs, left: nil, right: newVariable('x')}, false,
		},
		{
			"ln(x + y)", &node{atom: log, left: nil, right: &node{atom: plus, left: newVariable('x'), right: newVariable('y')}}, false,
		},
	}
	for _, tt := range tests {
		got, err := parseExpression(tt.args)
		if err != nil {
			_ = err.Error()
		}

		if (err != nil) != tt.wantErr {
			t.Errorf("%s : parseExpression() error = %v, wantErr %v", tt.args, err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("parseExpression() = %v, want %v", got, tt.want)
		}
	}
}
