package models

import (
	"errors"

	"github.com/zefrenchwan/perspectives.git/structures"
)

// EntityType defines the type of an entity to use.
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

// EntityTypeGroup is the type for slices of entities
const EntityTypeGroup EntityType = 4

// EntityTypeVariable is the type for variables
const EntityTypeVariable EntityType = 5

// ModelEntity is the general definition of an entity in the model we use.
// It decorates:
// links as pointers because we may modify them
// Objects as pointers because we may modify them too
// Group of entities
// Traits as immutable objects
// Variables as immutable objects
type ModelEntity interface {
	// GetType returns the type of the entity (trait ? link ? object ? )
	GetType() EntityType
}

// AsTrait returns the model entiy as a trait if possible, or an error
func AsTrait(e ModelEntity) (Trait, error) {
	if e == nil {
		return Trait{}, errors.New("nil value")
	} else if result, ok := e.(Trait); !ok {
		return Trait{}, errors.New("failed to cast as a trait")
	} else {
		return result, nil
	}
}

// AsLink returns the entity as a link, or raises an error if e was not a link
func AsLink(e ModelEntity) (*Link, error) {
	if e == nil {
		return nil, nil
	} else if result, ok := e.(*Link); !ok {
		return nil, errors.New("failed to cast as a link")
	} else {
		return result, nil
	}
}

// AsObject returns the entity as an object, or raises an error if e was not an object
func AsObject(e ModelEntity) (*Object, error) {
	if e == nil {
		return nil, nil
	} else if result, ok := e.(*Object); !ok {
		return nil, errors.New("failed to cast as an object")
	} else {
		return result, nil
	}
}

// AsVariable returns the entity as a variable, or raises an error if e was not a variable
func AsVariable(e ModelEntity) (Variable, error) {
	if e == nil {
		return Variable{}, errors.New("nil value")
	} else if result, ok := e.(Variable); !ok {
		return Variable{}, errors.New("failed to cast as a variable")
	} else {
		return result, nil
	}
}

// AsGroup returns the entity as a group, or an error if e was not a group
func AsGroup(e ModelEntity) ([]ModelEntity, error) {
	if e == nil {
		return nil, errors.New("nil value")
	} else if result, ok := e.(entitiesGroup); !ok {
		return nil, errors.New("failed to cast as a group")
	} else {
		return result, nil
	}
}

// TemporalEntity defines an entity with a duration.
// For links, it means the period a link is active during.
// For objects, it means the period the object is alive during.
type TemporalEntity interface {
	// We want a temporal entity to be an entity
	ModelEntity
	// ActivePeriod is the period the entity is active during
	ActivePeriod() structures.Period
	// SetActivity forces the period for that temporal entity
	SetActivity(newPeriod structures.Period)
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
		aLink, _ := AsLink(a)
		bLink, _ := AsLink(b)
		return aLink.Same(bLink)
	case EntityTypeGroup:
		aGroup, _ := AsGroup(a)
		bGroup, _ := AsGroup(b)
		return structures.SlicesEqualsAsSetsFunc(aGroup, bGroup, SameModelEntity)
	case EntityTypeObject:
		aObject, _ := AsObject(a)
		bObject, _ := AsObject(b)
		return aObject.Same(bObject)
	case EntityTypeVariable:
		aVar, _ := AsVariable(a)
		bVar, _ := AsVariable(b)
		return aVar.Same(bVar)
	case EntityTypeTrait:
		aTrait, _ := AsTrait(a)
		bTrait, _ := AsTrait(b)
		return aTrait.Equals(bTrait)
	default:
		return false
	}
}

// IdentifiableEntity defines an entity that has an id.
// An id should be globally unique : no link should have the same id as an object.
// An entity has an id if any observer may distinguish it from another.
type IdentifiableEntity interface {
	Id() string // Id returns the id of that entity.
}
