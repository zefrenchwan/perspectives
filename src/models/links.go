package models

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/zefrenchwan/perspectives.git/structures"
)

// Link will link objects together (0 level links) or links and object (higher level links).
// For instance Likes(Steve, Tiramisu) is a basic link and Knows(Paul, Likes(Steve, Tiramisu)) is an higher level link.
// There is a "any" type for operands.
// Still, accepted values are objects, slices of objects (groups), or links.
// Recursive definition meant an implementation using any.
// Other solutions were tested but ease of use meant using "any"
type Link struct {
	// id of the link
	id string
	// name defines the link semantic
	name string
	// operands are role based operands.
	// Usually, roles are "subject" or "object" or ...
	operands map[string]any
}

// LinkValue is the union type defintion of any operands
type LinkValue interface {
	// IsLink returns true for links
	IsLink() bool
	// IsGroup returns true for groups
	IsGroup() bool
	// IsObject returns true for objects
	IsObject() bool
	// AsLink casts the value as a link, or raises an error it underlying content is not a link
	AsLink() (Link, error)
	// AsGroup casts the value as a group of objects, or raises an error it underlying content is not a group
	AsGroup() ([]Object, error)
	// AsObject casts the value as an object, or raises an error it underlying content is not an object
	AsObject() (Object, error)
}

// LinkObject is an object as a link operand
type LinkObject Object

// LinkGroup is a group of objects as a link operand
type LinkGroup []Object

// NewLink builds a link, valid for a given period
// name is the semantic of that link (for instance "loves" or "knows")
// values are the values (role linked to operand)
// duration is the period the link is valid for
//
// Although formal parameter is any, expected types are:
// 1. Slices of objects
// 2. Objects
// 3. Links
//
// An error will raise if values do not match that constraint
func NewLink(name string, values map[string]any, duration structures.Period) (Link, error) {
	var link, empty Link
	link.id = uuid.NewString()
	link.name = name
	link.operands = make(map[string]any)

	for role, operand := range values {
		if operand == nil {
			continue
		} else if l, ok := operand.(Link); ok {
			link.operands[role] = l
		} else if g, ok := operand.([]Object); ok {
			link.operands[role] = g
		} else if o, ok := operand.(Object); ok {
			link.operands[role] = o
		} else {
			return empty, fmt.Errorf("unsupported type for role %s. Expecting either object or link or group of objects", role)
		}
	}

	return link, nil
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
	for role, operand := range l.operands {
		if operand == nil {
			continue
		} else if l, ok := operand.(Link); ok {
			result[role] = l
		} else if g, ok := operand.([]Object); ok {
			result[role] = LinkGroup(g)
		} else if o, ok := operand.(Object); ok {
			result[role] = LinkObject(o)
		}
	}

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
			switch {
			case value.IsLink():
				l, _ := value.AsLink()
				if !linksAlreadyVisited[l.id] {
					elements = append(elements, l)
				}
			case value.IsGroup():
				g, _ := value.AsGroup()
				for _, obj := range g {
					matches[obj.Id] = obj
				}
			case value.IsObject():
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

func (o LinkObject) IsLink() bool {
	return false
}

func (o LinkObject) IsGroup() bool {
	return false
}

func (o LinkObject) IsObject() bool {
	return true
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

func (o LinkGroup) IsLink() bool {
	return false
}

func (o LinkGroup) IsGroup() bool {
	return true
}

func (o LinkGroup) IsObject() bool {
	return false
}

func (o LinkGroup) AsLink() (Link, error) {
	return Link{}, errors.New("invalid value: expecting link, got group")
}

func (o LinkGroup) AsGroup() ([]Object, error) {
	return []Object(o), nil
}

func (o LinkGroup) AsObject() (Object, error) {
	var object Object
	return object, errors.New("invalid value: expecting object, got group")
}

func (o Link) IsLink() bool {
	return true
}

func (o Link) IsGroup() bool {
	return false
}

func (o Link) IsObject() bool {
	return false
}

func (o Link) AsLink() (Link, error) {
	return o, nil
}

func (o Link) AsGroup() ([]Object, error) {
	return nil, errors.New("invalid value: expecting group, got link")
}

func (o Link) AsObject() (Object, error) {
	var object Object
	return object, errors.New("invalid value: expecting object, got link")
}
