package expression

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
)

func TestEvalMissingVariable(t *testing.T) {
	e := mustParse(t, "x + y")
	_, err := e.Evaluate(Variables{NewVar('x'): NewNb(7)})
	if err == nil {
		t.Fatal()
	}
	_ = err.Error()
}

func TestRoundFloat(t *testing.T) {
	tests := []struct {
		args float64
		want float64
	}{
		{1, 1},
		{1.123, 1.123},
		{1 - 7./100, 0.93},
		{1908 * (1 - 68./100), 610.56},
	}
	for _, tt := range tests {
		if got := RoundFloat(tt.args); got != tt.want {
			t.Errorf("RoundFloat() = %v, want %v", got, tt.want)
		}
	}
}

func TestPrecision(t *testing.T) {
	if !AreFloatEqual(MustEvaluate("3 * (1 - (0.25 + 0.1)) / 4", nil), 0.4875) {
		t.Fatal("precision error")
	}
	if !AreFloatEqual(MustEvaluate("1 - 7/100", nil), 0.93) {
		t.Fatal("precision error")
	}
	if !AreFloatEqual(MustEvaluate("1908 * (1 - 68/100)", nil), 610.56) {
		t.Fatal("precision error")
	}

	// Note : we could use math.Big with a large precision to better handle
	// floating point arithmetic issues, but is seems to worth it, especially
	// since we almost never want the student to precise more than 8 digits

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
			"floor(-3)", nil, -3,
		},
		{
			"floor(-3.5)", nil, -4,
		},
		{
			"floor(3)", nil, 3,
		},
		{
			"floor(3.8)", nil, 3,
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
			"x + 2", Variables{NewVar('x'): NewNb(4)}, 6,
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
			"2 * randInt(8; a)", Variables{NewVar('a'): NewNb(8)}, 16,
		},
		{
			"2 * randChoice(8)", nil, 16,
		},
		{
			"0 * randChoice(8; -1)", nil, 0,
		},
		{
			"0 * randDecDen( )", nil, 0,
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
			"randLetter(A;B)", nil, 0, // 0 by convention
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
			"1 * isZero(a-1) + 2 * isZero(a-2) + 3*isZero(a-3)", Variables{NewVar('a'): NewNb(2)}, 2,
		},
		{
			"1 * isZero(a^2 - b^2 - c^2) + 2*isZero(b^2 - a^2 - c^2) + 3*isZero(c^2 - a^2 - b^2)", Variables{
				NewVar('a'): NewNb(8),  // BC
				NewVar('b'): NewNb(12), // AC
				NewVar('c'): NewNb(4),  // AB
			}, 0,
		},
		{
			"1 * isZero(a^2 - b^2 - c^2) + 2*isZero(b^2 - a^2 - c^2) + 3*isZero(c^2 - a^2 - b^2)", Variables{
				NewVar('a'): NewNb(3), // BC
				NewVar('b'): NewNb(4), // AC
				NewVar('c'): NewNb(5), // AB
			}, 3,
		},
		{
			"round(2.235; 2)", nil, 2.24,
		},
		{
			"min(2.235; 2)", nil, 2,
		},
		{
			"max(-2; 1.4; 5)", nil, 5,
		},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		got, err := expr.Evaluate(tt.bindings)
		if err != nil {
			t.Fatal(err)
		}
		if got != tt.want {
			t.Errorf("node.eval(%s) = %v, want %v", tt.expr, got, tt.want)
		}
	}
}

func TestExpression_Evaluate_err(t *testing.T) {
	tests := []struct {
		expr     string
		bindings Variables
	}{
		{
			"randInt(1;b)", nil,
		},
		{
			"randPrime(1;b)", nil,
		},
		{
			"randInt(a;b)", nil,
		},
		{
			"randInt(a;3)", Variables{NewVar('a'): NewNb(6)},
		},
		{
			"randPrime(a;3)", Variables{NewVar('a'): NewNb(-6)},
		},
		{
			"randPrime(a;9)", Variables{NewVar('a'): NewNb(8)},
		},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		_, err := expr.Evaluate(tt.bindings)
		if err == nil {
			t.Fatal("expected error on", tt.expr)
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

func Test_isPrime(t *testing.T) {
	primes := sieveOfEratosthenes(1, 1000)
	for _, p := range primes {
		if !isPrime(p) {
			t.Fatal(p)
		}
	}
}

func TestIsDecimal(t *testing.T) {
	atom := specialFunctionA{kind: randDenominator}
	for range [200]int{} {
		n, err := atom.eval(newRat(0), newRat(0), nil)
		if err != nil {
			t.Fatal(err)
		}
		if _, is := isInt(n.eval() * maxDecDen); !is {
			t.Fatal(n)
		}
		if n.eval() <= 0 || n.eval() > thresholdDecDen {
			t.Fatal(n)
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
		{
			"sqrt(x)", -2, 2, -1,
		},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		fn := FunctionDefinition{
			FunctionExpr: FunctionExpr{Function: expr, Variable: NewVar('x')},
			From:         tt.from, To: tt.to,
		}
		if got := fn.extrema(); got != tt.want {
			t.Errorf("Expression.Extrema() = %v, want %v", got, tt.want)
		}
	}
}

func Test_isFloatExceedingPrecision(t *testing.T) {
	tests := []struct {
		args float64
		want bool
	}{
		{1, false},
		{1.1456, false},
		{1.145678, false},
		{1.12345678, false},
		{1.12345678901, true},
		{103304.93, false},
		{1. / 3, true},
		{-1. / 3, true},
		{-0.55, false},
	}
	for _, tt := range tests {
		if got := isFloatExceedingPrecision(tt.args); got != tt.want {
			t.Errorf("isFloatExceedingPrecision() = %v, want %v", got, tt.want)
		}
	}

	for range [200]int{} {
		v := 1000 * rand.Float64()
		if vr := roundTo(v, 2); isFloatExceedingPrecision(vr) {
			t.Fatal(vr, v)
		}
	}
}

func Test_roundTo(t *testing.T) {
	tests := []struct {
		v      float64
		digits int
		want   float64
	}{
		{2.256, 1, 2.3},
		{2.224, 1, 2.2},
		{2, 1, 2},
		{-1.123, 2, -1.12},
		{-1.98, 0, -2},
	}
	for _, tt := range tests {
		if got := roundTo(tt.v, tt.digits); got != tt.want {
			t.Errorf("roundTo() = %v, want %v", got, tt.want)
		}
	}
}

func Test_sumRat(t *testing.T) {
	tests := []struct {
		r1   rat
		r2   rat
		want rat
	}{
		{
			rat{1, 2}, rat{3, 2}, rat{4, 2},
		},
		{
			rat{1.4, 2}, rat{3, 2}, rat{4.4, 2},
		},
		{
			rat{1.4, 2}, rat{3, 1.5}, rat{1.5*1.4 + 6, 3},
		},
	}
	for _, tt := range tests {
		if got := sumRat(tt.r1, tt.r2); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("sumRat() = %v, want %v", got, tt.want)
		}
	}
}

func Test_rat_reduce(t *testing.T) {
	tests := []struct {
		p    float64
		q    float64
		want rat
	}{
		{8, 4, rat{2, 1}},
		{12, 8, rat{3, 2}},
		{45.4, 2, rat{45.4, 2}},
	}
	for _, tt := range tests {
		r := rat{
			p: tt.p,
			q: tt.q,
		}
		r.reduce()
		if r != tt.want {
			t.Errorf("sumRat() = %v, want %v", r, tt.want)
		}
	}
}
