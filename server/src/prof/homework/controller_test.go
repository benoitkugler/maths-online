package homework

import (
	"database/sql"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/server/src/maths/questions"
	"github.com/benoitkugler/maths-online/server/src/pass"
	"github.com/benoitkugler/maths-online/server/src/sql/editor"
	ho "github.com/benoitkugler/maths-online/server/src/sql/homework"
	ta "github.com/benoitkugler/maths-online/server/src/sql/tasks"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	"github.com/benoitkugler/maths-online/server/src/tasks"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

type sample struct {
	userID   uID
	class    teacher.Classroom
	exe1     editor.Exercice // with one question, / 2
	exe2     editor.Exercice
	question editor.Question // used as mono question
}

func setupDB(t *testing.T) (db tu.TestDB, out sample) {
	t.Helper()

	db = tu.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../sql/editor/gen_create.sql", "../../sql/tasks/gen_create.sql", "../../sql/homework/gen_create.sql")

	_, err := teacher.Teacher{IsAdmin: true}.Insert(db)
	tu.AssertNoErr(t, err)
	out.userID = teacher.IdTeacher(1)

	out.class, err = teacher.Classroom{IdTeacher: out.userID, Name: "test"}.Insert(db)
	tu.AssertNoErr(t, err)

	group, err := editor.Exercicegroup{IdTeacher: out.userID}.Insert(db)
	tu.AssertNoErr(t, err)

	out.exe1, err = editor.Exercice{IdGroup: group.Id}.Insert(db)
	tu.AssertNoErr(t, err)
	out.exe2, err = editor.Exercice{IdGroup: group.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	exe1Qu, err := editor.Question{NeedExercice: out.exe1.Id.AsOptional(), Enonce: questions.Enonce{
		questions.NumberFieldBlock{Expression: "1"},
	}}.Insert(db)
	tu.AssertNoErr(t, err)
	tx, err := db.Begin()
	tu.AssertNoErr(t, err)
	err = editor.InsertManyExerciceQuestions(tx, editor.ExerciceQuestion{IdExercice: out.exe1.Id, IdQuestion: exe1Qu.Id, Bareme: 2, Index: 0})
	tu.AssertNoErr(t, err)
	err = tx.Commit()
	tu.AssertNoErr(t, err)

	quGroup, err := editor.Questiongroup{IdTeacher: out.userID}.Insert(db)
	tu.AssertNoErr(t, err)

	out.question, err = editor.Question{IdGroup: quGroup.Id.AsOptional()}.Insert(db)
	tu.AssertNoErr(t, err)

	return db, out
}

func TestCRUDSheet(t *testing.T) {
	db, sample := setupDB(t)
	defer db.Remove()
	userID, class, exe1, exe2 := sample.userID, sample.class, sample.exe1, sample.exe2
	qu := sample.question
	ct := NewController(db.DB, teacher.Teacher{Id: userID}, pass.Encrypter{})

	l, err := ct.getSheets(userID)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(l.Sheets) == 0)
	tu.Assert(t, len(l.Travaux) == 1) // one per classroom
	tu.Assert(t, len(l.Travaux[0].Travaux) == 0)

	sh, err := ct.createSheet(userID)
	tu.AssertNoErr(t, err)

	updated := ho.Sheet{}
	updated.Id = sh.Id
	updated.IdTeacher = userID
	updated.Level = string(editor.Seconde)
	err = ct.updateSheet(updated, userID)
	tu.AssertNoErr(t, err)

	_, err = ct.addExerciceTo(AddExerciceToTaskIn{IdSheet: sh.Id, IdExercice: exe1.Id}, userID)
	tu.AssertNoErr(t, err)
	_, err = ct.addExerciceTo(AddExerciceToTaskIn{IdSheet: sh.Id, IdExercice: exe2.Id}, userID)
	tu.AssertNoErr(t, err)
	_, err = ct.addExerciceTo(AddExerciceToTaskIn{IdSheet: sh.Id, IdExercice: exe1.Id}, userID)
	tu.AssertNoErr(t, err)

	_, err = ct.addMonoquestionTo(AddMonoquestionToTaskIn{IdSheet: sh.Id, IdQuestion: qu.Id}, userID)
	tu.AssertNoErr(t, err)

	l, err = ct.getSheets(userID)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(l.Sheets) == 1)
	tu.Assert(t, len(l.Travaux) == 1)
	tu.Assert(t, len(l.Travaux[0].Travaux) == 0)

	tr, err := ct.assignSheetTo(CreateTravailIn{IdSheet: sh.Id, IdClassroom: class.Id}, userID)
	tu.AssertNoErr(t, err)

	l, err = ct.getSheets(userID)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(l.Travaux[0].Travaux) == 1)

	out, err := ct.duplicateSheet(CopySheetIn{IdSheet: sh.Id}, userID)
	tu.AssertNoErr(t, err)
	tu.Assert(t, out.Sheet.Id != sh.Id)

	l, err = ct.getSheets(userID)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(l.Sheets) == 2)

	_, err = ct.copyTravailTo(CopyTravailIn{IdTravail: tr.Id, IdClassroom: class.Id}, userID)
	l, err = ct.getSheets(userID)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(l.Travaux[0].Travaux) == 2)

	err = ct.deleteSheet(sh.Id, userID)
	tu.AssertNoErr(t, err)
}

