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
