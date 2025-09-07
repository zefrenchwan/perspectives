package models

// Trait defines a general label for an object
type Trait struct {
	// Name of the label
	Name string
}

// NewTrait returns a new trait for that label
func NewTrait(label string) Trait {
	return Trait{Name: label}
}

// ObjectDescription describes the object
type ObjectDescription struct {
	// Id of the description (not the object)
	Id string
	// Id of the object
	IdObject string
	// Traits of the object
	Traits []string
	// Attributes of the object
	Attributes []string
}
