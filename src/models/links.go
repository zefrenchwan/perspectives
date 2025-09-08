package models

import (
	"errors"
	"fmt"
	"maps"

	"github.com/google/uuid"
	"github.com/zefrenchwan/perspectives.git/structures"
)

// Link will link objects together (0 level links) or links and object (higher level links).
// For instance Likes(Steve, Tiramisu) is a basic link and Knows(Paul, Likes(Steve, Tiramisu)) is an higher level link.
type Link struct {
	// id of the link
	id string
	// name defines the link semantic
	name string
	// operands are role based operands.
	// Usually, roles are "subject" or "object" or ...
	operands map[string]LinkValue
}

// LinkValueType defines the accepted types of a link value.
// So far, accepted types are:
// objects: for instance: John knows Jane
// traits: for instance John Likes Chocolate (with Chocolate a trait)
// links: for instance John knows (Marie likes Chocolate)
// groups: for instance, Mary and John (as a group) like Chocolate
type LinkValueType int

// LinkValueAsTrait says that operand is a trait
const LinkValueAsTrait = 1

// LinkValueAsLink says that operand is a link
const LinkValueAsLink = 2

// LinkValueAsGroup says that operand is a group
const LinkValueAsGroup = 3

// LinkValueAsObject says that operand is an object
const LinkValueAsObject = 4

// LinkValue is the union type defintion of any operands
type LinkValue interface {
	// GetType returns the actual type of the value
	GetType() LinkValueType
	// AsLink casts the value as a link, or raises an error it underlying content is not a link
	AsLink() (Link, error)
	// AsGroup casts the value as a group of objects, or raises an error it underlying content is not a group
	AsGroup() ([]Object, error)
	// AsObject casts the value as an object, or raises an error it underlying content is not an object
	AsObject() (Object, error)
	// AsTrait returns the value as a trait, or raises an error it underlying content is not a trait
	AsTrait() (Trait, error)
}

// LinkObject is an object as a link operand
type LinkObject Object

// LinkGroup is a group of objects as a link operand
type LinkGroup []Object

// LinkTrait is a trait as a link operand
type LinkTrait Trait

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
//
// An error will raise if values do not match that constraint
func NewLink(name string, values map[string]any, duration structures.Period) (Link, error) {
	var link, empty Link
	link.id = uuid.NewString()
	link.name = name
	link.operands = make(map[string]LinkValue)

	for role, operand := range values {
		if operand == nil {
			continue
		} else if l, ok := operand.(Link); ok {
			link.operands[role] = l
		} else if g, ok := operand.([]Object); ok {
			link.operands[role] = LinkGroup(g)
		} else if o, ok := operand.(Object); ok {
			link.operands[role] = LinkObject(o)
		} else if t, ok := operand.(Trait); ok {
			link.operands[role] = LinkTrait(t)
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
func (l Link) Id() string {
	return l.id
}

// Name returns the name of the link
func (l Link) Name() string {
	return l.name
}

// ValuesPerRole returns the role based map of values
func (l Link) ValuesPerRole() map[string]LinkValue {
	result := make(map[string]LinkValue)
	maps.Copy(result, l.operands)
	return result
}

// AllObjectsOperands returns the objects appearing recursively in the link.
// It means that if l is a link of links of objects, descendants objects will appear.
// Each object appears once per id
func (l Link) AllObjectsOperands() []Object {
	var result []Object
	matches := make(map[string]Object)
	linksAlreadyVisited := make(map[string]bool)

	elements := []Link{l}
	for len(elements) != 0 {
		current := elements[0]
		elements = elements[1:]

		if linksAlreadyVisited[current.id] {
			continue
		} else {
			linksAlreadyVisited[current.id] = true
		}

		for _, value := range current.ValuesPerRole() {
			switch value.GetType() {
			case LinkValueAsLink:
				l, _ := value.AsLink()
				if !linksAlreadyVisited[l.id] {
					elements = append(elements, l)
				}
			case LinkValueAsGroup:
				g, _ := value.AsGroup()
				for _, obj := range g {
					matches[obj.Id] = obj
				}
			case LinkValueAsObject:
				o, _ := value.AsObject()
				matches[o.Id] = o
			}
		}
	}

	for _, obj := range matches {
		result = append(result, obj)
	}

	return result
}

////////////////////////////////////////////////
// TECHNICAL IMPLEMENTATION OF LINKS OPERANDS //
////////////////////////////////////////////////

func (o LinkObject) GetType() LinkValueType {
	return LinkValueAsObject
}

func (o LinkObject) AsLink() (Link, error) {
	return Link{}, errors.New("invalid value: expecting link, got object")
}

func (o LinkObject) AsGroup() ([]Object, error) {
	return nil, errors.New("invalid value: expecting group, got object")
}

func (o LinkObject) AsObject() (Object, error) {
	return Object(o), nil
}

func (o LinkObject) AsTrait() (Trait, error) {
	var empty Trait
	return empty, errors.New("invalid value: expecting trait, got object")
}

func (g LinkGroup) GetType() LinkValueType {
	return LinkValueAsGroup
}

func (g LinkGroup) AsLink() (Link, error) {
	return Link{}, errors.New("invalid value: expecting link, got group")
}

func (g LinkGroup) AsGroup() ([]Object, error) {
	return []Object(g), nil
}

func (g LinkGroup) AsObject() (Object, error) {
	var object Object
	return object, errors.New("invalid value: expecting object, got group")
}

func (g LinkGroup) AsTrait() (Trait, error) {
	var empty Trait
	return empty, errors.New("invalid value: expecting trait, got group")
}

func (l Link) GetType() LinkValueType {
	return LinkValueAsLink
}

func (l Link) AsLink() (Link, error) {
	return l, nil
}

func (l Link) AsGroup() ([]Object, error) {
	return nil, errors.New("invalid value: expecting group, got link")
}

func (l Link) AsObject() (Object, error) {
	var object Object
	return object, errors.New("invalid value: expecting object, got link")
}

func (l Link) AsTrait() (Trait, error) {
	var trait Trait
	return trait, errors.New("invalid value: expecting trait, got link")
}

func (t LinkTrait) GetType() LinkValueType {
	return LinkValueAsTrait
}

func (t LinkTrait) AsLink() (Link, error) {
	var empty Link
	return empty, errors.New("invalid value: expecting link, got trait")
}

func (t LinkTrait) AsGroup() ([]Object, error) {
	return nil, errors.New("invalid value: expecting group, got link")
}

func (t LinkTrait) AsObject() (Object, error) {
	var object Object
	return object, errors.New("invalid value: expecting object, got link")
}

func (t LinkTrait) AsTrait() (Trait, error) {
	return Trait(t), nil
}
