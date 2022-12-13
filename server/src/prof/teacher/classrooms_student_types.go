package teacher

import tc "github.com/benoitkugler/maths-online/server/src/sql/teacher"

type StudentHeader struct {
	Id    tc.IdStudent
	Label string
}

func NewStudentHeader(st tc.Student) StudentHeader {
	return StudentHeader{Id: st.Id, Label: st.Name + " " + st.Surname}
}

type AttachStudentToClassroom1Out = []StudentHeader

type AttachStudentToClassroom2In struct {
	ClassroomCode string
	IdStudent     tc.IdStudent
	Birthday      string // 2006-01-02
}

type AttachStudentToClassroom2Out struct {
	ErrInvalidBirthday bool
	ErrAlreadyAttached bool
	IdCrypted          string
}

type CheckStudentClassroomOut struct {
	IsOK bool // if not, ignore `meta`
	Meta StudentClassroomHeader
}

type StudentClassroomHeader struct {
	Student       tc.Student
	ClassroomName string
	TeacherMail   string
}
