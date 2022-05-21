package editor

//go:generate ../../../../../structgen/structgen -source=models_sql.go -mode=sql:gen_scans.go -mode=sql_gen:gen_create.sql

// TeacherQuestion is a link table listing the questions owned
// by one teacher.
// Admin questions are owned by no one
//
// A question have at most one owner.
// sql:ADD UNIQUE(id_question)
type TeacherQuestion struct {
	IdTeacher  int64 `json:"id_teacher"`
	IdQuestion int64 `json:"id_question"`
	IsPublic   bool  `json:"is_public"` // If true, the question is available to other teachers (as read-only)
}
