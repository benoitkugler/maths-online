package editor

import (
	"database/sql/driver"
	"errors"
	"fmt"
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

// TagIndex summarize the classification induced by tags
type TagIndex struct {
	Level   LevelTag
	Chapter string
}

// TagGroup groups the tags for one question/exercice,
// used to resolve the tag hierachy
type TagGroup struct {
	TagIndex
	TrivMaths []string
}

// BySection classify the tag list according to section
func (ts Tags) BySection() (out TagGroup) {
	for _, tag := range ts {
		switch tag.Section {
		case Level:
			out.Level = LevelTag(tag.Tag)
		case Chapter:
			out.Chapter = tag.Tag
		case TrivMath:
			out.TrivMaths = append(out.TrivMaths, tag.Tag)
		}
	}
	return out
}

// DifficultyQuery is an union of tags. An empty slice means no selection : all variants are accepted.
type DifficultyQuery []DifficultyTag

func (s *DifficultyQuery) Scan(src interface{}) error  { return loadJSON(s, src) }
func (s DifficultyQuery) Value() (driver.Value, error) { return dumpJSON(s) }

// Match returns `true` if the query accepts `diff`.
// Questions with no difficulty always match
func (dq DifficultyQuery) Match(diff DifficultyTag) bool {
	if len(dq) == 0 || diff == "" {
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

// Crible is a set of tags (and section)
type Crible map[TagSection]bool

func (tags Tags) Crible() Crible {
	out := make(Crible, len(tags))
	for _, tag := range tags {
		out[tag] = true
	}
	return out
}

// HasAll returns `true` is all the `tags` are present in the crible.
func (cr Crible) HasAll(tags Tags) bool {
	for _, tag := range tags {
		if !cr[tag] {
			return false
		}
	}
	return true
}

// TagListSet is a map[[]TagSection]bool,
// where the order in the key list is ignored.
// Tags should be normalized before using this set.
type TagListSet struct {
	m map[string]Tags
}

func NewTagListSet() TagListSet {
	return TagListSet{m: make(map[string]Tags)}
}

// compute a hash, order invariant
func key(tags Tags) string {
	// some very unlikely pattern for tags, to make sure the key function
	// is injective
	const delim = "^-$-^"

	sort.Sort(tags)
	chunks := make([]string, len(tags))
	for i, tag := range tags {
		chunks[i] = fmt.Sprintf("%s:%d", tag.Tag, tag.Section)
	}
	return strings.Join(chunks, delim)
}

func (tls TagListSet) Add(tags Tags) {
	tls.m[key(tags)] = tags
}

func (tls TagListSet) Has(tags Tags) bool {
	_, has := tls.m[key(tags)]
	return has
}

// List returns the content of the set (in arbitrary order)
func (tls TagListSet) List() []Tags {
	out := make([]Tags, 0, len(tls.m))
	for _, v := range tls.m {
		out = append(out, v)
	}
	return out
}
