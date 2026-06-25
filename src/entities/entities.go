package entities

import (
	"time"

	"github.com/zefrenchwan/perspectives.git/periods"
)

// Entity is the unique concept of the system. It is the base of all the other concepts.
// It is ABSOLUTELY MANDATORY to make it as an immutable object.
type Entity interface {
	// Id of the entity, should be unique and immutable.
	Id() string
	// Same returns true if the entity is the same as the other entity.
	// It means that the entity has the same id as the other entity, same content too.
	Same(other Entity) bool
	// Activity is the lifetime of the entity.
	Activity() periods.Period
	// Attributes allows iteration over the attributes of the entity, by name.
	// It does not return the values of the attributes in a slice, to avoid multiple allocations.
	Attributes(yield func(attribute string) bool)
	// Attribute returns the details of the attribute with the given name.
	Attribute(attribute string) (AttributeDetails, bool)
	// Values allows an iteration over the attributes (by name) and values (for that name)
	Values(yield func(attribute string, value DynamicValues) bool)
	// Value returns the values of the attribute with the given name (if it exists).
	Value(attribute string) (DynamicValues, bool)
	// ValueAt returns the value of the attribute with the given name at the given moment (if it exists).
	ValueAt(attribute string, moment time.Time) (any, bool)
	// ValuesAt returns the values of the entity at the given moment.
	// It returns a map of attribute names to their values (as any, in reality a primitive type as defined earlier).
	ValuesAt(moment time.Time) (map[string]any, bool)
	// Roles allows iteration over the roles of the entity.
	// Same as attributes, we don't return the values of the roles in a slice, to avoid multiple allocations.
	Roles(yield func(role string) bool)
	// Role returns the entity associated with the given role (if it exists).
	Role(string) (Entity, bool)
	// Links allows iteration over the links of the entity as name and matching entity.
	Links(func(string, Entity) bool)
	// ToHashString returns the hash of the entity.
	// Because an entity is immutable, the hash string should be invariant.
	ToHashString() string
}

// localEntity is an entity that is stored locally.
// Its implementation is the default one, to prove the concept.
// Production ready code should manage long history.
type localEntity struct {
	// id of the entity
	id string
	// activity is the period during which the entity is valid.
	// In a nutshell, its lifetime.
	activity periods.Period
	// values are the temporal values associated with their attributes names
	values map[string]*valuesHandler
	// roles link this entity with another
	roles map[string]Entity
	// hashString is the hash of the entity
	hashString string
}

// Id returns the id of the entity
func (b *localEntity) Id() string {
	if b == nil {
		return ""
	}

	return b.id
}

// ToHashString returns the hash of the entity
func (b *localEntity) ToHashString() string {
	if b == nil {
		return ""
	}
	return b.hashString
}

// Same returns true if the entity is the same as the other entity : same class, same id, same period, same values
func (b *localEntity) Same(other Entity) bool {
	if b == nil && other == nil {
		return true
	} else if b == nil || other == nil {
		return false
	} else if b.Id() != other.Id() {
		return false
	}

	return b.ToHashString() == other.ToHashString()
}

// Activity returns the period during which the entity is valid
func (b *localEntity) Activity() periods.Period {
	if b == nil {
		return periods.NewEmptyPeriod()
	}
	return b.activity
}

// Attributes return the name of all attributes to iterate over
func (b *localEntity) Attributes(yield func(attribute string) bool) {
	for attribute := range b.values {
		if !yield(attribute) {
			return
		}
	}
}

// Attribute returns the attribute details by name
func (b *localEntity) Attribute(attribute string) (AttributeDetails, bool) {
	content, found := b.values[attribute]
	if !found {
		return AttributeDetails{}, false
	}

	return AttributeDetails{
		AttributeName:     attribute,
		AttributeType:     content.DataType(),
		AttributeValidity: content.Validity(),
		InstanceActivity:  b.activity,
	}, true
}

// Values iterates over the attributes values by name.
func (b *localEntity) Values(yield func(attributeName string, attributeValues DynamicValues) bool) {
	for attribute, content := range b.values {
		if !yield(attribute, content) {
			return
		}
	}
}

// Value returns the temporal values associated with the given attribute name, if it exists
func (b *localEntity) Value(attribute string) (DynamicValues, bool) {
	value, found := b.values[attribute]
	return value, found
}

// ValuesAt returns the content at a given time, as a map of attributes and values.
// If entity is not active at that moment, then it returns nil, false.
func (b *localEntity) ValuesAt(moment time.Time) (map[string]any, bool) {
	if b == nil {
		return nil, false
	} else if !b.activity.Contains(moment) {
		return nil, false
	}

	result := make(map[string]any)
	for attribute, content := range b.values {
		if value, exists := content.At(moment); exists {
			result[attribute] = value
		}
	}

	return result, true
}

// ValueAt returns the value of the attribute with the given name at the given moment (if it exists).
func (b *localEntity) ValueAt(attribute string, moment time.Time) (any, bool) {
	if b == nil {
		return nil, false
	}

	value, found := b.values[attribute]
	if !found {
		return nil, false
	}

	return value.At(moment)
}

// Roles allows an iteration over the name of each role
func (l *localEntity) Roles(yield func(string) bool) {
	if l == nil {
		return
	}

	for role := range l.roles {
		if !yield(role) {
			return
		}
	}
}

// Role returns, if any, the entity associated to the role
func (l *localEntity) Role(role string) (Entity, bool) {
	linkable, found := l.roles[role]
	return linkable, found
}

// Links allows an iteration over each linked entity as the name of each role and the associated entity
func (l *localEntity) Links(yield func(string, Entity) bool) {
	for role, linkable := range l.roles {
		if !yield(role, linkable) {
			return
		}
	}
}

// newLocalEntity returns an empty localEntity.
// It does NOT set the hash, pay attention
func newLocalEntity(id string) *localEntity {
	result := &localEntity{
		id:       id,
		activity: periods.NewEmptyPeriod(),
		values:   make(map[string]*valuesHandler),
		roles:    make(map[string]Entity),
	}

	// REMEMBER : AS AN INTERNAL CODE, IT DOES NOT SET THE HASH
	// result.hashString = hashEntity(result)
	return result
}

// localEntityLoad just copies values from other
func localEntityLoad(other Entity) *localEntity {
	result := new(localEntity)
	result.id = other.Id()
	result.activity = other.Activity()
	// manage attributes
	result.values = make(map[string]*valuesHandler)
	for attribute, content := range other.Values {
		handler := new(valuesHandler)
		handler.storedType = content.DataType()
		handler.equals = primitiveTypeEqualsFunc(handler.storedType)
		for period, value := range content.Range {
			handler.values = append(handler.values, valueNode{matchingPeriod: period, value: value})
		}

		result.values[attribute] = handler
	}

	// manage links
	newRoles := make(map[string]Entity)
	for role, linkable := range other.Links {
		newRoles[role] = linkable
	}

	result.roles = newRoles

	result.hashString = hashEntity(result)
	return result
}
