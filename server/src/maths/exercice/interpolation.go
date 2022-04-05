package exercice

import (
	"github.com/benoitkugler/maths-online/maths/expression"
)

// func NewTextBlock(text TextBlock) TextBlock {
// 	var chunks []string
// 	for _, part := range text.Parts {
// 		var chunk string
// 		switch part.Kind {
// 		case Text: // nothing to do
// 			chunk = part.Content
// 		case StaticMath: // wrap with $ $
// 			chunk = "$" + part.Content + "$"
// 		case Expression: // xrap with #{}
// 			chunk = "#{" + part.Content + "}"
// 		default:
// 			panic(ExhaustiveTextKind)
// 		}
// 		chunks = append(chunks, chunk)
// 	}
// 	return TextBlock{
// 		Parts:  strings.Join(chunks, ""),
// 		IsHint: text.IsHint,
// 	}
// }

type parser struct {
	src []rune
	pos int
}

func newParser(s string) *parser {
	return &parser{src: []rune(s)}
}

func (pr *parser) nextPart() (_ TextPart, eof bool, err error) {
	if pr.pos >= len(pr.src) {
		return TextPart{}, true, nil
	}

	switch pr.src[pr.pos] {
	case '$':
		pr.pos++
		start := pr.pos
		for ; pr.pos < len(pr.src) && pr.src[pr.pos] != '$'; pr.pos++ {
		}
		part := string(pr.src[start:pr.pos])
		pr.pos++ // skip the $
		return TextPart{Content: part, Kind: StaticMath}, false, nil
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
				return TextPart{}, false, err
			}
			return TextPart{Content: part, Kind: Expression}, false, nil
		}
	}

	// advance until reaching a delimiter
	start := pr.pos
	for ; pr.pos < len(pr.src); pr.pos++ {
		if pr.src[pr.pos] == '$' || (pr.pos+1 < len(pr.src) && pr.src[pr.pos] == '#' && pr.src[pr.pos+1] == '{') {
			break
		}
	}
	return TextPart{Content: string(pr.src[start:pr.pos]), Kind: Text}, false, nil
}

// Interpolated is a string with $<static math>$ or #{<expression>}
// delimiters.
type Interpolated string

// Parse extracts each parts of the interpolated string.
func (s Interpolated) Parse() (TextParts, error) {
	var out TextParts
	pr := newParser(string(s))
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
