package models

import (
	"errors"

	"github.com/zefrenchwan/perspectives.git/structures"
)

// objectsGroup decorates a slice of objects to match a model entity definition.
// It is not a public type because it is just a "sort of placeholder".
// This type represents an unnamed group of objects, something local.
// For instance, Marie and John went to drink a coffee.
// There is no need to make a special entity with properties and lifetime to represent it.
type objectsGroup []*Object

// GetType returns
func (g objectsGroup) GetType() EntityType {
	return EntityTypeGroup
}

// NewObjectGroup builds a group of objects (at least 1)
func NewObjectsGroup(objects []*Object) (ModelEntity, error) {
	if len(objects) == 0 {
		return nil, errors.New("empty group not allowed as object group")
	}

	result := structures.SliceDeduplicate(objects)
	return objectsGroup(result), nil
}

// NewGroupOfObjects builds a group of objects from single elements
func NewGroupOfObjects(objects ...*Object) (ModelEntity, error) {
	return NewObjectsGroup(objects)
}
