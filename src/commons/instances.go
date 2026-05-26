package commons

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"
)

// TemporalValues is the general contract for a time changing value :
// get value at a given time
// iterator over periods and values
// test if there is a value
type TemporalValues interface {
	At(t time.Time) (any, bool)             // At returns the value at a given time, if any, and a boolean to indicate if present
	Range(yield func(p Period, v any) bool) // Range iterates over periods and values during that period
	IsEmpty() bool                          // IsEmpty returns true if there are no values (avoids range iteration)
}

// Instance is the general contract for an entity with a lifetime and attributes.
// It has a lifetime and attributes.
// Lifetime defines the time span during which the instance exists.
// Attributes are key-value pairs that can change over time.
type Instance interface {
	Element // Element of the model

	Lifetime() Period   // Lifetime returns the time span during which the instance exists
	SetLifetime(Period) // SetLifetime changes the time span during which the instance exists

	Description() map[string]string // Description returns the structure of the instance as name of attribute and its type

	SetAttribute(name string, p Period, value any) error // SetAttribute changes value for that attribute (by name) during a period.
	Attribute(name string) TemporalValues                // Attribute returns the temporal values for that attribute
}

// periodValue is a container for a value and a period.
// Because a period contains slices, it was not possible to use map[Period]any
// so we deal with a slice of periodValue
type periodValue struct {
	value    any
	validity Period
}

// periodValues regroups all values for a given attribute.
// It contains a slice of periodValue and the declared type of the attribute.
// The declared type is used to check if the attribute is set with a compatible type.
// It should be consistent with the attribute values over time.
type periodValues struct {
	elements     []periodValue // elements are all the values over time
	declaredType reflect.Type  // declared type of the attribute, for each value
}

// At returns the value, if any, at that time
func (p *periodValues) At(t time.Time) (any, bool) {
	if p == nil {
		return nil, false
	}

	for _, element := range p.elements {
		if element.validity.Contains(t) {
			return element.value, true
		}
	}
	return nil, false
}

// Range allows iteration over all values as period and value.
// We may then perform : for p,v := range pv.Range()
func (p *periodValues) Range(yield func(p Period, v any) bool) {
	if p == nil {
		return
	}
	for _, element := range p.elements {
		if !yield(element.validity, element.value) {
			break
		}
	}
}

// wouldAccept checks if the given value is compatible with the declared type of the periodValues instance.
// if p is null, it raises a NPE on purpose
func (p *periodValues) wouldAccept(value any) bool {
	if p.declaredType == nil {
		return true
	}

	return p.declaredType == reflect.TypeOf(value)
}

// Add a value during that period to the periodValues instance.
// It will merge with existing values if they are compatible and recalculate matching periods for each previous element.
func (p *periodValues) Add(period Period, e any) error {
	if p == nil || period.IsEmpty() {
		return nil
	} else if len(p.elements) == 0 {
		p.elements = []periodValue{{validity: period, value: e}}
		p.declaredType = reflect.TypeOf(e)
		return nil
	}

	incomingType := reflect.TypeOf(e)
	if p.declaredType != nil && p.declaredType != incomingType {
		return fmt.Errorf("cannot add value with type %v to periodValues with declared type %v", incomingType, p.declaredType)
	}

	finalPeriod := period
	newElements := make([]periodValue, 0, len(p.elements)+1)
	for _, element := range p.elements {
		if reflect.DeepEqual(element.value, e) {
			finalPeriod = finalPeriod.Union(element.validity)
		} else {
			remaining := element.validity.Remove(period)
			if !remaining.IsEmpty() {
				newElements = append(newElements, periodValue{
					value:    element.value,
					validity: remaining,
				})
			}
		}
	}

	if finalPeriod.Equals(NewFullPeriod()) {
		p.elements = []periodValue{{validity: NewFullPeriod(), value: e}}
		return nil
	}

	newElements = append(newElements, periodValue{
		value:    e,
		validity: finalPeriod,
	})

	p.elements = newElements
	return nil
}

// IsEmpty checks if the periodValues has no element
func (p *periodValues) IsEmpty() bool {
	return p == nil || len(p.elements) == 0
}

