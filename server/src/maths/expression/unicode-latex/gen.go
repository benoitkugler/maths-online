package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
)

//go:embed unicode-math-symbols.txt
var symbols []byte

func importData() (map[rune]string, error) {
	out := make(map[rune]string)

	r := bufio.NewScanner(bytes.NewReader(symbols))
	for r.Scan() {
		if err := r.Err(); err != nil {
			return nil, err
		}

		line := r.Text()
		if strings.HasPrefix(line, "#") { // ignore comments
			continue
		}

		r := csv.NewReader(strings.NewReader(line))
		r.Comma = '^'
		fields, err := r.Read()
		if err != nil {
			return nil, err
		}

		charString, latex := fields[0], fields[2]

		if latex == "" {
			continue // ignore unknown runes
		}

		var char rune
		_, err = fmt.Sscanf(charString, "%05x", &char)
		if err != nil {
			return nil, fmt.Errorf("invalid rune litteral %s: %s", charString, err)
		}

		out[char] = latex
	}

	return out, nil
}

func generateCode(m map[rune]string) string {
	var lines []string
	for k, v := range m {
		lines = append(lines, fmt.Sprintf("0x%x: %q,", k, v))
	}

	sort.Slice(lines, func(i, j int) bool { return lines[i] < lines[j] })

	return fmt.Sprintf(`package expression

	// Code generated by unicode-latex/gen.go DO NOT EDIT

	var unicodeToLaTeX = map[rune]string{
		%s
	}
	`, strings.Join(lines, "\n"))
}

func main() {
	m, err := importData()
	if err != nil {
		log.Fatal(err)
	}

	code := generateCode(m)

	out := "printer_latex_unicode.go"
	err = os.WriteFile(out, []byte(code), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	err = exec.Command("goimports", "-w", out).Run()
	if err != nil {
		log.Fatal(err)
	}
}
