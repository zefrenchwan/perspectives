package entities

import (
	"iter"
	"time"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
)

// DynamicLink represents the direct neighborhood of an entity.
type DynamicLink interface {
	commons.Identifiable // Identifiable to get the id of the link.
	// Source defines the base entity
	Source() Entity
	// Role returns the role for that source to the link
	Role() string
	// Neighbors allows an iterator over the entities linked to this entity.
	// Each value is a pair : entity (second) and its link period association.
	Neighbors() iter.Seq2[periods.Period, Entity]
}

// Linkable defines the ability to link an entity with others.
// Linkable applies to entities.
// For instance, it alldws to define links such as :
// Knows(subject=Marie, object=Likes(subject=Lisa, object=chocolate))
type Linkable interface {
	// Links allows an iterator over the links as role and actual link value.
	Links() iter.Seq2[string, DynamicLink]
	// Link returns the entity associated with the given role (if it exists), for that period.
	Link(string) (DynamicLink, bool)
	// LinksAt returns the entities associated with the given roles for that moment.
	LinksAt(moment time.Time) iter.Seq2[string, Entity]
	// LinkAt returns the entity associated with the given role (if it exists), for that moment.
	LinkAt(role string, moment time.Time) (iter.Seq[Entity], bool)
}
