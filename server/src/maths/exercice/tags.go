package exercice

// List returns the tags from the relation table.
func (qus QuestionTags) List() []string {
	out := make([]string, len(qus))
	for index, qt := range qus {
		out[index] = qt.Tag
	}
	return out
}

// Crible is a set of tags.
type Crible map[string]bool

// Crible build a set from the tags
func (qus QuestionTags) Crible() Crible {
	out := make(Crible, len(qus))
	for _, qt := range qus {
		out[qt.Tag] = true
	}
	return out
}

// Difficulty returns the difficulty of the question,
// or an empty string.
func (cr Crible) Difficulty() DifficultyTag {
	if cr[string(Diff1)] {
		return Diff1
	} else if cr[string(Diff2)] {
		return Diff2
	} else if cr[string(Diff3)] {
		return Diff3
	}
	return ""
}

// HasAll returns `true` is all the `tags` are present in the crible.
func (cr Crible) HasAll(tags []string) bool {
	for _, tag := range tags {
		if !cr[tag] {
			return false
		}
	}
	return true
}
