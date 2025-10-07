package models

import (
	"errors"

	"github.com/zefrenchwan/perspectives.git/structures"
)

// objectsToEntities builds a slice of entities from a slice of objects
func objectsToEntities(values []*Object) []ModelEntity {
	var result []ModelEntity
	for _, value := range values {
		result = append(result, value)
	}

	return result
}

// linksToEntities builds a slice of entities from a slice of links
func linksToEntities(values []*Link) []ModelEntity {
	var result []ModelEntity
	for _, value := range values {
		result = append(result, value)
	}

	return result
}

// entitiesGroup is a group of entities.
// It is considered as an entity too.
// An instance of entitiesGroup is a sort of anonymous group.
// If group lasts, prefer an object as a placeholder.
type entitiesGroup []ModelEntity

// GetType returns EntityTypeGroup
func (g entitiesGroup) GetType() EntityType {
	return EntityTypeGroup
}

// NewObjectGroup builds a group of objects (at least 1)
func NewObjectsGroup(objects []*Object) (ModelEntity, error) {
	return newEntitiesGroup(objectsToEntities(objects))
}

// NewGroupOfObjects builds a group of objects from single elements
func NewGroupOfObjects(objects ...*Object) (ModelEntity, error) {
	return NewObjectsGroup(objects)
}

// NewLinksGroup builds an entity as a group of links
func NewLinksGroup(links []*Link) (ModelEntity, error) {
	return newEntitiesGroup(linksToEntities(links))
}

// newEntitiesGroup builds a new group of entities after deduplication
func newEntitiesGroup(elements []ModelEntity) (ModelEntity, error) {
	if len(elements) == 0 {
		return nil, errors.New("empty group not allowed")
	}

	return entitiesGroup(structures.SliceDeduplicateFunc(elements, SameModelEntity)), nil
}
