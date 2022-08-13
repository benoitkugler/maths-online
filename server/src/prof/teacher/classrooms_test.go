package teacher

import (
	"os"
	"testing"
	"time"

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
	classroom, err := Classroom{IdTeacher: 1}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}
	defer DeleteClassroomById(db, classroom.Id)

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

	DeleteStudentsByIdClassrooms(db, classroom.Id)
}

func TestStudentCRUD(t *testing.T) {
	db := testutils.NewTestDB(t, "gen_create.sql")
	defer db.Remove()

	teacher, err := Teacher{Id: 1}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	ct := Controller{db: db.DB, admin: Teacher{Id: teacher.Id}}
	classroom, err := Classroom{IdTeacher: teacher.Id}.Insert(db)
	if err != nil {
		t.Fatal(err)
	}

	st, err := ct.addStudent(classroom.Id, teacher.Id)
	if err != nil {
		t.Fatal(err)
	}

	st.Name = "sdlsl"
	st.Birthday = Date(time.Now())
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
