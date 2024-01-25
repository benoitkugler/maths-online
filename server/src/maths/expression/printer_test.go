package expression

import (
	"strings"
	"testing"

	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func generateLatex(t *testing.T, lines []string, outFile string) {
	tu.GenerateLatex(t, strings.Join(lines, "\n"), outFile)
}

func TestExpression_String(t *testing.T) {
	tests := []struct {
		expr string
		want string
	}{
		{"2 + x", "2 + x"},
		{"2 / x", "2/x"},
		{"2 * x", "2x"},
		{"a * x", "ax"},
		{"2 * 3", "2 * 3"},
		{"2 - x", "2 - x"},
		{"2 ^ x", "2 ^ x"},
		{"2 ^ (x+1)", "2 ^ (x + 1)"},
		{"e * \u03C0", "e\u03C0"},
		{"\uE001", "\uE001"},

		{"exp(2)", "exp(2)"},
		{"sin(2)", "sin(2)"},
		{"cos(2)", "cos(2)"},
		{"abs(2)", "abs(2)"},
		{"sqrt(2)", "sqrt(2)"},
		{"2 + x + log(10)", "2 + x + log(10)"},
		{"- x + 3", "-x + 3"},
		{"1x", "x"},
		{"+x", "x"},
		{"min(2) + max(3)", "min(2) + max(3)"},
		{"-(-a)", "a"},
		{"floor(4)", "floor(4)"},
		{"x + (-4 + y)", "x - 4 + y"},
		{"(1<2)+(3>4)+(5<=6)+(7>=8)", "(1 < 2) + (3 > 4) + (5 <= 6) + (7 >= 8)"},
		{"1 + 2<3", "1 + 2 < 3"},
		{"-inf", "-inf"},
		{"2*n!", "2n!"},
		{"A \u222A B", "A \u222A B"},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		got := expr.String()
		if got != tt.want {
			t.Errorf("Expression.String() = %v, want %v", got, tt.want)
		}

		expr2 := mustParse(t, got)
		if expr.String() != expr2.String() {
			t.Fatalf("inconsitent String() for %s", tt.expr)
		}
		if !AreExpressionsEquivalent(expr, expr2, SimpleSubstitutions) {
			t.Fatalf("inconsitent String() for %s:  %s", tt.expr, got)
		}
	}
}

func TestExpression_StringRoundtrip(t *testing.T) {
	for _, tt := range expressions {
		if tt.wantErr {
			continue
		}

		expr := mustParse(t, tt.expr)
		got := expr.String()
		expr2 := mustParse(t, got)
		if expr.String() != expr2.String() {
			t.Fatalf("inconsitent String() for %s", tt.expr)
		}
		if !AreExpressionsEquivalent(expr, expr2, SimpleSubstitutions) {
			t.Fatalf("inconsitent String() for %s:  %s", tt.expr, got)
		}
	}
}

// generate formulas.pdf in a temporary directory to perform visual tests
func TestExpression_AsLaTeX(t *testing.T) {
	var lines []string
	for _, expr := range []string{
		"2  + 3",
		"+3",
		"-4.789",
		"sqrt(4 +x^2)",
		"ln(4 +x^2)",
		"exp(4 +x^2)",
		"sin(4 +x^2)",
		"cos(4 +x^2)",
		"tan(4 +x^2)",
		"asin(4 +x^2)",
		"acos(4 +x^2)",
		"atan(4 +x^2)",
		"abs(4 +x^2)",
		"1+x+y",
		"2*(2 + x)",
		"(2+z)*(2 + x)",
		"(2+x)*(2 + z)",
		"(a+b)^2",
		"2 - (x - y)",
		"2 * 2.715",
		"x * y + 2*(x +z)*4.4",
		"2*e^(5*x + y)",
		"2*x*" + string(piRune),
		"2 + x/y + 3*(4+x)/(2 - y)",
		"x^y^z",
		"n^(2n)",
		"(x^y)^z^t",
		"n!",
		"2n!",
		"(n + 4)!",
		"\u0393 + \u0398 + \u03B8", // some greek letters
		"\uE000 +  \uE000 + \uE001 + \u0398 + \u03B8", // custom variables
		"x_A + y_B", // simple indices
		"randInt(3;14) + 2",
		"randPrime(3;14) + 2",
		"randChoice(A;B) + 2",
		"sgn(-8)",
		"isPrime(-8)",
		"8 % 3",
		"9 // 2",
		"24x^2 - 27x + 18",
		"round(x; 4)",
		"floor(x)",
		"forceDecimal(x)",
		`"acompletetext"`,
		`"\ge"`,
		"(1<2)+(3>4)+(5<=6)+(7>=8)+ (4==7)",
		"-inf",
		"inf",
		"+inf",
		"u_{2n+1}",
		"2u_{2n+1}",
		// bug #157
		"λx",
		// bug #172
		"a^n / 2^(2n) ",
		// matrix
		"[[1; 4; x]; [2; y; 4]]",
		"3[[1; 4; x]; [2; y; 4]]",
		"[[1; 4; x]; [2; y; 4]] * [[1; 4; x]; [2; y; 4]]",
		"[[1; 4; x]; [2; y; 4]] * ( [[1; 4; x]; [2; y; 4]] + 1)",
		"trace([[1; 4; x]; [2; y; 4]])",
		"det([[1; 4; x]; [2; y; 4]])",
		"A B trans(A) trans(A + B)",
		"trans([[1; 4; x]; [2; y; 4]])",
		"transpose([[1; 4; x]; [2; y; 4]])",
		"inv([[1; 4; x]; [2; y; 4]])",
		"inv(trans(A))",
		"trans(inv(A))",
		"det(inv(A))",
		"trace(inv(A))",
		"binom(x+5; 2)",
		// sets
		"\u00AC(A \u222A B_1)",
		"A \u222A B_1 \u2229 (\u00AC C \u222A D)",
	} {
		e, err := Parse(expr)
		if err != nil {
			t.Fatal(err)
		}

		code := e.AsLaTeX()
		lines = append(lines, "$$"+code+"$$")
	}

	generateLatex(t, lines, "formulas.tex")
}

func TestParenthesis(t *testing.T) {
	expr := mustParse(t, "((-1)/3)x + 2")

	if latex := expr.AsLaTeX(); strings.ContainsRune(latex, '(') {
		t.Fatal("unexpected parenthesis", latex)
	}
	if s := expr.String(); !strings.ContainsRune(s, '(') {
		t.Fatal("missing parenthesis", s)
	}

	expr = mustParse(t, "1 + 1 + a")
	expr.Substitute(Vars{NewVar('a'): newNb(-2)})
	if s := expr.String(); strings.ContainsRune(s, '(') {
		t.Fatal("unexpected parenthesis :", s)
	}

	expr = mustParse(t, "x ^ (2n)")
	latex := expr.AsLaTeX()
	tu.Assert(t, !strings.ContainsRune(latex, '('))
}

func Test0And1(t *testing.T) {
	for _, test := range []struct {
		expr  string
		latex string
	}{
		{"x + 0", "x"},
		{"x - 0", "x"},
		{"x^3 + 0x^2 - x", "{x}^{3} - x"},
		{"1x", "x"},
		{"1+x", "1 + x"},
		{"+2", "2"},
		{"-1x", "-x"},
		{"a -1x - b", "a - x - b"},
		{"-1x^2", "-{x}^{2}"},
		{"-1(4x + 3)", "-\\left(4 x + 3\\right)"},
		{"-1sqrt(100)", "-\\sqrt{100}"},
		{"x + (-y + 2 - 4 + 5)", "x - y + 2 - 4 + 5"},
	} {
		expr := mustParse(t, test.expr)
		latex := expr.AsLaTeX()
		if latex != test.latex {
			t.Fatalf("expected %s, got %s", test.latex, latex)
		}
	}
}

func TestFloatPrecision(t *testing.T) {
	v := Number(mustEvaluate("1908 * (1 - 68/100)", nil))
	if s := v.String(); s != "610,56" {
		t.Fatal(s)
	}
}

func TestPlusMinus(t *testing.T) {
	expr := mustParse(t, "x + (-5)")
	latex := expr.AsLaTeX()
	if strings.ContainsRune(latex, '+') {
		t.Fatal("unexpected +")
	}

	expr = mustParse(t, "x - (-5)")
	latex = expr.AsLaTeX()
	if strings.ContainsRune(latex, '-') {
		t.Fatal("unexpected -")
	}

	expr = mustParse(t, "y + (-x)")
	latex = expr.AsLaTeX()
	if strings.ContainsRune(latex, '+') {
		t.Fatal("unexpected +")
	}
}

func TestMinusMinus(t *testing.T) {
	expr := mustParse(t, "(-a/b)")
	expr.Substitute(Vars{NewVar('a'): newNb(-2), NewVar('b'): newNb(4)})
	if s := expr.String(); s != "2/4" {
		t.Fatalf("%#v: %s", expr.right, s)
	}

	expr = mustParse(t, "2 - a")
	expr.Substitute(Vars{NewVar('a'): newNb(-2)})
	if s := expr.String(); s != "2 + 2" {
		t.Fatalf("%#v: %s", expr.right, s)
	}

	expr = mustParse(t, "-(-2x)")
	latex := expr.AsLaTeX()
	if strings.ContainsRune(latex, '-') {
		t.Fatal("unexpected +")
	}

	expr = mustParse(t, "-(-2 + x)")
	if s := expr.String(); s != "-(-2 + x)" {
		t.Fatal(s)
	}
}

func TestInstantiateMinusZero(t *testing.T) {
	// related to issue #144
	exprB := mustParse(t, "0 * 1 / (-5)")
	tu.Assert(t, exprB.String() == "0")

	rp := RandomParameters{
		NewVar('b'): exprB,
	}
	vs, err := rp.Instantiate()
	tu.AssertNoErr(t, err)

	instance := vs[NewVar('b')]
	tu.Assert(t, instance.String() == "0")
}

func TestOmitTimes(t *testing.T) {
	e := mustParse(t, "2u_{n}")
	if latex := e.AsLaTeX(); strings.Contains(latex, "\\times") {
		t.Fatalf("times should ommited, got %s", latex)
	}
}

func TestIssue173(t *testing.T) {
	for _, expr := range []string{
		"6/49",
		"3/49",
		"1/4",
	} {
		e := mustParse(t, expr)
		v, err := e.evalReal(nil)
		tu.AssertNoErr(t, err)
		tu.Assert(t, v.toExpr().String() == expr) // should be printed as fraction
	}
}

func TestPrintFractions(t *testing.T) {
	for _, tt := range []struct {
		expr string
		Vars RandomParameters
		want string
	}{
		{"6 / 49", nil, "6/49"},
		{"3 / 49", nil, "3/49"},
		{"1 / 4", nil, "1/4"},
		{"0.25", nil, "0,25"},
		{"x", RandomParameters{NewVar('a'): newNb(1), NewVar('b'): newNb(3), NewVar('x'): mustParse(t, "a / b")}, "1/3"},
		{"x", RandomParameters{NewVar('a'): newNb(2), NewVar('b'): newNb(6), NewVar('x'): mustParse(t, "a / b")}, "1/3"},
		{"x", RandomParameters{NewVar('x'): mustParse(t, "forceDecimal(3 / 4)")}, "0,75"},
	} {
		e := mustParse(t, tt.expr)
		vars, err := tt.Vars.Instantiate()
		tu.AssertNoErr(t, err)
		e.Substitute(vars)

		tu.Assert(t, e.String() == tt.want)
	}
}
