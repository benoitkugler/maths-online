package teacher

type AskInscriptionIn struct {
	Mail     string
	Password string
}

type LogginIn struct {
	Mail     string
	Password string
}

type LogginOut struct {
	Error string // empty means success
	Token string // token to use in the next requests
}
