package models

import "github.com/zefrenchwan/perspectives.git/structures"

// FormalHierarchy deals with class inheritance and links inheritance
type FormalHierarchy struct {
	// classesHierarchy is the inheritance tree for classes
	classesHierarchy structures.Hierarchy[FormalClass]
	// relationsHierarchy is the inheritance tree for links
	relationsHierarchy structures.Hierarchy[FormalRelation]
}

// NewFormalHierarchy returns a new empty formal hierarchy
func NewFormalHierarchy() FormalHierarchy {
	return FormalHierarchy{
		classesHierarchy:   structures.NewHierarchy[FormalClass](),
		relationsHierarchy: structures.NewHierarchy[FormalRelation](),
	}
}

// SetClass registers a class in the hierarchy
func (h *FormalHierarchy) SetClass(c FormalClass) {
	h.classesHierarchy.SetValue(c.Name, c)
}

// AddChildClass links existing child to existing parent.
// If exclusive, then childs exclude each other
func (h *FormalHierarchy) AddChildClass(childName, parentName string, exclusive bool) error {
	if exclusive {
		return h.classesHierarchy.AddChildInPartition(childName, parentName)
	} else {
		return h.classesHierarchy.AddChildToParent(childName, parentName)
	}
}

// GetClassHierarchy gets all super classes of a class by name.
// It contains at least the class if it matches a name in the hierarchy, no element otherwise.
// For instance, asking for "dogs" when hierarchy is "dogs" < "animals" would return dogs and animals classes
func (h *FormalHierarchy) GetClassHierarchy(name string) []FormalClass {
	if values, found := h.classesHierarchy.AncestorsValues(name); !found {
		return nil
	} else {
		return values
	}
}
