package editor

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	que "github.com/benoitkugler/maths-online/server/src/maths/questions"
	"github.com/benoitkugler/maths-online/server/src/sql/teacher"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func selectOneTeacher(t *testing.T, db *sql.DB) teacher.Teacher {
	teachers, err := teacher.SelectAllTeachers(db)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(teachers) != 0)
	return teachers[teachers.IDs()[0]]
}

func TestRoot(t *testing.T) {
	// create a DB shared by all tests
	db := tu.NewTestDB(t, "../teacher/gen_create.sql", "gen_create.sql")
	defer db.Remove()

	t.Run("CRUD for Question", func(t *testing.T) { testQuestion(t, db.DB) })
	t.Run("Insert many Questions", func(t *testing.T) { testInsertQuestions(t, db.DB) })
	t.Run("Insert SignTable", func(t *testing.T) { testInsertSignTable(t, db.DB) })
	t.Run("CRUD for Exercice", func(t *testing.T) { testCRUDExercice(t, db.DB) })
}

func testQuestion(t *testing.T, db *sql.DB) {
	questions, err := SelectAllQuestions(db)
	tu.AssertNoErr(t, err)

	L := len(questions)

	user, err := teacher.Teacher{}.Insert(db)
	tu.AssertNoErr(t, err)

	group := randQuestiongroup()
	group.IdTeacher = user.Id
	group, err = group.Insert(db)
	tu.AssertNoErr(t, err)

	qu := randQuestion()
	qu.IdGroup = group.Id.AsOptional()
	qu.NeedExercice = OptionalIdExercice{}
	qu, err = qu.Insert(db)
	tu.AssertNoErr(t, err)

	questions, err = SelectAllQuestions(db)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(questions) == L+1)

	tx, err := db.Begin()
	tu.AssertNoErr(t, err)

	err = InsertManyQuestiongroupTags(tx,
		QuestiongroupTag{IdQuestiongroup: group.Id, Tag: "CALCUL", Section: Chapter},
		QuestiongroupTag{IdQuestiongroup: group.Id, Tag: "SECONDE", Section: Level},
	)
	tu.AssertNoErr(t, err)

	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}

	tags, err := SelectAllQuestiongroupTags(db)
	tu.AssertNoErr(t, err)
	tu.Assert(t, reflect.DeepEqual(tags.Tags(), Tags{{"SECONDE", Level}, {"CALCUL", Chapter}}))

	_, err = DeleteQuestionById(db, qu.Id)
	tu.AssertNoErr(t, err)
}

func testInsertSignTable(t *testing.T, db *sql.DB) {
	user := selectOneTeacher(t, db)

	group := Questiongroup{IdTeacher: user.Id}
	group, err := group.Insert(db)
	tu.AssertNoErr(t, err)

	qu := randQuestion()
	qu.IdGroup = group.Id.AsOptional()
	qu.NeedExercice = OptionalIdExercice{}
	qu.Enonce = que.Enonce{randque_SignTableBlock()}
	qu, err = qu.Insert(db)
	tu.AssertNoErr(t, err)
}

func TestLoadQuestions(t *testing.T) {
	db, err := tu.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", tu.DB, err)
		return
	}

	m, err := SelectAllQuestions(db)
	tu.AssertNoErr(t, err)

	fmt.Println("Questions :", len(m))
}

func testCRUDExercice(t *testing.T, db *sql.DB) {
	user := selectOneTeacher(t, db)

	group := randExercicegroup()
	group.IdTeacher = user.Id
	group, err := group.Insert(db)
	tu.AssertNoErr(t, err)

	ex := randExercice()
	ex.IdGroup = group.Id
	ex, err = ex.Insert(db)
	tu.AssertNoErr(t, err)

	qu1, err := Question{NeedExercice: ex.Id.AsOptional()}.Insert(db)
	tu.AssertNoErr(t, err)

	qu2, err := Question{NeedExercice: ex.Id.AsOptional()}.Insert(db)
	tu.AssertNoErr(t, err)

	tx, err := db.Begin()
	tu.AssertNoErr(t, err)

	err = InsertManyExerciceQuestions(tx,
		ExerciceQuestion{IdExercice: ex.Id, IdQuestion: qu1.Id, Bareme: 4, Index: 0},
		ExerciceQuestion{IdExercice: ex.Id, IdQuestion: qu2.Id, Bareme: 5, Index: 1},
	)
	tu.AssertNoErr(t, err)

	err = tx.Commit()
	tu.AssertNoErr(t, err)
}

func testInsertQuestions(t *testing.T, db *sql.DB) {
	user := selectOneTeacher(t, db)

	group, err := Questiongroup{IdTeacher: user.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	tx, err := db.Begin()
	tu.AssertNoErr(t, err)

	for i := 0; i < 50; i++ {
		qu := randQuestion()
		qu.IdGroup = group.Id.AsOptional()
		qu.NeedExercice = OptionalIdExercice{}

		_, err = qu.Insert(tx)
		tu.AssertNoErr(t, err)
	}

	err = tx.Commit()
	tu.AssertNoErr(t, err)
}
