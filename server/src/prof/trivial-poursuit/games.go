package trivialpoursuit

// gameSessionController monitor the games of one session (think one classroom)
// and broadcast the main advances from all the games to the teacher client
type gameSessionController struct {
	advances chan bool
}
