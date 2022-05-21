package teacher

import (
	"database/sql"
	"testing"

	"github.com/benoitkugler/maths-online/utils/testutils"
)

func TestRoot(t *testing.T) {
	// create a DB shared by all tests
	db := testutils.CreateDBDev(t, "gen_create.sql")
	defer testutils.RemoveDBDev()
	defer db.Close()

	t.Run("CRUD for Teacher", func(t *testing.T) { testTeacher(t, db) })
}

func testTeacher(t *testing.T, db *sql.DB) {
	teacher := randTeacher()
	teacher, err := teacher.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	teachears, err := SelectAllTeachers(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(teachears) != 1 {
		t.Fatal(err)
	}
}
