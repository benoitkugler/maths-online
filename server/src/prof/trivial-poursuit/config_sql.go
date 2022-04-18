package trivialpoursuit

import "github.com/benoitkugler/maths-online/trivial-poursuit/game"

//go:generate ../../../../../structgen/structgen -source=config_sql.go -mode=sql:gen_scans.go -mode=sql_test:gen_scans_test.go -mode=sql_gen:create_gen.sql  -mode=rand:gen_data_test.go  -mode=ts:../../../../prof/src/controller/trivial_config_gen.ts

// TrivialConfig is a trivial game configuration
// stored in the DB, one per activity.
type TrivialConfig struct {
	Id                 int64
	IsLaunched         bool
	Questions          CategoriesQuestions
	QuestionTimeout    int // in seconds
	GroupStrategy      GroupStrategy
	MaxPlayersPerGroup int
}

// CategoriesQuestions defines a union of intersection of tags,
// for every category.
type CategoriesQuestions [game.NbCategories][][]string

// GroupStrategy defines how players are matched
type GroupStrategy uint8
