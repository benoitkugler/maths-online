package teacher

import (
	"github.com/benoitkugler/maths-online/server/src/sql/events"
	tc "github.com/benoitkugler/maths-online/server/src/sql/teacher"
)

type StudentHeader struct {
	Id                 tc.IdStudent
	Label              string
	HasAttachedClients bool
}

func NewStudentHeader(st tc.Student) StudentHeader {
	return StudentHeader{Id: st.Id, Label: st.Name + " " + st.Surname, HasAttachedClients: len(st.Clients) != 0}
}

type AttachStudentToClassroom1Out []StudentHeader

type AttachStudentToClassroom2In struct {
	ClassroomCode string
	IdStudent     tc.IdStudent
	Birthday      string // 2006-01-02
	Device        string // the name of the device the student is using
}

type AttachStudentToClassroom2Out struct {
	ErrInvalidBirthday bool
	ErrAlreadyAttached bool // Deprecated
	IdCrypted          string
}

type CheckStudentClassroomOut struct {
	IsOK bool // if false, ignore the other fields
	Meta StudentClassroomHeader

	// Advance exposes the global advance of the student,
	// as defined by its events.
	Advance events.StudentAdvance
}

type StudentClient struct {
	Name    string
	Surname string

	Id               tc.IdStudent   // Depreacted
	Birthday         tc.Date        // Depreacted
	IdClassroom      tc.IdClassroom `json:"id_classroom"` // Deprecated
	TrivialSuccess   int            // Deprecated
	IsClientAttached bool           // Deprecated
}

type StudentClassroomHeader struct {
	Student           StudentClient
	ClassroomName     string
	TeacherMail       string // or contact, to be displayed
	TeacherContactURL string // optional, display a link if provided
}
