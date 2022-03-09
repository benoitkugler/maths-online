package exercice

import (
	"bytes"
	"database/sql"
	"fmt"
	"os/exec"
	"testing"
)

type logsDB struct {
	Host     string
	User     string
	Password string
	Name     string // of the database
	Port     int    // default to 5432
}

func ConnectDB(credences logsDB) (*sql.DB, error) {
	port := credences.Port
	if port == 0 {
		port = 5432
	}
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		credences.Host, port, credences.User, credences.Password, credences.Name)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("connexion DB : %s", err)
	}
	return db, nil
}

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

	return logsDB{
		Name:     "tmp_dev_test",
		Host:     "localhost",
		User:     getUserName(),
		Password: "dummy",
	}
}

func removeDBDev() {
	err := exec.Command("dropdb", "tmp_dev_test").Run()
	if err != nil {
		panic(err)
	}
}

func TestRoot(t *testing.T) {
	logs := createDBDev()
	defer removeDBDev()

	t.Run("basic CRUD", func(t *testing.T) { testSQL(t, logs) })
}

func testSQL(t *testing.T, logs logsDB) {
	db, err := ConnectDB(logs)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		t.Fatal(err)
	}

	ex := randExercice()
	ex, err = ex.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	qu1, qu2 := randQuestion(), randQuestion()
	qu1.IdExercice = ex.Id
	qu2.IdExercice = ex.Id
	err = InsertManyQuestions(tx, qu1, qu2)
	if err != nil {
		t.Fatal(err)
	}

	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}

	exQ, err := SelectExerciceQuestions(db, ex.Id)
	if err != nil {
		t.Fatal(err)
	}

	if len(exQ.Questions) != 2 {
		t.Fatal(err)
	}

	_, err = DeleteExerciceById(db, ex.Id)
	if err != nil {
		t.Fatal(err)
	}
}

// TODO: populate a "real" DB to support further testing
