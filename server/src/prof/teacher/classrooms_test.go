package teacher

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/benoitkugler/maths-online/server/src/pass"
	"github.com/benoitkugler/maths-online/server/src/sql/events"
	tc "github.com/benoitkugler/maths-online/server/src/sql/teacher"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
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
	tu.AssertNoErr(t, err)

	if len(out) != 31 {
		t.Fatal(len(out))
	}
}

func Test_importPronoteFile(t *testing.T) {
	f, err := os.Open("students_sample.csv")
	if err != nil {
		t.Skipf("Sample not available: %s", err)
	}

	db, err := tu.DB.ConnectPostgres()
	if err != nil {
		t.Skipf("DB %v not available : %s", tu.DB, err)
	}

	ct := Controller{db: db}
	classroom, err := tc.Classroom{IdTeacher: 1}.Insert(db)
	tu.AssertNoErr(t, err)

	defer tc.DeleteClassroomById(db, classroom.Id)

	err = ct.importPronoteFile(f, classroom.Id)
	tu.AssertNoErr(t, err)

	out, err := ct.getClassroomStudents(classroom.Id)
	tu.AssertNoErr(t, err)

	if len(out) != 31 {
		t.Fatal(len(out))
	}

	tc.DeleteStudentsByIdClassrooms(db, classroom.Id)
}

func TestStudentCRUD(t *testing.T) {
	db := tu.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../sql/events/gen_create.sql")
	defer db.Remove()

	teacher, err := tc.Teacher{Id: 1, FavoriteMatiere: tc.Mathematiques}.Insert(db)
	tu.AssertNoErr(t, err)

	ct := Controller{db: db.DB, admin: tc.Teacher{Id: teacher.Id}}
	classroom, err := tc.Classroom{IdTeacher: teacher.Id}.Insert(db)
	tu.AssertNoErr(t, err)

	st, err := ct.addStudent(classroom.Id, teacher.Id)
	tu.AssertNoErr(t, err)

	st.Student.Name = "sdlsl"
	st.Student.Birthday = tc.Date(time.Now())
	if err = ct.updateStudent(st.Student, teacher.Id); err != nil {
		t.Fatal(err)
	}

	if err = ct.updateStudent(st.Student, teacher.Id+1); err == nil {
		t.Fatal()
	}

	encID := ct.studentKey.EncryptID(int64(st.Student.Id))
	profile, err := ct.checkStudentClassroom(encID)
	tu.AssertNoErr(t, err)
	tu.Assert(t, profile.IsOK)
	tu.Assert(t, profile.Advance.Rank == 0)

	// add 300 * 10 = 3000 points
	for range [10]int{} {
		_, err = events.RegisterEvents(ct.db, st.Student.Id, events.E_IsyTriv_Win)
		tu.AssertNoErr(t, err)
	}
	profile, err = ct.checkStudentClassroom(encID)
	tu.AssertNoErr(t, err)
	tu.Assert(t, profile.Advance.Rank == 1)

	if err = ct.deleteStudent(st.Student.Id, teacher.Id); err != nil {
		t.Fatal(err)
	}
}

func TestDemoStudent(t *testing.T) {
	const DEMO_CODE = "1234"
	db := tu.NewTestDB(t, "../../sql/teacher/gen_create.sql", "../../../migrations/create_manual.sql")
	defer db.Remove()

	ct := NewController(db.DB, pass.SMTP{}, pass.Encrypter{}, pass.Encrypter{}, "localhost:1323", DEMO_CODE)

	_, err := ct.LoadAdminTeacher()
	tu.AssertNoErr(t, err)
	_, err = ct.LoadDemoClassroom()
	tu.AssertNoErr(t, err)

	_, err = ct.attachStudentCandidates("invalid code")
	tu.Assert(t, err != nil)

	l, err := ct.attachStudentCandidates(DEMO_CODE + ".1")
	tu.AssertNoErr(t, err)
	tu.Assert(t, len(l) == 1)

	out, err := ct.validAttachStudent(AttachStudentToClassroom2In{
		ClassroomCode: DEMO_CODE + ".1",
		IdStudent:     l[0].Id,
		Birthday:      "2000-01-01",
		Device:        "Xiaomi",
	})
	tu.AssertNoErr(t, err)
	tu.Assert(t, !out.ErrInvalidBirthday)
}

func TestClientJSON(t *testing.T) {
	cl := tc.Client{Device: "", Time: time.Now()}
	b, _ := json.Marshal(cl)
	fmt.Println(string(b))
}
