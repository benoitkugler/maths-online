package teacher

type StudentHeader struct {
	Id    IdStudent
	Label string
}

type AttachStudentToClassroom1Out = []StudentHeader

type AttachStudentToClassroom2In struct {
	ClassroomCode string
	IdStudent     IdStudent
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
	Student       Student
	ClassroomName string
	TeacherMail   string
}
