package testutils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func ShouldPanic(t *testing.T, f func()) {
	t.Helper()

	defer func() { recover() }()
	f()
	t.Errorf("should have panicked")
}

func Assert(t *testing.T, b bool) {
	t.Helper()
	if !b {
		t.Fatalf("assertion error")
	}
}

func AssertNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

// ReadEnv read [filename] and return its environnement variables.
func ReadEnv(filename string) map[string]string {
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(content), "\n")
	out := make(map[string]string)
	for _, line := range lines {
		chunks := strings.Split(line, "=")
		if len(chunks) != 2 {
			continue
		}
		k, v := chunks[0], chunks[1]
		out[k] = v
		fmt.Printf("Env. var. %s=%s\n", k, v)
	}
	return out
}

// GenerateLatex create a document, fill it with [body]
// and compile the latex code to create [outFile]
func GenerateLatex(t *testing.T, body string, outFile string) {
	code := fmt.Sprintf(`
		\documentclass{article}

		\usepackage{fullpage}
		\usepackage[utf8]{inputenc}
		\usepackage{amsmath}
		\usepackage[inline]{enumitem}
        \usepackage{amssymb}
		\usepackage[table]{xcolor}

		\begin{document}
		%s
		\end{document}
	`, body)

	dir := filepath.Join(os.TempDir(), "go-latex")

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		AssertNoErr(t, err)
	}

	err := os.WriteFile(filepath.Join(dir, outFile+".tex"), []byte(code), os.ModePerm)
	AssertNoErr(t, err)

	cmd := exec.Command("pdflatex", outFile+".tex")
	var buf bytes.Buffer
	cmd.Dir = dir
	cmd.Stdout = &buf
	err = cmd.Run()
	if _, ok := err.(*exec.ExitError); ok {
		fmt.Println(code)
		fmt.Println(buf.String())
	}
	AssertNoErr(t, err)

	t.Logf("PDF generated at \n file://%s", filepath.Join(dir, outFile+".pdf"))
}
