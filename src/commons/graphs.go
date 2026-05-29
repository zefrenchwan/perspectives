package commons

import (
	"slices"
	"sync"
)

// Linkable represents an entity that can be linked in a graph.
// A link, especially, is linkable.
// It means that it is possible to use links as operands
type Linkable interface {
	Element
}

// linkOperand represents a collection of linkable values that can be sorted or not.
type linkOperand struct {
	values []Linkable // values as the raw values, duplicates are possible
	sorted bool       // sorted is true for lists, false for sets
}

// addValue adds a new value to the operand collection, updating the sorted flag if necessary.
// It does nothing if value is nil, not even changing the sorted value.
func (op *linkOperand) addValue(value Linkable, sorted bool) {
	if value == nil {
		return
	}

	op.values = append(op.values, value)
	op.sorted = sorted
}

// copyValues returns a copy of the operand values (not values, just the slice)
func (op *linkOperand) copyValues() []Linkable {
	return slices.Clone(op.values)
}

// remove removes all values from the operand collection that satisfy the given predicate.
func (op *linkOperand) remove(predicate func(Linkable) bool) {
	result := make([]Linkable, 0)
	for _, value := range op.values {
		if !predicate(value) {
			result = append(result, value)
		}
	}
	op.values = result
}

// size returns the number of values in the operand collection.
func (op *linkOperand) size() int {
	if op == nil {
		return 0
	}

	return len(op.values)
}

// newLinkOperand creates a new link operand with an empty slice and sorted set to false.
func newLinkOperand() *linkOperand {
	return &linkOperand{
		values: make([]Linkable, 0),
		sorted: false,
	}
}

// Link represents a link between nodes in a graph.
type Link struct {
	locks    sync.RWMutex            // locks manage the locks for goroutines
	name     string                  // name of the link
	operands map[string]*linkOperand // operands of the link
	validity Period                  // moments the link is true
}

// NewLink creates a new empty link with provided name.
func NewLink(name string) *Link {
	return &Link{
		name:     name,
		operands: make(map[string]*linkOperand),
		validity: NewFullPeriod(),
	}
}

// Name returns the name of the link.
func (l *Link) Name() string {
	if l == nil {
		return ""
	}

	return l.name
}

// Same returns true if the link is the same as the other linkable.
func (l *Link) Same(other Element) bool {
	if l == nil && other == nil {
		return true
	} else if l == nil || other == nil {
		return false
	} else if !IsElementDeclaredInstance(other, l.DeclaringClass()) {
		return false
	}

	otherLink, ok := other.(*Link)
	if !ok {
		return false
	}
	// at this point, both are links.

	l.locks.RLock()
	otherLink.locks.RLock()
	defer l.locks.RUnlock()
	defer otherLink.locks.RUnlock()

	current := l
	matchingCurrent := otherLink

	if l.name != otherLink.name {
		return false
	} else if len(current.operands) != len(matchingCurrent.operands) {
		return false
	}

	// todo : end it
	return true

}

// DeclaringClass returns the class that declares the link, obviously including CLASS_LINK itself
func (l *Link) DeclaringClass() Class {
	return CLASS_LINK
}

// Validity returns the period the link is active
func (l *Link) Validity() Period {
	return l.validity.Copy()
}

// SetValidity sets the period the link is active
func (l *Link) SetValidity(p Period) {
	if l == nil {
		return
	}

	l.locks.Lock()
	defer l.locks.Unlock()

	l.validity = p
}

// Roles return the roles of the link.
// Result is sorted alphabetically
// For instance, Likes(subject=[Jean],Object=[Pizza]).Roles() will return ["object", "subject"]
func (l *Link) Roles() []string {
	if l == nil {
		return nil
	}

	l.locks.RLock()
	defer l.locks.RUnlock()

	size := len(l.operands)
	result := make([]string, size)
	index := 0
	for role := range l.operands {
		result[index] = role
		index++
	}

	slices.Sort(result)
	return result
}

// Has returns true and related value if the link has an operand with the given name.
// Otherwise, it returns nil, false.
func (l *Link) Has(name string) ([]Linkable, bool) {
	if l == nil {
		return nil, false
	}

	l.locks.RLock()
	defer l.locks.RUnlock()

	operands, ok := l.operands[name]
	return operands.copyValues(), ok
}

// Add adds a new operand to the link with the given name if order does not matter.
// For instance, Likes(subject=[Jean],Object=[Pizza]).Add("Object", "Tiramisu")
// will make Likes(subject=[Jean],Object=[Pizza, Tiramisu])
func (l *Link) Add(name string, operand Linkable) {
	l.addOperand(name, operand, true)
}

// Append adds a new operand to the link with the given name if order DOES matter.
// For instance, presidents of a country in order
func (l *Link) Append(name string, operand Linkable) {
	l.addOperand(name, operand, false)
}

// addOperand adds a new operand to the link with the given name and sorting preference.
// if value is nil, it does nothing, not even changing the sorted value.
func (l *Link) addOperand(name string, operand Linkable, sorted bool) {
	if l == nil || operand == nil {
		return
	}

	l.locks.Lock()
	defer l.locks.Unlock()

	if previous, ok := l.operands[name]; ok && previous != nil {
		previous.addValue(operand, sorted)
		l.operands[name] = previous
	} else {
		current := newLinkOperand()
		current.addValue(operand, sorted)
		l.operands[name] = current
	}
}

// Remove for a given role, specific values by predicate.
// If the predicate matches all values, the role is removed from the link.
func (l *Link) Remove(name string, op func(linkable Linkable) bool) {
	if l == nil || op == nil {
		return
	}

	l.locks.Lock()
	defer l.locks.Unlock()

	if operand, ok := l.operands[name]; ok {
		operand.remove(op)
		if operand.size() != 0 {
			l.operands[name] = operand
		} else {
			delete(l.operands, name)
		}
	}
}

// Operands return a copy of the role map
func (l *Link) Operands() map[string][]Linkable {
	if l == nil {
		return nil
	}

	l.locks.RLock()
	defer l.locks.RUnlock()

	result := make(map[string][]Linkable)
	for k, v := range l.operands {
		result[k] = v.copyValues()
	}

	return result
}
