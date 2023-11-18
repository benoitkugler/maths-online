package expression

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"

	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

var expressions = [...]struct {
	expr    string
	want    *Expr
	wantErr bool
}{
	{
		"1 + (x + z) * (x + 1)", &Expr{
			atom: plus,
			left: NewNb(1),
			right: &Expr{
				atom:  mult,
				left:  &Expr{atom: plus, left: newVarExpr('x'), right: newVarExpr('z')},
				right: &Expr{atom: plus, left: newVarExpr('x'), right: NewNb(1)},
			},
		}, false,
	},
	{
		"1 + 2 * 3 * 4", &Expr{
			atom: plus,
			left: NewNb(1),
			right: &Expr{
				atom:  mult,
				left:  &Expr{atom: mult, left: NewNb(2), right: NewNb(3)},
				right: NewNb(4),
			},
		}, false,
	},
	{
		"1 + (x + z) * (x + 1) * z ", &Expr{
			atom: plus,
			left: NewNb(1),
			right: &Expr{
				atom: mult,
				left: &Expr{
					atom:  mult,
					left:  &Expr{atom: plus, left: newVarExpr('x'), right: newVarExpr('z')},
					right: &Expr{atom: plus, left: newVarExpr('x'), right: NewNb(1)},
				},
				right: newVarExpr('z'),
			},
		}, false,
	},
	{
		"x", newVarExpr('x'), false,
	},
	{
		"t * q", &Expr{atom: mult, left: newVarExpr('t'), right: newVarExpr('q')}, false,
	},
	{
		" x ", newVarExpr('x'), false,
	},
	{
		" x_ ", newVarExpr('x'), false,
	},
	{
		" x_AB ", &Expr{atom: Variable{Name: 'x', Indice: "AB"}}, false,
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
		" e ", &Expr{atom: eConstant}, false,
	},
	// variable with indice
	{
		"(x_a  +x_b) /2", &Expr{atom: div, left: &Expr{
			atom:  plus,
			left:  &Expr{atom: Variable{Name: 'x', Indice: "a"}},
			right: &Expr{atom: Variable{Name: 'x', Indice: "b"}},
		}, right: NewNb(2)}, false,
	},
	// custom variables
	{
		" \uE000 + 2", &Expr{atom: plus, left: newVarExpr('\uE000'), right: NewNb(2)}, false,
	},
	{
		" \u03C0 ", &Expr{atom: piConstant}, false,
	},
	{
		" (e) ", &Expr{atom: eConstant}, false,
	},
	{
		" ((e)) ", &Expr{atom: eConstant}, false,
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
		" x + 3 ", &Expr{atom: plus, left: newVarExpr('x'), right: NewNb(3)}, false,
	},
	{
		" x - 3 ", &Expr{atom: minus, left: newVarExpr('x'), right: NewNb(3)}, false,
	},
	{
		" x ^ 3 ", &Expr{atom: pow, left: newVarExpr('x'), right: NewNb(3)}, false,
	},
	{
		" x ^ (-3) ", &Expr{atom: pow, left: newVarExpr('x'), right: NewNb(-3)}, false,
	},
	// extended support for negative power
	{
		" x ^- ", nil, true,
	},
	{
		" x ^-)", nil, true,
	},
	{
		" x ^-(a+)", nil, true,
	},
	{
		" x ^-3 ", &Expr{atom: pow, left: newVarExpr('x'), right: NewNb(-3)}, false,
	},
	{
		" x ^-3 + 1",
		&Expr{atom: plus, left: &Expr{atom: pow, left: newVarExpr('x'), right: NewNb(-3)}, right: newNb(1)}, false,
	},
	{
		" x ^-(x+y) ", &Expr{atom: pow, left: newVarExpr('x'), right: &Expr{atom: minus, right: &Expr{atom: plus, left: newVarExpr('x'), right: newVarExpr('y')}}}, false,
	},
	// { // = (a^n) / c
	// 	" a ^ -2n ", &Expr{atom: div, left: &Expr{atom: pow, left: newVarExpr('a'), right: newVarExpr('n')}, right: newVarExpr('c')}, false,
	// },
	// operator precedence
	{ // = (a^n) / c
		" a ^ n / c ", &Expr{atom: div, left: &Expr{atom: pow, left: newVarExpr('a'), right: newVarExpr('n')}, right: newVarExpr('c')}, false,
	},
	{ // = c * (a^n)
		" c * a ^ n ", &Expr{atom: mult, left: newVarExpr('c'), right: &Expr{atom: pow, left: newVarExpr('a'), right: newVarExpr('n')}}, false,
	},
	{ // = a + (n * c)
		" a + n * c ", &Expr{atom: plus, left: newVarExpr('a'), right: &Expr{atom: mult, left: newVarExpr('n'), right: newVarExpr('c')}}, false,
	},
	{ // = (a * n) + c
		" a * n + c ", &Expr{atom: plus, left: &Expr{atom: mult, left: newVarExpr('a'), right: newVarExpr('n')}, right: newVarExpr('c')}, false,
	},
	// ^ is not associative !
	{
		" a^b^c ", &Expr{
			atom: pow,
			left: newVarExpr('a'),
			right: &Expr{
				atom:  pow,
				left:  newVarExpr('b'),
				right: newVarExpr('c'),
			},
		}, false,
	},
	{
		" (a^b)^c ", &Expr{
			atom: pow,
			left: &Expr{
				atom:  pow,
				left:  newVarExpr('a'),
				right: newVarExpr('b'),
			},
			right: newVarExpr('c'),
		}, false,
	},
	{
		" a^b  * c ", &Expr{
			atom: mult,
			left: &Expr{
				atom:  pow,
				left:  newVarExpr('a'),
				right: newVarExpr('b'),
			},
			right: newVarExpr('c'),
		}, false,
	},
	{
		" x / 3 ", &Expr{atom: div, left: newVarExpr('x'), right: NewNb(3)}, false,
	},
	{
		" x * 3 ", &Expr{atom: mult, left: newVarExpr('x'), right: NewNb(3)}, false,
	},
	{
		"(x + 3 )", &Expr{atom: plus, left: newVarExpr('x'), right: NewNb(3)}, false,
	},
	// implicit multiplication référence
	{
		"(x + 3)*(x+4)", &Expr{
			atom:  mult,
			left:  &Expr{atom: plus, left: newVarExpr('x'), right: NewNb(3)},
			right: &Expr{atom: plus, left: newVarExpr('x'), right: NewNb(4)},
		}, false,
	},
	// implicit multiplication
	{
		"n!(n)", &Expr{
			atom:  mult,
			left:  &Expr{atom: factorial, left: newVarExpr('n')},
			right: newVarExpr('n'),
		}, false,
	},
	{
		"n!n", &Expr{
			atom:  mult,
			left:  &Expr{atom: factorial, left: newVarExpr('n')},
			right: newVarExpr('n'),
		}, false,
	},
	{
		"k! k!", &Expr{
			atom:  mult,
			left:  &Expr{atom: factorial, left: newVarExpr('k')},
			right: &Expr{atom: factorial, left: newVarExpr('k')},
		}, false,
	},
	{
		"xln(x)", &Expr{
			atom:  mult,
			left:  newVarExpr('x'),
			right: &Expr{atom: logFn, right: newVarExpr('x')},
		}, false,
	},
	{
		"(x + 3)(x+4)", &Expr{
			atom:  mult,
			left:  &Expr{atom: plus, left: newVarExpr('x'), right: NewNb(3)},
			right: &Expr{atom: plus, left: newVarExpr('x'), right: NewNb(4)},
		}, false,
	},
	{
		" (1+ 2 ) (x + 3) ", &Expr{
			atom:  mult,
			left:  &Expr{atom: plus, left: NewNb(1), right: NewNb(2)},
			right: &Expr{atom: plus, left: newVarExpr('x'), right: NewNb(3)},
		}, false,
	},
	{
		"(x−6)(4x−3)", &Expr{
			atom: mult,
			left: &Expr{atom: minus, left: newVarExpr('x'), right: NewNb(6)},
			right: &Expr{
				atom:  minus,
				left:  &Expr{atom: mult, left: NewNb(4), right: newVarExpr('x')},
				right: NewNb(3),
			},
		}, false,
	},
	// exponents
	{
		"x¹ + x²", &Expr{
			atom: plus, left: &Expr{
				atom:  pow,
				left:  newVarExpr('x'),
				right: newNb(1),
			},
			right: &Expr{
				atom:  pow,
				left:  newVarExpr('x'),
				right: newNb(2),
			},
		}, false,
	},
	{
		"24x^2 - 27x + 18", &Expr{
			atom: plus,
			left: &Expr{
				atom: minus,
				left: &Expr{
					atom: mult,
					left: NewNb(24),
					right: &Expr{
						atom:  pow,
						left:  newVarExpr('x'),
						right: NewNb(2),
					},
				},
				right: &Expr{
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
		"1 + 2 * 3", &Expr{
			atom:  plus,
			left:  NewNb(1),
			right: &Expr{atom: mult, left: NewNb(2), right: NewNb(3)},
		}, false,
	},
	{
		"1 + 2 * 3 ^ 4", &Expr{
			atom: plus,
			left: NewNb(1),
			right: &Expr{
				atom: mult,
				left: NewNb(2),
				right: &Expr{
					atom:  pow,
					left:  NewNb(3),
					right: NewNb(4),
				},
			},
		}, false,
	},
	{
		"1 + (x + z) * (x + 1) * z ", &Expr{
			atom: plus,
			left: NewNb(1),
			right: &Expr{
				atom: mult,
				left: &Expr{
					atom:  mult,
					left:  &Expr{atom: plus, left: newVarExpr('x'), right: newVarExpr('z')},
					right: &Expr{atom: plus, left: newVarExpr('x'), right: NewNb(1)},
				},
				right: newVarExpr('z'),
			},
		}, false,
	},
	{
		"1 + 2 * (x + 1)", &Expr{
			atom:  plus,
			left:  NewNb(1),
			right: &Expr{atom: mult, left: NewNb(2), right: &Expr{atom: plus, left: newVarExpr('x'), right: NewNb(1)}},
		}, false,
	},
	{
		"3 - (x + 3 )", &Expr{atom: minus, left: NewNb(3), right: &Expr{atom: plus, left: newVarExpr('x'), right: NewNb(3)}}, false,
	},
	{
		"(1 + y) / (3 - (x + 3 ))",

		&Expr{
			atom:  div,
			left:  &Expr{atom: plus, left: NewNb(1), right: newVarExpr('y')},
			right: &Expr{atom: minus, left: NewNb(3), right: &Expr{atom: plus, left: newVarExpr('x'), right: NewNb(3)}},
		},
		false,
	},
	{
		"sqrt(x)", &Expr{atom: sqrtFn, left: nil, right: newVarExpr('x')}, false,
	},
	{
		"ln(x)", &Expr{atom: logFn, left: nil, right: newVarExpr('x')}, false,
	},
	{
		"exp(x)", &Expr{atom: expFn, left: nil, right: newVarExpr('x')}, false,
	},
	{
		"sin(x)", &Expr{atom: sinFn, left: nil, right: newVarExpr('x')}, false,
	},
	{
		"cos(x)", &Expr{atom: cosFn, left: nil, right: newVarExpr('x')}, false,
	},
	{
		"tan(x)", &Expr{atom: tanFn, left: nil, right: newVarExpr('x')}, false,
	},
	{
		"asin(x)", &Expr{atom: asinFn, left: nil, right: newVarExpr('x')}, false,
	},
	{
		"acos(x)", &Expr{atom: acosFn, left: nil, right: newVarExpr('x')}, false,
	},
	{
		"atan(x)", &Expr{atom: atanFn, left: nil, right: newVarExpr('x')}, false,
	},
	{
		"abs(x)", &Expr{atom: absFn, left: nil, right: newVarExpr('x')}, false,
	},
	{
		"floor(x)", &Expr{atom: floorFn, left: nil, right: newVarExpr('x')}, false,
	},
	{
		"ln(x + y)", &Expr{atom: logFn, left: nil, right: &Expr{atom: plus, left: newVarExpr('x'), right: newVarExpr('y')}}, false,
	},
	{
		"forceDecimal(x)", &Expr{atom: forceDecimalFn, left: nil, right: newVarExpr('x')}, false,
	},
	{"3 + ln()", nil, true},
	{"2 , 5", nil, true},
	{"isPrime( )", nil, true},
	{"sgn( )", nil, true},
	{"sgn(-8)", &Expr{atom: sgnFn, right: &Expr{atom: minus, right: NewNb(8)}}, false},
	{"isZero( )", nil, true},
	{"%", nil, true},
	{"8 % 2", &Expr{atom: mod, left: NewNb(8), right: NewNb(2)}, false},
	{"//", nil, true},
	{"8 // 2", &Expr{atom: rem, left: NewNb(8), right: NewNb(2)}, false},
	// factorial
	{"n - !", nil, true},
	{"n!", &Expr{atom: factorial, left: newVarExpr('n'), right: nil}, false},
	{"(2 + n)!", &Expr{atom: factorial, left: &Expr{atom: plus, left: NewNb(2), right: newVarExpr('n')}, right: nil}, false},

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
	{"randInt(0; 1)", &Expr{atom: specialFunction{kind: randInt, args: []*Expr{NewNb(0), NewNb(1)}}}, false},
	{"randInt(2; 12)", &Expr{atom: specialFunction{kind: randInt, args: []*Expr{NewNb(2), NewNb(12)}}}, false},
	{"randInt(-1; 4)", &Expr{atom: specialFunction{kind: randInt, args: []*Expr{NewNb(-1), NewNb(4)}}}, false},
	{"randInt(a; b)", &Expr{atom: specialFunction{kind: randInt, args: []*Expr{newVarExpr('a'), newVarExpr('b')}}}, false},
	{"randPrime(0; 2)", &Expr{atom: specialFunction{kind: randPrime, args: []*Expr{NewNb(0), NewNb(2)}}}, false},
	{"randPrime(2; 12)", &Expr{atom: specialFunction{kind: randPrime, args: []*Expr{NewNb(2), NewNb(12)}}}, false},
	{"randChoice(1.2;4 ; -3)", &Expr{atom: specialFunction{kind: randChoice, args: []*Expr{NewNb(1.2), NewNb(4), NewNb(-3)}}}, false},
	{"randDecDen( )", &Expr{atom: specialFunction{kind: randDenominator, args: nil}}, false},
	{"randInt(15; 12)", nil, true},
	{"randChoice( )", nil, true},
	{"randChoice(2;", nil, true},
	// randChoice with arbitrary expression
	{"randChoice(A; x_A; b;  B; B)", &Expr{atom: specialFunction{
		kind: randChoice,
		args: []*Expr{
			newVarExpr('A'), NewVarExpr(Variable{Name: 'x', Indice: "A"}), newVarExpr('b'), newVarExpr('B'), newVarExpr('B'),
		},
	}}, false},
	{"randChoice( )", nil, true},
	{"randChoice)", nil, true},
	{"randChoice(x;", nil, true},
	{"randChoice(x,y)", nil, true},
	// choiceFrom
	{"choiceFrom(A; x_A; b;  B; B; sin(3))", &Expr{atom: specialFunction{
		kind: choiceFrom,
		args: []*Expr{
			newVarExpr('A'),
			NewVarExpr(Variable{Name: 'x', Indice: "A"}),
			newVarExpr('b'), newVarExpr('B'), newVarExpr('B'),
			{atom: sinFn, right: newNb(3)},
		},
	}}, false},
	{"choiceFrom( )", nil, true},
	{"choiceFrom((", nil, true},
	{"choiceFrom( 0.2 )", nil, true},
	{"choiceFrom( x;", nil, true},
	{"choiceFrom( x,y )", nil, true},
	{"choiceFrom(x;y); )", nil, true},
	{
		"2 + 3 * randInt(2; 12)", &Expr{
			atom:  plus,
			left:  NewNb(2),
			right: &Expr{atom: mult, left: NewNb(3), right: &Expr{atom: specialFunction{kind: randInt, args: []*Expr{NewNb(2), NewNb(12)}}}},
		}, false,
	},
	{
		"isPrime(2 * x)", &Expr{atom: isPrimeFn, left: nil, right: &Expr{atom: mult, left: NewNb(2), right: newVarExpr('x')}}, false,
	},
	// round
	{"round(x,y)", nil, true},
	{"round x", nil, true},
	{"round(x)", nil, true},
	{"round(x;2.2)", nil, true},
	{"round(x;2.2.)", nil, true},
	{"round(x;)", nil, true},
	{"round(x;2", nil, true},
	{"round(x;2)", &Expr{atom: roundFunc{2}, right: newVarExpr('x')}, false},
	{"round(x + randInt(1;5);2)", &Expr{
		atom: roundFunc{2},
		right: &Expr{
			atom: plus,
			left: newVarExpr('x'),
			right: &Expr{atom: specialFunction{
				kind: randInt,
				args: []*Expr{NewNb(1), NewNb(5)},
			}},
		},
	}, false},
	// space are optional
	{
		"(x−6)*(4*x−3)", &Expr{
			atom: mult,
			left: &Expr{atom: minus, left: newVarExpr('x'), right: NewNb(6)},
			right: &Expr{
				atom:  minus,
				left:  &Expr{atom: mult, left: NewNb(4), right: newVarExpr('x')},
				right: NewNb(3),
			},
		},
		false,
	},
	// infinity
	{
		"Inf", &Expr{atom: Number(math.Inf(1)), left: nil, right: nil}, false,
	},
	{
		"inf", &Expr{atom: Number(math.Inf(1)), left: nil, right: nil}, false,
	},
	{"max()", nil, true},
	{"min()", nil, true},
	{"sum()", nil, true},
	{"sum(1;2;3)", nil, true},
	{"sum(1;2;3;4;5)", nil, true},
	{"sum(k; 1; n; k^2)", &Expr{atom: specialFunction{kind: sumFn, args: []*Expr{
		newVarExpr('k'), newNb(1), newVarExpr('n'), {atom: pow, left: newVarExpr('k'), right: newNb(2)},
	}}}, false},
	// comparison
	{"1 <", nil, true},
	{">= 4", nil, true},
	{" 2 = 4 ", nil, true},
	{
		" 2 == 4 ",
		&Expr{
			atom:  equals,
			left:  newNb(2),
			right: newNb(4),
		},
		false,
	},
	{
		"(1<2) + (3>4)",
		&Expr{
			atom:  plus,
			left:  &Expr{atom: strictlyLesser, left: newNb(1), right: newNb(2)},
			right: &Expr{atom: strictlyGreater, left: newNb(3), right: newNb(4)},
		},
		false,
	},
	// indice node
	{
		"_{2+1}", nil, true, // missing name before _
	},
	{
		"a_1 _{2+1}", nil, true, // compound name are not accepted yet
	},
	{
		"a{2+1}", nil, true, // missing _ starter
	},
	{
		"a_{2+1", nil, true, // missing closing }
	},
	{
		"a_{2+}", nil, true, // invalid indice expr
	},
	{ // simple indice
		"a_{n}",
		&Expr{left: newVarExpr('a'), atom: indice{}, right: newVarExpr('n')},
		false,
	},
	{ // complexe indice
		"a_{n - k}",
		&Expr{left: newVarExpr('a'), atom: indice{}, right: &Expr{
			left:  newVarExpr('n'),
			atom:  minus,
			right: newVarExpr('k'),
		}},
		false,
	},
	{ // nested indice
		"a_{b_{n}}",
		&Expr{
			left:  newVarExpr('a'),
			atom:  indice{},
			right: &Expr{left: newVarExpr('b'), atom: indice{}, right: newVarExpr('n')},
		},
		false,
	},
	{
		"@_myVar", nil, true, // deprecated syntaxe
	},
	{
		// missing closing quote is accepted
		`"5+2`, NewVarExpr(NewVarI(0, "5+2")), false, // deprecated syntaxe
	},
	{
		// missing closing quote is accepted
		`randChoice("\ge"; "\le")`, &Expr{atom: specialFunction{kind: randChoice, args: []*Expr{
			NewVarExpr(NewVarI(0, `\ge`)),
			NewVarExpr(NewVarI(0, `\le`)),
		}}}, false, // deprecated syntaxe
	},
	// matrix
	{
		"[[1]]", &Expr{atom: matrix{{newNb(1)}}}, false,
	},
	{
		"[[1; 2]]", &Expr{atom: matrix{{newNb(1), newNb(2)}}}, false,
	},
	{
		"[[1; 2]; [3; 4]]", &Expr{atom: matrix{{newNb(1), newNb(2)}, {newNb(3), newNb(4)}}}, false,
	},
	{
		"[[ ]]", &Expr{atom: matrix{nil}}, false,
	},
	{
		"[[1 + ]]", nil, true,
	},
	{
		"[[1]; []]", nil, true,
	},
	{
		"[[1]; ]", nil, true,
	},
	{
		"[[1]; []", nil, true,
	},
	// matrices operations
	{
		"2 * [[1]]", &Expr{atom: mult, left: newNb(2), right: &Expr{atom: matrix{{newNb(1)}}}}, false,
	},
	{
		"2 [[1]]", &Expr{atom: mult, left: newNb(2), right: &Expr{atom: matrix{{newNb(1)}}}}, false,
	},
	{
		"[[1]] * [[3]]", &Expr{atom: mult, left: &Expr{atom: matrix{{newNb(1)}}}, right: &Expr{atom: matrix{{newNb(3)}}}}, false,
	},
	// matrices functions
	{
		"trace(A)", &Expr{atom: traceFn, right: newVarExpr('A')}, false,
	},
	{
		"det(A)", &Expr{atom: detFn, right: newVarExpr('A')}, false,
	},
	{
		"trans(A)", &Expr{atom: transposeFn, right: newVarExpr('A')}, false,
	},
	{
		"inv(A)", &Expr{atom: invertFn, right: newVarExpr('A')}, false,
	},
	{
		"coeff(A; i; j)", &Expr{atom: specialFunction{kind: matCoeff, args: []*Expr{newVarExpr('A'), newVarExpr('i'), newVarExpr('j')}}}, false,
	},
	{
		"coeff(A; i)", nil, true,
	},
	{
		"binom(1)", nil, true,
	},
	{
		"binom(1; 2 ; 3)", nil, true,
	},
	// real world examples
	{ // x is a variable, not a multiplication
		"104 + 31x11 ", nil, true,
	},
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

func Test_invalidrandChoice(t *testing.T) {
	expr := "randChoice(U;V"
	_, err := Parse(expr)
	if !strings.Contains(err.Error(), "délimiteur manquant") {
		t.Fatal(err)
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

func TestIsInt(t *testing.T) {
	if _, ok := IsInt(math.Inf(1)); ok {
		t.Fatal()
	}
	if _, ok := IsInt(math.Inf(-1)); ok {
		t.Fatal()
	}
	if _, ok := IsInt(math.NaN()); ok {
		t.Fatal()
	}
}

func TestExtendedPower(t *testing.T) {
	_, err := Parse("a^-2")
	tu.AssertNoErr(t, err)
	_, err = Parse("a^-(2+3)")
	tu.AssertNoErr(t, err)
	_, err = Parse("a^2.5")
	tu.AssertNoErr(t, err)
}
