package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/zefrenchwan/perspectives.git/structures"
)

// Attribute defines an attribute of an object.
// It is a name, a set of tags to define its semantic layer, and time dependant values
type Attribute struct {
	// name of the attribute (for instance "first name")
	name string
	// Semantics of the attribute (for instance "email")
	semantics []string
	// Values over time for that attribute
	values structures.Mapping[string]
}

// Object defines an object for a given lifetime with values
type Object struct {
	// Id of the object (assumed to be unique)
	Id string
	// traits of the object, by name
	traits []string
	// attributes of the object, key is attribute name
	attributes map[string]Attribute
	// lifetime of the object, that is the period that object "lives"
	lifetime structures.Period
}

// NewObject returns an object implementing those traits
func NewObject(traits []string) Object {
	return Object{
		Id:         uuid.NewString(),
		traits:     structures.SliceReduce(traits),
		attributes: make(map[string]Attribute),
		lifetime:   structures.NewFullPeriod(),
	}
}

// NewObjectSince returns an object that implements traits, valid since creationTime
func NewObjectSince(traits []string, creationTime time.Time) Object {
	base := NewObject(traits)
	base.lifetime = structures.NewPeriodSince(creationTime, true)
	return base
}

// NewObjectDuring returns an object that implements traits, valid during a period.
// It may raise an error if endTime is before startTime
func NewObjectDuring(traits []string, startTime, endTime time.Time) (Object, error) {
	if endTime.Before(startTime) {
		return Object{}, errors.New("cannot make an object with an end date before its start date")
	}

	base := NewObject(traits)
	base.lifetime = structures.NewFinitePeriod(startTime, endTime, true, true)
	return base, nil
}

// DeclaringTraits returns the declaring traits for that object
func (o *Object) DeclaringTraits() []string {
	return o.traits
}

// AddSemanticForAttribute flags this attribute for that particular meaning.
// If the attribute did not exist before, it is created
func (o *Object) AddSemanticForAttribute(attribute, meaning string) {
	if attr, found := o.attributes[attribute]; !found {
		o.attributes[attribute] = Attribute{name: attribute, semantics: []string{meaning}}
	} else {
		newValues := append(attr.semantics, meaning)
		attr.semantics = structures.SliceReduce(newValues)
		o.attributes[attribute] = attr
	}
}

// GetSemanticForAttribute returns semantic values for that attribute (if any)
func (o *Object) GetSemanticForAttribute(attribute string) ([]string, bool) {
	if attr, found := o.attributes[attribute]; !found {
		return nil, false
	} else {
		return attr.semantics, true
	}
}

// SetValue sets a value for that attribute
func (o *Object) SetValue(attribute, value string) {
	if attr, found := o.attributes[attribute]; !found {
		o.attributes[attribute] = Attribute{name: attribute, values: structures.NewValue(value)}
	} else if attr.values == nil {
		attr.values = make(structures.Mapping[string])
		attr.values.Set(value)
		o.attributes[attribute] = attr
	} else {
		attr.values.Set(value)
		o.attributes[attribute] = attr
	}
}

// GetAllValues returns all the values for all the attributes
func (o *Object) GetAllValues() map[string][]string {
	result := make(map[string][]string)

	for name, attr := range o.attributes {
		// GetValues deal with nil, no need to check for nil values
		values := attr.values.GetValues()
		if len(values) != 0 {
			result[name] = values
		}
	}

	return result
}
