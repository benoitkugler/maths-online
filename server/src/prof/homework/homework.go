// Package homework implements an activity for student
// consisting in personal, at home training on exercices given by the
// teacher.
package homework

import (
	"database/sql"
	"sort"

	ed "github.com/benoitkugler/maths-online/prof/editor"
	"github.com/benoitkugler/maths-online/prof/teacher"
	"github.com/benoitkugler/maths-online/utils"
)

type ClassroomSheets struct {
	Classroom teacher.Classroom
	Sheets    []SheetExt
}

func newClassroomSheets(cl teacher.Classroom, sheepMap map[IdSheet]SheetExt) ClassroomSheets {
	out := ClassroomSheets{Classroom: cl}
	for _, sheet := range sheepMap {
		if sheet.Sheet.IdClassroom == cl.Id {
			out.Sheets = append(out.Sheets, sheet)
		}
	}
	sort.Slice(out.Sheets, func(i, j int) bool { return out.Sheets[i].Sheet.Id < out.Sheets[j].Sheet.Id })
	return out
}

type SheetExt struct {
	Sheet     Sheet
	Exercices []ed.ExerciceHeader
}

func newSheetExt(sheet Sheet, links SheetExercices, exes map[ed.IdExercice]ed.ExerciceHeader) SheetExt {
	out := SheetExt{Sheet: sheet}
	links.ensureOrder()
	for _, link := range links {
		out.Exercices = append(out.Exercices, exes[link.IdExercice])
	}
	return out
}

func buildSheetExts(sheets Sheets, links SheetExercices, exes map[ed.IdExercice]ed.ExerciceHeader) map[IdSheet]SheetExt {
	out := make(map[IdSheet]SheetExt, len(sheets))
	m := links.ByIdSheet()
	for idSheet, v := range sheets {
		out[idSheet] = newSheetExt(v, m[idSheet], exes)
	}
	return out
}

func (l SheetExercices) ensureOrder() {
	sort.Slice(l, func(i, j int) bool { return l[i].Index < l[j].Index })
}

func updateSheetExercices(tx *sql.Tx, idSheet IdSheet, l []ed.IdExercice) error {
	links := make(SheetExercices, len(l))
	for i, idExe := range l {
		links[i] = SheetExercice{IdExercice: idExe, IdSheet: idSheet, Index: i}
	}

	_, err := DeleteSheetExercicesByIdSheets(tx, idSheet)
	if err != nil {
		return utils.SQLError(err)
	}

	err = InsertManySheetExercices(tx, links...)
	if err != nil {
		return utils.SQLError(err)
	}

	return nil
}

// Student API

type ExerciceProgressionHeader struct {
	Exercice       ed.Exercice `gomacro-extern:"editor:dart:../shared_gen.dart"`
	HasProgression bool
	// empty if HasProgression is false
	Progression  ed.ProgressionExt `gomacro-extern:"editor:dart:../shared_gen.dart"`
	Mark, Bareme int               // student mark / exercice total
}

// compute the student mark
func mark(questions ed.ExerciceQuestions, progression ed.ProgressionExt) int {
	questions.EnsureOrder()

	var out int
	for index, qu := range questions {
		results := progression.Questions[index]
		if results.Success() {
			out += qu.Bareme
		}
	}
	return out
}

// SheetProgression is the summary of the progression
// of one student for one sheet
type SheetProgression struct {
	Sheet     Sheet
	Exercices []ExerciceProgressionHeader
}

// sheetAndIndex is a key identidying an exercice in a sheet
type sheetAndIndex struct {
	IdSheet IdSheet
	Index   int
}

// assume links only concerne one student and use the unicity of (IdStudent, IdSheet, Index)
func (links StudentProgressions) bySheetAndIndex() (out map[sheetAndIndex]ed.IdProgression) {
	out = make(map[sheetAndIndex]ed.IdProgression)
	for _, link := range links {
		out[sheetAndIndex{IdSheet: link.IdSheet, Index: link.Index}] = link.IdProgression
	}
	return out
}

func (links SheetExercices) bySheetAndIndex() (out map[sheetAndIndex]ed.IdExercice) {
	out = make(map[sheetAndIndex]ed.IdExercice)
	for _, link := range links {
		out[sheetAndIndex{IdSheet: link.IdSheet, Index: link.Index}] = link.IdExercice
	}
	return out
}
