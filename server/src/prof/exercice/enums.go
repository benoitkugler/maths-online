package exercice

type Flow uint8

const (
	Parallel   Flow = iota // Questions indépendantes
	Sequencial             // Questions liées
)
