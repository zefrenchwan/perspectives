package structures

import (
	"errors"
	"fmt"
)

// parent_link defines an "extends" relation
const parent_link = 0

// Hierarchy is a DAG of named S , link is extends.
// Assumption is that a name is unique
type Hierarchy[S any] struct {
	// parents link childs to parents
	parents DVGraph[string, int]
	// values link instances of S by name
	values map[string]any
}

// NewHierarchy builds an empty hierarchy
func NewHierarchy[S any]() Hierarchy[S] {
	return Hierarchy[S]{
		parents: NewDVGraph[string, int](),
		values:  make(map[string]any),
	}
}

// Set value to a key (a name)
func (h Hierarchy[S]) Set(key string, value any) {
	h.values[key] = value
	h.parents.AddNode(key)
}

// LinkToParent links a child (assumed to exist) to a parent (assumed to exist)
func (h Hierarchy[S]) LinkToParent(child, parent string) error {
	_, sourceExists := h.values[child]
	_, destinationExists := h.values[parent]
	if child == parent {
		return errors.New("cannot link element with itself")
	} else if !sourceExists {
		return errors.New("child does not exist")
	} else if !destinationExists {
		return errors.New("parent does not exist")
	} else if !h.parents.LinkWithoutCycle(child, parent, parent_link) {
		return fmt.Errorf("linking %s to %s would make a cycle", child, parent)
	} else {
		return nil
	}
}
