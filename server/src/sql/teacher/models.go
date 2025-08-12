// Package teacher provides the data structures related to teacher and student accounts.
package teacher

//go:generate ../../../../../gomacro/cmd/gomacro models.go sql:gen_create.sql go/sqlcrud:gen_scans.go go/randdata:gen_randdata_test.go

type (
	IdTeacher   int64
	IdClassroom int64
	IdStudent   int64
)

// Teacher stores the data associated to one teacher account
// gomacro:SQL ADD UNIQUE(Mail)
type Teacher struct {
	Id                  IdTeacher `json:"id"`
	Mail                string    `json:"mail"`
	PasswordCrypted     []byte    `json:"password_crypted"`      // crypted
	IsAdmin             bool      `json:"is_admin"`              // almost always false
	HasSimplifiedEditor bool      `json:"has_simplified_editor"` // true will hide maths widgets in editor
	Contact             Contact   `json:"contact"`               // if empty, [Mail] is used
	FavoriteMatiere     MatiereTag
}

// Classroom is one group of student controlled by one or many teachers
type Classroom struct {
	Id               IdClassroom `json:"id"`
	Name             string      `json:"name"`
	MaxRankThreshold int         // for the last guilde, default to 40000
}

// TeacherClassroom is a link table describing
// the owners of a classroom.
//
// gomacro:SQL ADD UNIQUE(IdTeacher, IdClassroom)
type TeacherClassroom struct {
	IdTeacher   IdTeacher
	IdClassroom IdClassroom
}

// SelectClassroomsByIdTeacher does NOT wrap the error.
func SelectClassroomsByIdTeacher(db DB, id IdTeacher) (Classrooms, error) {
	links, err := SelectTeacherClassroomsByIdTeachers(db, id)
	if err != nil {
		return nil, err
	}
	return SelectClassrooms(db, links.IdClassrooms()...)
}

// ClassroomCode is a time limited, user friendly, code to access one class
// gomacro:SQL ADD UNIQUE(Code)
type ClassroomCode struct {
	IdClassroom IdClassroom `gomacro-sql-on-delete:"CASCADE"`
	Code        string
	ExpiresAt   Time
}

// Student is a student profile, always attached to a classroom.
type Student struct {
	Id       IdStudent
	Name     string
	Surname  string
	Birthday Date

	IdClassroom IdClassroom `json:"id_classroom" gomacro-sql-on-delete:"CASCADE"`

	Clients Clients
}
