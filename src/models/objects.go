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

// newAttribute returns a new attribute with no value based on parameters
func newAttribute(name string, semantics []string) Attribute {
	return Attribute{
		name:      name,
		semantics: structures.SliceDeduplicate(semantics),
		values:    make(structures.Mapping[string]),
	}
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

// Same returns true if objects share the same id (one id should be unique)
func (o *Object) Same(other *Object) bool {
	if o == nil && other == nil {
		return true
	} else if o == nil || other == nil {
		return false
	}

	return o.Id == other.Id
}

// IsEmpty returns true if the attribute contains no data
func (a *Attribute) IsEmpty() bool {
	return a == nil || len(a.values) == 0
}

// NewObject returns an object implementing those traits
func NewObject(traits []string) *Object {
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
	result := new(Object)
	result.Id = uuid.NewString()
	result.traits = objectTraits
	result.attributes = make(map[string]Attribute)
	result.lifetime = structures.NewFullPeriod()
	return result
}

// NewObjectSince returns an object that implements traits, valid since creationTime
func NewObjectSince(traits []string, creationTime time.Time) *Object {
	base := NewObject(traits)
	base.lifetime = structures.NewPeriodSince(creationTime, true)
	return base
}

// NewObjectDuring returns an object that implements traits, valid during a period.
// It may raise an error if endTime is before startTime
func NewObjectDuring(traits []string, startTime, endTime time.Time) (*Object, error) {
	if endTime.Before(startTime) {
		return nil, errors.New("cannot make an object with an end date before its start date")
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
func (o *Object) AsGroup() ([]*Object, error) {
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
	for name := range o.attributes {
		result = append(result, name)
	}

	if len(result) == 0 {
		result = make([]string, 0)
		return result
	} else {
		return structures.SliceReduce(result)
	}
}

// setValueDuringPeriod changes that attribute to set value during period.
// If object is nil or period is empty, no action.
// Else value changes during that period no matter the object's lifetime
func (o *Object) setValueDuringPeriod(attribute, value string, period structures.Period) {
	if o == nil {
		return
	} else if period.IsEmpty() {
		return
	}

	if attr, found := o.attributes[attribute]; !found {
		o.attributes[attribute] = Attribute{name: attribute, values: structures.NewValueDuringPeriod(value, period)}
	} else if attr.values == nil {
		attr.values = make(structures.Mapping[string])
		attr.values.SetDuringPeriod(value, period)
		o.attributes[attribute] = attr
	} else {
		attr.values.SetDuringPeriod(value, period)
		o.attributes[attribute] = attr
	}
}

// SetValue sets a value for that attribute
func (o *Object) SetValue(attribute, value string) {
	if o == nil {
		return
	}

	o.setValueDuringPeriod(attribute, value, structures.NewFullPeriod())
}

// SetValueSince sets the value for that attribute since startingTime
func (o *Object) SetValueSince(attribute, value string, startingTime time.Time, includeStartingTime bool) {
	if o == nil {
		return
	}

	period := structures.NewPeriodSince(startingTime, includeStartingTime)
	o.setValueDuringPeriod(attribute, value, period)
}

// SetValueUntil sets the value for that attribute until endingTime
func (o *Object) SetValueUntil(attribute, value string, endingTime time.Time, includeEndingTime bool) {
	if o == nil {
		return
	}

	period := structures.NewPeriodUntil(endingTime, includeEndingTime)
	o.setValueDuringPeriod(attribute, value, period)
}

// SetValueDuring sets value for that attribute during the interval [startingTime, endingTime] (both included)
func (o *Object) SetValueDuring(attribute, value string, startingTime, endingTime time.Time) {
	if o == nil {
		return
	}

	period := structures.NewFinitePeriod(startingTime, endingTime, true, true)
	o.setValueDuringPeriod(attribute, value, period)
}

// GetAllValues returns all the values for all attributes (including the ones with no value)
// Two options:
// Either reduceToObjectLifetime is true and we get values only during object lifetime
// Or reduceToObjectLifetime is false and we get all values
func (o *Object) GetAllValues(reduceToObjectLifetime bool) map[string][]string {
	if o == nil {
		return nil
	}

	result := make(map[string][]string)

	// for each attribute
	for name, attr := range o.attributes {
		// values contain all the possible values
		var values []string
		// for each value and then period for that value
		for value, period := range attr.values.Get() {
			if reduceToObjectLifetime {
				if !period.IsEmpty() && !period.Intersection(o.lifetime).IsEmpty() {
					values = append(values, value)
				}
			} else {
				values = append(values, value)
			}
		}

		// we made the values, so set for that attribute
		result[name] = structures.SliceDeduplicate(values)
	}

	return result
}

// GetValue returns the value for an attribute (by name) if any.
// Result (if any) is then the mapping value -> validity, true or nil, false for no match.
// Depending on reduceToObjectLifetime:
// Either it is true and then validity is the intersection of the object lifetime and the attribute validity
// Or we keep values and matching period as is
func (o *Object) GetValue(attribute string, reduceToObjectLifetime bool) (map[string]structures.Period, bool) {
	if o == nil {
		return nil, false
	}

	// values are the values from the attribute.
	var values map[string]structures.Period
	if attr, found := o.attributes[attribute]; !found {
		return nil, false
	} else {
		values = attr.values.Get()
	}

	// result contains the intersection with the object's lifetime
	result := make(map[string]structures.Period)
	for key, period := range values {
		if reduceToObjectLifetime {
			inter := period.Intersection(o.lifetime)
			if !inter.IsEmpty() {
				result[key] = inter
			}
		} else {
			result[key] = period
		}
	}

	return result, true
}

// Describe returns the structure of the object
func (o *Object) Describe() ObjectDescription {
	if o == nil {
		return ObjectDescription{
			Id: uuid.NewString(),
		}
	}

	attributes := make(map[string][]string)
	for name, attr := range o.attributes {
		semantics := structures.SliceDeduplicate(attr.semantics)
		attributes[name] = semantics
	}

	var traits []string
	for _, value := range o.traits {
		traits = append(traits, value.Name)
	}

	return ObjectDescription{
		Id:         uuid.NewString(),
		IdObject:   o.Id,
		Traits:     structures.SliceReduce(traits),
		Attributes: attributes,
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

// ActivePeriod returns the object's active period
func (o *Object) ActivePeriod() structures.Period {
	return o.lifetime
}

// objectsGroup decorates a slice of objects to match a model entity definition
type objectsGroup []*Object

// GetType returns
func (g objectsGroup) GetType() EntityType {
	return EntityTypeGroup
}

// NewObjectGroup builds a group of objects (at least 1)
func NewObjectsGroup(objects []*Object) (ModelEntity, error) {
	if len(objects) == 0 {
		return nil, errors.New("empty group not allowed as object group")
	}

	result := structures.SliceDeduplicate(objects)
	return objectsGroup(result), nil
}

// NewGroupOfObjects builds a group of objects from single elements
func NewGroupOfObjects(objects ...*Object) (ModelEntity, error) {
	return NewObjectsGroup(objects)
}

// AsLink raises an error
func (g objectsGroup) AsLink() (*Link, error) {
	return nil, errors.ErrUnsupported
}

// AsGroup returns the value as a slice of objects
func (g objectsGroup) AsGroup() ([]*Object, error) {
	return []*Object(g), nil
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
