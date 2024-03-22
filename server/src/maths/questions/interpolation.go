package questions

import (
	"regexp"
	"strings"

	"github.com/benoitkugler/maths-online/server/src/maths/expression"
	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
)

// Interpolated is a string with $<static math>$ or &<expression>&
// delimiters, with && are allowed in $$.
// Also, $$ ... $$ may be used to define a FormulaBlock, and
// # ... # to define NumberFieldBlock
type Interpolated string

const (
	iText uint8 = iota
	iFormula
	iNumberField
)

type textChunck struct {
	s    string
	kind uint8
}

func splitNumberField(text string) (out []textChunck) {
	for i := strings.IndexByte(text, '#'); i != -1; i = strings.IndexByte(text, '#') {
		before, after := text[:i], text[i+1:]
		close := strings.IndexByte(after, '#')
		if close == -1 {
			out = append(out, textChunck{s: text, kind: iText})
			return out
		}
		expr := after[:close]
		if before != "" {
			out = append(out, textChunck{s: before, kind: iText})
		}
		out = append(out, textChunck{s: expr, kind: iNumberField})
		text = after[close+1:]
	}
	if text != "" {
		out = append(out, textChunck{s: text, kind: iText})
	}
	return out
}

// parseFormula looks for $$ $$ lines and # # chunks
func (s Interpolated) parseFormula() (out []textChunck) {
	if s == "" {
		// always return at least one chunk
		return []textChunck{{"", iText}}
	}
	lines := strings.Split(string(s), "\n")
	var (
		currentLines []string
		tmp          []textChunck
	)
	for _, line := range lines {
		lineT := strings.TrimSpace(line)
		if lineT != "$$" && strings.HasPrefix(lineT, "$$") && strings.HasSuffix(lineT, "$$") {
			// found a formula
			if len(currentLines) != 0 {
				tmp = append(tmp, textChunck{strings.Join(currentLines, "\n"), iText})
				currentLines = nil
			}
			tmp = append(tmp, textChunck{lineT[2 : len(lineT)-2], iFormula})
		} else {
			currentLines = append(currentLines, line)
		}
	}
	if len(currentLines) != 0 {
		tmp = append(tmp, textChunck{strings.Join(currentLines, "\n"), iText})
	}

	// parse lines again to look for #
	for _, line := range tmp {
		if line.kind == iFormula {
			out = append(out, line)
			continue
		}
		out = append(out, splitNumberField(line.s)...)
	}
	return out
}

// parse extracts each parts of the interpolated string,
// as well as parsing expressions found.
// It returns an error for invalid expressions
func (s Interpolated) parse() (TextParts, error) {
	latex := splitByLaTeX(string(s))
	var out TextParts
	for _, c := range latex {
		newChunks := splitByExpression(c)
		for _, part := range newChunks {
			if part.Kind == Expression {
				_, err := expression.ParseCompound(part.Content)
				if err != nil {
					return nil, err
				}
			}
		}
		out = append(out, newChunks...)
	}
	return out, nil
}

func (s Interpolated) instantiate(params expression.Vars) (client.TextLine, error) {
	parsed, err := s.parse()
	if err != nil {
		return nil, err
	}
	return parsed.instantiate(params)
}

// instantiateAndMerge parse, instantiate and merge all the chunks back.
// this will produce the expected output on the client if all the contents are to be displayed in the
// same mode
func (s Interpolated) instantiateAndMerge(params expression.Vars) (string, error) {
	parsed, err := s.parse()
	if err != nil {
		return "", err
	}
	return parsed.instantiateAndMerge(params)
}

var (
	reLaTeX      = regexp.MustCompile(`\$([^$]+)\$`)
	reExpression = regexp.MustCompile(`&([^&\n]+)&`)
)

func splitByRegexp(re *regexp.Regexp, s string, kindMatch, kindDefault TextKind) (out []TextPart) {
	var cursor int
	for _, indexes := range re.FindAllStringSubmatchIndex(s, -1) {
		startOuter, endOuter, startInner, endInner := indexes[0], indexes[1], indexes[2], indexes[3]
		if startOuter > cursor {
			out = append(out, TextPart{Content: s[cursor:startOuter], Kind: kindDefault})
		}
		out = append(out, TextPart{Kind: kindMatch, Content: s[startInner:endInner]})
		cursor = endOuter
	}

	if len(s) > cursor {
		out = append(out, TextPart{Content: s[cursor:], Kind: kindDefault})
	}

	return out
}

// return either Text or StaticMath
func splitByLaTeX(s string) (out []TextPart) {
	return splitByRegexp(reLaTeX, s, StaticMath, Text)
}

func splitByExpression(t TextPart) []TextPart {
	return splitByRegexp(reExpression, t.Content, Expression, t.Kind)
}
