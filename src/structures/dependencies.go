package structures

import (
	"errors"
	"fmt"
)

// depends_link defines an "depends" relation
const depends_link = 0

// Dependencies is a DAG of named S , link is depends.
// Assumption is that a name is unique
type Dependencies[S any] struct {
	// parents link childs to parents
	parents DVGraph[string, int]
	// values link instances of S by name
	values map[string]S
}

// NewDependencies builds an empty Dependencies instance
func NewDependencies[S any]() Dependencies[S] {
	return Dependencies[S]{
		parents: NewDVGraph[string, int](),
		values:  make(map[string]S),
	}
}

// SetValue sets value to a key (a name)
func (h Dependencies[S]) SetValue(key string, value S) {
	h.values[key] = value
	h.parents.AddNode(key)
}

// GetValue returns the value associated to that key, if any.
// It returns the value, true if found, empty, false otherwise
func (h Dependencies[S]) GetValue(key string) (S, bool) {
	var empty S
	if v, found := h.values[key]; found {
		return v, found
	}
	return empty, false
}

// AddDependency links a child (assumed to exist) to a parent (assumed to exist)
func (h Dependencies[S]) AddDependency(source, dependency string) error {
	_, sourceExists := h.values[source]
	_, destinationExists := h.values[dependency]
	if source == dependency {
		return errors.New("cannot link element with itself")
	} else if !sourceExists {
		return errors.New("child does not exist")
	} else if !destinationExists {
		return errors.New("parent does not exist")
	} else if !h.parents.LinkWithoutCycle(source, dependency, depends_link) {
		return fmt.Errorf("linking %s to %s would make a cycle", source, dependency)
	} else {
		return nil
	}
}

// LoadWithDependencies returns all the dependencies from a node.
// For instance, if "a" -> X depends on "b" -> Y, then result for "a" would be "a" -> X , "b" -> Y
func (h Dependencies[S]) LoadWithDependencies(name string) map[string]S {
	result := make(map[string]S)
	h.parents.Walk(name, func(source string) {
		result[source] = h.values[source]
	})

	if len(result) == 0 {
		return nil
	}

	return result
}
