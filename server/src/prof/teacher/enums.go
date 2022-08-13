package teacher

// Visibility is the status of a ressource, among :
//	- personnal : read/write acces for the current teacher
//	- verified by admins : read access only
type Visibility uint8

const (
	Personnal Visibility = iota // Personnel
	Admin                       // Officiel
)

// NewVisibility returns the visilbity of the ressource owned by
// `ownerID` and requested by `userID`,
// or `false` if `userID` does not have access to it.
func NewVisibility(ownerID, userID, adminID IdTeacher, public bool) (Visibility, bool) {
	var vis Visibility
	if ownerID == userID {
		vis = Personnal
	} else if ownerID == adminID && public {
		vis = Admin
	} else {
		return 0, false
	}
	return vis, true
}
