package expression

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

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
		"(x^y)^z^t",
		"\u0393 + \u0398 + \u03B8", // some greek letters
		"\uE000 +  \uE000 + \uE001 + \u0398 + \u03B8", // custom variables
		"x_A + y_B", // indices
		"randInt(3;14) + 2",
		"randPrime(3;14) + 2",
		"sgn(-8)",
		"isZero(-8)",
		"isPrime(-8)",
		"8 % 3",
		"9 // 2",
		"24x^2 - 27x + 18",
	} {
		e, _, err := Parse(expr)
		if err != nil {
			t.Fatal(err)
		}

		_ = e.AsLaTeX(nil) // check for panic

		code := e.AsLaTeX(func(v Variable) string {
			if v.Name == '\uE000' {
				return "u_{n+1}"
			} else if v.Name == '\uE001' {
				return "v_{n+2}"
			}
			return DefaultLatexResolver(v)
		})
		lines = append(lines, "$$"+code+"$$")
	}

	code := fmt.Sprintf(`
		\documentclass{article}
		\usepackage[utf8]{inputenc}
		\usepackage{amsmath}

		\begin{document}
		%s
		\end{document}
	`, strings.Join(lines, "\n"))

	dir := filepath.Join(os.TempDir(), "go-latex")

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
	}

	err := os.WriteFile(filepath.Join(dir, "formulas.tex"), []byte(code), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("pdflatex", "formulas.tex")
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestParenthesis(t *testing.T) {
	expr := mustParse(t, "((-1)/3)x + 2")
	latex := expr.AsLaTeX(nil)
	if strings.ContainsRune(latex, '(') {
		t.Fatal("unexpected parenthesis")
	}
}

func TestPlusMinus(t *testing.T) {
	expr := mustParse(t, "x + (-5)")
	latex := expr.AsLaTeX(nil)
	if strings.ContainsRune(latex, '+') {
		t.Fatal("unexpected +")
	}

	expr = mustParse(t, "x - (-5)")
	latex = expr.AsLaTeX(nil)
	if strings.ContainsRune(latex, '-') {
		t.Fatal("unexpected -")
	}
}

func Test0And1(t *testing.T) {
	for _, test := range []struct {
		expr  string
		latex string
	}{
		{"x + 0", "x"},
		{"x - 0", "x"},
		{"1x", "x"},
		{"1+x", "1 + x"},
		{"+2", "2"},
	} {
		expr := mustParse(t, test.expr)
		latex := expr.AsLaTeX(nil)
		if latex != test.latex {
			t.Fatalf("expected %s, got %s", test.latex, latex)
		}
	}
}
