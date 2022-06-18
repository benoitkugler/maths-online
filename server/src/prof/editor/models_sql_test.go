package editor

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/maths/questions"
	"github.com/benoitkugler/maths-online/prof/teacher"
	"github.com/benoitkugler/maths-online/utils/testutils"
)

func TestRoot(t *testing.T) {
	// create a DB shared by all tests
	db := testutils.CreateDBDev(t, "gen_create.sql")
	defer testutils.RemoveDBDev()
	defer db.Close()

	// t.Run("CRUD for Exercice", func(t *testing.T) { testExercice(t, db) })
	t.Run("CRUD for Question", func(t *testing.T) { testQuestion(t, db) })
	t.Run("Insert SignTable", func(t *testing.T) { testInsertSignTable(t, db) })
}

func testQuestion(t *testing.T, db *sql.DB) {
	questions, err := SelectAllQuestions(db)
	if err != nil {
		t.Fatal(err)
	}
	L := len(questions)

	qu := randQuestion()
	qu, err = qu.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	questions, err = SelectAllQuestions(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(questions) != L+1 {
		t.Fatal()
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	err = InsertManyQuestionTags(tx, QuestionTag{IdQuestion: qu.Id, Tag: "seconde"}, QuestionTag{IdQuestion: qu.Id, Tag: "calcul"})
	if err != nil {
		t.Fatal(err)
	}
	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}

	tags, err := SelectAllQuestionTags(db)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(tags.List(), []string{"calcul", "seconde"}) {
		t.Fatal()
	}

	_, err = DeleteQuestionById(db, qu.Id)
	if err != nil {
		t.Fatal(err)
	}
}

func testInsertSignTable(t *testing.T, db *sql.DB) {
	qu := randQuestion()
	qu.Page.Enonce = questions.Enonce{randque_SignTableBlock()}
	qu, err := qu.Insert(db)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadQuestions(t *testing.T) {
	db, err := testutils.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", testutils.DB, err)
		return
	}

	m, err := SelectAllQuestions(db)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Questions :", len(m))
}

func TestValidation(t *testing.T) {
	db, err := testutils.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", testutils.DB, err)
		return
	}

	qu, err := SelectAllQuestions(db)
	if err != nil {
		t.Fatal(err)
	}

	ti := time.Now()
	err = validateAllQuestions(qu)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Validated in :", time.Since(ti), "average :", time.Since(ti)/time.Duration(len(qu)))
}

func BenchmarkValidation(b *testing.B) {
	db, err := testutils.DB.ConnectPostgres()
	if err != nil {
		b.Skipf("DB %v not available : %s", testutils.DB, err)
		return
	}

	qu, err := SelectAllQuestions(db)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		validateAllQuestions(qu)
	}
}

func TestCRUD(t *testing.T) {
	db := testutils.CreateDBDev(t, "../teacher/gen_create.sql", "gen_create.sql")
	defer testutils.RemoveDBDev()
	defer db.Close()

	user, err := teacher.Teacher{}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	ex := randExercice()
	ex.IdTeacher = user.Id
	ex, err = ex.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	qu1, err := Question{IdTeacher: user.Id}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}
	qu2, err := Question{IdTeacher: user.Id}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	err = InsertManyExerciceQuestions(tx,
		ExerciceQuestion{IdExercice: ex.Id, IdQuestion: qu1.Id, Bareme: 4},
		ExerciceQuestion{IdExercice: ex.Id, IdQuestion: qu2.Id, Bareme: 5},
	)
	if err != nil {
		t.Fatal(err)
	}

	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}
}
