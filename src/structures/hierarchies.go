package structures

import (
	"errors"
	"fmt"
)

// parent_link defines an "extends" relation
const parent_link = 10

// child_exclusive_link defines a link that says childs form a partition of parent.
// For instance, man and woman are exclusive childs of person
const child_exclusive_link = 100

// child_standard_link defines a link that says that childs are linked to parents.
// For instance, people may be worker or father (and both)
const child_standard_link = 1000

// Hierarchy is a DAG of named S , link is extends.
// Assumption is that a name is unique
type Hierarchy[S any] struct {
	// parents link childs to parents
	parents DVGraph[string, int]
	// values link instances of S by name
	values map[string]S
	// childs link parents to childs via a link : partition or standard
	childs DVGraph[string, int]
}

// NewHierarchy builds an empty hierarchy
func NewHierarchy[S any]() Hierarchy[S] {
	return Hierarchy[S]{
		parents: NewDVGraph[string, int](),
		values:  make(map[string]S),
		childs:  NewDVGraph[string, int](),
	}
}

// SetValue sets value to a key (a name)
func (h Hierarchy[S]) SetValue(key string, value S) {
	h.values[key] = value
	h.parents.AddNode(key)
}

// GetValue returns the value associated to that key, if any.
// It returns the value, true if found, empty, false otherwise
func (h Hierarchy[S]) GetValue(key string) (S, bool) {
	var empty S
	if v, found := h.values[key]; found {
		return v, found
	}
	return empty, false
}

// AddChildToParent adds a child to a parent
// If other childs form a partition, raise an error
func (h Hierarchy[S]) AddChildToParent(child, parent string) error {
	return h.addLink(child, parent, child_standard_link)
}

// AddChildInPartition adds a child in a partition
// If other childs form a standard union, raise an error
func (h Hierarchy[S]) AddChildInPartition(child, parent string) error {
	return h.addLink(child, parent, child_exclusive_link)
}

// addLink is private method to deal with links.
// It adds a
func (h Hierarchy[S]) addLink(child, parent string, link int) error {
	// test if values exist or not
	_, sourceExists := h.values[child]
	_, destinationExists := h.values[parent]
	if child == parent {
		return errors.New("cannot link element with itself")
	} else if !sourceExists {
		return errors.New("child does not exist")
	} else if !destinationExists {
		return errors.New("parent does not exist")
	}

	// at this point, parent and child exist, but link may not
	existingLinks, foundLinks := h.childs.Neighbors(parent)
	// if some previous links existed, ensure that new link has the same type
	var existingType int
	if foundLinks {
		// test if link is the same as the previous ones
		// THAT IS: don't add partition to simple union, and vice versa
		for _, existingLink := range existingLinks {
			existingType = existingLink
			break
		}

		// test if type of link to add matches the current type of the existing links
		if existingType != link {
			return fmt.Errorf("cannot create link with a different type than existing values for %s", child)
		}
	}

	// Add the link in the parents tree
	addedLink := h.parents.LinkWithoutCycle(child, parent, parent_link)
	if !addedLink {
		return fmt.Errorf("there would be a cycle linking %s to %s", child, parent)
	}

	// no link detected, so we may add
	h.childs.Link(parent, child, link)

	return nil
}

// Ancestors returns the ancestors of a value as a slice.
// It returns the ancestors (including the node), true or nil, false for not found
func (h Hierarchy[S]) Ancestors(source string) ([]string, bool) {
	var result []string
	h.parents.Walk(source, func(source string) {
		result = append(result, source)
	})

	return result, result != nil
}

// AncestorsValues return the ancestors values from a given source
func (h Hierarchy[S]) AncestorsValues(source string) ([]S, bool) {
	var result []S
	if ancestors, found := h.Ancestors(source); !found {
		return nil, false
	} else {
		for _, name := range ancestors {
			value := h.values[name]
			result = append(result, value)
		}

		return result, true
	}
}

// Childs returns the childs of a node by name
// Result is the cbilds (if any) and a boolean
// This boolean is true if childs are mutually exclusive
func (h Hierarchy[S]) Childs(source string) ([]string, bool) {
	if links, found := h.childs.Neighbors(source); !found {
		return nil, false
	} else {
		var result []string
		var exclusive bool
		for child, linkType := range links {
			exclusive = linkType == child_exclusive_link
			result = append(result, child)
		}

		return result, exclusive
	}
}
