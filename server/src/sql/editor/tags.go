package editor

import (
	"errors"
	"sort"
	"strings"

	"github.com/benoitkugler/maths-online/server/src/utils"
)

type TagSection struct {
	Tag     string
	Section Section
}

type Tags []TagSection

func (a Tags) Len() int      { return len(a) }
func (a Tags) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Tags) Less(i, j int) bool {
	ai, aj := a[i], a[j]
	if ai.Section < aj.Section {
		return true
	} else if ai.Section > aj.Section {
		return false
	}
	return ai.Tag < aj.Tag
}

// Tags returns the tags for the given [IdQuestiongroup],
// sorted by section then tag
func (qus QuestiongroupTags) Tags() Tags {
	out := make(Tags, len(qus))
	for i, tag := range qus {
		out[i] = TagSection{tag.Tag, tag.Section}
	}
	sort.Sort(out)
	return out
}

// Tags returns the tags for the given [IdQuestiongroup]
// sorted by section then tag
func (exes ExercicegroupTags) Tags() Tags {
	out := make(Tags, len(exes))
	for i, tag := range exes {
		out[i] = TagSection{tag.Tag, tag.Section}
	}
	sort.Sort(out)
	return out
}

func (ts Tags) asQuestionLinks(id IdQuestiongroup) QuestiongroupTags {
	out := make(QuestiongroupTags, len(ts))
	for i, tag := range ts {
		out[i] = QuestiongroupTag{IdQuestiongroup: id, Tag: tag.Tag, Section: tag.Section}
	}
	return out
}

func (ts Tags) asExerciceLinks(id IdExercicegroup) ExercicegroupTags {
	out := make(ExercicegroupTags, len(ts))
	for i, tag := range ts {
		out[i] = ExercicegroupTag{IdExercicegroup: id, Tag: tag.Tag, Section: tag.Section}
	}
	return out
}

func (ts Tags) normalize() (Tags, error) {
	var (
		out                Tags
		nbLevel, nbChapter int
	)
	for _, tag := range ts {
		// enforce proper tags
		tag.Tag = NormalizeTag(tag.Tag)
		if tag.Tag == "" {
			continue
		}

		switch tag.Section {
		case Level:
			nbLevel++
		case Chapter:
			nbChapter++
		}
		out = append(out, tag)
	}

	if nbLevel > 1 {
		return nil, errors.New("Un seul niveau est autorisé par ressource.")
	}
	if nbChapter > 1 {
		return nil, errors.New("Un seul chapitre est autorisé par ressource.")
	}
	return out, nil
}

// List returns the tags (in the same order).
func (ts Tags) List() []string {
	out := make([]string, len(ts))
	for i, tag := range ts {
		out[i] = tag.Tag
	}
	return out
}

// Index extract level and chapter information.
func (ts Tags) Index() (out TagIndex) {
	for _, tag := range ts {
		switch tag.Section {
		case Level:
			out.Level = LevelTag(tag.Tag)
		case Chapter:
			out.Chapter = tag.Tag
		}
	}
	return out
}

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

// TagIndex summarize the classification induced by tags
type TagIndex struct {
	Level   LevelTag
	Chapter string
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

// Crible is a set of tags, not using the section.
// TODO: check that
type Crible map[string]bool

func NewCrible(tags []string) Crible {
	out := make(Crible, len(tags))
	for _, tag := range tags {
		out[tag] = true
	}
	return out
}

// Crible build a set from the tags
func (qus QuestiongroupTags) Crible() Crible { return NewCrible(qus.Tags().List()) }

// Crible build a set from the tags
func (qus ExercicegroupTags) Crible() Crible { return NewCrible(qus.Tags().List()) }

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
