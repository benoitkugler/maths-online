package editor

import (
	"sort"
	"strings"

	"github.com/benoitkugler/maths-online/utils"
)

// DifficultyQuery is an union of tags. An empty slice means no selection : all variants are accepted.
type DifficultyQuery []DifficultyTag

// Match returns `true` if the query accepts `diff`.
func (dq DifficultyQuery) Match(diff DifficultyTag) bool {
	if len(dq) == 0 {
		return true
	}
	for _, query := range dq {
		if query == diff {
			return true
		}
	}
	return false
}

// TagQuery is an intersection of tags, with an optionnal
// one for the difficulty
type TagQuery struct {
	// Union, empty means no criterion
	Difficulties DifficultyQuery
	Tags         []string
}

// NormalizeTag returns `tag` with spaces and accents stripped
// and in upper case.
func NormalizeTag(tag string) string {
	return strings.ToUpper(utils.RemoveAccents(strings.TrimSpace((tag))))
}

// List returns the sorted tags from the `Tag` attribute.
func (qus QuestiongroupTags) List() []string {
	out := make([]string, len(qus))
	for i, tag := range qus {
		out[i] = tag.Tag
	}
	sort.Strings(out)
	return out
}

// List returns the sorted tags from the `Tag` attribute.
func (qus ExercicegroupTags) List() []string {
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

	sort.Strings(out)

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
func (qus QuestiongroupTags) Crible() Crible { return NewCrible(qus.List()) }

// Crible build a set from the tags
func (qus ExercicegroupTags) Crible() Crible { return NewCrible(qus.List()) }

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
