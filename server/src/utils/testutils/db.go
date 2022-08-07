package testutils

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/benoitkugler/maths-online/pass"
	"github.com/benoitkugler/maths-online/utils"
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

func runCmd(cmd *exec.Cmd) {
	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	err := cmd.Run()
	if err != nil {
		fmt.Println(stdOut.String())
		fmt.Println(stdErr.String())
		panic(err)
	}
}

type TestDB struct {
	*sql.DB
	name string // unique randomly generated
}

// NewTestDB creates a new database and add all the tables
// as defined in the `generateSQLFile` files.
func NewTestDB(t *testing.T, generateSQLFile ...string) TestDB {
	const userPassword = "dummy"

	name := "tmp_dev_" + utils.RandomString(true, 10)

	// cleanup if needed
	runCmd(exec.Command("dropdb", "--if-exists", name))

	runCmd(exec.Command("createdb", name))

	for _, file := range generateSQLFile {
		file, err := filepath.Abs(file)
		if err != nil {
			panic(err)
		}
		_, err = os.Stat(file)
		if err != nil {
			panic(err)
		}
		runCmd(exec.Command("bash", "-c", fmt.Sprintf("psql %s < %s", name, file)))
	}

	logs := pass.DB{
		Name:     name,
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

	return TestDB{DB: db, name: name}
}

// Remove closes the connection and remove the DB.
func (db TestDB) Remove() {
	db.Close()

	runCmd(exec.Command("dropdb", "--if-exists", db.name))
}

// DB is a test DB, usually build from importing the current production DB.
var DB = pass.DB{
	Host:     "localhost",
	User:     "benoit",
	Password: "dummy",
	Name:     "isyro_prod",
}

// WebsocketURL transforms `s` to set `ws` as scheme
// It panics on error.
func WebsocketURL(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	u.Scheme = "ws"
	return u.String()
}
