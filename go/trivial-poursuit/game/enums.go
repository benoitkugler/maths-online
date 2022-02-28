package game

// a categorie of questions, represented by a color
type categorie uint8

const (
	purple categorie = iota // purple
	green                   // green
	orange                  // orange
	yellow                  // yellow
	blue                    // blue

	nbCategories // the number of categories a player should complete
)
