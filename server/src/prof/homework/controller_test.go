package homework

import (
	"testing"

	"github.com/benoitkugler/maths-online/pass"
	"github.com/benoitkugler/maths-online/prof/editor"
	"github.com/benoitkugler/maths-online/prof/teacher"
	tu "github.com/benoitkugler/maths-online/utils/testutils"
)

type sample struct {
	userID     uID
	class      teacher.Classroom
	exe1, exe2 editor.Exercice
}

func setupDB(t *testing.T) (db tu.TestDB, out sample) {
	db = tu.NewTestDB(t, "../teacher/gen_create.sql", "../editor/gen_create.sql", "../../tasks/gen_create.sql", "gen_create.sql")

	_, err := teacher.Teacher{IsAdmin: true}.Insert(db)
	tu.Assert(t, err == nil)
	out.userID = teacher.IdTeacher(1)

	out.class, err = teacher.Classroom{IdTeacher: out.userID, Name: "test"}.Insert(db)
	tu.Assert(t, err == nil)

	out.exe1, err = editor.Exercice{IdTeacher: out.userID}.Insert(db)
	tu.Assert(t, err == nil)
	out.exe2, err = editor.Exercice{IdTeacher: out.userID}.Insert(db)
	tu.Assert(t, err == nil)

	return db, out
}

func TestCRUDSheet(t *testing.T) {
	db, sample := setupDB(t)
	defer db.Remove()
	userID, class, exe1, exe2 := sample.userID, sample.class, sample.exe1, sample.exe2
	ct := NewController(db.DB, teacher.Teacher{Id: userID}, pass.Encrypter{})

	l, err := ct.getSheets(userID)
	tu.Assert(t, err == nil)
	tu.Assert(t, len(l) == 1)
	tu.Assert(t, len(l[0].Sheets) == 0)

	sh, err := ct.createSheet(class.Id, userID)
	tu.Assert(t, err == nil)

	updated := randSheet()
	updated.Id = sh.Id
	updated.IdClassroom = class.Id
	err = ct.updateSheet(updated, userID)
	tu.Assert(t, err == nil)

	_, err = ct.addTaskTo(AddTaskIn{IdSheet: sh.Id, IdExercice: exe1.Id}, userID)
	tu.Assert(t, err == nil)
	_, err = ct.addTaskTo(AddTaskIn{IdSheet: sh.Id, IdExercice: exe2.Id}, userID)
	tu.Assert(t, err == nil)
	_, err = ct.addTaskTo(AddTaskIn{IdSheet: sh.Id, IdExercice: exe1.Id}, userID)
	tu.Assert(t, err == nil)

	l, err = ct.getSheets(userID)
	tu.Assert(t, err == nil)
	tu.Assert(t, len(l) == 1)
	tu.Assert(t, len(l[0].Sheets) == 1)

	out, err := ct.copySheetTo(CopySheetIn{IdSheet: sh.Id, IdClassroom: class.Id}, userID)
	tu.Assert(t, err == nil)
	tu.Assert(t, out.Sheet.Id != sh.Id)

	err = ct.deleteSheet(sh.Id, userID)
	tu.Assert(t, err == nil)
}

func TestStudentSheets(t *testing.T) {
	db, sample := setupDB(t)
	defer db.Remove()
	userID, class, exe1, exe2 := sample.userID, sample.class, sample.exe1, sample.exe2
	ct := NewController(db.DB, teacher.Teacher{Id: userID}, pass.Encrypter{})

	student, err := teacher.Student{IdClassroom: class.Id}.Insert(ct.db)
	tu.Assert(t, err == nil)

	sheets, err := ct.getStudentSheets(student.Id)
	tu.Assert(t, err == nil)
	tu.Assert(t, len(sheets) == 0)

	// create sheets with exercices...
	sh1, err := ct.createSheet(class.Id, userID)
	tu.Assert(t, err == nil)
	sh2, err := ct.createSheet(class.Id, userID)
	tu.Assert(t, err == nil)
	_, err = ct.createSheet(class.Id, userID)
	tu.Assert(t, err == nil)

	// open sheet 1 and 2 ...
	sh1.Activated, sh2.Activated = true, true
	err = ct.updateSheet(sh1, userID)
	tu.Assert(t, err == nil)
	err = ct.updateSheet(sh2, userID)
	tu.Assert(t, err == nil)
	// ... and add exercices
	_, err = ct.addTaskTo(AddTaskIn{IdSheet: sh1.Id, IdExercice: exe1.Id}, userID)
	tu.Assert(t, err == nil)
	_, err = ct.addTaskTo(AddTaskIn{IdSheet: sh1.Id, IdExercice: exe1.Id}, userID)
	tu.Assert(t, err == nil)
	_, err = ct.addTaskTo(AddTaskIn{IdSheet: sh1.Id, IdExercice: exe2.Id}, userID)
	tu.Assert(t, err == nil)
	_, err = ct.addTaskTo(AddTaskIn{IdSheet: sh2.Id, IdExercice: exe1.Id}, userID)
	tu.Assert(t, err == nil)

	sheets, err = ct.getStudentSheets(student.Id)
	tu.Assert(t, err == nil)
	tu.Assert(t, len(sheets) == 2)
	tu.Assert(t, len(sheets[0].Tasks) == 3)
	tu.Assert(t, len(sheets[1].Tasks) == 1)
}
