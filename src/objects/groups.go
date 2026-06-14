package objects

import (
	"slices"
	"strings"

	"github.com/zefrenchwan/perspectives.git/commons"
)

// SetOfInstances defines an immutable collection of unique Instance elements.
// It implements the Element interface and guarantees that instances are deduplicated.
type SetOfInstances interface {
	Element
	// Size returns the number of instances in the collection.
	Size() int
	// Range iterates over all instances in the collection, yielding each instance to a provided function.
	// Iteration stops if the yield function returns false.
	Range(func(Instance) bool)
	// Contains returns true if the collection contains the given instance.
	Contains(Instance) bool
	// SortedInstances returns a sorted slice of all instances in the collection.
	SortedInstances() []Instance
}

// baseCollection is the standard in-memory implementation of the SetOfInstances interface.
type baseCollection struct {
	// id is the unique identifier of the collection, computed from its deduplicated instances.
	id string
	// elements maps instance IDs to their corresponding Instance objects for O(1) access.
	elements map[string]Instance
}

// Id returns the unique identifier for the collection.
func (b *baseCollection) Id() string {
	return b.id
}

// DeclaringClass returns the class for this collection: a CLASS_INSTANCES_COLLECTION.
func (b *baseCollection) DeclaringClass() Class {
	return CLASS_INSTANCES_COLLECTION
}

// isLinkable uses the sealed pattern to ensure that baseCollection instances can satisfy the Element interface requirements.
func (b *baseCollection) isLinkable() bool {
	return true
}

// Same returns true if the collection is functionally equivalent to the other element:
// same class, same size, and containing the exact same instances.
func (b *baseCollection) Same(other Element) bool {
	if b == nil && other == nil {
		return true
	} else if b == nil || other == nil {
		return false
	} else if other.DeclaringClass() != CLASS_INSTANCES_COLLECTION {
		return false
	}

	otherCollection, ok := other.(SetOfInstances)
	if !ok {
		return false
	}

	// Fast path: if sizes differ, they cannot be the same.
	if b.Size() != otherCollection.Size() {
		return false
	}

	// Check if all elements in the other collection exist in this one.
	for instance := range otherCollection.Range {
		if instance == nil {
			return false
		} else if _, found := b.elements[instance.Id()]; !found {
			return false
		}
	}

	return true
}

// SortedInstances returns a sorted slice of all instances within the collection.
// The slice is sorted by instance ID to ensure idempotent results.
func (b *baseCollection) SortedInstances() []Instance {
	if b == nil || len(b.elements) == 0 {
		return nil
	}

	var elements []Instance
	for _, element := range b.elements {
		elements = append(elements, element)
	}

	slices.SortFunc(elements, func(a, b Instance) int {
		return strings.Compare(a.Id(), b.Id())
	})

	return elements
}

// Size returns the total number of unique instances within the collection.
func (b *baseCollection) Size() int {
	if b == nil {
		return 0
	}
	return len(b.elements)
}

// Contains checks whether the specified instance exists within the collection.
func (b *baseCollection) Contains(i Instance) bool {
	if i == nil || b == nil {
		return false
	}
	return b.elements[i.Id()] != nil
}

// Range iterates over all instances in the collection and yields each to the provided function.
func (b *baseCollection) Range(yield func(Instance) bool) {
	if b == nil {
		return
	}
	for _, instance := range b.elements {
		if !yield(instance) {
			return
		}
	}
}

// NewSetOfInstances creates a new immutable collection of instances.
// It deduplicates the provided instances and generates a stable ID based on the unique instances.
func NewSetOfInstances(instances []Instance) SetOfInstances {
	elements := make(map[string]Instance)

	// Step 1: Populate the map to guarantee deduplication and ignore nil values.
	for _, instance := range instances {
		if instance == nil {
			continue
		}
		elements[instance.Id()] = instance
	}

	// Step 2: Generate the ID from the deduplicated elements.
	allIds := make([]string, 0, len(elements))
	for id := range elements {
		allIds = append(allIds, id)
	}

	slices.Sort(allIds)
	commonId := commons.HashString(strings.Join(allIds, ","))

	return &baseCollection{
		id:       commonId,
		elements: elements,
	}
}
