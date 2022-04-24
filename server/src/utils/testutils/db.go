package testutils

import (
	"bytes"
	"database/sql"
	"os/exec"
	"testing"

	"github.com/benoitkugler/maths-online/pass"
)

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

// CreateDBDev creates a new database and add all the tables
// as defined in `generateSQLFile`.
func CreateDBDev(t *testing.T, generateSQLFile string) *sql.DB {
	const userPassword = "dummy"

	// cleanup if needed
	err := exec.Command("dropdb", "--if-exists", "tmp_dev_test").Run()
	if err != nil {
		panic(err)
	}

	err = exec.Command("createdb", "tmp_dev_test").Run()
	if err != nil {
		panic(err)
	}

	err = exec.Command("bash", "-c", "psql tmp_dev_test < "+generateSQLFile).Run()
	if err != nil {
		panic(err)
	}

	logs := pass.DB{
		Name:     "tmp_dev_test",
		Host:     "localhost",
		User:     getUserName(),
		Password: userPassword,
	}
	db, err := logs.ConnectPostgres()
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	t.Log("Successfully created dev DB")

	return db
}

func RemoveDBDev() {
	err := exec.Command("dropdb", "--if-exists", "tmp_dev_test").Run()
	if err != nil {
		panic(err)
	}
}
