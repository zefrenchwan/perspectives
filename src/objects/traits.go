package objects

import (
	"fmt"
)

// Trait is the immutable definition of a concept that can be instantiated into objects (instances of traits).
type Trait struct {
	name       string            // name of the trait, should be unique
	attributes map[string]string // attributes of the trait, as name and type
}

// NewTrait creates a new trait with the given name
func NewTrait(name string) *Trait {
	return &Trait{
		name:       name,
		attributes: make(map[string]string),
	}
}

// Name returns the name of the trait
func (t *Trait) Name() string {
	if t == nil {
		return ""
	}
	return t.name
}

// DeclaringClass returns the class that declares this trait.
// It is, at the very least, the CLASS_TRAIT class.
func (t *Trait) DeclaringClass() Class {
	return CLASS_TRAIT
}

// String returns a string representation of the trait to include its name
func (t *Trait) String() string {
	if t == nil {
		return "Trait{nil}"
	}
	return fmt.Sprintf("Trait{name: %s}", t.name)
}

// Same returns true if the other element is a Trait pointer with the same name
func (t *Trait) Same(other Element) bool {
	if other == nil && t == nil {
		return true
	} else if other == nil || t == nil {
		return false
	} else if !IsElementDeclaredInstance(other, CLASS_TRAIT) {
		return false
	} else if value, ok := other.(*Trait); !ok {
		return false
	} else if t.name != value.name {
		return false
	}
	return true
}

// Attributes returns a copy of all attributes of the trait (key) and type (value)
// A defensive copy is returned to enforce immutability.
func (t *Trait) Attributes() map[string]string {
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
func (t *Trait) WithAttribute(name, expectedType string) *Trait {
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

	return &Trait{
		name:       t.name,
		attributes: newAttributes,
	}
}

// WithoutAttribute returns a new Trait instance without the given attribute.
// If the attribute does not exist, it returns the current Trait instance.
func (t *Trait) WithoutAttribute(name string) *Trait {
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

	return &Trait{
		name:       t.name,
		attributes: newAttributes,
	}
}
