package commons

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TemporalValues interface {
	At(t time.Time) (any, bool)
	Range(yield func(p Period, v any) bool)
	IsEmpty() bool
}

type Instance interface {
	Element
	Lifetime() Period
	SetLifetime(Period)

	Description() map[string]string

	SetAttribute(name string, p Period, value any) error
	Attribute(name string) TemporalValues
}

type periodValue struct {
	value    any
	validity Period
}

type periodValues struct {
	elements     []periodValue
	declaredType reflect.Type
}

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

func (p *periodValues) wouldAccept(value any) bool {
	if p.declaredType == nil {
		return true
	}

	return p.declaredType == reflect.TypeOf(value)
}

func (p *periodValues) Add(period Period, e any) error {
	if p == nil || period.IsEmpty() {
		return nil
	} else if len(p.elements) == 0 {
		p.elements = []periodValue{{validity: period, value: e}}
		p.declaredType = reflect.TypeOf(e)
		return nil
	}

	incomingType := reflect.TypeOf(e)
	if p.declaredType != incomingType {
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

func (p *periodValues) IsEmpty() bool {
	return p == nil || len(p.elements) == 0
}

func (p *periodValues) copy() *periodValues {
	result := new(periodValues)
	result.elements = make([]periodValue, len(p.elements))
	copy(result.elements, p.elements)
	result.declaredType = p.declaredType
	return result
}

func newPeriodValues() *periodValues {
	result := new(periodValues)
	result.elements = make([]periodValue, 0)
	return result
}

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

type temporalInstance struct {
	id         string
	locks      sync.RWMutex
	attributes map[string]*periodValues
	lifetime   Period
}

func (t *temporalInstance) Id() string {
	return t.id
}

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

func (t *temporalInstance) DeclaringClasses() []Class {
	return []Class{CLASS_INSTANCE}
}

func (t *temporalInstance) Lifetime() Period {
	return t.lifetime
}

func (t *temporalInstance) SetLifetime(p Period) {
	if t == nil {
		return
	}

	t.locks.Lock()
	defer t.locks.Unlock()
	t.lifetime = p
}

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

func NewTemporalInstance() *temporalInstance {
	return &temporalInstance{
		id:         uuid.NewString(),
		attributes: make(map[string]*periodValues),
		lifetime:   NewFullPeriod(),
	}
}
