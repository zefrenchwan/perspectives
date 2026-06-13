package objects

import (
	"errors"
	"fmt"
	"maps"
)

// Trait is the immutable definition of a concept that can be instantiated into objects (instances of traits).
type Trait interface {
	Element
	// Name returns the name of the trait, should be unique
	Name() string
	// Attributes returns the attributes of the trait, as name and type
	Attributes() map[string]string
}

// TraitBuilder is the interface for building traits.
// It allows to build traits by adding name and attributes without overloading GC.
// One may modify the trait builder by adding or removing attributes,
// but errors encountered during the building of the trait are cumulative.
// It means that even if an error is encountered once, it stays as an error.
type TraitBuilder interface {
	// WithName changes the name of the trait.
	WithName(string) TraitBuilder
	// WithAttribute adds or changes the attribute of the trait.
	// Name is the name of the attribute, and expectedType is the type for that attribute.
	// It should be a primitive type.
	WithAttribute(name, expectedType string) TraitBuilder
	// WithoutAttribute removes the attribute of the trait.
	WithoutAttribute(name string) TraitBuilder
	// Errors returns the errors encountered during the building of the trait.
	Errors() error
	// Build returns the trait and the errors encountered during the building of the trait.
	// It resets the builder to its initial state.
	Build() (Trait, error)
}

// baseTrait is the standard implementation of the Trait interface.
type baseTrait struct {
	name       string            // name of the trait, should be unique
	attributes map[string]string // attributes of the trait, as name and type
}

// Id returns the name of the trait, should be unique
func (t *baseTrait) Id() string {
	return t.name
}

// Name returns the name of the trait
func (t *baseTrait) Name() string {
	if t == nil {
		return ""
	}
	return t.name
}

// isLinkable is a SEALED INTERFACE pattern implementation.
// It allows traits to be linked to other elements.
func (t *baseTrait) isLinkable() bool {
	return true
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
	} else if !IsInstanceOfClass(other, CLASS_TRAIT) {
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

// baseTraitBuilder is a builder for Trait instances that avoids reconstructing the trait many times.
type baseTraitBuilder struct {
	// name of the trait to build
	name string
	// attributes of the trait to build
	attributes map[string]string
	// globalErrors is a list of errors that occurred during the construction of the trait.
	// Note that they are cumulative
	globalErrors error
}

// WithName sets the name of the trait.
func (t *baseTraitBuilder) WithName(name string) TraitBuilder {
	if name == "" {
		t.globalErrors = errors.Join(t.globalErrors, errors.New("trait name cannot be empty"))
		return t
	}

	t.name = name
	return t
}

// WithAttribute adds an attribute to the trait builder.
// Name is the name of the attribute to add, and expectedType is the type of the attribute.
// Only primitive types are allowed.
func (t *baseTraitBuilder) WithAttribute(name, expectedType string) TraitBuilder {
	if name == "" {
		t.globalErrors = errors.Join(t.globalErrors, errors.New("attribute name cannot be empty"))
		return t
	} else if expectedType == "" {
		t.globalErrors = errors.Join(t.globalErrors, errors.New("attribute type cannot be empty"))
		return t
	} else if !IsPrimitiveTypeName(expectedType) {
		t.globalErrors = errors.Join(t.globalErrors, errors.New("attribute type must be a primitive type"))
		return t
	}

	t.attributes[name] = expectedType
	return t
}

// WithoutAttribute removes the attribute by name
func (t *baseTraitBuilder) WithoutAttribute(name string) TraitBuilder {
	if name == "" {
		// we may avoid an error, but it is an error
		//t.globalErrors = errors.New("attribute name cannot be empty")
		return t
	}

	delete(t.attributes, name)
	return t
}

// Errors returns the errors during trait builder, so far.
func (t *baseTraitBuilder) Errors() error {
	return t.globalErrors
}

// Build returns the trait, or all errors during its generation.
func (t *baseTraitBuilder) Build() (Trait, error) {
	if t.globalErrors != nil {
		return nil, t.globalErrors
	}

	attributes := make(map[string]string)
	maps.Copy(attributes, t.attributes)

	result := &baseTrait{
		name:       t.name,
		attributes: attributes,
	}

	t.name = ""
	t.globalErrors = nil
	t.attributes = make(map[string]string)

	return result, nil
}

// NewTraitBuilder creates a new empty trait builder.
func NewTraitBuilder() TraitBuilder {
	return &baseTraitBuilder{
		attributes: make(map[string]string),
	}
}

// TraitBuilderLoad loads a trait into a trait builder.
// This way, we may modify the trait again by rebuilding the new trait builder.
func TraitBuilderLoad(other Trait) TraitBuilder {
	if other == nil {
		return &baseTraitBuilder{
			globalErrors: errors.New("base trait cannot be nil"),
			attributes:   make(map[string]string),
		}
	}

	attributes := make(map[string]string)
	maps.Copy(attributes, other.Attributes())

	return &baseTraitBuilder{
		name:       other.Name(),
		attributes: attributes,
	}
}
