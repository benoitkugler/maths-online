package teacher

import "github.com/benoitkugler/maths-online/server/src/sql/teacher"

// Visibility is the status of a ressource, among :
//   - personnal : read/write acces for the current teacher
//   - verified by admins : read access only
type Visibility uint8

const (
	Hidden    Visibility = iota // not accessible by the user, except in reviews
	Personnal                   // Personnel
	Admin                       // Officiel
)

// NewVisibility returns the visilbity of the ressource owned by
// `ownerID` and requested by `userID`,
// or `false` if `userID` does not have access to it.
func NewVisibility(ownerID, userID, adminID teacher.IdTeacher, public bool) Visibility {
	if ownerID == userID {
		return Personnal
	} else if ownerID == adminID && public {
		return Admin
	} else {
		return Hidden
	}
}

// Restricted returns true if the item access if forbidden.
func (vis Visibility) Restricted() bool { return vis == Hidden }
