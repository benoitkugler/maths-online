// Package homework implements an activity for student
// consisting in personal, at home training on exercices given by the
// teacher.
package homework

import (
	"database/sql"
	"sort"

	"github.com/benoitkugler/maths-online/prof/editor"
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
	Exercices []editor.Exercice
}

func buildSheetExts(sheets Sheets, links SheetExercices, exes editor.Exercices) map[IdSheet]SheetExt {
	out := make(map[IdSheet]SheetExt, len(sheets))
	m := links.ByIdSheet()
	for idSheet, v := range sheets {
		ext := SheetExt{Sheet: v}
		links := m[idSheet]
		links.ensureOrder()
		for _, link := range links {
			ext.Exercices = append(ext.Exercices, exes[link.IdExercice])
		}
		out[idSheet] = ext
	}
	return out
}

func (l SheetExercices) ensureOrder() {
	sort.Slice(l, func(i, j int) bool { return l[i].Index < l[j].Index })
}

func updateSheetExercices(tx *sql.Tx, idSheet IdSheet, l []editor.IdExercice) error {
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
