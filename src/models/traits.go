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

// Equals returns true if names match
func (t Trait) Equals(other Trait) bool {
	return t.Name == other.Name
}

// GetType returns EntityTypeTrait
func (t Trait) GetType() EntityType {
	return EntityTypeTrait
}
