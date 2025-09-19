package models

import (
	"errors"
	"fmt"
	"slices"

	"github.com/zefrenchwan/perspectives.git/structures"
)

// Variable defines a variable that may be replaced by any other link value
type Variable struct {
	// Name of the variable (usually "x","y","z")
	name string
	// ValidTypes are the union of the valid types to replace this variable for
	validTypes []EntityType
	// ValidTraits contain the union of acceptable traits
	validTraits []Trait
}

// Name returns the name of the variable
func (lv Variable) Name() string {
	return lv.name
}

// GetType returns EntityTypeVariable
func (lv Variable) GetType() EntityType {
	return EntityTypeVariable
}

// AsLink raises an error
func (lv Variable) AsLink() (*Link, error) {
	return nil, errors.ErrUnsupported
}

// AsGroup raises an error
func (lv Variable) AsGroup() ([]Object, error) {
	return nil, errors.ErrUnsupported
}

// AsObject raises an error
func (lv Variable) AsObject() (*Object, error) {
	return nil, errors.ErrUnsupported
}

// AsTrait raises an error
func (lv Variable) AsTrait() (Trait, error) {
	return Trait{}, errors.ErrUnsupported
}

// AsVariable returns the value as a variable
func (lv Variable) AsVariable() (Variable, error) {
	return lv, nil
}

// NewVariableForObject returns a new variable for that object
func NewVariableForObject(name string, traits []string) Variable {
	var matches []Trait
	for _, trait := range structures.SliceReduce(traits) {
		matches = append(matches, NewTrait(trait))
	}

	return Variable{
		name:        name,
		validTypes:  []EntityType{EntityTypeObject},
		validTraits: matches,
	}
}

// NewVariableForGroup returns the variable for a group of objects matching traits
func NewVariableForGroup(name string, traits []string) Variable {
	var matches []Trait
	for _, trait := range structures.SliceReduce(traits) {
		matches = append(matches, NewTrait(trait))
	}

	return Variable{
		name:        name,
		validTypes:  []EntityType{EntityTypeGroup},
		validTraits: matches,
	}
}

// NewVariableForTrait returns a new variable for that trait
func NewVariableForTrait(name string) Variable {
	return Variable{
		name:        name,
		validTypes:  []EntityType{EntityTypeTrait},
		validTraits: nil,
	}
}

// NewVariableForSpecificTraits returns a variable for trait that may take only specific values
func NewVariableForSpecificTraits(name string, traits []string) Variable {
	var matches []Trait
	for _, trait := range structures.SliceReduce(traits) {
		matches = append(matches, NewTrait(trait))
	}

	return Variable{
		name:        name,
		validTypes:  []EntityType{EntityTypeTrait},
		validTraits: matches,
	}
}

// NewVariableForLink returns a new variable for that link
func NewVariableForLink(name string) Variable {
	return Variable{
		name:        name,
		validTypes:  []EntityType{EntityTypeLink},
		validTraits: nil,
	}
}

// MatchesTraits returns true if traits match accepted traits for that variable
func (lv Variable) MatchesTraits(traits []Trait) bool {
	if len(lv.validTraits) == 0 {
		// no prerequisite
		return true
	} else {
		// prerequisites are set, and then should match
		return structures.SliceCommonElementFunc(lv.validTraits, traits, func(a, b Trait) bool { return a.Equals(b) })
	}
}

// MapAs transforms a variable to a value.
// Accepted values are slices of objects, objects, traits, links and related link values
func (lv Variable) MapAs(other any) (ModelEntity, error) {
	expectedTypes := lv.validTypes

	if v, ok := other.(Object); ok {
		if !slices.Contains(expectedTypes, EntityTypeObject) {
			return nil, errors.New("object does not match expected type")
		}

		// test if object matches the definition.
		// Accept if there is a matching trait
		if lv.MatchesTraits(v.traits) {
			return &v, nil
		} else {
			return nil, errors.New("no matching trait compatible with type definition")
		}
	} else if v, ok := other.([]Object); ok {
		if !slices.Contains(expectedTypes, EntityTypeGroup) {
			return nil, errors.New("group does not match expected type")
		}

		// test if each object within the group matches the trait condition
		for index, obj := range v {
			if !lv.MatchesTraits(obj.traits) {
				return nil, fmt.Errorf("value at index %d does not match traits condition", index)
			}
		}

		return objectsGroup(v), nil
	} else if v, ok := other.(Trait); ok {
		if !slices.Contains(expectedTypes, EntityTypeTrait) {
			return nil, errors.New("group does not match expected type")
		}

		// for traits, either variable does not specify any, or traits match
		if len(lv.validTraits) == 0 {
			return v, nil
		} else if !slices.ContainsFunc(lv.validTraits, func(t Trait) bool { return v.Equals(t) }) {
			return nil, errors.New("trait value does not match expected traits")
		} else {
			return &v, nil
		}
	} else if v, ok := other.(Link); ok {
		if !slices.Contains(expectedTypes, EntityTypeLink) {
			return nil, errors.New("group does not match expected type")
		} else {
			return &v, nil
		}
	} else {
		return nil, errors.New("invalid value to map")
	}
}
