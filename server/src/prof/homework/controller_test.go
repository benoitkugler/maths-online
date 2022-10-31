package homework

import (
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/maths/questions"
	"github.com/benoitkugler/maths-online/pass"
	"github.com/benoitkugler/maths-online/sql/editor"
	ho "github.com/benoitkugler/maths-online/sql/homework"
	ta "github.com/benoitkugler/maths-online/sql/tasks"
	"github.com/benoitkugler/maths-online/sql/teacher"
	"github.com/benoitkugler/maths-online/tasks"
	tu "github.com/benoitkugler/maths-online/utils/testutils"
)

type sample struct {
	userID     uID
	class      teacher.Classroom
	exe1, exe2 editor.Exercice
	question   editor.Question
}

func setupDB(t *testing.T) (db tu.TestDB, out sample) {
	db = tu.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../sql/editor/gen_create.sql", "../../sql/tasks/gen_create.sql", "../../sql/homework/gen_create.sql")

	_, err := teacher.Teacher{IsAdmin: true}.Insert(db)
	tu.Assert(t, err == nil)
	out.userID = teacher.IdTeacher(1)

	out.class, err = teacher.Classroom{IdTeacher: out.userID, Name: "test"}.Insert(db)
	tu.Assert(t, err == nil)

	group, err := editor.Exercicegroup{IdTeacher: out.userID}.Insert(db)
	tu.Assert(t, err == nil)

	out.exe1, err = editor.Exercice{IdGroup: group.Id}.Insert(db)
	tu.Assert(t, err == nil)
	out.exe2, err = editor.Exercice{IdGroup: group.Id}.Insert(db)
	tu.Assert(t, err == nil)

	exe1Qu, err := editor.Question{NeedExercice: out.exe1.Id.AsOptional(), Page: questions.QuestionPage{
		Enonce: questions.Enonce{
			questions.NumberFieldBlock{Expression: "1"},
		},
	}}.Insert(db)
	tu.Assert(t, err == nil)
	tx, err := db.Begin()
	tu.Assert(t, err == nil)
	err = editor.InsertManyExerciceQuestions(tx, editor.ExerciceQuestion{IdExercice: out.exe1.Id, IdQuestion: exe1Qu.Id, Bareme: 2, Index: 0})
	tu.Assert(t, err == nil)
	err = tx.Commit()
	tu.Assert(t, err == nil)

	quGroup, err := editor.Questiongroup{IdTeacher: out.userID}.Insert(db)
	tu.Assert(t, err == nil)

	out.question, err = editor.Question{IdGroup: quGroup.Id.AsOptional()}.Insert(db)
	tu.Assert(t, err == nil)

	return db, out
}

