package models

import (
	"github.com/zefrenchwan/perspectives.git/structures"
)

// ObjectDescription describes the object
type ObjectDescription struct {
	// Id of the description (not the object)
	Id string
	// Id of the object
	IdObject string
	// Traits of the object
	Traits []string
	// Attributes of the object
	Attributes map[string][]string
}

// BuildEmptyObjectFromDescription returns an EMPTY object from a description.
// Result has:
// the id provided in parameter
// a lifetime set to full period,
// and attributes with provided semantics
//
// Once built, it is strongly encouraged to change lifetime
func (d ObjectDescription) BuildEmptyObjectFromDescription(newId string) *Object {
	result := new(Object)
	result.Id = newId
	result.attributes = make(map[string]Attribute)
	result.lifetime = structures.NewFullPeriod()

	for _, trait := range structures.SliceDeduplicate(d.Traits) {
		newTrait := NewTrait(trait)
		result.traits = append(result.traits, newTrait)
	}

	for name, semantics := range d.Attributes {
		result.attributes[name] = newAttribute(name, semantics)
	}

	return result
}

// BuildObjectFromDescription builds an object from a description, changes its lifetime and sets values.
// Result is an object with:
// Id as the newId
// Attributes set from description or coming from values
// Values set as constant values (no period for attributes)
// Object's lifetime is lifetime
func (d ObjectDescription) BuildObjectFromDescription(newId string, lifetime structures.Period, values map[string]string) *Object {
	result := d.BuildEmptyObjectFromDescription(newId)
	result.Id = newId
	result.lifetime = lifetime
	for k, v := range values {
		result.SetValue(k, v)
	}

	return result
}
