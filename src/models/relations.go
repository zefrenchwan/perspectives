package models

import (
	"errors"
	"maps"

	"github.com/zefrenchwan/perspectives.git/structures"
)

// RelationTerm is an union type for a relation content:
// Either objects (as a group) for operands
// Or a relation (verb + roles linked to parameters)
type RelationTerm struct {
	// objects contains the objects as a term. If set, relation should be nil
	objects []*Object
	// relation contains the relation as a term. If set, objects should be nil
	relation *Relation
}

// Relation is a tree with nodes being terms and links being roles.
// For instance Likes(John, Cheese) is a relation.
// But Knows(John, Likes(Marie, Cheese)) is a relation too
type Relation struct {
	// Link defines the semantic of the relation
	Link string
	// Parameters links roles to terms.
	Parameters map[string]RelationTerm
	// Lifetime defines the period during the relation is true
	Lifetime structures.Period
}

// AsObjects returns the objects within the term (may be nil)
func (t RelationTerm) AsObjects() []*Object {
	return t.objects
}

// Build returns a relation from that term
// It raises an error if term is an object (or a group)
func (t RelationTerm) Build() (*Relation, error) {
	if t.relation == nil {
		return nil, errors.New("term is only an object, no relation")
	}

	return t.relation, nil
}

// NewObjectTerm returns a term for that single object
func NewObjectTerm(object Object) RelationTerm {
	return RelationTerm{
		objects: []*Object{&object},
	}
}

// NewGroupTerm builds a group of objects as a term.
// For instance "Marie and John like pudding" would be represented as:
// Like(Group ( {Marie, John} ), Pudding)
func NewGroupTerm(objects []Object) RelationTerm {
	var result RelationTerm
	result.objects = make([]*Object, 0)
	for _, obj := range objects {
		result.objects = append(result.objects, &obj)
	}

	return result
}

// NewRelationTerm builds a new relation (as a term) from a link, roles and true for a given duration
func NewRelationTerm(link string, parameters map[string]RelationTerm, duration structures.Period) RelationTerm {
	relation := Relation{Link: link, Parameters: make(map[string]RelationTerm), Lifetime: duration}
	maps.Copy(relation.Parameters, parameters)

	return RelationTerm{
		relation: &relation,
	}
}