func TestStudentSheets(t *testing.T) {
	db, sample := setupDB(t)
	defer db.Remove()
	userID, class, exe1, exe2 := sample.userID, sample.class, sample.exe1, sample.exe2
	ct := NewController(db.DB, teacher.Teacher{Id: userID}, pass.Encrypter{})

	student, err := teacher.Student{IdClassroom: class.Id}.Insert(ct.db)
	tu.AssertNoErr(t, err)

	sheets, err := ct.getStudentSheets(student.Id, true)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(sheets) == 0)

	// create sheets with exercices...
	sh1, err := ct.createSheet(userID)
	tu.AssertNoErr(t, err)
	sh2, err := ct.createSheet(userID)
	tu.AssertNoErr(t, err)
	_, err = ct.createSheet(userID)
	tu.AssertNoErr(t, err)

	task1, err := ct.addExerciceTo(AddExerciceToTaskIn{IdSheet: sh1.Id, IdExercice: exe1.Id}, userID)
	tu.AssertNoErr(t, err)
	_, err = ct.addExerciceTo(AddExerciceToTaskIn{IdSheet: sh1.Id, IdExercice: exe1.Id}, userID)
	tu.AssertNoErr(t, err)
	_, err = ct.addExerciceTo(AddExerciceToTaskIn{IdSheet: sh1.Id, IdExercice: exe2.Id}, userID)
	tu.AssertNoErr(t, err)
	_, err = ct.addExerciceTo(AddExerciceToTaskIn{IdSheet: sh2.Id, IdExercice: exe1.Id}, userID)
	tu.AssertNoErr(t, err)

	// open sheet 1 and 2 ...
	_, err = ct.assignSheetTo(CreateTravailIn{IdSheet: sh1.Id, IdClassroom: class.Id}, userID)
	tu.AssertNoErr(t, err)
	_, err = ct.assignSheetTo(CreateTravailIn{IdSheet: sh2.Id, IdClassroom: class.Id}, userID)
	tu.AssertNoErr(t, err)

	sheets, err = ct.getStudentSheets(student.Id, true)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(sheets) == 2)
	tu.Assert(t, len(sheets[0].Tasks) == 3)
	tu.Assert(t, len(sheets[1].Tasks) == 1)

	// travaux are noted by default
	sheets, err = ct.getStudentSheets(student.Id, false)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(sheets) == 0)

	err = ct.removeTask(task1.Id, userID)
	tu.AssertNoErr(t, err)

	_, err = ct.addExerciceTo(AddExerciceToTaskIn{IdSheet: sh1.Id, IdExercice: exe1.Id}, userID)
	tu.AssertNoErr(t, err)
}

func loadProgression(t *testing.T, db editor.DB, student teacher.IdStudent, taskID ta.IdTask) tasks.TaskProgressionHeader {
	m, err := tasks.LoadTasksProgression(db, student, []ta.IdTask{taskID})
	tu.AssertNoErr(t, err)
	return m[taskID]
}

