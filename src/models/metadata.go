package models

import (
	"github.com/google/uuid"
	"github.com/zefrenchwan/perspectives.git/structures"
)

// FormalAttribute defines the main characteristics of an attribute.
// So far, it is a name (for instance "weblink") and some semantic tags
type FormalAttribute struct {
	// Name of the attribute within a class
	Name string
	// Semantics of the attribute to flag a specific meaning.
	// For instance "email", "social network account", etc
	Semantics []string
}

// FormalClass is a formal class that would match the traits in objects
type FormalClass struct {
	// Id of the class
	Id string
	// Name of the class, should be unique per topic
	Name string
	// Attributes of the class
	Attributes map[string]FormalAttribute
}

// FormalRelation is a formal link definition.
// It defines relations main chararestics: how to use it (transitive ? Symetric ?) and matching roles.
// Roles may be: "subject", "location", etc.
type FormalRelation struct {
	// Id of the formal relation
	Id string
	// Name of the relation
	Link string
	// Transitive means R(a,b) and R(b,c) implies R(a,c)
	Transitive bool
	// Symetric means R(a,b) equivalent to R(b,a)
	Symetric bool
	// Roles may be subject, object, location, etc
	Roles map[string][]string
}

// possible roles
// RELATION_ROLE_SUBJECT for a subect (for instance Love(Subject, ...))
const RELATION_ROLE_SUBJECT = "subject"

// RELATION_ROLE_LOCATION for a location (for instance Eat(..., location, ...))
const RELATION_ROLE_LOCATION = "location"

// RELATION_ROLE_OBJECT for an object (for instance Eat(..., what (object), ...))
const RELATION_ROLE_OBJECT = "object"

// RELATION_ROLE_TIME for a time, duration (for instance, Meet(..,..., time,...))
const RELATION_ROLE_TIME = "time"

// NewFormalClass returns a new class
func NewFormalClass(name string) FormalClass {
	return FormalClass{
		Id:         uuid.NewString(),
		Name:       name,
		Attributes: make(map[string]FormalAttribute),
	}
}

// NewFormalCompleteClass returns a new class with attributes set
func NewFormalCompleteClass(name string, attributes map[string][]string) FormalClass {
	base := NewFormalClass(name)

	for attr, values := range attributes {
		base.SetAttribute(attr, values)
	}

	return base
}

// SetAttribute sets the values for that attribute (upsert)
func (c *FormalClass) SetAttribute(name string, tags []string) {
	var semantic []string
	if len(tags) != 0 {
		semantic = structures.SliceReduce(tags)
	}

	c.Attributes[name] = FormalAttribute{
		Name:      name,
		Semantics: semantic,
	}
}
