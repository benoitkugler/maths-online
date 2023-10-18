package questions

import (
	"regexp"
	"strings"

	"github.com/benoitkugler/maths-online/server/src/maths/expression"
	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
)

// Interpolated is a string with $<static math>$ or &<expression>&
// delimiters.
// && are allowed in $$.
type Interpolated string

type textOrFormula struct {
	s         string
	isFormula bool
}

// parseFormula looks for $$ $$ lines
func (s Interpolated) parseFormula() (out []textOrFormula) {
	lines := strings.Split(string(s), "\n")
	var currentLines []string
	for _, line := range lines {
		lineT := strings.TrimSpace(line)
		if lineT != "$$" && strings.HasPrefix(lineT, "$$") && strings.HasSuffix(lineT, "$$") {
			// found a formula
			if len(currentLines) != 0 {
				out = append(out, textOrFormula{strings.Join(currentLines, "\n"), false})
				currentLines = nil
			}
			out = append(out, textOrFormula{lineT[2 : len(lineT)-2], true})
		} else {
			currentLines = append(currentLines, line)
		}
	}
	if len(currentLines) != 0 {
		out = append(out, textOrFormula{strings.Join(currentLines, "\n"), false})
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
