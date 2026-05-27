package commons

import (
	"fmt"
	"sync"
)

// Trait is the definition of a concept that can be instantiated into objects.
type Trait struct {
	name       string            // name of the trait, should be unique
	locks      sync.RWMutex      // locks to manage concurrent access to attributes
	attributes map[string]string // attributes of the trait, as name and type
}

// Id returns the name as the unique identifier of the trait
func (t *Trait) Id() string {
	return t.name
}

// Name returns the name of the trait
func (t *Trait) Name() string {
	return t.name
}

// DeclaringClass returns the class that declares this trait.
// It is, at the very least, the CLASS_TRAIT class.
func (t *Trait) DeclaringClass() Class {
	return CLASS_TRAIT
}

// String returns a string representation of the trait to include its name
func (t *Trait) String() string {
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

// WithAttribute sets an attribute for the trait with the given name and expected type.
// To chain, return the trait instance
func (t *Trait) WithAttribute(name, expectedType string) *Trait {
	if name == "" || expectedType == "" {
		return t
	}
	t.locks.Lock()
	defer t.locks.Unlock()
	t.attributes[name] = expectedType
	return t
}

// RemoveAttribute removes an attribute from the trait with the given name
func (t *Trait) RemoveAttribute(name string) {
	if name == "" {
		return
	}

	t.locks.Lock()
	defer t.locks.Unlock()
	delete(t.attributes, name)
}

// Attributes returns all attributes of the trait (key) and type (value)
func (t *Trait) Attributes() map[string]string {
	t.locks.RLock()
	defer t.locks.RUnlock()
	return t.attributes
}

// NewTrait creates a new trait with the given name
func NewTrait(name string) *Trait {
	return &Trait{
		attributes: make(map[string]string),
		name:       name,
	}
}
