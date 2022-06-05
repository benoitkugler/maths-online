package teacher

import "github.com/benoitkugler/maths-online/prof/students"

//go:generate ../../../../../structgen/structgen -source=classrooms_student_types.go -mode=dart:../../../../eleve/lib/shared/students.gen.dart

type AttachStudentToClassroomOut = []students.Student
