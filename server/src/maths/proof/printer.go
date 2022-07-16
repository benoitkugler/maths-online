package proof

import "strings"

// String returns a multi-line, indented text version of the proof.
func (pr Sequence) String() string {
	var chunks []string
	for _, part := range pr.Parts {
		chunks = append(chunks, part.String())
	}
	return strings.Join(chunks, "\n  donc\n")
}

func indentBlock(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = "\t" + line
	}
	return strings.Join(lines, "\n")
}

func (s Statement) String() string { return s.Content }
func (s Equality) String() string  { return strings.Join(s.Terms, " = ") }
func (s Node) String() string {
	left := s.Left.String()
	right := s.Right.String()
	var sep string
	switch s.Op {
	case And:
		sep = "et"
	case Or:
		sep = "ou"
	default:
		panic("invalid operator")
	}
	return indentBlock(left + "\n  " + sep + "\n" + right)
}
