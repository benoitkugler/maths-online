package exercice

import (
	"bytes"
	"database/sql"
	"os/exec"
	"reflect"
	"testing"

	"github.com/benoitkugler/maths-online/pass"
)

type logsDB = pass.DB

func getUserName() string {
	var buf bytes.Buffer
	cmd := exec.Command("whoami")
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	return buf.String()
}

// create a new database and add tables
func createDBDev() logsDB {
	err := exec.Command("createdb", "tmp_dev_test").Run()
	if err != nil {
		panic(err)
	}

	err = exec.Command("bash", "-c", "psql tmp_dev_test < create_gen.sql").Run()
	if err != nil {
		panic(err)
	}

	const userPassword = "dummy"
	return logsDB{
		Name:     "tmp_dev_test",
		Host:     "localhost",
		User:     getUserName(),
		Password: userPassword,
	}
}

func removeDBDev() {
	err := exec.Command("dropdb", "tmp_dev_test").Run()
	if err != nil {
		panic(err)
	}
}

func TestRoot(t *testing.T) {
	// create a DB shared by all tests
	logs := createDBDev()
	defer removeDBDev()

	db, err := logs.ConnectPostgres()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		t.Fatal(err)
	}

	// t.Run("CRUD for Exercice", func(t *testing.T) { testExercice(t, db) })
	t.Run("CRUD for Question", func(t *testing.T) { testQuestion(t, db) })
}

// func testExercice(t *testing.T, db *sql.DB) {
// 	exes, err := SelectAllExercices(db)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	L := len(exes)

// 	ex := randExercice()
// 	ex, err = ex.Insert(db)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	exes, err = SelectAllExercices(db)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(exes) != L+1 {
// 		t.Fatal()
// 	}

// 	_, err = DeleteExerciceById(db, ex.Id)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

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

	tags, err := SelectAllTags(db)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(tags, []string{"calcul", "seconde"}) {
		t.Fatal()
	}

	d, err := SelectQuestionByTags(db, "calcul")
	if err != nil {
		t.Fatal(err)
	}
	if len(d) != 1 {
		t.Fatal()
	}

	d, err = SelectQuestionByTags(db, "calcul", "seconde")
	if err != nil {
		t.Fatal(err)
	}
	if len(d) != 1 {
		t.Fatal()
	}

	d, err = SelectQuestionByTags(db, "calcul", "XXX")
	if err != nil {
		t.Fatal(err)
	}
	if len(d) != 0 {
		t.Fatal()
	}

	_, err = DeleteQuestionById(db, qu.Id)
	if err != nil {
		t.Fatal(err)
	}
}

// TODO: populate a "real" DB to support further testing
func setupDBDev(db *sql.DB) {
}
