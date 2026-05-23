package commons

import (
	"fmt"
)

// Trait is the definition of a concept that can be instantiated into objects.
type Trait struct {
	name string // name of the trait, should be unique
}

// Id returns the name as the unique identifier of the trait
func (t Trait) Id() string {
	return t.name
}

// Name returns the name of the trait
func (t Trait) Name() string {
	return t.name
}

// DeclaringClasses returns the classes that declare this trait.
// It is, at the very least, the CLASS_TRAIT class.
func (t Trait) DeclaringClasses() []Class {
	return []Class{CLASS_TRAIT}
}

// String returns a string representation of the trait to include its name
func (t Trait) String() string {
	return fmt.Sprintf("Trait{name: %s}", t.name)
}

// Same returns true if the other element is a Trait with the same name
func (t Trait) Same(other Element) bool {
	if other == nil {
		return false
	} else if !IsElementDeclaredInstance(other, CLASS_TRAIT) {
		return false
	} else if value, ok := other.(Trait); !ok {
		return false
	} else if t.name != value.name {
		return false
	}
	return true
}

// NewTrait creates a new trait with the given name
func NewTrait(name string) Trait {
	return Trait{name: name}
}
