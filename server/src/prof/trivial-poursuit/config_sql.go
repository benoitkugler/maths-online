package trivialpoursuit

//go:generate ../../../../../structgen/structgen -source=config_sql.go -mode=sql:gen_scans.go -mode=sql_gen:gen_create.sql -mode=ts:../../../../prof/src/controller/trivial_config_gen.ts

// TrivialConfig is a trivial game configuration
// stored in the DB, one per activity.
type TrivialConfig struct {
	Id              int64
	Questions       CategoriesQuestions
	QuestionTimeout int // in seconds
	ShowDecrassage  bool
}

// TeacherTrivialConfig is a link table listing the trivials owned
// by one teacher.
//
// A config have at most one owner.
// sql:ADD UNIQUE(id_trivial_config)
type TeacherTrivialConfig struct {
	IdTeacher       int64 `json:"id_teacher"`
	IdTrivialConfig int64 `json:"id_trivial_config"`
	IsPublic        bool  `json:"is_public"` // If true, the config is available to other teachers (as read-only)
}
