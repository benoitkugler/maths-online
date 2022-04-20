package students

//go:generate ../../../../../structgen/structgen -source=models_sql.go -mode=sql:gen_scans.go -mode=sql_test:gen_scans_test.go -mode=sql_gen:create_gen.sql  -mode=rand:gen_data_test.go

type Student struct {
	Id      int64
	Name    string
	Surname string
}
