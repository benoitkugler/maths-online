package expression

import (
	"reflect"
	"testing"
)

var expressions = [...]struct {
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
	{"3 + ln()", nil, true},
	{"2 , 5", nil, true},
	{"sgn( )", nil, true},
	{"sgn(-8)", &Expression{atom: sgnFn, right: &Expression{atom: minus, right: newNumber(8)}}, false},
	{"randInt(1.5, )", nil, true},
	{"randInt 1.5, )", nil, true},
	{"randInt(1.., )", nil, true},
	{"randInt(1.5, )", nil, true},
	{"randInt(0, 2.2)", nil, true},
	{"randInt(0, 1", nil, true},
	{"randInt(2 * 4, 1)", nil, true}, // not supported for now
	{"randInt(0, 1)", &Expression{atom: random{0, 1}}, false},
	{"randInt(2, 12)", &Expression{atom: random{2, 12}}, false},
	{"randInt(15, 12)", nil, true},
	{
		"2 + 3 * randInt(2, 12)", &Expression{atom: plus, left: newNumber(2), right: &Expression{atom: mult, left: newNumber(3), right: &Expression{atom: random{2, 12}}}}, false,
	},
	// space are optional
	{
		"(x−6)*(4*x−3)", &Expression{
			atom: mult,
			left: &Expression{atom: minus, left: newVariable('x'), right: newNumber(6)},
			right: &Expression{
				atom:  minus,
				left:  &Expression{atom: mult, left: newNumber(4), right: newVariable('x')},
				right: newNumber(3),
			},
		},
		false,
	},
}

func Test_parseExpression(t *testing.T) {
	for _, tt := range expressions {
		got, _, err := Parse(tt.expr)
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

func TestVarMap(t *testing.T) {
	for _, tt := range []struct {
		want VarMap
		expr string
	}{
		{VarMap{}, "2 + 3"},
		{VarMap{0: 'a', 4: 'b'}, "a + b"},
		{VarMap{0: 'a', 4: 'b', 9: 'a'}, "a + b * (a + 2)"},
	} {
		_, vm, err := Parse(tt.expr)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(vm, tt.want) {
			t.Fatalf("for expr %s, expected %v, got %v", tt.expr, tt.want, vm)
		}
	}
}

func TestVarMap_Positions(t *testing.T) {
	tests := []struct {
		expr string
		args []Variable
		want []int
	}{
		{
			"a + (x * b)", []Variable{}, nil,
		},
		{
			"a + (x * b)", []Variable{'a', 'b'}, []int{0, 9},
		},
		{
			"a + (x * b - b)", []Variable{'a', 'b'}, []int{0, 9, 13},
		},
	}
	for _, tt := range tests {
		rv := make(RandomParameters)
		for _, v := range tt.args {
			rv[v] = nil
		}

		_, vm, err := Parse(tt.expr)
		if err != nil {
			t.Fatal(err)
		}
		if got := vm.Positions(rv); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("VarMap.Positions() = %v, want %v", got, tt.want)
		}
	}
}
