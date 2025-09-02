package structures

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Attribute defines an attribute of an object.
// It is a name, a set of tags to define its semantic layer, and time dependant values
type Attribute struct {
	// Name of the attribute (for instance "first name")
	Name string
	// Semantics of the attribute (for instance "email")
	Semantics []string
	// Values over time for that attribute
	Values Mapping[string]
}

// Object defines an object for a given lifetime with values
type Object struct {
	// Id of the object (assumed to be unique)
	Id string
	// Traits of the object, by name
	Traits []string
	// Attributes of the object, key is attribute name
	Attributes map[string]Attribute
	// Lifetime of the object, that is the period that object "lives"
	Lifetime Period
}

// NewObject returns an object implementing those traits
func NewObject(traits []string) Object {
	return Object{
		Id:         uuid.NewString(),
		Traits:     SliceReduce(traits),
		Attributes: make(map[string]Attribute),
		Lifetime:   NewFullPeriod(),
	}
}

// NewObjectSince returns an object that implements traits, valid since creationTime
func NewObjectSince(traits []string, creationTime time.Time) Object {
	base := NewObject(traits)
	base.Lifetime = NewPeriodSince(creationTime, true)
	return base
}

// NewObjectDuring returns an object that implements traits, valid during a period.
// It may raise an error if endTime is before startTime
func NewObjectDuring(traits []string, startTime, endTime time.Time) (Object, error) {
	if endTime.Before(startTime) {
		return Object{}, errors.New("cannot make an object with an end date before its start date")
	}

	base := NewObject(traits)
	base.Lifetime = NewFinitePeriod(startTime, endTime, true, true)
	return base, nil
}
