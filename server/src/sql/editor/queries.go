package editor

import (
	"database/sql"
	"sort"

	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/benoitkugler/maths-online/server/src/utils"
)

type OptionalIdExercice struct {
	Valid bool
	ID    IdExercice
}

func (id IdExercice) AsOptional() OptionalIdExercice {
	return OptionalIdExercice{ID: id, Valid: true}
}

type OptionalIdQuestiongroup struct {
	ID    IdQuestiongroup
	Valid bool
}

func (id IdQuestiongroup) AsOptional() OptionalIdQuestiongroup {
	return OptionalIdQuestiongroup{ID: id, Valid: true}
}

// selectQuestiongroupsByTag returns the question groups matching the given tag.
// It is meant to avoid loading the whole tags table.
func selectQuestiongroupsByTag(db DB, tag TagSection) (Questiongroups, error) {
	rs, err := db.Query(`SELECT questiongroups.*
	FROM questiongroups
	JOIN questiongroup_tags ON questiongroups.id = questiongroup_tags.IdQuestiongroup
   	WHERE questiongroup_tags.tag = $1 AND questiongroup_tags.section = $2`, tag.Tag, tag.Section)
	if err != nil {
		return nil, err
	}
	return ScanQuestiongroups(rs)
}

// IsVisibleBy returns `true` if the question is public or
// owned by `userID`
func (qu Questiongroup) IsVisibleBy(userID teacher.IdTeacher) bool {
	return qu.Public || qu.IdTeacher == userID
}

// RestrictVisible remove the questions not visible by `userID`
func (qus Questiongroups) RestrictVisible(userID teacher.IdTeacher) {
	for id, qu := range qus {
		if !qu.IsVisibleBy(userID) {
			delete(qus, id)
		}
	}
}

// RestrictNeedExercice remove the questions marked as requiring an
// exercice
func (qus Questions) RestrictNeedExercice() {
	for id, question := range qus {
		if question.NeedExercice.Valid {
			delete(qus, id)
		}
	}
}

func (qus Questions) ByGroup() map[IdQuestiongroup][]Question {
	out := make(map[IdQuestiongroup][]Question)
	for _, question := range qus {
		idGroup := question.IdGroup.ID
		out[idGroup] = append(out[idGroup], question)
	}
	return out
}

// UpdateQuestiongroupTags sets the tags of [id], normalizing and validating [tags]
// It does NOT commit or rollback.
func UpdateQuestiongroupTags(tx *sql.Tx, tags Tags, id IdQuestiongroup) error {
	tags, err := tags.normalize()
	if err != nil {
		return err
	}

	_, err = DeleteQuestiongroupTagsByIdQuestiongroups(tx, id)
	if err != nil {
		return utils.SQLError(err)
	}
	err = InsertManyQuestiongroupTags(tx, tags.asQuestionLinks(id)...)
	if err != nil {
		return utils.SQLError(err)
	}
	return nil
}

// UpdateExercicegroupTags enforces proper IdExercicegroup, mutating `tags`.
// It does NOT commit or rollback.
func UpdateExercicegroupTags(tx *sql.Tx, tags Tags, id IdExercicegroup) error {
	tags, err := tags.normalize()
	if err != nil {
		return err
	}

	_, err = DeleteExercicegroupTagsByIdExercicegroups(tx, id)
	if err != nil {
		return utils.SQLError(err)
	}
	err = InsertManyExercicegroupTags(tx, tags.asExerciceLinks(id)...)
	if err != nil {
		return utils.SQLError(err)
	}
	return nil
}

// UpdateExerciceQuestionList set the questions for the given exercice,
// overiding `IdExercice` and `index` fields of the list items.
// Do NOT commit or rollback
func UpdateExerciceQuestionList(tx *sql.Tx, idExercice IdExercice, l ExerciceQuestions) error {
	// enforce fields
	for i := range l {
		l[i].Index = i
		l[i].IdExercice = idExercice
	}

	_, err := DeleteExerciceQuestionsByIdExercices(tx, idExercice)
	if err != nil {
		return utils.SQLError(err)
	}

	err = InsertManyExerciceQuestions(tx, l...)
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}

// IsVisibleBy returns `true` if the exercice group is public or
// owned by `userID`
func (qu Exercicegroup) IsVisibleBy(userID teacher.IdTeacher) bool {
	return qu.Public || qu.IdTeacher == userID
}

// RestrictVisible remove the questions not visible by `userID`
func (qus Exercicegroups) RestrictVisible(userID teacher.IdTeacher) {
	for id, ex := range qus {
		if !ex.IsVisibleBy(userID) {
			delete(qus, id)
		}
	}
}

func (qus Exercices) ByGroup() map[IdExercicegroup][]Exercice {
	out := make(map[IdExercicegroup][]Exercice)
	for _, exercice := range qus {
		idGroup := exercice.IdGroup
		out[idGroup] = append(out[idGroup], exercice)
	}
	return out
}

// SelectQuestiongroupByTags returns the question groups matching the given query,
// and available to `userID`, returning their tags.
// It panics if `pattern` is empty.
func SelectQuestiongroupByTags(db DB, userID teacher.IdTeacher, pattern Tags) (map[IdQuestiongroup]QuestiongroupTags, error) {
	firstSelection, err := selectQuestiongroupsByTag(db, pattern[0])
	if err != nil {
		return nil, err
	}

	quTags, err := SelectQuestiongroupTagsByIdQuestiongroups(db, firstSelection.IDs()...)
	if err != nil {
		return nil, err
	}

	var (
		selectedIDs []IdQuestiongroup
		tagsByGroup = quTags.ByIdQuestiongroup()
	)
	// remove questions not matching all the tags
	for idGroup, cr := range tagsByGroup {
		hasAll := cr.Tags().Crible().HasAll(pattern)
		if !hasAll {
			delete(tagsByGroup, idGroup)
		} else {
			selectedIDs = append(selectedIDs, idGroup)
		}
	}

	groups, err := SelectQuestiongroups(db, selectedIDs...)
	if err != nil {
		return nil, utils.SQLError(err)
	}

	// restrict to user visible groups
	for _, group := range groups {
		if !group.IsVisibleBy(userID) {
			delete(tagsByGroup, group.Id)
		}
	}

	return tagsByGroup, nil
}

// EnsureOrder must be call on the questions of one exercice,
// to make sure the order in the slice is consistent with the one
// indicated by `Index`
func (l ExerciceQuestions) EnsureOrder() {
	sort.Slice(l, func(i, j int) bool { return l[i].Index < l[j].Index })
}
