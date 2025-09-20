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
	traits []Trait
	// attributes of the object, key is attribute name
	attributes map[string]Attribute
	// lifetime of the object, that is the period that object "lives"
	lifetime structures.Period
}

// IsEmpty returns true if the attribute contains no data
func (a *Attribute) IsEmpty() bool {
	return a == nil || len(a.values) == 0
}

// NewObject returns an object implementing those traits
func NewObject(traits []string) Object {
	// map and deduplicate traits
	declaringTraits := make(map[string]Trait)

	for _, trait := range traits {
		declaringTraits[trait] = Trait{Name: trait}
	}

	var objectTraits []Trait
	for _, value := range declaringTraits {
		objectTraits = append(objectTraits, value)
	}

	// then, build the object
	return Object{
		Id:         uuid.NewString(),
		traits:     objectTraits,
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
	var result []string
	for _, trait := range o.traits {
		result = append(result, trait.Name)
	}

	return structures.SliceReduce(result)
}

// GetType returns the type of the entity, that is EntityTypeObject
func (o *Object) GetType() EntityType {
	return EntityTypeObject
}

// AsLink would raise an error
func (o *Object) AsLink() (*Link, error) {
	return nil, errors.ErrUnsupported
}

// AsGroup would raise an error
func (o *Object) AsGroup() ([]Object, error) {
	return nil, errors.ErrUnsupported
}

// AsObject returns the object
func (o *Object) AsObject() (*Object, error) {
	return o, nil
}

// AsTrait raises an error
func (o *Object) AsTrait() (Trait, error) {
	return Trait{}, errors.ErrUnsupported
}

// AsVariable raises an error
func (o *Object) AsVariable() (Variable, error) {
	return Variable{}, errors.ErrUnsupported
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

// Attributes return the attributes of the object
func (o *Object) Attributes() []string {
	var result []string
	for name, attr := range o.attributes {
		if attr.IsEmpty() {
			result = append(result, name)
		}
	}

	if len(result) == 0 {
		result = make([]string, 0)
		return result
	} else {
		return structures.SliceReduce(result)
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

// Describe returns the structure of the object
func (o *Object) Describe() ObjectDescription {
	if o == nil {
		return ObjectDescription{
			Id: uuid.NewString(),
		}
	}

	var attributes []string
	for name, attr := range o.attributes {
		if !attr.IsEmpty() {
			attributes = append(attributes, name)
		}
	}

	var traits []string
	for _, value := range o.traits {
		traits = append(traits, value.Name)
	}

	return ObjectDescription{
		Id:         uuid.NewString(),
		IdObject:   o.Id,
		Traits:     structures.SliceReduce(traits),
		Attributes: structures.SliceReduce(attributes),
	}
}

// Equals returns true for same object based on id
func (o *Object) Equals(other *Object) bool {
	if o == nil && other == nil {
		return true
	} else if o == nil || other == nil {
		return false
	}
	return o.Id == other.Id
}

// Duration returns the object's active period
func (o *Object) Duration() structures.Period {
	return o.lifetime
}

// objectsGroup decorates a slice of objects to match a model entity definition
type objectsGroup []Object

// GetType returns
func (g objectsGroup) GetType() EntityType {
	return EntityTypeGroup
}

// AsLink raises an error
func (g objectsGroup) AsLink() (*Link, error) {
	return nil, errors.ErrUnsupported
}

// AsGroup returns the value as a slice of objects
func (g objectsGroup) AsGroup() ([]Object, error) {
	return []Object(g), nil
}

// AsObject raises an error
func (g objectsGroup) AsObject() (*Object, error) {
	return nil, errors.ErrUnsupported
}

// AsTrait raises an error
func (g objectsGroup) AsTrait() (Trait, error) {
	return Trait{}, errors.ErrUnsupported
}

// AsVariable raises an error
func (g objectsGroup) AsVariable() (Variable, error) {
	return Variable{}, errors.ErrUnsupported
}
