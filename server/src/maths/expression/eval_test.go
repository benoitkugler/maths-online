package expression

import (
	"math"
	"reflect"
	"testing"
)

func TestEvalMissingVariable(t *testing.T) {
	e := mustParse(t, "x + y")
	_, err := e.Evaluate(Variables{NewVariable('x'): 7})
	if err == nil {
		t.Fatal()
	}
	_ = err.Error()
}

func TestPrecision(t *testing.T) {
	if !AreFloatEqual(MustEvaluate("3 * (1 - (0.25 + 0.1)) / 4", nil), 0.4875) {
		t.Fatal("precision error")
	}

	// Note : we could use math.Big with a large precision to better handle
	// floating point arithmetic issues
	// b := big.NewFloat(0.25)
	// b.SetPrec(100)
	// b.Add(b, big.NewFloat(0.1))
	// fmt.Println(b)
	// b.Sub(big.NewFloat(1), b)
	// fmt.Println(b)
	// b.Mul(big.NewFloat(0).SetPrec(100).SetInt(big.NewInt(3)), b)
	// fmt.Println(b)
	// b.Quo(b, big.NewFloat(4))
	// fmt.Println(b.Float64())
}

func Test_Expression_eval(t *testing.T) {
	tests := []struct {
		expr     string
		bindings ValueResolver
		want     float64
	}{
		{
			"3 + 2", nil, 5,
		},
		// {
		// 	"3 * (1 - (0.25 + 0.1)) / 4", nil, 0.4875,
		// },
		{
			"3 + exp(0)", nil, 4,
		},
		{
			"sin(0)", nil, 0,
		},
		{
			"cos(0)", nil, 1,
		},
		{
			"tan(0)", nil, 0,
		},
		{
			"asin(0)", nil, 0,
		},
		{
			"acos(1)", nil, 0,
		},
		{
			"atan(0)", nil, 0,
		},
		{
			"abs(-3)", nil, 3,
		},
		{
			"ln(e)", nil, 1,
		},
		{
			"4/2", nil, 2,
		},
		{
			"4/3", nil, 4. / 3,
		},
		{
			"4/3", nil, 4. / 3,
		},
		{
			"4 * 3", nil, 12,
		},
		{
			"4 ^ 3", nil, 64,
		},
		{
			"\u03C0 / 2", nil, math.Pi / 2,
		},
		{
			"1 + 2 * (3 + 2)", nil, 11,
		},
		{
			"1 + 1 * 3 ^ 3 * 2 - 1", nil, 54,
		},
		{
			"x + 2", Variables{NewVariable('x'): 4}, 6,
		},
		{
			"2 + 0 * randInt(1;3)", nil, 2,
		},
		{
			"4 * sgn(-1)", nil, -4,
		},
		{
			"sqrt(16) * sqrt(9)", nil, 4 * 3,
		},
		{
			"2 * sqrt(16) * sqrt(9) * sqrt(25)", nil, 2 * 4 * 3 * 5,
		},
		{
			"4 * sgn(-1) * sgn(1) * sgn(0)", nil, 0,
		},
		{
			"2 * randPrime(8; 12)", nil, 22,
		},
		{
			"2 * isPrime(8)", nil, 0,
		},
		{
			"2 * isPrime(11)", nil, 2,
		},
		{
			"2 * isPrime(-11)", nil, 2,
		},
		{
			"2 * isPrime(11.4)", nil, 0,
		},
		{
			"8 % 3", nil, 2,
		},
		{
			"8.2 % 3", nil, 0, // error, actually
		},
		{
			"11 // 3", nil, 3,
		},
		{
			"8.2 // 3", nil, 0, // error, actually
		},
		{
			"acos(7/sqrt(98)) * 180 / " + string(piRune), nil, 45,
		},
		{
			"sqrt(sqrt(98)^2 - 7^2)", nil, 7,
		},
		{
			"1 * isZero(a-1) + 2 * isZero(a-2) + 3*isZero(a-3)", Variables{NewVariable('a'): 2}, 2,
		},
		{
			"1 * isZero(a^2 - b^2 - c^2) + 2*isZero(b^2 - a^2 - c^2) + 3*isZero(c^2 - a^2 - b^2)", Variables{
				NewVariable('a'): 8,  // BC
				NewVariable('b'): 12, // AC
				NewVariable('c'): 4,  // AB
			}, 0,
		},
		{
			"1 * isZero(a^2 - b^2 - c^2) + 2*isZero(b^2 - a^2 - c^2) + 3*isZero(c^2 - a^2 - b^2)", Variables{
				NewVariable('a'): 3, // BC
				NewVariable('b'): 4, // AC
				NewVariable('c'): 5, // AB
			}, 3,
		},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		got, err := expr.Evaluate(tt.bindings)
		if err != nil {
			t.Fatal(err)
		}
		if got != tt.want {
			t.Errorf("node.eval() = %v, want %v", got, tt.want)
		}
	}
}

func Test_Expression_simplifyNumbers(t *testing.T) {
	tests := []struct {
		expr string
		want string
	}{
		{"3 + x", "3 + x"}, // no op
		{"3 + 2", "5"},
		{"3 - 2", "1"},
		{"x - 0", "x"},
		{"x + 0", "x"},
		{"0+x", "x"},
		{"x * 0", "0"},
		{"0*x", "0"},
		{"0/x", "0"},
		{"x * 1", "x"},
		{"x / 1", "x"},
		{"x ^ 1", "x"},
		{"1 ^ x", "1"},
		{"- 2", "-2"},
		{"3 / 4", "0.75"},
		{"1 + 2*(5 - 3 + 4)", "13"},
		{"1 + x + 2", "1 + x + 2"}, // need commutativity, not handled by simplifyNumbers
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		expr.simplify0And1()
		expr.simplifyNumbers()

		want := mustParse(t, tt.want)
		want.simplifyNumbers()
		if !reflect.DeepEqual(expr, want) {
			t.Errorf("node.simplifyNumbers() = %v, want %v", expr, tt.want)
		}
	}
}

func TestExpression_Substitute(t *testing.T) {
	tests := []struct {
		expr string
		vars Variables
		want string
	}{
		{"a + b", Variables{}, "a+b"},
		{"a + b", Variables{NewVariable('a'): 4}, "4+b"},
		{"a + b / 2*a", Variables{NewVariable('a'): 4}, "4+b/2*4"},
		{"a + b", Variables{NewVariable('a'): 4, NewVariable('b'): 5}, "4+5"},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		expr.Substitute(tt.vars)

		want := mustParse(t, tt.want)
		if !reflect.DeepEqual(expr, want) {
			t.Errorf("Substitute(%s) = %v, want %v", tt.expr, expr, tt.want)
		}
	}
}

func Test_isPrime(t *testing.T) {
	primes := sieveOfEratosthenes(1, 1000)
	for _, p := range primes {
		if !isPrime(p) {
			t.Fatal(p)
		}
	}
}

func TestExpression_Extrema(t *testing.T) {
	tests := []struct {
		expr string
		from float64
		to   float64
		want float64
	}{
		{
			"x", -2, 2, 2,
		},
		{
			"cos(x)", -2, 2, 1,
		},
		{
			"exp(x)", -2, 10, math.Exp(10),
		},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		fn := FunctionDefinition{
			FunctionExpr: FunctionExpr{Function: expr, Variable: NewVariable('x')},
			From:         tt.from, To: tt.to,
		}
		if got := fn.Extrema(); got != tt.want {
			t.Errorf("Expression.Extrema() = %v, want %v", got, tt.want)
		}
	}
}
