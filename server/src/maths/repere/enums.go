package repere

// LabelPos chooses the relative positionning
// of the label from its "natural" location
type LabelPos uint8

const (
	Top         LabelPos = iota // Au dessus
	Bottom                      // En dessous
	Left                        // A gauche
	Right                       // A droite
	TopLeft                     // Au dessus, à gauche
	TopRight                    // Au dessus, à droite
	BottomRight                 // En dessous, à droite
	BottomLeft                  // En dessous, à gauche
)
