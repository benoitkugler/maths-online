package expression

import (
	"fmt"
	"math"
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
			left: NewNb(1),
			right: &Expression{
				atom:  mult,
				left:  &Expression{atom: plus, left: newVarExpr('x'), right: newVarExpr('z')},
				right: &Expression{atom: plus, left: newVarExpr('x'), right: NewNb(1)},
			},
		}, false,
	},
	{
		"1 + 2 * 3 * 4", &Expression{
			atom: plus,
			left: NewNb(1),
			right: &Expression{
				atom:  mult,
				left:  &Expression{atom: mult, left: NewNb(2), right: NewNb(3)},
				right: NewNb(4),
			},
		}, false,
	},
	{
		"1 + (x + z) * (x + 1) * z ", &Expression{
			atom: plus,
			left: NewNb(1),
			right: &Expression{
				atom: mult,
				left: &Expression{
					atom:  mult,
					left:  &Expression{atom: plus, left: newVarExpr('x'), right: newVarExpr('z')},
					right: &Expression{atom: plus, left: newVarExpr('x'), right: NewNb(1)},
				},
				right: newVarExpr('z'),
			},
		}, false,
	},
	{
		"x", newVarExpr('x'), false,
	},
	{
		"t * q", &Expression{atom: mult, left: newVarExpr('t'), right: newVarExpr('q')}, false,
	},
	{
		" x ", newVarExpr('x'), false,
	},
	{
		" x_ ", newVarExpr('x'), false,
	},
	{
		" x_AB ", &Expression{atom: Variable{Name: 'x', Indice: "AB"}}, false,
	},
	{
		" 3 ", NewNb(3), false,
	},
	{
		" 3.14 ", NewNb(3.14), false,
	},
	{
		" 3,14 ", NewNb(3.14), false,
	},
	{
		" e ", &Expression{atom: eConstant}, false,
	},
	// variable with indice
	{
		"(x_a  +x_b) /2", &Expression{atom: div, left: &Expression{
			atom:  plus,
			left:  &Expression{atom: Variable{Name: 'x', Indice: "a"}},
			right: &Expression{atom: Variable{Name: 'x', Indice: "b"}},
		}, right: NewNb(2)}, false,
	},
	// custom variables
	{
		" \uE000 + 2", &Expression{atom: plus, left: newVarExpr('\uE000'), right: NewNb(2)}, false,
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
		"(;)", nil, true,
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
		" log 3 ", nil, true,
	},
	{
		" log( /2)", nil, true,
	},
	{
		" x + 3 ", &Expression{atom: plus, left: newVarExpr('x'), right: NewNb(3)}, false,
	},
	{
		" x - 3 ", &Expression{atom: minus, left: newVarExpr('x'), right: NewNb(3)}, false,
	},
	{
		" x ^ 3 ", &Expression{atom: pow, left: newVarExpr('x'), right: NewNb(3)}, false,
	},
	// ^ is not associative !
	{
		" a^b^c ", &Expression{
			atom: pow,
			left: newVarExpr('a'),
			right: &Expression{
				atom:  pow,
				left:  newVarExpr('b'),
				right: newVarExpr('c'),
			},
		}, false,
	},
	{
		" (a^b)^c ", &Expression{
			atom: pow,
			left: &Expression{
				atom:  pow,
				left:  newVarExpr('a'),
				right: newVarExpr('b'),
			},
			right: newVarExpr('c'),
		}, false,
	},
	{
		" a^b  * c ", &Expression{
			atom: mult,
			left: &Expression{
				atom:  pow,
				left:  newVarExpr('a'),
				right: newVarExpr('b'),
			},
			right: newVarExpr('c'),
		}, false,
	},
	{
		" x / 3 ", &Expression{atom: div, left: newVarExpr('x'), right: NewNb(3)}, false,
	},
	{
		" x * 3 ", &Expression{atom: mult, left: newVarExpr('x'), right: NewNb(3)}, false,
	},
	{
		"(x + 3 )", &Expression{atom: plus, left: newVarExpr('x'), right: NewNb(3)}, false,
	},
	// implicit multiplication référence
	{
		"(x + 3)*(x+4)", &Expression{
			atom:  mult,
			left:  &Expression{atom: plus, left: newVarExpr('x'), right: NewNb(3)},
			right: &Expression{atom: plus, left: newVarExpr('x'), right: NewNb(4)},
		}, false,
	},
	// implicit multiplication
	{
		"xln(x)", &Expression{
			atom:  mult,
			left:  newVarExpr('x'),
			right: &Expression{atom: logFn, right: newVarExpr('x')},
		}, false,
	},
	{
		"(x + 3)(x+4)", &Expression{
			atom:  mult,
			left:  &Expression{atom: plus, left: newVarExpr('x'), right: NewNb(3)},
			right: &Expression{atom: plus, left: newVarExpr('x'), right: NewNb(4)},
		}, false,
	},
	{
		" (1+ 2 ) (x + 3) ", &Expression{
			atom:  mult,
			left:  &Expression{atom: plus, left: NewNb(1), right: NewNb(2)},
			right: &Expression{atom: plus, left: newVarExpr('x'), right: NewNb(3)},
		}, false,
	},
	{
		"(x−6)(4x−3)", &Expression{
			atom: mult,
			left: &Expression{atom: minus, left: newVarExpr('x'), right: NewNb(6)},
			right: &Expression{
				atom:  minus,
				left:  &Expression{atom: mult, left: NewNb(4), right: newVarExpr('x')},
				right: NewNb(3),
			},
		}, false,
	},
	{
		"24x^2 - 27x + 18", &Expression{
			atom: plus,
			left: &Expression{
				atom: minus,
				left: &Expression{
					atom: mult,
					left: NewNb(24),
					right: &Expression{
						atom:  pow,
						left:  newVarExpr('x'),
						right: NewNb(2),
					},
				},
				right: &Expression{
					atom:  mult,
					left:  NewNb(27),
					right: newVarExpr('x'),
				},
			},
			right: NewNb(18),
		}, false,
	},
	{
		"x4", nil, true, // invalid implicit multiplication
	},
	{
		"1 + 2 * 3", &Expression{
			atom:  plus,
			left:  NewNb(1),
			right: &Expression{atom: mult, left: NewNb(2), right: NewNb(3)},
		}, false,
	},
	{
		"1 + 2 * 3 ^ 4", &Expression{
			atom: plus,
			left: NewNb(1),
			right: &Expression{
				atom: mult,
				left: NewNb(2),
				right: &Expression{
					atom:  pow,
					left:  NewNb(3),
					right: NewNb(4),
				},
			},
		}, false,
	},
	{
		"1 + (x + z) * (x + 1) * z ", &Expression{
			atom: plus,
			left: NewNb(1),
			right: &Expression{
				atom: mult,
				left: &Expression{
					atom:  mult,
					left:  &Expression{atom: plus, left: newVarExpr('x'), right: newVarExpr('z')},
					right: &Expression{atom: plus, left: newVarExpr('x'), right: NewNb(1)},
				},
				right: newVarExpr('z'),
			},
		}, false,
	},
	{
		"1 + 2 * (x + 1)", &Expression{
			atom:  plus,
			left:  NewNb(1),
			right: &Expression{atom: mult, left: NewNb(2), right: &Expression{atom: plus, left: newVarExpr('x'), right: NewNb(1)}},
		}, false,
	},
	{
		"3 - (x + 3 )", &Expression{atom: minus, left: NewNb(3), right: &Expression{atom: plus, left: newVarExpr('x'), right: NewNb(3)}}, false,
	},
	{
		"(1 + y) / (3 - (x + 3 ))",

		&Expression{
			atom:  div,
			left:  &Expression{atom: plus, left: NewNb(1), right: newVarExpr('y')},
			right: &Expression{atom: minus, left: NewNb(3), right: &Expression{atom: plus, left: newVarExpr('x'), right: NewNb(3)}},
		},
		false,
	},
	{
		"sqrt(x)", &Expression{atom: sqrtFn, left: nil, right: newVarExpr('x')}, false,
	},
	{
		"ln(x)", &Expression{atom: logFn, left: nil, right: newVarExpr('x')}, false,
	},
	{
		"exp(x)", &Expression{atom: expFn, left: nil, right: newVarExpr('x')}, false,
	},
	{
		"sin(x)", &Expression{atom: sinFn, left: nil, right: newVarExpr('x')}, false,
	},
	{
		"cos(x)", &Expression{atom: cosFn, left: nil, right: newVarExpr('x')}, false,
	},
	{
		"tan(x)", &Expression{atom: tanFn, left: nil, right: newVarExpr('x')}, false,
	},
	{
		"asin(x)", &Expression{atom: asinFn, left: nil, right: newVarExpr('x')}, false,
	},
	{
		"acos(x)", &Expression{atom: acosFn, left: nil, right: newVarExpr('x')}, false,
	},
	{
		"atan(x)", &Expression{atom: atanFn, left: nil, right: newVarExpr('x')}, false,
	},
	{
		"abs(x)", &Expression{atom: absFn, left: nil, right: newVarExpr('x')}, false,
	},
	{
		"ln(x + y)", &Expression{atom: logFn, left: nil, right: &Expression{atom: plus, left: newVarExpr('x'), right: newVarExpr('y')}}, false,
	},
	{"3 + ln()", nil, true},
	{"2 , 5", nil, true},
	{"isPrime( )", nil, true},
	{"sgn( )", nil, true},
	{"sgn(-8)", &Expression{atom: sgnFn, right: &Expression{atom: minus, right: NewNb(8)}}, false},
	{"isZero( )", nil, true},
	{"isZero(-8)", &Expression{atom: isZeroFn, right: &Expression{atom: minus, right: NewNb(8)}}, false},
	{"%", nil, true},
	{"8 % 2", &Expression{atom: mod, left: NewNb(8), right: NewNb(2)}, false},
	{"//", nil, true},
	{"8 // 2", &Expression{atom: rem, left: NewNb(8), right: NewNb(2)}, false},
	{"randInt(-a, )", nil, true},
	{"randInt(1.5; )", nil, true},
	{"randInt 1.5; )", nil, true},
	{"randInt(1..; )", nil, true},
	{"randInt(1.5; )", nil, true},
	{"randInt(0; 2.2)", nil, true},
	{"randInt(0; 1", nil, true},
	{"randInt(2 * 4; 1)", nil, true}, // not supported for now
	{"randPrime(-1; 12)", nil, true},
	{"randPrime(4; 4)", nil, true},
	{"randInt(0; 1)", &Expression{atom: specialFunctionA{kind: randInt, args: []*Expression{NewNb(0), NewNb(1)}}}, false},
	{"randInt(2; 12)", &Expression{atom: specialFunctionA{kind: randInt, args: []*Expression{NewNb(2), NewNb(12)}}}, false},
	{"randInt(-1; 4)", &Expression{atom: specialFunctionA{kind: randInt, args: []*Expression{NewNb(-1), NewNb(4)}}}, false},
	{"randInt(a; b)", &Expression{atom: specialFunctionA{kind: randInt, args: []*Expression{newVarExpr('a'), newVarExpr('b')}}}, false},
	{"randPrime(0; 2)", &Expression{atom: specialFunctionA{kind: randPrime, args: []*Expression{NewNb(0), NewNb(2)}}}, false},
	{"randPrime(2; 12)", &Expression{atom: specialFunctionA{kind: randPrime, args: []*Expression{NewNb(2), NewNb(12)}}}, false},
	{"randChoice(1.2;4 ; -3)", &Expression{atom: specialFunctionA{kind: randChoice, args: []*Expression{NewNb(1.2), NewNb(4), NewNb(-3)}}}, false},
	{"randDecDen( )", &Expression{atom: specialFunctionA{kind: randDenominator, args: nil}}, false},
	{"randInt(15; 12)", nil, true},
	{"randChoice( )", nil, true},
	{"randChoice(2;", nil, true},
	{"randLetter(A; x_A; b;  B; B)", &Expression{atom: randVariable{
		NewVar('A'), Variable{Name: 'x', Indice: "A"}, NewVar('b'), NewVar('B'), NewVar('B'),
	}}, false},
	{"randLetter( )", nil, true},
	{"randLetter)", nil, true},
	{"randLetter(0.2 )", nil, true},
	{"randLetter(x;", nil, true},
	{"randLetter(x,y)", nil, true},
	{
		"2 + 3 * randInt(2; 12)", &Expression{
			atom:  plus,
			left:  NewNb(2),
			right: &Expression{atom: mult, left: NewNb(3), right: &Expression{atom: specialFunctionA{kind: randInt, args: []*Expression{NewNb(2), NewNb(12)}}}},
		}, false,
	},
	{
		"isPrime(2 * x)", &Expression{atom: isPrimeFn, left: nil, right: &Expression{atom: mult, left: NewNb(2), right: newVarExpr('x')}}, false,
	},
	// round
	{"round(x,y)", nil, true},
	{"round x", nil, true},
	{"round(x)", nil, true},
	{"round(x;2.2)", nil, true},
	{"round(x;2.2.)", nil, true},
	{"round(x;)", nil, true},
	{"round(x;2", nil, true},
	{"round(x;2)", &Expression{atom: roundFn{2}, right: newVarExpr('x')}, false},
	{"round(x + randInt(1;5);2)", &Expression{
		atom: roundFn{2},
		right: &Expression{
			atom: plus,
			left: newVarExpr('x'),
			right: &Expression{atom: specialFunctionA{
				kind: randInt,
				args: []*Expression{NewNb(1), NewNb(5)},
			}},
		},
	}, false},
	// space are optional
	{
		"(x−6)*(4*x−3)", &Expression{
			atom: mult,
			left: &Expression{atom: minus, left: newVarExpr('x'), right: NewNb(6)},
			right: &Expression{
				atom:  minus,
				left:  &Expression{atom: mult, left: NewNb(4), right: newVarExpr('x')},
				right: NewNb(3),
			},
		},
		false,
	},
	// infinity
	{
		"Inf", &Expression{atom: Number(math.Inf(1)), left: nil, right: nil}, false,
	},
	{
		"inf", &Expression{atom: Number(math.Inf(1)), left: nil, right: nil}, false,
	},
	{"max()", nil, true},
	{"min()", nil, true},
}

