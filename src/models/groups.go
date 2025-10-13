package models

import (
	"errors"

	"github.com/zefrenchwan/perspectives.git/commons"
)

// objectsToEntities builds a slice of entities from a slice of objects
func objectsToEntities(values []*Object) []Entity {
	var result []Entity
	for _, value := range values {
		result = append(result, value)
	}

	return result
}

// linksToEntities builds a slice of entities from a slice of links
func linksToEntities(values []*Link) []Entity {
	var result []Entity
	for _, value := range values {
		result = append(result, value)
	}

	return result
}

// entitiesGroup is a group of entities.
// It is considered as an entity too.
// An instance of entitiesGroup is a sort of anonymous group.
// If group lasts, prefer an object as a placeholder.
type entitiesGroup []Entity

// GetType returns EntityTypeGroup
func (g entitiesGroup) GetType() EntityType {
	return EntityTypeGroup
}

// NewObjectGroup builds a group of objects (at least 1)
func NewObjectsGroup(objects []*Object) (Entity, error) {
	return newEntitiesGroup(objectsToEntities(objects))
}

// NewGroupOfObjects builds a group of objects from single elements
func NewGroupOfObjects(objects ...*Object) (Entity, error) {
	return NewObjectsGroup(objects)
}

// NewLinksGroup builds an entity as a group of links
func NewLinksGroup(links []*Link) (Entity, error) {
	return newEntitiesGroup(linksToEntities(links))
}

// newEntitiesGroup builds a new group of entities after deduplication
func newEntitiesGroup(elements []Entity) (Entity, error) {
	if len(elements) == 0 {
		return nil, errors.New("empty group not allowed")
	}

	return entitiesGroup(commons.SliceDeduplicateFunc(elements, SameEntity)), nil
}
