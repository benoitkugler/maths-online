package editor

import (
	"strings"

	ex "github.com/benoitkugler/maths-online/maths/exercice"
	"github.com/benoitkugler/maths-online/maths/expression"
)

func NewTextBlock(text ex.TextBlock) TextBlock {
	var chunks []string
	for _, part := range text.Parts {
		var chunk string
		switch part.Kind {
		case ex.Text: // nothing to do
			chunk = part.Content
		case ex.StaticMath: // wrap with $ $
			chunk = "$" + part.Content + "$"
		case ex.Expression: // xrap with #{}
			chunk = "#{" + part.Content + "}"
		default:
			panic(ex.ExhaustiveTextKind)
		}
		chunks = append(chunks, chunk)
	}
	return TextBlock{
		Parts:  strings.Join(chunks, ""),
		IsHint: text.IsHint,
	}
}

type parser struct {
	src []rune
	pos int
}

func newParser(s string) *parser {
	return &parser{src: []rune(s)}
}

func (pr *parser) nextPart() (_ ex.TextPart, eof bool, err error) {
	if pr.pos >= len(pr.src) {
		return ex.TextPart{}, true, nil
	}

	switch pr.src[pr.pos] {
	case '$':
		pr.pos++
		start := pr.pos
		for ; pr.pos < len(pr.src) && pr.src[pr.pos] != '$'; pr.pos++ {
		}
		part := string(pr.src[start:pr.pos])
		pr.pos++ // skip the $
		return ex.TextPart{Content: part, Kind: ex.StaticMath}, false, nil
	case '#':
		pr.pos++
		if pr.pos < len(pr.src) && pr.src[pr.pos] == '{' {
			pr.pos++
			start := pr.pos
			for ; pr.pos < len(pr.src) && pr.src[pr.pos] != '}'; pr.pos++ {
			}
			part := string(pr.src[start:pr.pos])

			pr.pos++ // skip the }

			_, _, err := expression.Parse(part) // TODO: adjust feedback offset
			if err != nil {
				return ex.TextPart{}, false, err
			}
			return ex.TextPart{Content: part, Kind: ex.Expression}, false, nil
		}
	}

	// advance until reaching a delimiter
	start := pr.pos
	for ; pr.pos < len(pr.src); pr.pos++ {
		if pr.src[pr.pos] == '$' || (pr.pos+1 < len(pr.src) && pr.src[pr.pos] == '#' && pr.src[pr.pos+1] == '{') {
			break
		}
	}
	return ex.TextPart{Content: string(pr.src[start:pr.pos]), Kind: ex.Text}, false, nil
}

// ParseInterpolatedString expects a string with $<static math>$ or #{<expression>}
// delimiters.
func ParseInterpolatedString(s string) (ex.TextParts, error) {
	var out ex.TextParts
	pr := newParser(s)
	for {
		part, eof, err := pr.nextPart()
		if err != nil {
			return nil, err
		}
		if eof {
			return out, nil
		}
		out = append(out, part)
	}
}
