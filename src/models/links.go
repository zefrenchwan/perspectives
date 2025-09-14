package models

import (
	"errors"
	"fmt"
	"maps"
	"slices"

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

// Equals returns true for same links based on id
func (l Link) Equals(other Link) bool {
	return l.id == other.id
}

// LinkValueType defines the accepted types of a link value.
// So far, accepted types are:
// objects: for instance: John knows Jane
// traits: for instance John Likes Chocolate (with Chocolate a trait)
// links: for instance John knows (Marie likes Chocolate)
// groups: for instance, Mary and John (as a group) like Chocolate
// variables: to be replaced by any previous type
type LinkValueType int

// LinkValueAsTrait says that operand is a trait
const LinkValueAsTrait = 1

// LinkValueAsLink says that operand is a link
const LinkValueAsLink = 2

// LinkValueAsGroup says that operand is a group
const LinkValueAsGroup = 3

// LinkValueAsObject says that operand is an object
const LinkValueAsObject = 4

// LinkValueAsVariable says that operand is a variable
const LinkValueAsVariable = 5

// RoleSubject is the constant value for the subject role
const RoleSubject = "subject"

// RoleObject is the constant value for the object role
const RoleObject = "object"

// LinkValue is the union type defintion of any operands
type LinkValue interface {
	// UniqueId returns the unique id of the value in the link
	UniqueId() string
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
	// AsVariable returns the value as a variable, or raises an error if underlying content is not a variable
	AsVariable() (LinkVariable, error)
}

// LinkObject is an object as a link operand
type LinkObject struct {
	id      string
	content Object
}

// NewLinkObject builds a new link object to decorate an object
func NewLinkObject(value Object) LinkObject {
	return LinkObject{id: uuid.NewString(), content: value}
}

// LinkGroup is a group of objects as a link operand
type LinkGroup struct {
	id      string
	content []Object
}

// NewLinkGroup decorates values as a link group
func NewLinkGroup(values []Object) LinkGroup {
	return LinkGroup{id: uuid.NewString(), content: values}
}

// LinkTrait is a trait as a link operand
type LinkTrait struct {
	id      string
	content Trait
}

func NewLinkTrait(trait Trait) LinkTrait {
	return LinkTrait{id: uuid.NewString(), content: trait}
}

// LinkVariable defines a variable that may be replaced by any other link value
type LinkVariable struct {
	// id as an unique id for that variable
	id string
	// Name of the variable (usually "x","y","z")
	name string
	// ValidTypes are the union of the valid types to replace this variable for
	validTypes []LinkValueType
	// ValidTraits contain the union of acceptable traits
	validTraits []Trait
}

// NewLinkVariableForObject returns a new variable for that object
func NewLinkVariableForObject(name string, traits []string) LinkVariable {
	var matches []Trait
	for _, trait := range structures.SliceReduce(traits) {
		matches = append(matches, NewTrait(trait))
	}

	return LinkVariable{
		id:          uuid.NewString(),
		name:        name,
		validTypes:  []LinkValueType{LinkValueAsObject},
		validTraits: matches,
	}
}

// NewLinkVariableForGroup returns the variable for a group of objects matching traits
func NewLinkVariableForGroup(name string, traits []string) LinkVariable {
	var matches []Trait
	for _, trait := range structures.SliceReduce(traits) {
		matches = append(matches, NewTrait(trait))
	}

	return LinkVariable{
		id:          uuid.NewString(),
		name:        name,
		validTypes:  []LinkValueType{LinkValueAsGroup},
		validTraits: matches,
	}
}

// NewLinkVariableForTrait returns a new variable for that trait
func NewLinkVariableForTrait(name string) LinkVariable {
	return LinkVariable{
		id:          uuid.NewString(),
		name:        name,
		validTypes:  []LinkValueType{LinkValueAsTrait},
		validTraits: nil,
	}
}

// NewLinkVariableForSpecificTraits returns a variable for trait that may take only specific values
func NewLinkVariableForSpecificTraits(name string, traits []string) LinkVariable {
	var matches []Trait
	for _, trait := range structures.SliceReduce(traits) {
		matches = append(matches, NewTrait(trait))
	}

	return LinkVariable{
		id:          uuid.NewString(),
		name:        name,
		validTypes:  []LinkValueType{LinkValueAsTrait},
		validTraits: matches,
	}
}

// NewLinkVariableForLink returns a new variable for that link
func NewLinkVariableForLink(name string) LinkVariable {
	return LinkVariable{
		id:          uuid.NewString(),
		name:        name,
		validTypes:  []LinkValueType{LinkValueAsLink},
		validTraits: nil,
	}
}

// MatchesTraits returns true if traits match accepted traits for that variable
func (lv LinkVariable) MatchesTraits(traits []Trait) bool {
	if len(lv.validTraits) == 0 {
		// no prerequisite
		return true
	} else {
		// prerequisites are set, and then should match
		return structures.SliceCommonElementFunc(lv.validTraits, traits, func(a, b Trait) bool { return a.Equals(b) })
	}
}

// MapAs transforms a variable to a value.
// Accepted values for other are the same as the link values, except variables.
// That is: slices of objects, objects, traits, links and related link values
func (lv LinkVariable) MapAs(other any) (LinkValue, error) {
	var empty LinkValue
	expectedTypes := lv.validTypes

	if v, ok := other.(Object); ok {
		if !slices.Contains(expectedTypes, LinkValueAsObject) {
			return empty, errors.New("object does not match expected type")
		}

		// test if object matches the definition.
		// Accept if there is a matching trait
		if lv.MatchesTraits(v.traits) {
			return NewLinkObject(v), nil
		} else {
			return empty, errors.New("no matching trait compatible with type definition")
		}
	} else if v, ok := other.(LinkObject); ok {
		if !slices.Contains(expectedTypes, LinkValueAsObject) {
			return empty, errors.New("object does not match expected type")
		}

		// test if object matches the definition.
		// Accept if there is a matching trait
		if lv.MatchesTraits(v.content.traits) {
			return v, nil
		} else {
			return empty, errors.New("no matching trait compatible with type definition")
		}

	} else if v, ok := other.([]Object); ok {
		if !slices.Contains(expectedTypes, LinkValueAsGroup) {
			return empty, errors.New("group does not match expected type")
		}

		// test if each object within the group matches the trait condition
		for index, obj := range v {
			if !lv.MatchesTraits(obj.traits) {
				return empty, fmt.Errorf("value at index %d does not match traits condition", index)
			}
		}

		return NewLinkGroup(v), nil
	} else if v, ok := other.(LinkGroup); ok {
		if !slices.Contains(expectedTypes, LinkValueAsGroup) {
			return empty, errors.New("group does not match expected type")
		}

		// test if each object within the group matches the trait condition
		for index, obj := range v.content {
			if !lv.MatchesTraits(obj.traits) {
				return empty, fmt.Errorf("value at index %d does not match traits condition", index)
			}
		}

		return v, nil
	} else if v, ok := other.(LinkTrait); ok {
		if !slices.Contains(expectedTypes, LinkValueAsTrait) {
			return empty, errors.New("group does not match expected type")
		}

		// for traits, either variable does not specify any, or traits match
		if len(lv.validTraits) == 0 {
			return v, nil
		} else if !slices.ContainsFunc(lv.validTraits, func(t Trait) bool { return v.content.Equals(t) }) {
			return empty, errors.New("trait value does not match expected traits")
		} else {
			return v, nil
		}
	} else if v, ok := other.(Trait); ok {
		if !slices.Contains(expectedTypes, LinkValueAsTrait) {
			return empty, errors.New("group does not match expected type")
		}

		// for traits, either variable does not specify any, or traits match
		if len(lv.validTraits) == 0 {
			return NewLinkTrait(v), nil
		} else if !slices.ContainsFunc(lv.validTraits, func(t Trait) bool { return v.Equals(t) }) {
			return empty, errors.New("trait value does not match expected traits")
		} else {
			return NewLinkTrait(v), nil
		}
	} else if v, ok := other.(Link); ok {
		if !slices.Contains(expectedTypes, LinkValueAsLink) {
			return empty, errors.New("group does not match expected type")
		} else {
			return v, nil
		}
	} else {
		return empty, errors.New("invalid value to map")
	}
}

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
	link.operands = make(map[string]LinkValue)

	for role, operand := range values {
		if operand == nil {
			continue
		} else if l, ok := operand.(Link); ok {
			link.operands[role] = l
		} else if g, ok := operand.([]Object); ok {
			link.operands[role] = NewLinkGroup(g)
		} else if o, ok := operand.(Object); ok {
			link.operands[role] = NewLinkObject(o)
		} else if t, ok := operand.(Trait); ok {
			link.operands[role] = NewLinkTrait(t)
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

// FindAllMatchingCondition goes through the full link and find elements matching condition
func (l Link) FindAllMatchingCondition(acceptance func(LinkValue) bool) []LinkValue {
	matches := make([]LinkValue, 0)
	linksAlreadyVisited := make(map[string]bool)

	elements := []LinkValue{l}
	for len(elements) != 0 {
		current := elements[0]
		elements = elements[1:]

		// STEP ONE: DEAL WITH THE WALKTHROUGH
		if current.GetType() == LinkValueAsLink {
			link, _ := current.AsLink()
			if linksAlreadyVisited[link.id] {
				continue
			} else {
				linksAlreadyVisited[link.id] = true
			}

			for _, value := range link.ValuesPerRole() {
				elements = append(elements, value)
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
func (l Link) AllObjectsOperands() []Object {
	acceptValueAsObject := func(v LinkValue) bool {
		matchingTypes := []LinkValueType{LinkValueAsGroup, LinkValueAsObject}
		return slices.Contains(matchingTypes, v.GetType())
	}

	var matches []Object
	values := l.FindAllMatchingCondition(acceptValueAsObject)
	for _, value := range values {
		switch value.GetType() {
		case LinkValueAsGroup:
			g, _ := value.AsGroup()
			matches = append(matches, g...)
		case LinkValueAsObject:
			o, _ := value.AsObject()
			matches = append(matches, o)
		}
	}

	return structures.SliceDeduplicateFunc(matches, func(a, b Object) bool { return a.Id == b.Id })
}

// LocalLinkValueMapper defines a mapping from a value to another.
// Accepted transformations are:
// IF value is anything but a link, THEN its image is also anything but a link
type LocalLinkValueMapper func(LinkValue) (LinkValue, bool, error)

// localLinkCaller calls a mapper but ensures invariants are respected
func localLinkValueCaller(value LinkValue, mapper LocalLinkValueMapper) (LinkValue, bool, error) {
	return mapper(value)
}

////////////////////////////////////////////////
// TECHNICAL IMPLEMENTATION OF LINKS OPERANDS //
////////////////////////////////////////////////

func (o LinkObject) UniqueId() string {
	return o.id
}

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
	return o.content, nil
}

func (o LinkObject) AsTrait() (Trait, error) {
	var empty Trait
	return empty, errors.New("invalid value: expecting trait, got object")
}

func (o LinkObject) AsVariable() (LinkVariable, error) {
	var empty LinkVariable
	return empty, errors.New("invalid value: expecting variable, got object")
}

func (g LinkGroup) UniqueId() string {
	return g.id
}

func (g LinkGroup) GetType() LinkValueType {
	return LinkValueAsGroup
}

func (g LinkGroup) AsLink() (Link, error) {
	return Link{}, errors.New("invalid value: expecting link, got group")
}

func (g LinkGroup) AsGroup() ([]Object, error) {
	return g.content, nil
}

func (g LinkGroup) AsObject() (Object, error) {
	var object Object
	return object, errors.New("invalid value: expecting object, got group")
}

func (g LinkGroup) AsTrait() (Trait, error) {
	var empty Trait
	return empty, errors.New("invalid value: expecting trait, got group")
}

func (g LinkGroup) AsVariable() (LinkVariable, error) {
	var empty LinkVariable
	return empty, errors.New("invalid value: expecting variable, got group")
}

func (l Link) UniqueId() string {
	return l.id
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

func (l Link) AsVariable() (LinkVariable, error) {
	var empty LinkVariable
	return empty, errors.New("invalid value: expecting variable, got link")
}

func (t LinkTrait) UniqueId() string {
	return t.id
}

func (t LinkTrait) GetType() LinkValueType {
	return LinkValueAsTrait
}

func (t LinkTrait) AsLink() (Link, error) {
	var empty Link
	return empty, errors.New("invalid value: expecting link, got trait")
}

func (t LinkTrait) AsGroup() ([]Object, error) {
	return nil, errors.New("invalid value: expecting group, got trait")
}

func (t LinkTrait) AsObject() (Object, error) {
	var object Object
	return object, errors.New("invalid value: expecting object, got trait")
}

func (t LinkTrait) AsTrait() (Trait, error) {
	return t.content, nil
}

func (t LinkTrait) AsVariable() (LinkVariable, error) {
	var empty LinkVariable
	return empty, errors.New("invalid value: expecting variable, got trait")
}

func (v LinkVariable) UniqueId() string {
	return v.id
}

func (v LinkVariable) GetType() LinkValueType {
	return LinkValueAsVariable
}

func (v LinkVariable) AsLink() (Link, error) {
	return Link{}, errors.New("invalid value: expecting link, got variable")
}

func (v LinkVariable) AsGroup() ([]Object, error) {
	return nil, errors.New("invalid value: expecting group, got variable")
}

func (v LinkVariable) AsObject() (Object, error) {
	var object Object
	return object, errors.New("invalid value: expecting object, got variable")
}

func (v LinkVariable) AsTrait() (Trait, error) {
	var empty Trait
	return empty, errors.New("invalid value: expecting trait, got variable")
}

func (v LinkVariable) AsVariable() (LinkVariable, error) {
	return v, nil
}
