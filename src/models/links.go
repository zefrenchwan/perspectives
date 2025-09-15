package models

import (
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/zefrenchwan/perspectives.git/structures"
)

// linkValue is just a container for values in link.
// Here is the problem:
// We want to have an unique representation of a node in a link even if values matches.
// For instance, consider the link:
// Is(Cheese, Cheese)
// Values are the same, nodes are not
type linkValue struct {
	// uniqueId, assumed to be unique in the link
	uniqueId string
	// content is the actual value of the node in the link
	content ModelEntity
}

// contentType returns the type of the underlying content
func (v linkValue) contentType() EntityType {
	return v.content.GetType()
}

// newLinkValue builds a new node in a link based on a content
func newLinkValue(content ModelEntity) linkValue {
	return linkValue{
		uniqueId: uuid.NewString(),
		content:  content,
	}
}

// newLinkValueForObjects builds a link value as a group of objects
func newLinkValueForObjects(values []Object) linkValue {
	return linkValue{uuid.NewString(), objectsGroup(values)}
}

// Link will link objects together (0 level links) or links and object (higher level links).
// For instance Likes(Steve, Tiramisu) is a basic link and Knows(Paul, Likes(Steve, Tiramisu)) is an higher level link.
type Link struct {
	// id of the link
	id string
	// name defines the link semantic
	name string
	// operands are role based operands.
	// Usually, roles are "subject" or "object" or ...
	operands map[string]linkValue
	// Lifetime is the duration of the link
	lifetime structures.Period
}

// GetType returns EntityTypeLink
func (l *Link) GetType() EntityType {
	return EntityTypeLink
}

// AsLink returns the link
func (l *Link) AsLink() (*Link, error) {
	return l, nil
}

// AsGroup raises an error
func (l *Link) AsGroup() ([]Object, error) {
	return nil, errors.ErrUnsupported
}

// AsObject raises an error
func (l *Link) AsObject() (*Object, error) {
	return nil, errors.ErrUnsupported
}

// AsTrait raises an error
func (l *Link) AsTrait() (Trait, error) {
	return Trait{}, errors.ErrUnsupported
}

// AsVariable raises an error
func (l *Link) AsVariable() (Variable, error) {
	return Variable{}, errors.ErrUnsupported
}

// RoleSubject is the constant value for the subject role
const RoleSubject = "subject"

// RoleObject is the constant value for the object role
const RoleObject = "object"

// NewLink builds a link, valid for a given period
// name is the semantic of that link (for instance "loves" or "knows")
// values are the values (role linked to operand)
// duration is the period the link is valid for
//
// Although formal parameter is any, expected types are:
// 1. Slices of objects
// 2. Objects
// 3. Links
// 4. Traits
// 5. Variables representing previous mentioned types
//
// An error will raise if values do not match that constraint
func NewLink(name string, values map[string]any, duration structures.Period) (Link, error) {
	var link, empty Link
	link.id = uuid.NewString()
	link.name = name
	link.operands = make(map[string]linkValue)
	link.lifetime = duration

	for role, operand := range values {
		if operand == nil {
			continue
		} else if l, ok := operand.(Link); ok {
			link.operands[role] = newLinkValue(&l)
		} else if g, ok := operand.([]Object); ok {
			link.operands[role] = newLinkValueForObjects(g)
		} else if o, ok := operand.(Object); ok {
			link.operands[role] = newLinkValue(&o)
		} else if t, ok := operand.(Trait); ok {
			link.operands[role] = newLinkValue(t)
		} else {
			return empty, fmt.Errorf("unsupported type for role %s. Expecting either trait or object or link or group of objects", role)
		}
	}

	return link, nil
}

// NewSimpleLink is a shortcut to declare a link(subject, object) valid for the full time
func NewSimpleLink(link string, subject, object any) (Link, error) {
	return NewLink(link, map[string]any{RoleSubject: subject, RoleObject: object}, structures.NewFullPeriod())
}

// Id returns the globally unique id for that link
func (l *Link) Id() string {
	return l.id
}

// Name returns the name of the link
func (l *Link) Name() string {
	return l.name
}

// findAllMatchingCondition goes through the full link and find elements matching condition
func (l *Link) findAllMatchingCondition(acceptance func(ModelEntity) bool) []ModelEntity {
	matches := make([]ModelEntity, 0)
	linksAlreadyVisited := make(map[string]bool)

	elements := []*Link{l}
	for len(elements) != 0 {
		current := elements[0]
		elements = elements[1:]

		// STEP ONE: DEAL WITH THE WALKTHROUGH
		if current.GetType() == EntityTypeLink {
			link, _ := current.AsLink()
			if linksAlreadyVisited[link.id] {
				continue
			} else {
				linksAlreadyVisited[link.id] = true
			}

			for _, value := range link.operands {
				if value.contentType() == EntityTypeLink {
					childLlink, _ := value.content.AsLink()
					elements = append(elements, childLlink)
				}
			}
		}
		// END OF WALKTHROUGH

		// STEP TWO: TEST MATCH AND ADD IN MATCHES ACCORDINGLY
		if acceptance(current) {
			matches = append(matches, current)
		}

	}

	return matches
}

// AllObjectsOperands returns the objects appearing recursively in the link.
// It means that if l is a link of links of objects, descendants objects will appear.
// Each object appears once per id
func (l *Link) AllObjectsOperands() []Object {
	acceptValueAsObject := func(v ModelEntity) bool {
		matchingTypes := []EntityType{EntityTypeGroup, EntityTypeObject}
		return slices.Contains(matchingTypes, v.GetType())
	}

	var matches []Object
	values := l.findAllMatchingCondition(acceptValueAsObject)
	for _, value := range values {
		switch value.GetType() {
		case EntityTypeGroup:
			g, _ := value.AsGroup()
			matches = append(matches, g...)
		case EntityTypeObject:
			o, _ := value.AsObject()
			matches = append(matches, *o)
		}
	}

	return structures.SliceDeduplicateFunc(matches, func(a, b Object) bool { return a.Id == b.Id })
}

// LocalLinkValueMapper defines a mapping from a value to another.
// Accepted transformations are:
// IF value is anything but a link, THEN its image is also anything but a link
type LocalLinkValueMapper func(linkValue) (linkValue, bool, error)

// localLinkCaller calls a mapper but ensures invariants are respected
func localLinkValueCaller(value linkValue, mapper LocalLinkValueMapper) (linkValue, bool, error) {
	return mapper(value)
}
