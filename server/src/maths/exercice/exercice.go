// Package exercice models the content (question, explanation, hints,..)
// of an exercice
package exercice

// Parcours modelize a sequence of exercice, with
// multiples paths.
// It is represented as a directed graph.
type Parcours struct {
	Id int64
}

type ParcoursEtape struct {
	Id int64

	Condition  string // TODO:
	IdExercice int64
}

type ParcoursLien struct {
	IdParcoursEtape int64
}
