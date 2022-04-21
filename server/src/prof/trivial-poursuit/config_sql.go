package trivialpoursuit

//go:generate ../../../../../structgen/structgen -source=config_sql.go -mode=sql:gen_scans.go -mode=sql_test:gen_scans_test.go -mode=sql_gen:create_gen.sql -mode=rand:gen_data_test.go  -mode=ts:../../../../prof/src/controller/trivial_config_gen.ts

// TrivialConfig is a trivial game configuration
// stored in the DB, one per activity.
type TrivialConfig struct {
	Id              int64
	LaunchSessionID string // empty before launch
	Questions       CategoriesQuestions
	QuestionTimeout int // in seconds
}

func (tc *TrivialConfig) IsLaunched() bool { return tc.LaunchSessionID != "" }
