// package questions models the content (question, explanation, hints,..)
// of an exercice.
// An exercice is roughly a list of questions, and a question is represented
// in this package in 3 forms:
//	- as written by a teacher, containing random parameters
//	- instanciated, where random parameters are generated
//	- as send to the client : the answer are stripped and the expression are formatted
package questions

// // Parcours modelize a sequence of exercice, with
// // multiples paths.
// // It is represented as a directed graph.
// type Parcours struct {
// 	Id int64
// }

// type ParcoursEtape struct {
// 	Id int64

// 	Condition  string // TODO:
// 	IdExercice int64
// }

// type ParcoursLien struct {
// 	IdParcoursEtape int64
// }

type Exercice struct {
	Id          int64      `json:"id"`
	Title       string     `json:"title"`       // name for the exercice
	Description string     `json:"description"` // overall description for all questions (optional)
	Parameters  Parameters `json:"parameters"`  // random parameters shared by the questions
}
