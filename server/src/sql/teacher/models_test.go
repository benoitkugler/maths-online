package teacher

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func TestTime(t *testing.T) {
	ti := time.Now()
	tu.Assert(t, ti.String()[0:10] == ti.Format(DateLayout))

	d := Date(ti)
	s, _ := json.Marshal(d)

	var d2 Date
	err := json.Unmarshal(s, &d2)
	tu.AssertNoErr(t, err)

	tu.Assert(t, ti.Equal(time.Time(d2)))
}

func TestSQLTime(t *testing.T) {
	db := tu.NewTestDB(t, "gen_create.sql")
	defer db.Remove()

	teacher, _ := Teacher{}.Insert(db)
	classromm, _ := Classroom{IdTeacher: teacher.Id}.Insert(db)

	st, err := Student{Birthday: Date(time.Now()), IdClassroom: classromm.Id}.Insert(db)
	tu.AssertNoErr(t, err)
	tu.Assert(t, !time.Time(st.Birthday).IsZero())
}

func TestRoot(t *testing.T) {
	// create a DB shared by all tests
	db := tu.NewTestDB(t, "gen_create.sql")
	defer db.Remove()

	t.Run("CRUD for Teacher", func(t *testing.T) { testTeacher(t, db.DB) })
	t.Run("CRUD for Classroom", func(t *testing.T) { testClassroom(t, db.DB) })
	t.Run("CRUD for Student", func(t *testing.T) { testStudent(t, db.DB) })
}

func testTeacher(t *testing.T, db *sql.DB) {
	teacher := randTeacher()
	teacher, err := teacher.Insert(db)
	tu.AssertNoErr(t, err)

	teachers, err := SelectAllTeachers(db)
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(teachers) == 1)

	_, err = DeleteTeacherById(db, teacher.Id)
	tu.AssertNoErr(t, err)
}

func testClassroom(t *testing.T, db *sql.DB) {
	tc, err := randTeacher().Insert(db)
	tu.AssertNoErr(t, err)

	classroom := randClassroom()
	classroom.IdTeacher = tc.Id
	classroom, err = classroom.Insert(db)
	tu.AssertNoErr(t, err)

	classrooms, err := SelectAllClassrooms(db)
	tu.AssertNoErr(t, err)

	if len(classrooms) != 1 {
		t.Fatal(err)
	}

	_, err = DeleteClassroomById(db, classroom.Id)
	tu.AssertNoErr(t, err)
}

func testStudent(t *testing.T, db *sql.DB) {
	tc, err := randTeacher().Insert(db)
	tu.AssertNoErr(t, err)

	classroom := randClassroom()
	classroom.IdTeacher = tc.Id
	classroom, err = classroom.Insert(db)
	tu.AssertNoErr(t, err)

	student := randStudent()
	student.IdClassroom = classroom.Id
	student, err = student.Insert(db)
	tu.AssertNoErr(t, err)

	students, err := SelectAllStudents(db)
	tu.AssertNoErr(t, err)

	if len(students) != 1 {
		t.Fatal(err)
	}

	_, err = DeleteStudentById(db, student.Id)
	tu.AssertNoErr(t, err)
}
