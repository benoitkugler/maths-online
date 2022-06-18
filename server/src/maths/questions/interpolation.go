package questions

import (
	"regexp"

	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/questions/client"
)

// Interpolated is a string with $<static math>$ or #{<expression>}
// delimiters.
// #{} are allowed in $$.
type Interpolated string

// Parse extracts each parts of the interpolated string,
// as well as parsing expressions found.
// It returns an error for invalid expressions
func (s Interpolated) Parse() (TextParts, error) {
	latex := splitByLaTeX(string(s))
	var out TextParts
	for _, c := range latex {
		newChunks := splitByExpression(c)
		for _, part := range newChunks {
			if part.Kind == Expression {
				_, err := expression.Parse(part.Content)
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
	parsed, err := s.Parse()
	if err != nil {
		return nil, err
	}
	return parsed.instantiate(params)
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
