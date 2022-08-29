package editor

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/benoitkugler/maths-online/maths/questions"
	"github.com/benoitkugler/maths-online/sql/teacher"
	tu "github.com/benoitkugler/maths-online/utils/testutils"
)

func selectOneTeacher(t *testing.T, db *sql.DB) teacher.Teacher {
	teachers, err := teacher.SelectAllTeachers(db)
	tu.Assert(t, err == nil)
	tu.Assert(t, len(teachers) != 0)
	return teachers[teachers.IDs()[0]]
}

func TestRoot(t *testing.T) {
	// create a DB shared by all tests
	db := tu.NewTestDB(t, "../teacher/gen_create.sql", "gen_create.sql")
	defer db.Remove()

	t.Run("CRUD for Question", func(t *testing.T) { testQuestion(t, db.DB) })
	t.Run("Insert SignTable", func(t *testing.T) { testInsertSignTable(t, db.DB) })
	t.Run("CRUD for Exercice", func(t *testing.T) { testCRUDExercice(t, db.DB) })
}

func testQuestion(t *testing.T, db *sql.DB) {
	questions, err := SelectAllQuestions(db)
	tu.Assert(t, err == nil)

	L := len(questions)

	user, err := teacher.Teacher{}.Insert(db)
	tu.Assert(t, err == nil)

	group := randQuestiongroup()
	group.IdTeacher = user.Id
	group, err = group.Insert(db)
	tu.Assert(t, err == nil)

	qu := randQuestion()
	qu.IdGroup = group.Id.AsOptional()
	qu, err = qu.Insert(db)
	tu.Assert(t, err == nil)

	questions, err = SelectAllQuestions(db)
	tu.Assert(t, err == nil)
	tu.Assert(t, len(questions) == L+1)

	tx, err := db.Begin()
	tu.Assert(t, err == nil)

	err = InsertManyQuestiongroupTags(tx,
		QuestiongroupTag{IdQuestiongroup: group.Id, Tag: "seconde"},
		QuestiongroupTag{IdQuestiongroup: group.Id, Tag: "calcul"},
	)
	tu.Assert(t, err == nil)

	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}

	tags, err := SelectAllQuestiongroupTags(db)
	tu.Assert(t, err == nil)
	tu.Assert(t, reflect.DeepEqual(tags.List(), []string{"calcul", "seconde"}))

	_, err = DeleteQuestionById(db, qu.Id)
	tu.Assert(t, err == nil)
}

func testInsertSignTable(t *testing.T, db *sql.DB) {
	user := selectOneTeacher(t, db)

	group := Questiongroup{IdTeacher: user.Id}
	group, err := group.Insert(db)
	tu.Assert(t, err == nil)

	qu := randQuestion()
	qu.IdGroup = group.Id.AsOptional()
	qu.NeedExercice = OptionalIdExercice{}
	qu.Page.Enonce = questions.Enonce{randque_SignTableBlock()}
	qu, err = qu.Insert(db)
	tu.Assert(t, err == nil)
}

func TestLoadQuestions(t *testing.T) {
	db, err := tu.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", tu.DB, err)
		return
	}

	m, err := SelectAllQuestions(db)
	tu.Assert(t, err == nil)

	fmt.Println("Questions :", len(m))
}

func testCRUDExercice(t *testing.T, db *sql.DB) {
	user := selectOneTeacher(t, db)

	group := randExercicegroup()
	group.IdTeacher = user.Id
	group, err := group.Insert(db)
	tu.Assert(t, err == nil)

	ex := randExercice()
	ex.IdGroup = group.Id
	ex, err = ex.Insert(db)
	tu.Assert(t, err == nil)

	qu1, err := Question{NeedExercice: ex.Id.AsOptional()}.Insert(db)
	tu.Assert(t, err == nil)

	qu2, err := Question{NeedExercice: ex.Id.AsOptional()}.Insert(db)
	tu.Assert(t, err == nil)

	tx, err := db.Begin()
	tu.Assert(t, err == nil)

	err = InsertManyExerciceQuestions(tx,
		ExerciceQuestion{IdExercice: ex.Id, IdQuestion: qu1.Id, Bareme: 4, Index: 0},
		ExerciceQuestion{IdExercice: ex.Id, IdQuestion: qu2.Id, Bareme: 5, Index: 1},
	)
	tu.Assert(t, err == nil)

	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}
}