// copy creates a copy of the periodValues instance : elements and declaredType
func (p *periodValues) copy() *periodValues {
	result := new(periodValues)
	result.elements = make([]periodValue, len(p.elements))
	copy(result.elements, p.elements)
	result.declaredType = p.declaredType
	return result
}

// newPeriodValues creates a new empty periodValues instance.
func newPeriodValues() *periodValues {
	result := new(periodValues)
	result.elements = make([]periodValue, 0)
	return result
}

// Initializes a new periodValues instance with a given period and value.
// Forces the declaring type too
func initPeriodValues(period Period, value any) *periodValues {
	result := newPeriodValues()
	if period.IsEmpty() {
		return result
	}

	result.elements = append(result.elements, periodValue{
		validity: period,
		value:    value,
	})

	result.declaredType = reflect.TypeOf(value)
	return result
}

// temporalInstance represents an object with a lifetime period, and attributes that may change over time.
type temporalInstance struct {
	id         string                   // unique identifier for the instance
	locks      sync.RWMutex             // synchronization mechanism for concurrent access
	attributes map[string]*periodValues // map of attribute names to their time-dependent values
	lifetime   Period                   // lifetime period during which the instance exists
}

// Id returns the unique identifier of the temporal instance
func (t *temporalInstance) Id() string {
	return t.id
}

// Same checks if two temporal instances are the same by comparing their IDs
func (t *temporalInstance) Same(other Element) bool {
	if t == nil && other == nil {
		return true
	} else if t == nil || other == nil {
		return false
	} else if !IsElementDeclaredInstance(other, CLASS_INSTANCE) {
		return false
	}

	return t.id == other.Id()
}

// DeclaringClasses declares that an instance's class is CLASS_INSTANCE
func (t *temporalInstance) DeclaringClasses() []Class {
	return []Class{CLASS_INSTANCE}
}

// Lifetime returns the time span during which the instance exists
func (t *temporalInstance) Lifetime() Period {
	return t.lifetime
}

// SetLifetime changes the time span during which the instance exists
func (t *temporalInstance) SetLifetime(p Period) {
	if t == nil {
		return
	}

	t.locks.Lock()
	defer t.locks.Unlock()
	t.lifetime = p
}

// Description returns the structure of the instance as name of attribute and its type.
// For instance, if an instance has an attribute "name" of type string, the description will contain "name": "string".
func (t *temporalInstance) Description() map[string]string {
	if t == nil {
		return nil
	}
	t.locks.RLock()
	defer t.locks.RUnlock()

	result := make(map[string]string)
	for name, attr := range t.attributes {
		result[name] = attr.declaredType.Name()
	}

	return result
}

// SetAttribute changes value for that attribute (by name) during a period.
// We may put any value, but once a type has been chosen, it must be consistent.
// If it is not, it returns an error.
func (t *temporalInstance) SetAttribute(name string, p Period, value any) error {
	if t == nil {
		return nil
	}
	if p.IsEmpty() {
		return nil
	}

	t.locks.Lock()
	defer t.locks.Unlock()

	if matchingAttribute, ok := t.attributes[name]; !ok {
		t.attributes[name] = initPeriodValues(p, value)
	} else if !matchingAttribute.wouldAccept(value) {
		return fmt.Errorf("attribute %s does not accept type %T", name, value)
	} else {
		matchingAttribute.Add(p, value)
		t.attributes[name] = matchingAttribute
	}
	return nil
}

// Attribute returns the values for that name during the full time.
// Note that a copy is returned to avoid concurrent modification
func (t *temporalInstance) Attribute(name string) TemporalValues {
	if t == nil {
		return nil
	}

	t.locks.RLock()
	defer t.locks.RUnlock()

	if res, ok := t.attributes[name]; !ok {
		return nil
	} else {
		return res.copy()
	}
}

// NewTemporalInstance creates a new temporal instance with default lifetime (full period).
// A *temporalInstance is an instance, so you may use it as one
func NewTemporalInstance() Instance {
	return &temporalInstance{
		id:         uuid.NewString(),
		attributes: make(map[string]*periodValues),
		lifetime:   NewFullPeriod(),
	}
}
