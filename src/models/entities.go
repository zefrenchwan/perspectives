package models

import "github.com/zefrenchwan/perspectives.git/structures"

// EntityType defines the type of an entity to use
// So far, accepted types are:
// objects: for instance: John knows Jane
// traits: for instance John Likes Chocolate (with Chocolate a trait)
// links: for instance John knows (Marie likes Chocolate)
// groups: for instance, Mary and John (as a group) like Chocolate
// variables: to be replaced by any previous type
type EntityType int

// EntityTypeTrait is the type for traits
const EntityTypeTrait EntityType = 1

// EntityTypeLink is the type for links
const EntityTypeLink EntityType = 2

// EntityTypeObject is the type for objects
const EntityTypeObject EntityType = 3

// EntityTypeGroup is the type for slices of objects
const EntityTypeGroup EntityType = 4

// EntityTypeVariable is the type for variables
const EntityTypeVariable EntityType = 5

// ModelEntity is the general definition of an entity in the model we use.
// It decorates:
// links as pointers because we may modify them
// Objects as pointers because we may modify them too
// Group of objects (as pointers for the same reason)
// Traits as immutable objects
// Variables as immutable objects
type ModelEntity interface {
	// GetType returns the type of the entity (trait ? link ? object ? )
	GetType() EntityType
	// AsLink casts the value as a link, or raises an error it underlying content is not a link
	AsLink() (*Link, error)
	// AsGroup casts the value as a group of objects, or raises an error it underlying content is not a group
	AsGroup() ([]*Object, error)
	// AsObject casts the value as an object, or raises an error it underlying content is not an object
	AsObject() (*Object, error)
	// AsTrait returns the value as a trait, or raises an error it underlying content is not a trait
	AsTrait() (Trait, error)
	// AsVariable returns the value as a variable, or raises an error if underlying content is not a variable
	AsVariable() (Variable, error)
}

// TemporalEntity defines an entity with a duration.
// For links, it means the period a link is active during.
// For objects, it means the period the object is alive during.
type TemporalEntity interface {
	// We want a temporal entity to be an entity
	ModelEntity
	// ActivePeriod is the period the entity is active during
	ActivePeriod() structures.Period
}

// SameModelEntity tests if two model entities are the same based on their own definition of same
func SameModelEntity(a, b ModelEntity) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil || b == nil {
		return false
	}

	aType := a.GetType()
	bType := b.GetType()
	if aType != bType {
		return false
	}

	switch aType {
	case EntityTypeLink:
		aLink, _ := a.AsLink()
		bLink, _ := b.AsLink()
		return aLink.Same(bLink)
	case EntityTypeGroup:
		aGroup, _ := a.AsGroup()
		bGroup, _ := b.AsGroup()
		return structures.SlicesEqualsAsSetsFunc(aGroup, bGroup, func(a, b *Object) bool { return a.Equals(b) })
	case EntityTypeObject:
		aObject, _ := a.AsObject()
		bObject, _ := b.AsObject()
		return aObject.Same(bObject)
	case EntityTypeVariable:
		aVar, _ := a.AsVariable()
		bVar, _ := b.AsVariable()
		return aVar.Same(bVar)
	case EntityTypeTrait:
		aTrait, _ := a.AsTrait()
		bTrait, _ := b.AsTrait()
		return aTrait.Equals(bTrait)
	default:
		return false
	}
}
