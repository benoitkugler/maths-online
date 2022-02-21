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
	} {
		e, err := Parse(expr)
		if err != nil {
			t.Fatal(err)
		}

		_ = e.AsLaTeX(nil) // check for panic

		code := e.AsLaTeX(func(v Variable) string {
			if v == '\uE000' {
				return "u_{n}"
			} else if v == '\uE001' {
				return "v_{n}"
			}
			return DefaultLatexResolver(v)
		})
		lines = append(lines, "$$"+code+"$$")
	}

	code := fmt.Sprintf(`
		\documentclass{article}
		\usepackage[utf8]{inputenc}
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
