package homework

import (
	"database/sql"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/server/src/maths/questions"
	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	"github.com/benoitkugler/maths-online/server/src/pass"
	"github.com/benoitkugler/maths-online/server/src/sql/editor"
	"github.com/benoitkugler/maths-online/server/src/sql/events"
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

	db = tu.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../sql/editor/gen_create.sql", "../../sql/tasks/gen_create.sql",
		"../../sql/homework/gen_create.sql", "../../sql/reviews/gen_create.sql", "../../sql/events/gen_create.sql")

	_, err := teacher.Teacher{IsAdmin: true, FavoriteMatiere: teacher.Mathematiques}.Insert(db)
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

	l, err := ct.getSheets(userID, teacher.Mathematiques)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(l.Sheets) == 0)
	tu.Assert(t, len(l.Travaux) == 1) // one per classroom
	tu.Assert(t, len(l.Travaux[0].Travaux) == 0)

	sh_, err := ct.createSheet(userID)
	tu.AssertNoErr(t, err)
	sh := sh_.Sheet

	updated := ho.Sheet{}
	updated.Id = sh.Id
	updated.IdTeacher = userID
	updated.Level = string(editor.Seconde)
	updated.Matiere = teacher.Mathematiques
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

	l, err = ct.getSheets(userID, teacher.Francais)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(l.Sheets) == 0)

	l, err = ct.getSheets(userID, teacher.Mathematiques)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(l.Sheets) == 1)
	tu.Assert(t, len(l.Travaux) == 1)
	tu.Assert(t, len(l.Travaux[0].Travaux) == 0)

	tr, err := ct.assignSheetTo(CreateTravailWithIn{IdSheet: sh.Id, IdClassroom: class.Id}, userID)
	tu.AssertNoErr(t, err)

	l, err = ct.getSheets(userID, teacher.Mathematiques)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(l.Travaux[0].Travaux) == 1)

	out, err := ct.duplicateSheet(sh.Id, userID)
	tu.AssertNoErr(t, err)
	tu.Assert(t, out.Sheet.Id != sh.Id)

	l, err = ct.getSheets(userID, teacher.Mathematiques)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(l.Sheets) == 2)

	_, err = ct.copyTravailTo(CopyTravailIn{IdTravail: tr.Id, IdClassroom: class.Id}, userID)
	l, err = ct.getSheets(userID, teacher.Mathematiques)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(l.Travaux[0].Travaux) == 2)

	err = ct.deleteSheet(sh.Id, userID)
	tu.AssertNoErr(t, err)

	out2, err := ct.createTravail(class.Id, userID)
	tu.AssertNoErr(t, err)
	tu.Assert(t, out2.Sheet.Sheet.Anonymous.ID == out2.Travail.Id)

	l, err = ct.getSheets(userID, teacher.Mathematiques)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(l.Sheets) == 2)

	_, err = ct.copyTravailTo(CopyTravailIn{out2.Travail.Id, class.Id}, userID)
	tu.AssertNoErr(t, err)
	l, err = ct.getSheets(userID, teacher.Mathematiques)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(l.Sheets) == 3) // anonymous sheet is duplicated

	err = ct.deleteTravail(out2.Travail.Id, userID)
	tu.AssertNoErr(t, err)

	l, err = ct.getSheets(userID, teacher.Mathematiques)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(l.Sheets) == 2) // anonymous sheet deleted by cascade
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
	sh2_, err := ct.createSheet(userID)
	tu.AssertNoErr(t, err)
	sh1 := sh2_.Sheet
	sh2_, err = ct.createSheet(userID)
	tu.AssertNoErr(t, err)
	sh2 := sh2_.Sheet
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
	tr1, err := ct.assignSheetTo(CreateTravailWithIn{IdSheet: sh1.Id, IdClassroom: class.Id}, userID)
	tu.AssertNoErr(t, err)
	tr2, err := ct.assignSheetTo(CreateTravailWithIn{IdSheet: sh2.Id, IdClassroom: class.Id}, userID)
	tu.AssertNoErr(t, err)

	tr1.ShowAfter = ho.Time(time.Now().Add(-time.Second))
	err = ct.updateTravail(tr1, userID)
	tu.AssertNoErr(t, err)
	tr2.ShowAfter = ho.Time(time.Now().Add(-time.Second))
	err = ct.updateTravail(tr2, userID)
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
	sh_, err := ct.createSheet(sp.userID)
	tu.AssertNoErr(t, err)
	sh := sh_.Sheet
	task, err := ct.addExerciceTo(AddExerciceToTaskIn{IdSheet: sh.Id, IdExercice: sp.exe1.Id}, sp.userID)
	tu.AssertNoErr(t, err)

	tr, err := ct.assignSheetTo(CreateTravailWithIn{IdSheet: sh.Id, IdClassroom: class.Id}, sp.userID)
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
			ID:          task.IdWork,
			AnswerIndex: 0,
			Answer:      tasks.AnswerP{},
		},
		IdTravail: tr.Id,
	})
	tu.AssertNoErr(t, err)

	// the sheet is not expired, check that a new progression has been added
	tu.Assert(t, len(out.Ex.Progression.Questions[0]) == 1)
	taHeader := loadProgression(t, ct.db, student.Id, task.Id)
	tu.Assert(t, taHeader.HasProgression)
	tu.Assert(t, len(taHeader.Progression.Questions[0]) == 1)

	// check the events
	evL, err := events.SelectEventsByIdStudents(ct.db, student.Id)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(evL) == 1 && evL[0].Event == events.E_All_QuestionWrong)
	tu.Assert(t, len(out.Advance.Events) == 1)

	// now expire the sheet ...
	tr.Deadline = ho.Time(time.Now().Add(-time.Hour))
	err = ct.updateTravail(tr, sp.userID)
	tu.AssertNoErr(t, err)

	args := StudentEvaluateTaskIn{
		StudentID: studentKey.EncryptID(int64(student.Id)),
		IdTask:    task.Id,
		Ex: tasks.EvaluateWorkIn{
			ID:          task.IdWork,
			AnswerIndex: 0,
			Answer:      tasks.AnswerP{Answer: client.QuestionAnswersIn{Data: client.Answers{0: client.NumberAnswer{Value: 1}}}},
			Progression: out.Ex.Progression,
		},
		IdTravail: tr.Id,
	}
	out, err = ct.studentEvaluateTask(args)
	tu.AssertNoErr(t, err)

	// check the events
	tu.Assert(t, len(out.Advance.Events) == 2)
	tu.Assert(t, out.Advance.Events[0] == events.E_All_QuestionRight)
	tu.Assert(t, out.Advance.Events[1] == events.E_Homework_TaskDone)

	// ... and check that no new progression has been added,
	// despite the runtime progression correctly updated
	tu.Assert(t, len(out.Ex.Progression.Questions[0]) == 2)
	taHeader = loadProgression(t, ct.db, student.Id, task.Id)
	tu.Assert(t, len(taHeader.Progression.Questions[0]) == 1)

	// finally add an exception...
	exc := ho.TravailException{IdStudent: student.Id, IdTravail: tr.Id, Deadline: sql.NullTime{Valid: true, Time: time.Now().Add(time.Minute)}}
	err = exc.Insert(ct.db)
	tu.AssertNoErr(t, err)
	// ... and evaluate again
	out, err = ct.studentEvaluateTask(args)
	tu.AssertNoErr(t, err)
	// the sheet is not expired for this student, check that a new progression has been added
	taHeader = loadProgression(t, ct.db, student.Id, task.Id)
	tu.Assert(t, len(taHeader.Progression.Questions[0]) == 2)
}

func insertProgression(db *sql.DB, idTask ta.IdTask, idStudent teacher.IdStudent, questions []ta.QuestionHistory) error {
	links := make(ta.Progressions, len(questions))
	for i, qu := range questions {
		links[i] = ta.Progression{IdStudent: idStudent, IdTask: idTask, Index: int16(i), History: qu}
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
