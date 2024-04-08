// Package teacher provides the data structures related to teacher and student accounts.
package teacher

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

// Classroom is one group of student controlled by a teacher
// gomacro:SQL ADD UNIQUE(Id, IdTeacher)
type Classroom struct {
	Id               IdClassroom `json:"id"`
	IdTeacher        IdTeacher   `json:"id_teacher" gomacro-sql-on-delete:"CASCADE"`
	Name             string      `json:"name"`
	MaxRankThreshold int         // for the last guilde, default to 40000
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
