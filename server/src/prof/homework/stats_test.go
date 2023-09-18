package homework

import (
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/server/src/pass"
	ho "github.com/benoitkugler/maths-online/server/src/sql/homework"
	ta "github.com/benoitkugler/maths-online/server/src/sql/tasks"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func TestGetMarks(t *testing.T) {
	db, sp := setupDB(t)
	defer db.Remove()
	class := sp.class
	studentKey := pass.Encrypter{}
	ct := NewController(db.DB, teacher.Teacher{Id: sp.userID}, studentKey)

	// setup the sheet and exercices
	sh, err := ct.createSheet(sp.userID)
	tu.AssertNoErr(t, err)
	task1, err := ct.addExerciceTo(AddExerciceToTaskIn{IdSheet: sh.Id, IdExercice: sp.exe1.Id}, sp.userID)
	tu.AssertNoErr(t, err)
	task2, err := ct.addMonoquestionTo(AddMonoquestionToTaskIn{IdSheet: sh.Id, IdQuestion: sp.question.Id}, sp.userID)
	tu.AssertNoErr(t, err)

	tr, err := ct.assignSheetTo(CreateTravailWithIn{IdSheet: sh.Id, IdClassroom: class.Id}, sp.userID)
	tu.AssertNoErr(t, err)
	tr.Noted = true
	tr.Deadline = ho.Time(time.Now().Add(time.Hour))
	err = ct.updateTravail(tr, sp.userID)
	tu.AssertNoErr(t, err)

	out, err := ct.getMarks(HowemorkMarksIn{
		IdClassroom: class.Id,
		IdTravaux:   []ho.IdTravail{tr.Id},
	}, sp.userID)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(out.Students) == 0)

	// add students and progressions
	student1, err := teacher.Student{IdClassroom: sp.class.Id}.Insert(ct.db)
	tu.AssertNoErr(t, err)
	student2, err := teacher.Student{IdClassroom: sp.class.Id}.Insert(ct.db)
	tu.AssertNoErr(t, err)

	// task1 has one question
	err = insertProgression(ct.db, task1.Id, student1.Id, []ta.QuestionHistory{
		{false, false, true},
	})
	tu.AssertNoErr(t, err)
	// do not create progression for student 2

	// task2 has 3 questions
	err = insertProgression(ct.db, task2.Id, student1.Id, []ta.QuestionHistory{
		{false, false, true},
		{true},
		{},
	})
	tu.AssertNoErr(t, err)
	err = insertProgression(ct.db, task2.Id, student2.Id, []ta.QuestionHistory{
		{false, false, false},
		{false, false, true},
		{false, true, false},
	})
	tu.AssertNoErr(t, err)

	err = ho.InsertTravailException(ct.db, ho.TravailException{
		IdStudent:     student2.Id,
		IdTravail:     tr.Id,
		IgnoreForMark: true,
	})
	tu.AssertNoErr(t, err)

	out, err = ct.getMarks(HowemorkMarksIn{
		IdClassroom: class.Id,
		IdTravaux:   []ho.IdTravail{tr.Id},
	}, sp.userID)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(out.Students) == 2)
	// student1 : 4/5 => 16/20
	// student2 : 2/5 => 8 /20
	ma := out.Marks[tr.Id]
	tu.Assert(t, ma.Marks[student1.Id].Mark == 16)
	tu.Assert(t, ma.Marks[student1.Id].NbTries == 7)
	tu.Assert(t, ma.Marks[student2.Id].Mark == 8)
	tu.Assert(t, ma.Marks[student2.Id].NbTries == 9)
	tu.Assert(t, ma.Marks[student2.Id].Dispensed)
}

func TestGetStats(t *testing.T) {
	db, err := tu.DB.ConnectPostgres()
	if err != nil {
		t.Skip(err)
	}

	ct := NewController(db, teacher.Teacher{Id: 1}, pass.Encrypter{})
	out, err := ct.getMarks(HowemorkMarksIn{IdClassroom: 27, IdTravaux: []ho.IdTravail{77, 79, 82, 86, 87, 88}}, 4)
	tu.AssertNoErr(t, err)

	_ = out.Marks
}
