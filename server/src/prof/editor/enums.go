package editor

const (
	Diff1 DifficultyTag = "\u2605"             // 1 étoile
	Diff2 DifficultyTag = "\u2605\u2605"       // 2 étoiles
	Diff3 DifficultyTag = "\u2605\u2605\u2605" // 3 étoiles
)

const (
	Seconde   LevelTag = "2NDE" // Seconde
	Premiere  LevelTag = "1ERE" // Première
	Terminale LevelTag = "TERM" // Terminale
)

type Flow uint8

const (
	Parallel   Flow = iota // Questions indépendantes
	Sequencial             // Questions liées
)
