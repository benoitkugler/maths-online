package proof

import "strings"

// String returns a multi-line, indented text version of the proof.
func (pr ProofPart) String() string {
	var chunks []string
	for _, part := range pr {
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

func (s Statement) String() string { return string(s) }
func (s Equality) String() string  { return strings.Join(s, " = ") }
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
