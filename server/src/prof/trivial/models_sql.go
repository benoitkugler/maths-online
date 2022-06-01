package trivial

//go:generate ../../../../../structgen/structgen -source=models_sql.go -mode=sql:gen_scans.go -mode=sql_gen:gen_create.sql -mode=ts:../../../../prof/src/controller/trivial_config_gen.ts

// TrivialConfig is a trivial game configuration
// stored in the DB, one per activity.
type TrivialConfig struct {
	Id              int64
	Questions       CategoriesQuestions
	QuestionTimeout int // in seconds
	ShowDecrassage  bool
	Public          bool
	IdTeacher       int64 `json:"id_teacher"`
	Name            string
}
