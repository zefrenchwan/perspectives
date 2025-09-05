package models

import (
	"errors"
	"maps"

	"github.com/google/uuid"
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

// ObjectsOperands returns all the objects the term uses.
// One object appears once as a pointer (deduplication)
func (r RelationTerm) ObjectsOperands() []*Object {
	matches := make(map[string]*Object)
	seen := make(map[*Relation]bool)
	var relations []*Relation

	if r.relation != nil {
		relations = append(relations, r.relation)
	} else {
		return structures.SliceDeduplicate(r.objects)
	}

	for len(relations) != 0 {
		current := relations[0]
		relations = relations[1:]
		seen[current] = true
		for _, term := range current.Parameters {
			if term.objects != nil {
				for _, object := range term.objects {
					matches[object.Id] = object
				}
			} else if term.relation != nil {
				if !seen[term.relation] {
					relations = append(relations, term.relation)
				}
			}
		}
	}

	var objects []*Object
	for _, val := range matches {
		objects = append(objects, val)
	}

	return objects
}

// Relation is a tree with nodes being terms and links being roles.
// For instance Likes(John, Cheese) is a relation.
// But Knows(John, Likes(Marie, Cheese)) is a relation too
type Relation struct {
	// Id defines an unique relation
	Id string
	// Link defines the semantic of the relation
	Link string
	// Parameters links roles to terms.
	Parameters map[string]RelationTerm
	// Lifetime defines the period during the relation is true
	Lifetime structures.Period
}

// NewRelation builds a new relation.
// A relation links content as roles and values for a given duration
func NewRelation(link string, parameters map[string]RelationTerm, duration structures.Period) Relation {
	relation := Relation{
		Id:         uuid.NewString(),
		Link:       link,
		Parameters: make(map[string]RelationTerm),
		Lifetime:   duration,
	}

	maps.Copy(relation.Parameters, parameters)
	return relation
}

// AsRelationTerm returns a term for that relation.
// Use it to include this relation again as a term to compose
func (r Relation) AsRelationTerm() RelationTerm {
	return RelationTerm{
		relation: &r,
	}
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
	relation := NewRelation(link, parameters, duration)

	return RelationTerm{
		relation: &relation,
	}
}
