package teacher

// Visibility is the status of a ressource, among :
//	- personnal : read/write acces for the current teacher
//	- public from other teachers : read access for everyone, write for the owner
//	- verified by admins : read access only
type Visibility uint8

const (
	Personnal Visibility = iota // Personnel
	Shared                      // Partagé par la communauté
	Admin                       // Officiel
)

// NewVisibility returns the visilbity of the ressource owned by
// `ownerID` and requested by `userID`,
// or `false` if `userID` does not have access to it.
func NewVisibility(ownerID, userID, adminID int64, public bool) (Visibility, bool) {
	var vis Visibility
	if ownerID == userID {
		vis = Personnal
	} else if ownerID == adminID {
		vis = Admin
	} else if public {
		vis = Shared
	} else {
		return 0, false
	}
	return vis, true
}
