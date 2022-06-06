package teacher

//go:generate ../../../../../structgen/structgen -source=classrooms_student_types.go -mode=dart:../../../../eleve/lib/shared/students.gen.dart

type StudentHeader struct {
	Id    int64
	Label string
}

type AttachStudentToClassroom1Out = []StudentHeader

type AttachStudentToClassroom2In struct {
	ClassroomCode string
	IdStudent     int64
	Birthday      string // 2006-01-02
}

type AttachStudentToClassroom2Out struct {
	ErrInvalidBirthday bool
	ErrAlreadyAttached bool
	IdCrypted          string
}
