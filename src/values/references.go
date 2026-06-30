package values

import "strconv"

// ReferenceValue is a reference to an id.
// For instance, on a given entity.
type ReferenceValue struct {
	// referenceId is the id of the reference.
	referenceId string
	// serialized is the serialized form of the reference.
	serialized string
}

// Datatype returns the type of the value : a reference.
func (r ReferenceValue) Datatype() string {
	return REFERENCE_TYPE
}

// Equals compares value with other. If true, it means a reference to same id
func (r ReferenceValue) Equals(other any) bool {
	if other == nil {
		return false
	} else if otherValue, ok := other.(ReferenceValue); !ok {
		return false
	} else if r.Datatype() != otherValue.Datatype() {
		return false
	} else {
		return r.referenceId == otherValue.referenceId
	}
}

// Content returns the underlying value : the reference id.
func (r ReferenceValue) Content() any {
	return r.referenceId
}

// Serialize returns the serialized form of the value.
func (r ReferenceValue) Serialize() string {
	return r.serialized
}

// NewReference creates a new reference value for that id.
func NewReference(otherId string) ReferenceValue {
	result := ReferenceValue{referenceId: otherId}
	result.serialized = REFERENCE_TYPE + "|" + strconv.Itoa(len(otherId)) + "|" + otherId
	return result
}