func Test_parseExpression(t *testing.T) {
	for _, tt := range expressions {
		got, err := Parse(tt.expr)
		if err != nil {
			_ = err.Error()
			_ = err.(ErrInvalidExpr).Portion()
		}

		if (err != nil) != tt.wantErr {
			t.Fatalf("parseExpression(%s) error = %v, wantErr %v", tt.expr, err, tt.wantErr)
			return
		}

		if !reflect.DeepEqual(got, tt.want) {
			fmt.Printf("%#v \n%#v\n", got, tt.want)
			t.Fatalf("parseExpression(%s) = %v, want %v", tt.expr, got, tt.want)
		}
	}
}

func TestVarMap(t *testing.T) {
	for _, tt := range []struct {
		want varMap
		expr string
	}{
		{varMap{}, "2 + 3"},
		{varMap{0: NewVar('a'), 4: NewVar('b')}, "a + b"},
		{varMap{0: NewVar('a'), 4: NewVar('b'), 9: NewVar('a')}, "a + b * (a + 2)"},
	} {
		_, vm, err := parseBytes([]byte(tt.expr))
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
			"a + (x * b)", []Variable{NewVar('a'), NewVar('b')}, []int{0, 9},
		},
		{
			"a + (x * b - b)", []Variable{NewVar('a'), NewVar('b')}, []int{0, 9, 13},
		},
	}
	for _, tt := range tests {
		rv := make(RandomParameters)
		for _, v := range tt.args {
			rv[v] = nil
		}

		_, vm, err := parseBytes([]byte(tt.expr))
		if err != nil {
			t.Fatal(err)
		}
		if got := vm.Positions(rv); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("VarMap.Positions() = %v, want %v", got, tt.want)
		}
	}
}
