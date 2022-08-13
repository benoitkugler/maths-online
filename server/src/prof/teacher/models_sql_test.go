package teacher

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/utils/testutils"
)

func TestTime(t *testing.T) {
	ti := time.Now()
	if ti.String()[0:10] != ti.Format(DateLayout) {
		t.Fatal()
	}

	d := Date(ti)
	s, _ := json.Marshal(d)

	var d2 Date
	err := json.Unmarshal(s, &d2)
	if err != nil {
		t.Fatal(err)
	}
	if !ti.Equal(time.Time(d2)) {
		t.Fatal("invalid json", string(s))
	}
}

func TestSQLTime(t *testing.T) {
	db := testutils.NewTestDB(t, "gen_create.sql")
	defer db.Remove()

	teacher, _ := Teacher{}.Insert(db)
	classromm, _ := Classroom{IdTeacher: teacher.Id}.Insert(db)

	st, err := Student{Birthday: Date(time.Now()), IdClassroom: classromm.Id}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	if time.Time(st.Birthday).IsZero() {
		t.Fatal()
	}
}

func TestRoot(t *testing.T) {
	// create a DB shared by all tests
	db := testutils.NewTestDB(t, "gen_create.sql")
	defer db.Remove()

	t.Run("CRUD for Teacher", func(t *testing.T) { testTeacher(t, db.DB) })
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
