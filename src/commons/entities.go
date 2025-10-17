package commons

import (
	"errors"
	"iter"
	"slices"
)

// ModelEntity is either an object, or a group of entities.
// An entity does not exist per se, it is just a modeling option.
// It implements composite: either it is a group of entities (non leaf), or an object (leaf).
type ModelEntity interface {
	// An entiy is an element of a model
	Modelable
	// We may link entities together.
	// For instance, (Jean and Marie) went to this place.
	Linkable
	// IsGroup returns true if object is composite
	IsGroup() bool
	// AsGroup returns the object as a group, an error is is is not
	AsGroup() (ModelGroup, error)
	// AsObject returns the entity as an object
	AsObject() (ModelObject, error)
}

// ModelGroup is a group of objects.
// Groups are just a declaration of an anonymous set that have no special meaning.
// If you want to add meaning in that group, for instance a couple, then create a link.
// It is Modelable, then, but not a model component.
// Reason is a group is a model element, it does not exist per se as a concrete object.
// It makes no sense to consider a group that would be temporal.

type ModelGroup interface {
	// A group is an entity too by definition
	ModelEntity
	// Elements return the elements of the group as a sequence to iterate over
	Elements() iter.Seq[ModelEntity]
	// Content returns the group by values
	Content() []ModelEntity
}

// simpleGroup enriches a slice of model entities
type simpleGroup []ModelEntity

// GetType returns TypeGroup by definition
func (s simpleGroup) GetType() ModelableType {
	return TypeGroup
}

// Elements decorates the slice as an iterator (using slices.Values)
func (s simpleGroup) Elements() iter.Seq[ModelEntity] {
	return slices.Values(s)
}

// Content is exactly s, but to avoid any side effect, we copy
func (s simpleGroup) Content() []ModelEntity {
	result := make([]ModelEntity, len(s))
	copy(result, s)
	return result
}

// IsGroup returns true
func (s simpleGroup) IsGroup() bool {
	return true
}

// AsGroup returns the object as a group, an error is is is not
func (s simpleGroup) AsGroup() (ModelGroup, error) {
	return s, nil
}

// AsObject raises an error
func (s simpleGroup) AsObject() (ModelObject, error) {
	return nil, errors.New("expected object, got group")
}

// NewModelGroup builds a new group decorating values
func NewModelGroup(values []ModelEntity) ModelGroup {
	result := make(simpleGroup, len(values))
	copy(result, values)
	return simpleGroup(values)
}

// ModelObject is the component that runs in the structure.
// An entiy defines an objet or a group, but the actual component is an object.
type ModelObject interface {
	// An object is well defined
	Identifiable
	// An object is a component of a model
	ModelComponent
}

// TemporalObject is an object that is active during a given period.
// For instance, an human activity is called life.
type TemporalObject interface {
	// A temporal object is an object
	ModelObject
	// By definition, a temporal object is a Temporal
	Temporal
}

// objectTemporalDecorator adds a period as a lifetime to an object
type objectTemporalDecorator struct {
	// period of activity for that object
	period Period
	// value is the decorated object
	value ModelObject
}

// ActivePeriod returns the object activity
func (o *objectTemporalDecorator) ActivePeriod() Period {
	// will panic for nil and that is OK
	return o.period
}

// SetActivePeriod changes period for that object
func (o *objectTemporalDecorator) SetActivePeriod(p Period) {
	if o != nil {
		o.period = p
	}
}

// Id returns empty for nil, the object id otherwise
func (o *objectTemporalDecorator) Id() string {
	if o == nil {
		return ""
	} else {
		return o.Id()
	}
}

// GetType returns TypeObject if o is not nil, otherwise it is unmanaged.
// Temporal object with no period makes no sense, so prefer unmanaged
func (o *objectTemporalDecorator) GetType() ModelableType {
	if o == nil {
		return TypeUnmanaged
	}

	return o.GetType()
}

// NewTemporalObject returns an object active during that period
func NewTemporalObject(period Period, object ModelObject) TemporalObject {
	if object == nil {
		return nil
	}

	result := new(objectTemporalDecorator)
	result.period = period
	result.value = object
	return result
}
