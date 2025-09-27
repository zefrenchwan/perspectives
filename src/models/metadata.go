package models

import "errors"

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

// AsLink raises an error by definition
func (t Trait) AsLink() (*Link, error) {
	return nil, errors.ErrUnsupported
}

// AsGroup raises an error by definition
func (t Trait) AsGroup() ([]*Object, error) {
	return nil, errors.ErrUnsupported
}

// AsObject raises an error by definition
func (t Trait) AsObject() (*Object, error) {
	return nil, errors.ErrUnsupported
}

// AsTrait returns the value as a trait
func (t Trait) AsTrait() (Trait, error) {
	return t, nil
}

// AsVariable raises an error by definition
func (t Trait) AsVariable() (Variable, error) {
	return Variable{}, errors.ErrUnsupported
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
