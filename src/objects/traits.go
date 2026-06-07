package objects

import (
	"fmt"
)

// Trait is the immutable definition of a concept that can be instantiated into objects (instances of traits).
type Trait interface {
	Element
	// Name returns the name of the trait, should be unique
	Name() string
	// Attributes returns the attributes of the trait, as name and type
	Attributes() map[string]string
	// WithAttribute returns a new trait with the given attribute added.
	// Name is the name of the attribute to add, and its type is expectedType
	WithAttribute(name, expectedType string) Trait
	// WithoutAttribute returns a new trait with the given attribute removed
	WithoutAttribute(name string) Trait
}

// baseTrait is the standard implementation of the Trait interface.
type baseTrait struct {
	name       string            // name of the trait, should be unique
	attributes map[string]string // attributes of the trait, as name and type
}

// NewTrait creates a new trait with the given name
func NewTrait(name string) Trait {
	return &baseTrait{
		name:       name,
		attributes: make(map[string]string),
	}
}

// Name returns the name of the trait
func (t *baseTrait) Name() string {
	if t == nil {
		return ""
	}
	return t.name
}

// DeclaringClass returns the class that declares this trait.
// It is, at the very least, the CLASS_TRAIT class.
func (t *baseTrait) DeclaringClass() Class {
	return CLASS_TRAIT
}

// String returns a string representation of the trait to include its name
func (t *baseTrait) String() string {
	if t == nil {
		return "Trait{nil}"
	}
	return fmt.Sprintf("Trait{name: %s}", t.name)
}

// Same returns true if the other element is a Trait with the same name and attributes
func (t *baseTrait) Same(other Element) bool {
	if other == nil && t == nil {
		return true
	} else if other == nil || t == nil {
		return false
	} else if !IsElementDeclaredInstance(other, CLASS_TRAIT) {
		return false
	}

	// Safely assert against the new Trait interface
	otherTrait, ok := other.(Trait)
	if !ok {
		return false
	}

	// Compare against the interface method rather than internal fields
	if t.name != otherTrait.Name() {
		return false
	}

	// now, test attributes
	if len(t.attributes) != len(otherTrait.Attributes()) {
		return false
	}

	otherAttributes := otherTrait.Attributes()
	for attr, attrType := range t.attributes {
		if otherType, found := otherAttributes[attr]; !found || attrType != otherType {
			return false
		}
	}

	return true
}

// Attributes returns a copy of all attributes of the trait (key) and type (value)
// A defensive copy is returned to enforce immutability.
func (t *baseTrait) Attributes() map[string]string {
	if t == nil || t.attributes == nil {
		return nil
	}

	result := make(map[string]string, len(t.attributes))
	for k, v := range t.attributes {
		result[k] = v
	}
	return result
}

// WithAttribute returns a new Trait instance with the given attribute added or updated.
// The original Trait remains unchanged.
func (t *baseTrait) WithAttribute(name, expectedType string) Trait {
	if t == nil {
		return nil
	}
	if name == "" || expectedType == "" {
		return t
	}

	// Copy existing attributes
	newAttributes := t.Attributes()
	if newAttributes == nil {
		newAttributes = make(map[string]string)
	}

	// Add/Update the new attribute
	newAttributes[name] = expectedType

	return &baseTrait{
		name:       t.name,
		attributes: newAttributes,
	}
}

// WithoutAttribute returns a new Trait instance without the given attribute.
// If the attribute does not exist, it returns the current Trait instance.
func (t *baseTrait) WithoutAttribute(name string) Trait {
	if t == nil || name == "" {
		return t
	}

	// Fast return if the attribute is not present (avoids useless allocation)
	if _, exists := t.attributes[name]; !exists {
		return t
	}

	// Copy all attributes except the one to remove
	newAttributes := make(map[string]string, len(t.attributes)-1)
	for k, v := range t.attributes {
		if k != name {
			newAttributes[k] = v
		}
	}

	return &baseTrait{
		name:       t.name,
		attributes: newAttributes,
	}
}
