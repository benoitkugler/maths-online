package trivial

// a Categorie of questions, represented by a color
type Categorie uint8

const (
	Purple Categorie = iota // purple
	Green                   // green
	Orange                  // orange
	Yellow                  // yellow
	Blue                    // blue

	nbCategories // the number of categories a player should complete
)
