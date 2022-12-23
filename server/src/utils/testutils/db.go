package testutils

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/server/src/pass"
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

var (
	dbCount      int
	dbCountMutex sync.Mutex
)

// NewTestDB creates a new database and add all the tables
// as defined in the `generateSQLFile` files.
func NewTestDB(t *testing.T, generateSQLFile ...string) TestDB {
	t.Helper()

	const userPassword = "dummy"

	dbCountMutex.Lock()
	name := fmt.Sprintf("tmp_dev_%d_%d", time.Now().UnixNano(), dbCount)
	dbCount++
	dbCountMutex.Unlock()

	// cleanup if needed
	runCmd(exec.Command("dropdb", "--if-exists", name))

	runCmd(exec.Command("createdb", name))

	for _, file := range generateSQLFile {
		file, err := filepath.Abs(file)
		if err != nil {
			t.Fatal(err)
		}
		_, err = os.Stat(file)
		if err != nil {
			t.Fatal(err)
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
	err := db.DB.Close()
	if err != nil {
		panic(err)
	}

	runCmd(exec.Command("dropdb", "--if-exists", "--force", db.name))
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
