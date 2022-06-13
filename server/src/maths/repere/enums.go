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
	Hide                        // Masque la légende
)

type SegmentKind uint8

const (
	SKSegment SegmentKind = iota // Segment
	SKVector                     // Vecteur
	SKLine                       // Droite (infinie)
)
