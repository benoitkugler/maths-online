package teacher

import (
	"os"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/pass"
	tc "github.com/benoitkugler/maths-online/sql/teacher"
	"github.com/benoitkugler/maths-online/utils/testutils"
)

func Test_parsePronoteName(t *testing.T) {
	tests := []struct {
		args        string
		wantName    string
		wantSurname string
	}{
		{"DEIARE Matthéa", "DEIARE", "Matthéa"},
		{"DEMANS-HAUC Jode", "DEMANS-HAUC", "Jode"},
		{"PONCLVES ROHA Oceli", "PONCLVES ROHA", "Oceli"},
	}
	for _, tt := range tests {
		gotName, gotSurname := parsePronoteName(tt.args)
		if gotName != tt.wantName {
			t.Errorf("parsePronoteName() gotName = %v, want %v", gotName, tt.wantName)
		}
		if gotSurname != tt.wantSurname {
			t.Errorf("parsePronoteName() gotSurname = %v, want %v", gotSurname, tt.wantSurname)
		}
	}
}

func Test_parsePronoteStudentList(t *testing.T) {
	f, err := os.Open("students_sample.csv")
	if err != nil {
		t.Skipf("Sample not available: %s", err)
	}

	out, err := parsePronoteStudentList(f)
	if err != nil {
		t.Fatal(err)
	}

	if len(out) != 31 {
		t.Fatal(len(out))
	}
}

func Test_importPronoteFile(t *testing.T) {
	f, err := os.Open("students_sample.csv")
	if err != nil {
		t.Skipf("Sample not available: %s", err)
	}

	db, err := testutils.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", testutils.DB, err)
	}

	ct := Controller{db: db}
	classroom, err := tc.Classroom{IdTeacher: 1}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}
	defer tc.DeleteClassroomById(db, classroom.Id)

	err = ct.importPronoteFile(f, classroom.Id)
	if err != nil {
		t.Fatal(err)
	}

	out, err := ct.getClassroomStudents(classroom.Id)
	if err != nil {
		t.Fatal(err)
	}

	if len(out) != 31 {
		t.Fatal(len(out))
	}

	tc.DeleteStudentsByIdClassrooms(db, classroom.Id)
}

func TestStudentCRUD(t *testing.T) {
	db := testutils.NewTestDB(t, "../../sql/teacher/gen_create.sql")
	defer db.Remove()

	teacher, err := tc.Teacher{Id: 1}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	ct := Controller{db: db.DB, admin: tc.Teacher{Id: teacher.Id}}
	classroom, err := tc.Classroom{IdTeacher: teacher.Id}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	st, err := ct.addStudent(classroom.Id, teacher.Id)
	if err != nil {
		t.Fatal(err)
	}

	st.Name = "sdlsl"
	st.Birthday = tc.Date(time.Now())
	if err = ct.updateStudent(st, teacher.Id); err != nil {
		t.Fatal(err)
	}

	if err = ct.updateStudent(st, teacher.Id+1); err == nil {
		t.Fatal()
	}

	if err = ct.deleteStudent(st.Id, teacher.Id); err != nil {
		t.Fatal(err)
	}
}

func TestDemoStudent(t *testing.T) {
	const DEMO_CODE = "1234"
	db := testutils.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../../migrations/create_manual.sql")
	defer db.Remove()

	ct := NewController(db.DB, pass.SMTP{}, pass.Encrypter{}, pass.Encrypter{}, "localhost:1323", DEMO_CODE)

	_, err := ct.LoadAdminTeacher()
	testutils.Assert(t, err == nil)
	_, err = ct.LoadDemoClassroom()
	testutils.Assert(t, err == nil)

	_, err = ct.attachStudentCandidates("invalid code")
	testutils.Assert(t, err != nil)

	l, err := ct.attachStudentCandidates(DEMO_CODE + ".1")
	testutils.Assert(t, err == nil)
	testutils.Assert(t, len(l) == 1)

	out, err := ct.validAttachStudent(AttachStudentToClassroom2In{
		ClassroomCode: DEMO_CODE + ".1",
		IdStudent:     l[0].Id,
		Birthday:      "2000-01-01",
	})
	testutils.Assert(t, err == nil)
	testutils.Assert(t, !out.ErrAlreadyAttached && !out.ErrInvalidBirthday)
}
