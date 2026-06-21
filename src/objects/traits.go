package objects

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zefrenchwan/perspectives.git/periods"
)

// Trait is a formal definition of attributes with a name.
// For instance, a trait can be "student" with attributes "student id".
// A trait is defined by its name and its attributes.
// As any key type here, Trait is immutable.
// It means that, to build one, use a TraitBuilder.
type Trait struct {
	// id of the trait (calculated from the name)
	id string
	// name of the trait, should be unique
	name string
	// attributes of the trait.
	// Keys are attribute names and values are attribute types (should be primitive types)
	attributes map[string]string
	// hashValue (calculated once) is the hash value of the trait.
	hashValue string
}

// Id returns the id of the trait.
// Id is NOT the name of the trait.
func (t Trait) Id() string {
	return t.id
}

// Name returns the name of the trait
func (t Trait) Name() string {
	return t.name
}

// Same returns true if the trait is the same as the other element.
// It means a trait, with the same attributes.
func (t Trait) Same(other Element) bool {
	if other == nil {
		return false
	} else if other.DeclaringClass() != CLASS_DEFINITION {
		return false
	} else if other.Id() != t.Id() {
		return false
	}

	return other.toHashString() == t.hashValue
}

// DeclaringClass returns the class for that element.
func (t Trait) DeclaringClass() Class {
	return CLASS_DEFINITION
}

// Attributes returns the attributes of the trait as a range iterator.
func (t Trait) Attributes(yield func(attributeName string, attributeType string) bool) {
	for attributeName, attributeType := range t.attributes {
		if !yield(attributeName, attributeType) {
			return
		}
	}
}

// isLinkable defines the trait as a linkable
func (t Trait) isLinkable() bool {
	return true
}

// toHashString returns the hash value of the trait
func (t Trait) toHashString() string {
	return t.hashValue
}

// String returns a string representation of the trait
func (t Trait) String() string {
	return fmt.Sprintf("Trait(%s)", t.name)
}

// Matches returns true if the instance matches the trait.
// It means that all attributes of the trait are present in the instance and have the same type.
// Result is the period during which the instance matches the trait.
// For instance, if the trait is "student" with attributes "student id",
// and the instance has attributes "student id" for five years,
// then the instance matches the trait during those five years.
func (t Trait) Matches(i Instance) (periods.Period, bool) {
	result := i.Activity()
	for attributeName, attributeType := range t.attributes {
		if matchingAttr, ok := i.Attribute(attributeName); ok {
			if attributeType != matchingAttr.AttributeType {
				return periods.Period{}, false
			} else {
				result = result.Intersection(matchingAttr.AttributeValidity)
			}
		} else {
			return periods.Period{}, false
		}

		if result.IsEmpty() {
			return periods.Period{}, false
		}
	}
	return result, true
}

// =========================================================================
// TRAIT BUILDER IMPLEMENTATION
// =========================================================================

// TraitBuilder is a builder for traits.
// One cannot use traits from scratch, they should be built using the builder.
// Builder checks for basic errors, and errors are accumulated.
// Conventionally, methods return the builder itself to allow method chaining.
type TraitBuilder interface {
	// WithName sets the name for that trait to build.
	WithName(name string) TraitBuilder
	// WithAttribute adds an attribute (or changes its type) for the trait to build.
	// Name should not be blank, and type must be a valid primitive type name.
	WithAttribute(attrName, attrType string) TraitBuilder
	// WithoutAttribute removes the given attribute from the builder.
	// If attribute is blank or does not exist, it is a no-op.
	WithoutAttribute(attrName string) TraitBuilder
	// Errors returns the errors encountered while building the trait.
	// Errors are cumulative.
	Errors() error
	// Build returns a trait built from the builder, or combined errors if any.
	// It resets the builder to its initial state after building.
	Build() (Trait, error)
}

// baseTraitBuilder manages in-memory trait creation and modification.
// It implements the TraitBuilder interface using pointer receivers
// to ensure map references and error states are mutated safely.
type baseTraitBuilder struct {
	// name of the trait to build
	name string
	// attributes of the trait to build, keys are names, values are types
	attributes map[string]string
	// globalErrors contains all global errors accumulated during the building process
	globalErrors error
}

// NewTraitBuilder returns a new empty trait builder.
func NewTraitBuilder() TraitBuilder {
	return &baseTraitBuilder{attributes: make(map[string]string)}
}

// TraitBuilderLoad returns a trait builder loaded with the given trait.
func TraitBuilderLoad(trait Trait) TraitBuilder {
	builder := &baseTraitBuilder{attributes: make(map[string]string)}
	builder.name = trait.name
	for attrName, attrValue := range trait.attributes {
		builder.attributes[attrName] = attrValue
	}

	return builder
}

// WithName sets the name for that trait to build.
func (b *baseTraitBuilder) WithName(name string) TraitBuilder {
	if strings.TrimSpace(name) == "" {
		b.globalErrors = errors.Join(b.globalErrors, fmt.Errorf("trait name cannot be empty"))
		return b
	}

	b.name = name
	return b
}

// WithAttribute adds an attribute (or change type) for the trait to build.
// Name should not be blank, and type must be a primitive value.
func (b *baseTraitBuilder) WithAttribute(attrName, attrType string) TraitBuilder {
	if strings.TrimSpace(attrName) == "" {
		b.globalErrors = errors.Join(b.globalErrors, fmt.Errorf("attribute name cannot be empty"))
		return b
	} else if !IsPrimitiveTypeName(attrType) {
		b.globalErrors = errors.Join(b.globalErrors, fmt.Errorf("attribute type must be a valid primitive type name"))
		return b
	}

	if b.attributes == nil {
		b.attributes = make(map[string]string)
	}
	b.attributes[attrName] = attrType
	return b
}

// WithoutAttribute removes the given attribute from the builder.
// If attribute is blank or does not exist, it is a no-op.
func (b *baseTraitBuilder) WithoutAttribute(attrName string) TraitBuilder {
	if b.attributes != nil {
		delete(b.attributes, attrName)
	}
	return b
}

// Errors returns the errors encountered while building the trait.
func (b *baseTraitBuilder) Errors() error {
	return b.globalErrors
}

// Build returns a trait built from the builder, or combined errors if any.
func (b *baseTraitBuilder) Build() (Trait, error) {
	// deal with either empty or default builder
	if b.globalErrors == nil && b.name == "" {
		b.attributes = make(map[string]string)
		return Trait{}, errors.New("trait name cannot be empty")
	}

	// normal process starts here :
	if b.globalErrors != nil {
		globalErr := b.globalErrors
		// Reset the builder
		b.name = ""
		b.attributes = make(map[string]string)
		b.globalErrors = nil
		return Trait{}, globalErr
	}

	id := string(CLASS_DEFINITION) + " => TRAIT : " + b.name

	// DEFENSIVE COPY: We copy the map before assigning it to the Trait.
	// This ensures strict immutability. Even though the builder is reset,
	// copying guarantees no lingering references can mutate the Trait attributes.
	attributesCopy := make(map[string]string, len(b.attributes))
	for k, v := range b.attributes {
		attributesCopy[k] = v
	}

	result := Trait{
		id:         id,
		name:       b.name,
		attributes: attributesCopy,
	}

	result.hashValue = hashTrait(result)

	// Reset the builder
	b.name = ""
	b.attributes = make(map[string]string)

	return result, nil
}
