// Package teacher implements the logic related to teacher accounts.
package teacher

//go:generate ../../../../../structgen/structgen -source=models_sql.go -mode=sql:gen_scans.go -mode=sql_gen:gen_create.sql -mode=rand:gen_data_test.go

// Teacher stores the data associated to one teacher account
// sql:ADD UNIQUE(mail)
type Teacher struct {
	Id              int64  `json:"id"`
	Mail            string `json:"mail"`
	PasswordCrypted []byte `json:"password_crypted"` // crypted
	IsAdmin         bool   `json:"is_admin"`         // almost always false
}

// Classroom is one group of student controlled by a teacher
type Classroom struct {
	Id        int64  `json:"id"`
	IdTeacher int64  `json:"id_teacher" sql_on_delete:"CASCADE"`
	Name      string `json:"name"`
}
