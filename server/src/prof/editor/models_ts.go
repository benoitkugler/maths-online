package editor

//go:generate ../../../../../structgen/structgen -source=models_ts.go -mode=ts:../../../../prof/src/controller/exercice_gen.ts

// LevelTag are special question tags used to indicate the
// level (class) for the question.
type LevelTag string