func TestCRUDSheet(t *testing.T) {
	db, sample := setupDB(t)
	defer db.Remove()
	userID, class, exe1, exe2 := sample.userID, sample.class, sample.exe1, sample.exe2
	qu := sample.question
	ct := NewController(db.DB, teacher.Teacher{Id: userID}, pass.Encrypter{})

	l, err := ct.getSheets(userID)
	tu.Assert(t, err == nil)
	tu.Assert(t, len(l) == 1)
	tu.Assert(t, len(l[0].Sheets) == 0)

	sh, err := ct.createSheet(class.Id, userID)
	tu.Assert(t, err == nil)

	updated := ho.Sheet{}
	updated.Id = sh.Id
	updated.IdClassroom = class.Id
	err = ct.updateSheet(updated, userID)
	tu.Assert(t, err == nil)

	_, err = ct.addExerciceTo(AddExerciceToTaskIn{IdSheet: sh.Id, IdExercice: exe1.Id}, userID)
	tu.Assert(t, err == nil)
	_, err = ct.addExerciceTo(AddExerciceToTaskIn{IdSheet: sh.Id, IdExercice: exe2.Id}, userID)
	tu.Assert(t, err == nil)
	_, err = ct.addExerciceTo(AddExerciceToTaskIn{IdSheet: sh.Id, IdExercice: exe1.Id}, userID)
	tu.Assert(t, err == nil)

	_, err = ct.addMonoquestionTo(AddMonoquestionToTaskIn{IdSheet: sh.Id, IdQuestion: qu.Id}, userID)
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
	task1, err := ct.addExerciceTo(AddExerciceToTaskIn{IdSheet: sh1.Id, IdExercice: exe1.Id}, userID)
	tu.Assert(t, err == nil)
	_, err = ct.addExerciceTo(AddExerciceToTaskIn{IdSheet: sh1.Id, IdExercice: exe1.Id}, userID)
	tu.Assert(t, err == nil)
	_, err = ct.addExerciceTo(AddExerciceToTaskIn{IdSheet: sh1.Id, IdExercice: exe2.Id}, userID)
	tu.Assert(t, err == nil)
	_, err = ct.addExerciceTo(AddExerciceToTaskIn{IdSheet: sh2.Id, IdExercice: exe1.Id}, userID)
	tu.Assert(t, err == nil)

	sheets, err = ct.getStudentSheets(student.Id)
	tu.Assert(t, err == nil)
	tu.Assert(t, len(sheets) == 2)
	tu.Assert(t, len(sheets[0].Tasks) == 3)
	tu.Assert(t, len(sheets[1].Tasks) == 1)

	err = ct.removeTask(task1.Id, userID)
	tu.Assert(t, err == nil)

	_, err = ct.addExerciceTo(AddExerciceToTaskIn{IdSheet: sh1.Id, IdExercice: exe1.Id}, userID)
	tu.Assert(t, err == nil)
}

func loadProgression(t *testing.T, db editor.DB, student teacher.IdStudent, taskID ta.IdTask) tasks.TaskProgressionHeader {
	m, err := tasks.LoadTasksProgression(db, student, []ta.IdTask{taskID})
	tu.Assert(t, err == nil)
	return m[taskID]
}

func TestEvaluateTask(t *testing.T) {
	db, sp := setupDB(t)
	defer db.Remove()
	class := sp.class
	studentKey := pass.Encrypter{}
	ct := NewController(db.DB, teacher.Teacher{Id: sp.userID}, studentKey)

	// setup the sheet and exercices
	sh, err := ct.createSheet(class.Id, sp.userID)
	tu.Assert(t, err == nil)
	sh.Activated = true
	sh.Deadline = ho.Time(time.Now().Add(time.Hour))
	err = ct.updateSheet(sh, sp.userID)
	tu.Assert(t, err == nil)
	task, err := ct.addExerciceTo(AddExerciceToTaskIn{IdSheet: sh.Id, IdExercice: sp.exe1.Id}, sp.userID)
	tu.Assert(t, err == nil)

	// setup a student
	student, err := teacher.Student{IdClassroom: sp.class.Id}.Insert(ct.db)
	tu.Assert(t, err == nil)

	out, err := ct.studentEvaluateTask(StudentEvaluateTaskIn{
		StudentID: studentKey.EncryptID(int64(student.Id)),
		IdTask:    task.Id,
		Ex: tasks.EvaluateWorkIn{
			ID: task.IdWork,
			Answers: map[int]tasks.Answer{
				0: {},
			},
		},
	})
	tu.Assert(t, err == nil)

	// the sheet is not expired, check that a new progression has been added
	if len(out.Ex.Progression.Questions[0]) != 1 {
		t.Fatal()
	}
	taHeader := loadProgression(t, ct.db, student.Id, task.Id)
	if !taHeader.HasProgression {
		t.Fatal()
	}

	// now expire the sheet ...
	sh.Deadline = ho.Time(time.Now().Add(-time.Hour))
	err = ct.updateSheet(sh, sp.userID)
	tu.Assert(t, err == nil)

	out, err = ct.studentEvaluateTask(StudentEvaluateTaskIn{
		StudentID: studentKey.EncryptID(int64(student.Id)),
		IdTask:    task.Id,
		Ex: tasks.EvaluateWorkIn{
			ID: task.IdWork,
			Answers: map[int]tasks.Answer{
				0: {},
			},
			Progression: out.Ex.Progression,
		},
	})
	tu.Assert(t, err == nil)

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
