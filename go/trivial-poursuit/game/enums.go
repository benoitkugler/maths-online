package game

// a categorie of questions, represented by a color
type categorie uint8

const (
	Purple categorie = iota // purple
	Green                   // green
	Orange                  // orange
	Yellow                  // yellow
	Blue                    // blue

	nbCategories // the number of categories a player should complete
)
