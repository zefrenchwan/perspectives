package values

const PRIMITIVE_TYPE_BOOL = "bool"
const PRIMITIVE_TYPE_INT = "int"
const PRIMITIVE_TYPE_STRING = "string"
const PRIMITIVE_TYPE_TIME = "time"
const PRIMITIVE_TYPE_FLOAT = "float64"
const REFERENCE_TYPE = "reference"

// Value is the generic interface for any value, primitive or reference.
type Value interface {
	// Datatype returns the type of the value.
	Datatype() string
	// Equals compares value with other.
	// If true, it means same type and content.
	Equals(any) bool
	// Content returns the underlying value.
	Content() any
	// Serialize returns the serialized form of the value.
	Serialize() string
}