func TestEvaluateTask(t *testing.T) {
	db, sp := setupDB(t)
	defer db.Remove()
	class := sp.class
	studentKey := pass.Encrypter{}
	ct := NewController(db.DB, teacher.Teacher{Id: sp.userID}, studentKey)

	// setup the sheet and exercices
	sh, err := ct.createSheet(sp.userID)
	tu.AssertNoErr(t, err)
	task, err := ct.addExerciceTo(AddExerciceToTaskIn{IdSheet: sh.Id, IdExercice: sp.exe1.Id}, sp.userID)
	tu.AssertNoErr(t, err)

	tr, err := ct.assignSheetTo(CreateTravailIn{IdSheet: sh.Id, IdClassroom: class.Id}, sp.userID)
	tu.AssertNoErr(t, err)

	tr.Noted = true
	tr.Deadline = ho.Time(time.Now().Add(time.Hour))
	err = ct.updateTravail(tr, sp.userID)
	tu.AssertNoErr(t, err)

	// setup a student
	student, err := teacher.Student{IdClassroom: sp.class.Id}.Insert(ct.db)
	tu.AssertNoErr(t, err)

	out, err := ct.studentEvaluateTask(StudentEvaluateTaskIn{
		StudentID: studentKey.EncryptID(int64(student.Id)),
		IdTask:    task.Id,
		Ex: tasks.EvaluateWorkIn{
			ID: task.IdWork,
			Answers: map[int]tasks.AnswerP{
				0: {},
			},
		},
		IdTravail: tr.Id,
	})
	tu.AssertNoErr(t, err)

	// the sheet is not expired, check that a new progression has been added
	if len(out.Ex.Progression.Questions[0]) != 1 {
		t.Fatal()
	}
	taHeader := loadProgression(t, ct.db, student.Id, task.Id)
	if !taHeader.HasProgression {
		t.Fatal()
	}

	// now expire the sheet ...
	tr.Deadline = ho.Time(time.Now().Add(-time.Hour))
	err = ct.updateTravail(tr, sp.userID)
	tu.AssertNoErr(t, err)

	out, err = ct.studentEvaluateTask(StudentEvaluateTaskIn{
		StudentID: studentKey.EncryptID(int64(student.Id)),
		IdTask:    task.Id,
		Ex: tasks.EvaluateWorkIn{
			ID: task.IdWork,
			Answers: map[int]tasks.AnswerP{
				0: {},
			},
			Progression: out.Ex.Progression,
		},
		IdTravail: tr.Id,
	})
	tu.AssertNoErr(t, err)

	// ... and check that no new progression has been added,
	// despite the runtime progression correctly updated
	if len(out.Ex.Progression.Questions[0]) != 2 {
		t.Fatal()
	}
	taHeader = loadProgression(t, ct.db, student.Id, task.Id)
	if len(taHeader.Progression.Questions[0]) != 1 {
		t.Fatal()
	}
}

func insertProgression(db *sql.DB, idTask ta.IdTask, idStudent teacher.IdStudent, questions []ta.QuestionHistory) error {
	links := make(ta.Progressions, len(questions))
	for i, qu := range questions {
		links[i] = ta.Progression{IdStudent: idStudent, IdTask: idTask, Index: i, History: qu}
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	err = ta.InsertManyProgressions(tx, links...)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

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

	tr, err := ct.assignSheetTo(CreateTravailIn{IdSheet: sh.Id, IdClassroom: class.Id}, sp.userID)
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

	out, err = ct.getMarks(HowemorkMarksIn{
		IdClassroom: class.Id,
		IdTravaux:   []ho.IdTravail{tr.Id},
	}, sp.userID)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(out.Students) == 2)
	// student1 : 4/5 => 16/20
	// student2 : 2/5 => 8 /20
	ma := out.Marks[tr.Id]
	tu.Assert(t, ma[student1.Id] == 16)
	tu.Assert(t, ma[student2.Id] == 8)
}
