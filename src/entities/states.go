package entities

import (
	"iter"
	"slices"
	"strconv"
	"strings"

	"github.com/zefrenchwan/perspectives.git/commons"
	"github.com/zefrenchwan/perspectives.git/periods"
	"github.com/zefrenchwan/perspectives.git/values"
)

// State is the immutable description of an entity at a given time.
type State interface {
	// Identifiable to define a unique identifier for the state.
	commons.Identifiable
	// Hashable to get the hash of the state.
	// States are immutables, so hash sums up the current state.
	commons.Hashable
	// TimeBounded to define a time period during which the entity exists.
	// It may vary, because, for instance, X is alive so far, until death (and then end of period).
	periods.TimeBounded
	// Attributes describe the state of an element.
	// Keys are names, and values are attributes (basically a map[period]primitives)
	Attributes() iter.Seq2[string, values.ImmutableValuesMapping[values.PrimitiveValue]]
	// Roles describe the relationships between elements.
	// Keys are names, and values are roles (basically a map[period]references)
	Roles() iter.Seq2[string, values.ImmutableValuesMapping[values.ReferenceValue]]
}

// localState is the in memory implementation of states.
// It contains the most basic implementation based on maps for roles and attributes.
type localState struct {
	// id of the state, it should be unique.
	id string
	// activity period of the state.
	activity periods.Period
	// attributes as a map of names and related values as a mapping.
	// Values are primitives.
	attributes map[string]values.ImmutableValuesMapping[values.PrimitiveValue]
	// roles as a map of names and related values as a mapping.
	// Values are references to other entities.
	roles map[string]values.ImmutableValuesMapping[values.ReferenceValue]
	// hashString is the hash of the state, calculated once
	hashString string
}

// ToHashString returns the hash string of the state.
// It is constant because it is calculated once and does not change.
func (l localState) ToHashString() string {
	return l.hashString
}

// Id of current state
func (l localState) Id() string {
	return l.id
}

// Activity of the state
func (l localState) Activity() periods.Period {
	return l.activity
}

// Attributes of the state as an iterator (to avoid defensive copies)
func (l localState) Attributes() iter.Seq2[string, values.ImmutableValuesMapping[values.PrimitiveValue]] {
	return func(yield func(string, values.ImmutableValuesMapping[values.PrimitiveValue]) bool) {
		for attr, mapper := range l.attributes {
			if !yield(attr, mapper) {
				return
			}
		}

		return
	}
}

// Roles of the state as an iterator (to avoid defensive copies)
func (l localState) Roles() iter.Seq2[string, values.ImmutableValuesMapping[values.ReferenceValue]] {
	return func(yield func(string, values.ImmutableValuesMapping[values.ReferenceValue]) bool) {
		for role, mapper := range l.roles {
			if !yield(role, mapper) {
				return
			}
		}

		return
	}
}

// localStateHash returns a hash of the local state (long, should be done once).
// REMEMBER TO SET THE HASH FOR YOUR INSTANCE !
func localStateHash(l *localState) string {
	var base strings.Builder
	base.WriteString("LocalState id =")
	base.WriteString(l.id)
	base.WriteString("\n")
	base.WriteString("activity =")
	base.WriteString(l.activity.AsRawString())
	base.WriteString("\n")

	rolesSize := len(l.roles)
	roleValues := make([]string, rolesSize)
	index := 0
	base.WriteString("Roles : ")
	base.WriteString(strconv.Itoa(rolesSize))
	base.WriteString("\n")
	var rolesString strings.Builder
	for role, mapper := range l.roles {
		rolesString.WriteString(strconv.Itoa(len(role)))
		rolesString.WriteString(":")
		rolesString.WriteString(role)
		rolesString.WriteString("=>")
		rolesString.WriteString(mapper.ToHashString())
		roleValues[index] = rolesString.String()
		index++
		rolesString.Reset()
	}

	slices.Sort(roleValues)
	base.WriteString("roles =")
	base.WriteString(strings.Join(roleValues, ","))
	base.WriteString("\n")

	// same for attributes
	attrSize := len(l.attributes)
	attrValues := make([]string, attrSize)
	index = 0
	base.WriteString("Attributes : ")
	base.WriteString(strconv.Itoa(attrSize))
	base.WriteString("\n")
	var attrString strings.Builder
	for attr, mapper := range l.attributes {
		attrString.WriteString(strconv.Itoa(len(attr)))
		attrString.WriteString(":")
		attrString.WriteString(attr)
		attrString.WriteString("->")
		attrString.WriteString(mapper.ToHashString())
		attrValues[index] = attrString.String()
		index++
		attrString.Reset()
	}

	slices.Sort(attrValues)
	base.WriteString("attributes =")
	base.WriteString(strings.Join(attrValues, ","))
	base.WriteString("\n")

	return commons.HashString(base.String())
}

// NewLocalState creates a new state in memory.
// It contains an id, the activity period, attributes and roles.
// Attributes and roles are stored as a map, but mappings are not necessarily the in-memory implementation.
// This function is used to create a new state with a few attributes and roles.
func NewLocalState(
	id string, // id of the state
	activity periods.Period, // activity period : when is the state valid
	attributes map[string]values.ImmutableValuesMapping[values.PrimitiveValue], // name of attributes linked to immutable values
	roles map[string]values.ImmutableValuesMapping[values.ReferenceValue], // name of roles linked to immutable references
) State {
	result := localState{
		id:         id,
		activity:   activity,
		attributes: attributes,
		roles:      roles,
		hashString: "",
	}

	// hash calculation once content is set
	result.hashString = localStateHash(&result)

	return result
}
