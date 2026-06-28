package commons

// Identifiable is an interface that defines an oobject we may distinguish from another by an id.
// Different id implies different object.
// Because we deal with time, same id DOES NOT IMPLY same version of the object.
type Identifiable interface {
	// Id of the object as a string (uuid for instance)
	Id() string
}
