package homework

import (
	"testing"

	"github.com/benoitkugler/maths-online/prof/editor"
	"github.com/benoitkugler/maths-online/prof/teacher"
	tu "github.com/benoitkugler/maths-online/utils/testutils"
)

func TestCRUDSheet(t *testing.T) {
	db := tu.NewTestDB(t, "../teacher/gen_create.sql", "../editor/gen_create.sql", "gen_create.sql")
	defer db.Remove()

	_, err := teacher.Teacher{IsAdmin: true}.Insert(db)
	tu.Assert(t, err == nil)
	userID := teacher.IdTeacher(1)

	class, err := teacher.Classroom{IdTeacher: userID, Name: "test"}.Insert(db)
	tu.Assert(t, err == nil)

	exe1, err := editor.Exercice{IdTeacher: userID}.Insert(db)
	tu.Assert(t, err == nil)
	exe2, err := editor.Exercice{IdTeacher: userID}.Insert(db)
	tu.Assert(t, err == nil)

	ct := NewController(db.DB, teacher.Teacher{Id: userID})
	l, err := ct.getSheets(userID)
	tu.Assert(t, err == nil)
	tu.Assert(t, len(l) == 1)
	tu.Assert(t, len(l[0].Sheets) == 0)

	sh, err := ct.createSheet(class.Id, userID)
	tu.Assert(t, err == nil)

	updated := randSheet()
	updated.Id = sh.Id
	updated.IdClassroom = class.Id
	err = ct.updateSheet(UpdateSheetIn{
		Sheet:     updated,
		Exercices: []editor.IdExercice{exe1.Id, exe2.Id, exe1.Id},
	}, userID)
	tu.Assert(t, err == nil)

	l, err = ct.getSheets(userID)
	tu.Assert(t, err == nil)
	tu.Assert(t, len(l) == 1)
	tu.Assert(t, len(l[0].Sheets) == 1)

	err = ct.deleteSheet(sh.Id, userID)
	tu.Assert(t, err == nil)
}
