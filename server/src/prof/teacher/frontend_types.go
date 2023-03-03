package teacher

type AskInscriptionIn struct {
	Mail                string
	Password            string
	HasEditorSimplified bool
}

type AskInscriptionOut struct {
	Error           string // empty for no error
	IsPasswordError bool
}

type LogginIn struct {
	Mail     string
	Password string
}

type LogginOut struct {
	Error           string // empty means success
	IsPasswordError bool
	Token           string // token to use in the next requests
}
