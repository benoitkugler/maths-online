package students

//go:generate ../../../../../structgen/structgen -source=models_sql.go -mode=sql:gen_scans.go -mode=sql_test:gen_scans_test.go -mode=sql_gen:gen_create.sql  -mode=rand:gen_data_test.go

type Student struct {
	Id               int64
	Name             string
	Surname          string
	Birthday         Date
	TrivialSuccess   int
	IsClientAttached bool // true if at least one student appli has claimed this profile
}
