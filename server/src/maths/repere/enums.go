package repere

// LabelPos chooses the relative positionning
// of the label from its "natural" location
type LabelPos uint8

const (
	Top         LabelPos = iota // Top
	Bottom                      // Bottom
	Left                        // Left
	Right                       // Right
	TopLeft                     // TopLeft
	TopRight                    // TopRight
	BottomRight                 // BottomRight
	BottomLeft                  // BottomLeft
)
