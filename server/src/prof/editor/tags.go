package editor

import (
	"sort"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// NormalizeTag returns `tag` with accents stripped
// and in upper case.
func NormalizeTag(tag string) string {
	return strings.ToUpper(removeAccents(strings.TrimSpace((tag))))
}

func NormalizeTags(tags []string) {
	for i := range tags {
		tags[i] = NormalizeTag(tags[i])
	}
}

func removeAccents(s string) string {
	noAccent := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(noAccent, s)
	if e != nil {
		return s
	}
	return output
}

// List returns the sorted tags from the `Tag` attribute.
func (qus QuestionTags) List() []string {
	out := make([]string, len(qus))
	for i, tag := range qus {
		out[i] = tag.Tag
	}
	sort.Strings(out)
	return out
}

// CommonTags returns the tags found in every list.
func CommonTags(tags [][]string) []string {
	L := len(tags)
	crible := make(map[string][]bool)

	for index, inter := range tags {
		for _, tag := range inter {
			list := crible[tag]
			if list == nil {
				list = make([]bool, L)
			}
			list[index] = true
			crible[tag] = list
		}
	}
	var out []string
	for tag, occurences := range crible {
		isEverywhere := true
		for _, b := range occurences {
			if !b {
				isEverywhere = false
				break
			}
		}
		if isEverywhere {
			out = append(out, tag)
		}
	}
	return out
}

// Crible is a set of tags.
type Crible map[string]bool

func NewCrible(tags []string) Crible {
	out := make(Crible, len(tags))
	for _, tag := range tags {
		out[tag] = true
	}
	return out
}

// Crible build a set from the tags
func (qus QuestionTags) Crible() Crible {
	tmp := qus.List()
	NormalizeTags(tmp)
	return NewCrible(tmp)
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

// TagListSet is a map[[]string]bool,
// where the order in the key list is ignored.
// Tags should be normalized before using this set.
type TagListSet struct {
	m map[string]bool
}

func NewTagListSet() TagListSet {
	return TagListSet{m: make(map[string]bool)}
}

const delim = "^-$-^" // some very unlikely pattern for tags

func key(tags []string) string {
	sort.Strings(tags)
	return strings.Join(tags, delim)
}

func (tls TagListSet) Add(tags []string) {
	tls.m[key(tags)] = true
}

func (tls TagListSet) Has(tags []string) bool {
	return tls.m[key(tags)]
}

func (tls TagListSet) List() [][]string {
	out := make([][]string, 0, len(tls.m))
	for k := range tls.m {
		out = append(out, strings.Split(k, delim))
	}
	return out
}
